# Bookmark Service Deployment

Deployment repository for the Bookmark Service platform. This repository contains the infrastructure configuration required to run the API, Redis, Nginx reverse proxy, and portal together using Docker Compose.

This repository does not include the application source code. It is intended to be used alongside the application repository in [bookmark-service](bookmark-service).

---

## Table of Contents

- [System Architecture](#system-architecture)
- [Components](#components)
- [Repository Structure](#repository-structure)
- [Requirements](#requirements)
- [Setup & Deployment](#setup--deployment)
- [Environment Variables](#environment-variables)
- [Nginx Routing](#nginx-routing)
- [Post-Deployment Verification](#post-deployment-verification)
- [Common Operations](#common-operations)
- [Scaling](#scaling)
- [Troubleshooting](#troubleshooting)
- [Security Notes](#security-notes)
- [License](#license)

---

## System Architecture

The deployment exposes a single entry point through Nginx, which routes traffic to the portal frontend and to the Bookmark Service API. The API is deployed as three separate instances for load balancing, while Redis provides persistent storage for shortened URLs.

<img width="1688" height="969" alt="image" src="https://github.com/user-attachments/assets/29748fc4-edce-45a9-adac-3eecbbfb54fe" />

### Request Flow

1. The client sends requests to Nginx on port 80.
2. Nginx routes requests based on the URL path:
   - `/` is forwarded to the portal service.
   - `/api/bookmark_service/` is forwarded to the backend API upstream.
3. Each API instance processes the request and reads/writes data to Redis.
4. The portal frontend consumes the API through the same Nginx entry point.

---

## Components

| Service           | Image                         | Role                                           | Internal Port |
| ----------------- | ----------------------------- | ---------------------------------------------- | ------------- |
| nginx             | nginx:alpine                  | Public reverse proxy and load balancer         | 80            |
| bookmark_service  | khaivutri/shorten_link:test   | API instance #1                                | 8080          |
| bookmark_service2 | khaivutri/shorten_link:test   | API instance #2                                | 8080          |
| bookmark_service3 | khaivutri/shorten_link:test   | API instance #3                                | 8080          |
| redis             | redis:alpine                  | Short-link storage and health-check dependency | 6379          |
| portal            | ebvn/bookmark-app-portal:mono | Frontend portal application                    | 3000          |

> The three API services share the same image and the same environment file so their runtime configuration stays consistent.

---

## Repository Structure

The compose file is located in [deployment/docker-compose.yaml](deployment/docker-compose.yaml), and the Nginx configuration is in [deployment/nginx/nginx.conf](deployment/nginx/nginx.conf). The runtime environment file is in [deployment/bookmark_service/.env](deployment/bookmark_service/.env).

```text
bookmark-deploy/
├── bookmark-service/              # Application source repository
├── deployment/
│   ├── docker-compose.yaml
│   ├── bookmark_service/
│   │   └── .env
│   └── nginx/
│       └── nginx.conf
├── LICENSE
└── README.md
```

---

## Requirements

- Docker Engine 20.x or newer
- Docker Compose v2 or Docker Compose Classic
- Git
- The application repository [bookmark-service](bookmark-service) located in the same parent directory as this deployment repository

---

## Setup & Deployment

### 1. Clone the repositories

```bash
mkdir workspace && cd workspace
git clone https://github.com/khaivutri/bookmark-service.git
git clone https://github.com/khaivutri/bookmark-service-deployment.git
cd bookmark-service-deployment
```

### 2. Review the environment configuration

The deployment uses [deployment/bookmark_service/.env](deployment/bookmark_service/.env) for the API services.

```bash
nano deployment/bookmark_service/.env
```

Ensure the values match your environment, especially:

- REDIS_ADDRESS
- BASE_PATH

### 3. Start the full stack

```bash
docker compose -f deployment/docker-compose.yaml up --build -d
```

### 4. Follow the logs

```bash
docker compose -f deployment/docker-compose.yaml logs -f
```

### 5. Stop the stack

```bash
docker compose -f deployment/docker-compose.yaml down
```

---

## Environment Variables

The shared environment file for all API instances is [deployment/bookmark_service/.env](deployment/bookmark_service/.env).

| Variable      | Example Value         | Description                                        |
| ------------- | --------------------- | -------------------------------------------------- |
| REDIS_ADDRESS | redis:6379            | Redis hostname and port used by the API containers |
| BASE_PATH     | /api/bookmark_service | Base path expected by the application behind Nginx |

Additional variables can be added here if the application requires them.

---

## Nginx Routing

Nginx listens on port 80 and forwards requests as follows:

| Path                   | Forwarded To              | Upstream                                                              |
| ---------------------- | ------------------------- | --------------------------------------------------------------------- |
| /                      | portal                    | portal:3000                                                           |
| /api/bookmark_service/ | bookmark_service upstream | bookmark_service:8080, bookmark_service2:8080, bookmark_service3:8080 |

Example requests through Nginx:

```bash
# Health check
curl http://localhost/api/bookmark_service/health-check

# Create a short link
curl -X POST http://localhost/api/bookmark_service/v1/links/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","exp":60}'

# Redirect a short link
curl -i http://localhost/api/bookmark_service/v1/links/redirect/{code}
```

---

## Post-Deployment Verification

After the stack starts successfully, verify the following:

1. Open http://localhost/ in a browser to confirm the portal is reachable.
2. Run the health check:

```bash
curl http://localhost/api/bookmark_service/health-check
```

Expected result: HTTP 200 and a response indicating Redis is healthy.

3. Confirm container status:

```bash
docker compose -f deployment/docker-compose.yaml ps
```

All services should be listed as running.

---

## Common Operations

| Task                       | Command                                                                                                             |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| Start the stack            | docker compose -f deployment/docker-compose.yaml up --build -d                                                      |
| View logs                  | docker compose -f deployment/docker-compose.yaml logs -f                                                            |
| View logs for one service  | docker compose -f deployment/docker-compose.yaml logs -f nginx                                                      |
| Restart one service        | docker compose -f deployment/docker-compose.yaml restart bookmark_service                                           |
| Rebuild the API containers | docker compose -f deployment/docker-compose.yaml up --build -d bookmark_service bookmark_service2 bookmark_service3 |
| Stop the stack             | docker compose -f deployment/docker-compose.yaml down                                                               |
| Stop and remove volumes    | docker compose -f deployment/docker-compose.yaml down -v                                                            |

---

## Scaling

The current deployment uses three separate API services so Nginx can target each container explicitly by DNS name. To add another instance:

1. Add a new service in [deployment/docker-compose.yaml](deployment/docker-compose.yaml) following the same pattern as the existing API services.
2. Add a matching upstream server entry in [deployment/nginx/nginx.conf](deployment/nginx/nginx.conf).
3. Recreate the stack with the same compose command.

---

## Troubleshooting

| Symptom                                           | Likely Cause                                                             | Suggested Fix                                                                                                                       |
| ------------------------------------------------- | ------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------- |
| 502 Bad Gateway from Nginx                        | One or more API containers are still starting or have crashed            | Check logs with docker compose -f deployment/docker-compose.yaml logs bookmark_service                                              |
| Health check reports Redis as down                | Redis is not ready yet or REDIS_ADDRESS is incorrect                     | Confirm the Redis container is running and review the value in [deployment/bookmark_service/.env](deployment/bookmark_service/.env) |
| Missing environment file inside the API container | The environment file was not created or the path is incorrect            | Verify [deployment/bookmark_service/.env](deployment/bookmark_service/.env) exists and is referenced correctly by the compose file  |
| Nginx configuration issues                        | The mounted config file is not available or the container cannot read it | Confirm [deployment/nginx/nginx.conf](deployment/nginx/nginx.conf) exists and the volume mount is valid                             |

---

## Security Notes

- Do not commit real secrets or production credentials in [deployment/bookmark_service/.env](deployment/bookmark_service/.env).
- In production, restrict access to Redis and expose only the Nginx entry point to the public internet.
- Consider adding TLS termination at Nginx for secure HTTPS traffic.

---

## License

This repository is provided as deployment infrastructure for the Bookmark Service stack and may be adapted to fit your environment.
