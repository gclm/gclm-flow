import { describe, it, expect } from 'vitest';
import { COMPONENTS, LANGUAGES, getClaudeDir, getGclmFlowDir } from '../src/utils.js';

describe('utils', () => {
  describe('COMPONENTS', () => {
    it('should have all required components', () => {
      expect(COMPONENTS).toHaveProperty('agents');
      expect(COMPONENTS).toHaveProperty('rules');
      expect(COMPONENTS).toHaveProperty('commands');
      expect(COMPONENTS).toHaveProperty('skills');
      expect(COMPONENTS).toHaveProperty('hooks');
    });

    it('should have correct structure for each component', () => {
      for (const [key, component] of Object.entries(COMPONENTS)) {
        expect(component).toHaveProperty('name');
        expect(component).toHaveProperty('description');
        expect(component).toHaveProperty('source');
        expect(component).toHaveProperty('target');
        expect(component).toHaveProperty('pattern');
      }
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

  describe('getGclmFlowDir', () => {
    it('should return path ending with .gclm-flow', () => {
      const dir = getGclmFlowDir();
      expect(dir).toMatch(/\.gclm-flow$/);
    });
  });
});
