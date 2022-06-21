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
	e func(ui lorca.UI)
}

func (ami *appModeImp) Event(e func(ui lorca.UI)) *appModeImp {
	ami.e = e
	return ami
}
func (ami *appModeImp) Run(w, h int, hes ...*HttpEngine) {
	port := Helper().GetFreePort()
	gofunc.New(func() {
		dirData := helper.AbsPath("lorca-data")
		if !helper.PathExists(dirData) {
			dirData = ""
		}
		ui, err := lorca.New(fmt.Sprint("http://127.0.0.1:", port), dirData, w, h)
		if err != nil {
			fmt.Println(err)
			return
		}
		if ami.e != nil {
			ami.e(ui)
		}
		gofunc.New(func() {
			<-ui.Done()
			ui.Close()
			os.Exit(0)
		})
		gofunc.New(func() {
			for {
				ui.Eval(`
				if (!window["__interval"]) {
					window["__interval"] = setInterval(function () {
						var xmlhttp = new XMLHttpRequest();
						xmlhttp.open("GET", "/", true);
						console.log(xmlhttp.send());
						xmlhttp.onreadystatechange = function () { 
							if (!xmlhttp.status) {
								window.close()
							}
						}
					},1000)
				}
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
