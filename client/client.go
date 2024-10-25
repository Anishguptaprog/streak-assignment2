package main

import (
	"context"
	"fmt"
	"log"
	"streak/user"

	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)

	var username, password string
	var action string
	var loggedIn bool

	for {
		// Ask user to create or login
		fmt.Print("Would you like to (c)reate a user or (l)ogin? ")
		fmt.Scan(&action)

		switch action {
		case "c":
			// User creation
			fmt.Print("Enter username: ")
			fmt.Scan(&username)
			fmt.Print("Enter password: ")
			fmt.Scan(&password)

			response, err := client.CreateUser(context.Background(), &user.CreateUserRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				log.Fatalf("Error calling CreateUser: %v", err)
			}
			log.Printf("CreateUser Response: %v", response)

		case "l":
			// User login
			if loggedIn {
				log.Println("You are already logged in. Please logout first.")
				continue
			}

			fmt.Print("Enter username: ")
			fmt.Scan(&username)
			fmt.Print("Enter password: ")
			fmt.Scan(&password)

			loginResponse, err := client.LoginUser(context.Background(), &user.LoginUserRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				log.Fatalf("Error calling LoginUser: %v", err)
			}

			if !loginResponse.Success {
				log.Println("Login failed:", loginResponse.Message)
				continue
			}

			loggedIn = true
			log.Println("Login successful!")

			// Ask for logout
			fmt.Print("Do you want to logout? (y/n): ")
			var logoutChoice string
			fmt.Scan(&logoutChoice)

			if logoutChoice == "y" {
				logoutResponse, err := client.LogoutUser(context.Background(), &user.LogoutUserRequest{
					Username: username,
				})
				if err != nil {
					log.Fatalf("Error calling LogoutUser: %v", err)
				}
				log.Printf("LogoutUser Response: %v", logoutResponse)
				loggedIn = false
			}

		default:
			log.Println("Invalid option. Please choose 'c' or 'l'.")
		}
	}
}
