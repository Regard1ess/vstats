#!/bin/bash
# ===========================================
# VStats Cloud - Quick Check Script
# 检查 API 路由是否正确配置
# ===========================================

set -e

DOMAIN="${1:-vstats.zsoft.cc}"

echo "========================================="
echo "VStats Cloud - API Route Check"
echo "Domain: $DOMAIN"
echo "========================================="
echo ""

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

check_endpoint() {
    local name="$1"
    local url="$2"
    local expect_type="$3"  # "json", "text", "binary"
    
    echo -n "Checking $name... "
    
    response=$(curl -s -w "\n%{http_code}" "$url" 2>/dev/null)
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "200" ]; then
        case "$expect_type" in
            "json")
                if echo "$body" | jq . >/dev/null 2>&1; then
                    echo -e "${GREEN}✓ OK${NC}"
                    return 0
                else
                    echo -e "${RED}✗ FAIL (not JSON)${NC}"
                    return 1
                fi
                ;;
            "text")
                if [ -n "$body" ] && [ "$body" != "404 page not found" ]; then
                    echo -e "${GREEN}✓ OK${NC} (version: $body)"
                    return 0
                else
                    echo -e "${RED}✗ FAIL (empty or 404)${NC}"
                    return 1
                fi
                ;;
            "binary")
                content_type=$(curl -sI "$url" | grep -i "content-type" | cut -d' ' -f2-)
                if [[ "$content_type" == *"octet-stream"* ]] || [[ "$content_type" == *"application/x-"* ]]; then
                    echo -e "${GREEN}✓ OK${NC} (content-type: $content_type)"
                    return 0
                else
                    echo -e "${RED}✗ FAIL${NC} (content-type: $content_type)"
                    echo "  Expected: application/octet-stream"
                    echo "  Got: $content_type"
                    return 1
                fi
                ;;
        esac
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP $http_code)"
        if [ "$http_code" = "404" ]; then
            echo "  Error: Route not found"
            echo "  This usually means:"
            echo "    1. Nginx config missing the route"
            echo "    2. API container not running"
            echo "    3. Need to restart nginx: docker compose restart nginx"
        fi
        return 1
    fi
}

# 基础健康检查
echo "1. Basic Health Checks"
check_endpoint "Health" "https://$DOMAIN/health" "text"
check_endpoint "Version" "https://$DOMAIN/version" "json"
echo ""

# Release API 检查
echo "2. Release API Endpoints"
check_endpoint "Latest Version (JSON)" "https://$DOMAIN/api/release/latest" "json"
check_endpoint "Latest Version (Text)" "https://$DOMAIN/api/release/version" "text"
check_endpoint "Assets List" "https://$DOMAIN/api/release/assets" "json"
echo ""

# Download 端点检查
echo "3. Download Endpoints"
check_endpoint "Download (Server Binary)" "https://$DOMAIN/download/latest/vstats-server-linux-amd64" "binary"
check_endpoint "Download (Agent Binary)" "https://$DOMAIN/download/latest/vstats-agent-linux-amd64" "binary"
echo ""

# Docker 容器状态
echo "4. Docker Container Status"
if command -v docker &> /dev/null; then
    echo -n "Checking API container... "
    if docker ps | grep -q vstats-api; then
        container_id=$(docker ps | grep vstats-api | awk '{print $1}')
        uptime=$(docker ps | grep vstats-api | awk '{print $10, $11}')
        echo -e "${GREEN}✓ Running${NC} (uptime: $uptime)"
        
        echo -n "Checking API container health... "
        health=$(docker inspect --format='{{.State.Health.Status}}' $container_id 2>/dev/null || echo "unknown")
        if [ "$health" = "healthy" ]; then
            echo -e "${GREEN}✓ Healthy${NC}"
        else
            echo -e "${YELLOW}⚠ $health${NC}"
        fi
    else
        echo -e "${RED}✗ Not running${NC}"
    fi
    
    echo -n "Checking Nginx container... "
    if docker ps | grep -q vstats-nginx; then
        uptime=$(docker ps | grep vstats-nginx | awk '{print $10, $11}')
        echo -e "${GREEN}✓ Running${NC} (uptime: $uptime)"
    else
        echo -e "${RED}✗ Not running${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Docker not available (run on deployment server)${NC}"
fi
echo ""

echo "========================================="
echo "Check completed"
echo "========================================="
