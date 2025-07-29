package main

import (
	"errors"
	di "gofkgo/internal/app"
)

func main() {
	app := di.NewApp(":9090")
	if app == nil {
		panic(errors.New("app is nil"))
	}
	app.RegisterDependencies()
	app.Run()
}
