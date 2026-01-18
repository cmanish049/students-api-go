# Students API

A RESTful API service for managing student records built with Go and SQLite. This project demonstrates clean architecture principles with clear separation of concerns between HTTP handlers, business logic, and data storage.

## Features

- ✅ Create, Read, Update, and Delete (CRUD) operations for students
- ✅ SQLite database for data persistence
- ✅ Request validation using validator/v10
- ✅ Structured JSON responses
- ✅ YAML-based configuration
- ✅ Graceful server shutdown
- ✅ Structured logging with slog
- ✅ Clean architecture with dependency injection

## Tech Stack

- **Language**: Go 1.25.5
- **Database**: SQLite
- **HTTP Router**: Go standard library (http.ServeMux)
- **Configuration**: cleanenv (YAML/ENV)
- **Validation**: go-playground/validator
- **Logging**: Go standard library (log/slog)

## Project Structure

```
students-api/
├── cmd/
│   └── students-api/
│       └── main.go              # Application entry point
├── config/
│   └── local.yaml               # Local configuration file
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading logic
│   ├── http/
│   │   └── handlers/
│   │       └── student/
│   │           └── student.go   # HTTP handlers for student operations
│   ├── storage/
│   │   ├── storage.go           # Storage interface definition
│   │   ├── postgres/            # PostgreSQL implementation (placeholder)
│   │   └── sqlite/
│   │       └── sqlite.go        # SQLite implementation
│   ├── types/
│   │   └── types.go             # Data type definitions
│   └── utils/
│       └── response/
│           └── response.go      # HTTP response utilities
├── storage/                     # Database file location
├── go.mod                       # Go module dependencies
└── README.md                    # This file
```

## Installation

### Prerequisites

- Go 1.25.5 or higher
- SQLite3

### Setup

1. Clone the repository:
```bash
git clone https://github.com/cmanish049/students-api.git
cd students-api
```

2. Install dependencies:
```bash
go mod download
```

3. Create the storage directory:
```bash
mkdir -p storage
```

## Configuration

The application uses YAML configuration files. Create or modify `config/local.yaml`:

```yaml
env: "dev"
storage_path: "storage/storage.db"
http_server:
  address: "localhost:8082"
```

### Configuration Options

- `env`: Environment name (dev, production)
- `storage_path`: Path to SQLite database file
- `http_server.address`: Server address and port

### Configuration Loading

The application loads configuration in the following priority:

1. `CONFIG_PATH` environment variable
2. `--config` command-line flag

Example:
```bash
# Using environment variable
export CONFIG_PATH=config/local.yaml
go run cmd/students-api/main.go

# Using command-line flag
go run cmd/students-api/main.go --config=config/local.yaml
```

## Running the Application

### Development Mode

```bash
go run cmd/students-api/main.go --config=config/local.yaml
```

### Build and Run

```bash
# Build the binary
go build -o bin/students-api cmd/students-api/main.go

# Run the binary
./bin/students-api --config=config/local.yaml
```

The server will start on `http://localhost:8082` (or the address specified in your config file).

## Deployment

This section provides comprehensive guidance for deploying the Students API to production environments.

### Production Build

#### 1. Build the Binary

For production deployment, build a static binary optimized for your target platform:

```bash
# For Linux (most common for servers)
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o students-api cmd/students-api/main.go

# For macOS
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o students-api cmd/students-api/main.go

# For Windows
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o students-api.exe cmd/students-api/main.go
```

**Note**: `CGO_ENABLED=1` is required for SQLite support.

**Build flags explanation**:
- `-ldflags="-s -w"`: Strips debug information, reducing binary size
- `-o`: Specifies output binary name

#### 2. Create Production Configuration

Create a production configuration file `config/production.yaml`:

```yaml
env: "production"
storage_path: "/var/lib/students-api/storage.db"
http_server:
  address: "0.0.0.0:8082"
```

### Deployment Options

#### Option 1: Traditional VPS/Server Deployment

##### Step 1: Prepare the Server

```bash
# Update system packages
sudo apt update && sudo apt upgrade -y

# Install SQLite (if not already installed)
sudo apt install sqlite3 -y

# Create application user
sudo useradd -r -s /bin/false students-api

# Create application directories
sudo mkdir -p /opt/students-api
sudo mkdir -p /var/lib/students-api
sudo mkdir -p /var/log/students-api

# Set permissions
sudo chown -R students-api:students-api /opt/students-api
sudo chown -R students-api:students-api /var/lib/students-api
sudo chown -R students-api:students-api /var/log/students-api
```

##### Step 2: Upload Application Files

```bash
# From your local machine, upload the binary
scp students-api user@your-server:/tmp/
scp -r config user@your-server:/tmp/

# On the server, move files to appropriate locations
sudo mv /tmp/students-api /opt/students-api/
sudo mv /tmp/config /opt/students-api/
sudo chmod +x /opt/students-api/students-api
sudo chown -R students-api:students-api /opt/students-api
```

##### Step 3: Create systemd Service

Create a systemd service file `/etc/systemd/system/students-api.service`:

```ini
[Unit]
Description=Students API Service
After=network.target

[Service]
Type=simple
User=students-api
Group=students-api
WorkingDirectory=/opt/students-api
ExecStart=/opt/students-api/students-api --config=/opt/students-api/config/production.yaml
Restart=on-failure
RestartSec=5s

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/students-api /var/log/students-api

# Environment variables (optional)
Environment="CONFIG_PATH=/opt/students-api/config/production.yaml"

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=students-api

[Install]
WantedBy=multi-user.target
```

##### Step 4: Start and Enable Service

```bash
# Reload systemd daemon
sudo systemctl daemon-reload

# Start the service
sudo systemctl start students-api

# Enable service to start on boot
sudo systemctl enable students-api

# Check service status
sudo systemctl status students-api

# View logs
sudo journalctl -u students-api -f
```

##### Step 5: Configure Nginx Reverse Proxy

Install and configure Nginx:

```bash
# Install Nginx
sudo apt install nginx -y

# Create Nginx configuration
sudo nano /etc/nginx/sites-available/students-api
```

Nginx configuration (`/etc/nginx/sites-available/students-api`):

```nginx
upstream students_api {
    server 127.0.0.1:8082;
}

server {
    listen 80;
    server_name api.yourdomain.com;

    # Request size limit
    client_max_body_size 10M;

    # Logging
    access_log /var/log/nginx/students-api-access.log;
    error_log /var/log/nginx/students-api-error.log;

    # Security headers
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-XSS-Protection "1; mode=block" always;

    location / {
        proxy_pass http://students_api;
        proxy_http_version 1.1;
        
        # Proxy headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint
    location /health {
        access_log off;
        proxy_pass http://students_api;
    }
}
```

Enable the site:

```bash
# Create symbolic link
sudo ln -s /etc/nginx/sites-available/students-api /etc/nginx/sites-enabled/

# Test Nginx configuration
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx
```

##### Step 6: Configure SSL with Let's Encrypt

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain SSL certificate
sudo certbot --nginx -d api.yourdomain.com

# Certbot will automatically configure SSL in Nginx

# Test automatic renewal
sudo certbot renew --dry-run
```

After SSL setup, your Nginx config will be automatically updated with HTTPS configuration.

#### Option 2: Docker Deployment

##### Step 1: Create Dockerfile

Create `Dockerfile` in the project root:

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o students-api cmd/students-api/main.go

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache sqlite-libs ca-certificates

# Create non-root user
RUN addgroup -g 1000 students && \
    adduser -D -u 1000 -G students students

# Create necessary directories
RUN mkdir -p /app/config /app/storage && \
    chown -R students:students /app

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/students-api .

# Copy configuration files
COPY --chown=students:students config/ ./config/

# Switch to non-root user
USER students

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

# Run the application
CMD ["./students-api", "--config=config/production.yaml"]
```

##### Step 2: Create .dockerignore

Create `.dockerignore`:

```
storage/
.git
.gitignore
README.md
*.md
bin/
.env
*.log
```

##### Step 3: Create docker-compose.yml

```yaml
version: '3.8'

services:
  students-api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: students-api
    restart: unless-stopped
    ports:
      - "8082:8082"
    volumes:
      - ./storage:/app/storage
      - ./config:/app/config:ro
    environment:
      - CONFIG_PATH=/app/config/production.yaml
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8082/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
    networks:
      - students-api-network

networks:
  students-api-network:
    driver: bridge
```

##### Step 4: Build and Run with Docker

```bash
# Build the Docker image
docker build -t students-api:latest .

# Run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

##### Step 5: Docker Production Deployment

For production, use Docker with a reverse proxy:

```yaml
version: '3.8'

services:
  students-api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: students-api
    restart: unless-stopped
    expose:
      - "8082"
    volumes:
      - students-data:/app/storage
      - ./config:/app/config:ro
    environment:
      - CONFIG_PATH=/app/config/production.yaml
    networks:
      - students-api-network

  nginx:
    image: nginx:alpine
    container_name: students-api-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - students-api
    networks:
      - students-api-network

volumes:
  students-data:
    driver: local

networks:
  students-api-network:
    driver: bridge
```

#### Option 3: Cloud Platform Deployment

##### AWS EC2

1. **Launch EC2 Instance**:
   - Choose Ubuntu Server 22.04 LTS
   - Instance type: t2.micro (free tier) or larger
   - Configure security group: Allow HTTP (80), HTTPS (443), SSH (22)

2. **Follow VPS deployment steps** (Option 1) above

3. **Configure AWS-specific settings**:
   ```bash
   # Install CloudWatch agent for monitoring
   wget https://s3.amazonaws.com/amazoncloudwatch-agent/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb
   sudo dpkg -i amazon-cloudwatch-agent.deb
   ```

##### Google Cloud Platform (GCP)

1. **Create Compute Engine instance**:
   ```bash
   gcloud compute instances create students-api \
     --image-family=ubuntu-2204-lts \
     --image-project=ubuntu-os-cloud \
     --machine-type=e2-micro \
     --tags=http-server,https-server
   ```

2. **Configure firewall**:
   ```bash
   gcloud compute firewall-rules create allow-http \
     --allow tcp:80 --target-tags http-server
   
   gcloud compute firewall-rules create allow-https \
     --allow tcp:443 --target-tags https-server
   ```

3. **Follow VPS deployment steps**

##### DigitalOcean Droplet

1. **Create Droplet**:
   - Choose Ubuntu 22.04
   - Select appropriate size (Basic $6/month recommended)
   - Enable monitoring

2. **Follow VPS deployment steps**

##### Heroku (Using Docker)

Create `heroku.yml`:

```yaml
build:
  docker:
    web: Dockerfile
run:
  web: ./students-api --config=config/production.yaml
```

Deploy:

```bash
# Login to Heroku
heroku login

# Create app
heroku create your-students-api

# Set stack to container
heroku stack:set container

# Deploy
git push heroku main

# View logs
heroku logs --tail
```

### Database Backup and Recovery

#### Automated Backup Script

Create `/opt/students-api/backup.sh`:

```bash
#!/bin/bash

# Configuration
DB_PATH="/var/lib/students-api/storage.db"
BACKUP_DIR="/var/backups/students-api"
RETENTION_DAYS=7

# Create backup directory
mkdir -p $BACKUP_DIR

# Create backup with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/students_api_$TIMESTAMP.db"

# Perform backup
sqlite3 $DB_PATH ".backup '$BACKUP_FILE'"

# Compress backup
gzip $BACKUP_FILE

# Remove old backups
find $BACKUP_DIR -name "*.db.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: $BACKUP_FILE.gz"
```

Make it executable and schedule:

```bash
# Make executable
sudo chmod +x /opt/students-api/backup.sh

# Add to crontab for daily backups at 2 AM
sudo crontab -e

# Add this line:
0 2 * * * /opt/students-api/backup.sh >> /var/log/students-api/backup.log 2>&1
```

#### Restore from Backup

```bash
# Stop the service
sudo systemctl stop students-api

# Restore database
gunzip -c /var/backups/students-api/students_api_TIMESTAMP.db.gz > /var/lib/students-api/storage.db

# Fix permissions
sudo chown students-api:students-api /var/lib/students-api/storage.db

# Start the service
sudo systemctl start students-api
```

### Monitoring and Logging

#### Application Logs

```bash
# View real-time logs
sudo journalctl -u students-api -f

# View logs from last hour
sudo journalctl -u students-api --since "1 hour ago"

# View logs with specific priority
sudo journalctl -u students-api -p err
```

#### System Monitoring

Install and configure monitoring tools:

```bash
# Install htop for process monitoring
sudo apt install htop -y

# Install netdata for comprehensive monitoring
bash <(curl -Ss https://my-netdata.io/kickstart.sh)
```

#### Log Rotation

Create `/etc/logrotate.d/students-api`:

```
/var/log/students-api/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 students-api students-api
    sharedscripts
    postrotate
        systemctl reload students-api > /dev/null 2>&1 || true
    endscript
}
```

### Security Best Practices

1. **Firewall Configuration**:
   ```bash
   # Install UFW
   sudo apt install ufw -y
   
   # Allow SSH
   sudo ufw allow 22/tcp
   
   # Allow HTTP and HTTPS
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   
   # Enable firewall
   sudo ufw enable
   ```

2. **Keep System Updated**:
   ```bash
   # Enable automatic security updates
   sudo apt install unattended-upgrades -y
   sudo dpkg-reconfigure --priority=low unattended-upgrades
   ```

3. **Fail2Ban for SSH Protection**:
   ```bash
   sudo apt install fail2ban -y
   sudo systemctl enable fail2ban
   sudo systemctl start fail2ban
   ```

4. **Disable Root Login**:
   ```bash
   sudo nano /etc/ssh/sshd_config
   # Set: PermitRootLogin no
   sudo systemctl restart sshd
   ```

### Performance Optimization

#### 1. Enable Response Compression in Nginx

Add to Nginx configuration:

```nginx
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types application/json text/plain text/css application/javascript;
```

#### 2. Connection Pooling

Already handled by SQLite driver, but ensure proper configuration in production.

#### 3. Rate Limiting in Nginx

```nginx
# Define rate limit zone
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

server {
    # Apply rate limit
    location /api/ {
        limit_req zone=api_limit burst=20 nodelay;
        # ... rest of configuration
    }
}
```

### Health Checks and Uptime Monitoring

#### Create Health Check Endpoint

Add to your application (if not already present):

```go
router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
})
```

#### External Monitoring Services

- **UptimeRobot**: Free monitoring service (https://uptimerobot.com)
- **Pingdom**: Comprehensive monitoring (https://www.pingdom.com)
- **StatusCake**: Free tier available (https://www.statuscake.com)

### Troubleshooting

#### Service Won't Start

```bash
# Check service status
sudo systemctl status students-api

# Check logs for errors
sudo journalctl -u students-api -n 50

# Verify binary permissions
ls -la /opt/students-api/students-api

# Check configuration file
cat /opt/students-api/config/production.yaml
```

#### Database Permission Issues

```bash
# Fix database file permissions
sudo chown students-api:students-api /var/lib/students-api/storage.db
sudo chmod 644 /var/lib/students-api/storage.db

# Fix directory permissions
sudo chown students-api:students-api /var/lib/students-api
sudo chmod 755 /var/lib/students-api
```

#### High Memory Usage

```bash
# Check process memory
ps aux | grep students-api

# Monitor with htop
htop -p $(pgrep students-api)
```

### Continuous Deployment

#### GitHub Actions Example

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Production

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.23'
    
    - name: Build
      run: |
        CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o students-api cmd/students-api/main.go
    
    - name: Deploy to Server
      uses: appleboy/scp-action@master
      with:
        host: ${{ secrets.SERVER_HOST }}
        username: ${{ secrets.SERVER_USER }}
        key: ${{ secrets.SSH_PRIVATE_KEY }}
        source: "students-api,config/"
        target: "/tmp/"
    
    - name: Restart Service
      uses: appleboy/ssh-action@master
      with:
        host: ${{ secrets.SERVER_HOST }}
        username: ${{ secrets.SERVER_USER }}
        key: ${{ secrets.SSH_PRIVATE_KEY }}
        script: |
          sudo mv /tmp/students-api /opt/students-api/
          sudo chown students-api:students-api /opt/students-api/students-api
          sudo chmod +x /opt/students-api/students-api
          sudo systemctl restart students-api
```

## API Endpoints

### Student Model

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

### Endpoints

#### Create a Student

```http
POST /api/students
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

**Success Response** (201 Created):
```json
{
  "id": 1
}
```

**Error Response** (400 Bad Request):
```json
{
  "status": "Error",
  "error": "field Name is required field"
}
```

#### Get Student by ID

```http
GET /api/students/{id}
```

**Success Response** (200 OK):
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

#### Get All Students

```http
GET /api/students
```

**Success Response** (200 OK):
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 20
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "email": "jane@example.com",
    "age": 22
  }
]
```

#### Update a Student

```http
PUT /api/students/{id}
Content-Type: application/json

{
  "name": "John Updated",
  "email": "john.updated@example.com",
  "age": 21
}
```

**Success Response** (200 OK):
```json
{
  "message": "student updated successfully"
}
```

#### Delete a Student

```http
DELETE /api/students/{id}
```

**Success Response** (200 OK):
```json
{
  "message": "student deleted successfully"
}
```

## Testing with cURL

### Create a student
```bash
curl -X POST http://localhost:8082/api/students \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","age":20}'
```

### Get all students
```bash
curl http://localhost:8082/api/students
```

### Get student by ID
```bash
curl http://localhost:8082/api/students/1
```

### Update a student
```bash
curl -X PUT http://localhost:8082/api/students/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com","age":21}'
```

### Delete a student
```bash
curl -X DELETE http://localhost:8082/api/students/1
```

## Validation Rules

The API validates student data with the following rules:

- `name`: Required field
- `email`: Required field (must be unique in database)
- `age`: Required field

## Error Handling

The API returns consistent error responses:

```json
{
  "status": "Error",
  "error": "error description"
}
```

Common HTTP status codes:
- `200 OK`: Successful GET/PUT/DELETE operation
- `201 Created`: Successful POST operation
- `400 Bad Request`: Invalid input or validation error
- `500 Internal Server Error`: Server-side error

## Database Schema

The SQLite database contains a single `students` table:

```sql
CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    age INTEGER NOT NULL
);
```

## Architecture

### Clean Architecture Principles

The project follows clean architecture patterns:

1. **Handlers Layer** (`internal/http/handlers`): HTTP request/response handling
2. **Storage Interface** (`internal/storage`): Abstraction for data persistence
3. **Storage Implementation** (`internal/storage/sqlite`): Concrete database implementation
4. **Types** (`internal/types`): Domain models
5. **Config** (`internal/config`): Configuration management
6. **Utils** (`internal/utils`): Shared utilities

### Dependency Injection

The application uses dependency injection to maintain loose coupling:

```go
// Storage interface is injected into handlers
router.HandleFunc("POST /api/students", student.New(db))
```

This allows for easy testing and swapping of storage implementations (e.g., SQLite to PostgreSQL).

## Graceful Shutdown

The application implements graceful shutdown to handle SIGINT and SIGTERM signals:

- Active requests complete processing
- Server shutdown timeout: 5 seconds
- Database connections are properly closed

## Development

### Adding New Features

1. **Add new storage method**: Update `internal/storage/storage.go` interface
2. **Implement in SQLite**: Add method to `internal/storage/sqlite/sqlite.go`
3. **Create handler**: Add handler in `internal/http/handlers/student/student.go`
4. **Register route**: Add route in `cmd/students-api/main.go`

### Adding PostgreSQL Support

The project structure includes a `postgres/` directory for future PostgreSQL implementation:

1. Implement the `Storage` interface in `internal/storage/postgres/`
2. Update the main.go to switch between SQLite and PostgreSQL based on configuration

## Dependencies

- `github.com/go-playground/validator/v10`: Request validation
- `github.com/ilyakaznacheev/cleanenv`: Configuration management
- `github.com/mattn/go-sqlite3`: SQLite driver
- Go standard library for HTTP server and logging

## License

This project is available for educational and personal use.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

[@cmanish049](https://github.com/cmanish049)

## Roadmap

- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Implement PostgreSQL support
- [ ] Add authentication and authorization
- [ ] Add pagination for GET /api/students
- [ ] Add filtering and sorting
- [ ] Add API documentation with Swagger
- [ ] Add Docker support
- [ ] Add CI/CD pipeline
