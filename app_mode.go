package nuwa

import (
	"log"
	"os"
	"time"

	"github.com/wailovet/gofunc"
	"github.com/wailovet/lorca"
	"golang.org/x/exp/errors/fmt"
)

var appMode = appModeImp{}

func AppMode() *appModeImp {
	return &appMode
}

type appModeImp struct {
}

func (*appModeImp) Run(w, h int, hes ...*HttpEngine) {
	port := Helper().GetFreePort()
	gofunc.New(func() {
		ui, err := lorca.New(fmt.Sprint("http://127.0.0.1:", port), Helper().GetSelfFilePath()+"/.cache/", w, h)
		if err != nil {
			fmt.Println(err)
			return
		}

		ui.Bind("setItem", func(k, v string) {
			NutsDB().Bucket("cache").Set(k, v)
		})
		ui.Bind("getItem", func(k string) string {
			return NutsDB().Bucket("cache").Get(k).String()
		})
		gofunc.New(func() {
			<-ui.Done()
			ui.Close()
			os.Exit(0)
		})
		gofunc.New(func() {
			for {
				ui.Eval(`
				setTimeout(function () {
					var xmlhttp = new XMLHttpRequest();
					xmlhttp.open("GET", "/", true);
					console.log(xmlhttp.send());
					xmlhttp.onreadystatechange = function () { 
						if (!xmlhttp.status) {
							window.close()
						}
					}
				},2000)
				`).Err()
				time.Sleep(time.Second)
			}
		})
		gofunc.Pause()

	}).Catch(func(i interface{}) {
		log.Println("UI Catch:", i)
	})
	var he *HttpEngine
	if len(hes) > 0 {
		he = hes[0]
	} else {
		he = Http()
	}
	he.InstanceConfig.Port = fmt.Sprint(port)
	he.Run()
}
