#!/bin/bash
# OpenTelemetry Collector 配置更新脚本

echo "=========================================="
echo "更新 OpenTelemetry Collector 配置"
echo "=========================================="
echo ""

# 容器名称
CONTAINER_NAME="otel-collector"

# 检查容器是否存在
if ! docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "❌ 错误：找不到容器 ${CONTAINER_NAME}"
    exit 1
fi

echo "✅ 找到容器：${CONTAINER_NAME}"
echo ""

# 备份当前配置
echo "1. 备份当前配置..."
docker exec ${CONTAINER_NAME} cat /etc/otel-collector-config.yaml > otel-collector-config.backup.yaml 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ 配置已备份到 otel-collector-config.backup.yaml"
else
    echo "⚠️  无法备份配置（可能路径不同）"
fi
echo ""

# 复制新配置到容器
echo "2. 复制新配置到容器..."
docker cp ./deploy/otel-collector-config.yaml ${CONTAINER_NAME}:/etc/otel-collector-config.yaml
if [ $? -eq 0 ]; then
    echo "✅ 配置文件已复制"
else
    echo "❌ 复制失败，请检查配置文件路径"
    exit 1
fi
echo ""

# 重启容器
echo "3. 重启容器..."
docker restart ${CONTAINER_NAME}
if [ $? -eq 0 ]; then
    echo "✅ 容器已重启"
else
    echo "❌ 重启失败"
    exit 1
fi
echo ""

# 等待容器启动
echo "4. 等待容器启动..."
sleep 5

# 检查容器状态
if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "✅ 容器运行正常"
else
    echo "❌ 容器未运行"
    exit 1
fi
echo ""

# 查看日志
echo "5. 查看最新日志..."
echo "----------------------------------------"
docker logs --tail 20 ${CONTAINER_NAME}
echo "----------------------------------------"
echo ""

echo "=========================================="
echo "✅ 配置更新完成！"
echo "=========================================="
echo ""
echo "请测试追踪功能："
echo "1. 调用接口产生追踪数据"
echo "2. 访问 Jaeger UI: http://14.103.153.242:16686/search"
echo "3. 在 Service 列表中应该能看到 'ShunFengParcel'"
echo ""

