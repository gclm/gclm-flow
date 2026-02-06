package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gclm/gclm-flow/gclm-engine/internal/api/websocket"
	"github.com/gclm/gclm-flow/gclm-engine/internal/logger"
)

// ============================================================================
// Task Handlers
// ============================================================================

// listTasks 列出所有任务
func (s *Server) listTasks(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := s.taskSvc.ListTasks(ctx, nil, 0)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list tasks")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

// createTask 创建新任务
func (s *Server) createTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		Prompt       string `json:"prompt" binding:"required"`
		WorkflowType string `json:"workflow_type,omitempty"`
		WorkflowID   string `json:"workflow_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用默认工作流类型
	workflowType := req.WorkflowType
	if workflowType == "" {
		workflowType = "feat"
	}

	task, err := s.taskSvc.CreateTask(ctx, req.Prompt, workflowType)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// 广播任务创建事件
	if s.wsHub != nil {
		s.wsHub.Broadcast(&websocket.Event{
			Type:   "task_created",
			TaskID: task.ID,
			Data:   task,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"task": task,
	})
}

// getTask 获取任务详情
func (s *Server) getTask(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	response, err := s.taskSvc.GetTaskStatus(ctx, taskID)
	if err != nil {
		logger.Error().Err(err).Str("task_id", taskID).Msg("Failed to get task")
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": response,
	})
}

// getTaskPhases 获取任务的所有阶段
func (s *Server) getTaskPhases(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	response, err := s.taskSvc.GetTaskStatus(ctx, taskID)
	if err != nil {
		logger.Error().Err(err).Str("task_id", taskID).Msg("Failed to get task phases")
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"phases": response.Phases,
	})
}

// getTaskEvents 获取任务的事件日志
// TODO: 需要在 TaskService 中添加 GetEvents 方法
func (s *Server) getTaskEvents(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// 当前 TaskService 没有提供获取事件列表的方法
	c.JSON(http.StatusOK, gin.H{
		"events":   []interface{}{},
		"message": "Event listing not fully implemented",
	})
}

// pauseTask 暂停任务
func (s *Server) pauseTask(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	if err := s.taskSvc.PauseTask(ctx, taskID); err != nil {
		logger.Error().Err(err).Str("task_id", taskID).Msg("Failed to pause task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pause task"})
		return
	}

	// 广播暂停事件
	if s.wsHub != nil {
		s.wsHub.Broadcast(&websocket.Event{
			Type:   "task_paused",
			TaskID: taskID,
			Data:   map[string]string{"task_id": taskID},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task paused successfully",
	})
}

// resumeTask 恢复任务
func (s *Server) resumeTask(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	if err := s.taskSvc.ResumeTask(ctx, taskID); err != nil {
		logger.Error().Err(err).Str("task_id", taskID).Msg("Failed to resume task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resume task"})
		return
	}

	// 广播恢复事件
	if s.wsHub != nil {
		s.wsHub.Broadcast(&websocket.Event{
			Type:   "task_resumed",
			TaskID: taskID,
			Data:   map[string]string{"task_id": taskID},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task resumed successfully",
	})
}

// cancelTask 取消任务
func (s *Server) cancelTask(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	if err := s.taskSvc.CancelTask(ctx, taskID); err != nil {
		logger.Error().Err(err).Str("task_id", taskID).Msg("Failed to cancel task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel task"})
		return
	}

	// 广播取消事件
	if s.wsHub != nil {
		s.wsHub.Broadcast(&websocket.Event{
			Type:   "task_cancelled",
			TaskID: taskID,
			Data:   map[string]string{"task_id": taskID},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task cancelled successfully",
	})
}

// ============================================================================
// Phase Handlers
// ============================================================================

// completePhase 完成阶段
// TODO: 需要实现通过 phaseID 获取 taskID 的逻辑
func (s *Server) completePhase(c *gin.Context) {
	phaseID := c.Param("id")

	if phaseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phase ID is required"})
		return
	}

	var req struct {
		Output string `json:"output"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 当前 ReportPhaseOutput 需要 taskID，需要先查询 phase 获取 taskID
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Phase completion not fully implemented",
	})
}

// failPhase 标记阶段失败
// TODO: 需要实现通过 phaseID 获取 taskID 的逻辑
func (s *Server) failPhase(c *gin.Context) {
	phaseID := c.Param("id")

	if phaseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phase ID is required"})
		return
	}

	var req struct {
		Error string `json:"error" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 当前 ReportPhaseError 需要 taskID，需要先查询 phase 获取 taskID
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Phase failure reporting not fully implemented",
	})
}

// ============================================================================
// Workflow Handlers
// ============================================================================

// listWorkflows 列出所有工作流
func (s *Server) listWorkflows(c *gin.Context) {
	ctx := c.Request.Context()

	workflows, err := s.workflowSvc.ListWorkflows(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list workflows")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list workflows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflows": workflows,
	})
}

// getWorkflow 获取工作流详情
func (s *Server) getWorkflow(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow name is required"})
		return
	}

	workflow, err := s.workflowSvc.GetWorkflow(ctx, name)
	if err != nil {
		logger.Error().Err(err).Str("workflow", name).Msg("Failed to get workflow")
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflow": workflow,
	})
}

// getWorkflowYAML 获取工作流 YAML 源文件
func (s *Server) getWorkflowYAML(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow name is required"})
		return
	}

	yamlData, err := s.workflowSvc.ExportWorkflow(ctx, name)
	if err != nil {
		logger.Error().Err(err).Str("workflow", name).Msg("Failed to get workflow YAML")
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow YAML not found"})
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Data(http.StatusOK, "text/plain; charset=utf-8", yamlData)
}

// getWorkflowByType 按类型获取工作流详情
func (s *Server) getWorkflowByType(c *gin.Context) {
	ctx := c.Request.Context()
	workflowType := c.Param("type")

	if workflowType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow type is required"})
		return
	}

	workflow, err := s.workflowSvc.GetWorkflowByType(ctx, workflowType)
	if err != nil {
		logger.Error().Err(err).Str("workflowType", workflowType).Msg("Failed to get workflow by type")
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflow": workflow,
	})
}
