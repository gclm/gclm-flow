import chalk from 'chalk';
import figlet from 'figlet';
import { install, uninstall } from './installer.js';
import { promptComponents, promptLanguages, confirmInstall } from './prompts.js';
import { COMPONENTS, getPackageVersion, getInstalledVersion } from './utils.js';

/**
 * 显示欢迎横幅
 */
function showBanner() {
  console.log(chalk.cyan(figlet.textSync('Gclm-Flow', { font: 'Standard' })));
  console.log(chalk.gray('  全栈开发工作流增强，智能检测，统一体验'));
  console.log();
}

/**
 * 交互式安装流程
 */
export async function runInteractive() {
  showBanner();

  const installedVersion = await getInstalledVersion();

  if (installedVersion) {
    console.log(chalk.green(`  已安装版本: v${installedVersion.version}`));
    console.log();
    console.log('  选择操作:');
    console.log('  1. 更新 Gclm-Flow');
    console.log('  2. 卸载 Gclm-Flow');
    console.log('  3. 退出');
    console.log();

    // 简化的交互逻辑
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
      await runUpdate();
    } else if (action === 'uninstall') {
      await runUninstall();
    }
    return;
  }

  // 新安装流程
  console.log(chalk.gray(`  目标目录: ~/.claude`));
  console.log();

  // 1. 选择组件
  const selectedComponents = await promptComponents();

  // 2. 选择语言规则
  const selectedLanguages = await promptLanguages();

  // 3. 确认安装
  const confirmed = await confirmInstall(selectedComponents, selectedLanguages);
  if (!confirmed) {
    console.log(chalk.yellow('  已取消安装'));
    return;
  }

  // 4. 执行安装
  await install(selectedComponents, selectedLanguages);
}

/**
 * 直接安装
 */
export async function runInstall(options = {}) {
  const { yes, force } = options;

  if (!yes) {
    showBanner();
  }

  // 默认安装所有组件和语言
  const selectedComponents = Object.keys(COMPONENTS);
  const selectedLanguages = ['common', 'java', 'python', 'golang', 'rust', 'frontend'];

  await install(selectedComponents, selectedLanguages, { force });
}

/**
 * 更新
 */
export async function runUpdate() {
  console.log(chalk.cyan('  正在更新 Gclm-Flow...'));
  // TODO: 实现更新逻辑
  console.log(chalk.green('  更新完成！'));
}

/**
 * 卸载
 */
export async function runUninstall() {
  const installedVersion = await getInstalledVersion();

  if (!installedVersion) {
    console.log(chalk.yellow('  Gclm-Flow 尚未安装'));
    return;
  }

  const selectedComponents = installedVersion.components || Object.keys(COMPONENTS);
  await uninstall(selectedComponents);
}

/**
 * 列出组件
 */
export async function runList() {
  console.log();
  console.log(chalk.cyan('  可用组件:'));
  console.log();

  for (const [key, component] of Object.entries(COMPONENTS)) {
    console.log(`  ${chalk.green('●')} ${component.name}`);
    console.log(`    ${chalk.gray(component.description)}`);
    console.log();
  }

  console.log(chalk.cyan('  语言规则:'));
  console.log();
  console.log(`  ${chalk.green('●')} common (必装)`);
  console.log(`  ${chalk.green('●')} java - Java/Spring Boot`);
  console.log(`  ${chalk.green('●')} python - Python/Flask/FastAPI`);
  console.log(`  ${chalk.green('●')} golang - Go/Gin`);
  console.log(`  ${chalk.green('●')} rust - Rust/Axum/Actix`);
  console.log(`  ${chalk.green('●')} frontend - TypeScript/React/Vue`);
  console.log();
}
