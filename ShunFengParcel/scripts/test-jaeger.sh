#!/bin/bash

# Jaeger 集成测试脚本

echo "=========================================="
echo "Jaeger 链路追踪集成测试"
echo "=========================================="
echo ""

JAEGER_HOST="14.103.153.242"
SERVICE_HOST="localhost:8000"

# 1. 检查 Jaeger 服务器连通性
echo "1. 检查 Jaeger 服务器连通性..."
echo "   - UI 端口 (16686):"
if nc -zv $JAEGER_HOST 16686 2>&1 | grep -q succeeded; then
    echo "     ✅ 可访问"
else
    echo "     ❌ 不可访问"
fi

echo "   - OTLP HTTP 端口 (4318):"
if nc -zv $JAEGER_HOST 4318 2>&1 | grep -q succeeded; then
    echo "     ✅ 可访问"
else
    echo "     ❌ 不可访问"
fi

echo ""

# 2. 发送测试请求
echo "2. 发送测试请求到服务..."
echo "   POST http://$SERVICE_HOST/admin/register"

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
  "http://$SERVICE_HOST/admin/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "password": "test_pass"
  }')

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ 请求成功 (HTTP $HTTP_CODE)"
else
    echo "   ⚠️  请求返回 HTTP $HTTP_CODE"
fi

echo ""

# 3. 等待数据上报
echo "3. 等待链路数据上报到 Jaeger (10秒)..."
sleep 10

# 4. 检查 Jaeger UI
echo ""
echo "4. 查看链路追踪数据:"
echo "   打开浏览器访问: http://$JAEGER_HOST:16686"
echo "   - 选择服务: ShunFengParcel"
echo "   - 点击 'Find Traces'"
echo "   - 查看最新的 trace"

echo ""
echo "=========================================="
echo "测试完成！"
echo "=========================================="
