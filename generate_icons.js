const fs = require('fs');
const { createCanvas } = require('canvas');

function createIcon(size, name) {
  const canvas = createCanvas(size, size);
  const ctx = canvas.getContext('2d');

  // Background
  ctx.fillStyle = '#4F46E5'; // Indigo-600
  ctx.fillRect(0, 0, size, size);

  // Text
  ctx.fillStyle = 'white';
  ctx.font = `bold ${size/2}px sans-serif`;
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillText('DJ', size/2, size/2);

  const buffer = canvas.toBuffer('image/png');
  fs.writeFileSync(`web/public/${name}`, buffer);
  console.log(`Generated ${name}`);
}

// Check if canvas is available, otherwise just copy a file or create a dummy file
try {
  createIcon(192, 'pwa-192x192.png');
  createIcon(512, 'pwa-512x512.png');
} catch (e) {
  console.error("Canvas not available, creating dummy files");
  // Fallback: Just create a valid 1x1 PNG or copy vite.svg if we can't generate
  // Since we can't easily convert svg to png without canvas/imagemagick in this restricted env
  // We will create a minimal valid PNG using base64

  const minimalPng = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==', 'base64');

  fs.writeFileSync('web/public/pwa-192x192.png', minimalPng);
  fs.writeFileSync('web/public/pwa-512x512.png', minimalPng);
}
