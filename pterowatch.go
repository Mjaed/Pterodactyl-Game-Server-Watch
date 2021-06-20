package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gamemann/Pterodactyl-Game-Server-Watch/config"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/pterodactyl"
	"github.com/gamemann/Pterodactyl-Game-Server-Watch/query"
)

// Timer function.
func ServerWatch(server config.Server, timer *time.Ticker, fails *int, restarts *int, nextscan *int64, conn *net.UDPConn, apiURL string, apiToken string) {
	destroy := make(chan struct{})

	for {
		select {
		case <-timer.C:

			// Check if container status is 'on'.
			if !pterodactyl.CheckStatus(apiURL, apiToken, server.UID) {

				continue
			}

			// Send A2S_INFO request.
			query.SendRequest(conn)

			//fmt.Println("[" + server.IP + ":" + strconv.Itoa(server.Port) + "] A2S_INFO sent.")

			// Check for response. If no response, increase fail count. Otherwise, reset fail count to 0.
			if !query.CheckResponse(conn) {
				// Increase fail count.
				*fails++

				//fmt.Println("[" + server.IP + ":" + strconv.Itoa(server.Port) + "] Fails => " + strconv.Itoa(*fails))

				// Check to see if we want to restart the server.
				if *fails >= server.MaxFails && *restarts < server.MaxRestarts && *nextscan < time.Now().Unix() {
					//fmt.Println("[" + server.IP + ":" + strconv.Itoa(server.Port) + "] Fails exceeded.")

					// Attempt to kill container.
					pterodactyl.KillServer(apiURL, apiToken, server.UID)

					// Now attempt to start it again.
					pterodactyl.StartServer(apiURL, apiToken, server.UID)

					// Increment restarts count.
					*restarts++

					// Set next scan time and ensure the restart interval is at least 1.
					restartint := server.RestartInt

					if restartint < 1 {
						restartint = 120
					}

					// Get new scan time.
					*nextscan = time.Now().Unix() + int64(restartint)

					// Debug.
					fmt.Println(server.IP + ":" + strconv.Itoa(server.Port) + " was found down. Attempting to restart. Fail Count => " + strconv.Itoa(*fails) + ". Restart Count => " + strconv.Itoa(*restarts) + ".")
				}
			} else {
				// Reset everything.
				*fails = 0
				*restarts = 0
				*nextscan = 0
				*restarts = 0
			}

		case <-destroy:
			conn.Close()
			timer.Stop()
			return
		}
	}
}

func main() {
	// Specify config file path.
	configFile := "/etc/pterowatch/pterowatch.conf"

	// Create config struct.
	cfg := config.Config{}

	// Attempt to read config.
	config.ReadConfig(&cfg, configFile)

	// Check if we want to automatically add servers.
	if cfg.AddServers {
		pterodactyl.AddServers(&cfg)
	}

	// Loop through each container from the config.
	for i := 0; i < len(cfg.Servers); i++ {
		// Check if server is enabled for scanning.
		if !cfg.Servers[i].Enable {
			continue
		}

		// Specify server-specific variable.s
		var fails int = 0
		var restarts int = 0
		var nextscan int64 = 0

		// Get scan time.
		stime := cfg.Servers[i].ScanTime

		if stime < 1 {
			stime = 5
		}

		// Let's create the connection now.
		conn, err := query.CreateConnection(cfg.Servers[i].IP, cfg.Servers[i].Port)

		if err != nil {
			fmt.Println("Error creating UDP connection for " + cfg.Servers[i].IP + ":" + strconv.Itoa(cfg.Servers[i].Port))
			fmt.Println(err)

			return
		}

		// Create repeating timer.
		ticker := time.NewTicker(time.Duration(stime) * time.Second)
		go ServerWatch(cfg.Servers[i], ticker, &fails, &restarts, &nextscan, conn, cfg.APIURL, cfg.Token)
	}

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT)

	x := 0

	// Create a loop so the program doesn't exit. Look for signals and if SIGINT, stop the program.
	for x < 1 {
		kill := false
		s := <-sigc

		switch s {
		case os.Interrupt:
			kill = true
		}

		if kill {
			break
		}

		// Sleep every second to avoid unnecessary CPU consumption.
		time.Sleep(time.Duration(1) * time.Second)
	}
}
