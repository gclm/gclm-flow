/**
 * Pre Compact Hook
 * 上下文压缩前保存状态
 */

let data = '';
process.stdin.on('data', chunk => data += chunk);
process.stdin.on('end', () => {
  try {
    const fs = require('fs');
    const path = require('path');
    const os = require('os');

    const gclmDir = path.join(os.homedir(), '.gclm-flow');
    const compactFile = path.join(gclmDir, 'compact-history.json');

    // 读取历史
    let history = [];
    if (fs.existsSync(compactFile)) {
      history = JSON.parse(fs.readFileSync(compactFile, 'utf-8'));
    }

    // 添加新记录
    history.push({
      timestamp: new Date().toISOString(),
      project: process.cwd()
    });

    // 只保留最近 50 条
    if (history.length > 50) {
      history = history.slice(-50);
    }

    fs.writeFileSync(compactFile, JSON.stringify(history, null, 2));
    console.error('[GCLM] 上下文压缩前已保存状态');

  } catch (e) {
    console.error(`[GCLM] Pre-compact error: ${e.message}`);
  }

  console.log(data);
});
