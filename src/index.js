import chalk from 'chalk';
import figlet from 'figlet';
import { install, uninstall, update, installAll } from './installer.js';
import { promptComponents, promptLanguages, confirmInstall, promptPlatform } from './prompts.js';
import { getComponents, PLATFORMS, getPackageVersion, getInstalledVersion } from './utils.js';

/**
 * 显示欢迎横幅
 */
function showBanner() {
  console.log(chalk.cyan(figlet.textSync('Gclm-Flow', { font: 'Standard' })));
  console.log(chalk.gray('  全栈开发工作流增强，支持 Claude Code 和 Codex CLI'));
  console.log();
}

/**
 * 交互式安装流程
 */
export async function runInteractive() {
  showBanner();

  // 选择平台
  const platform = await promptPlatform();

  // 处理"全部平台"选项
  if (platform === 'all') {
    console.log(chalk.cyan('  将安装到所有平台：'));
    console.log(chalk.gray('    - Claude Code (~/.claude)'));
    console.log(chalk.gray('    - Codex CLI (~/.codex + ~/.agents/skills)'));
    console.log();

    // 使用默认配置安装
    const claudeComponents = ['agents', 'rules', 'skills', 'hooks', 'userConfig'];
    const languages = ['common', 'java', 'python', 'golang', 'rust', 'frontend'];
    const codexComponents = ['skills', 'userConfig'];

    console.log(chalk.cyan('  Claude Code 组件:'));
    claudeComponents.forEach(c => console.log(chalk.gray(`    - ${c}`)));
    console.log();
    console.log(chalk.cyan('  Codex CLI 组件:'));
    codexComponents.forEach(c => console.log(chalk.gray(`    - ${c}`)));
    console.log();

    const { confirmed } = await import('inquirer').then(m => m.default.prompt([{
      type: 'confirm',
      name: 'confirmed',
      message: '确认安装到所有平台？',
      default: true
    }]));

    if (!confirmed) {
      console.log(chalk.yellow('  已取消安装'));
      return;
    }

    await installAll(claudeComponents, languages, { force: false });
    return;
  }

  const platformConfig = PLATFORMS[platform];
  const components = getComponents(platform);
  const installedVersion = await getInstalledVersion(platform);

  if (installedVersion) {
    console.log(chalk.green(`  已安装版本: v${installedVersion.version} (${platformConfig.name})`));
    console.log();
    console.log('  选择操作:');
    console.log('  1. 更新 Gclm-Flow');
    console.log('  2. 卸载 Gclm-Flow');
    console.log('  3. 退出');
    console.log();

    const { action } = await import('inquirer').then(m => m.default.prompt([{
      type: 'list',
      name: 'action',
      message: '请选择:',
      choices: [
        { name: '更新 Gclm-Flow', value: 'update' },
        { name: '卸载 Gclm-Flow', value: 'uninstall' },
        { name: '退出', value: 'exit' }
      ]
    }]));

    if (action === 'update') {
      await runUpdate(platform);
    } else if (action === 'uninstall') {
      await runUninstall(platform);
    }
    return;
  }

  // 新安装流程
  console.log(chalk.gray(`  目标目录: ~/${platformConfig.configDir}`));
  console.log();

  // 1. 选择组件
  const selectedComponents = await promptComponents(platform);

  // 2. 选择语言规则（仅 Claude Code）
  let selectedLanguages = [];
  if (platform === 'claude-code') {
    selectedLanguages = await promptLanguages();
  }

  // 3. 确认安装
  const confirmed = await confirmInstall(selectedComponents, selectedLanguages, platform);
  if (!confirmed) {
    console.log(chalk.yellow('  已取消安装'));
    return;
  }

  // 4. 执行安装
  await install(platform, selectedComponents, selectedLanguages);
}

/**
 * 直接安装
 */
export async function runInstall(options = {}) {
  const { yes, force, platform: platformOpt, all } = options;

  if (!yes) {
    showBanner();
  }

  // 安装所有平台
  if (all) {
    const selectedComponents = ['agents', 'rules', 'skills', 'hooks', 'userConfig'];
    const selectedLanguages = ['common', 'java', 'python', 'golang', 'rust', 'frontend'];
    await installAll(selectedComponents, selectedLanguages, { force });
    return;
  }

  // 指定平台
  const platform = platformOpt || 'claude-code';
  const components = getComponents(platform);

  // 默认安装所有组件
  const selectedComponents = Object.keys(components);
  const selectedLanguages = platform === 'claude-code'
    ? ['common', 'java', 'python', 'golang', 'rust', 'frontend']
    : [];

  await install(platform, selectedComponents, selectedLanguages, { force });
}

/**
 * 更新
 */
export async function runUpdate(platform = 'claude-code') {
  console.log(chalk.cyan(`  正在更新 Gclm-Flow (${PLATFORMS[platform].name})...`));
  await update(platform);
  console.log(chalk.green('  更新完成！'));
}

/**
 * 卸载
 */
export async function runUninstall(platform = 'claude-code') {
  const installedVersion = await getInstalledVersion(platform);

  if (!installedVersion) {
    console.log(chalk.yellow(`  Gclm-Flow 尚未安装 (${PLATFORMS[platform].name})`));
    return;
  }

  await uninstall(platform);
}

/**
 * 列出组件
 */
export async function runList(platformOpt) {
  console.log();

  if (platformOpt === 'all' || !platformOpt) {
    // 显示所有平台
    for (const [key, config] of Object.entries(PLATFORMS)) {
      console.log(chalk.cyan(`  ${config.name} (~/${config.configDir}):`));
      const components = getComponents(key);

      for (const [compKey, component] of Object.entries(components)) {
        console.log(`    ${chalk.green('●')} ${component.name}`);
        console.log(`      ${chalk.gray(component.description)}`);
      }
      console.log();
    }
  } else {
    const platform = platformOpt || 'claude-code';
    const components = getComponents(platform);

    console.log(chalk.cyan(`  ${PLATFORMS[platform].name} 可用组件:`));
    console.log();

    for (const [key, component] of Object.entries(components)) {
      console.log(`  ${chalk.green('●')} ${component.name}`);
      console.log(`    ${chalk.gray(component.description)}`);
      console.log();
    }
  }

  // 语言规则（仅 Claude Code）
  if (!platformOpt || platformOpt === 'all' || platformOpt === 'claude-code') {
    console.log(chalk.cyan('  语言规则 (仅 Claude Code):'));
    console.log();
    console.log(`  ${chalk.green('●')} common (必装)`);
    console.log(`  ${chalk.green('●')} java - Java/Spring Boot`);
    console.log(`  ${chalk.green('●')} python - Python/Flask/FastAPI`);
    console.log(`  ${chalk.green('●')} golang - Go/Gin`);
    console.log(`  ${chalk.green('●')} rust - Rust/Axum/Actix`);
    console.log(`  ${chalk.green('●')} frontend - TypeScript/React/Vue`);
    console.log();
  }
}

/**
 * 显示状态
 */
export async function runStatus() {
  console.log();
  console.log(chalk.cyan('  Gclm-Flow 安装状态:'));
  console.log();

  for (const [key, config] of Object.entries(PLATFORMS)) {
    const installedVersion = await getInstalledVersion(key);
    const status = installedVersion
      ? chalk.green(`✓ 已安装 v${installedVersion.version}`)
      : chalk.gray('○ 未安装');
    console.log(`  ${config.name}: ${status}`);
  }
  console.log();
}
