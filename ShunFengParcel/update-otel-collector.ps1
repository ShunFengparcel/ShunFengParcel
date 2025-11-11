# OpenTelemetry Collector 配置更新脚本 (Windows PowerShell)

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "更新 OpenTelemetry Collector 配置" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

$CONTAINER_NAME = "otel-collector"

# 检查 Docker 是否可用
try {
    docker ps | Out-Null
} catch {
    Write-Host "❌ 错误：Docker 不可用或未安装" -ForegroundColor Red
    exit 1
}

# 检查容器是否存在
$containerExists = docker ps -a --format "{{.Names}}" | Select-String -Pattern "^${CONTAINER_NAME}$"
if (-not $containerExists) {
    Write-Host "❌ 错误：找不到容器 ${CONTAINER_NAME}" -ForegroundColor Red
    Write-Host "请在服务器 14.103.153.242 上运行此脚本" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ 找到容器：${CONTAINER_NAME}" -ForegroundColor Green
Write-Host ""

# 备份当前配置
Write-Host "1. 备份当前配置..." -ForegroundColor Yellow
try {
    docker exec $CONTAINER_NAME cat /etc/otel-collector-config.yaml > otel-collector-config.backup.yaml 2>$null
    Write-Host "✅ 配置已备份到 otel-collector-config.backup.yaml" -ForegroundColor Green
} catch {
    Write-Host "⚠️  无法备份配置（可能路径不同）" -ForegroundColor Yellow
}
Write-Host ""

# 复制新配置到容器
Write-Host "2. 复制新配置到容器..." -ForegroundColor Yellow
$result = docker cp .\deploy\otel-collector-config.yaml ${CONTAINER_NAME}:/etc/otel-collector-config.yaml 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 配置文件已复制" -ForegroundColor Green
} else {
    Write-Host "❌ 复制失败：$result" -ForegroundColor Red
    Write-Host "请检查配置文件路径" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# 重启容器
Write-Host "3. 重启容器..." -ForegroundColor Yellow
docker restart $CONTAINER_NAME | Out-Null
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 容器已重启" -ForegroundColor Green
} else {
    Write-Host "❌ 重启失败" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 等待容器启动
Write-Host "4. 等待容器启动..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# 检查容器状态
$running = docker ps --format "{{.Names}}" | Select-String -Pattern "^${CONTAINER_NAME}$"
if ($running) {
    Write-Host "✅ 容器运行正常" -ForegroundColor Green
} else {
    Write-Host "❌ 容器未运行" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 查看日志
Write-Host "5. 查看最新日志..." -ForegroundColor Yellow
Write-Host "----------------------------------------" -ForegroundColor Gray
docker logs --tail 20 $CONTAINER_NAME
Write-Host "----------------------------------------" -ForegroundColor Gray
Write-Host ""

Write-Host "==========================================" -ForegroundColor Green
Write-Host "✅ 配置更新完成！" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""
Write-Host "请测试追踪功能：" -ForegroundColor Cyan
Write-Host "1. 调用接口产生追踪数据" -ForegroundColor White
Write-Host "2. 访问 Jaeger UI: http://14.103.153.242:16686/search" -ForegroundColor White
Write-Host "3. 在 Service 列表中应该能看到 'ShunFengParcel'" -ForegroundColor White
Write-Host ""

