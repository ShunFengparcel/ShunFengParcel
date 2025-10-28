# TabBar 图标下载脚本
# 此脚本会从 iconfont.cn 下载图标（需要手动操作）

Write-Host "==================================" -ForegroundColor Green
Write-Host "TabBar 图标下载指南" -ForegroundColor Green
Write-Host "==================================" -ForegroundColor Green
Write-Host ""

Write-Host "请按照以下步骤操作：" -ForegroundColor Yellow
Write-Host ""

Write-Host "1. 打开浏览器，访问：https://www.iconfont.cn/" -ForegroundColor Cyan
Write-Host ""

Write-Host "2. 搜索并下载以下图标：" -ForegroundColor Cyan
Write-Host "   - home (首页图标)" -ForegroundColor White
Write-Host "   - search (搜索图标)" -ForegroundColor White
Write-Host "   - package (包裹图标)" -ForegroundColor White
Write-Host "   - user (用户图标)" -ForegroundColor White
Write-Host ""

Write-Host "3. 每个图标下载两个颜色版本：" -ForegroundColor Cyan
Write-Host "   - 灰色版本：#7A7E83" -ForegroundColor White
Write-Host "   - 红色版本：#D81E06" -ForegroundColor White
Write-Host ""

Write-Host "4. 图标规格：" -ForegroundColor Cyan
Write-Host "   - 尺寸：81px × 81px" -ForegroundColor White
Write-Host "   - 格式：PNG" -ForegroundColor White
Write-Host "   - 背景：透明" -ForegroundColor White
Write-Host ""

Write-Host "5. 文件命名：" -ForegroundColor Cyan
Write-Host "   home.png, home-active.png" -ForegroundColor White
Write-Host "   search.png, search-active.png" -ForegroundColor White
Write-Host "   send.png, send-active.png" -ForegroundColor White
Write-Host "   user.png, user-active.png" -ForegroundColor White
Write-Host ""

$targetDir = "web\user-client\static\tabbar"
Write-Host "6. 将下载的图标放到目录：" -ForegroundColor Cyan
Write-Host "   $targetDir" -ForegroundColor White
Write-Host ""

# 检查目录是否存在
if (Test-Path $targetDir) {
    Write-Host "✓ 目标目录存在" -ForegroundColor Green
    
    # 检查现有文件
    $files = Get-ChildItem "$targetDir\*.png" -ErrorAction SilentlyContinue
    if ($files) {
        Write-Host ""
        Write-Host "当前目录中的PNG文件：" -ForegroundColor Yellow
        foreach ($file in $files) {
            $size = $file.Length
            if ($size -eq 0) {
                Write-Host "  ✗ $($file.Name) - 空文件 (0 字节)" -ForegroundColor Red
            } else {
                $sizeKB = [math]::Round($size / 1KB, 2)
                Write-Host "  ✓ $($file.Name) - $sizeKB KB" -ForegroundColor Green
            }
        }
    }
} else {
    Write-Host "✗ 目标目录不存在" -ForegroundColor Red
}

Write-Host ""
Write-Host "==================================" -ForegroundColor Green
Write-Host "快速链接：" -ForegroundColor Yellow
Write-Host ""
Write-Host "iconfont.cn: https://www.iconfont.cn/" -ForegroundColor Cyan
Write-Host "图片调整工具: https://www.iloveimg.com/zh-cn/resize-image" -ForegroundColor Cyan
Write-Host "在线PS: https://www.photopea.com/" -ForegroundColor Cyan
Write-Host ""

Write-Host "按任意键打开 iconfont.cn..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

Start-Process "https://www.iconfont.cn/"

Write-Host ""
Write-Host "图标下载完成后，运行以下命令重新编译：" -ForegroundColor Yellow
Write-Host "cd web/user-client" -ForegroundColor Cyan
Write-Host "npm run dev:mp-weixin" -ForegroundColor Cyan
