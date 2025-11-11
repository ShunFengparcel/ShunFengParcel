# ShunFengParcel 服务启动脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "启动 ShunFengParcel 服务" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 停止旧服务
Write-Host "1. 停止旧服务..." -ForegroundColor Yellow
$oldProcess = Get-Process -Name "ShunFengParcel" -ErrorAction SilentlyContinue
if ($oldProcess) {
    Stop-Process -Name "ShunFengParcel" -Force
    Start-Sleep -Seconds 1
    Write-Host "   ✅ 旧服务已停止" -ForegroundColor Green
} else {
    Write-Host "   ℹ️  没有运行中的服务" -ForegroundColor Gray
}

# 启动新服务
Write-Host ""
Write-Host "2. 启动新服务..." -ForegroundColor Yellow
Start-Process -FilePath ".\bin\ShunFengParcel.exe" `
    -ArgumentList "-conf", ".\configs" `
    -WorkingDirectory "." `
    -WindowStyle Hidden

Start-Sleep -Seconds 3

# 检查服务状态
$process = Get-Process -Name "ShunFengParcel" -ErrorAction SilentlyContinue
if ($process) {
    Write-Host "   ✅ 服务启动成功 (PID: $($process.Id))" -ForegroundColor Green
} else {
    Write-Host "   ❌ 服务启动失败" -ForegroundColor Red
    exit 1
}

# 测试服务
Write-Host ""
Write-Host "3. 测试服务..." -ForegroundColor Yellow
try {
    $response = curl.exe -X POST http://localhost:8000/admin/register `
        -H "Content-Type: application/json" `
        -d '{"username":"testuser","password":"test123"}' `
        -s 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "   ✅ 服务响应正常" -ForegroundColor Green
        Write-Host "   响应: $response" -ForegroundColor Gray
    } else {
        Write-Host "   ⚠️  服务可能还在初始化" -ForegroundColor Yellow
    }
} catch {
    Write-Host "   ⚠️  无法连接到服务" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "服务信息" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host "HTTP 服务: http://localhost:8000" -ForegroundColor Cyan
Write-Host "gRPC 服务: localhost:9000" -ForegroundColor Cyan
Write-Host "Jaeger UI: http://14.103.153.242:16686/search" -ForegroundColor Cyan
Write-Host ""
Write-Host "查看日志: Get-Process ShunFengParcel" -ForegroundColor Gray
Write-Host "停止服务: Stop-Process -Name ShunFengParcel -Force" -ForegroundColor Gray
Write-Host ""

