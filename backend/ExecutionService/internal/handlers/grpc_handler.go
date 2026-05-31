package handler

import (
	"context"
	"log"

	"execution-service/internal/models"
	"execution-service/internal/service"

	pb "github.com/AllaxSydia/trenager/proto/execution"
)

type ExecutionGrpcHandler struct {
	pb.UnimplementedExecutionServiceServer
	executionService *service.ExecutionService
}

func NewExecutionGrpcHandler(executionService *service.ExecutionService) *ExecutionGrpcHandler {
	return &ExecutionGrpcHandler{
		executionService: executionService,
	}
}

func (h *ExecutionGrpcHandler) ExecuteCode(ctx context.Context, req *pb.ExecuteCodeRequest) (*pb.ExecuteCodeResponse, error) {
	log.Printf("ExecuteCode: language=%s", req.Language)

	execReq := &models.ExecutionRequest{
		Code:          req.Code,
		Language:      req.Language,
		Input:         req.Input,
		TimeLimitMs:   int(req.TimeLimitMs),
		MemoryLimitMB: int(req.MemoryLimitMb),
	}

	result, err := h.executionService.ExecuteCode(ctx, execReq)
	if err != nil {
		return &pb.ExecuteCodeResponse{
			Success: false,
			Error:   err.Error(),
			Status:  "failed",
		}, nil
	}

	return &pb.ExecuteCodeResponse{
		Success:         result.Success,
		Output:          result.Output,
		Error:           result.Error,
		ExecutionTimeMs: result.ExecutionTimeMs,
		MemoryUsedBytes: result.MemoryUsedBytes,
		Status:          result.Status,
	}, nil
}

func (h *ExecutionGrpcHandler) ExecuteTest(ctx context.Context, req *pb.ExecuteTestRequest) (*pb.ExecuteTestResponse, error) {
	log.Printf("ExecuteTest: language=%s, tests_count=%d", req.Language, len(req.Tests))

	var testCases []models.TestCase
	for _, test := range req.Tests {
		testCases = append(testCases, models.TestCase{
			Input:          test.Input,
			ExpectedOutput: test.ExpectedOutput,
		})
	}

	results, err := h.executionService.ExecuteTests(ctx, req.Code, req.Language, testCases)
	if err != nil {
		return nil, err
	}

	var pbResults []*pb.TestResult
	passedCount := 0

	for _, res := range results {
		pbResults = append(pbResults, &pb.TestResult{
			TestId:          res.TestID,
			Passed:          res.Passed,
			ActualOutput:    res.ActualOutput,
			Error:           res.Error,
			ExecutionTimeMs: res.ExecutionTimeMs,
		})
		if res.Passed {
			passedCount++
		}
	}

	return &pb.ExecuteTestResponse{
		AllPassed:   passedCount == len(results),
		Results:     pbResults,
		PassedCount: int32(passedCount),
		TotalCount:  int32(len(results)),
	}, nil
}

func (h *ExecutionGrpcHandler) GetExecutionStatus(ctx context.Context, req *pb.GetExecutionStatusRequest) (*pb.ExecutionStatus, error) {
	log.Printf("GetExecutionStatus: id=%s", req.ExecutionId)

	status, err := h.executionService.GetStatus(ctx, req.ExecutionId)
	if err != nil {
		return nil, err
	}

	return &pb.ExecutionStatus{
		ExecutionId: status.ID,
		Status:      status.Status,
		Result:      status.Result.Output,
		Error:       status.Result.Error,
	}, nil
}
