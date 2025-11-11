# 检查 Jaeger 服务器配置

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "检查 Jaeger 服务器配置" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$JaegerHost = "14.103.153.242"

# 1. 检查 Jaeger 版本和配置
Write-Host "1. 检查 Jaeger 服务信息..." -ForegroundColor Yellow

try {
    $response = Invoke-WebRequest -Uri "http://${JaegerHost}:16686/api/services" -UseBasicParsing
    $services = $response.Content | ConvertFrom-Json
    
    Write-Host "   当前服务列表:" -ForegroundColor Green
    foreach ($service in $services.data) {
        Write-Host "   - $service" -ForegroundColor White
    }
} catch {
    Write-Host "   ❌ 无法获取服务列表: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# 2. 测试不同的 OTLP 端点
Write-Host "2. 测试 OTLP 端点..." -ForegroundColor Yellow

$endpoints = @(
    @{Url="http://${JaegerHost}:4318/v1/traces"; Name="OTLP HTTP /v1/traces"},
    @{Url="http://${JaegerHost}:4318/api/traces"; Name="OTLP HTTP /api/traces"},
    @{Url="http://${JaegerHost}:14268/api/traces"; Name="Jaeger HTTP (旧)"}
)

foreach ($ep in $endpoints) {
    Write-Host "   测试 $($ep.Name)..." -NoNewline
    try {
        $response = Invoke-WebRequest -Uri $ep.Url `
            -Method POST `
            -ContentType "application/json" `
            -Body '{"resourceSpans":[]}' `
            -UseBasicParsing `
            -TimeoutSec 3
        Write-Host " ✅ (HTTP $($response.StatusCode))" -ForegroundColor Green
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode) {
            Write-Host " ⚠️  (HTTP $statusCode)" -ForegroundColor Yellow
        } else {
            Write-Host " ❌ (连接失败)" -ForegroundColor Red
        }
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "诊断结果" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "可能的问题:" -ForegroundColor Yellow
Write-Host ""
Write-Host "1. Jaeger 使用的是 all-in-one 模式" -ForegroundColor White
Write-Host "   → 可能没有启用 OTLP receiver" -ForegroundColor Gray
Write-Host ""
Write-Host "2. Jaeger 版本太旧" -ForegroundColor White
Write-Host "   → 需要 v1.35+ 才支持 OTLP" -ForegroundColor Gray
Write-Host ""
Write-Host "3. OTLP receiver 没有正确配置" -ForegroundColor White
Write-Host "   → 需要在启动参数中启用" -ForegroundColor Gray
Write-Host ""

Write-Host "建议的解决方案:" -ForegroundColor Yellow
Write-Host ""
Write-Host "方案 1: 使用 Jaeger 原生协议（已废弃但可用）" -ForegroundColor White
Write-Host "  端口: 14268" -ForegroundColor Gray
Write-Host ""
Write-Host "方案 2: 升级 Jaeger 并启用 OTLP" -ForegroundColor White
Write-Host "  docker run -d --name jaeger \\" -ForegroundColor Gray
Write-Host "    -e COLLECTOR_OTLP_ENABLED=true \\" -ForegroundColor Gray
Write-Host "    -p 16686:16686 \\" -ForegroundColor Gray
Write-Host "    -p 4318:4318 \\" -ForegroundColor Gray
Write-Host "    jaegertracing/all-in-one:latest" -ForegroundColor Gray
Write-Host ""
Write-Host "方案 3: 使用 OpenTelemetry Collector 作为中间层" -ForegroundColor White
Write-Host "  应用 → OTLP → Collector → Jaeger" -ForegroundColor Gray
Write-Host ""
