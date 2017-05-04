package main

import (
	"github.com/evanlinjin/bbs-server/cmd"
)

func main() {
	config := cmd.MakeConfig()
	cmd.Run(config)
}
