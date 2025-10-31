#!/bin/bash
echo "🚀 Building frontend..."
npm run build

echo "📦 Serving static files..."
npx serve -s dist -l 3000