package shard

import (
	"bytes"
	"testing"
)

func TestShardManagerSplitData(t *testing.T) {
	sm := NewShardManager(3, "test-node")
	testData := []byte("Hello, World!")

	t.Run("Split and Reconstruct", func(t *testing.T) {
		shards := sm.SplitData(testData)
		if len(shards) != 3 {
			t.Fatalf("Expected 3 shards, got %d", len(shards))
		}
		var reconstructed []byte
		for i := 0; i < 3; i++ {
			reconstructed = append(reconstructed, shards[i]...)
		}
		if !bytes.Equal(testData, reconstructed) {
			t.Error("Data reconstruction failed")
		}
	})
}

func TestShardManagerEncryption(t *testing.T) {
	sm := NewShardManager(3, "test-node")
	testData := []byte("Secret message")
	key := "12345678901234567890123456789012" // 32-byte key

	t.Run("Encryption/Decryption Success", func(t *testing.T) {
		encrypted, err := sm.EncryptShard(testData, key)
		if err != nil {
			t.Fatalf("Encryption failed: %v", err)
		}
		decrypted, err := sm.DecryptShard(encrypted, key)
		if err != nil {
			t.Fatalf("Decryption failed: %v", err)
		}
		if !bytes.Equal(testData, decrypted) {
			t.Error("Encryption/decryption cycle failed")
		}
	})

	t.Run("Decryption with Wrong Key", func(t *testing.T) {
		encrypted, err := sm.EncryptShard(testData, key)
		if err != nil {
			t.Fatalf("Encryption failed: %v", err)
		}
		_, err = sm.DecryptShard(encrypted, "invalid-key-incorrect-length")
		if err == nil {
			t.Error("Expected error when decrypting with wrong key")
		}
	})
}
