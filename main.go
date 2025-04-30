package main

import (
	"github.com/denstiny/golang-language-server/pkg/engine"
	"github.com/denstiny/golang-language-server/pkg/route"
)

func main() {
	r := route.NewRoute()
	client := engine.NewClient(r)
	client.Start()
}
