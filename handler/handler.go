package handler

import (
	"github.com/yookoala/gofast"
	"net/http/fcgi"
	"os"
	"time"
	"strings"
	"fmt"
)

func NewRouterHandler(routeFile string) gofast.Middleware {
	router, err := loadRouter(routeFile)
	if err != nil {
		panic(err)
	}

	go func() {
		var prevTime time.Time
		for {
			select {
			case <-time.After(1 * time.Second):
				fd, err := os.Open(routeFile)
				if err != nil {
					continue
				}
				info, err := fd.Stat()
				if err != nil {
					continue
				}
				modTime := info.ModTime()
				if !modTime.After(prevTime) {
					continue
				}
				prevTime = modTime
				router, _ = loadRouter(routeFile)
				fmt.Println("router reloaded.")
			}
		}
	}()

	return func(inner gofast.SessionHandler) gofast.SessionHandler {
		return func(client gofast.Client, req *gofast.Request) (resp *gofast.ResponsePipe, err error) {
			env := fcgi.ProcessEnv(req.Raw)

			for k, v := range env {
				req.Params[k] = v
			}

			documentUri := env["DOCUMENT_URI"]

			resolved, args := router.Resolve(req.Raw.Method, documentUri)
			if resolved == "" {
				resolved, args = router.Resolve("GET", "__404__")
			}

			for k, v := range args {
				req.Params["ROUTER_ARG_"+strings.ToUpper(k)] = v
			}

			req.Params["SCRIPT_FILENAME"] = req.Params["SCRIPT_DIRECTORY"] + resolved
			req.Params["PATH_INFO"] = documentUri


			return inner(client, req)
		}
	}
}
