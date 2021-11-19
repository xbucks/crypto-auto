package views

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"

	"github.com/hirokimoto/crypto-auto/services"
	"github.com/zserge/lorca"
)

func (v *Views) OpenStables() error {
	v.WaitGroup.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		args := []string{}
		if runtime.GOOS == "linux" {
			args = append(args, "--class=Lorca")
		}
		ui, err := lorca.New("", "", 480, 320, args...)
		if err != nil {
			log.Fatal(err)
		}
		defer ui.Close()

		// Create and bind Go object to the UI
		c := &services.Tokens{}
		ui.Bind("addPair", func(t string) {
			fmt.Println(t)
		})
		ui.Bind("getPairs", c.Get)
		ui.Bind("getLength", c.GetLength)
		ui.Bind("getItem", c.GetItem)

		// A simple way to know when UI is ready (uses body.onload event in JS)
		ui.Bind("start", func() {
			log.Println("UI is ready")
			services.TrackStable(c)
		})

		// Load HTML.
		// You may also use `data:text/html,<base64>` approach to load initial HTML,
		// e.g: ui.Load("data:text/html," + url.PathEscape(html))

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()
		go http.Serve(ln, http.FileServer(http.FS(fs)))
		ui.Load(fmt.Sprintf("http://%s/www/stables.html", ln.Addr()))

		// You may use console.log to debug your JS code, it will be printed via
		// log.Println(). Also exceptions are printed in a similar manner.
		ui.Eval(`
			console.log("Hello, world!");
			console.log('Multiple values:', [1, false, {"x":5}]);
		`)

		// Wait until the interrupt signal arrives or browser window is closed
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)
		select {
		case <-sigc:
		case <-ui.Done():
		}

		log.Println("exiting...")
	}(v.WaitGroup)

	return nil
}
