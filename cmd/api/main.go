package main

import (
	"fmt"
	"rr-backend/internal/auth"
	"rr-backend/internal/server"
)

func main() {
	auth.NewAuth()
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
