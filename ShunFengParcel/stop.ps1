# ShunFengParcel 服务停止脚本

Write-Host "停止 ShunFengParcel 服务..." -ForegroundColor Yellow

$process = Get-Process -Name "ShunFengParcel" -ErrorAction SilentlyContinue
if ($process) {
    Stop-Process -Name "ShunFengParcel" -Force
    Start-Sleep -Seconds 1
    
    $stillRunning = Get-Process -Name "ShunFengParcel" -ErrorAction SilentlyContinue
    if ($stillRunning) {
        Write-Host "❌ 服务停止失败" -ForegroundColor Red
        exit 1
    } else {
        Write-Host "✅ 服务已停止" -ForegroundColor Green
    }
} else {
    Write-Host "ℹ️  服务未运行" -ForegroundColor Gray
}

