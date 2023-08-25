package main

import (
	"bp-transaction-api/configs"
	"bp-transaction-api/server"
	"flag"
)

func main() {
	configPath := flag.String(
		"config",
		"configs",
		"set configs path, default as: 'configs'",
	)
	flag.Parse()

	env, err := configs.Read(*configPath)
	if err != nil {
		panic(err)
	}

	server, err := server.New(env.Config)
	if err != nil {
		panic(err)
	}

	server.Start()
}
