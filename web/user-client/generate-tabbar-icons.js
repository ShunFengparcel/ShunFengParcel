// 生成 TabBar 图标的脚本
// 使用 Canvas 生成纯色图标

const fs = require('fs');
const { createCanvas } = require('canvas');

const icons = [
  { name: 'home', symbol: '🏠', label: '首页' },
  { name: 'search', symbol: '🔍', label: '查快递' },
  { name: 'send', symbol: '📦', label: '寄快递' },
  { name: 'user', symbol: '👤', label: '我的' }
];

const size = 81;
const grayColor = '#7A7E83';
const redColor = '#D81E06';

function generateIcon(name, symbol, color, suffix = '') {
  const canvas = createCanvas(size, size);
  const ctx = canvas.getContext('2d');
  
  // 设置字体
  ctx.font = `${size * 0.6}px Arial`;
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillStyle = color;
  
  // 绘制符号
  ctx.fillText(symbol, size / 2, size / 2);
  
  // 保存文件
  const buffer = canvas.toBuffer('image/png');
  const filename = `static/tabbar/${name}${suffix}.png`;
  fs.writeFileSync(filename, buffer);
  console.log(`Generated: ${filename}`);
}

// 生成所有图标
icons.forEach(icon => {
  generateIcon(icon.name, icon.symbol, grayColor, '');
  generateIcon(icon.name, icon.symbol, redColor, '-active');
});

console.log('All icons generated!');
