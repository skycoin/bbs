package main

import (
	"github.com/evanlinjin/bbs/cmd"
)

func main() {
	config := cmd.MakeConfig()
	config.Parse()
	cmd.Run(config)
}
