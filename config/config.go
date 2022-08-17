package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LruSize       int    `json:"lru_size"`
	AppBaseDir    string `json:"app_base_dir"`
	IsDistributed bool   `json:"is_distributed"`
	Addr          string `json:"addr"`
	Port          int    `json:"port"`
	PlatformAddr  string `json:"platform_addr"`
	PlatformPort  int    `json:"platform_port"`
}

var config *Config

func init() {
	// baseDir, _ := os.Getwd()
	baseDir := "C:/work/chen/code/distribute_cache"
	ymlFilePath := fmt.Sprintf("%s/%s", baseDir, "config/cache.yml")
	rb, _ := os.ReadFile(ymlFilePath)
	config = new(Config)
	yaml.Unmarshal(rb, config)
	config.AppBaseDir = baseDir
}
func GetConfig() *Config {
	return config
}
