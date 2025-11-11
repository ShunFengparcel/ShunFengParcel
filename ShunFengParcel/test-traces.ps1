# 测试服务并生成 traces
Write-Host "=== Testing Service and Generating Traces ===" -ForegroundColor Cyan

$serviceUrl = "http://localhost:8000"
$testEndpoint = "$serviceUrl/helloworld/ShunFeng"

Write-Host "`n1. Checking if service is running..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri $testEndpoint -Method GET -UseBasicParsing -TimeoutSec 3
    Write-Host "   ✓ Service is running" -ForegroundColor Green
    Write-Host "   Response: $($response.Content)" -ForegroundColor White
} catch {
    Write-Host "   ✗ Service is not running or not responding" -ForegroundColor Red
    Write-Host "   Please start the service first with: kratos run" -ForegroundColor Yellow
    exit 1
}

Write-Host "`n2. Sending multiple requests to generate traces..." -ForegroundColor Yellow
$names = @("Alice", "Bob", "Charlie", "David", "Eva")
foreach ($name in $names) {
    try {
        $url = "$serviceUrl/helloworld/$name"
        $response = Invoke-RestMethod -Uri $url -Method GET
        Write-Host "   ✓ Request to /helloworld/$name - Status: OK" -ForegroundColor Green
        Start-Sleep -Milliseconds 500
    } catch {
        Write-Host "   ✗ Request to /helloworld/$name failed" -ForegroundColor Red
    }
}

Write-Host "`n3. Waiting for traces to be exported..." -ForegroundColor Yellow
Write-Host "   (Traces are batched and sent periodically)" -ForegroundColor Gray
Start-Sleep -Seconds 3

Write-Host "`n4. Check Jaeger UI at: http://14.103.153.242:16686" -ForegroundColor Cyan
Write-Host "   - Service: ShunFengParcel" -ForegroundColor White
Write-Host "   - Operation: /helloworld.v1.Greeter/SayHello" -ForegroundColor White
Write-Host "   - Lookback: Last Hour" -ForegroundColor White

Write-Host "`n=== Test Complete ===" -ForegroundColor Cyan


