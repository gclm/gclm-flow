#!/usr/bin/env python3
import argparse
import json
import re
from collections import Counter, defaultdict
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path


WORKFLOW_TOPICS = {
    "skills": [r"\bskill\b", r"\bskills\b", r"references", r"SKILL\.md", r"updating-domain-skills"],
    "hooks": [r"\bhook\b", r"\bhooks\b"],
    "config": [r"config\.toml", r"\bconfig\b", r"provider", r"model", r"persistence", r"multi_agent", r"history\.jsonl"],
    "docs": [r"README", r"AGENTS\.md", r"\bdocs\b", r"文档"],
    "git": [r"\bgit\b", r"commit", r"push", r"branch", r"worktree"],
    "review": [r"\breview\b", r"代码审查", r"审计", r"verify"],
    "testing": [r"\btest\b", r"testing", r"pytest", r"vitest", r"go test", r"cargo test"],
    "agents": [r"\bagent\b", r"\bagents\b", r"agent team", r"agent teams", r"multi-agent", r"multi agent", r"team"],
}

DOMAIN_TOPICS = {
    "devops": [r"docker", r"k8s", r"kubernetes", r"terraform", r"ci/cd", r"github actions", r"deploy"],
    "frontend-stack": [r"react", r"vue", r"frontend", r"tsx", r"vitest", r"playwright"],
    "python-stack": [r"python", r"fastapi", r"flask", r"pytest"],
    "go-stack": [r"\bgo\b", r"gin", r"echo", r"go test"],
    "java-stack": [r"java", r"spring", r"quarkus", r"junit"],
    "rust-stack": [r"rust", r"axum", r"actix", r"cargo test"],
    "database": [r"sql", r"postgres", r"mysql", r"redis", r"mongo", r"migration", r"schema"],
}

IGNORE_REPEATED_PROMPTS = {
    "继续",
    "确认",
    "你好",
    "你是谁",
    "restart",
}


@dataclass
class Entry:
    session_id: str
    ts: int
    text: str


def load_entries(path: Path) -> list[Entry]:
    entries = []
    with path.open() as fh:
        for lineno, line in enumerate(fh, start=1):
            line = line.strip()
            if not line:
                continue
            obj = json.loads(line)
            entries.append(Entry(session_id=obj["session_id"], ts=int(obj["ts"]), text=str(obj["text"])))
    return entries


def normalize_text(text: str) -> str:
    text = re.sub(r"\s+", " ", text.strip())
    return text


def short(text: str, limit: int = 88) -> str:
    text = normalize_text(text)
    if len(text) <= limit:
        return text
    return text[: limit - 1] + "..."


def match_topics(text: str, mapping: dict[str, list[str]]) -> set[str]:
    found = set()
    lowered = text.lower()
    for topic, patterns in mapping.items():
        if any(re.search(pattern, lowered, flags=re.IGNORECASE) for pattern in patterns):
            found.add(topic)
    return found


def summarize(entries: list[Entry], top_n: int, recent_n: int) -> str:
    sessions: dict[str, list[Entry]] = defaultdict(list)
    exact_prompt_counter = Counter()
    workflow_counter = Counter()
    domain_counter = Counter()
    day_counter = Counter()

    for entry in entries:
        sessions[entry.session_id].append(entry)
        exact_prompt_counter[normalize_text(entry.text)] += 1
        workflow_counter.update(match_topics(entry.text, WORKFLOW_TOPICS))
        domain_counter.update(match_topics(entry.text, DOMAIN_TOPICS))
        day_counter[datetime.fromtimestamp(entry.ts, tz=timezone.utc).date().isoformat()] += 1

    session_stats = []
    for session_id, items in sessions.items():
        items.sort(key=lambda item: item.ts)
        session_stats.append(
            {
                "session_id": session_id,
                "messages": len(items),
                "start": items[0].ts,
                "end": items[-1].ts,
                "duration_minutes": round((items[-1].ts - items[0].ts) / 60, 1),
                "first_prompt": short(items[0].text, 70),
            }
        )
    session_stats.sort(key=lambda item: (-item["messages"], -item["end"]))

    repeated_prompts = [
        (prompt, count)
        for prompt, count in exact_prompt_counter.most_common()
        if count >= 2 and prompt and prompt not in IGNORE_REPEATED_PROMPTS
    ][:top_n]

    recent_entries = entries[-recent_n:]
    recent_workflow = Counter()
    recent_domain = Counter()
    for entry in recent_entries:
        recent_workflow.update(match_topics(entry.text, WORKFLOW_TOPICS))
        recent_domain.update(match_topics(entry.text, DOMAIN_TOPICS))

    lines = []
    lines.append("# Codex History Retrospective")
    lines.append("")
    lines.append("## Overview")
    lines.append(f"- Entries: {len(entries)}")
    lines.append(f"- Sessions: {len(sessions)}")
    if entries:
        start = datetime.fromtimestamp(entries[0].ts, tz=timezone.utc).isoformat()
        end = datetime.fromtimestamp(entries[-1].ts, tz=timezone.utc).isoformat()
        lines.append(f"- Range: {start} -> {end}")
    lines.append("")
    lines.append("## Top Sessions")
    for item in session_stats[:top_n]:
        lines.append(
            f"- {item['session_id']}: {item['messages']} messages, {item['duration_minutes']} min, first prompt: {item['first_prompt']}"
        )
    lines.append("")
    lines.append("## Repeated Prompts")
    if repeated_prompts:
        for prompt, count in repeated_prompts:
            lines.append(f"- x{count}: {short(prompt, 120)}")
    else:
        lines.append("- No repeated prompts above threshold")
    lines.append("")
    lines.append("## Workflow Topic Counts")
    for topic, count in workflow_counter.most_common():
        lines.append(f"- {topic}: {count}")
    lines.append("")
    lines.append("## Domain Topic Counts")
    if domain_counter:
        for topic, count in domain_counter.most_common():
            lines.append(f"- {topic}: {count}")
    else:
        lines.append("- No domain-heavy signals detected")
    lines.append("")
    lines.append(f"## Recent Focus (last {len(recent_entries)} entries)")
    if recent_workflow:
        lines.append("- Workflow: " + ", ".join(f"{topic}={count}" for topic, count in recent_workflow.most_common()))
    else:
        lines.append("- Workflow: none")
    if recent_domain:
        lines.append("- Domain: " + ", ".join(f"{topic}={count}" for topic, count in recent_domain.most_common()))
    else:
        lines.append("- Domain: none")
    lines.append("")
    lines.append("## Suggested Action Routing")
    workflow_actions = build_workflow_actions(workflow_counter, repeated_prompts)
    domain_actions = build_domain_actions(domain_counter)
    if workflow_actions:
        lines.append("### Workflow / Config")
        for action in workflow_actions:
            lines.append(f"- {action}")
    if domain_actions:
        lines.append("### Domain Skills")
        for action in domain_actions:
            lines.append(f"- {action}")
    if not workflow_actions and not domain_actions:
        lines.append("- No strong writeback candidates yet")
    return "\n".join(lines) + "\n"


def build_workflow_actions(workflow_counter: Counter, repeated_prompts: list[tuple[str, int]]) -> list[str]:
    actions = []
    if workflow_counter["skills"] >= 8:
        actions.append("`skills` 主题高度集中，优先沉淀到 workflow skills、README、hooks 或发布脚本，而不是直接写入 `updating-domain-skills`。")
    if workflow_counter["hooks"] >= 4 or workflow_counter["config"] >= 6:
        actions.append("hooks/config 反复出现，优先把高频判断做成轻量检查、默认配置或启动脚本。")
    if workflow_counter["review"] >= 4 and workflow_counter["testing"] >= 4:
        actions.append("review/testing 同时高频，说明质量闭环重要，优先维护全局 `code-review` 与 `testing` skill，而不是在各栈重复。")
    if repeated_prompts:
        actions.append("存在重复提问或重复需求，可优先新增脚本、skill 或维护约定，减少再次口述同样指令。")
    return actions


def build_domain_actions(domain_counter: Counter) -> list[str]:
    actions = []
    for topic, count in domain_counter.most_common():
        if count >= 4:
            actions.append(f"`{topic}` 出现 {count} 次，若这些结论已在真实任务中验证，可通过 `updating-domain-skills` 回写到对应领域 skill。")
    return actions


def main() -> int:
    parser = argparse.ArgumentParser(description="Analyze Codex history.jsonl and generate a retrospective report.")
    parser.add_argument("--input", default=str(Path.home() / ".codex" / "history.jsonl"))
    parser.add_argument("--top", type=int, default=8)
    parser.add_argument("--recent", type=int, default=40)
    parser.add_argument("--output", default="-")
    args = parser.parse_args()

    path = Path(args.input)
    entries = load_entries(path)
    report = summarize(entries, top_n=max(args.top, 1), recent_n=max(args.recent, 1))
    if args.output == "-":
        print(report, end="")
    else:
        Path(args.output).write_text(report)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
