/**
 * TabBar 图标生成脚本
 * 使用 canvas 生成简单的PNG图标
 * 
 * 安装依赖：npm install canvas
 * 运行：node generate-icons.js
 */

const fs = require('fs');
const path = require('path');

// 检查是否安装了 canvas
let Canvas;
try {
  Canvas = require('canvas');
} catch (e) {
  console.log('❌ 未安装 canvas 模块');
  console.log('请先安装：npm install canvas');
  console.log('');
  console.log('如果安装失败，可以使用在线工具手动下载图标：');
  console.log('https://www.iconfont.cn/');
  process.exit(1);
}

const { createCanvas } = Canvas;

// 配置
const SIZE = 81;
const GRAY = '#7A7E83';
const RED = '#D81E06';
const OUTPUT_DIR = path.join(__dirname, 'static', 'tabbar');

// 确保输出目录存在
if (!fs.existsSync(OUTPUT_DIR)) {
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

/**
 * 绘制首页图标（房子）
 */
function drawHome(ctx, color) {
  ctx.strokeStyle = color;
  ctx.fillStyle = color;
  ctx.lineWidth = 3;
  ctx.lineCap = 'round';
  ctx.lineJoin = 'round';

  // 屋顶
  ctx.beginPath();
  ctx.moveTo(20, 45);
  ctx.lineTo(40.5, 25);
  ctx.lineTo(61, 45);
  ctx.stroke();

  // 房子主体
  ctx.strokeRect(25, 45, 31, 31);

  // 门
  ctx.fillRect(35, 60, 11, 16);
}

/**
 * 绘制搜索图标（放大镜）
 */
function drawSearch(ctx, color) {
  ctx.strokeStyle = color;
  ctx.lineWidth = 3;
  ctx.lineCap = 'round';

  // 圆圈
  ctx.beginPath();
  ctx.arc(35, 35, 15, 0, Math.PI * 2);
  ctx.stroke();

  // 手柄
  ctx.beginPath();
  ctx.moveTo(46, 46);
  ctx.lineTo(58, 58);
  ctx.stroke();
}

/**
 * 绘制包裹图标（盒子）
 */
function drawPackage(ctx, color) {
  ctx.strokeStyle = color;
  ctx.fillStyle = color;
  ctx.lineWidth = 3;
  ctx.lineCap = 'round';
  ctx.lineJoin = 'round';

  // 盒子
  ctx.strokeRect(25, 35, 31, 26);

  // 盖子
  ctx.beginPath();
  ctx.moveTo(25, 35);
  ctx.lineTo(40.5, 25);
  ctx.lineTo(56, 35);
  ctx.stroke();

  // 中线
  ctx.beginPath();
  ctx.moveTo(40.5, 25);
  ctx.lineTo(40.5, 61);
  ctx.stroke();
}

/**
 * 绘制用户图标（人像）
 */
function drawUser(ctx, color) {
  ctx.strokeStyle = color;
  ctx.fillStyle = color;
  ctx.lineWidth = 3;
  ctx.lineCap = 'round';

  // 头
  ctx.beginPath();
  ctx.arc(40.5, 32, 10, 0, Math.PI * 2);
  ctx.stroke();

  // 身体
  ctx.beginPath();
  ctx.arc(40.5, 65, 18, Math.PI, 0, true);
  ctx.stroke();
}

/**
 * 生成图标
 */
function generateIcon(name, drawFunc, color, suffix) {
  const canvas = createCanvas(SIZE, SIZE);
  const ctx = canvas.getContext('2d');

  // 清空画布
  ctx.clearRect(0, 0, SIZE, SIZE);

  // 绘制图标
  drawFunc(ctx, color);

  // 保存文件
  const filename = `${name}${suffix}.png`;
  const filepath = path.join(OUTPUT_DIR, filename);
  const buffer = canvas.toBuffer('image/png');
  fs.writeFileSync(filepath, buffer);

  console.log(`✓ 生成: ${filename}`);
}

/**
 * 主函数
 */
function main() {
  console.log('开始生成 TabBar 图标...\n');

  const icons = [
    { name: 'home', draw: drawHome, label: '首页' },
    { name: 'search', draw: drawSearch, label: '搜索' },
    { name: 'send', draw: drawPackage, label: '寄件' },
    { name: 'user', draw: drawUser, label: '用户' }
  ];

  icons.forEach(icon => {
    console.log(`生成 ${icon.label} 图标...`);
    generateIcon(icon.name, icon.draw, GRAY, '');
    generateIcon(icon.name, icon.draw, RED, '-active');
  });

  console.log('\n✅ 所有图标生成完成！');
  console.log(`输出目录: ${OUTPUT_DIR}`);
  console.log('\n下一步：');
  console.log('1. 检查生成的图标');
  console.log('2. 重新编译项目');
  console.log('3. 在微信开发者工具中查看效果');
}

// 运行
main();
