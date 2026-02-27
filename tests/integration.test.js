import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import fs from 'fs-extra';
import path from 'path';
import os from 'os';
import { install, uninstall } from '../src/installer.js';
import {
  getConfigDir,
  getPlatformConfig,
  getComponents,
  getInstalledVersion
} from '../src/utils.js';

describe('installer integration', () => {
  const testDir = path.join(os.tmpdir(), 'gclm-flow-test-' + Date.now());

  beforeAll(async () => {
    await fs.ensureDir(testDir);
  });

  afterAll(async () => {
    await fs.remove(testDir);
  });

  describe('install', () => {
    it('should install skills to codex-cli', async () => {
      const platform = 'codex-cli';
      const components = getComponents(platform);
      const platformConfig = getPlatformConfig(platform);

      // 只安装 skills 组件
      await install(platform, ['skills'], [], { silent: true, force: true });

      // 验证版本文件
      const version = await getInstalledVersion(platform);
      expect(version).not.toBeNull();
      expect(version.platform).toBe(platform);
      expect(version.components).toContain('skills');

      // 验证 skills 目录（Codex CLI 使用 ~/.agents/skills）
      const skillsDir = path.join(os.homedir(), platformConfig.skillsTarget);
      const exists = await fs.pathExists(skillsDir);
      expect(exists).toBe(true);

      // 清理
      await uninstall(platform);
    });

    it('should install userConfig with correct template', async () => {
      const platform = 'codex-cli';
      const configDir = getConfigDir(platform);
      const platformConfig = getPlatformConfig(platform);

      await install(platform, ['userConfig'], [], { silent: true, force: true });

      // 验证 AGENTS.md 已创建
      const agentsMd = path.join(configDir, platformConfig.configFileName);
      const exists = await fs.pathExists(agentsMd);
      expect(exists).toBe(true);

      // 验证内容
      const content = await fs.readFile(agentsMd, 'utf-8');
      expect(content).toContain('Gclm-Flow');
      expect(content).toContain('gclm-core');

      // 清理
      await uninstall(platform);
    });
  });
});
