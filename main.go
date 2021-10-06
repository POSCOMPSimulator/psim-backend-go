package main

import (
	"poscomp-simulator.com/backend/app"
)

func main() {
	a := app.App{}
	a.Initialize()
	a.Run(":8010")
}
