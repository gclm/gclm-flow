---
name: browser
description: 浏览器自动化 - 使用 Playwright 进行网页操作、截图、测试
allowed-tools: ["Bash"]
---

# browser Skill

使用 Playwright MCP 进行浏览器自动化操作。

## 安装 Playwright

```bash
# 全局安装 Playwright MCP
npm install -g @microsoft/playwright-mcp

# 安装浏览器
npx playwright install
```

## MCP 配置

```json
{
  "mcpServers": {
    "playwright": {
      "command": "npx",
      "args": ["-y", "@microsoft/playwright-mcp"],
      "env": {
        "PLAYWRIGHT_HEADLESS": "true"
      }
    }
  }
}
```

## 使用场景

### 1. 网页测试
```markdown
打开 https://example.com 并验证标题
```

### 2. 截图
```markdown
截取 https://example.com 的屏幕截图
```

### 3. 表单填写
```markdown
访问登录页面，填写用户名和密码，提交表单
```

### 4. 数据抓取
```markdown
访问网页并提取所有链接
```

### 5. E2E 测试
```markdown
测试用户注册流程：访问注册页 → 填写表单 → 提交 → 验证成功消息
```

## Playwright MCP 工具

### 导航工具
- `playwright_navigate` - 打开网页
- `playwright_goto` - 跳转到 URL

### 交互工具
- `playwright_click` - 点击元素
- `playwright_fill` - 填写表单
- `playwright_type` - 输入文本

### 信息工具
- `playwright_screenshot` - 截图
- `playwright_content` - 获取页面内容
- `playwright_evaluate` - 执行 JavaScript

## 使用示例

### 示例 1: 验证网页
```markdown
使用 Playwright:
1. 访问 https://claude.ai
2. 截取首页截图
3. 检查页面标题
```

### 示例 2: 表单测试
```markdown
使用 Playwright 测试登录:
1. 访问 /login
2. 填写 username: test@example.com
3. 填写 password: ********
4. 点击登录按钮
5. 验证跳转到 /dashboard
```

## 注意事项

1. **Headless 模式**: 默认使用无头模式
2. **等待时间**: 使用 `playwright_wait_for_load_state` 等待页面加载
3. **选择器**: 使用 CSS 选择器或文本选择器
4. **调试**: 设置 `PLAYWRIGHT_HEADLESS=false` 查看浏览器操作

## 常见命令

```bash
# 安装浏览器
npx playwright install

# 调试模式
export PLAYWRIGHT_HEADLESS=false

# 运行 Playwright 代码生成
npx playwright codegen https://example.com
```
