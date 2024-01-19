package main

import "github.com/vidarandrebo/oncetree"

func main() {
	client := oncetree.NewClient("test-client", []string{":8080"})
	client.Run()
}
