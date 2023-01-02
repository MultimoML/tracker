package main

import (
	"context"

	"github.com/multimoml/tracker/internal/server"
)

// @title Tracker API
// @version 1.0.0
// @host localhost:6003
// @BasePath /tracker
func main() {
	ctx := context.Background()
	server.Run(ctx)
}
