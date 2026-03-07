#!/usr/bin/env python3
import re
import sys
from pathlib import Path

ALLOWED_FRONTMATTER_KEYS = {"name", "description"}
DESCRIPTION_PREFIXES = ("Use when", "Use before", "Use after")
PRIMARY_REFERENCE_HEADINGS = (
    "## 重点做法",
    "## 重点检查",
    "## 检查清单",
    "## 重点关注",
    "## 处理顺序",
    "## 模板",
    "## 诊断顺序",
    "## 最小骨架建议",
    "## 推荐结构",
    "## 选择顺序",
    "## 常用类型",
    "## P0",
)
TAIL_REFERENCE_HEADINGS = (
    "## 注意事项",
    "## 使用要求",
    "## 使用建议",
    "## 优先更新",
    "## 高风险信号",
    "## 发布补充",
    "## 设计目标",
    "## Scope 选择",
    "## 选择信号",
    "## 常见输出",
    "## Open Questions",
)


def parse_frontmatter(path: Path, text: str) -> dict[str, str]:
    if not text.startswith("---\n"):
        raise ValueError("missing YAML frontmatter start")
    parts = text.split("---\n", 2)
    if len(parts) < 3:
        raise ValueError("missing YAML frontmatter end")
    raw = parts[1]
    data: dict[str, str] = {}
    for line in raw.splitlines():
        if not line.strip():
            continue
        match = re.match(r"^([A-Za-z0-9_-]+):\s*(.+)$", line)
        if not match:
            raise ValueError(f"invalid frontmatter line: {line}")
        key, value = match.groups()
        data[key] = value.strip()
    missing = ALLOWED_FRONTMATTER_KEYS - data.keys()
    extra = data.keys() - ALLOWED_FRONTMATTER_KEYS
    if missing:
        raise ValueError(f"missing frontmatter keys: {', '.join(sorted(missing))}")
    if extra:
        raise ValueError(f"unexpected frontmatter keys: {', '.join(sorted(extra))}")
    return data


def lint_skill(skill_dir: Path) -> list[str]:
    errors: list[str] = []
    skill_md = skill_dir / "SKILL.md"
    if not skill_md.exists():
        return [f"{skill_dir}: missing SKILL.md"]
    text = skill_md.read_text()
    try:
        frontmatter = parse_frontmatter(skill_md, text)
    except ValueError as exc:
        return [f"{skill_md}: {exc}"]
    if frontmatter["name"] != skill_dir.name:
        errors.append(f"{skill_md}: frontmatter name must match directory name '{skill_dir.name}'")
    if not frontmatter["description"].startswith(DESCRIPTION_PREFIXES):
        errors.append(f"{skill_md}: description must start with 'Use when', 'Use before', or 'Use after'")
    links = set(re.findall(r"\]\((references/[^)]+\.md)\)", text))
    refs_dir = skill_dir / "references"
    if refs_dir.exists():
        refs = sorted(refs_dir.glob("*.md"))
        linked_files = {Path(link).name for link in links}
        for link in links:
            if not (skill_dir / link).exists():
                errors.append(f"{skill_md}: broken reference link {link}")
        for ref in refs:
            if ref.name not in linked_files:
                errors.append(f"{skill_md}: unlinked reference file {ref.name}")
            errors.extend(lint_reference(ref))
    return errors


def lint_reference(path: Path) -> list[str]:
    errors: list[str] = []
    text = path.read_text()
    stripped = [line.strip() for line in text.splitlines() if line.strip()]
    if not stripped or not stripped[0].startswith("# "):
        errors.append(f"{path}: first non-empty line must be an H1 title")
        return errors
    if len(stripped) < 2 or not stripped[1].startswith("用于"):
        errors.append(f"{path}: second non-empty line must be a purpose sentence starting with '用于'")
    if "## 何时查看" not in text:
        errors.append(f"{path}: missing '## 何时查看'")
    if not any(heading in text for heading in PRIMARY_REFERENCE_HEADINGS):
        errors.append(f"{path}: missing a primary structure section such as '重点做法' or '检查清单'")
    if not any(heading in text for heading in TAIL_REFERENCE_HEADINGS):
        errors.append(f"{path}: missing a tail section such as '注意事项' or '使用建议'")
    return errors


def main() -> int:
    script_path = Path(__file__).resolve()
    candidates = [
        script_path.parents[1] / "skills",
        script_path.parents[2] / "codex" / "skills",
    ]
    skills_root = next((path for path in candidates if path.exists()), None)
    if skills_root is None:
        expected = ", ".join(str(path) for path in candidates)
        print(f"skills directory not found; checked: {expected}", file=sys.stderr)
        return 2
    errors: list[str] = []
    for skill_dir in sorted(p for p in skills_root.iterdir() if p.is_dir() and not p.name.startswith(".")):
        errors.extend(lint_skill(skill_dir))
    if errors:
        print("SKILL_LINT_FAILED")
        for error in errors:
            print(f"- {error}")
        return 1
    print("SKILL_LINT_OK")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
