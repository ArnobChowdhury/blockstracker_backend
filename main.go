package main

import (
	"fmt"
	"net/http"

	"blockstracker_backend/routes"

	"github.com/gin-gonic/gin"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Air hot reload is working!")
}

func main() {
	r := gin.Default()

	routes.RegisterTaskRoutes(r)

	fmt.Println("Server started on :8080")
	r.Run(":8080")
}
