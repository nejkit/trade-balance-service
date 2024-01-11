package main

import (
	"context"
	"trade-balance-service/app"
)

func main() {
	app.StartProgram(context.Background(), "", "")
}
