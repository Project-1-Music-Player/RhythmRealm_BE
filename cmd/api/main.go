package main

import (
	"fmt"
	"rr-backend/internal/auth"
	"rr-backend/internal/server"
)

func main() {
	auth.NewAuth()
	server := server.NewServer()

	fmt.Println("Server is running on port 3000")
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
