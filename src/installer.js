import fs from 'fs-extra';
import path from 'path';
import os from 'os';
import chalk from 'chalk';
import ora from 'ora';
import { glob } from 'glob';
import {
  getConfigDir,
  getGclmFlowDir,
  getSourceDir,
  getPackageVersion,
  saveInstalledVersion,
  getComponents,
  getPlatformConfig,
  LANGUAGES
} from './utils.js';

/**
 * 安装组件
 */
export async function install(platform, selectedComponents, selectedLanguages, options = {}) {
  const { force = false, silent = false } = options;
  const configDir = getConfigDir(platform);
  const sourceDir = getSourceDir();
  const version = getPackageVersion();
  const platformConfig = getPlatformConfig(platform);
  const components = getComponents(platform);

  const spinner = silent ? null : ora(`正在安装 Gclm-Flow 到 ${platformConfig.name}...`).start();
  const messages = [];

  try {
    // 确保目录存在
    await fs.ensureDir(configDir);
    await fs.ensureDir(getGclmFlowDir());

    const installedFiles = [];

    // 安装组件
    for (const componentKey of selectedComponents) {
      const component = components[componentKey];
      if (!component) continue;

      if (spinner) spinner.text = `正在安装 ${component.name}...`;

      const result = await installComponent(component, sourceDir, configDir, force, platform);
      installedFiles.push(...result.files);
      if (result.messages) messages.push(...result.messages);
    }

    // 安装语言规则（仅 Claude Code）
    if (platform === 'claude-code') {
      for (const langKey of selectedLanguages) {
        const lang = LANGUAGES[langKey];
        if (!lang) continue;

        if (spinner) spinner.text = `正在安装 ${lang.name} 规则...`;

        const files = await installLanguageRule(lang, sourceDir, configDir, force);
        installedFiles.push(...files);
      }
    }

    // 保存版本信息
    await saveInstalledVersion(platform, version, selectedComponents, selectedLanguages, installedFiles);

    if (spinner) spinner.succeed(chalk.green(`安装完成！(${platformConfig.name})`));

    // 显示收集的消息
    if (messages.length > 0) {
      console.log();
      messages.forEach(msg => console.log(chalk.gray(`  ${msg}`)));
    }

    console.log();
    console.log(chalk.cyan('  已安装组件:'));
    selectedComponents.forEach(key => {
      if (components[key]) {
        console.log(`    ${chalk.green('✓')} ${components[key].name}`);
      }
    });

    if (platform === 'claude-code' && selectedLanguages.length > 0) {
      console.log();
      console.log(chalk.cyan('  已安装语言规则:'));
      selectedLanguages.forEach(key => {
        console.log(`    ${chalk.green('✓')} ${LANGUAGES[key].name}`);
      });
    }

    console.log();
    console.log(chalk.gray(`  共安装 ${installedFiles.length} 个文件`));
    console.log();
    console.log(chalk.cyan('  快速开始:'));
    if (platform === 'claude-code') {
      console.log(chalk.gray('    1. 重启 Claude Code'));
      console.log(chalk.gray('    2. 使用 /gclm 智能编排'));
      console.log(chalk.gray('    3. 使用 /gclm-init 初始化项目'));
    } else {
      console.log(chalk.gray('    1. 重启 Codex CLI'));
      console.log(chalk.gray('    2. 使用 /gclm 智能编排'));
      console.log(chalk.gray('    3. 使用 /gclm-init 初始化项目'));
    }
    console.log();

  } catch (error) {
    if (spinner) spinner.fail(chalk.red('安装失败'));
    console.error(error);
    throw error;
  }
}

/**
 * 安装单个组件
 */
async function installComponent(component, sourceDir, targetDir, force, platform) {
  // 处理特殊组件（userConfig）
  if (component.special) {
    return await installSpecialComponent(component, sourceDir, targetDir, force, platform);
  }

  const installedFiles = [];
  const messages = [];
  const sourcePath = path.join(sourceDir, component.source);
  const homeDir = os.homedir();

  // 处理 skills 组件的平台特定路径
  let targetPath;
  let relativeBase; // 用于计算相对路径的基准目录
  if (component.source === 'skills') {
    const platformConfig = getPlatformConfig(platform);
    // skillsTarget 是相对于 home 目录的路径，如 '.claude/skills' 或 '.agents/skills'
    targetPath = path.join(homeDir, platformConfig.skillsTarget);
    relativeBase = homeDir;
  } else {
    targetPath = path.join(targetDir, component.target);
    relativeBase = targetDir;
  }

  await fs.ensureDir(path.dirname(targetPath));

  if (component.recursive) {
    const pattern = path.join(sourcePath, component.pattern).replace(/\\/g, '/');
    const files = await glob(pattern, { nodir: true });

    for (const file of files) {
      const relativePath = path.relative(sourcePath, file);
      const destPath = path.join(targetPath, relativePath);

      await fs.ensureDir(path.dirname(destPath));

      if (!force && await fs.pathExists(destPath)) {
        await fs.copy(destPath, `${destPath}.backup`);
      }

      await fs.copy(file, destPath);
      installedFiles.push(path.relative(relativeBase, destPath));
    }
  } else {
    const pattern = path.join(sourcePath, component.pattern).replace(/\\/g, '/');
    const files = await glob(pattern, { nodir: true });

    for (const file of files) {
      const destPath = path.join(targetPath, path.basename(file));

      if (!force && await fs.pathExists(destPath)) {
        await fs.copy(destPath, `${destPath}.backup`);
      }

      await fs.copy(file, destPath);
      installedFiles.push(path.relative(relativeBase, destPath));
    }
  }

  // Codex CLI: 同时安装到 ~/.codex/skills 以兼容不同版本
  if (platform === 'codex-cli' && component.source === 'skills') {
    const codexSkillsDir = path.join(homeDir, '.codex', 'skills');
    const codexInstalledFiles = await copySkillsToDir(sourcePath, codexSkillsDir, component, force);
    installedFiles.push(...codexInstalledFiles);
    if (codexInstalledFiles.length > 0) {
      messages.push(`同时安装到 ~/.codex/skills (${codexInstalledFiles.length} 个文件)`);
    }
  }

  return { files: installedFiles, messages };
}

/**
 * 复制 skills 到指定目录（用于 Codex 多目录兼容）
 */
async function copySkillsToDir(sourcePath, targetDir, component, force) {
  const installedFiles = [];
  const homeDir = os.homedir();

  await fs.ensureDir(targetDir);

  if (component.recursive) {
    const pattern = path.join(sourcePath, component.pattern).replace(/\\/g, '/');
    const files = await glob(pattern, { nodir: true });

    for (const file of files) {
      const relativePath = path.relative(sourcePath, file);
      const destPath = path.join(targetDir, relativePath);

      await fs.ensureDir(path.dirname(destPath));

      if (!force && await fs.pathExists(destPath)) {
        await fs.copy(destPath, `${destPath}.backup`);
      }

      await fs.copy(file, destPath);
      installedFiles.push(path.relative(homeDir, destPath));
    }
  } else {
    const pattern = path.join(sourcePath, component.pattern).replace(/\\/g, '/');
    const files = await glob(pattern, { nodir: true });

    for (const file of files) {
      const destPath = path.join(targetDir, path.basename(file));

      if (!force && await fs.pathExists(destPath)) {
        await fs.copy(destPath, `${destPath}.backup`);
      }

      await fs.copy(file, destPath);
      installedFiles.push(path.relative(homeDir, destPath));
    }
  }

  return installedFiles;
}

/**
 * 安装特殊组件（userConfig）
 */
async function installSpecialComponent(component, sourceDir, targetDir, force, platform) {
  const installedFiles = [];
  const messages = [];
  const sourcePath = path.join(sourceDir, component.source);
  const platformConfig = getPlatformConfig(platform);
  const homeDir = os.homedir();

  // 安装配置文件（CLAUDE.md 或 AGENTS.md）
  const configSource = path.join(sourcePath, platformConfig.templateName);
  const configTarget = path.join(targetDir, platformConfig.configFileName);

  if (await fs.pathExists(configSource)) {
    if (!force && await fs.pathExists(configTarget)) {
      await fs.copy(configTarget, `${configTarget}.backup`);
    }
    await fs.copy(configSource, configTarget);
    // 使用相对于 home 目录的路径
    installedFiles.push(path.relative(homeDir, configTarget));
    messages.push(`安装 ${platformConfig.configFileName} → ${configTarget}`);
  }

  // 合并配置到 settings.json（仅 Claude Code）
  if (platform === 'claude-code') {
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
          settings.permissions = settings.permissions || {};
          if (permissionsConfig.permissions.allow) {
            const existingAllow = new Set(settings.permissions.allow || []);
            permissionsConfig.permissions.allow.forEach(item => existingAllow.add(item));
            settings.permissions.allow = Array.from(existingAllow);
          }
          if (permissionsConfig.permissions.deny) {
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
      // 使用相对于 home 目录的路径
      installedFiles.push(path.relative(homeDir, settingsTarget));
    }
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
  const homeDir = os.homedir();

  if (!await fs.pathExists(sourcePath)) {
    console.log(chalk.yellow(`  警告: 规则文件不存在 ${lang.source}`));
    return installedFiles;
  }

  await fs.ensureDir(path.dirname(targetPath));

  if (!force && await fs.pathExists(targetPath)) {
    await fs.copy(targetPath, `${targetPath}.backup`);
  }

  await fs.copy(sourcePath, targetPath);
  // 使用 home 目录作为基准，保持与其他组件一致
  installedFiles.push(path.relative(homeDir, targetPath));

  return installedFiles;
}

/**
 * 卸载组件
 */
export async function uninstall(platform) {
  const configDir = getConfigDir(platform);
  const platformConfig = getPlatformConfig(platform);
  const versionFile = path.join(configDir, platformConfig.versionFile);
  const homeDir = os.homedir();

  const spinner = ora(`正在卸载 Gclm-Flow (${platformConfig.name})...`).start();

  try {
    let installedFiles = [];

    if (await fs.pathExists(versionFile)) {
      const versionInfo = await fs.readJson(versionFile);
      installedFiles = versionInfo.installedFiles || [];
    }

    let deletedCount = 0;
    for (const file of installedFiles) {
      // 文件路径是相对于 home 目录存储的
      const filePath = path.join(homeDir, file);
      if (await fs.pathExists(filePath)) {
        await fs.remove(filePath);
        deletedCount++;
      }
    }

    await cleanupEmptyDirs(configDir);
    // 清理 skills 目录（如果与 configDir 不同）
    const skillsDir = path.join(homeDir, platformConfig.skillsTarget);
    if (skillsDir !== configDir) {
      await cleanupEmptyDirs(skillsDir);
    }
    // Codex CLI: 同时清理 ~/.codex/skills
    if (platform === 'codex-cli') {
      const codexSkillsDir = path.join(homeDir, '.codex', 'skills');
      await cleanupEmptyDirs(codexSkillsDir);
    }
    await fs.remove(versionFile);

    spinner.succeed(chalk.green(`卸载完成！(${platformConfig.name})`));
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

    for (const dir of dirs.sort().reverse()) {
      const fullPath = path.join(rootDir, dir);

      try {
        const stat = await fs.stat(fullPath);
        if (!stat.isDirectory()) continue;

        const files = await fs.readdir(fullPath);
        if (files.length === 0) {
          await fs.rmdir(fullPath);
        }
      } catch {
        continue;
      }
    }
  } catch {
    // 忽略清理错误
  }
}

/**
 * 更新组件
 */
export async function update(platform, selectedComponents, selectedLanguages) {
  const configDir = getConfigDir(platform);
  const platformConfig = getPlatformConfig(platform);
  const versionFile = path.join(configDir, platformConfig.versionFile);

  if (await fs.pathExists(versionFile)) {
    const versionInfo = await fs.readJson(versionFile);
    selectedComponents = selectedComponents || versionInfo.components;
    selectedLanguages = selectedLanguages || versionInfo.languages;
  }

  await install(platform, selectedComponents, selectedLanguages, { force: true });
}

/**
 * 安装所有平台
 */
export async function installAll(selectedComponents, selectedLanguages, options = {}) {
  console.log(chalk.cyan('安装到所有平台...\n'));

  // Claude Code
  await install('claude-code', selectedComponents, selectedLanguages, options);
  console.log();

  // Codex CLI
  await install('codex-cli', ['skills', 'userConfig'], [], options);
}
