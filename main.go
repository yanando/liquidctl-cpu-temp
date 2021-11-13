package main

import (
	"github.com/yanando/cpu-temp/config"
	"github.com/yanando/cpu-temp/monitor"
)

func main() {
	config := config.GetConfig()

	tempMonitor := monitor.Monitor{
		Config: config,
	}

	tempMonitor.Start()
}
