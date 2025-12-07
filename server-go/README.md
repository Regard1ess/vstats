# vStats (Go Implementation)

这是 vStats 的 Go 语言实现，包含服务器端和代理端的代码。

## 目录结构

```
server-go/
├── cmd/
│   ├── server/     # 服务器端代码
│   └── agent/      # 代理端代码
├── internal/
│   └── common/     # 共享代码
├── go.mod
└── go.sum
```

## 构建

### 构建服务器

```bash
cd server-go
go mod tidy
go build -o vstats-server ./cmd/server
```

### 构建代理

```bash
cd server-go
go mod tidy
go build -o vstats-agent ./cmd/agent
```

## 运行

### 运行服务器

```bash
./vstats-server
```

### 运行代理

```bash
./vstats-agent
```

## 服务器命令行选项

- `--check`: 显示诊断信息
- `--reset-password`: 重置管理员密码

## 代理命令行选项

- `version`, `--version`, `-v`: 显示版本信息
- `register`: 注册代理到服务器
- `install`: 安装代理服务
- `uninstall`: 卸载代理服务
- `show-config`: 显示配置信息

## 环境变量

### 服务器

- `VSTATS_PORT`: 服务器端口（默认: 3001）

## API 端点

- `GET /health` - 健康检查
- `GET /api/metrics` - 获取本地服务器指标
- `GET /api/metrics/all` - 获取所有服务器指标
- `GET /api/history/:server_id?range=1h|24h|7d|30d` - 获取历史数据
- `POST /api/auth/login` - 登录
- `GET /api/auth/verify` - 验证令牌
- `GET /ws` - Dashboard WebSocket
- `GET /ws/agent` - Agent WebSocket

## 配置文件

配置文件位置：与可执行文件同目录下的 `vstats-config.json`

## 数据库

SQLite 数据库位置：与可执行文件同目录下的 `vstats.db`
