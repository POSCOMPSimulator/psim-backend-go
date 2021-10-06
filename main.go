package main

import (
	"poscomp-simulator.com/backend/controllers"
)

func main() {
	a := controllers.App{}
	a.Initialize()
	a.Run(":8060")
}
