package main

import (
	"github.com/wailovet/nuwa"
)

func main() {
	nuwa.Config().Host = "0.0.0.0"
	nuwa.Http().HandleFunc("/hello", func(ctx nuwa.HttpContext) {
		if len(ctx.REQUEST) > 0 {
			ctx.DisplayByData(ctx.REQUEST)
		}
		ctx.DisplayByData("Hello world!")
	})

	_ = nuwa.Http().Run()
}
