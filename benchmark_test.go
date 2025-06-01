package firefly

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// BenchmarkClientCreation benchmarks client creation performance
func BenchmarkClientCreation(b *testing.B) {
	baseURL := "https://demo.firefly-iii.org"
	token := "test-token"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := NewFireflyClient(baseURL, token)
		if err != nil {
			b.Fatal(err)
		}
		_ = client
	}
}

// BenchmarkConvenienceClientCreation benchmarks convenience function performance
func BenchmarkConvenienceClientCreation(b *testing.B) {
	baseURL := "https://demo.firefly-iii.org"
	token := "test-token"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewFirefly(baseURL, token)
		_ = client
	}
}

// BenchmarkTransactionModelCreation benchmarks transaction model creation
func BenchmarkTransactionModelCreation(b *testing.B) {
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transaction := TransactionModel{
			ID:          "benchmark-test",
			Currency:    "USD",
			Amount:      100.50,
			TransType:   "deposit",
			Description: "Benchmark transaction",
			Date:        now,
			Category:    "Testing",
		}
		_ = transaction
	}
}

// BenchmarkContextCreation benchmarks context creation for API calls
func BenchmarkContextCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		_ = ctx
		cancel()
	}
}

// TestBenchmarkValidation tests that benchmark functions work correctly
func TestBenchmarkValidation(t *testing.T) {
	// Test that benchmark functions don't panic

	t.Run("ClientCreation", func(t *testing.T) {
		assert.NotPanics(t, func() {
			result := testing.Benchmark(BenchmarkClientCreation)
			assert.Greater(t, result.N, 0)
		})
	})

	t.Run("ConvenienceClientCreation", func(t *testing.T) {
		assert.NotPanics(t, func() {
			result := testing.Benchmark(BenchmarkConvenienceClientCreation)
			assert.Greater(t, result.N, 0)
		})
	})

	t.Run("TransactionModelCreation", func(t *testing.T) {
		assert.NotPanics(t, func() {
			result := testing.Benchmark(BenchmarkTransactionModelCreation)
			assert.Greater(t, result.N, 0)
		})
	})

	t.Run("ContextCreation", func(t *testing.T) {
		assert.NotPanics(t, func() {
			result := testing.Benchmark(BenchmarkContextCreation)
			assert.Greater(t, result.N, 0)
		})
	})
}
