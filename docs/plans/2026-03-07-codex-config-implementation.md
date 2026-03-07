# Codex Configuration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Git-managed Codex configuration source of truth under `codex/`, then publish it safely to `~/.codex`.

**Architecture:** Keep declarative configuration in-repo and runtime state in `~/.codex`. Use a one-way sync script to copy only managed files. Enforce high-risk guardrails with hooks and keep workflow guidance in `AGENTS.md`.

**Tech Stack:** TOML, Markdown, Python 3, POSIX shell, Git

---

### Task 1: Create the managed Codex directory structure

**Files:**
- Create: `codex/README.md`
- Create: `codex/.gitignore`
- Create: `codex/agents/`
- Create: `codex/hooks/`
- Create: `codex/bin/`

**Step 1: Create directories**

Run: `mkdir -p codex/agents codex/hooks codex/bin`
Expected: directories exist

**Step 2: Add repository metadata files**

Write `codex/README.md` and `codex/.gitignore`.
Expected: source-of-truth layout is documented

**Step 3: Verify structure**

Run: `find codex -maxdepth 2 -type f | sort`
Expected: README and base files appear under `codex/`

### Task 2: Write the global configuration

**Files:**
- Create: `codex/config.toml`

**Step 1: Register the local provider**

Add `model_provider`, `model`, and `[model_providers.local]`.
Expected: provider points at `http://localhost:8319/v1`

**Step 2: Add safety defaults**

Add `approval_policy`, `sandbox_mode`, and history persistence settings.
Expected: fork-first defaults are explicit

**Step 3: Add profiles and agent registry**

Define `fast`, `deep`, and `review`, plus `[agents.*]` entries.
Expected: runtime can select specialized profiles and roles

**Step 4: Verify TOML parses**

Run: `python3 - <<'PY'
import tomllib, pathlib
print(tomllib.loads(pathlib.Path('codex/config.toml').read_text()))
PY`
Expected: parse succeeds

### Task 3: Write the global guidance document

**Files:**
- Create: `codex/AGENTS.md`

**Step 1: Convert the long-form policy into Codex-friendly guidance**

Keep principles and remove executable enforcement details.
Expected: guidance is concise and stable

**Step 2: Verify overlap with hooks is minimal**

Review `codex/AGENTS.md` against planned hooks.
Expected: no major duplicate enforcement logic

### Task 4: Convert gclm-flow roles into Codex agent configs

**Files:**
- Create: `codex/agents/planner.toml`
- Create: `codex/agents/investigator.toml`
- Create: `codex/agents/builder.toml`
- Create: `codex/agents/reviewer.toml`
- Create: `codex/agents/recorder.toml`
- Create: `codex/agents/remember.toml`

**Step 1: Map each existing role to a Codex-native config**

Use the existing markdown role responsibilities as source material.
Expected: each role has focused instructions and sensible runtime defaults

**Step 2: Verify every role file parses as TOML**

Run: `python3 - <<'PY'
import tomllib, pathlib
for path in pathlib.Path('codex/agents').glob('*.toml'):
    tomllib.loads(path.read_text())
    print('OK', path)
PY`
Expected: all role files parse successfully

### Task 5: Implement mixed-mode hooks

**Files:**
- Create: `codex/hooks/pre_tool_risk_guard.py`
- Create: `codex/hooks/pre_tool_tmux_advice.py`
- Create: `codex/hooks/post_tool_git_push_hint.py`
- Create: `codex/hooks/session_start_context.py`
- Create: `codex/hooks/stop_self_check.py`

**Step 1: Implement one responsibility per script**

Expected: blocking logic only in `pre_tool_risk_guard.py`

**Step 2: Make scripts executable**

Run: `chmod +x codex/hooks/*.py`
Expected: scripts are runnable

**Step 3: Verify Python syntax**

Run: `python3 -m py_compile codex/hooks/*.py`
Expected: no syntax errors

### Task 6: Add publish and diff scripts

**Files:**
- Create: `codex/bin/diff-home.sh`
- Create: `codex/bin/sync-to-home.sh`

**Step 1: Implement diff script**

Expected: compare managed files with `~/.codex`

**Step 2: Implement sync script with backups**

Expected: publish only managed files and back up overwritten targets

**Step 3: Verify shell syntax**

Run: `bash -n codex/bin/diff-home.sh codex/bin/sync-to-home.sh`
Expected: no syntax errors

### Task 7: Validate the full managed tree

**Files:**
- Verify: `codex/**`

**Step 1: Run config and script verification**

Run all parse and syntax checks.
Expected: all pass

**Step 2: Inspect the Git diff**

Run: `git status --short && git diff --stat`
Expected: only intended new files are present

### Task 8: Commit and push the codex branch

**Files:**
- Commit all managed Codex files

**Step 1: Stage the new files**

Run: `git add codex docs/plans`
Expected: only intended files staged

**Step 2: Commit**

Run: `git commit -m "feat(codex): add managed codex configuration source"`
Expected: commit succeeds

**Step 3: Push**

Run: `git push -u origin codex`
Expected: remote `codex` branch updated

## Runtime helpers

- `codex/bin/serve-local.sh` starts `codex serve` with stable local defaults.
- `codex/bin/github-webhook-local.sh` starts `codex github` with env-driven defaults.
- `codex/bin/smoke-test-hooks.sh` runs a reproducible real hook smoke test.
