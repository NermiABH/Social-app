package main

import (
	"Social-app/internal/apiserver"
	"log"
)

func main() {
	apiConfig := apiserver.NewConfig()
	if err := apiserver.Start(apiConfig); err != nil {
		log.Fatalln(err)
	}
}
