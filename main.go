package main

import (
	"fmt"
	"log"
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
		log.Fatal("RUNNER_AUTH_TOKEN environment variable not set")
	}
	return token
}

func execCommand(command string) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command '%s' failed: %s\n", command, err)
	}
	log.Printf("Output for '%s': %s\n", command, output)
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
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/update", updateHandler(config))

	log.Printf("Server running on port %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}


