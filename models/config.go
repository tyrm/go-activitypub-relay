package models

type Config struct {
	ID    int
	Key   string
	Value string
}

func GetConfig(k string) *Config {
	return nil
}