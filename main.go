package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Air hot reload is working!")
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Wow! This works :)")
	})

	fmt.Println("Server started on :8080")
	r.Run(":8080")
}
