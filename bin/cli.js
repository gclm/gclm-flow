#!/usr/bin/env node

import { program } from 'commander';
import { runInteractive, runInstall, runUpdate, runUninstall, runList } from '../src/index.js';

program
  .name('gclm-flow')
  .description('Gclm-Flow - 全栈开发工作流增强')
  .version('1.0.0');

// 安装命令（默认命令）
program
  .command('install', { isDefault: true })
  .description('安装 Gclm-Flow（无参数时进入交互模式）')
  .option('-y, --yes', '跳过确认，使用默认配置')
  .option('-f, --force', '强制覆盖已存在的文件')
  .action(async (options) => {
    const { yes, force } = options;

    if (yes) {
      // -y 参数：直接安装，使用默认配置
      await runInstall({ yes: true, force });
    } else {
      // 无 -y 参数：进入交互模式
      await runInteractive();
    }
  });

// 更新命令
program
  .command('update')
  .description('更新 Gclm-Flow 到最新版本')
  .action(async () => {
    await runUpdate();
  });

// 卸载命令
program
  .command('uninstall')
  .description('卸载 Gclm-Flow')
  .action(async () => {
    await runUninstall();
  });

// 列出组件
program
  .command('list')
  .description('列出所有可用组件')
  .action(async () => {
    await runList();
  });

program.parse();
