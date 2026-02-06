package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gclm/gclm-flow/gclm-engine/internal/api/websocket"
	"github.com/gclm/gclm-flow/gclm-engine/internal/logger"
)

// sendError sends a standardized error response
func sendError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// sendInternalError sends an internal server error response with logging
func sendInternalError(c *gin.Context, err error, message string) {
	logger.Error().Err(err).Str("context", message).Msg("Request failed")
	sendError(c, http.StatusInternalServerError, message)
}

// sendInternalErrorWithFields sends an internal server error response with logging and fields
func sendInternalErrorWithFields(c *gin.Context, err error, message string, fields map[string]interface{}) {
	event := logger.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(message)
	sendError(c, http.StatusInternalServerError, message)
}

// getTaskIDParam extracts and validates the task ID from URL params
func getTaskIDParam(c *gin.Context) (string, bool) {
	taskID := c.Param("id")
	if taskID == "" {
		sendError(c, http.StatusBadRequest, "Task ID is required")
		return "", false
	}
	return taskID, true
}

// getPhaseIDParam extracts and validates the phase ID from URL params
func getPhaseIDParam(c *gin.Context) (string, bool) {
	phaseID := c.Param("id")
	if phaseID == "" {
		sendError(c, http.StatusBadRequest, "Phase ID is required")
		return "", false
	}
	return phaseID, true
}

// getWorkflowNameParam extracts and validates the workflow name from URL params
func getWorkflowNameParam(c *gin.Context) (string, bool) {
	name := c.Param("name")
	if name == "" {
		sendError(c, http.StatusBadRequest, "Workflow name is required")
		return "", false
	}
	return name, true
}

// broadcastTaskEvent broadcasts a task-related event via WebSocket
func (s *Server) broadcastTaskEvent(eventType string, taskID string, data interface{}) {
	if s.wsHub == nil {
		return
	}
	s.wsHub.Broadcast(&websocket.Event{
		Type:   eventType,
		TaskID: taskID,
		Data:   data,
	})
}

// bindJSON binds JSON request body to a struct and returns error if failed
func bindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
