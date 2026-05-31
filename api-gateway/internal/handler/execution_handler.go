package handler

import (
	"encoding/json"
	"net/http"

	"api-gateway/internal/client"

	executionpb "github.com/AllaxSydia/trenager/proto/execution"
)

type ExecutionHandler struct {
	clients *client.GRPCClients
}

func NewExecutionHandler(clients *client.GRPCClients) *ExecutionHandler {
	return &ExecutionHandler{clients: clients}
}

func (h *ExecutionHandler) ExecuteCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code        string `json:"code"`
		Language    string `json:"language"`
		Input       string `json:"input,omitempty"`
		TimeLimitMs int32  `json:"time_limit_ms,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Code == "" {
		writeError(w, http.StatusBadRequest, "code is required")
		return
	}

	if req.Language == "" {
		writeError(w, http.StatusBadRequest, "language is required")
		return
	}

	resp, err := h.clients.Execution.ExecuteCode(r.Context(), &executionpb.ExecuteCodeRequest{
		Code:        req.Code,
		Language:    req.Language,
		Input:       req.Input,
		TimeLimitMs: req.TimeLimitMs,
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ExecutionHandler) ExecuteTest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code     string `json:"code"`
		Language string `json:"language"`
		Tests    []struct {
			Input          string `json:"input"`
			ExpectedOutput string `json:"expected_output"`
		} `json:"tests"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var tests []*executionpb.TestCase
	for _, test := range req.Tests {
		tests = append(tests, &executionpb.TestCase{
			Input:          test.Input,
			ExpectedOutput: test.ExpectedOutput,
		})
	}

	resp, err := h.clients.Execution.ExecuteTest(r.Context(), &executionpb.ExecuteTestRequest{
		Code:     req.Code,
		Language: req.Language,
		Tests:    tests,
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
