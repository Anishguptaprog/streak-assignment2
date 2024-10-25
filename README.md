# Streak-Assignment

This is a gRPC-based User Management Service written in Go, providing functionalities to create, login, and logout users.

## Features

- **User Creation**: Allows users to register with a unique username and a strong password.
- **User Login**: Users can log in with their credentials, provided they are valid.
- **User Logout**: Users can log out of their session.
- **Password Validation**: Ensures passwords meet specified security criteria, including:
  - Minimum length of 8 characters
  - Inclusion of uppercase and lowercase letters
  - At least one numeric digit
  - At least one special character

## Technologies Used

- Go (Golang)
- gRPC
- Protocol Buffers (protobuf)

## Getting Started
- Start the server by running the command "go run main.go". Make sure you are in the server directory.
- Start the client side by running the command "go run client.go". Make sure you are in the client directory.
- Now you can create user, login and logout, or create multiple users.
