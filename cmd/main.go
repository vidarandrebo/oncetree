package main

import (
	"github.com/vidarandrebo/once-tree/pkg/server"
)

func main() {
	srv := server.NewServer(":8080")
	srv.Run()
}
