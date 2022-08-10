package config

import "os"

type Config struct {
	LruSize    int    `json:"lru_size"`
	AppBaseDir string `json:"app_base_dir"`
}

var config *Config

func init() {
	baseDir, _ := os.Getwd()
	config = &Config{
		LruSize:    50,
		AppBaseDir: baseDir,
	}
}
func GetConfig() *Config {
	return config
}
