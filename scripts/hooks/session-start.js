/**
 * Session Start Hook
 * 加载上一次会话的上下文，检测包管理器
 */

let data = '';
process.stdin.on('data', chunk => data += chunk);
process.stdin.on('end', async () => {
  try {
    const fs = require('fs');
    const path = require('path');
    const os = require('os');

    const gclmDir = path.join(os.homedir(), '.gclm-flow');
    const sessionFile = path.join(gclmDir, 'session.json');

    // 确保 gclm-flow 目录存在
    if (!fs.existsSync(gclmDir)) {
      fs.mkdirSync(gclmDir, { recursive: true });
    }

    // 加载上一次会话状态
    if (fs.existsSync(sessionFile)) {
      const session = JSON.parse(fs.readFileSync(sessionFile, 'utf-8'));
      if (session.lastProject) {
        console.error(`[GCLM] 上次项目: ${session.lastProject}`);
      }
      if (session.packageManager) {
        console.error(`[GCLM] 包管理器: ${session.packageManager}`);
      }
    }

    // 检测当前项目的包管理器
    const cwd = process.cwd();
    let packageManager = 'npm';

    if (fs.existsSync(path.join(cwd, 'pnpm-lock.yaml'))) {
      packageManager = 'pnpm';
    } else if (fs.existsSync(path.join(cwd, 'yarn.lock'))) {
      packageManager = 'yarn';
    } else if (fs.existsSync(path.join(cwd, 'bun.lockb'))) {
      packageManager = 'bun';
    } else if (fs.existsSync(path.join(cwd, 'package-lock.json'))) {
      packageManager = 'npm';
    } else if (fs.existsSync(path.join(cwd, 'go.mod'))) {
      packageManager = 'go';
    } else if (fs.existsSync(path.join(cwd, 'Cargo.toml'))) {
      packageManager = 'cargo';
    } else if (fs.existsSync(path.join(cwd, 'pom.xml'))) {
      packageManager = 'maven';
    } else if (fs.existsSync(path.join(cwd, 'build.gradle'))) {
      packageManager = 'gradle';
    }

    // 保存当前会话状态
    const session = {
      lastProject: cwd,
      packageManager,
      startedAt: new Date().toISOString()
    };
    fs.writeFileSync(sessionFile, JSON.stringify(session, null, 2));

    console.error(`[GCLM] 会话已初始化 (${packageManager})`);

  } catch (e) {
    console.error(`[GCLM] Session start error: ${e.message}`);
  }

  // 输出原始数据
  console.log(data);
});
