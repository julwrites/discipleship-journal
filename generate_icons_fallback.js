const fs = require('fs');

const minimalPng = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==', 'base64');

fs.writeFileSync('web/public/pwa-192x192.png', minimalPng);
fs.writeFileSync('web/public/pwa-512x512.png', minimalPng);
fs.writeFileSync('web/public/apple-touch-icon.png', minimalPng);
fs.writeFileSync('web/public/favicon.ico', minimalPng); // Technically png but browsers often handle it
console.log("Created dummy icons");
