package main

import (
	"github.com/wailovet/nuwa"
)

func main() {
	nuwa.Config().Host = "0.0.0.0"
	nuwa.Http().HandleFunc("/hello", func(request nuwa.Request, response nuwa.Response) {
		if len(request.REQUEST) > 0 {
			response.DisplayByData(request.REQUEST)
		}
		response.DisplayByData("Hello world!")
	})

	_ = nuwa.Http().Run()
}
