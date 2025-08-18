package main

import (
	"voteflix/api/internal/app"
	"voteflix/api/internal/routes"
)

func main() {
	app.Init(false).Serve(routes.Router)
}
