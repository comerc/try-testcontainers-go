package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithRedisByGenericContainer(t *testing.T) {
	cfg := testcontainers.ReadConfig()
	fmt.Printf("%#v\n", cfg) // go test -v  ./...
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start redis: %s", err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatalf("Could not stop redis: %s", err)
		}
	}()
}

func TestWithRedisByModule(t *testing.T) {
	ctx := context.Background()
	redisContainer, err := redis.Run(ctx,
		"docker.io/redis:7", // or "docker.io/redis/redis-stack:latest",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
		redis.WithConfigFile(filepath.Join("testdata", "redis7.conf")),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	// Clean up the container
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
}
