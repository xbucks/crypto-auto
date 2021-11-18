package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/hirokimoto/crypto-auto/services"
	"github.com/hirokimoto/crypto-auto/utils"
	"github.com/zserge/lorca"
)

type stables struct {
	sync.Mutex
	data []string
}

func (c *stables) Add(e string) {
	c.Lock()
	defer c.Unlock()
	c.data = append(c.data, e)
}

func (c *stables) Value() []string {
	c.Lock()
	defer c.Unlock()
	return c.data
}

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
		c := &stables{}
		ui.Bind("addPair", c.Add)
		ui.Bind("getPairs", c.Value)

		// A simple way to know when UI is ready (uses body.onload event in JS)
		ui.Bind("start", func() {
			log.Println("UI is ready")
			token := make(chan services.Token, 1)
			go trackStable(token)
			msg := <-token
			c.Add(msg.Get())
			fmt.Println("New token!!!!!   ", msg.Get())
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

func trackStable(t chan services.Token) {
	pc := make(chan string, 1)
	for {
		go utils.Post(pc, "pairs", 1000, "")

		msg1 := <-pc
		var pairs utils.Pairs
		json.Unmarshal([]byte(msg1), &pairs)
		counts := len(pairs.Data.Pairs)
		fmt.Println("Counts of Pairs: ", counts)
		if counts > 0 {
			var wg sync.WaitGroup
			wg.Add(counts)
			services.StableTokens(&wg, pairs, t)
			wg.Wait()
		}

		time.Sleep(time.Minute * 10)
	}
}
