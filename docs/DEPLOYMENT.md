# Deployment Guide

## 📋 Table of Contents
- [Overview](#overview)
- [Environment Requirements](#environment-requirements)
- [Docker Deployment](#docker-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Production Configuration](#production-configuration)
- [Monitoring Setup](#monitoring-setup)
- [Security Considerations](#security-considerations)
- [Maintenance](#maintenance)

## 🌐 Overview

This guide covers deploying the Go MVC application to production environments. The application supports multiple deployment strategies:

- **Docker Compose**: Single-server deployment
- **Kubernetes**: Scalable container orchestration
- **Traditional**: Direct server deployment
- **Cloud Platforms**: AWS, GCP, Azure

## 🔧 Environment Requirements

### Minimum System Requirements

| Component | Requirement |
|-----------|-------------|
| **CPU** | 2 cores |
| **Memory** | 4GB RAM |
| **Storage** | 20GB SSD |
| **Network** | 1Gbps |

### Production Requirements

| Component | Requirement |
|-----------|-------------|
| **CPU** | 4+ cores |
| **Memory** | 8GB+ RAM |
| **Storage** | 100GB+ SSD |
| **Network** | 10Gbps |

### Software Dependencies

```bash
# Required
Docker 20.0+
Docker Compose 2.0+

# Optional (for Kubernetes)
Kubernetes 1.20+
Helm 3.0+

# Monitoring
Prometheus
Grafana
Jaeger
```

## 🐳 Docker Deployment

### Single Server Deployment

#### 1. Prepare Environment

```bash
# Clone repository
git clone https://github.com/tranvuongduy2003/go-mvc.git
cd go-mvc

# Create production environment file
cp .env.example .env.production

# Edit production settings
vim .env.production
```

#### 2. Production Environment Variables

```bash
# .env.production
ENV=production

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_TIMEOUT=30s

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=go_mvc_prod
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_SSL_MODE=require

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# Security
JWT_SECRET=your-super-secure-jwt-secret-key
ENCRYPTION_KEY=your-32-character-encryption-key

# External Services
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=noreply@example.com
SMTP_PASSWORD=smtp-password

# Monitoring
PROMETHEUS_ENABLED=true
JAEGER_ENABLED=true
JAEGER_ENDPOINT=http://jaeger:14268/api/traces

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

#### 3. Production Docker Compose

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    image: go-mvc:latest
    container_name: go-mvc-app
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - ENV=production
    env_file:
      - .env.production
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    container_name: go-mvc-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: go_mvc_prod
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: go-mvc-redis
    restart: unless-stopped
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Reverse Proxy
  nginx:
    image: nginx:alpine
    container_name: go-mvc-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./configs/nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    networks:
      - app-network

  # Monitoring Stack
  prometheus:
    image: prom/prometheus:latest
    container_name: go-mvc-prometheus
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./configs/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'
    networks:
      - app-network

  grafana:
    image: grafana/grafana:latest
    container_name: go-mvc-grafana
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
      - ./configs/grafana/provisioning:/etc/grafana/provisioning
      - ./configs/grafana/dashboards:/var/lib/grafana/dashboards
    networks:
      - app-network

  jaeger:
    image: jaegertracing/all-in-one:latest
    container_name: go-mvc-jaeger
    restart: unless-stopped
    ports:
      - "16686:16686"
      - "14268:14268"
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:

networks:
  app-network:
    driver: bridge
```

#### 4. Nginx Configuration

Create `configs/nginx.conf`:

```nginx
events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;
    }

    server {
        listen 80;
        server_name your-domain.com;
        
        # Redirect HTTP to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        # SSL Configuration
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
        ssl_prefer_server_ciphers off;

        # Security Headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header Referrer-Policy "no-referrer-when-downgrade" always;
        add_header Content-Security-Policy "default-src 'self'" always;

        # Gzip Compression
        gzip on;
        gzip_vary on;
        gzip_min_length 1024;
        gzip_proxied expired no-cache no-store private auth;
        gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;

        # Rate Limiting
        limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
        limit_req zone=api burst=20 nodelay;

        # Proxy Settings
        location / {
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            # Timeouts
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Health check
        location /health {
            proxy_pass http://app/health;
            access_log off;
        }

        # Static files (if any)
        location /static/ {
            alias /var/www/static/;
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

#### 5. Deploy to Production

```bash
# Build and deploy
docker-compose -f docker-compose.prod.yml up -d --build

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f app

# Run migrations
docker-compose -f docker-compose.prod.yml exec app ./bin/migrate up
```

### Docker Swarm Deployment

#### 1. Initialize Swarm

```bash
# Initialize swarm
docker swarm init

# Create overlay network
docker network create --driver overlay app-network
```

#### 2. Docker Stack File

Create `docker-stack.yml`:

```yaml
version: '3.8'

services:
  app:
    image: go-mvc:latest
    deploy:
      replicas: 3
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      update_config:
        parallelism: 1
        delay: 10s
        failure_action: rollback
      resources:
        limits:
          cpus: '1'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
    ports:
      - "8080:8080"
    networks:
      - app-network
    secrets:
      - db_password
      - jwt_secret

secrets:
  db_password:
    external: true
  jwt_secret:
    external: true

networks:
  app-network:
    external: true
```

#### 3. Deploy Stack

```bash
# Create secrets
echo "your-db-password" | docker secret create db_password -
echo "your-jwt-secret" | docker secret create jwt_secret -

# Deploy stack
docker stack deploy -c docker-stack.yml go-mvc

# Check services
docker service ls
docker service logs go-mvc_app
```

## ☸️ Kubernetes Deployment

### 1. Namespace and ConfigMap

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: go-mvc

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-mvc-config
  namespace: go-mvc
data:
  ENV: "production"
  SERVER_HOST: "0.0.0.0"
  SERVER_PORT: "8080"
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
```

### 2. Secrets

```yaml
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: go-mvc-secrets
  namespace: go-mvc
type: Opaque
data:
  DB_PASSWORD: <base64-encoded-password>
  JWT_SECRET: <base64-encoded-jwt-secret>
  REDIS_PASSWORD: <base64-encoded-redis-password>
```

### 3. PostgreSQL Deployment

```yaml
# k8s/postgres.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: go-mvc
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        env:
        - name: POSTGRES_DB
          value: go_mvc_prod
        - name: POSTGRES_USER
          value: postgres
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: go-mvc-secrets
              key: DB_PASSWORD
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1"
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
  namespace: go-mvc
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: go-mvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
```

### 4. Application Deployment

```yaml
# k8s/app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-mvc-app
  namespace: go-mvc
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-mvc-app
  template:
    metadata:
      labels:
        app: go-mvc-app
    spec:
      containers:
      - name: go-mvc
        image: go-mvc:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: go-mvc-config
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: go-mvc-secrets
              key: DB_PASSWORD
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: go-mvc-secrets
              key: JWT_SECRET
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: go-mvc-service
  namespace: go-mvc
spec:
  selector:
    app: go-mvc-app
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-mvc-ingress
  namespace: go-mvc
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - your-domain.com
    secretName: go-mvc-tls
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-mvc-service
            port:
              number: 80
```

### 5. Deploy to Kubernetes

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n go-mvc
kubectl get services -n go-mvc
kubectl get ingress -n go-mvc

# View logs
kubectl logs -f deployment/go-mvc-app -n go-mvc

# Scale deployment
kubectl scale deployment go-mvc-app --replicas=5 -n go-mvc
```

### 6. Helm Chart (Optional)

Create a Helm chart for easier management:

```bash
# Create Helm chart
helm create go-mvc-chart

# Install chart
helm install go-mvc ./go-mvc-chart

# Upgrade
helm upgrade go-mvc ./go-mvc-chart

# Rollback
helm rollback go-mvc 1
```

## ⚙️ Production Configuration

### Application Configuration

```yaml
# configs/production.yaml
server:
  host: "0.0.0.0"
  port: 8080
  timeout: 30s
  graceful_timeout: 15s

database:
  host: "${DB_HOST}"
  port: ${DB_PORT}
  name: "${DB_NAME}"
  user: "${DB_USER}"
  password: "${DB_PASSWORD}"
  ssl_mode: "require"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: "5m"

redis:
  host: "${REDIS_HOST}"
  port: ${REDIS_PORT}
  password: "${REDIS_PASSWORD}"
  db: ${REDIS_DB}
  pool_size: 10
  min_idle_conns: 5

logging:
  level: "info"
  format: "json"
  output: "stdout"

security:
  jwt:
    secret: "${JWT_SECRET}"
    expiry: "24h"
  cors:
    allowed_origins: ["https://your-domain.com"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allowed_headers: ["*"]

monitoring:
  prometheus:
    enabled: true
    port: 9090
  jaeger:
    enabled: true
    endpoint: "${JAEGER_ENDPOINT}"
    service_name: "go-mvc"

external_services:
  smtp:
    host: "${SMTP_HOST}"
    port: ${SMTP_PORT}
    user: "${SMTP_USER}"
    password: "${SMTP_PASSWORD}"
```

### Environment Variables Management

#### Using Docker Secrets

```bash
# Create secrets
echo "db-password" | docker secret create db_password -
echo "jwt-secret" | docker secret create jwt_secret -

# Use in Docker Compose
services:
  app:
    secrets:
      - db_password
      - jwt_secret
    environment:
      - DB_PASSWORD_FILE=/run/secrets/db_password
      - JWT_SECRET_FILE=/run/secrets/jwt_secret
```

#### Using Kubernetes Secrets

```bash
# Create secrets from command line
kubectl create secret generic go-mvc-secrets \
  --from-literal=DB_PASSWORD=your-password \
  --from-literal=JWT_SECRET=your-jwt-secret \
  -n go-mvc

# Create from file
kubectl create secret generic go-mvc-secrets \
  --from-env-file=.env.production \
  -n go-mvc
```

## 📊 Monitoring Setup

### Prometheus Configuration

```yaml
# configs/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'go-mvc'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']

  - job_name: 'nginx'
    static_configs:
      - targets: ['nginx-exporter:9113']
```

### Grafana Dashboard

Dashboard includes:
- Application metrics (requests, latency, errors)
- System metrics (CPU, memory, disk)
- Database metrics (connections, queries)
- Custom business metrics

### Alerting Rules

```yaml
# configs/alert-rules.yml
groups:
  - name: go-mvc-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"

      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"

      - alert: DatabaseDown
        expr: up{job="postgres"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database is down"
```

## 🔒 Security Considerations

### SSL/TLS Configuration

```bash
# Generate SSL certificates (Let's Encrypt)
certbot certonly --standalone -d your-domain.com

# Or use cert-manager in Kubernetes
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.8.0/cert-manager.yaml
```

### Security Headers

Nginx configuration includes:
- HSTS (HTTP Strict Transport Security)
- XSS Protection
- Content Security Policy
- X-Frame-Options
- X-Content-Type-Options

### Network Security

```bash
# Docker network isolation
docker network create --driver bridge --internal app-internal

# Kubernetes NetworkPolicies
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: go-mvc-netpol
spec:
  podSelector:
    matchLabels:
      app: go-mvc-app
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: nginx
```

### Secret Management

- Use external secret management (HashiCorp Vault, AWS Secrets Manager)
- Rotate secrets regularly
- Never log sensitive information
- Use secure environment variable injection

## 🔧 Maintenance

### Health Checks

```bash
# Application health
curl -f http://your-domain.com/health

# Database health
docker exec postgres pg_isready -U postgres

# Redis health
docker exec redis redis-cli ping
```

### Backup Strategy

#### Database Backup

```bash
# Automated backup script
#!/bin/bash
BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup
docker exec postgres pg_dump -U postgres go_mvc_prod > $BACKUP_DIR/backup_$DATE.sql

# Compress
gzip $BACKUP_DIR/backup_$DATE.sql

# Upload to S3 (optional)
aws s3 cp $BACKUP_DIR/backup_$DATE.sql.gz s3://your-backup-bucket/

# Cleanup old backups (keep 30 days)
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete
```

#### Application Data Backup

```bash
# Volume backup
docker run --rm -v go-mvc_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz /data
```

### Log Management

```yaml
# Logging configuration
version: '3.8'
services:
  app:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Update Strategy

#### Rolling Updates

```bash
# Docker Compose
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d --no-deps app

# Kubernetes
kubectl set image deployment/go-mvc-app go-mvc=go-mvc:v2.0.0 -n go-mvc
kubectl rollout status deployment/go-mvc-app -n go-mvc
```

#### Blue-Green Deployment

```bash
# Deploy to staging environment first
docker-compose -f docker-compose.staging.yml up -d

# Test staging
curl -f http://staging.your-domain.com/health

# Switch traffic (update load balancer/DNS)
# Monitor and rollback if needed
```

### Troubleshooting

#### Common Issues

1. **Application Won't Start**
   ```bash
   # Check logs
   docker logs go-mvc-app
   
   # Check environment variables
   docker exec go-mvc-app env | grep -E "(DB_|REDIS_|JWT_)"
   ```

2. **Database Connection Issues**
   ```bash
   # Test database connectivity
   docker exec go-mvc-app nc -zv postgres 5432
   
   # Check database logs
   docker logs go-mvc-postgres
   ```

3. **High Memory Usage**
   ```bash
   # Monitor memory
   docker stats
   
   # Check application metrics
   curl http://localhost:8080/metrics | grep memory
   ```

#### Performance Tuning

```bash
# Optimize database connections
# Adjust connection pool settings in production.yaml

# Enable HTTP/2
# Update nginx configuration

# Use CDN for static assets
# Configure cache headers
```

This deployment guide ensures a robust, scalable, and secure production deployment of the Go MVC application.