# Jaeger 连接诊断脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Jaeger 连接诊断工具" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$JaegerHost = "14.103.153.242"

# 1. 测试端口连通性
Write-Host "1. 测试端口连通性..." -ForegroundColor Yellow
Write-Host ""

$ports = @(
    @{Port=16686; Name="Jaeger UI"},
    @{Port=4318; Name="OTLP HTTP"},
    @{Port=4317; Name="OTLP gRPC"},
    @{Port=14268; Name="Jaeger HTTP (旧)"}
)

foreach ($p in $ports) {
    Write-Host "   测试 $($p.Name) (端口 $($p.Port))..." -NoNewline
    $test = Test-NetConnection -ComputerName $JaegerHost -Port $p.Port -WarningAction SilentlyContinue -InformationLevel Quiet
    if ($test) {
        Write-Host " ✅" -ForegroundColor Green
    } else {
        Write-Host " ❌" -ForegroundColor Red
    }
}

# 2. 测试 OTLP 端点
Write-Host ""
Write-Host "2. 测试 OTLP HTTP 端点..." -ForegroundColor Yellow

try {
    $response = Invoke-WebRequest -Uri "http://$JaegerHost:4318/v1/traces" `
        -Method POST `
        -ContentType "application/json" `
        -Body '{"resourceSpans":[]}' `
        -UseBasicParsing `
        -TimeoutSec 5
    
    Write-Host "   ✅ OTLP 端点响应正常 (HTTP $($response.StatusCode))" -ForegroundColor Green
} catch {
    Write-Host "   ❌ OTLP 端点测试失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 3. 运行测试程序
Write-Host ""
Write-Host "3. 运行 OTLP 测试程序..." -ForegroundColor Yellow
Write-Host ""

go run test-otlp.go

# 4. 检查 Jaeger UI
Write-Host ""
Write-Host "4. 检查 Jaeger UI..." -ForegroundColor Yellow
Write-Host "   正在打开浏览器..."

Start-Process "http://$JaegerHost:16686/search?service=TestService&lookback=15m"

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "诊断完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "请在 Jaeger UI 中检查:" -ForegroundColor Yellow
Write-Host "  1. 是否能看到 'TestService' 服务" -ForegroundColor White
Write-Host "  2. 是否能看到 'test-span' 操作" -ForegroundColor White
Write-Host ""
Write-Host "如果能看到 TestService:" -ForegroundColor Yellow
Write-Host "  → 说明 OTLP 连接正常" -ForegroundColor Green
Write-Host "  → 问题可能在主服务配置" -ForegroundColor Green
Write-Host ""
Write-Host "如果看不到 TestService:" -ForegroundColor Yellow
Write-Host "  → 说明 Jaeger 服务器配置有问题" -ForegroundColor Red
Write-Host "  → 检查 Jaeger 是否启用了 OTLP receiver" -ForegroundColor Red
Write-Host ""
