/**
 * Session End Hook
 * 持久化会话状态
 */

let data = '';
process.stdin.on('data', chunk => data += chunk);
process.stdin.on('end', () => {
  try {
    const fs = require('fs');
    const path = require('path');
    const os = require('os');

    const gclmDir = path.join(os.homedir(), '.gclm-flow');
    const sessionFile = path.join(gclmDir, 'session.json');

    if (fs.existsSync(sessionFile)) {
      const session = JSON.parse(fs.readFileSync(sessionFile, 'utf-8'));
      session.endedAt = new Date().toISOString();

      // 计算会话时长
      const start = new Date(session.startedAt);
      const end = new Date(session.endedAt);
      const duration = Math.round((end - start) / 1000 / 60); // 分钟

      session.durationMinutes = duration;
      fs.writeFileSync(sessionFile, JSON.stringify(session, null, 2));

      console.error(`[GCLM] 会话已保存 (时长: ${duration} 分钟)`);
    }

  } catch (e) {
    console.error(`[GCLM] Session end error: ${e.message}`);
  }

  console.log(data);
});
