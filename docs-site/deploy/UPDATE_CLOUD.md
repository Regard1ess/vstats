# 更新 VStats Cloud 服务

## 问题诊断

如果遇到以下问题：
- `/api/release/version` 返回 404
- `/download/latest/*` 返回网页而不是文件
- Docker 容器镜像未更新

## 解决步骤

### 1. 更新 Nginx 配置

确保 `nginx/conf.d/default.conf` 包含 `/download/` 路由配置：

```nginx
# Download Endpoints (for binary distribution)
location /download/ {
    proxy_pass http://api_server;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header CF-Connecting-IP $http_cf_connecting_ip;
    proxy_set_header Connection "";

    # Longer timeouts for large binary downloads
    proxy_connect_timeout 120s;
    proxy_send_timeout 300s;
    proxy_read_timeout 300s;

    # Disable buffering for large files
    proxy_buffering off;
    proxy_request_buffering off;

    # Allow large files
    client_max_body_size 200M;
}
```

### 2. 拉取最新镜像并重启

```bash
cd /path/to/vstats/docs-site/deploy

# 方式1：使用部署脚本（推荐）
./scripts/deploy.sh update

# 方式2：手动更新
docker compose pull api
docker compose up -d --force-recreate api
docker compose restart nginx
```

### 3. 验证更新

```bash
# 检查容器状态
docker compose ps

# 检查 API 容器日志
docker compose logs -f api

# 验证版本信息
curl https://vstats.zsoft.cc/api/release/version

# 验证下载功能
curl -I https://vstats.zsoft.cc/download/latest/vstats-server-linux-amd64
```

### 4. 检查 API 容器镜像

```bash
# 查看当前运行的镜像版本
docker images | grep vstats-cloud

# 应该看到类似这样的输出：
# zsai001/vstats-cloud   latest   <IMAGE_ID>   X hours ago   XX MB
```

### 5. 强制重新构建（如果使用本地构建）

如果使用本地构建而不是 Docker Hub 镜像：

```bash
cd /path/to/vstats

# 构建 Cloud 镜像
docker build -f Dockerfile.cloud -t zsai001/vstats-cloud:latest .

# 或使用预构建版本
docker build -f Dockerfile.cloud.prebuilt -t zsai001/vstats-cloud:latest .

# 重启服务
cd docs-site/deploy
docker compose up -d --force-recreate api
```

## 常见问题

### Q1: 容器无法启动

```bash
# 检查日志
docker compose logs api

# 检查网络连接
docker compose exec api wget -q -O- http://localhost:3001/health
```

### Q2: Nginx 配置不生效

```bash
# 重新加载 Nginx 配置
docker compose restart nginx

# 或进入容器检查配置
docker compose exec nginx nginx -t
docker compose exec nginx nginx -s reload
```

### Q3: Redis 或 PostgreSQL 连接失败

```bash
# 检查数据库服务状态
docker compose ps postgres redis

# 检查连接
docker compose exec postgres psql -U vstats -d vstats_cloud -c "SELECT 1;"
docker compose exec redis redis-cli -a <password> ping
```

### Q4: 镜像版本不更新

```bash
# 清理旧镜像
docker image prune -a -f

# 强制拉取最新镜像
docker compose pull --ignore-pull-failures
docker compose up -d --force-recreate
```

## 完整重启流程

如果需要完全重启服务：

```bash
cd /path/to/vstats/docs-site/deploy

# 1. 停止所有服务
docker compose down

# 2. 拉取最新镜像
docker compose pull

# 3. 重新启动
docker compose up -d

# 4. 检查状态
docker compose ps
docker compose logs -f
```

## 健康检查

```bash
# 运行内置健康检查
./scripts/health-check.sh

# 或手动检查
curl https://vstats.zsoft.cc/health
curl https://vstats.zsoft.cc/health/detailed
curl https://vstats.zsoft.cc/version
curl https://vstats.zsoft.cc/api/release/version
```

## 监控日志

```bash
# 实时查看所有服务日志
docker compose logs -f

# 只查看 API 日志
docker compose logs -f api

# 只查看 Nginx 日志
docker compose logs -f nginx

# 查看最近 100 行
docker compose logs --tail=100 api
```
