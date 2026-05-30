package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/AllaxSydia/trenager/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("=== Testing AuthService ===\n")

	// 1. Register
	fmt.Println("1. Registering user...")
	registerResp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Success: %v\n", registerResp.Success)
	fmt.Printf("   📝 Message: %s\n", registerResp.Message)

	// 2. Login
	fmt.Println("\n2. Logging in...")
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Success: %v\n", loginResp.Success)
	fmt.Printf("   👤 User: %s\n", loginResp.User.Username)

	// 3. Validate Token
	fmt.Println("\n3. Validating token...")
	validateResp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		AccessToken: loginResp.AccessToken,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Valid: %v\n", validateResp.Valid)

	// 4. Get User
	fmt.Println("\n4. Getting user info...")
	userResp, err := client.GetUser(ctx, &pb.GetUserRequest{
		UserId: "user-123",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ User: %s (%s)\n", userResp.User.Username, userResp.User.Email)

	// 5. Health Check
	fmt.Println("\n5. Health check...")
	healthResp, err := client.HealthCheck(ctx, &pb.HealthRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✅ Status: %s\n", healthResp.Status)

	fmt.Println("\n🎉 All tests passed!")
}
