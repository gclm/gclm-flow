import fs from 'fs-extra';
import path from 'path';
import chalk from 'chalk';
import ora from 'ora';
import { glob } from 'glob';
import {
  getClaudeDir,
  getGclmFlowDir,
  getSourceDir,
  getPackageVersion,
  saveInstalledVersion
} from './utils.js';
import { COMPONENTS, LANGUAGES } from './utils.js';

/**
 * 安装组件
 */
export async function install(selectedComponents, selectedLanguages, options = {}) {
  const { force = false, silent = false } = options;
  const claudeDir = getClaudeDir();
  const sourceDir = getSourceDir();
  const version = getPackageVersion();

  const spinner = silent ? null : ora('正在安装 Gclm-Flow...').start();
  const messages = []; // 收集安装过程中的消息

  try {
    // 确保目录存在
    await fs.ensureDir(claudeDir);
    await fs.ensureDir(getGclmFlowDir());

    const installedFiles = [];

    // 安装组件
    for (const componentKey of selectedComponents) {
      const component = COMPONENTS[componentKey];
      if (!component) continue;

      if (spinner) spinner.text = `正在安装 ${component.name}...`;

      const result = await installComponent(component, sourceDir, claudeDir, force);
      installedFiles.push(...result.files);
      if (result.messages) messages.push(...result.messages);
    }

    // 安装语言规则
    for (const langKey of selectedLanguages) {
      const lang = LANGUAGES[langKey];
      if (!lang) continue;

      if (spinner) spinner.text = `正在安装 ${lang.name} 规则...`;

      const files = await installLanguageRule(lang, sourceDir, claudeDir, force);
      installedFiles.push(...files);
    }

    // 保存版本信息
    await saveInstalledVersion(version, selectedComponents, selectedLanguages, installedFiles);

    if (spinner) spinner.succeed(chalk.green('安装完成！'));

    // 显示收集的消息
    if (messages.length > 0) {
      console.log();
      messages.forEach(msg => console.log(chalk.gray(`  ${msg}`)));
    }

    console.log();
    console.log(chalk.cyan('  已安装组件:'));
    selectedComponents.forEach(key => {
      console.log(`    ${chalk.green('✓')} ${COMPONENTS[key].name}`);
    });
    console.log();
    console.log(chalk.cyan('  已安装语言规则:'));
    selectedLanguages.forEach(key => {
      console.log(`    ${chalk.green('✓')} ${LANGUAGES[key].name}`);
    });
    console.log();
    console.log(chalk.gray(`  共安装 ${installedFiles.length} 个文件`));
    console.log();
    console.log(chalk.cyan('  快速开始:'));
    console.log(chalk.gray('    1. 重启 Claude Code'));
    console.log(chalk.gray('    2. 使用 /gclm:init 初始化项目'));
    console.log(chalk.gray('    3. 使用 /gclm:plan 规划任务'));
    console.log();

  } catch (error) {
    if (spinner) spinner.fail(chalk.red('安装失败'));
    console.error(error);
    throw error;
  }
}

/**
 * 安装单个组件
 * @returns {{ files: string[], messages: string[] }}
 */
async function installComponent(component, sourceDir, targetDir, force) {
  // 处理特殊组件（userConfig）
  if (component.special) {
    return await installSpecialComponent(component, sourceDir, targetDir, force);
  }

  const installedFiles = [];
  const sourcePath = path.join(sourceDir, component.source);
  const targetPath = path.join(targetDir, component.target);

  // 确保目标目录存在
  await fs.ensureDir(path.dirname(targetPath));

  if (component.recursive) {
    // 递归安装目录
    const pattern = path.join(sourcePath, component.pattern).replace(/\\/g, '/');
    const files = await glob(pattern, { nodir: true });

    for (const file of files) {
      const relativePath = path.relative(sourcePath, file);
      const destPath = path.join(targetPath, relativePath);

      await fs.ensureDir(path.dirname(destPath));

      if (!force && await fs.pathExists(destPath)) {
        // 备份已存在的文件
        await fs.copy(destPath, `${destPath}.backup`);
      }

      await fs.copy(file, destPath);
      installedFiles.push(path.relative(targetDir, destPath));
    }
  } else {
    // 单文件或通配符
    const pattern = path.join(sourcePath, component.pattern).replace(/\\/g, '/');
    const files = await glob(pattern, { nodir: true });

    for (const file of files) {
      const destPath = path.join(targetPath, path.basename(file));

      if (!force && await fs.pathExists(destPath)) {
        await fs.copy(destPath, `${destPath}.backup`);
      }

      await fs.copy(file, destPath);
      installedFiles.push(path.relative(targetDir, destPath));
    }
  }

  return { files: installedFiles, messages: [] };
}

/**
 * 安装特殊组件（userConfig）
 * @returns {{ files: string[], messages: string[] }}
 */
async function installSpecialComponent(component, sourceDir, targetDir, force) {
  const installedFiles = [];
  const messages = [];
  const sourcePath = path.join(sourceDir, component.source);

  // 安装 CLAUDE.md
  const claudeMdSource = path.join(sourcePath, 'user-CLAUDE.md');
  const claudeMdTarget = path.join(targetDir, 'CLAUDE.md');

  if (await fs.pathExists(claudeMdSource)) {
    if (!force && await fs.pathExists(claudeMdTarget)) {
      await fs.copy(claudeMdTarget, `${claudeMdTarget}.backup`);
    }
    await fs.copy(claudeMdSource, claudeMdTarget);
    installedFiles.push('CLAUDE.md');
    messages.push(`安装 CLAUDE.md → ${claudeMdTarget}`);
  }

  // 合并 statusline.json、mcp-servers.json、permissions.json 到 settings.json
  const statuslineSource = path.join(sourcePath, 'statusline.json');
  const mcpSource = path.join(sourcePath, 'mcp-servers.json');
  const permissionsSource = path.join(sourcePath, 'permissions.json');
  const settingsTarget = path.join(targetDir, 'settings.json');

  const hasConfigUpdates = await fs.pathExists(statuslineSource) ||
                           await fs.pathExists(mcpSource) ||
                           await fs.pathExists(permissionsSource);

  if (hasConfigUpdates) {
    let settings = {};
    if (await fs.pathExists(settingsTarget)) {
      if (!force) {
        await fs.copy(settingsTarget, `${settingsTarget}.backup`);
      }
      settings = await fs.readJson(settingsTarget);
    }

    // 合并 statusLine 配置
    if (await fs.pathExists(statuslineSource)) {
      const statuslineConfig = await fs.readJson(statuslineSource);
      if (statuslineConfig.statusLine) {
        settings.statusLine = statuslineConfig.statusLine;
      }
    }

    // 合并 mcpServers 配置
    if (await fs.pathExists(mcpSource)) {
      const mcpConfig = await fs.readJson(mcpSource);
      if (mcpConfig.mcpServers) {
        settings.mcpServers = { ...settings.mcpServers, ...mcpConfig.mcpServers };
      }
    }

    // 合并 permissions 配置
    if (await fs.pathExists(permissionsSource)) {
      const permissionsConfig = await fs.readJson(permissionsSource);
      if (permissionsConfig.permissions) {
        // 合并 allow 列表（去重）
        if (permissionsConfig.permissions.allow) {
          settings.permissions = settings.permissions || {};
          const existingAllow = new Set(settings.permissions.allow || []);
          permissionsConfig.permissions.allow.forEach(item => existingAllow.add(item));
          settings.permissions.allow = Array.from(existingAllow);
        }
        // 合并 deny 列表（去重）
        if (permissionsConfig.permissions.deny) {
          settings.permissions = settings.permissions || {};
          const existingDeny = new Set(settings.permissions.deny || []);
          permissionsConfig.permissions.deny.forEach(item => existingDeny.add(item));
          settings.permissions.deny = Array.from(existingDeny);
        }
      }
      messages.push('更新 settings.json（statusLine + mcpServers + permissions）');
    } else {
      messages.push('更新 settings.json（statusLine + mcpServers）');
    }

    await fs.writeJson(settingsTarget, settings, { spaces: 2 });
    installedFiles.push('settings.json');
  }

  return { files: installedFiles, messages };
}

/**
 * 安装语言规则
 */
async function installLanguageRule(lang, sourceDir, targetDir, force) {
  const installedFiles = [];
  const sourcePath = path.join(sourceDir, lang.source);
  const targetPath = path.join(targetDir, lang.target);

  // 检查源文件是否存在
  if (!await fs.pathExists(sourcePath)) {
    console.log(chalk.yellow(`  警告: 规则文件不存在 ${lang.source}`));
    return installedFiles;
  }

  await fs.ensureDir(path.dirname(targetPath));

  if (!force && await fs.pathExists(targetPath)) {
    await fs.copy(targetPath, `${targetPath}.backup`);
  }

  await fs.copy(sourcePath, targetPath);
  installedFiles.push(path.relative(targetDir, targetPath));

  return installedFiles;
}

/**
 * 卸载组件
 */
export async function uninstall(selectedComponents) {
  const claudeDir = getClaudeDir();
  const spinner = ora('正在卸载 Gclm-Flow...').start();

  try {
    const versionFile = path.join(claudeDir, '.gclm-version');
    let installedFiles = [];

    // 读取已安装文件列表
    if (await fs.pathExists(versionFile)) {
      const versionInfo = await fs.readJson(versionFile);
      installedFiles = versionInfo.installedFiles || [];
    }

    // 删除已安装的文件
    let deletedCount = 0;
    for (const file of installedFiles) {
      const filePath = path.join(claudeDir, file);
      if (await fs.pathExists(filePath)) {
        await fs.remove(filePath);
        deletedCount++;
      }
    }

    // 清理空目录
    await cleanupEmptyDirs(claudeDir);

    // 删除版本文件
    await fs.remove(versionFile);

    spinner.succeed(chalk.green('卸载完成！'));
    console.log();
    console.log(chalk.gray(`  已删除 ${deletedCount} 个文件`));

  } catch (error) {
    spinner.fail(chalk.red('卸载失败'));
    console.error(error);
    throw error;
  }
}

/**
 * 清理空目录
 */
async function cleanupEmptyDirs(rootDir) {
  try {
    const dirs = await glob('**/*/', { cwd: rootDir });

    // 从最深层开始删除空目录
    for (const dir of dirs.sort().reverse()) {
      const fullPath = path.join(rootDir, dir);

      // 检查是否是目录
      try {
        const stat = await fs.stat(fullPath);
        if (!stat.isDirectory()) continue;

        const files = await fs.readdir(fullPath);
        if (files.length === 0) {
          await fs.rmdir(fullPath);
        }
      } catch {
        // 忽略错误，继续处理其他目录
        continue;
      }
    }
  } catch (error) {
    // 忽略清理错误
  }
}

/**
 * 更新组件
 */
export async function update(selectedComponents, selectedLanguages) {
  const claudeDir = getClaudeDir();
  const versionFile = path.join(claudeDir, '.gclm-version');

  // 读取当前安装信息
  if (await fs.pathExists(versionFile)) {
    const versionInfo = await fs.readJson(versionFile);
    selectedComponents = selectedComponents || versionInfo.components;
    selectedLanguages = selectedLanguages || versionInfo.languages;
  }

  // 强制覆盖安装
  await install(selectedComponents, selectedLanguages, { force: true });
}
