package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-gateway/internal/client"

	taskpb "github.com/AllaxSydia/trenager/proto/task"
)

type TaskHandler struct {
	clients *client.GRPCClients
}

func NewTaskHandler(clients *client.GRPCClients) *TaskHandler {
	return &TaskHandler{clients: clients}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	resp, err := h.clients.Task.ListTasks(r.Context(), &taskpb.ListTasksRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req taskpb.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.clients.Task.CreateTask(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
