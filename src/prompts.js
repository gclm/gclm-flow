import inquirer from 'inquirer';
import chalk from 'chalk';
import { COMPONENTS, LANGUAGES } from './utils.js';

/**
 * 选择组件
 */
export async function promptComponents() {
  console.log(chalk.cyan('  选择要安装的组件:'));
  console.log(chalk.gray('  (按空格选择，回车确认)'));
  console.log();

  const choices = Object.entries(COMPONENTS).map(([key, component]) => ({
    name: `${component.name} - ${chalk.gray(component.description)}`,
    value: key,
    checked: true // 默认全选
  }));

  const { components } = await inquirer.prompt([
    {
      type: 'checkbox',
      name: 'components',
      message: '选择组件',
      choices,
      pageSize: 10
    }
  ]);

  if (components.length === 0) {
    console.log(chalk.yellow('  至少需要选择一个组件'));
    return promptComponents();
  }

  return components;
}

/**
 * 选择语言规则
 */
export async function promptLanguages() {
  console.log();
  console.log(chalk.cyan('  选择语言规则:'));
  console.log(chalk.gray('  common 为必装，其他根据需要选择'));
  console.log();

  const choices = Object.entries(LANGUAGES)
    .filter(([key]) => key !== 'common') // common 不在列表中显示
    .map(([key, lang]) => ({
      name: `${lang.name} - ${chalk.gray(lang.description)}`,
      value: key,
      checked: ['java', 'python', 'golang', 'frontend'].includes(key) // 默认选择常用语言
    }));

  const { languages } = await inquirer.prompt([
    {
      type: 'checkbox',
      name: 'languages',
      message: '选择语言',
      choices,
      pageSize: 10
    }
  ]);

  // common 必装
  return ['common', ...languages];
}

/**
 * 确认安装
 */
export async function confirmInstall(components, languages) {
  console.log();
  console.log(chalk.cyan('  安装确认:'));
  console.log();
  console.log(chalk.white('  将安装以下组件:'));

  components.forEach(key => {
    const component = COMPONENTS[key];
    if (component) {
      console.log(`    ${chalk.green('●')} ${component.name}`);
    }
  });

  console.log();
  console.log(chalk.white('  将安装以下语言规则:'));

  languages.forEach(key => {
    const lang = LANGUAGES[key];
    if (lang) {
      const required = lang.required ? chalk.red('(必装)') : '';
      console.log(`    ${chalk.green('●')} ${lang.name} ${required}`);
    }
  });

  console.log();
  console.log(chalk.gray('  目标目录: ~/.claude'));
  console.log(chalk.gray('  数据目录: ~/.gclm-flow'));
  console.log();

  const { confirmed } = await inquirer.prompt([
    {
      type: 'confirm',
      name: 'confirmed',
      message: '确认安装？',
      default: true
    }
  ]);

  return confirmed;
}

/**
 * 选择更新操作
 */
export async function promptUpdateAction() {
  const { action } = await inquirer.prompt([
    {
      type: 'list',
      name: 'action',
      message: '选择操作',
      choices: [
        { name: '更新所有组件', value: 'all' },
        { name: '选择性更新', value: 'select' },
        { name: '取消', value: 'cancel' }
      ]
    }
  ]);

  return action;
}

/**
 * 确认卸载
 */
export async function confirmUninstall() {
  console.log(chalk.yellow('  警告: 这将删除所有 Gclm-Flow 组件'));
  console.log();

  const { confirmed } = await inquirer.prompt([
    {
      type: 'confirm',
      name: 'confirmed',
      message: '确认卸载？',
      default: false
    }
  ]);

  return confirmed;
}
