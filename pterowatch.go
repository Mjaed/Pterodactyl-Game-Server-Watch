package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Mjaed/Pterodactyl-Game-Server-Watch/config"
	"github.com/Mjaed/Pterodactyl-Game-Server-Watch/pterodactyl"
	"github.com/Mjaed/Pterodactyl-Game-Server-Watch/servers"
	"github.com/Mjaed/Pterodactyl-Game-Server-Watch/update"
)

func main() {
	// Look for 'cfg' flag in command line arguments (default path: /etc/pterowatch/pterowatch.conf).
	configFile := flag.String("cfg", "/etc/pterowatch/pterowatch.conf", "The path to the Pterowatch config file.")
	flag.Parse()

	// Create config struct.
	cfg := config.Config{}

	// Set config defaults.
	cfg.SetDefaults()

	// Attempt to read config.
	cfg.ReadConfig(*configFile)

	// Level 1 debug.
	if cfg.DebugLevel > 0 {
		fmt.Println("[D1] Found config with API URL => " + cfg.APIURL + ". Token => " + cfg.Token + ". Auto Add Servers => " + strconv.FormatBool(cfg.AddServers) + ". Debug level => " + strconv.Itoa(cfg.DebugLevel) + ". Reload time => " + strconv.Itoa(cfg.ReloadTime))
	}

	// Level 2 debug.
	if cfg.DebugLevel > 1 {
		fmt.Println("[D2] Config default server values. Enable => " + strconv.FormatBool(cfg.DefEnable) + ". Scan time => " + strconv.Itoa(cfg.DefScanTime) + ". Max Fails => " + strconv.Itoa(cfg.DefMaxFails) + ". Max Restarts => " + strconv.Itoa(cfg.DefMaxRestarts) + ". Restart Interval => " + strconv.Itoa(cfg.DefRestartInt) + ". Report Only => " + strconv.FormatBool(cfg.DefReportOnly) + ". A2S Timeout => " + strconv.Itoa(cfg.DefA2STimeout) + ". Mentions => " + cfg.DefMentions + ".")
	}

	// Check if we want to automatically add servers.
	if cfg.AddServers {
		pterodactyl.AddServers(&cfg)
	}

	// Handle all servers (create timers, etc.).
	servers.HandleServers(&cfg, false)

	// Set config file for use later (e.g. updating/reloading).
	cfg.ConfLoc = *configFile

	// Initialize updater/reloader.
	update.Init(&cfg)

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
}
