# vStats OAuth Proxy

Cloudflare Worker 实现的 OAuth 2.0 代理服务，为自部署的 vStats 实例提供统一的 OAuth 认证。

## 功能

- 支持 GitHub OAuth
- 支持 Google OAuth
- 免费（Cloudflare Workers 免费额度：100,000 请求/天）
- 全球边缘部署，低延迟

## 部署步骤

### 1. 安装依赖

```bash
cd oauth-proxy
npm install
```

### 2. 配置 Cloudflare

首先确保你有 Cloudflare 账号，然后登录：

```bash
npx wrangler login
```

### 3. 创建 OAuth Apps

#### GitHub OAuth App

1. 访问 https://github.com/settings/developers
2. 点击 "New OAuth App"
3. 填写信息：
   - Application name: `vStats OAuth Proxy`
   - Homepage URL: `https://vstats.zsoft.cc`
   - Authorization callback URL: `https://vstats-oauth-proxy.zsai001.workers.dev/oauth/github/callback`
4. 保存 Client ID 和 Client Secret

#### Google OAuth Client

1. 访问 https://console.cloud.google.com/apis/credentials
2. 创建 "OAuth 2.0 Client ID"（Web 应用）
3. 添加授权重定向 URI: `https://vstats-oauth-proxy.zsai001.workers.dev/oauth/google/callback`
4. 保存 Client ID 和 Client Secret

### 4. 配置 Secrets

```bash
# 设置 GitHub OAuth
npx wrangler secret put GITHUB_CLIENT_ID
npx wrangler secret put GITHUB_CLIENT_SECRET

# 设置 Google OAuth
npx wrangler secret put GOOGLE_CLIENT_ID
npx wrangler secret put GOOGLE_CLIENT_SECRET
```

### 5. 部署

```bash
npm run deploy
```

### 6. 配置自定义域名（可选）

在 Cloudflare Dashboard 中：
1. 进入 Workers & Pages
2. 选择 `vstats-oauth-proxy`
3. 点击 "Custom Domains"
4. 添加 `auth.vstats.zsoft.cc`

## API 端点

| 端点 | 描述 |
|------|------|
| `GET /oauth/github?redirect_uri=...&state=...` | 开始 GitHub OAuth 流程 |
| `GET /oauth/google?redirect_uri=...&state=...` | 开始 Google OAuth 流程 |
| `GET /health` | 健康检查 |

## 工作流程

```
用户的 vStats 实例                    OAuth Proxy                         GitHub/Google
      │                                   │                                    │
      │ 1. /oauth/github                  │                                    │
      │    ?redirect_uri=xxx&state=yyy    │                                    │
      ├──────────────────────────────────►│                                    │
      │                                   │                                    │
      │                                   │ 2. 重定向到 GitHub/Google          │
      │                                   ├───────────────────────────────────►│
      │                                   │                                    │
      │                                   │ 3. 用户授权                        │
      │                                   │◄───────────────────────────────────┤
      │                                   │    code                            │
      │                                   │                                    │
      │                                   │ 4. 交换 token，获取用户信息        │
      │                                   │◄──────────────────────────────────►│
      │                                   │                                    │
      │ 5. 回调返回用户信息               │                                    │
      │◄──────────────────────────────────┤                                    │
      │    ?state=yyy&provider=github     │                                    │
      │    &user=username                 │                                    │
```

## 本地开发

```bash
# 启动本地开发服务器
npm run dev
```

访问 http://localhost:8787 测试。

## 安全说明

- OAuth Client Secret 通过 Cloudflare Secrets 安全存储
- 使用 state 参数防止 CSRF 攻击
- redirect_uri 验证确保只重定向到有效的 URL

