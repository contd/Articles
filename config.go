package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Represents database server and credentials
type Config struct {
	Server     string
	Database   string
	Collection string
}

// Read and parse the configuration file
func (c *Config) Read(conf string) {
	if _, err := toml.DecodeFile(conf, &c); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Using database: %s on %s", c.Database, c.Server)
	}
}
