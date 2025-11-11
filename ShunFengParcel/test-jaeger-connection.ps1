# 测试 Jaeger OTLP 端口连接
Write-Host "=== Testing Jaeger OTLP Connection ===" -ForegroundColor Cyan

$jaegerHost = "14.103.153.242"
$otlpPort = 4317
$httpPort = 16686

Write-Host "`n1. Testing Jaeger UI (HTTP $httpPort)..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://${jaegerHost}:${httpPort}" -TimeoutSec 5 -UseBasicParsing
    Write-Host "   ✓ Jaeger UI is accessible" -ForegroundColor Green
} catch {
    Write-Host "   ✗ Jaeger UI not accessible: $_" -ForegroundColor Red
}

Write-Host "`n2. Testing OTLP gRPC port ($otlpPort)..." -ForegroundColor Yellow
try {
    $tcpClient = New-Object System.Net.Sockets.TcpClient
    $tcpClient.Connect($jaegerHost, $otlpPort)
    if ($tcpClient.Connected) {
        Write-Host "   ✓ OTLP gRPC port $otlpPort is open" -ForegroundColor Green
        $tcpClient.Close()
    }
} catch {
    Write-Host "   ✗ OTLP gRPC port $otlpPort is NOT accessible" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
    Write-Host "`n   Possible solutions:" -ForegroundColor Yellow
    Write-Host "   - Check if Jaeger was started with OTLP support" -ForegroundColor White
    Write-Host "   - Try using Jaeger collector HTTP port 14268 instead" -ForegroundColor White
    Write-Host "   - Verify firewall rules allow port 4317" -ForegroundColor White
}

Write-Host "`n3. Checking alternative ports..." -ForegroundColor Yellow
$ports = @{
    "14268" = "Jaeger Collector HTTP"
    "14250" = "Jaeger Collector gRPC"
    "4318" = "OTLP HTTP"
}

foreach ($port in $ports.Keys) {
    try {
        $tcpClient = New-Object System.Net.Sockets.TcpClient
        $tcpClient.Connect($jaegerHost, $port)
        if ($tcpClient.Connected) {
            Write-Host "   ✓ Port $port ($($ports[$port])) is open" -ForegroundColor Green
            $tcpClient.Close()
        }
    } catch {
        Write-Host "   ✗ Port $port ($($ports[$port])) is closed" -ForegroundColor Gray
    }
}

Write-Host "`n=== Test Complete ===" -ForegroundColor Cyan


