package main

import (
	"chat/internal"
	"flag"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

var (
	port    = flag.Int("port", 3333, "the port to listen ")
	timeout = flag.Int("timeout", 0, "timeout of milliseconds ")
	cpu     = flag.Int("cpu", 500, "cpu threshold ")
)

func main() {
	flag.Parse()

	logx.Disable()
	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Host:         "localhost",
		Port:         *port,
		Timeout:      int64(*timeout),
		CpuThreshold: int64(*cpu),
	})
	defer engine.Stop()

	hub := internal.NewHub()
	go hub.Run()

	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			if r.Method != "GET" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			http.ServeFile(w, r, "home.html")
		},
	})

	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/ws",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			internal.ServeWs(hub, w, r)
		},
	})

	engine.Start()
}
