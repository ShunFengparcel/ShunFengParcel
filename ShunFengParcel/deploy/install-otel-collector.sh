#!/bin/bash

# OpenTelemetry Collector 安装脚本
# 适用于 Linux 系统

set -e

echo "=========================================="
echo "OpenTelemetry Collector 安装脚本"
echo "=========================================="

# 检测系统架构
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "不支持的架构: $ARCH"
        exit 1
        ;;
esac

# 版本号
VERSION="0.111.0"
DOWNLOAD_URL="https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v${VERSION}/otelcol_${VERSION}_linux_${ARCH}.tar.gz"

echo "系统架构: $ARCH"
echo "下载版本: $VERSION"

# 创建目录
echo "创建目录..."
sudo mkdir -p /opt/otel-collector
sudo mkdir -p /etc/otel-collector
sudo mkdir -p /var/log/otel-collector

# 下载 OpenTelemetry Collector
echo "下载 OpenTelemetry Collector..."
cd /tmp
wget -O otelcol.tar.gz "$DOWNLOAD_URL"

# 解压
echo "解压文件..."
tar -xzf otelcol.tar.gz
sudo mv otelcol /opt/otel-collector/

# 设置权限
sudo chmod +x /opt/otel-collector/otelcol

# 复制配置文件
echo "复制配置文件..."
if [ -f "otel-collector-config.yaml" ]; then
    sudo cp otel-collector-config.yaml /etc/otel-collector/config.yaml
else
    echo "警告: 配置文件不存在，请手动创建 /etc/otel-collector/config.yaml"
fi

# 创建 systemd 服务
echo "创建 systemd 服务..."
sudo tee /etc/systemd/system/otel-collector.service > /dev/null <<EOF
[Unit]
Description=OpenTelemetry Collector
After=network.target

[Service]
Type=simple
User=root
ExecStart=/opt/otel-collector/otelcol --config=/etc/otel-collector/config.yaml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

# 重载 systemd
echo "重载 systemd..."
sudo systemctl daemon-reload

# 启动服务
echo "启动 OpenTelemetry Collector..."
sudo systemctl start otel-collector
sudo systemctl enable otel-collector

# 检查状态
echo ""
echo "=========================================="
echo "安装完成！"
echo "=========================================="
echo ""
echo "服务状态:"
sudo systemctl status otel-collector --no-pager

echo ""
echo "端口监听:"
sudo netstat -tlnp | grep otelcol || sudo ss -tlnp | grep otelcol

echo ""
echo "常用命令:"
echo "  查看状态: sudo systemctl status otel-collector"
echo "  查看日志: sudo journalctl -u otel-collector -f"
echo "  重启服务: sudo systemctl restart otel-collector"
echo "  停止服务: sudo systemctl stop otel-collector"
echo ""
