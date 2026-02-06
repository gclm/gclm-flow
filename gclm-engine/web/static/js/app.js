// API base URL
const API_BASE = '/api';
const WS_BASE = 'ws://' + window.location.host;

// ============================================================================
// State
// ============================================================================

let currentView = 'dashboard';
let allTasks = [];
let allWorkflows = [];
let wsConnection = null;
let autoRefreshInterval = null;
let currentWorkflowName = null;
let currentTaskPhases = null;
let currentTaskId = null;
let currentTaskWorkflowType = null;
let phaseViewMode = 'graph'; // 'graph' or 'list'

// ============================================================================
// API Functions
// ============================================================================

async function apiGet(url) {
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
    }
    return response.json();
}

async function apiPost(url, data) {
    const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
    }
    return response.json();
}

// Get all tasks
async function getTasks() {
    return apiGet(`${API_BASE}/tasks`);
}

// Get task details
async function getTask(taskId) {
    return apiGet(`${API_BASE}/tasks/${taskId}`);
}

// Get task phases
async function getTaskPhases(taskId) {
    return apiGet(`${API_BASE}/tasks/${taskId}/phases`);
}

// Get task events
async function getTaskEvents(taskId) {
    return apiGet(`${API_BASE}/tasks/${taskId}/events`);
}

// Pause task
async function pauseTask(taskId) {
    return apiPost(`${API_BASE}/tasks/${taskId}/pause`, {});
}

// Resume task
async function resumeTask(taskId) {
    return apiPost(`${API_BASE}/tasks/${taskId}/resume`, {});
}

// Cancel task
async function cancelTask(taskId) {
    return apiPost(`${API_BASE}/tasks/${taskId}/cancel`, {});
}

// Get workflows
async function getWorkflows() {
    return apiGet(`${API_BASE}/workflows`);
}

// Get workflow detail
async function getWorkflow(name) {
    return apiGet(`${API_BASE}/workflows/${name}`);
}

// Get workflow by type
async function getWorkflowByType(workflowType) {
    return apiGet(`${API_BASE}/workflows/type/${workflowType}`);
}

// Get workflow YAML
async function getWorkflowYAML(name) {
    const response = await fetch(`${API_BASE}/workflows/${name}/yaml`);
    if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
    }
    return response.text();
}

// ============================================================================
// Navigation
// ============================================================================

function initNavigation() {
    const navItems = document.querySelectorAll('.nav-item');

    navItems.forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            const view = item.dataset.view;
            switchView(view);
        });
    });

    // Handle browser back/forward
    window.addEventListener('popstate', () => {
        const hash = window.location.hash.slice(1) || 'dashboard';
        switchView(hash, false);
    });

    // Initial view from hash
    const initialView = window.location.hash.slice(1) || 'dashboard';
    switchView(initialView, false);
}

function switchView(viewName, updateHash = true) {
    currentView = viewName;

    // Update nav active state
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.toggle('active', item.dataset.view === viewName);
    });

    // Update views
    document.querySelectorAll('.view').forEach(view => {
        view.classList.toggle('active', view.id === viewName + 'View');
    });

    // Update URL hash
    if (updateHash) {
        history.pushState({ view: viewName }, '', `#${viewName}`);
    }

    // Load view data
    switch (viewName) {
        case 'dashboard':
            loadDashboard();
            break;
        case 'tasks':
            loadTasksView();
            break;
        case 'workflows':
            loadWorkflowsView();
            break;
    }
}

// ============================================================================
// Dashboard View
// ============================================================================

async function loadDashboard() {
    try {
        const data = await getTasks();
        const tasks = data.tasks || [];

        // Update stats
        updateStats(tasks);

        // Update recent activity
        updateRecentActivity(tasks);

        // Update system info
        document.getElementById('systemVersion').textContent = 'v0.2.0';
        document.getElementById('serverAddr').textContent = window.location.origin;
        document.getElementById('systemPort').textContent = '9988';

    } catch (error) {
        console.error('Failed to load dashboard:', error);
    }
}

function updateStats(tasks) {
    const stats = {
        total: tasks.length,
        running: 0,
        completed: 0,
        failed: 0
    };

    tasks.forEach(task => {
        switch (task.status) {
            case 'running':
                stats.running++;
                break;
            case 'completed':
                stats.completed++;
                break;
            case 'failed':
                stats.failed++;
                break;
        }
    });

    document.getElementById('statTotal').textContent = stats.total;
    document.getElementById('statRunning').textContent = stats.running;
    document.getElementById('statCompleted').textContent = stats.completed;
    document.getElementById('statFailed').textContent = stats.failed;
}

function updateRecentActivity(tasks) {
    const container = document.getElementById('recentActivity');

    // Sort by createdAt descending, get top 5
    const recentTasks = tasks
        .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
        .slice(0, 5);

    if (recentTasks.length === 0) {
        container.innerHTML = '<p class="empty">暂无活动</p>';
        return;
    }

    container.innerHTML = recentTasks.map(task => {
        const icon = getActivityIcon(task.status);
        const time = formatRelativeTime(task.createdAt);
        return `
            <div class="activity-item">
                <div class="activity-icon">${icon}</div>
                <div class="activity-content">
                    <div class="activity-title">${escapeHtml(task.prompt.substring(0, 60))}${task.prompt.length > 60 ? '...' : ''}</div>
                    <div class="activity-time">${time}</div>
                </div>
            </div>
        `;
    }).join('');
}

function getActivityIcon(status) {
    const icons = {
        'created': '📝',
        'running': '🔄',
        'paused': '⏸️',
        'completed': '✅',
        'failed': '❌',
        'cancelled': '🚫'
    };
    return icons[status] || '📋';
}

function refreshDashboard() {
    loadDashboard();
    showToast('仪表板已刷新', 'success');
}

// ============================================================================
// Tasks View
// ============================================================================

async function loadTasksView() {
    await loadTasks();

    // Setup filter listeners
    document.getElementById('statusFilter').addEventListener('change', loadTasks);
    document.getElementById('searchInput').addEventListener('input', debounce(loadTasks, 300));
}

async function loadTasks() {
    const container = document.getElementById('tasksList');
    container.innerHTML = '<p class="loading">加载中...</p>';

    try {
        const data = await getTasks();
        allTasks = data.tasks || [];

        const statusFilter = document.getElementById('statusFilter').value;
        const searchQuery = document.getElementById('searchInput').value.toLowerCase();

        // Filter tasks
        let filteredTasks = allTasks;

        if (statusFilter) {
            filteredTasks = filteredTasks.filter(t => t.status === statusFilter);
        }

        if (searchQuery) {
            filteredTasks = filteredTasks.filter(t =>
                t.prompt.toLowerCase().includes(searchQuery) ||
                t.id.toLowerCase().includes(searchQuery)
            );
        }

        if (filteredTasks.length === 0) {
            container.innerHTML = '<p class="empty">没有找到匹配的任务</p>';
            return;
        }

        container.innerHTML = filteredTasks.map(task => renderTaskItem(task)).join('');
    } catch (error) {
        container.innerHTML = `<p class="error">加载失败: ${error.message}</p>`;
    }
}

function renderTaskItem(task) {
    const statusClass = `status-${task.status}`;
    const statusLabel = getStatusLabel(task.status);
    const time = formatRelativeTime(task.createdAt);

    // Create phase progress indicator
    const phaseProgress = renderPhaseProgress(task.currentPhase, task.totalPhases, task.status);

    return `
        <div class="task-item" onclick="openTaskModal('${task.id}')">
            <div class="task-header-row">
                <div class="task-title-row">
                    <span class="status-badge ${statusClass}">${statusLabel}</span>
                    <span class="task-title">${escapeHtml(task.prompt)}</span>
                </div>
                <span class="task-time">${time}</span>
            </div>

            <!-- Workflow Info Row -->
            <div class="task-workflow-row">
                <span class="task-workflow-type">🔧 ${escapeHtml(task.workflowType)}</span>
                <span class="task-phase-progress">
                    ${phaseProgress}
                    <span class="phase-text">${task.currentPhase}/${task.totalPhases}</span>
                </span>
            </div>

            ${task.status === 'running' || task.status === 'paused' ? `
                <div class="task-actions-row" onclick="event.stopPropagation()">
                    ${task.status === 'running' ? `<button onclick="pauseTaskUI('${task.id}')">暂停</button>` : ''}
                    ${task.status === 'paused' ? `<button onclick="resumeTaskUI('${task.id}')">恢复</button>` : ''}
                    <button onclick="cancelTaskUI('${task.id}')">取消</button>
                </div>
            ` : ''}
        </div>
    `;
}

// Render compact phase progress indicator
function renderPhaseProgress(current, total, status) {
    const dots = [];
    for (let i = 0; i < total; i++) {
        let dotClass = 'phase-dot-pending';
        if (i < current) {
            dotClass = 'phase-dot-completed';
        } else if (i === current && status === 'running') {
            dotClass = 'phase-dot-running';
        } else if (i === current && status === 'failed') {
            dotClass = 'phase-dot-failed';
        }
        dots.push(`<span class="phase-dot ${dotClass}"></span>`);
    }
    return dots.join('');
}

function refreshTasks() {
    loadTasks();
    showToast('任务列表已刷新', 'success');
}

// ============================================================================
// Workflows View
// ============================================================================

async function loadWorkflowsView() {
    const container = document.getElementById('workflowsList');
    container.innerHTML = '<p class="loading">加载中...</p>';

    try {
        const data = await getWorkflows();
        const workflows = data.workflows || [];

        if (workflows.length === 0) {
            container.innerHTML = '<p class="empty">暂无工作流</p>';
            return;
        }

        // Get detailed workflow info including nodes
        const workflowDetails = await Promise.all(
            workflows.map(wf => getWorkflow(wf.name))
        );

        container.innerHTML = `
            <div class="workflows-list">
                ${workflowDetails.map(d => renderWorkflowCard(d.workflow)).join('')}
            </div>
        `;
    } catch (error) {
        container.innerHTML = `<p class="error">加载失败: ${error.message}</p>`;
    }
}

// Render Agent Distribution Overview - REMOVED (redundant with visual cards)
function renderAgentDistribution(workflows) {
    // Agent distribution is now shown on each card visually
    // No need for separate table
    return '';
}

// Count agents in workflow nodes
function countAgents(nodes) {
    const counts = {};
    nodes.forEach(node => {
        const agent = node.agent || 'unknown';
        counts[agent] = (counts[agent] || 0) + 1;
    });
    return counts;
}

// Get agent icon/emoji
function getAgentIcon(agent) {
    const icons = {
        'investigator': '🔍',
        'architect': '🏗️',
        'worker': '🔧',
        'tdd-guide': '🧪',
        'code-simplifier': '✨',
        'security-guidance': '🛡️',
        'code-reviewer': '👀',
        'llmdoc': '📚',
        'spec-guide': '📋',
        'unknown': '❓'
    };
    return icons[agent] || '🤖';
}

function renderWorkflowCard(workflow) {
    const nodes = workflow.nodes || [];
    const agentCounts = countAgents(nodes);
    const totalPhases = nodes.length;

    return `
        <div class="workflow-card" onclick="openWorkflowModal('${workflow.name}')">
            <!-- Header: Name and Type -->
            <div class="workflow-card-header">
                <div class="workflow-name">${escapeHtml(workflow.displayName || workflow.name)}</div>
                <div class="workflow-type">${escapeHtml(workflow.workflowType || workflow.name)}</div>
            </div>

            <!-- Centerpiece: Workflow Graph -->
            <div class="workflow-graph-container">
                ${renderWorkflowGraph(workflow, 'small')}
            </div>

            <!-- Agent Badges -->
            <div class="workflow-agents">
                ${Object.entries(agentCounts).map(([agent, count]) => `
                    <span class="agent-badge" title="${agent} x ${count}">
                        ${getAgentIcon(agent)} ${count}
                    </span>
                `).join('')}
                <span class="agent-total">${totalPhases} 阶段</span>
            </div>

            <!-- Footer: Meta info -->
            <div class="workflow-meta">
                <span>v${escapeHtml(workflow.version || '-')}</span>
                ${workflow.isBuiltin ? '<span class="builtin-badge">内置</span>' : ''}
            </div>
        </div>
    `;
}

// ============================================================================
// Task Modal
// ============================================================================

async function openTaskModal(taskId) {
    const modal = document.getElementById('taskModal');
    const title = document.getElementById('modalTaskTitle');
    const body = document.getElementById('modalTaskBody');
    const actions = document.getElementById('taskModalActions');

    currentTaskId = taskId;
    phaseViewMode = 'graph'; // Reset to graph view

    title.textContent = `任务详情 (${taskId.substring(0, 8)})`;
    body.innerHTML = '<p class="loading">加载中...</p>';
    actions.innerHTML = ''; // Clear actions
    modal.classList.add('show');

    try {
        const [taskResp, phasesResp, eventsResp] = await Promise.all([
            getTask(taskId),
            getTaskPhases(taskId),
            getTaskEvents(taskId)
        ]);

        // Extract data from API responses
        const task = taskResp.task || {};
        const phases = phasesResp.phases || [];
        const events = eventsResp.events || [];

        currentTaskPhases = phases;
        currentTaskWorkflowType = task.workflowType;

        // Render actions in header
        renderTaskActions(task, actions);

        // Render detail (async for graph loading)
        await renderTaskDetailAsync(task, phases, events, body);
    } catch (error) {
        body.innerHTML = `<p class="error">加载失败: ${error.message}</p>`;
    }
}

// Async version of renderTaskDetail that loads workflow graph
async function renderTaskDetailAsync(task, phases, events, container) {
    const statusClass = `status-${task.status}`;
    const statusLabel = getStatusLabel(task.status);

    // Generate the phase graph HTML (async)
    let phaseGraphHtml = '';
    if (phases.length > 0) {
        try {
            const graphSvg = await renderTaskPhasesGraph(phases, task.workflowType);
            phaseGraphHtml = `
                <div class="workflow-graph-container" style="margin-bottom: 20px;">
                    ${graphSvg}
                </div>
            `;
        } catch (error) {
            console.error('Failed to render phase graph:', error);
            phaseGraphHtml = '<p class="error">图示加载失败</p>';
        }
    }

    container.innerHTML = `
        <div class="detail-section">
            <h4>基本信息</h4>
            <div class="detail-grid">
                <div class="detail-item">
                    <span class="detail-label">任务 ID</span>
                    <span class="detail-value">${escapeHtml(task.id || '-')}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">状态</span>
                    <span class="detail-value"><span class="status-badge ${statusClass}">${statusLabel}</span></span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">工作流</span>
                    <span class="detail-value">${escapeHtml(task.workflowType || '-')}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">进度</span>
                    <span class="detail-value">${task.currentPhase}/${task.totalPhases}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">创建时间</span>
                    <span class="detail-value">${formatDateTime(task.createdAt)}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">更新时间</span>
                    <span class="detail-value">${formatDateTime(task.updatedAt)}</span>
                </div>
            </div>
            <div style="margin-top: 12px;">
                <span class="detail-label">任务描述</span>
                <p style="margin-top: 4px; color: var(--text-primary);">${escapeHtml(task.prompt || '')}</p>
            </div>
        </div>

        <div class="detail-section">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
                <h4>阶段列表 (${phases.length})</h4>
                ${phases.length > 0 ? `
                    <div class="view-toggle">
                        <button class="toggle-btn active" onclick="switchPhaseView('graph')">图示</button>
                        <button class="toggle-btn" onclick="switchPhaseView('list')">列表</button>
                    </div>
                ` : ''}
            </div>
            ${phases.length > 0 ? `
                <div id="phaseViewContainer">
                    ${phaseGraphHtml}
                </div>
            ` : '<p class="empty">暂无阶段</p>'}
        </div>

        <div class="detail-section">
            <h4>事件日志 (${events.length})</h4>
            ${events.length > 0 ? `
                <div class="event-list">
                    ${events.map(e => renderEventItem(e)).join('')}
                </div>
            ` : '<p class="empty">暂无事件</p>'}
        </div>
    `;
}

// Render task action buttons in modal header
function renderTaskActions(task, container) {
    if (task.status === 'running' || task.status === 'paused') {
        let buttons = '';
        if (task.status === 'running') {
            buttons += `<button class="btn-small" onclick="pauseTaskUI('${task.id}')">暂停</button>`;
        }
        if (task.status === 'paused') {
            buttons += `<button class="btn-small" onclick="resumeTaskUI('${task.id}')">恢复</button>`;
        }
        buttons += `<button class="btn-small" onclick="cancelTaskUI('${task.id}')" style="color: var(--error-color);">取消</button>`;
        container.innerHTML = buttons;
    } else {
        container.innerHTML = '';
    }
}

function renderTaskDetail(task, phases, events) {
    const statusClass = `status-${task.status}`;
    const statusLabel = getStatusLabel(task.status);

    return `
        <div class="detail-section">
            <h4>基本信息</h4>
            <div class="detail-grid">
                <div class="detail-item">
                    <span class="detail-label">任务 ID</span>
                    <span class="detail-value">${escapeHtml(task.id || '-')}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">状态</span>
                    <span class="detail-value"><span class="status-badge ${statusClass}">${statusLabel}</span></span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">工作流</span>
                    <span class="detail-value">${escapeHtml(task.workflowType || '-')}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">进度</span>
                    <span class="detail-value">${task.currentPhase}/${task.totalPhases}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">创建时间</span>
                    <span class="detail-value">${formatDateTime(task.createdAt)}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">更新时间</span>
                    <span class="detail-value">${formatDateTime(task.updatedAt)}</span>
                </div>
            </div>
            <div style="margin-top: 12px;">
                <span class="detail-label">任务描述</span>
                <p style="margin-top: 4px; color: var(--text-primary);">${escapeHtml(task.prompt || '')}</p>
            </div>
        </div>

        <div class="detail-section">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
                <h4>阶段列表 (${phases.length})</h4>
                ${phases.length > 0 ? `
                    <div class="view-toggle">
                        <button class="toggle-btn ${phaseViewMode === 'graph' ? 'active' : ''}" onclick="switchPhaseView('graph')">图示</button>
                        <button class="toggle-btn ${phaseViewMode === 'list' ? 'active' : ''}" onclick="switchPhaseView('list')">列表</button>
                    </div>
                ` : ''}
            </div>
            ${phases.length > 0 ? `
                <div id="phaseViewContainer">
                    ${phaseViewMode === 'graph' ? renderTaskPhasesGraphView(phases) : renderTaskPhasesListView(phases)}
                </div>
            ` : '<p class="empty">暂无阶段</p>'}
        </div>

        <div class="detail-section">
            <h4>事件日志 (${events.length})</h4>
            ${events.length > 0 ? `
                <div class="event-list">
                    ${events.map(e => renderEventItem(e)).join('')}
                </div>
            ` : '<p class="empty">暂无事件</p>'}
        </div>
    `;
}

// Switch between graph and list view for phases
async function switchPhaseView(mode) {
    phaseViewMode = mode;
    const container = document.getElementById('phaseViewContainer');
    if (container && currentTaskPhases) {
        if (mode === 'graph') {
            container.innerHTML = '<p class="loading">加载中...</p>';
            container.innerHTML = await renderTaskPhasesGraphView(currentTaskPhases);
        } else {
            container.innerHTML = renderTaskPhasesListView(currentTaskPhases);
        }

        // Update toggle button styles
        document.querySelectorAll('.toggle-btn').forEach(btn => {
            btn.classList.toggle('active', btn.textContent === (mode === 'graph' ? '图示' : '列表'));
        });
    }
}

// Render phases as graph view
async function renderTaskPhasesGraphView(phases) {
    const graphHtml = await renderTaskPhasesGraph(phases, currentTaskWorkflowType);
    return `
        <div class="workflow-graph-container" style="margin-bottom: 20px;">
            ${graphHtml}
        </div>
    `;
}

// Render phases as list view
function renderTaskPhasesListView(phases) {
    return `
        <div class="phase-list">
            ${phases.map(p => renderPhaseItem(p)).join('')}
        </div>
    `;
}

function renderPhaseItem(phase) {
    const statusClass = `status-${phase.status}`;
    const statusLabel = getStatusLabel(phase.status);

    return `
        <div class="phase-list-item">
            <div class="phase-list-header">
                <span class="phase-list-name">${escapeHtml(phase.displayName)}</span>
                <span class="status-badge ${statusClass}">${statusLabel}</span>
            </div>
            <div class="phase-list-meta">
                <span>${getAgentIcon(phase.agentName)} ${escapeHtml(phase.agentName)}</span>
                <span>Model: ${escapeHtml(phase.modelName)}</span>
            </div>
            ${phase.outputText ? `<div class="phase-output">${escapeHtml(phase.outputText)}</div>` : ''}
            ${phase.error ? `<div class="phase-error">${escapeHtml(phase.error)}</div>` : ''}
        </div>
    `;
}

function renderEventItem(event) {
    return `
        <div class="event-list-item">
            <span class="event-time">${formatDateTime(event.occurredAt)}</span>
            <span class="event-type">${escapeHtml(event.eventType)}</span>
            <span class="event-data">${event.data ? escapeHtml(event.data) : ''}</span>
        </div>
    `;
}

function closeModal() {
    document.getElementById('taskModal').classList.remove('show');
    document.getElementById('taskModalActions').innerHTML = '';
    currentTaskId = null;
    currentTaskPhases = null;
    currentTaskWorkflowType = null;
    phaseViewMode = 'graph';
}

// ============================================================================
// Workflow Modal
// ============================================================================

async function openWorkflowModal(workflowName) {
    const modal = document.getElementById('workflowModal');
    const title = document.getElementById('modalWorkflowTitle');
    const body = document.getElementById('modalWorkflowBody');
    const yamlBtn = document.getElementById('viewYamlBtn');

    // Save current workflow name for YAML button
    currentWorkflowName = workflowName;

    title.textContent = `工作流详情 - ${workflowName}`;
    body.innerHTML = '<p class="loading">加载中...</p>';

    // Show YAML button when viewing workflow detail
    yamlBtn.style.display = 'inline-block';

    modal.classList.add('show');

    try {
        const data = await getWorkflow(workflowName);
        const workflow = data.workflow || {};

        body.innerHTML = renderWorkflowDetail(workflow);
    } catch (error) {
        body.innerHTML = `<p class="error">加载失败: ${error.message}</p>`;
    }
}

function renderWorkflowDetail(workflow) {
    const nodes = workflow.nodes || [];
    const agentCounts = countAgents(nodes);

    return `
        <div class="detail-section">
            <h4>基本信息</h4>
            <div class="detail-grid">
                <div class="detail-item">
                    <span class="detail-label">名称</span>
                    <span class="detail-value">${escapeHtml(workflow.name || '-')}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">显示名称</span>
                    <span class="detail-value">${escapeHtml(workflow.displayName || '-')}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">类型</span>
                    <span class="detail-value">${escapeHtml(workflow.workflowType || '-')}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">版本</span>
                    <span class="detail-value">${escapeHtml(workflow.version || '-')}</span>
                </div>
            </div>
            <div style="margin-top: 12px;">
                <span class="detail-label">描述</span>
                <p style="margin-top: 4px; color: var(--text-primary);">${escapeHtml(workflow.description || '暂无描述')}</p>
            </div>
        </div>

        <div class="detail-section">
            <h4>工作流结构图</h4>
            <div class="workflow-graph-container">
                ${renderWorkflowGraph(workflow, 'normal')}
            </div>
        </div>

        <div class="detail-section">
            <h4>Agent 分布</h4>
            <div class="agent-distribution-detail">
                ${Object.entries(agentCounts).map(([agent, count]) => `
                    <div class="agent-dist-item">
                        <span class="agent-icon">${getAgentIcon(agent)}</span>
                        <span class="agent-name">${escapeHtml(agent)}</span>
                        <span class="agent-count">${count} 个阶段</span>
                    </div>
                `).join('')}
            </div>
        </div>

        <div class="detail-section">
            <h4>节点列表 (${nodes.length})</h4>
            ${nodes.length > 0 ? `
                <div class="phase-list">
                    ${nodes.map(node => `
                        <div class="phase-list-item">
                            <div class="phase-list-header">
                                <span class="phase-list-name">${escapeHtml(node.displayName || node.ref || 'Unknown')}</span>
                                ${node.required ? '<span class="status-badge status-completed">必需</span>' : '<span class="status-badge status-cancelled">可选</span>'}
                            </div>
                            <div class="phase-list-meta">
                                <span>${getAgentIcon(node.agent)} ${escapeHtml(node.agent)}</span>
                                <span>Model: ${escapeHtml(node.model)}</span>
                                <span>超时: ${node.timeout}s</span>
                            </div>
                            ${node.dependsOn && node.dependsOn.length > 0 ? `
                                <div style="margin-top: 8px; font-size: 0.8125rem; color: var(--text-secondary);">
                                    依赖: ${Array.isArray(node.dependsOn) ? node.dependsOn.map(escapeHtml).join(', ') : escapeHtml(String(node.dependsOn))}
                                </div>
                            ` : ''}
                        </div>
                    `).join('')}
                </div>
            ` : '<p class="empty">暂无节点</p>'}
        </div>
    `;
}

async function viewWorkflowYAML(workflowName) {
    try {
        const yaml = await getWorkflowYAML(workflowName);

        // Create YAML modal content
        const modal = document.getElementById('workflowModal');
        const body = document.getElementById('modalWorkflowBody');
        const yamlBtn = document.getElementById('viewYamlBtn');

        // Hide YAML button when viewing YAML
        yamlBtn.style.display = 'none';

        body.innerHTML = `
            <div class="detail-section">
                <h4>YAML 源文件 - ${workflowName}</h4>
                <div style="margin-bottom: 12px;">
                    <button class="btn-small" onclick="closeWorkflowYAML()">返回详情</button>
                    <button class="btn-small" onclick="copyWorkflowYAML()">复制</button>
                </div>
                <pre style="background: #1e293b; color: #e2e8f0; padding: 16px; border-radius: 8px; overflow-x: auto; font-size: 0.8125rem; max-height: 400px; overflow-y: auto;">${escapeHtml(yaml)}</pre>
            </div>
        `;
    } catch (error) {
        showToast(`加载 YAML 失败: ${error.message}`, 'error');
    }
}

function closeWorkflowYAML() {
    // Reload workflow detail
    openWorkflowModal(currentWorkflowName);
}

// View YAML from modal header button
function viewWorkflowYAMLFromModal() {
    if (currentWorkflowName) {
        viewWorkflowYAML(currentWorkflowName);
    }
}

function copyWorkflowYAML() {
    const yamlContent = document.querySelector('#modalWorkflowBody pre').textContent;
    navigator.clipboard.writeText(yamlContent).then(() => {
        showToast('YAML 已复制到剪贴板', 'success');
    }).catch(() => {
        showToast('复制失败', 'error');
    });
}

function closeWorkflowModal() {
    document.getElementById('workflowModal').classList.remove('show');
    // Hide YAML button when modal is closed
    document.getElementById('viewYamlBtn').style.display = 'none';
    currentWorkflowName = null;
}

// ============================================================================
// Task Actions
// ============================================================================

async function pauseTaskUI(taskId) {
    try {
        await pauseTask(taskId);
        showToast('任务已暂停', 'success');
        loadTasks();
        if (currentView === 'dashboard') loadDashboard();
        closeModal();
    } catch (error) {
        showToast(`暂停失败: ${error.message}`, 'error');
    }
}

async function resumeTaskUI(taskId) {
    try {
        await resumeTask(taskId);
        showToast('任务已恢复', 'success');
        loadTasks();
        if (currentView === 'dashboard') loadDashboard();
        closeModal();
    } catch (error) {
        showToast(`恢复失败: ${error.message}`, 'error');
    }
}

async function cancelTaskUI(taskId) {
    if (!confirm('确定要取消此任务吗？')) return;

    try {
        await cancelTask(taskId);
        showToast('任务已取消', 'success');
        loadTasks();
        if (currentView === 'dashboard') loadDashboard();
        closeModal();
    } catch (error) {
        showToast(`取消失败: ${error.message}`, 'error');
    }
}

// ============================================================================
// Toast Notifications
// ============================================================================

function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.textContent = message;
    document.body.appendChild(toast);

    setTimeout(() => toast.classList.add('show'), 10);
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// ============================================================================
// Utility Functions
// ============================================================================

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDateTime(dateStr) {
    return new Date(dateStr).toLocaleString('zh-CN');
}

function formatRelativeTime(dateStr) {
    const date = new Date(dateStr);
    const now = new Date();
    const diff = Math.floor((now - date) / 1000);

    if (diff < 60) return '刚刚';
    if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`;
    if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`;
    if (diff < 604800) return `${Math.floor(diff / 86400)}天前`;
    return formatDateTime(dateStr);
}

function formatUptime() {
    // In a real implementation, this would be calculated from server start time
    return '运行中';
}

function getStatusLabel(status) {
    const labels = {
        'created': '已创建',
        'running': '运行中',
        'paused': '已暂停',
        'completed': '已完成',
        'failed': '失败',
        'cancelled': '已取消'
    };
    return labels[status] || status;
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// ============================================================================
// Workflow Graph Visualization
// ============================================================================

// Render workflow graph as SVG
function renderWorkflowGraph(workflow, size = 'normal') {
    const nodes = workflow.nodes || [];
    if (nodes.length === 0) return '<p class="empty">暂无节点</p>';

    const config = {
        small: { nodeWidth: 80, nodeHeight: 40, fontSize: 10, iconSize: 12, spacingX: 30, spacingY: 20 },
        normal: { nodeWidth: 140, nodeHeight: 60, fontSize: 12, iconSize: 16, spacingX: 50, spacingY: 30 }
    }[size];

    // Calculate node positions using topological sort layering
    const layers = calculateNodeLayers(nodes);
    const nodePositions = calculateNodePositions(layers, config);

    // Build SVG
    const width = calculateGraphWidth(nodePositions, config);
    const height = calculateGraphHeight(layers, config);

    let svg = `<svg class="workflow-graph" width="100%" height="${height}" viewBox="0 0 ${width} ${height}" preserveAspectRatio="xMidYMid meet">`;

    // Add arrowhead marker
    svg += `
        <defs>
            <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
                <polygon points="0 0, 10 3.5, 0 7" fill="#94a3b8" />
            </marker>
        </defs>
    `;

    // Draw edges (connections)
    svg += '<g class="graph-edges">';
    nodes.forEach(node => {
        const fromPos = nodePositions[node.ref];
        if (!fromPos) return;

        (node.dependsOn || []).forEach(depRef => {
            const toPos = nodePositions[depRef];
            if (toPos) {
                svg += drawEdge(fromPos, toPos, config);
            }
        });
    });
    svg += '</g>';

    // Draw nodes
    svg += '<g class="graph-nodes">';
    nodes.forEach(node => {
        const pos = nodePositions[node.ref];
        if (pos) {
            svg += drawNode(node, pos, config, 'pending');
        }
    });
    svg += '</g>';

    svg += '</svg>';
    return svg;
}

// Render task phases graph with status colors
// This function now loads the workflow to get proper structure with dependencies
async function renderTaskPhasesGraph(phases, workflowType) {
    if (!phases || phases.length === 0) return '<p class="empty">暂无阶段</p>';

    try {
        // Fetch workflow by type (not by name)
        const workflowData = await getWorkflowByType(workflowType);
        console.log('Workflow data from API:', workflowData);

        const workflow = workflowData.workflow || {};
        console.log('Extracted workflow:', workflow);

        const workflowNodes = workflow.nodes || [];
        console.log('Workflow nodes:', workflowNodes);

        // Create a map of phase status with multiple matching keys
        const phaseStatusMap = {};
        phases.forEach(p => {
            // Map by phaseName (primary key) - this should match node.ref
            phaseStatusMap[p.phaseName] = p.status;
            console.log(`Phase: phaseName="${p.phaseName}", displayName="${p.displayName}", status="${p.status}"`);
        });

        console.log('Final phase status map:', phaseStatusMap);

        // Use same rendering logic as workflow graph but with status colors
        return renderWorkflowGraphWithStatus(workflow, phaseStatusMap);
    } catch (error) {
        console.error('Failed to load workflow for graph:', error);
        // Fallback to simple horizontal layout
        return renderTaskPhasesGraphFallback(phases);
    }
}

// Render workflow graph with custom status colors for each node
function renderWorkflowGraphWithStatus(workflow, phaseStatusMap) {
    const nodes = workflow.nodes || [];
    if (nodes.length === 0) return '<p class="empty">暂无节点</p>';

    // Use normal size like workflow detail
    const config = { nodeWidth: 140, nodeHeight: 60, fontSize: 12, iconSize: 16, spacingX: 50, spacingY: 30 };

    // Calculate node positions using topological sort layering (same as workflow graph)
    const layers = calculateNodeLayers(nodes);
    const nodePositions = calculateNodePositions(layers, config);

    // Build SVG
    const width = calculateGraphWidth(nodePositions, config);
    const height = calculateGraphHeight(layers, config);

    let svg = `<svg class="workflow-graph" width="100%" viewBox="0 0 ${width} ${height}" style="max-width: 100%; height: auto;">`;

    // Add arrowhead marker
    svg += `
        <defs>
            <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
                <polygon points="0 0, 10 3.5, 0 7" fill="#94a3b8" />
            </marker>
        </defs>
    `;

    // Draw edges (connections)
    svg += '<g class="graph-edges">';
    nodes.forEach(node => {
        const fromPos = nodePositions[node.ref];
        if (!fromPos) return;

        (node.dependsOn || []).forEach(depRef => {
            const toPos = nodePositions[depRef];
            if (toPos) {
                svg += drawEdge(fromPos, toPos, config);
            }
        });
    });
    svg += '</g>';

    // Draw nodes with status colors
    svg += '<g class="graph-nodes">';
    nodes.forEach(node => {
        const pos = nodePositions[node.ref];
        if (pos) {
            // Try to match status by ref first, then by displayName, default to pending
            let status = phaseStatusMap[node.ref];
            if (status === undefined) {
                status = phaseStatusMap[node.displayName];
            }
            if (status === undefined) {
                status = 'pending';
            }
            svg += drawNode(node, pos, config, status);
        }
    });
    svg += '</g>';

    svg += '</svg>';
    return svg;
}

// Fallback: simple horizontal layout (if workflow loading fails)
function renderTaskPhasesGraphFallback(phases) {
    const config = { nodeWidth: 140, nodeHeight: 60, fontSize: 12, iconSize: 16, spacingX: 50, spacingY: 30 };

    const nodePositions = {};
    phases.forEach((phase, i) => {
        nodePositions[phase.phaseName] = {
            x: config.spacingX + i * (config.nodeWidth + config.spacingX),
            y: config.spacingY
        };
    });

    const width = phases.length * (config.nodeWidth + config.spacingX) + config.spacingX;
    const height = config.nodeHeight + 2 * config.spacingY;

    let svg = `<svg class="workflow-graph" width="100%" viewBox="0 0 ${width} ${height}" style="max-width: 100%; height: auto;">`;

    svg += `
        <defs>
            <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
                <polygon points="0 0, 10 3.5, 0 7" fill="#94a3b8" />
            </marker>
        </defs>
    `;

    svg += '<g class="graph-nodes">';
    phases.forEach(phase => {
        const pos = nodePositions[phase.phaseName];
        if (pos) {
            svg += drawNode(
                { ref: phase.phaseName, displayName: phase.displayName, agent: phase.agentName },
                pos,
                config,
                phase.status
            );
        }
    });
    svg += '</g>';

    svg += '</svg>';
    return svg;
}

// Calculate node layers using topological sort
function calculateNodeLayers(nodes) {
    const nodeMap = {};
    nodes.forEach(n => nodeMap[n.ref] = { ...n, layer: 0 });

    const layers = [];
    const visited = new Set();
    const inDegree = {};

    // Calculate in-degree
    nodes.forEach(n => {
        inDegree[n.ref] = (n.dependsOn || []).length;
    });

    // Topological sort with layering
    let remaining = [...nodes];
    let currentLayer = 0;

    while (remaining.length > 0) {
        const layerNodes = remaining.filter(n => (inDegree[n.ref] || 0) === 0);
        if (layerNodes.length === 0) break; // Circular dependency

        layers.push(layerNodes);
        layerNodes.forEach(n => {
            visited.add(n.ref);
            nodeMap[n.ref].layer = currentLayer;
            remaining = remaining.filter(r => r.ref !== n.ref);
        });

        // Update in-degree
        layerNodes.forEach(n => {
            (n.dependsOn || []).forEach(dep => {
                // This doesn't make sense - we need to find nodes that depend on current node
            });
        });

        // Better approach: reduce in-degree of dependent nodes
        nodes.forEach(n => {
            if (layerNodes.some(ln => ln.ref === n.ref)) return;
            const newDeps = (n.dependsOn || []).filter(d => !visited.has(d));
            inDegree[n.ref] = newDeps.length;
        });

        currentLayer++;
    }

    // Handle remaining nodes (circular deps) - put them in last layer
    if (remaining.length > 0) {
        layers.push(remaining);
    }

    return layers;
}

// Calculate positions for nodes based on layers
function calculateNodePositions(layers, config) {
    const positions = {};
    const layerWidths = [];

    // Calculate width needed for each layer
    layers.forEach((layer, i) => {
        layerWidths[i] = layer.length * config.nodeWidth + (layer.length - 1) * config.spacingX;
    });

    const maxLayerWidth = Math.max(...layerWidths, 0);

    layers.forEach((layer, layerIndex) => {
        const layerWidth = layerWidths[layerIndex];
        const startX = (maxLayerWidth - layerWidth) / 2;

        layer.forEach((node, nodeIndex) => {
            positions[node.ref] = {
                x: startX + nodeIndex * (config.nodeWidth + config.spacingX),
                y: layerIndex * (config.nodeHeight + config.spacingY)
            };
        });
    });

    return positions;
}

// Calculate total graph width
function calculateGraphWidth(positions, config) {
    const maxX = Math.max(...Object.values(positions).map(p => p.x), 0);
    return maxX + config.nodeWidth + config.spacingX;
}

// Calculate total graph height
function calculateGraphHeight(layers, config) {
    return layers.length * (config.nodeHeight + config.spacingY) + config.spacingY;
}

// Draw a single node
function drawNode(node, pos, config, status) {
    const colors = {
        'pending': '#94a3b8',
        'running': '#3b82f6',
        'completed': '#10b981',
        'failed': '#ef4444',
        'cancelled': '#6b7280',
        'created': '#94a3b8'
    };
    const strokeColor = colors[status] || colors.pending;
    const fillColor = status === 'running' ? '#dbeafe' : '#ffffff';
    const agentIcon = getAgentIcon(node.agent);

    return `
        <g class="graph-node" transform="translate(${pos.x}, ${pos.y})">
            <rect
                x="0" y="0"
                width="${config.nodeWidth}" height="${config.nodeHeight}"
                rx="6" ry="6"
                fill="${fillColor}"
                stroke="${strokeColor}"
                stroke-width="2"
            />
            <text x="${config.nodeWidth / 2}" y="${config.fontSize + 8}"
                text-anchor="middle"
                font-size="${config.fontSize}"
                font-weight="500"
                fill="#1e293b">
                ${escapeHtml(node.displayName || node.ref).substring(0, 12)}
            </text>
            <text x="${config.nodeWidth / 2}" y="${config.nodeHeight - 10}"
                text-anchor="middle"
                font-size="${config.fontSize - 2}"
                fill="#64748b">
                ${agentIcon} ${escapeHtml(node.agent || '').substring(0, 10)}
            </text>
        </g>
    `;
}

// Draw an edge between two nodes
function drawEdge(from, to, config) {
    const fromX = from.x + config.nodeWidth;
    const fromY = from.y + config.nodeHeight / 2;
    const toX = to.x;
    const toY = to.y + config.nodeHeight / 2;

    // Draw curved path
    const midX = (fromX + toX) / 2;
    const path = `M ${fromX} ${fromY} C ${midX} ${fromY}, ${midX} ${toY}, ${toX} ${toY}`;

    return `
        <path d="${path}"
            fill="none"
            stroke="#94a3b8"
            stroke-width="1.5"
            marker-end="url(#arrowhead)"
        />
    `;
}

// ============================================================================
// Auto-refresh
// ============================================================================

function startAutoRefresh() {
    if (autoRefreshInterval) clearInterval(autoRefreshInterval);

    autoRefreshInterval = setInterval(() => {
        switch (currentView) {
            case 'dashboard':
                loadDashboard();
                break;
            case 'tasks':
                loadTasks();
                break;
        }
    }, 30000); // Refresh every 30 seconds
}

// ============================================================================
// Initialize
// ============================================================================

document.addEventListener('DOMContentLoaded', () => {
    initNavigation();
    startAutoRefresh();

    // Close modals on outside click
    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal')) {
            e.target.classList.remove('show');
        }
    });

    // Update server status
    const serverStatus = document.getElementById('serverStatus');
    serverStatus.querySelector('.status-text').textContent = '已连接';
});
