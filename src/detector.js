import fs from 'fs-extra';
import path from 'path';
import { getGclmFlowDir } from './utils.js';

/**
 * 项目检测配置
 */
const DETECTION_PATTERNS = {
  java: {
    files: ['pom.xml', 'build.gradle', 'build.gradle.kts'],
    frameworks: {
      springboot: {
        files: ['pom.xml', 'build.gradle'],
        contentPatterns: ['spring-boot', 'org.springframework']
      },
      quarkus: {
        files: ['pom.xml', 'build.gradle'],
        contentPatterns: ['quarkus']
      }
    }
  },
  python: {
    files: ['requirements.txt', 'pyproject.toml', 'setup.py', 'Pipfile'],
    frameworks: {
      flask: {
        files: ['requirements.txt', 'pyproject.toml'],
        contentPatterns: ['flask', 'Flask']
      },
      fastapi: {
        files: ['requirements.txt', 'pyproject.toml'],
        contentPatterns: ['fastapi', 'FastAPI']
      },
      django: {
        files: ['requirements.txt', 'pyproject.toml'],
        contentPatterns: ['django', 'Django']
      }
    }
  },
  golang: {
    files: ['go.mod', 'go.sum'],
    frameworks: {
      gin: {
        files: ['go.mod'],
        contentPatterns: ['github.com/gin-gonic/gin']
      },
      echo: {
        files: ['go.mod'],
        contentPatterns: ['github.com/labstack/echo']
      },
      fiber: {
        files: ['go.mod'],
        contentPatterns: ['github.com/gofiber/fiber']
      }
    }
  },
  rust: {
    files: ['Cargo.toml'],
    frameworks: {
      axum: {
        files: ['Cargo.toml'],
        contentPatterns: ['axum']
      },
      actix: {
        files: ['Cargo.toml'],
        contentPatterns: ['actix-web']
      },
      rocket: {
        files: ['Cargo.toml'],
        contentPatterns: ['rocket']
      }
    }
  },
  frontend: {
    files: ['package.json'],
    frameworks: {
      react: {
        files: ['package.json'],
        contentPatterns: ['"react"', '"next"']
      },
      vue: {
        files: ['package.json'],
        contentPatterns: ['"vue"', '"nuxt"']
      },
      angular: {
        files: ['package.json'],
        contentPatterns: ['"@angular/core"']
      },
      svelte: {
        files: ['package.json'],
        contentPatterns: ['"svelte"', '"@sveltejs/kit"']
      }
    }
  }
};

/**
 * 构建工具检测配置
 */
const BUILD_TOOLS = {
  java: {
    maven: { files: ['pom.xml'] },
    gradle: { files: ['build.gradle', 'build.gradle.kts'] }
  },
  python: {
    pip: { files: ['requirements.txt'] },
    poetry: { files: ['pyproject.toml'] },
    pipenv: { files: ['Pipfile'] }
  },
  golang: {
    go_modules: { files: ['go.mod'] }
  },
  rust: {
    cargo: { files: ['Cargo.toml'] }
  },
  frontend: {
    npm: { files: ['package.json'] },
    yarn: { files: ['yarn.lock'] },
    pnpm: { files: ['pnpm-lock.yaml'] }
  }
};

/**
 * 测试框架检测配置
 */
const TEST_FRAMEWORKS = {
  java: {
    junit: { patterns: ['junit', 'JUnit'] },
    testng: { patterns: ['testng'] },
    mockito: { patterns: ['mockito'] }
  },
  python: {
    pytest: { patterns: ['pytest', 'Pytest'] },
    unittest: { patterns: ['unittest'] }
  },
  golang: {
    testing: { patterns: ['testing'] } // Go 内置测试
  },
  rust: {
    cargo_test: { patterns: [] } // Rust 内置测试
  },
  frontend: {
    jest: { patterns: ['jest', 'Jest'] },
    vitest: { patterns: ['vitest', 'Vitest'] },
    mocha: { patterns: ['mocha', 'Mocha'] },
    playwright: { patterns: ['playwright'] },
    cypress: { patterns: ['cypress'] }
  }
};

/**
 * 项目检测器
 */
export class ProjectDetector {
  constructor(projectPath) {
    this.projectPath = projectPath;
    this.cache = null;
  }

  /**
   * 检测项目
   */
  async detect() {
    if (this.cache) {
      return this.cache;
    }

    const result = {
      path: this.projectPath,
      name: path.basename(this.projectPath),
      type: await this.detectProjectType(),
      languages: await this.detectLanguages(),
      frameworks: await this.detectFrameworks(),
      buildTools: await this.detectBuildTools(),
      testFrameworks: await this.detectTestFrameworks(),
      detectedAt: new Date().toISOString()
    };

    this.cache = result;
    return result;
  }

  /**
   * 检测项目类型
   */
  async detectProjectType() {
    const languages = await this.detectLanguages();

    if (languages.includes('frontend') && languages.length === 1) {
      return 'frontend';
    }

    if (languages.includes('frontend') && languages.length > 1) {
      return 'fullstack';
    }

    if (languages.length > 0) {
      return 'backend';
    }

    return 'unknown';
  }

  /**
   * 检测语言
   */
  async detectLanguages() {
    const languages = [];

    for (const [lang, config] of Object.entries(DETECTION_PATTERNS)) {
      for (const file of config.files) {
        if (await fs.pathExists(path.join(this.projectPath, file))) {
          if (!languages.includes(lang)) {
            languages.push(lang);
          }
          break;
        }
      }
    }

    return languages;
  }

  /**
   * 检测框架
   */
  async detectFrameworks() {
    const frameworks = {};

    for (const [lang, config] of Object.entries(DETECTION_PATTERNS)) {
      // 检查是否有该语言的文件
      let hasLanguage = false;
      for (const file of config.files) {
        if (await fs.pathExists(path.join(this.projectPath, file))) {
          hasLanguage = true;
          break;
        }
      }

      if (!hasLanguage) continue;

      frameworks[lang] = [];

      // 检测框架
      for (const [frameworkName, frameworkConfig] of Object.entries(config.frameworks)) {
        if (await this.detectFramework(frameworkConfig)) {
          frameworks[lang].push(frameworkName);
        }
      }
    }

    return frameworks;
  }

  /**
   * 检测单个框架
   */
  async detectFramework(config) {
    for (const file of config.files) {
      const filePath = path.join(this.projectPath, file);

      if (!await fs.pathExists(filePath)) {
        continue;
      }

      if (config.contentPatterns && config.contentPatterns.length > 0) {
        try {
          const content = await fs.readFile(filePath, 'utf-8');
          for (const pattern of config.contentPatterns) {
            if (content.includes(pattern)) {
              return true;
            }
          }
        } catch {
          continue;
        }
      } else {
        return true;
      }
    }

    return false;
  }

  /**
   * 检测构建工具
   */
  async detectBuildTools() {
    const buildTools = {};

    for (const [lang, tools] of Object.entries(BUILD_TOOLS)) {
      buildTools[lang] = [];

      for (const [toolName, config] of Object.entries(tools)) {
        for (const file of config.files) {
          if (await fs.pathExists(path.join(this.projectPath, file))) {
            buildTools[lang].push(toolName);
            break;
          }
        }
      }
    }

    return buildTools;
  }

  /**
   * 检测测试框架
   */
  async detectTestFrameworks() {
    const testFrameworks = {};

    for (const [lang, frameworks] of Object.entries(TEST_FRAMEWORKS)) {
      testFrameworks[lang] = [];

      // 检查依赖文件
      const depFiles = DETECTION_PATTERNS[lang]?.files || [];

      for (const [frameworkName, config] of Object.entries(frameworks)) {
        // Go 和 Rust 有内置测试
        if (lang === 'golang' && frameworkName === 'testing') {
          // 检查是否有 _test.go 文件
          const testFiles = await this.findTestFiles('go');
          if (testFiles.length > 0) {
            testFrameworks[lang].push(frameworkName);
          }
          continue;
        }

        if (lang === 'rust' && frameworkName === 'cargo_test') {
          // Rust 内置测试
          testFrameworks[lang].push(frameworkName);
          continue;
        }

        // 其他框架从依赖文件中检测
        for (const file of depFiles) {
          const filePath = path.join(this.projectPath, file);

          if (!await fs.pathExists(filePath)) {
            continue;
          }

          try {
            const content = await fs.readFile(filePath, 'utf-8');
            for (const pattern of config.patterns) {
              if (content.includes(pattern)) {
                testFrameworks[lang].push(frameworkName);
                break;
              }
            }
          } catch {
            continue;
          }
        }
      }
    }

    return testFrameworks;
  }

  /**
   * 查找测试文件
   */
  async findTestFiles(lang) {
    const patterns = {
      go: '**/*_test.go',
      rust: '**/tests/**/*.rs',
      java: '**/src/test/**/*.java',
      python: '**/test_*.py',
      frontend: '**/*.test.{js,ts,jsx,tsx}'
    };

    // 这里简化处理，实际可以使用 glob
    return [];
  }
}

/**
 * 检测缓存管理
 */
export class DetectionCache {
  constructor() {
    this.cacheDir = path.join(getGclmFlowDir(), 'cache');
    this.indexFile = path.join(this.cacheDir, 'projects.json');
  }

  /**
   * 获取项目缓存
   */
  async get(projectPath) {
    await this.ensureCacheDir();

    if (await fs.pathExists(this.indexFile)) {
      const index = await fs.readJson(this.indexFile);
      const entry = index.projects[projectPath];

      if (entry && !this.isExpired(entry)) {
        return entry.data;
      }
    }

    return null;
  }

  /**
   * 保存项目缓存
   */
  async set(projectPath, data) {
    await this.ensureCacheDir();

    let index = { version: '1.0', projects: {} };

    if (await fs.pathExists(this.indexFile)) {
      index = await fs.readJson(this.indexFile);
    }

    index.projects[projectPath] = {
      data,
      cachedAt: new Date().toISOString(),
      ttl: 86400000 // 24 小时
    };

    await fs.writeJson(this.indexFile, index, { spaces: 2 });
  }

  /**
   * 清除缓存
   */
  async clear() {
    if (await fs.pathExists(this.cacheDir)) {
      await fs.remove(this.cacheDir);
    }
  }

  /**
   * 确保缓存目录存在
   */
  async ensureCacheDir() {
    await fs.ensureDir(this.cacheDir);
  }

  /**
   * 检查缓存是否过期
   */
  isExpired(entry) {
    if (!entry.ttl) return false;

    const cachedAt = new Date(entry.cachedAt).getTime();
    const now = Date.now();

    return (now - cachedAt) > entry.ttl;
  }
}

/**
 * 快速检测项目（带缓存）
 */
export async function detectProject(projectPath) {
  const cache = new DetectionCache();

  // 尝试从缓存获取
  const cached = await cache.get(projectPath);
  if (cached) {
    return cached;
  }

  // 执行检测
  const detector = new ProjectDetector(projectPath);
  const result = await detector.detect();

  // 保存缓存
  await cache.set(projectPath, result);

  return result;
}

/**
 * 生成检测报告
 */
export function generateDetectionReport(result) {
  const lines = [];

  lines.push('# 项目检测报告');
  lines.push('');
  lines.push(`**项目**: ${result.name}`);
  lines.push(`**路径**: ${result.path}`);
  lines.push(`**类型**: ${result.type}`);
  lines.push('');

  lines.push('## 检测结果');
  lines.push('');

  lines.push('### 语言');
  result.languages.forEach(lang => {
    lines.push(`- ${lang}`);
  });
  lines.push('');

  lines.push('### 框架');
  for (const [lang, frameworks] of Object.entries(result.frameworks)) {
    if (frameworks.length > 0) {
      lines.push(`- ${lang}: ${frameworks.join(', ')}`);
    }
  }
  lines.push('');

  lines.push('### 构建工具');
  for (const [lang, tools] of Object.entries(result.buildTools)) {
    if (tools.length > 0) {
      lines.push(`- ${lang}: ${tools.join(', ')}`);
    }
  }
  lines.push('');

  lines.push('### 测试框架');
  for (const [lang, frameworks] of Object.entries(result.testFrameworks)) {
    if (frameworks.length > 0) {
      lines.push(`- ${lang}: ${frameworks.join(', ')}`);
    }
  }

  return lines.join('\n');
}
