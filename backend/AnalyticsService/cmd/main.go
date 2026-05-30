package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"analytics-service/pkg/database"

	pb "github.com/AllaxSydia/trenager/proto/analytics"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
	db *database.Database
}

func NewAnalyticsServer(db *database.Database) *AnalyticsServer {
	return &AnalyticsServer{db: db}
}

func (s *AnalyticsServer) TrackEvent(ctx context.Context, req *pb.TrackEventRequest) (*pb.TrackEventResponse, error) {
	log.Printf("📊 TrackEvent: user=%s, event=%s", req.UserId, req.EventType)
	eventID := uuid.New()
	return &pb.TrackEventResponse{Success: true, EventId: eventID.String()}, nil
}

func (s *AnalyticsServer) GetUserProgress(ctx context.Context, req *pb.GetUserProgressRequest) (*pb.UserProgress, error) {
	log.Printf("📈 GetUserProgress: user=%s", req.UserId)
	return &pb.UserProgress{
		UserId:              req.UserId,
		TotalTasksAttempted: 0,
		TotalTasksCompleted: 0,
		OverallScore:        0,
		TotalXp:             0,
		CurrentStreakDays:   0,
		TasksByDifficulty:   make(map[string]int32),
		SkillsProgress:      make(map[string]float32),
	}, nil
}

func (s *AnalyticsServer) GetLeaderboard(ctx context.Context, req *pb.GetLeaderboardRequest) (*pb.Leaderboard, error) {
	log.Printf("🏆 GetLeaderboard: timeframe=%s", req.Timeframe)
	return &pb.Leaderboard{Entries: []*pb.LeaderboardEntry{}, Timeframe: req.Timeframe}, nil
}

func (s *AnalyticsServer) GetUserActivity(ctx context.Context, req *pb.GetUserActivityRequest) (*pb.UserActivity, error) {
	log.Printf("📅 GetUserActivity: user=%s", req.UserId)
	return &pb.UserActivity{UserId: req.UserId, Activities: []*pb.Activity{}}, nil
}

func (s *AnalyticsServer) GetTaskStats(ctx context.Context, req *pb.GetTaskStatsRequest) (*pb.TaskStats, error) {
	log.Printf("📊 GetTaskStats: task=%s", req.TaskId)
	return &pb.TaskStats{
		TaskId:             req.TaskId,
		TotalAttempts:      0,
		SuccessfulAttempts: 0,
		SuccessRate:        0,
		AverageScore:       0,
		ScoreDistribution:  make(map[string]int32),
	}, nil
}

func main() {
	log.Println("📊 Starting AnalyticsService gRPC server...")

	cfg := database.LoadConfig()
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterAnalyticsServiceServer(grpcServer, NewAnalyticsServer(db))
	reflection.Register(grpcServer) // Добавляем reflection

	lis, err := net.Listen("tcp", ":50056")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("✅ AnalyticsService running on :50056 with PostgreSQL")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
