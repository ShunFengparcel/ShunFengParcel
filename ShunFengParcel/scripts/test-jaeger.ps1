# Jaeger 集成测试脚本 (PowerShell)

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Jaeger 链路追踪集成测试" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

$JaegerHost = "14.103.153.242"
$ServiceHost = "localhost:8000"

# 1. 检查 Jaeger 服务器连通性
Write-Host "1. 检查 Jaeger 服务器连通性..." -ForegroundColor Yellow

Write-Host "   - UI 端口 (16686):" -NoNewline
$uiTest = Test-NetConnection -ComputerName $JaegerHost -Port 16686 -WarningAction SilentlyContinue
if ($uiTest.TcpTestSucceeded) {
    Write-Host " ✅ 可访问" -ForegroundColor Green
} else {
    Write-Host " ❌ 不可访问" -ForegroundColor Red
}

Write-Host "   - OTLP HTTP 端口 (4318):" -NoNewline
$otlpTest = Test-NetConnection -ComputerName $JaegerHost -Port 4318 -WarningAction SilentlyContinue
if ($otlpTest.TcpTestSucceeded) {
    Write-Host " ✅ 可访问" -ForegroundColor Green
} else {
    Write-Host " ❌ 不可访问" -ForegroundColor Red
}

Write-Host ""

# 2. 发送测试请求
Write-Host "2. 发送测试请求到服务..." -ForegroundColor Yellow
Write-Host "   POST http://$ServiceHost/admin/register"

$body = @{
    username = "test_user"
    password = "test_pass"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "http://$ServiceHost/admin/register" `
        -Method POST `
        -ContentType "application/json" `
        -Body $body `
        -UseBasicParsing
    
    if ($response.StatusCode -eq 200) {
        Write-Host "   ✅ 请求成功 (HTTP $($response.StatusCode))" -ForegroundColor Green
    } else {
        Write-Host "   ⚠️  请求返回 HTTP $($response.StatusCode)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "   ❌ 请求失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# 3. 等待数据上报
Write-Host "3. 等待链路数据上报到 Jaeger..." -ForegroundColor Yellow
for ($i = 10; $i -gt 0; $i--) {
    Write-Host "   倒计时: $i 秒" -NoNewline
    Start-Sleep -Seconds 1
    Write-Host "`r" -NoNewline
}
Write-Host "   ✅ 等待完成                    "

# 4. 打开 Jaeger UI
Write-Host ""
Write-Host "4. 查看链路追踪数据:" -ForegroundColor Yellow
Write-Host "   正在打开 Jaeger UI..."

$jaegerUrl = "http://$JaegerHost:16686"
Start-Process $jaegerUrl

Write-Host ""
Write-Host "   在 Jaeger UI 中:" -ForegroundColor Cyan
Write-Host "   1. 选择服务: ShunFengParcel"
Write-Host "   2. 点击 'Find Traces'"
Write-Host "   3. 查看最新的 trace"

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "测试完成！" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
