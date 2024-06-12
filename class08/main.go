package main

import "class08/routes"

func main() {
	server := routes.NewServer()

	server.Run(":8000")
}
