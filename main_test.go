package main

import (
	"os"
	"testing"

	"github.com/posthog/posthog-go"
)

func BenchmarkEventsPipeline(b *testing.B) {
	posthogKey := os.Getenv("POSTHOG_KEY")
	testEndpoint := os.Getenv("POSTHOG_ENDPOINT")
	config := posthog.Config{
		Endpoint: testEndpoint,
	}
	client, err := posthog.NewWithConfig(posthogKey, config)
	if err != nil {
		panic("oh no")
	}
	defer client.Close()

	for i := 0; i < b.N; i++ {
		runOneEventPipeline(i, config, client)
	}
}
