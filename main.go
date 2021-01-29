package main

import (
	"os"

	"github.com/sethvargo/go-signalcontext"
)

func main() {
	ctx, cancel := signalcontext.OnInterrupt()
	defer cancel()

	cli, err := getCLI(os.Args)
	if err != nil {
		panic(err)
	}

	config, err := loadConfig(cli)
	if err != nil {
		panic(err)
	}

	start(ctx, config)
}
