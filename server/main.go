package main

import (
	"context"
	"log"
	"net"
	"regexp"
	"sync"

	pb "streak/user"

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
	users map[string]*User
}

var userStore = &UserStore{
	users: make(map[string]*User),
}

type server struct {
	pb.UnimplementedUserServiceServer
}

func isValidPassword(password string) bool {

	if len(password) < 8 {
		return false
	}

	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>]`, password); !matched {
		return false
	}

	return true
}

// CreateUser Method
func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userStore.Lock()
	defer userStore.Unlock()

	if _, exists := userStore.users[req.Username]; exists {
		return &pb.CreateUserResponse{Success: false, Message: "User already exists"}, nil
	}
	if !isValidPassword(req.Password) {
		return &pb.CreateUserResponse{Success: false, Message: "Password must be at least 8 characters long and include upper case letters, lower case letters, digits, and special characters"}, nil
	}

	userStore.users[req.Username] = &User{
		Username: req.Username,
		Password: req.Password,
		LoggedIn: false,
	}
	return &pb.CreateUserResponse{Success: true, Message: "User created successfully"}, nil
}

// LoginUser method
func (s *server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	userStore.Lock()
	defer userStore.Unlock()

	user, exists := userStore.users[req.Username]
	if !exists || user.Password != req.Password {
		return &pb.LoginUserResponse{Success: false, Message: "Invalid credentials"}, nil
	}

	if user.LoggedIn {
		return &pb.LoginUserResponse{Success: false, Message: "User already logged in"}, nil
	}

	// Mark user as logged in
	user.LoggedIn = true
	return &pb.LoginUserResponse{Success: true, Message: "User logged in successfully"}, nil
}

// LogoutUser method
func (s *server) LogoutUser(ctx context.Context, req *pb.LogoutUserRequest) (*pb.LogoutUserResponse, error) {
	userStore.Lock()
	defer userStore.Unlock()

	user, exists := userStore.users[req.Username]
	if !exists || !user.LoggedIn {
		return &pb.LogoutUserResponse{Success: false, Message: "User is not logged in"}, nil
	}

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
