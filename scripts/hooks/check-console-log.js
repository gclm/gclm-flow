/**
 * Check Console Log Hook
 * 检查代码中的 console.log/debug.Println 等调试语句
 */

let data = '';
process.stdin.on('data', chunk => data += chunk);
process.stdin.on('end', () => {
  try {
    const fs = require('fs');
    const { execSync } = require('child_process');

    // 获取修改的文件
    let modifiedFiles = [];
    try {
      const output = execSync('git diff --name-only HEAD 2>/dev/null', { encoding: 'utf-8' });
      modifiedFiles = output.trim().split('\n').filter(f => f);
    } catch (e) {
      // 不是 git 仓库
    }

    if (modifiedFiles.length === 0) {
      console.log(data);
      return;
    }

    const debugPatterns = {
      '.js': ['console.log', 'console.debug', 'console.warn', 'debugger'],
      '.jsx': ['console.log', 'console.debug', 'console.warn', 'debugger'],
      '.ts': ['console.log', 'console.debug', 'console.warn', 'debugger'],
      '.tsx': ['console.log', 'console.debug', 'console.warn', 'debugger'],
      '.py': ['print(', 'pprint(', 'breakpoint()'],
      '.go': ['fmt.Println', 'fmt.Printf', 'log.Println'],
      '.rs': ['println!', 'dbg!'],
      '.java': ['System.out.print', 'System.err.print']
    };

    let found = false;

    for (const file of modifiedFiles) {
      if (!fs.existsSync(file)) continue;

      const ext = file.substring(file.lastIndexOf('.'));
      const patterns = debugPatterns[ext];
      if (!patterns) continue;

      const content = fs.readFileSync(file, 'utf-8');
      const lines = content.split('\n');

      for (let i = 0; i < lines.length; i++) {
        for (const pattern of patterns) {
          if (lines[i].includes(pattern)) {
            if (!found) {
              console.error('[GCLM] 发现调试语句:');
              found = true;
            }
            console.error(`  ${file}:${i + 1} - ${pattern}`);
          }
        }
      }
    }

    if (found) {
      console.error('[GCLM] 建议在提交前移除调试语句');
    }

  } catch (e) {
    // 忽略错误
  }

  console.log(data);
});
