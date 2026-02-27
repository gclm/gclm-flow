import { describe, it, expect } from 'vitest';
import {
  PLATFORMS,
  CLAUDE_COMPONENTS,
  CODEX_COMPONENTS,
  LANGUAGES,
  getClaudeDir,
  getGclmFlowDir,
  getConfigDir,
  getComponents,
  getPlatformConfig
} from '../src/utils.js';

describe('utils', () => {
  describe('PLATFORMS', () => {
    it('should have claude-code and codex-cli platforms', () => {
      expect(PLATFORMS).toHaveProperty('claude-code');
      expect(PLATFORMS).toHaveProperty('codex-cli');
    });

    it('should have correct config for each platform', () => {
      for (const [key, config] of Object.entries(PLATFORMS)) {
        expect(config).toHaveProperty('name');
        expect(config).toHaveProperty('configDir');
        expect(config).toHaveProperty('configFileName');
        expect(config).toHaveProperty('templateName');
        expect(config).toHaveProperty('skillsTarget');
        expect(config).toHaveProperty('versionFile');
      }
    });

    it('should use same version file for both platforms', () => {
      expect(PLATFORMS['claude-code'].versionFile).toBe('.gclm-flow-version');
      expect(PLATFORMS['codex-cli'].versionFile).toBe('.gclm-flow-version');
    });
  });

  describe('CLAUDE_COMPONENTS', () => {
    it('should have all required components', () => {
      expect(CLAUDE_COMPONENTS).toHaveProperty('agents');
      expect(CLAUDE_COMPONENTS).toHaveProperty('rules');
      expect(CLAUDE_COMPONENTS).toHaveProperty('skills');
      expect(CLAUDE_COMPONENTS).toHaveProperty('hooks');
      expect(CLAUDE_COMPONENTS).toHaveProperty('userConfig');
    });

    it('should have correct structure for each component', () => {
      for (const [key, component] of Object.entries(CLAUDE_COMPONENTS)) {
        expect(component).toHaveProperty('name');
        expect(component).toHaveProperty('description');
        expect(component).toHaveProperty('source');
        expect(component).toHaveProperty('target');
        expect(component).toHaveProperty('pattern');
      }
    });
  });

  describe('CODEX_COMPONENTS', () => {
    it('should have skills and userConfig', () => {
      expect(CODEX_COMPONENTS).toHaveProperty('skills');
      expect(CODEX_COMPONENTS).toHaveProperty('userConfig');
    });
  });

  describe('LANGUAGES', () => {
    it('should have common as required', () => {
      expect(LANGUAGES.common.required).toBe(true);
    });

    it('should have all supported languages', () => {
      expect(LANGUAGES).toHaveProperty('common');
      expect(LANGUAGES).toHaveProperty('java');
      expect(LANGUAGES).toHaveProperty('python');
      expect(LANGUAGES).toHaveProperty('golang');
      expect(LANGUAGES).toHaveProperty('rust');
      expect(LANGUAGES).toHaveProperty('frontend');
    });
  });

  describe('getClaudeDir', () => {
    it('should return path ending with .claude', () => {
      const dir = getClaudeDir();
      expect(dir).toMatch(/\.claude$/);
    });
  });

  describe('getConfigDir', () => {
    it('should return .claude for claude-code platform', () => {
      const dir = getConfigDir('claude-code');
      expect(dir).toMatch(/\.claude$/);
    });

    it('should return .codex for codex-cli platform', () => {
      const dir = getConfigDir('codex-cli');
      expect(dir).toMatch(/\.codex$/);
    });
  });

  describe('getComponents', () => {
    it('should return CLAUDE_COMPONENTS for claude-code', () => {
      const components = getComponents('claude-code');
      expect(components).toBe(CLAUDE_COMPONENTS);
    });

    it('should return CODEX_COMPONENTS for codex-cli', () => {
      const components = getComponents('codex-cli');
      expect(components).toBe(CODEX_COMPONENTS);
    });
  });

  describe('getPlatformConfig', () => {
    it('should return correct config for claude-code', () => {
      const config = getPlatformConfig('claude-code');
      expect(config.name).toBe('Claude Code');
      expect(config.configFileName).toBe('CLAUDE.md');
    });

    it('should return correct config for codex-cli', () => {
      const config = getPlatformConfig('codex-cli');
      expect(config.name).toBe('Codex CLI');
      expect(config.configFileName).toBe('AGENTS.md');
    });
  });

  describe('getGclmFlowDir', () => {
    it('should return path ending with .gclm-flow', () => {
      const dir = getGclmFlowDir();
      expect(dir).toMatch(/\.gclm-flow$/);
    });
  });
});
