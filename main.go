package main

import (
	"github.com/evanlinjin/bbs/cmd"
)

func main() {
	config := cmd.NewConfig()
	config.Parse()
	cmd.Run(config)
}
