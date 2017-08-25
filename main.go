package main

import (
	"log"
	"os"
)

func main() {
	conf := os.Getenv("CONFIG_PATH")
	port := os.Getenv("SERVER_PORT")
	if conf == "" {
		conf = "config.toml"
	}
	if port == "" {
		port = "3000"
	}
	config := Config{}
	config.Read(conf)

	conn := ""
	if len(os.Getenv("MONGODB_USERNAME")) > 0 {
		conn += os.Getenv("MONGODB_USERNAME")

		if len(os.Getenv("MONGODB_PASSWORD")) > 0 {
			conn += ":" + os.Getenv("MONGODB_PASSWORD")
		}

		conn += "@"
	}
	if len(os.Getenv("MONGODB_PORT_27017_TCP_ADDR")) > 0 {
		conn += os.Getenv("MONGODB_PORT_27017_TCP_ADDR")
	} else {
		conn += config.Server
	}
	config.Server = conn

	a := App{}
	log.Println("Connecting to database...")
	a.Connect(config)
	log.Printf("Starting server...")
	a.Run(port)
}
