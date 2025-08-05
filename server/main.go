package main

import (
	"gollama/routes"
	"gollama/config"
)

func main() {
	router := routes.Master()
	router.Run(":"+config.ENV.Port)
}
