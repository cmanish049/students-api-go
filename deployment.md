# Students API Deployment Guide

**Stack:** Go REST API + SQLite
**Environment:** Ubuntu 22.04 on AWS EC2
**Database Path:** `/var/lib/students-api/storage.db`
**HTTP Port:** `8082` (proxied through Nginx)

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [EC2 Setup](#ec2-setup)
3. [Project Directory Layout](#project-directory-layout)
4. [Install Go](#install-go)
5. [Clone and Build Application](#clone-and-build-application)
6. [Database Setup](#database-setup)
7. [systemd Service Setup](#systemd-service-setup)
8. [Nginx Reverse Proxy](#nginx-reverse-proxy)
9. [Service Verification](#service-verification)
10. [Updating the API](#updating-the-api)
11. [Rollback Procedure](#rollback-procedure)
12. [SQLite Best Practices](#sqlite-best-practices)
13. [Optional: Automating Deployment](#optional-automating-deployment)

---

## Prerequisites

* AWS account with an EC2 Ubuntu 22.04 instance
* Security Group allowing:

  * SSH (22)
  * HTTP (80)
  * HTTPS (443 – optional)
* SSH key for EC2 access
* Git installed locally and on EC2

---

## EC2 Setup

Connect to your EC2 instance:

```bash
chmod 400 your-key.pem
ssh -i your-key.pem ubuntu@<EC2_PUBLIC_IP>
```

Update system and install required packages:

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y git curl build-essential sqlite3 nginx
```

---

## Project Directory Layout

Recommended structure:

```text
/var/www/students-api/       # Application code
├── app                     # compiled binary
├── config/
│   └── production.yaml
├── go.mod
└── go.sum

/var/lib/students-api/       # Persistent SQLite database
└── storage.db
```

Create directories and set permissions:

```bash
sudo mkdir -p /var/www/students-api
sudo mkdir -p /var/lib/students-api
sudo chown -R ubuntu:ubuntu /var/www/students-api
sudo chown -R ubuntu:ubuntu /var/lib/students-api
```

---

## Install Go

Download and install the latest Go version:

```bash
curl -OL https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
```

Add Go to your PATH:

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
go version
```

---

## Clone and Build Application

```bash
cd /var/www
git clone https://github.com/yourname/students-api.git
cd students-api
go mod tidy
go mod download
go build -o app ./cmd/server  # Adjust path if main.go is elsewhere
```

---

## Database Setup

Create SQLite DB if it doesn’t exist:

```bash
sqlite3 /var/lib/students-api/storage.db
sqlite> PRAGMA journal_mode=WAL;
sqlite> .exit
```

Set correct ownership and permissions:

```bash
sudo chown -R ubuntu:ubuntu /var/lib/students-api
chmod 664 /var/lib/students-api/storage.db
chmod 775 /var/lib/students-api
```

---

## systemd Service Setup

Create systemd service:

```bash
sudo nano /etc/systemd/system/students-api.service
```

```ini
[Unit]
Description=Students API (Go + SQLite)
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/var/www/students-api
ExecStart=/var/www/students-api/app -config config/production.yaml
Restart=always
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

Enable and start service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable students-api
sudo systemctl start students-api
sudo systemctl status students-api
```

---

## Nginx Reverse Proxy

Create Nginx site configuration:

```bash
sudo nano /etc/nginx/sites-available/students-api
```

```nginx
server {
    listen 80;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:8082;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

Enable site and reload:

```bash
sudo ln -s /etc/nginx/sites-available/students-api /etc/nginx/sites-enabled/
sudo rm /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

---

## Service Verification

Test service locally:

```bash
curl http://127.0.0.1:8082/health
```

Test externally via EC2 public IP:

```bash
curl http://<EC2_PUBLIC_IP>
```

Check logs:

```bash
journalctl -u students-api -f
```

---

## Updating the API

1. Pull new changes:

```bash
cd /var/www/students-api
git pull origin main
```

2. Install new dependencies:

```bash
go mod tidy
go mod download
```

3. Build updated binary:

```bash
go build -o app ./cmd/server
```

4. Restart systemd service:

```bash
sudo systemctl restart students-api
sudo systemctl status students-api
journalctl -u students-api -f
```

5. Verify API endpoints.

---

## Rollback Procedure

If the new build fails:

```bash
mv app.old app
sudo systemctl restart students-api
```

SQLite DB is unaffected.

---

## SQLite Best Practices

* Always use **WAL mode**:

```sql
PRAGMA journal_mode=WAL;
```

* Ensure `/var/lib/students-api` is **writable** by the service user
* Backup regularly:

```bash
sqlite3 /var/lib/students-api/storage.db ".backup /var/lib/students-api/backup-$(date +%F).sqlite"
```

* Avoid long-running transactions to prevent locks

---

## Optional: Automating Deployment

Create `deploy.sh`:

```bash
#!/bin/bash
cd /var/www/students-api
git pull origin main
go mod tidy
go mod download
go build -o app ./cmd/server
sudo systemctl restart students-api
```

Make it executable:

```bash
chmod +x deploy.sh
```

Run to deploy updates in one command:

```bash
./deploy.sh
```

---

This workflow ensures **safe deployments, easy updates, and production reliability**.

---

If you want, I can also create an **enhanced version of this Markdown** that includes:

* Automatic **SQLite backup before each deploy**
* Automatic **migration step**
* Optional **Docker alternative**

This would make your EC2 deployment almost fully hands-off.

Do you want me to create that enhanced version?
