package main

import (
	"gitlab.intelligrape.net/tothenew/tcp-server-integration/util"
	"gitlab.intelligrape.net/tothenew/tcp-server-integration/tcp"
)

func main() {
	tcp.RegisterTCPListener(util.Args.Port)
}
