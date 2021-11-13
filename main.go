package main

import (
	"flag"

	"github.com/yanando/cpu-temp/config"
	"github.com/yanando/cpu-temp/liquidctl"
	"github.com/yanando/cpu-temp/monitor"
)

func main() {
	initLiquidctl := flag.Bool("initialize-liquidctl", false, "Initializes all liquidctl devices on startup")
	flag.Parse()

	if *initLiquidctl {
		liquidctl.InitLiquidctl()
	}

	config := config.GetConfig()

	tempMonitor := monitor.Monitor{
		Config: config,
	}

	tempMonitor.Start()
}
