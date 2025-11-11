# Jaeger 集成测试脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Jaeger 集成测试" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 1. 检查 OpenTelemetry Collector 端口
Write-Host "1. 检查 OpenTelemetry Collector 连接..." -ForegroundColor Yellow
$collectorTest = Test-NetConnection -ComputerName 14.103.153.242 -Port 4317 -WarningAction SilentlyContinue
if ($collectorTest.TcpTestSucceeded) {
    Write-Host "   ✅ Collector 端口 4317 可访问" -ForegroundColor Green
} else {
    Write-Host "   ❌ Collector 端口 4317 不可访问" -ForegroundColor Red
    exit 1
}

# 2. 检查 Jaeger UI
Write-Host ""
Write-Host "2. 检查 Jaeger UI..." -ForegroundColor Yellow
$jaegerTest = Test-NetConnection -ComputerName 14.103.153.242 -Port 16686 -WarningAction SilentlyContinue
if ($jaegerTest.TcpTestSucceeded) {
    Write-Host "   ✅ Jaeger UI 端口 16686 可访问" -ForegroundColor Green
} else {
    Write-Host "   ❌ Jaeger UI 端口 16686 不可访问" -ForegroundColor Red
}

# 3. 重启服务
Write-Host ""
Write-Host "3. 重启 ShunFengParcel 服务..." -ForegroundColor Yellow
Stop-Process -Name "ShunFengParcel" -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 1

Start-Process -FilePath ".\bin\ShunFengParcel.exe" `
    -ArgumentList "-conf", ".\configs" `
    -WorkingDirectory "." `
    -WindowStyle Hidden

Start-Sleep -Seconds 3

$process = Get-Process -Name "ShunFengParcel" -ErrorAction SilentlyContinue
if ($process) {
    Write-Host "   ✅ 服务已启动 (PID: $($process.Id))" -ForegroundColor Green
} else {
    Write-Host "   ❌ 服务启动失败" -ForegroundColor Red
    exit 1
}

# 4. 调用接口产生追踪数据
Write-Host ""
Write-Host "4. 调用接口产生追踪数据..." -ForegroundColor Yellow
for ($i = 1; $i -le 5; $i++) {
    $response = curl.exe -X POST http://localhost:8000/admin/findadminbyusername `
        -H "Content-Type: application/json" `
        -d "{`"username`":`"user$i`",`"password`":`"pass123`"}" `
        -s
    Write-Host "   请求 $i : $response"
    Start-Sleep -Milliseconds 500
}

# 5. 等待追踪数据发送
Write-Host ""
Write-Host "5. 等待追踪数据发送到 Collector..." -ForegroundColor Yellow
Start-Sleep -Seconds 3
Write-Host "   ✅ 等待完成" -ForegroundColor Green

# 6. 检查 Jaeger 服务列表
Write-Host ""
Write-Host "6. 检查 Jaeger 服务列表..." -ForegroundColor Yellow
try {
    $services = curl.exe -s "http://14.103.153.242:16686/api/services" | ConvertFrom-Json
    Write-Host "   当前服务列表:" -ForegroundColor Cyan
    foreach ($service in $services.data) {
        if ($service -eq "ShunFengParcel") {
            Write-Host "   - $service" -ForegroundColor Green
        } else {
            Write-Host "   - $service" -ForegroundColor Gray
        }
    }
    
    if ($services.data -contains "ShunFengParcel") {
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Green
        Write-Host "✅✅✅ 成功！ShunFengParcel 已出现在 Jaeger 中！" -ForegroundColor Green
        Write-Host "========================================" -ForegroundColor Green
        Write-Host ""
        Write-Host "请访问 Jaeger UI 查看追踪数据:" -ForegroundColor Cyan
        Write-Host "http://14.103.153.242:16686/search" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "操作步骤:" -ForegroundColor Yellow
        Write-Host "1. 在 Service 下拉框中选择 'ShunFengParcel'" -ForegroundColor White
        Write-Host "2. 点击 'Find Traces' 按钮" -ForegroundColor White
        Write-Host "3. 查看追踪详情和 Span 信息" -ForegroundColor White
    } else {
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Red
        Write-Host "❌ ShunFengParcel 服务未出现在 Jaeger 中" -ForegroundColor Red
        Write-Host "========================================" -ForegroundColor Red
        Write-Host ""
        Write-Host "可能的原因:" -ForegroundColor Yellow
        Write-Host "1. OpenTelemetry Collector 未配置 Jaeger 导出器" -ForegroundColor White
        Write-Host "2. Collector 未重启以加载新配置" -ForegroundColor White
        Write-Host "3. Collector 到 Jaeger 的连接有问题" -ForegroundColor White
        Write-Host ""
        Write-Host "请按照 JAEGER_SETUP_GUIDE.md 中的说明:" -ForegroundColor Cyan
        Write-Host "1. 更新 Collector 配置文件" -ForegroundColor White
        Write-Host "2. 重启 OpenTelemetry Collector" -ForegroundColor White
        Write-Host "3. 重新运行此测试脚本" -ForegroundColor White
    }
} catch {
    Write-Host "   ❌ 无法连接到 Jaeger API" -ForegroundColor Red
    Write-Host "   错误: $_" -ForegroundColor Red
}

Write-Host ""

