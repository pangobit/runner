# Runner

A simple Go web server that executes docker-compose commands via webhook.

## Usage

1. Configure the server port and commands in `config.yaml`
2. Set the `RUNNER_AUTH_TOKEN` environment variable with a secure token:
   ```bash
   export RUNNER_AUTH_TOKEN="your-secure-token-here"
   ```
3. Run the server with `go run main.go`
4. The server exposes:
   - GET `/health` - Health check endpoint
   - POST `/update` - Webhook to trigger docker-compose commands (requires authentication)

### Authentication

The `/update` endpoint requires authentication using a Bearer token in the Authorization header. The token must match the value set in the `RUNNER_AUTH_TOKEN` environment variable.

Example authenticated request:
```bash
curl -H "Authorization: Bearer your-secure-token-here" http://localhost:8080/update
```

## Configuration

Edit `config.yaml` to set the port and list of docker-compose commands to run when the update endpoint is called.

## Docker

You can also build and run the application using Docker:

```bash
docker build -t runner-app .
docker run -v ./config:/app/config -v ./compose:/app/compose -p 8080:8080 runner-app
```

This will mount your local `config` directory to `/app/config` and `compose` directory to `/app/compose` in the container, allowing the application to access your configuration and docker-compose files.

## Docker Compose

Alternatively, use the provided `docker-compose.yml` to run the application:

```bash
docker-compose up -d
```

This will build and start the Runner application in the background with the same volume mounts as the Docker command.

**Important Note:** The Runner application executes `docker-compose` commands, which must be installed and available on the host system or in the environment where the commands are expected to run. If you are running the application in a Docker container, ensure that `docker-compose` is installed on the host, or consider mounting the Docker socket (`/var/run/docker.sock`) into the container to interact with the host's Docker daemon directly (with appropriate permissions and security considerations). The provided `docker-compose.yml` includes this mounting by default.

Be sure to read up on what mounting /var/run/docker.sock means from a security standpoint before running!

Here's a sample `docker-compose.yml` for running the Runner app:

```yaml
version: '3.8'
services:
  runner:
    image: ghcr.io/pangobit/runner:latest
    container_name: runner-app
    ports:
      - "${PORT:-8080}:8080"
    volumes:
      - ./config:/app/config
      - ./compose:/app/compose
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - PORT=${PORT:-8080}
    restart: unless-stopped
```

You can copy this configuration into your own `docker-compose.yml` file or modify it as needed to suit your environment. 