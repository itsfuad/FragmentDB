package shard

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

type ShardManager struct {
	TotalShards int
	NodeID      string
}

func NewShardManager(shards int, nodeID string) *ShardManager {
	return &ShardManager{
		TotalShards: shards,
		NodeID:      nodeID,
	}
}

func (sm *ShardManager) GetShard(key string) int {
	hash := sha256.Sum256([]byte(key))
	hexHash := hex.EncodeToString(hash[:])
	return int(hexHash[0]) % sm.TotalShards
}

func (sm *ShardManager) SplitData(data []byte) map[int][]byte {
	shards := make(map[int][]byte)
	chunkSize := len(data) / sm.TotalShards

	for i := 0; i < sm.TotalShards; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == sm.TotalShards-1 {
			end = len(data)
		}
		shards[i] = data[start:end]
	}
	return shards
}

func (sm *ShardManager) EncryptShard(data []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (sm *ShardManager) DecryptShard(data []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
