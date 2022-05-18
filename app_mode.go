package nuwa

import (
	"log"
	"os"
	"time"

	"fmt"

	"github.com/wailovet/gofunc"
	"github.com/wailovet/lorca"
)

var appMode = appModeImp{}

func AppMode() *appModeImp {
	return &appMode
}

type appModeImp struct {
}

func (*appModeImp) Run(he *HttpEngine, w, h int) {
	port := Helper().GetFreePort()
	gofunc.New(func() {
		ui, err := lorca.New(fmt.Sprint("http://127.0.0.1:", port), "", w, h)
		if err != nil {
			fmt.Println(err)
			return
		}

		gofunc.New(func() {
			<-ui.Done()
			ui.Close()
			os.Exit(0)
		})
		time.Sleep(time.Second)

		ui.Eval(`
		setInterval(function () {
			var xmlhttp = new XMLHttpRequest();
			xmlhttp.open("GET", "/", true);
			console.log(xmlhttp.send());
			xmlhttp.onreadystatechange = function () { 
				if (!xmlhttp.status) {
					window.close()
				}
			}
		}, 1000) 
		`).Err()
		gofunc.Pause()

	}).Catch(func(i interface{}) {
		log.Println("UI Catch:", i)
	})

	he.InstanceConfig.Port = fmt.Sprint(port)
	he.Run()
}
