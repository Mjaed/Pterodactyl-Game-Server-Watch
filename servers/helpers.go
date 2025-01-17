package servers

import (
	"net"
	"time"

	"github.com/Mjaed/Pterodactyl-Game-Server-Watch/config"
)

type Tuple struct {
	IP   string
	Port int
	UID  string
}

type Stats struct {
	Fails    *int
	Restarts *int
	NextScan *int64
}

type TickerHolder struct {
	Info      Tuple
	Ticker    *time.Ticker
	Conn      *net.UDPConn
	ScanTime  int
	Destroyer *chan bool
	Idx       *int
	Stats     Stats
}

func RemoveTicker(t *[]TickerHolder, idx int) {
	copy((*t)[idx:], (*t)[idx+1:])
	*t = (*t)[:len(*t)-1]
}

func RemoveServer(cfg *config.Config, idx int) {
	copy(cfg.Servers[idx:], cfg.Servers[idx+1:])
	cfg.Servers = cfg.Servers[:len(cfg.Servers)-1]
}
