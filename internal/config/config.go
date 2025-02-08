package config

import (
	"flag"
)

type Config struct {
	Env         string
	StoragePath string
	Address     string
}



func MustLoad() *Config {
	addr := flag.String("addr", ":8080", "USAGE: :PORT, EX: \":8080\"")
	env := flag.String("env", "dev", "USAGE: DEV, EX: DEV|STAGE|PROD")
	dsn := flag.String("dsn", "./data/storage.db", "USAGE: STORAGE PATH, EX: ./data/storage.db")

	flag.Parse()

	cfg := Config{
		Env:         *env,
		Address:     *addr,
		StoragePath: *dsn,
	}

	return &cfg
}
