# 启动服务并测试 Jaeger 连接

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "启动 ShunFengParcel 服务并测试 Jaeger" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 1. 停止旧进程
Write-Host "1. 检查并停止旧进程..." -ForegroundColor Yellow
$oldProcess = Get-Process -Name "ShunFengParcel" -ErrorAction SilentlyContinue
if ($oldProcess) {
    Stop-Process -Name "ShunFengParcel" -Force
    Write-Host "   ✅ 已停止旧进程" -ForegroundColor Green
    Start-Sleep -Seconds 2
} else {
    Write-Host "   ℹ️  没有运行中的进程" -ForegroundColor Gray
}

# 2. 编译服务
Write-Host ""
Write-Host "2. 编译服务..." -ForegroundColor Yellow
go build -o bin/ShunFengParcel.exe ./cmd/ShunFengParcel
if ($LASTEXITCODE -eq 0) {
    Write-Host "   ✅ 编译成功" -ForegroundColor Green
} else {
    Write-Host "   ❌ 编译失败" -ForegroundColor Red
    exit 1
}

# 3. 启动服务（后台运行）
Write-Host ""
Write-Host "3. 启动服务..." -ForegroundColor Yellow
$process = Start-Process -FilePath ".\bin\ShunFengParcel.exe" -WorkingDirectory "." -PassThru -WindowStyle Normal
Write-Host "   ✅ 服务已启动 (PID: $($process.Id))" -ForegroundColor Green

# 4. 等待服务启动
Write-Host ""
Write-Host "4. 等待服务初始化..." -ForegroundColor Yellow
Start-Sleep -Seconds 5
Write-Host "   ✅ 服务应该已经启动" -ForegroundColor Green

# 5. 检查端口
Write-Host ""
Write-Host "5. 检查服务端口..." -ForegroundColor Yellow
$httpPort = Get-NetTCPConnection -LocalPort 8000 -ErrorAction SilentlyContinue
$grpcPort = Get-NetTCPConnection -LocalPort 9000 -ErrorAction SilentlyContinue

if ($httpPort) {
    Write-Host "   ✅ HTTP 端口 8000 已监听" -ForegroundColor Green
} else {
    Write-Host "   ❌ HTTP 端口 8000 未监听" -ForegroundColor Red
}

if ($grpcPort) {
    Write-Host "   ✅ gRPC 端口 9000 已监听" -ForegroundColor Green
} else {
    Write-Host "   ❌ gRPC 端口 9000 未监听" -ForegroundColor Red
}

# 6. 发送测试请求
Write-Host ""
Write-Host "6. 发送测试请求..." -ForegroundColor Yellow
Start-Sleep -Seconds 2

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8000/admin/register" `
        -Method POST `
        -ContentType "application/json" `
        -Body '{"username":"jaeger_test","password":"test123"}' `
        -UseBasicParsing `
        -ErrorAction Stop
    
    Write-Host "   ✅ 请求成功 (HTTP $($response.StatusCode))" -ForegroundColor Green
} catch {
    Write-Host "   ⚠️  请求失败: $($_.Exception.Message)" -ForegroundColor Yellow
    Write-Host "   (这可能是正常的，如果接口需要特定参数)" -ForegroundColor Gray
}

# 7. 等待数据上报
Write-Host ""
Write-Host "7. 等待链路数据上报到 Jaeger..." -ForegroundColor Yellow
for ($i = 10; $i -gt 0; $i--) {
    Write-Host "   倒计时: $i 秒" -NoNewline
    Start-Sleep -Seconds 1
    Write-Host "`r" -NoNewline
}
Write-Host "   ✅ 等待完成                    " -ForegroundColor Green

# 8. 打开 Jaeger UI
Write-Host ""
Write-Host "8. 打开 Jaeger UI..." -ForegroundColor Yellow
$jaegerUrl = "http://14.103.153.242:16686/search?service=ShunFengParcel"
Start-Process $jaegerUrl
Write-Host "   ✅ 已打开浏览器" -ForegroundColor Green

# 9. 显示说明
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "✅ 测试完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "在 Jaeger UI 中查看:" -ForegroundColor Yellow
Write-Host "  1. 确认 Service 选择的是 'ShunFengParcel'" -ForegroundColor White
Write-Host "  2. 查找 'jaeger-connection-test' 测试链路" -ForegroundColor White
Write-Host "  3. 查找 '/admin/register' 请求链路" -ForegroundColor White
Write-Host ""
Write-Host "服务信息:" -ForegroundColor Yellow
Write-Host "  - PID: $($process.Id)" -ForegroundColor White
Write-Host "  - HTTP: http://localhost:8000" -ForegroundColor White
Write-Host "  - gRPC: localhost:9000" -ForegroundColor White
Write-Host ""
Write-Host "停止服务:" -ForegroundColor Yellow
Write-Host "  Stop-Process -Id $($process.Id) -Force" -ForegroundColor White
Write-Host ""
