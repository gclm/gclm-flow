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
    "devops": [r"docker", r"k8s", r"kubernetes", r"terraform", r"ci/cd", r"github actions", r"deploy", r"aliyun", r"cloudflare", r"vercel"],
    "frontend-stack": [r"react", r"vue", r"frontend", r"tsx", r"vitest", r"playwright"],
    "python-stack": [r"python", r"fastapi", r"flask", r"pytest"],
    "go-stack": [r"\bgo\b", r"gin", r"echo", r"go test"],
    "java-stack": [r"java", r"spring", r"quarkus", r"junit"],
    "rust-stack": [r"rust", r"axum", r"actix", r"cargo test"],
    "database": [r"sql", r"sqlite", r"postgres", r"mysql", r"redis", r"mongo", r"migration", r"schema", r"query"],
}

IGNORE_REPEATED_PROMPTS = {
    "继续",
    "确认",
    "你好",
    "你是谁",
    "restart",
}

DOMAIN_SAMPLE_NOISE_PATTERNS = [
    r"history\.jsonl",
    r"reviewing-codex-history",
    r"updating-domain-skills",
    r"topic counts",
    r"candidate clusters",
    r"output-template",
    r"agents/remember\.toml",
    r"entry-template",
]

DOMAIN_META_GOVERNANCE_PATTERNS = [
    r"\bskill\b",
    r"\bskills\b",
    r"SKILL\.md",
    r"references?",
    r"薄入口",
    r"模板式",
    r"模板",
    r"维护约定",
    r"全量 skills",
    r"references 内容风格",
    r"skill 结构",
    r"轻量 hook",
    r"reviewing-codex-history",
    r"updating-domain-skills",
    r"history\.jsonl",
]


@dataclass
class Entry:
    session_id: str
    ts: int
    text: str


def load_entries(path: Path) -> list[Entry]:
    entries = []
    with path.open() as fh:
        for line in fh:
            line = line.strip()
            if not line:
                continue
            obj = json.loads(line)
            entries.append(Entry(session_id=obj["session_id"], ts=int(obj["ts"]), text=str(obj["text"])))
    return entries


def normalize_text(text: str) -> str:
    return re.sub(r"\s+", " ", text.strip())


def short(text: str, limit: int = 120) -> str:
    text = normalize_text(text)
    if len(text) <= limit:
        return text
    return text[: limit - 1] + "..."


def match_topics(text: str, mapping: dict[str, list[str]]) -> set[str]:
    lowered = text.lower()
    found = set()
    for topic, patterns in mapping.items():
        if any(re.search(pattern, lowered, flags=re.IGNORECASE) for pattern in patterns):
            found.add(topic)
    return found


def session_summary(entries: list[Entry], top_n: int) -> list[dict[str, object]]:
    sessions: dict[str, list[Entry]] = defaultdict(list)
    for entry in entries:
        sessions[entry.session_id].append(entry)
    summary = []
    for session_id, items in sessions.items():
        items.sort(key=lambda item: item.ts)
        summary.append(
            {
                "session_id": session_id,
                "message_count": len(items),
                "start_ts": items[0].ts,
                "end_ts": items[-1].ts,
                "duration_minutes": round((items[-1].ts - items[0].ts) / 60, 1),
                "first_prompt": short(items[0].text, 90),
            }
        )
    summary.sort(key=lambda item: (-int(item["message_count"]), -int(item["end_ts"])))
    return summary[:top_n]


def repeated_prompts(entries: list[Entry], top_n: int) -> list[dict[str, object]]:
    counter = Counter(normalize_text(entry.text) for entry in entries)
    prompts = []
    for prompt, count in counter.most_common():
        if count < 2 or not prompt or prompt in IGNORE_REPEATED_PROMPTS:
            continue
        prompts.append({"text": short(prompt, 160), "count": count})
        if len(prompts) >= top_n:
            break
    return prompts


def is_domain_meta_discussion(text: str) -> bool:
    lowered = text.lower()
    if not any(topic in lowered for topic in DOMAIN_TOPICS):
        return False
    return any(re.search(pattern, lowered, flags=re.IGNORECASE) for pattern in DOMAIN_META_GOVERNANCE_PATTERNS)


def is_domain_sample_noise(text: str) -> bool:
    lowered = text.lower()
    if is_domain_meta_discussion(lowered):
        return True
    return any(re.search(pattern, lowered, flags=re.IGNORECASE) for pattern in DOMAIN_SAMPLE_NOISE_PATTERNS)


def topic_counts(
    entries: list[Entry],
    mapping: dict[str, list[str]],
    *,
    exclude_domain_meta: bool = False,
) -> Counter:
    counter = Counter()
    for entry in entries:
        if exclude_domain_meta and is_domain_sample_noise(entry.text):
            continue
        counter.update(match_topics(entry.text, mapping))
    return counter


def topic_samples(
    entries: list[Entry],
    mapping: dict[str, list[str]],
    sample_per_topic: int,
    *,
    topic_filter: str | None = None,
    exclude_domain_meta: bool = False,
) -> dict[str, list[dict[str, object]]]:
    topics = [topic_filter] if topic_filter else list(mapping)
    by_topic: dict[str, list[dict[str, object]]] = {}
    for topic in topics:
        if topic not in mapping:
            continue
        matched = [entry for entry in entries if topic in match_topics(entry.text, mapping)]
        if not matched:
            continue
        samples = []
        seen = set()
        for entry in sorted(matched, key=lambda item: item.ts, reverse=True):
            if exclude_domain_meta and is_domain_sample_noise(entry.text):
                continue
            snippet = short(entry.text, 180)
            if snippet in seen:
                continue
            seen.add(snippet)
            samples.append(
                {
                    "session_id": entry.session_id,
                    "ts": entry.ts,
                    "text": snippet,
                }
            )
            if len(samples) >= sample_per_topic:
                break
        by_topic[topic] = samples
    return by_topic


def recent_focus(entries: list[Entry], recent_n: int) -> dict[str, dict[str, int]]:
    recent_entries = entries[-recent_n:]
    workflow = topic_counts(recent_entries, WORKFLOW_TOPICS)
    domain = topic_counts(recent_entries, DOMAIN_TOPICS, exclude_domain_meta=True)
    return {
        "window_size": len(recent_entries),
        "workflow": dict(workflow.most_common()),
        "domain": dict(domain.most_common()),
    }


def resolve_topic_mapping(topic: str, topic_kind: str) -> tuple[str, dict[str, list[str]]]:
    if topic_kind == "workflow":
        if topic not in WORKFLOW_TOPICS:
            raise ValueError(f"Unknown workflow topic: {topic}")
        return topic_kind, WORKFLOW_TOPICS
    if topic_kind == "domain":
        if topic not in DOMAIN_TOPICS:
            raise ValueError(f"Unknown domain topic: {topic}")
        return topic_kind, DOMAIN_TOPICS
    if topic in DOMAIN_TOPICS:
        return "domain", DOMAIN_TOPICS
    if topic in WORKFLOW_TOPICS:
        return "workflow", WORKFLOW_TOPICS
    raise ValueError(f"Unknown topic: {topic}")


def focused_topic_report(entries: list[Entry], topic: str, topic_kind: str, sample_limit: int) -> dict[str, object]:
    resolved_kind, mapping = resolve_topic_mapping(topic, topic_kind)
    exclude_domain_meta = resolved_kind == "domain"
    counts = topic_counts(entries, mapping, exclude_domain_meta=exclude_domain_meta)
    return {
        "kind": resolved_kind,
        "topic": topic,
        "count": counts.get(topic, 0),
        "samples": topic_samples(
            entries,
            mapping,
            sample_per_topic=sample_limit,
            topic_filter=topic,
            exclude_domain_meta=exclude_domain_meta,
        ).get(topic, []),
    }


def build_report(
    entries: list[Entry],
    top_n: int,
    recent_n: int,
    sample_per_topic: int,
    *,
    topic: str | None = None,
    topic_kind: str = "auto",
    topic_samples_limit: int = 12,
) -> dict[str, object]:
    workflow = topic_counts(entries, WORKFLOW_TOPICS)
    domain = topic_counts(entries, DOMAIN_TOPICS, exclude_domain_meta=True)
    report = {
        "metadata": {
            "generated_at": datetime.now(timezone.utc).isoformat(),
            "source": "history.jsonl",
            "entries": len(entries),
            "sessions": len({entry.session_id for entry in entries}),
            "range": {
                "start_ts": entries[0].ts if entries else None,
                "end_ts": entries[-1].ts if entries else None,
                "start_iso": datetime.fromtimestamp(entries[0].ts, tz=timezone.utc).isoformat() if entries else None,
                "end_iso": datetime.fromtimestamp(entries[-1].ts, tz=timezone.utc).isoformat() if entries else None,
            },
        },
        "top_sessions": session_summary(entries, top_n=top_n),
        "repeated_prompts": repeated_prompts(entries, top_n=top_n),
        "workflow_topics": {
            "counts": dict(workflow.most_common()),
            "samples": topic_samples(entries, WORKFLOW_TOPICS, sample_per_topic=sample_per_topic),
        },
        "domain_topics": {
            "counts": dict(domain.most_common()),
            "samples": topic_samples(entries, DOMAIN_TOPICS, sample_per_topic=sample_per_topic, exclude_domain_meta=True),
        },
        "recent_focus": recent_focus(entries, recent_n=recent_n),
    }
    if topic:
        report["focused_topic"] = focused_topic_report(entries, topic, topic_kind, topic_samples_limit)
    return report


def render_markdown(report: dict[str, object]) -> str:
    meta = report["metadata"]
    lines = ["# Codex History Facts", "", "## Overview"]
    lines.append(f"- Entries: {meta['entries']}")
    lines.append(f"- Sessions: {meta['sessions']}")
    lines.append(f"- Range: {meta['range']['start_iso']} -> {meta['range']['end_iso']}")
    lines.append("")
    lines.append("## Top Sessions")
    for item in report["top_sessions"]:
        lines.append(
            f"- {item['session_id']}: {item['message_count']} messages, {item['duration_minutes']} min, first prompt: {item['first_prompt']}"
        )
    lines.append("")
    lines.append("## Repeated Prompts")
    if report["repeated_prompts"]:
        for item in report["repeated_prompts"]:
            lines.append(f"- x{item['count']}: {item['text']}")
    else:
        lines.append("- No repeated prompts above threshold")
    lines.append("")
    lines.append("## Workflow Topics")
    for topic, count in report["workflow_topics"]["counts"].items():
        lines.append(f"- {topic}: {count}")
    lines.append("")
    lines.append("## Domain Topics")
    for topic, count in report["domain_topics"]["counts"].items():
        lines.append(f"- {topic}: {count}")
    lines.append("")
    lines.append("## Recent Focus")
    lines.append(f"- Window: last {report['recent_focus']['window_size']} entries")
    lines.append("- Workflow: " + ", ".join(f"{k}={v}" for k, v in report["recent_focus"]["workflow"].items()))
    lines.append("- Domain: " + ", ".join(f"{k}={v}" for k, v in report["recent_focus"]["domain"].items()))
    focused = report.get("focused_topic")
    if focused:
        lines.append("")
        lines.append(f"## Focused Topic: {focused['kind']}/{focused['topic']}")
        lines.append(f"- Count: {focused['count']}")
        if focused["samples"]:
            for item in focused["samples"]:
                lines.append(f"- {item['ts']}: {item['text']}")
        else:
            lines.append("- No samples matched after filtering")
    return "\n".join(lines) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(description="Export structured Codex history facts from history.jsonl.")
    parser.add_argument("--input", default=str(Path.home() / ".codex" / "history.jsonl"))
    parser.add_argument("--top", type=int, default=8)
    parser.add_argument("--recent", type=int, default=40)
    parser.add_argument("--samples", type=int, default=3)
    parser.add_argument("--topic")
    parser.add_argument("--topic-kind", choices=["auto", "workflow", "domain"], default="auto")
    parser.add_argument("--topic-samples", type=int, default=12)
    parser.add_argument("--format", choices=["json", "markdown"], default="json")
    parser.add_argument("--output", default="-")
    args = parser.parse_args()

    entries = load_entries(Path(args.input))
    report = build_report(
        entries,
        top_n=max(args.top, 1),
        recent_n=max(args.recent, 1),
        sample_per_topic=max(args.samples, 1),
        topic=args.topic,
        topic_kind=args.topic_kind,
        topic_samples_limit=max(args.topic_samples, 1),
    )
    if args.format == "json":
        payload = json.dumps(report, ensure_ascii=False, indent=2) + "\n"
    else:
        payload = render_markdown(report)
    if args.output == "-":
        print(payload, end="")
    else:
        Path(args.output).write_text(payload)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
