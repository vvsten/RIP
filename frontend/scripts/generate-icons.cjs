// Простой скрипт для генерации PWA иконок
// Требует: npm install sharp
const sharp = require('sharp');
const fs = require('fs');
const path = require('path');

const publicDir = path.join(__dirname, '../public');

// Создаем простую иконку с грузовиком
const createIcon = async (size, filename) => {
  const svg = `
    <svg width="${size}" height="${size}" xmlns="http://www.w3.org/2000/svg">
      <rect width="${size}" height="${size}" fill="#0d6efd" rx="${size * 0.2}"/>
      <text x="${size / 2}" y="${size * 0.65}" font-size="${size * 0.4}" 
            font-family="Arial, sans-serif" font-weight="bold" 
            text-anchor="middle" fill="white">🚚</text>
    </svg>
  `;
  
  await sharp(Buffer.from(svg))
    .png()
    .toFile(path.join(publicDir, filename));
  
  console.log(`Created ${filename}`);
};

(async () => {
  try {
    await createIcon(192, 'pwa-192x192.png');
    await createIcon(512, 'pwa-512x512.png');
    console.log('Icons generated successfully!');
  } catch (error) {
    console.error('Error generating icons:', error);
    console.log('Note: Install sharp with: npm install sharp');
  }
})();

