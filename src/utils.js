import os from 'os';
import path from 'path';
import fs from 'fs-extra';
import { fileURLToPath } from 'url';

/**
 * 支持的平台
 */
export const PLATFORMS = {
  'claude-code': {
    name: 'Claude Code',
    configDir: '.claude',
    configFileName: 'CLAUDE.md',
    templateName: 'CLAUDE-template.md',
    skillsTarget: '.claude/skills',   // ~/.claude/skills
    versionFile: '.gclm-flow-version'
  },
  'codex-cli': {
    name: 'Codex CLI',
    configDir: '.codex',              // 配置文件目录
    configFileName: 'AGENTS.md',
    templateName: 'AGENTS-template.md',
    skillsTarget: '.agents/skills',   // ~/.agents/skills
    versionFile: '.gclm-flow-version'
  }
};

/**
 * Claude Code 组件定义
 */
export const CLAUDE_COMPONENTS = {
  agents: {
    name: 'Agents（代理）',
    description: '6个专用代理',
    source: 'agents',
    target: 'agents',
    pattern: '*.md'
  },
  rules: {
    name: 'Rules（规则）',
    description: '分层规则架构',
    source: 'rules',
    target: 'rules',
    pattern: '**/*.md',
    recursive: true
  },
  skills: {
    name: 'Skills（技能）',
    description: '工作流和语言栈技能',
    source: 'skills',
    target: 'skills',
    pattern: '**/*.md',
    recursive: true
  },
  hooks: {
    name: 'Hooks（钩子）',
    description: '自动化钩子',
    source: 'hooks',
    target: 'hooks',
    pattern: 'hooks.json'
  },
  userConfig: {
    name: 'User Config（用户配置）',
    description: '全局配置',
    source: 'templates',
    target: '.',
    pattern: '*.md',
    special: true
  }
};

/**
 * Codex CLI 组件定义
 */
export const CODEX_COMPONENTS = {
  skills: {
    name: 'Skills（技能）',
    description: '工作流和语言栈技能',
    source: 'skills',
    target: 'skills',
    pattern: '**/*.md',
    recursive: true
  },
  userConfig: {
    name: 'User Config（用户配置）',
    description: '全局配置',
    source: 'templates',
    target: '.',
    pattern: '*.md',
    special: true
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
 * 获取平台配置
 */
export function getPlatformConfig(platform) {
  return PLATFORMS[platform] || PLATFORMS['claude-code'];
}

/**
 * 获取组件定义
 */
export function getComponents(platform) {
  return platform === 'codex-cli' ? CODEX_COMPONENTS : CLAUDE_COMPONENTS;
}

/**
 * 获取配置目录
 */
export function getConfigDir(platform) {
  const config = getPlatformConfig(platform);
  return path.join(os.homedir(), config.configDir);
}

/**
 * 获取 Claude 配置目录（兼容旧 API）
 */
export function getClaudeDir() {
  return getConfigDir('claude-code');
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
export function getVersionFilePath(platform) {
  const config = getPlatformConfig(platform);
  return path.join(getConfigDir(platform), config.versionFile);
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
export async function getInstalledVersion(platform) {
  const versionFile = getVersionFilePath(platform);
  if (await fs.pathExists(versionFile)) {
    return fs.readJson(versionFile);
  }
  return null;
}

/**
 * 保存版本信息
 */
export async function saveInstalledVersion(platform, version, components, languages, installedFiles) {
  const versionFile = getVersionFilePath(platform);
  await fs.writeJson(versionFile, {
    platform,
    version,
    components,
    languages,
    installedFiles,
    installedAt: new Date().toISOString()
  }, { spaces: 2 });
}

/**
 * 获取源目录
 */
export function getSourceDir() {
  const __filename = fileURLToPath(import.meta.url);
  const __dirname = path.dirname(__filename);
  return path.dirname(__dirname);
}

/**
 * 检测平台是否已安装
 */
export async function isPlatformInstalled(platform) {
  const configDir = getConfigDir(platform);
  const versionFile = getVersionFilePath(platform);
  return await fs.pathExists(configDir) && await fs.pathExists(versionFile);
}
