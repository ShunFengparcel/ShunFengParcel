# 启动 ShunFengParcel 服务（自动处理端口占用）

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "启动 ShunFengParcel 服务" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 1. 停止所有旧进程
Write-Host "1. 清理旧进程..." -ForegroundColor Yellow
$oldProcesses = Get-Process -Name "ShunFengParcel" -ErrorAction SilentlyContinue
if ($oldProcesses) {
    $oldProcesses | ForEach-Object {
        Write-Host "   停止进程 PID: $($_.Id)" -ForegroundColor Gray
        Stop-Process -Id $_.Id -Force
    }
    Start-Sleep -Seconds 2
    Write-Host "   ✅ 已清理" -ForegroundColor Green
} else {
    Write-Host "   ℹ️  没有旧进程" -ForegroundColor Gray
}

# 2. 检查并释放端口
Write-Host ""
Write-Host "2. 检查端口占用..." -ForegroundColor Yellow

$ports = @(8000, 9000)
foreach ($port in $ports) {
    $conn = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
    if ($conn) {
        Write-Host "   端口 $port 被占用，正在释放..." -ForegroundColor Yellow
        $conn | ForEach-Object {
            Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue
        }
        Start-Sleep -Seconds 1
        Write-Host "   ✅ 端口 $port 已释放" -ForegroundColor Green
    } else {
        Write-Host "   ✅ 端口 $port 可用" -ForegroundColor Green
    }
}

# 3. 编译服务
Write-Host ""
Write-Host "3. 编译服务..." -ForegroundColor Yellow
go build -o bin/ShunFengParcel.exe ./cmd/ShunFengParcel
if ($LASTEXITCODE -eq 0) {
    Write-Host "   ✅ 编译成功" -ForegroundColor Green
} else {
    Write-Host "   ❌ 编译失败" -ForegroundColor Red
    exit 1
}

# 4. 启动服务
Write-Host ""
Write-Host "4. 启动服务..." -ForegroundColor Yellow
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "服务正在启动，请查看下方日志..." -ForegroundColor Green
Write-Host "按 Ctrl+C 停止服务" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 在当前窗口运行，这样可以看到日志
& ".\bin\ShunFengParcel.exe"
