#!/usr/bin/env node

import { program } from 'commander';
import { runInteractive, runInstall, runUpdate, runUninstall, runList, runStatus } from '../src/index.js';

program
  .name('gclm-flow')
  .description('Gclm-Flow - 全栈开发工作流增强，支持 Claude Code 和 Codex CLI')
  .version('1.0.0');

// 安装命令（默认命令）
program
  .command('install', { isDefault: true })
  .description('安装 Gclm-Flow（无参数时进入交互模式）')
  .option('-y, --yes', '跳过确认，使用默认配置')
  .option('-f, --force', '强制覆盖已存在的文件')
  .option('-p, --platform <platform>', '目标平台 (claude-code | codex-cli | all)')
  .option('--all', '安装到所有平台')
  .action(async (options) => {
    const { yes, force, platform, all } = options;

    if (yes || platform || all) {
      await runInstall({ yes: true, force, platform, all });
    } else {
      await runInteractive();
    }
  });

// 更新命令
program
  .command('update')
  .description('更新 Gclm-Flow 到最新版本')
  .option('-p, --platform <platform>', '目标平台 (claude-code | codex-cli)')
  .action(async (options) => {
    await runUpdate(options.platform);
  });

// 卸载命令
program
  .command('uninstall')
  .description('卸载 Gclm-Flow')
  .option('-p, --platform <platform>', '目标平台 (claude-code | codex-cli)')
  .action(async (options) => {
    await runUninstall(options.platform);
  });

// 列出组件
program
  .command('list')
  .description('列出所有可用组件')
  .option('-p, --platform <platform>', '目标平台 (claude-code | codex-cli | all)')
  .action(async (options) => {
    await runList(options.platform);
  });

// 状态命令
program
  .command('status')
  .description('显示安装状态')
  .action(async () => {
    await runStatus();
  });

program.parse();
