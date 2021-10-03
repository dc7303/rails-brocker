package main

import (
	"github.com/dc7303/rails-brocker/api"
)

func main() {
	server := api.New()
	server.HandleRequests()
}
