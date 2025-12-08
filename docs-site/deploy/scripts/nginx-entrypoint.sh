#!/bin/sh
# ===========================================
# Nginx Entrypoint Script
# Checks for SSL certificates and generates if missing
# ===========================================

set -e

SSL_DIR="/etc/nginx/ssl"
CERT_FILE="$SSL_DIR/cert.pem"
KEY_FILE="$SSL_DIR/key.pem"

# 创建 SSL 目录（如果不存在）
mkdir -p "$SSL_DIR"

# 检查证书是否存在
if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    echo "======================================"
    echo "⚠️  SSL certificates not found!"
    echo "======================================"
    echo "Generating self-signed certificate..."
    echo ""
    
    # 检查 openssl 是否可用
    if ! command -v openssl >/dev/null 2>&1; then
        echo "❌ Error: openssl not found in nginx:alpine image"
        echo "   Please install openssl or provide SSL certificates manually"
        exit 1
    fi
    
    # 生成自签名证书
    # 尝试使用 -addext（OpenSSL 1.1.1+），如果不支持则使用旧方法
    if openssl req -help 2>&1 | grep -q "\-addext"; then
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
            -keyout "$KEY_FILE" \
            -out "$CERT_FILE" \
            -subj "/CN=localhost/O=VStats/C=US" \
            -addext "subjectAltName=DNS:localhost,DNS:*.localhost,IP:127.0.0.1,IP:::1" 2>/dev/null || \
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
            -keyout "$KEY_FILE" \
            -out "$CERT_FILE" \
            -subj "/CN=localhost/O=VStats/C=US"
    else
        # 旧版本 OpenSSL，使用配置文件
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
            -keyout "$KEY_FILE" \
            -out "$CERT_FILE" \
            -subj "/CN=localhost/O=VStats/C=US"
    fi
    
    # 设置权限
    chmod 644 "$CERT_FILE" 2>/dev/null || true
    chmod 600 "$KEY_FILE" 2>/dev/null || true
    
    echo "✅ Self-signed certificate generated"
    echo "   Certificate: $CERT_FILE"
    echo "   Private Key: $KEY_FILE"
    echo ""
    echo "⚠️  This is a self-signed certificate for development only!"
    echo "   For production, use Let's Encrypt or Cloudflare Origin Certificate."
    echo ""
else
    echo "✅ SSL certificates found"
fi

# 执行原始的 nginx entrypoint
exec /docker-entrypoint.sh "$@"
