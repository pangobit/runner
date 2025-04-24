package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Actions ActionsConfig `yaml:"actions"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int `yaml:"port"`
}

// ActionsConfig holds action-related configuration
type ActionsConfig struct {
	ComposeFiles []string `yaml:"composeFiles"`
}

var config Config

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	configPath := flag.String("config", "config/config.yaml", "Path to configuration YAML file")
	flag.Parse()

	if err := loadConfig(*configPath); err != nil {
		slog.Error("Failed to load config", "error", err, "path", *configPath)
		os.Exit(1)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/health", healthHandler)
	handler.HandleFunc("/update", updateHandler)

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Server.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("Starting server", "port", config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}

func loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &config)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received update webhook")

	for _, filePath := range config.Actions.ComposeFiles {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			slog.Error("Compose file does not exist", "file", filePath)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Compose file does not exist:", filePath)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Update process started")

	go func() {
		for _, filePath := range config.Actions.ComposeFiles {
			slog.Info("Executing docker-compose command", "file", filePath)
			cmd := exec.Command("docker-compose", "-f", filePath, "up", "-d")
			output, err := cmd.CombinedOutput()
			if err != nil {
				slog.Error("Error executing command", "file", filePath, "error", err, "output", string(output))
				continue
			}
			slog.Info("Command executed successfully", "file", filePath, "output", string(output))
		}
	}()
}
