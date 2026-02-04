# 安装指南

## 前置要求

- **Claude Code**: 已安装 Claude Code CLI
- **Bash/Zsh**: 支持 Bash 或 Zsh shell
- **Git**: 用于版本控制 (可选)

---

## 安装步骤

### 1. 克隆或下载项目

```bash
# 克隆项目
git clone https://github.com/your-username/gclm-flow.git
cd gclm-flow

# 或下载并解压
wget https://github.com/your-username/gclm-flow/archive/main.zip
unzip main.zip
cd gclm-flow-main
```

---

### 2. 运行安装脚本

```bash
./install.sh
```

安装脚本会：
- 检查 Claude Code 目录
- 复制 agents/ 到 `~/.claude/agents/`
- 复制 commands/ 到 `~/.claude/commands/`
- 复制 skills/ 到 `~/.claude/skills/`
- 复制 rules/ 到 `~/.claude/rules/`
- 复制 hooks/ 到 `~/.claude/hooks/`
- 设置可执行权限

---

### 3. 配置 Claude Code

#### 3.1 复制核心配置

```bash
# 复制 CLAUDE.example.md 到 ~/.claude/CLAUDE.md
cp CLAUDE.example.md ~/.claude/CLAUDE.md

# 或追加到现有配置
cat CLAUDE.example.md >> ~/.claude/CLAUDE.md
```

#### 3.2 配置 MCP (可选)

```bash
# 复制 MCP 配置
cp settings.example.json ~/.claude/settings.local.json
```

---

### 4. 安装 auggie (推荐)

```bash
npm install -g @augmentcode/auggie@prerelease
```

**验证安装**:
```bash
which auggie
# 应输出: /usr/local/bin/auggie 或类似路径
```

---

### 5. 验证安装

#### 检查文件

```bash
# 检查 agents
ls ~/.claude/agents/
# 应包含: investigator.md, architect.md, worker.md, tdd-guide.md, spec-guide.md, code-reviewer.md

# 检查 commands
ls ~/.claude/commands/
# 应包含: gclm.md, investigate.md, tdd.md, spec.md, llmdoc.md

# 检查 skills
ls ~/.claude/skills/
# 应包含: gclm/, file-naming-helper/

# 检查 hooks
ls ~/.claude/hooks/
# 应包含: notify.sh, stop-gclm-loop.sh
```

#### 检查权限

```bash
ls -la ~/.claude/hooks/
# notify.sh 和 stop-gclm-loop.sh 应该有可执行权限 (755 或 rwxr-xr-x)
```

---

## 配置选项

### 环境变量

| 变量 | 描述 | 默认值 |
|:---|:---|:---|
| `CLAUDE_DIR` | Claude Code 配置目录 | `~/.claude` |
| `PROJECT_DIR` | gclm-flow 项目目录 | 自动检测 |

### MCP 设置

如果使用 MCP 服务器，编辑 `~/.claude/settings.local.json`:

```json
{
  "mcpServers": {
    "auggie": {
      "command": "npx",
      "args": ["@augmentcode/auggie@prerelease"]
    }
  }
}
```

---

## 故障排查

### 问题: install.sh 报错 "Permission denied"

**解决**:
```bash
chmod +x install.sh
./install.sh
```

---

### 问题: hooks/ 目录为空

**解决**: 检查 `install.sh` 中的路径是否正确

---

### 问题: auggie 命令未找到

**解决**:
1. 检查 npm 全局路径: `npm config get prefix`
2. 添加到 PATH: `export PATH="$PATH:$(npm config get prefix)/bin"`
3. 添加到 shell 配置文件 (`~/.bashrc` 或 `~/.zshrc`)

---

### 问题: Claude Code 无法识别 skills

**解决**:
1. 检查 skills 目录权限
2. 确保 SKILL.md 文件存在
3. 重启 Claude Code

---

## 卸载

### 手动删除

```bash
# 删除 agents
rm -rf ~/.claude/agents/gclm-*

# 删除 commands
rm -rf ~/.claude/commands/gclm.md
rm -rf ~/.claude/commands/investigate.md
rm -rf ~/.claude/commands/tdd.md
rm -rf ~/.claude/commands/spec.md
rm -rf ~/.claude/commands/llmdoc.md

# 删除 skills
rm -rf ~/.claude/skills/gclm
rm -rf ~/.claude/skills/file-naming-helper

# 删除 hooks
rm -rf ~/.claude/hooks/notify.sh
rm -rf ~/.claude/hooks/stop-gclm-loop.sh

# 删除 rules (如果只包含 gclm-flow 规则)
rm -rf ~/.claude/rules/

# 删除配置
rm ~/.claude/CLAUDE.md
```

---

## 更新

```bash
cd gclm-flow
git pull
./install.sh
```
