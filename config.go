package main

type Config struct {
	Server struct {
		APIHost   string `mapstructure:"api_host"`
		DebugHost string `mapstructure:"debug_host"`
		Version   string `mapstructure:"version"`
	}
}
