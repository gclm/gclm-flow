#!/usr/bin/env node

import fs from 'fs-extra';
import os from 'os';
import path from 'path';
import { glob } from 'glob';

function parseArgs(argv) {
  const args = {
    platform: 'codex-cli',
    dryRun: false,
    prune: false,
    agentsSource: null,
    claudeSource: null
  };

  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    if (arg === '--dry-run') args.dryRun = true;
    else if (arg === '--prune') args.prune = true;
    else if (arg === '--platform' || arg === '-p') args.platform = argv[i + 1];
    else if (arg.startsWith('--platform=')) args.platform = arg.split('=')[1];
    else if (arg === '--agents-source') args.agentsSource = argv[i + 1];
    else if (arg.startsWith('--agents-source=')) args.agentsSource = arg.split('=')[1];
    else if (arg === '--claude-source') args.claudeSource = argv[i + 1];
    else if (arg.startsWith('--claude-source=')) args.claudeSource = arg.split('=')[1];
  }

  return args;
}

function usage() {
  console.log('Usage: node scripts/sync-local-to-repo.js [--platform codex-cli|claude-code|all] [--dry-run] [--prune] [--agents-source <path>] [--claude-source <path>]');
}

async function fileDiffers(source, target) {
  if (!await fs.pathExists(target)) return true;
  const [a, b] = await Promise.all([fs.readFile(source), fs.readFile(target)]);
  return !a.equals(b);
}

async function ensureCopy(source, target, dryRun, changed, actionLabel) {
  if (!await fileDiffers(source, target)) return;
  if (dryRun) {
    changed.push(`[DRY-RUN] ${actionLabel}: ${source} -> ${target}`);
    return;
  }
  await fs.ensureDir(path.dirname(target));
  await fs.copy(source, target);
  changed.push(`${actionLabel}: ${source} -> ${target}`);
}

async function collectBestSkillFiles(sourceDirs) {
  const selected = new Map();

  for (const sourceDir of sourceDirs) {
    if (!await fs.pathExists(sourceDir)) continue;
    const files = await glob('**/*.md', { cwd: sourceDir, nodir: true });
    for (const relPath of files) {
      const absPath = path.join(sourceDir, relPath);
      const stat = await fs.stat(absPath);
      const existing = selected.get(relPath);
      if (!existing || stat.mtimeMs > existing.mtimeMs) {
        selected.set(relPath, { absPath, mtimeMs: stat.mtimeMs });
      }
    }
  }

  return selected;
}

async function syncSkills(sourceDirs, targetDir, dryRun, prune, changed) {
  const selected = await collectBestSkillFiles(sourceDirs);

  for (const [relPath, meta] of selected.entries()) {
    const targetPath = path.join(targetDir, relPath);
    await ensureCopy(meta.absPath, targetPath, dryRun, changed, 'SYNC');
  }

  if (!prune) return;
  const targetFiles = await glob('**/*.md', { cwd: targetDir, nodir: true });
  for (const relPath of targetFiles) {
    if (selected.has(relPath)) continue;
    const targetPath = path.join(targetDir, relPath);
    if (dryRun) changed.push(`[DRY-RUN] DELETE: ${targetPath}`);
    else {
      await fs.remove(targetPath);
      changed.push(`DELETE: ${targetPath}`);
    }
  }
}

async function syncConfigFile(sourceFile, targetFile, dryRun, changed) {
  if (!await fs.pathExists(sourceFile)) return;
  await ensureCopy(sourceFile, targetFile, dryRun, changed, 'SYNC');
}

async function run() {
  const args = parseArgs(process.argv.slice(2));
  const validPlatforms = new Set(['codex-cli', 'claude-code', 'all']);
  if (!validPlatforms.has(args.platform)) {
    usage();
    process.exitCode = 1;
    return;
  }

  const repoRoot = process.cwd();
  const home = os.homedir();
  const changed = [];

  if (args.platform === 'codex-cli' || args.platform === 'all') {
    await syncSkills(
      [path.join(home, '.codex', 'skills'), path.join(home, '.agents', 'skills')],
      path.join(repoRoot, 'skills'),
      args.dryRun,
      args.prune,
      changed
    );
    await syncConfigFile(
      args.agentsSource || path.join(home, '.codex', 'AGENTS.md'),
      path.join(repoRoot, 'templates', 'AGENTS-template.md'),
      args.dryRun,
      changed
    );
  }

  if (args.platform === 'claude-code' || args.platform === 'all') {
    await syncSkills(
      [path.join(home, '.claude', 'skills')],
      path.join(repoRoot, 'skills'),
      args.dryRun,
      args.prune,
      changed
    );
    await syncConfigFile(
      args.claudeSource || path.join(home, '.claude', 'CLAUDE.md'),
      path.join(repoRoot, 'templates', 'CLAUDE-template.md'),
      args.dryRun,
      changed
    );
  }

  if (changed.length === 0) {
    console.log('No changes to sync.');
    return;
  }

  console.log(`Synced ${changed.length} item(s):`);
  changed.forEach(item => console.log(`- ${item}`));
  console.log('\nNext: run `git status` to review changes.');
}

run().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
