package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/AllaxSydia/trenager/proto/analytics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50056",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewAnalyticsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("=== Testing AnalyticsService ===\n")

	// 1. Track Event
	fmt.Println("1. Tracking event...")
	trackResp, err := client.TrackEvent(ctx, &pb.TrackEventRequest{
		UserId:    "user-123",
		EventType: "task_completed",
		Metadata:  map[string]string{"task_id": "task-456"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Event tracked: %s\n", trackResp.EventId)

	// 2. Get User Progress
	fmt.Println("\n2. Getting user progress...")
	progress, err := client.GetUserProgress(ctx, &pb.GetUserProgressRequest{
		UserId: "user-123",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Completed: %d/%d tasks\n", progress.TotalTasksCompleted, progress.TotalTasksAttempted)
	fmt.Printf("   📊 Score: %.1f%%\n", progress.OverallScore)
	fmt.Printf("   ⭐ XP: %d\n", progress.TotalXp)
	fmt.Printf("   🔥 Streak: %d days\n", progress.CurrentStreakDays)

	// 3. Get Leaderboard
	fmt.Println("\n3. Getting leaderboard...")
	leaderboard, err := client.GetLeaderboard(ctx, &pb.GetLeaderboardRequest{
		Timeframe: "weekly",
		Limit:     5,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Found %d leaders\n", len(leaderboard.Entries))
	for _, entry := range leaderboard.Entries {
		fmt.Printf("      %d. %s - %d points (%d tasks)\n",
			entry.Rank, entry.Username, entry.Score, entry.TasksCompleted)
	}

	// 4. Get User Activity
	fmt.Println("\n4. Getting user activity...")
	activity, err := client.GetUserActivity(ctx, &pb.GetUserActivityRequest{
		UserId: "user-123",
		Days:   7,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Last %d days: %d activities\n", len(activity.Activities), len(activity.Activities))

	// 5. Get Task Stats
	fmt.Println("\n5. Getting task statistics...")
	taskStats, err := client.GetTaskStats(ctx, &pb.GetTaskStatsRequest{
		TaskId: "task-456",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Attempts: %d, Success rate: %.1f%%\n",
		taskStats.TotalAttempts, taskStats.SuccessRate)
	fmt.Printf("   📊 Average score: %.1f\n", taskStats.AverageScore)

	fmt.Println("\n🎉 All tests passed!")
}
