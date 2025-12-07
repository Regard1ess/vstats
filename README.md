# vStats - Server Monitoring Dashboard

[![GitHub Release](https://img.shields.io/github/v/release/zsai001/vstats?style=flat-square)](https://github.com/zsai001/vstats/releases)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=flat-square&logo=react&logoColor=61DAFB)
![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=flat-square&logo=typescript&logoColor=white)

极简美观的服务器探针监控系统。Go 驱动，毫秒级延迟，一键部署。

## 💝 赞助商

<div align="center">

感谢以下赞助商对本项目的支持！

[TOHU Cloud](https://www.tohu.cloud) | [Debee](https://debee.io/)

</div>

## 📸 预览

<table>
  <tr>
    <td align="center">
      <img src="https://vstats.zsoft.cc/theme/1.png" alt="预览图 1" width="100%"/>
    </td>
    <td align="center">
      <img src="https://vstats.zsoft.cc/theme/2.png" alt="预览图 2" width="100%"/>
    </td>
  </tr>
  <tr>
    <td align="center">
      <img src="https://vstats.zsoft.cc/theme/3.png" alt="预览图 3" width="100%"/>
    </td>
    <td align="center">
      <img src="https://vstats.zsoft.cc/theme/4.png" alt="预览图 4" width="100%"/>
    </td>
  </tr>
</table>

## ✨ 特性

- 🚀 **实时监控** - WebSocket 实时推送系统指标
- 🖥️ **多服务器管理** - 支持监控多台服务器
- 💻 **CPU / 内存 / 磁盘 / 网络** - 全方位监控
- 🎨 **现代 UI** - 玻璃拟态设计，流畅动画
- 🔐 **安全认证** - JWT 认证保护管理接口
- ⚡ **一键部署** - Docker / 脚本一键安装

## 📚 文档与资源

| 资源 | 链接 |
|------|------|
| 📖 **完整文档** | [vstats.zsoft.cc](https://vstats.zsoft.cc) |
| 🎯 **在线演示** | [vps.zsoft.cc](https://vps.zsoft.cc/) |
| 🐳 **Docker Hub** | [zsai001/vstats-server](https://hub.docker.com/r/zsai001/vstats-server) |
| 📦 **GitHub Releases** | [下载页面](https://github.com/zsai001/vstats/releases) |

## 🚀 快速开始

```bash
# Docker 一键部署
docker run -d --name vstats-server -p 3001:3001 \
  -v $(pwd)/data:/app/data zsai001/vstats-server:latest
```

更多安装方式请访问 **[文档站点](https://vstats.zsoft.cc/docs)**

## ⭐ Star History

<a href="https://star-history.com/#zsai001/vstats&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=zsai001/vstats&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=zsai001/vstats&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=zsai001/vstats&type=Date" />
 </picture>
</a>

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
