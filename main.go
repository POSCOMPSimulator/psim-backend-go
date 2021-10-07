package main

import (
	"os"

	"poscomp-simulator.com/backend/controllers"
)

func main() {
	a := controllers.App{}
	a.Initialize()
	a.Run(":" + os.Getenv("PORT"))
}
