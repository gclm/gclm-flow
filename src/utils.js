import os from 'os';
import path from 'path';
import fs from 'fs-extra';
import { fileURLToPath } from 'url';

/**
 * 组件定义
 */
export const COMPONENTS = {
  agents: {
    name: 'Agents（代理）',
    description: '6个专用代理（planner, builder, reviewer, investigator, recorder, remember）',
    source: 'agents',
    target: 'agents',
    pattern: '*.md'
  },
  rules: {
    name: 'Rules（规则）',
    description: '分层规则架构（core + 语言特定）',
    source: 'rules',
    target: 'rules',
    pattern: '**/*.md',
    recursive: true
  },
  commands: {
    name: 'Commands（命令）',
    description: '12个斜杠命令（/gclm:init, /gclm:plan, /gclm:do 等）',
    source: 'commands/gclm',
    target: 'commands/gclm',
    pattern: '*.md'
  },
  skills: {
    name: 'Skills（技能）',
    description: '工作流定义和领域知识',
    source: 'skills',
    target: 'skills',
    pattern: '**/*.md',
    recursive: true
  },
  hooks: {
    name: 'Hooks（钩子）',
    description: '基于触发器的自动化',
    source: 'hooks',
    target: 'hooks',
    pattern: 'hooks.json'
  },
  userConfig: {
    name: 'User Config（用户配置）',
    description: '全局用户级配置（CLAUDE.md + statusline + MCP servers）',
    source: 'templates',
    target: '.',
    pattern: '*.md',
    special: true  // 特殊标记，需要特殊处理
  }
};

/**
 * 语言规则定义
 */
export const LANGUAGES = {
  common: {
    name: 'common（必装）',
    description: '通用规则',
    source: 'rules/core.md',
    target: 'rules/core.md',
    required: true
  },
  java: {
    name: 'Java/Spring Boot',
    description: 'Java 和 Spring Boot 规则',
    source: 'rules/languages/java.md',
    target: 'rules/languages/java.md'
  },
  python: {
    name: 'Python/Flask/FastAPI',
    description: 'Python、Flask、FastAPI 规则',
    source: 'rules/languages/python.md',
    target: 'rules/languages/python.md'
  },
  golang: {
    name: 'Go/Gin',
    description: 'Go 和 Gin 规则',
    source: 'rules/languages/go.md',
    target: 'rules/languages/go.md'
  },
  rust: {
    name: 'Rust/Axum/Actix',
    description: 'Rust、Axum、Actix 规则',
    source: 'rules/languages/rust.md',
    target: 'rules/languages/rust.md'
  },
  frontend: {
    name: '前端（TypeScript/React/Vue）',
    description: '前端规则',
    source: 'rules/languages/frontend.md',
    target: 'rules/languages/frontend.md'
  }
};

/**
 * 获取 Claude 配置目录
 */
export function getClaudeDir() {
  return path.join(os.homedir(), '.claude');
}

/**
 * 获取 Gclm-Flow 数据目录
 */
export function getGclmFlowDir() {
  return path.join(os.homedir(), '.gclm-flow');
}

/**
 * 获取版本文件路径
 */
export function getVersionFilePath() {
  return path.join(getClaudeDir(), '.gclm-version');
}

/**
 * 获取包版本
 */
export function getPackageVersion() {
  try {
    const pkgPath = path.join(getSourceDir(), '..', 'package.json');
    const pkg = fs.readJsonSync(pkgPath);
    return pkg.version;
  } catch {
    return '1.0.0';
  }
}

/**
 * 获取已安装的版本信息
 */
export async function getInstalledVersion() {
  const versionFile = getVersionFilePath();
  if (await fs.pathExists(versionFile)) {
    return fs.readJson(versionFile);
  }
  return null;
}

/**
 * 保存版本信息
 */
export async function saveInstalledVersion(version, components, languages, installedFiles) {
  const versionFile = getVersionFilePath();
  await fs.writeJson(versionFile, {
    version,
    components,
    languages,
    installedFiles,
    installedAt: new Date().toISOString()
  }, { spaces: 2 });
}

/**
 * 获取源目录（组件所在目录）
 */
export function getSourceDir() {
  // 获取当前模块的目录路径，然后返回项目根目录
  const __filename = fileURLToPath(import.meta.url);
  const __dirname = path.dirname(__filename);
  // src 的父目录就是项目根目录
  return path.dirname(__dirname);
}
