package main

import (
	"automatic-fiesta-go/main.go/firebase"
	server "automatic-fiesta-go/main.go/server"
	"fmt"
)

func main() {
	// Code here
	firebase.Init()
	server.Http_Server()
	fmt.Println("Let's Go!")

}
