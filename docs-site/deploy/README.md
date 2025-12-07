# VStats Cloud 部署指南

本目录包含 VStats Cloud 服务的部署配置文件。

## 目录结构

```
deploy/
├── docker-compose.yml    # Docker Compose 配置 (PostgreSQL + Redis)
├── db/
│   └── schema.sql        # 数据库初始化 Schema
├── nginx.conf            # Nginx 配置示例
├── env.example           # 环境变量示例
└── README.md             # 本文件
```

## 快速开始

### 1. 准备环境变量

```bash
# 复制环境变量示例文件
cp env.example .env

# 编辑配置（必须修改密码！）
vim .env
```

**重要：** 请务必修改以下配置：
- `POSTGRES_PASSWORD` - PostgreSQL 密码
- `REDIS_PASSWORD` - Redis 密码
- `SESSION_SECRET` - Session 密钥
- `JWT_SECRET` - JWT 密钥
- OAuth 相关配置

### 2. 启动数据库服务

```bash
# 启动 PostgreSQL 和 Redis
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f
```

### 3. 验证服务

```bash
# 测试 PostgreSQL 连接
docker exec -it vstats-postgres psql -U vstats -d vstats_cloud -c "SELECT 1;"

# 测试 Redis 连接
docker exec -it vstats-redis redis-cli -a your_redis_password ping
```

### 4. 配置 Nginx

```bash
# 复制 Nginx 配置
sudo cp nginx.conf /etc/nginx/sites-available/vstats-cloud

# 修改域名和证书路径
sudo vim /etc/nginx/sites-available/vstats-cloud

# 创建软链接
sudo ln -s /etc/nginx/sites-available/vstats-cloud /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载 Nginx
sudo systemctl reload nginx
```

### 5. 配置 SSL 证书 (Let's Encrypt)

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d vstats.example.com

# 自动续期 (已自动配置)
sudo certbot renew --dry-run
```

## 数据库 Schema

### 核心表结构

| 表名 | 说明 |
|------|------|
| `users` | 用户账户 |
| `oauth_providers` | OAuth 登录信息 (GitHub/Google) |
| `sessions` | 用户会话 |
| `servers` | 监控的服务器 |
| `server_metrics` | 服务器指标数据 |
| `alert_rules` | 告警规则 |
| `alert_history` | 告警历史 |
| `api_keys` | API 密钥 |
| `audit_logs` | 审计日志 |
| `subscriptions` | 订阅信息 |

### ER 图简述

```
users (1) ----< (N) oauth_providers
users (1) ----< (N) sessions
users (1) ----< (N) servers
users (1) ----< (N) alert_rules
users (1) ----< (N) api_keys
users (1) ----< (N) subscriptions

servers (1) ----< (N) server_metrics
servers (1) ----< (N) alert_history

alert_rules (1) ----< (N) alert_history
```

## Redis 使用说明

Redis 用于以下场景：

1. **Session 存储** - 用户登录会话缓存
2. **API 限流** - 请求频率限制
3. **实时数据** - WebSocket 连接状态、实时指标
4. **缓存** - 频繁查询的数据缓存

### Key 命名规范

```
vstats:session:{session_id}     # 用户会话
vstats:user:{user_id}:cache     # 用户缓存数据
vstats:server:{server_id}:live  # 服务器实时状态
vstats:ratelimit:{ip}:{endpoint} # API 限流计数器
vstats:ws:connections           # WebSocket 连接数
```

## 运维命令

### Docker Compose

```bash
# 启动服务
docker compose up -d

# 停止服务
docker compose down

# 重启服务
docker compose restart

# 查看日志
docker compose logs -f postgres
docker compose logs -f redis

# 进入容器
docker exec -it vstats-postgres bash
docker exec -it vstats-redis sh
```

### 数据库管理

```bash
# 连接数据库
docker exec -it vstats-postgres psql -U vstats -d vstats_cloud

# 备份数据库
docker exec vstats-postgres pg_dump -U vstats vstats_cloud > backup.sql

# 恢复数据库
docker exec -i vstats-postgres psql -U vstats vstats_cloud < backup.sql

# 清理过期会话
docker exec -it vstats-postgres psql -U vstats -d vstats_cloud -c "SELECT cleanup_expired_sessions();"

# 清理旧指标数据 (保留30天)
docker exec -it vstats-postgres psql -U vstats -d vstats_cloud -c "SELECT cleanup_old_metrics(30);"
```

### Redis 管理

```bash
# 连接 Redis CLI
docker exec -it vstats-redis redis-cli -a your_redis_password

# 查看内存使用
docker exec -it vstats-redis redis-cli -a your_redis_password INFO memory

# 清空缓存 (谨慎!)
docker exec -it vstats-redis redis-cli -a your_redis_password FLUSHDB
```

## 生产环境建议

### 安全性

1. **修改默认密码** - 务必修改所有默认密码
2. **限制端口访问** - 建议不对外暴露 5432 和 6379 端口
3. **启用 SSL** - 确保所有流量走 HTTPS
4. **定期备份** - 设置自动备份策略

### 性能优化

1. **PostgreSQL 调优**
   ```bash
   # 在 docker-compose.yml 中添加
   command: >
     postgres
     -c shared_buffers=256MB
     -c effective_cache_size=768MB
     -c maintenance_work_mem=64MB
     -c checkpoint_completion_target=0.9
   ```

2. **Redis 调优**
   - 已配置 maxmemory 和 LRU 策略
   - 根据实际使用调整内存限制

3. **指标数据分区**
   - 对于大量数据，考虑按时间分区 `server_metrics` 表

### 监控

建议使用以下工具监控服务状态：
- PostgreSQL: pgAdmin 或 Grafana + Prometheus
- Redis: Redis Commander 或 RedisInsight

## 故障排除

### PostgreSQL 无法启动

```bash
# 检查日志
docker compose logs postgres

# 检查数据目录权限
ls -la ./data/postgres
```

### Redis 连接失败

```bash
# 检查密码是否正确
docker exec -it vstats-redis redis-cli -a wrong_password ping

# 检查服务状态
docker compose ps redis
```

### 数据库初始化失败

```bash
# 重新初始化 (会清空数据!)
docker compose down -v
docker compose up -d
```

## 联系支持

如有问题，请通过以下方式联系：
- GitHub Issues: https://github.com/zsai001/vstats/issues
- 文档站点: https://vstats.zsoft.cc
