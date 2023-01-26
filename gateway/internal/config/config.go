package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Database DBConfig
}

type DBConfig struct {
	Driver string
	DSN    string
}
