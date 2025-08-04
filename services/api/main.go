package main

import (
	"fmt"
	"log"
	"net/http"
	"voteflix/api/src/routes"
	"voteflix/api/src/utils"
)

func main() {
	log.Println("Starting API")

	serverErr := http.ListenAndServe(fmt.Sprintf(":%d", utils.GetAppConfig().Port), routes.Router())

	if nil != serverErr {
		log.Fatal(serverErr)
	}
}
