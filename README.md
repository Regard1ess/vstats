# vStats - Server Monitoring Dashboard

[![GitHub Release](https://img.shields.io/github/v/release/Regard1ess/vstats?style=flat-square)](https://github.com/Regard1ess/vstats/releases)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)
![Rust](https://img.shields.io/badge/Rust-000000?style=flat-square&logo=rust&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=flat-square&logo=react&logoColor=61DAFB)
![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=flat-square&logo=typescript&logoColor=white)

极简美观的服务器探针监控系统。Rust 驱动，毫秒级延迟，一键部署。

**文档网站**: [vstats.zsoft.cc](https://vstats.zsoft.cc)

**在线示例**: [vps.zsoft.cc](https://vps.zsoft.cc/)

## 📸 预览

<table>
  <tr>
    <td align="center">
      <img src="doc/1.png" alt="预览图 1" width="100%"/>
    </td>
    <td align="center">
      <img src="doc/2.png" alt="预览图 2" width="100%"/>
    </td>
  </tr>
  <tr>
    <td align="center">
      <img src="doc/3.png" alt="预览图 3" width="100%"/>
    </td>
    <td align="center">
      <img src="doc/4.png" alt="预览图 4" width="100%"/>
    </td>
  </tr>
</table>

## ✨ 特性

- 🚀 **实时监控** - WebSocket 实时推送系统指标
- 🖥️ **多服务器管理** - 支持监控多台服务器
- 💻 **CPU 监控** - 总体使用率和每核心负载可视化
- 🧠 **内存监控** - RAM 和 Swap 使用情况
- 💾 **磁盘监控** - 挂载点和使用率
- 🌐 **网络监控** - 实时上传/下载速度
- 📊 **负载平均** - 1/5/15 分钟负载
- 🎨 **现代 UI** - 玻璃拟态设计，流畅动画
- 🔐 **安全认证** - JWT 认证保护管理接口
- ⚡ **一键部署** - 提供自动化安装脚本

## 🚀 一键安装

### 安装主控端 (Dashboard)

```bash
curl -fsSL https://github.com/Regard1ess/vstats/raw/refs/heads/main/install.sh | sudo bash
```

### 安装探针 (Agent)

登录 Dashboard 后，进入 **Settings** 页面获取安装命令，或直接运行：

#### Linux / macOS

```bash
curl -fsSL https://github.com/Regard1ess/vstats/raw/refs/heads/main/agent.sh | sudo bash -s -- \
  --server http://YOUR_DASHBOARD_IP:30001 \
  --token "your-jwt-token" \
  --name "$(hostname)" \
  --location "US" \
  --provider "Vultr"
```

#### Windows (PowerShell 管理员模式)

```powershell
# 下载安装脚本
Invoke-WebRequest -Uri "https://vstats.zsoft.cc/agent.ps1" -OutFile "agent.ps1"

# 运行安装
.\agent.ps1 -Server "http://YOUR_DASHBOARD_IP:30001" -Token "your-jwt-token" -Name $env:COMPUTERNAME
```

或一行命令安装：

```powershell
irm https://vstats.zsoft.cc/agent.ps1 -OutFile agent.ps1; .\agent.ps1 -Server "http://YOUR_DASHBOARD_IP:30001" -Token "your-jwt-token"
```

### 升级

#### Linux / macOS

```bash
# 升级主控端
curl -fsSL https://vstats.zsoft.cc/install.sh | sudo bash -s -- --upgrade

# 升级探针
curl -fsSL https://vstats.zsoft.cc/agent.sh | sudo bash -s -- --upgrade
```

#### Windows (PowerShell 管理员模式)

```powershell
# 升级探针
.\agent.ps1 -Upgrade
```

### 卸载

#### Linux / macOS

```bash
# 卸载主控端
curl -fsSL https://github.com/Regard1ess/vstats/raw/refs/heads/main/install.sh | sudo bash -s -- --uninstall

# 卸载探针
curl -fsSL https://github.com/Regard1ess/vstats/raw/refs/heads/main/agent.sh | sudo bash -s -- --uninstall
```

#### Windows (PowerShell 管理员模式)

```powershell
# 卸载探针
.\agent.ps1 -Uninstall
```

## 🏗️ 架构

```
┌─────────────────────────────────────────────────────────────┐
│                         Dashboard                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   Web UI    │  │  REST API   │  │  WebSocket  │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│                          │                                  │
│              ┌───────────┴───────────┐                      │
│              │    Rust Backend       │                      │
│              │   (Axum + Tokio)      │                      │
│              └───────────────────────┘                      │
└─────────────────────────────────────────────────────────────┘
                           │
           ┌───────────────┼───────────────┐
           ▼               ▼               ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │   Agent 1   │ │   Agent 2   │ │   Agent N   │
    └─────────────┘ └─────────────┘ └─────────────┘
```

## 🛠️ 手动开发环境

### 启动后端服务

```bash
cd server
cargo run --release
```

服务器将在 `http://localhost:30001` 启动。

### 启动前端开发服务器

```bash
cd web
npm install
npm run dev
```

前端将在 `http://localhost:5173` 启动。

## 📁 项目结构

```
vstats/
├── server/                 # Rust 后端 (Dashboard)
│   ├── src/
│   │   ├── main.rs        # 主程序
│   │   ├── handlers.rs    # API 处理
│   │   ├── websocket.rs   # WebSocket 处理
│   │   └── ...
│   └── Cargo.toml
├── agent/                  # Rust 探针 (Agent)
│   ├── src/
│   │   ├── main.rs        # 主程序
│   │   ├── metrics.rs     # 指标采集
│   │   ├── websocket.rs   # WebSocket 客户端
│   │   └── config.rs      # 配置管理
│   └── Cargo.toml
├── web/                    # React 前端
│   ├── src/
│   │   ├── pages/         # 页面组件
│   │   ├── components/    # UI 组件
│   │   ├── hooks/         # 自定义 Hooks
│   │   └── context/       # React Context
│   └── package.json
├── scripts/                # 安装脚本
│   ├── install.sh         # 主控端安装脚本
│   ├── agent.sh           # 探针安装脚本 (Linux/macOS)
│   └── agent.ps1          # 探针安装脚本 (Windows)
├── docs/                   # GitHub Pages 文档站
│   ├── index.html         # 落地页
│   ├── install.sh         # 安装脚本 (镜像)
│   ├── agent.sh           # 探针脚本 (镜像)
│   └── agent.ps1          # Windows 探针脚本 (镜像)
└── .github/
    └── workflows/
        ├── release.yml    # 构建发布工作流
        └── pages.yml      # GitHub Pages 部署
```

## 🔌 API 接口

### 公开接口

| 端点 | 方法 | 描述 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/api/metrics` | GET | 获取当前系统指标 |
| `/api/servers` | GET | 获取服务器列表 |
| `/api/auth/login` | POST | 用户登录 |
| `/api/auth/verify` | GET | 验证 Token |
| `/ws` | WebSocket | 实时指标推送 (1秒/次) |

### 需要认证的接口

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/servers` | POST | 添加服务器 |
| `/api/servers/{id}` | DELETE | 删除服务器 |
| `/api/auth/password` | POST | 修改密码 |

## 📊 系统指标

```typescript
interface SystemMetrics {
  timestamp: string;
  hostname: string;
  os: { name, version, kernel, arch };
  cpu: { brand, cores, usage, frequency, per_core[] };
  memory: { total, used, available, swap_total, swap_used, usage_percent };
  disks: [{ name, mount_point, fs_type, total, used, available, usage_percent }];
  network: { interfaces[], total_rx, total_tx };
  uptime: number;
  load_average: { one, five, fifteen };
}
```

## 🔐 默认凭据

- **默认密码**: `admin` (或安装时生成的随机密码)
- 首次登录后请立即修改密码

## 🛠️ 服务管理

### Linux (systemd)

```bash
# 查看状态
systemctl status vstats

# 重启服务
systemctl restart vstats

# 查看日志
journalctl -u vstats -f

# 停止服务
systemctl stop vstats
```

### Windows (管理员模式)

```powershell
# 查看状态
sc query vstats-agent
# 或
Get-Service vstats-agent

# 重启服务
Restart-Service vstats-agent

# 停止服务
Stop-Service vstats-agent

# 启动服务
Start-Service vstats-agent

# 查看日志
Get-EventLog -LogName Application -Source vstats-agent -Newest 50
```

## 🔧 技术栈

### 后端 (Rust)
- **Axum** - 高性能 Web 框架
- **Tokio** - 异步运行时
- **sysinfo** - 系统信息采集
- **WebSocket** - 实时数据推送
- **JWT** - 认证授权
- **bcrypt** - 密码加密

### 前端 (React)
- **Vite** - 快速构建工具
- **TypeScript** - 类型安全
- **Tailwind CSS** - 原子化 CSS
- **React Router** - 路由管理
- **自定义组件** - 进度环、进度条等

## 📦 发布流程

1. 创建新的 Git tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions 自动:
   - 构建多平台二进制文件 (Linux x86_64/aarch64, macOS x86_64/aarch64, Windows x86_64/aarch64)
   - 构建 Web 资源
   - 创建 GitHub Release
   - 上传所有构建产物

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
