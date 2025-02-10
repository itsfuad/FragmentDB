package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServerConfig struct {
	NodeID     string   `json:"node_id"`
	Port       int      `json:"port"`
	PeerNodes  []string `json:"peer_nodes"`
	DataPath   string   `json:"data_path"`
	SecretKey  string   `json:"secret_key"`
	ShardCount int      `json:"shard_count"`
}

func LoadConfig(path string) (*ServerConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if config.ShardCount < 3 {
		config.ShardCount = 3 // Minimum shard count
	}

	if config.SecretKey == "" {
		return nil, fmt.Errorf("secret_key cannot be empty")
	}

	return &config, nil
}
