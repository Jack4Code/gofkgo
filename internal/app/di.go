package di

import (
	"gofkgo/internal/ports"
	"gofkgo/internal/protocol"
	"gofkgo/internal/server"
)

type App struct {
	addr   string
	Server ports.Server
}

func NewApp(addr string) *App {
	return &App{addr: addr}
}

func (app *App) RegisterDependencies() {
	var handler ports.ProtocolPort = protocol.NewProtocolPort()
	app.Server = server.NewServer(app.addr, handler)
}

func (app *App) Run() {
	err := app.Server.Start()
	if err != nil {
		panic(err)
	}
}
