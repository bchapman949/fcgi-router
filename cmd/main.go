package main

import (
	"github.com/yookoala/gofast"
	"net"
	"net/http/fcgi"
	"github.com/kr9ly/fcgirouter/handler"
	"os"
	"strconv"
)

func main() {
	listenAddr := os.Getenv("ROUTER_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":9000"
	}
	backendAddr := os.Getenv("ROUTER_BACKEND_ADDR")
	if backendAddr == "" {
		backendAddr = ":9001"
	}
	clientLimit := os.Getenv("ROUTER_CLIENT_LIMIT")
	limitInt, err := strconv.Atoi(clientLimit)
	if err != nil {
		limitInt = 100
	}
	configPath := os.Getenv("ROUTER_CONFIG_PATH")
	if configPath == "" {
		configPath = "routes.yml"
	}

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	connFactory := gofast.SimpleConnFactory("tcp", backendAddr)
	fcgiClient := gofast.SimpleClientFactory(connFactory, uint32(limitInt))

	fcgi.Serve(l, gofast.NewHandler(
		gofast.Chain(
			gofast.BasicParamsMap,
			gofast.MapHeader,
			handler.NewRouterHandler(configPath),
		)(gofast.BasicSession),
		fcgiClient,
	))
}
