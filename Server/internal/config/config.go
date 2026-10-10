// Package config holds server-wide settings.
package config

import (
	"flag"
	"time"
)

type Config struct {
	Addr         string        // HTTP/WebSocket listen address
	ClientDir    string        // directory with the static HTML client
	DataDir      string        // directory for account/character files
	TickRate     time.Duration // simulation step
	SnapshotEach int           // send a world snapshot every N ticks
	AutosaveEach time.Duration
	Dev          bool // enables /tp, /level, /gold ... test commands
}

// Load parses command line flags.
func Load() Config {
	c := Config{
		TickRate:     50 * time.Millisecond,
		SnapshotEach: 2,
		AutosaveEach: 30 * time.Second,
	}
	flag.StringVar(&c.Addr, "addr", ":8080", "listen address")
	flag.StringVar(&c.ClientDir, "client", "../Client", "path to client directory (empty = WebSocket API only)")
	flag.BoolVar(&c.Dev, "dev", false, "enable developer chat commands (/tp x z, /level n, /gold n, /unlock)")
	flag.StringVar(&c.DataDir, "data", "data", "path to data directory")
	flag.Parse()
	return c
}
