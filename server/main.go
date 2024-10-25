package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "streak/user" // Adjust to your actual module path

	"google.golang.org/grpc"
)

// User structure to hold user details
type User struct {
	Username string
	Password string
	LoggedIn bool
}

// UserStore to hold registered users
type UserStore struct {
	sync.RWMutex
	users map[string]*User // map to store users by username
}

var userStore = &UserStore{
	users: make(map[string]*User),
}

// Define the server struct
type server struct {
	pb.UnimplementedUserServiceServer
}

// Implement CreateUser method
func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userStore.Lock()
	defer userStore.Unlock()

	// Check if user already exists
	if _, exists := userStore.users[req.Username]; exists {
		return &pb.CreateUserResponse{Success: false, Message: "User already exists"}, nil
	}

	// Add new user
	userStore.users[req.Username] = &User{
		Username: req.Username,
		Password: req.Password,
		LoggedIn: false,
	}
	return &pb.CreateUserResponse{Success: true, Message: "User created successfully"}, nil
}

// Implement LoginUser method
func (s *server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	userStore.Lock()
	defer userStore.Unlock()

	user, exists := userStore.users[req.Username]
	if !exists || user.Password != req.Password {
		return &pb.LoginUserResponse{Success: false, Message: "Invalid credentials"}, nil
	}

	// Check if user is already logged in
	if user.LoggedIn {
		return &pb.LoginUserResponse{Success: false, Message: "User already logged in"}, nil
	}

	// Mark user as logged in
	user.LoggedIn = true
	return &pb.LoginUserResponse{Success: true, Message: "User logged in successfully"}, nil
}

// Implement LogoutUser method
func (s *server) LogoutUser(ctx context.Context, req *pb.LogoutUserRequest) (*pb.LogoutUserResponse, error) {
	userStore.Lock()
	defer userStore.Unlock()

	user, exists := userStore.users[req.Username]
	if !exists || !user.LoggedIn {
		return &pb.LogoutUserResponse{Success: false, Message: "User is not logged in"}, nil
	}

	// Mark user as logged out
	user.LoggedIn = false
	return &pb.LogoutUserResponse{Success: true, Message: "User logged out successfully"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})

	log.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
