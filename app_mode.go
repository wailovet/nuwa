package nuwa

import (
	"log"
	"os"

	"github.com/wailovet/gofunc"
	"github.com/wailovet/lorca"
	"golang.org/x/exp/errors/fmt"
)

var appMode = appModeImp{}

func AppMode() *appModeImp {
	return &appMode
}

type appModeImp struct {
	e    func(ui lorca.UI)
	port int
	url  string
	ui   lorca.UI
}

func (ami *appModeImp) Event(e func(ui lorca.UI)) *appModeImp {
	ami.e = e
	return ami
}

func (ami *appModeImp) UI() lorca.UI {
	return ami.ui
}
func (ami *appModeImp) Load(url string) *appModeImp {
	ami.url = url
	return ami
}

func (ami *appModeImp) Port(port int) *appModeImp {
	ami.port = port
	return ami
}

func (ami *appModeImp) Run(w, h int, hes ...*HttpEngine) {
	if ami.port == 0 {
		ami.port = Helper().GetFreePort()
	}
	gofunc.New(func() {
		dirData := helper.AbsPath("lorca-data")
		if !helper.PathExists(dirData) {
			dirData = ""
		}

		if ami.url == "" {
			ami.url = fmt.Sprint("http://127.0.0.1:", ami.port)
		}

		ui, err := lorca.New(ami.url, dirData, w, h)
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
		// gofunc.New(func() {
		// 	for {
		// 		ui.Eval(`
		// 		if (!window["__interval"]) {
		// 			window["__interval"] = setInterval(function () {
		// 				var xmlhttp = new XMLHttpRequest();
		// 				xmlhttp.open("GET", "/", true);
		// 				console.log(xmlhttp.send());
		// 				xmlhttp.onreadystatechange = function () {
		// 					if (!xmlhttp.status) {
		// 						window.close()
		// 					}
		// 				}
		// 			},1000)
		// 		}
		// 		`).Err()
		// 		time.Sleep(time.Second)
		// 	}
		// })
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
	he.InstanceConfig.Port = fmt.Sprint(ami.port)
	he.Run()
}
