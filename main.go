package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     int      `yaml:"port"`
	Commands []string `yaml:"commands"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	return &config, err
}

func getAuthToken() string {
	token := os.Getenv("RUNNER_AUTH_TOKEN")
	if token == "" {
		slog.Error("RUNNER_AUTH_TOKEN environment variable not set")
		os.Exit(1)
	}
	return token
}

func execCommand(command string) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Command execution failed", "command", command, "error", err)
	}
	slog.Info("Command output", "command", command, "output", string(output))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func updateHandler(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		expectedToken := "Bearer " + getAuthToken()

		if token != expectedToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		for _, cmd := range config.Commands {
			execCommand(cmd)
		}

		fmt.Fprintln(w, "Update triggered successfully")
	}
}

func main() {
	// Setup structured JSON logger for OTEL and log scraper compatibility
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	}
	logHandler := slog.NewJSONHandler(os.Stdout, logOpts)
	slog.SetDefault(slog.New(logHandler))

	config, err := loadConfig("config.yaml")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/update", updateHandler(config))

	slog.Info("Server started", "port", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}


