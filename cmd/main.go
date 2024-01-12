package main

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/core"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

// @title fantasy api doc
// main godoc
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	fx.New(core.Core()).Run()

}
