/**
 * Post Edit Format Hook
 * 编辑后自动格式化代码
 */

let data = '';
process.stdin.on('data', chunk => data += chunk);
process.stdin.on('end', () => {
  try {
    const input = JSON.parse(data);
    const filePath = input.tool_input?.file_path || '';

    if (!filePath) {
      console.log(data);
      return;
    }

    const fs = require('fs');
    const { execSync } = require('child_process');
    const path = require('path');

    if (!fs.existsSync(filePath)) {
      console.log(data);
      return;
    }

    const ext = path.extname(filePath);

    // JS/TS 文件用 Prettier
    if (['.js', '.jsx', '.ts', '.tsx'].includes(ext)) {
      try {
        execSync(`npx prettier --write "${filePath}"`, { stdio: 'pipe' });
        console.error(`[GCLM] 已格式化: ${path.basename(filePath)}`);
      } catch (e) {
        // Prettier 可能不存在，忽略错误
      }
    }

    // Python 文件用 ruff
    if (ext === '.py') {
      try {
        execSync(`ruff format "${filePath}"`, { stdio: 'pipe' });
        console.error(`[GCLM] 已格式化: ${path.basename(filePath)}`);
      } catch (e) {
        // ruff 可能不存在，忽略错误
      }
    }

    // Go 文件
    if (ext === '.go') {
      try {
        execSync(`gofmt -w "${filePath}"`, { stdio: 'pipe' });
        console.error(`[GCLM] 已格式化: ${path.basename(filePath)}`);
      } catch (e) {
        // gofmt 可能不存在，忽略错误
      }
    }

    // Rust 文件
    if (ext === '.rs') {
      try {
        execSync(`rustfmt "${filePath}"`, { stdio: 'pipe' });
        console.error(`[GCLM] 已格式化: ${path.basename(filePath)}`);
      } catch (e) {
        // rustfmt 可能不存在，忽略错误
      }
    }

  } catch (e) {
    // 忽略所有错误
  }

  console.log(data);
});
