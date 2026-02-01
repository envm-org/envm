package config

type Config struct {
	DatabaseURI string
	Addr string
	db dbConfig
}


type dbConfig struct {
	DSN string
}

