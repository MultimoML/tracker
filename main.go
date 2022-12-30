package main

import (
	"context"

	"github.com/multimoml/tracker/internal/server"
)

func main() {
	ctx := context.Background()
	server.Run(ctx)
}
