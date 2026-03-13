# rules/

按语言分层的编码规范，部署到 `~/.claude/rules/`。

## 来源

| 目录 | 来源 | 更新方式 |
|---|---|---|
| `common/` | ECC submodule | `git submodule update --remote` |
| `golang/` | ECC submodule | 同上 |
| `python/` | ECC submodule | 同上 |
| `typescript/` | ECC submodule | 同上（Next.js） |
| `java/` | 本仓库 `rules/java/` | 直接编辑 |
| `rust/` | 本仓库 `rules/rust/` | 直接编辑 |

ECC submodule 路径：`vendor/everything-claude-code/`

## 更新 ECC rules

```bash
git submodule update --remote vendor/everything-claude-code
git add vendor/everything-claude-code
git commit -m "chore(vendor): 更新 ECC submodule"
bash claude/bin/sync-to-home.sh
```

## Rules vs Skills

- **rules/**：始终生效的编码标准和检查清单（what to do）
- **skills/**：任务级深度指南，按需加载（how to do it）

语言 rule 文件会在适当位置引用对应 skill。

## 结构

```
rules/
├── java/
│   ├── coding-style.md
│   ├── testing.md
│   ├── security.md
│   └── patterns.md
└── rust/
    ├── coding-style.md
    ├── testing.md
    ├── security.md
    └── patterns.md
```
