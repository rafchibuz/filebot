const express = require('express');
const cors = require('cors');
const path = require('path');
require('dotenv').config();

const twitchService = require('./services/twitchService');
const clipsController = require('./controllers/clipsController');
const { initDatabase } = require('./database/db');

const app = express();
const PORT = process.env.PORT || 3001;

// Middleware
app.use(cors());
app.use(express.json());

// Initialize database
initDatabase();

// Routes
app.get('/api/clips', clipsController.getClips);
app.get('/api/clips/refresh', clipsController.refreshClips);
app.get('/api/streamers', clipsController.getStreamers);

// Health check
app.get('/api/health', (req, res) => {
  res.json({ status: 'OK', timestamp: new Date().toISOString() });
});

// Serve static files from React build
if (process.env.NODE_ENV === 'production') {
  app.use(express.static(path.join(__dirname, '../client/build')));
  
  app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, '../client/build/index.html'));
  });
}

// Start server
app.listen(PORT, () => {
  console.log(`🚀 Сервер запущен на порту ${PORT}`);
  console.log(`📺 Отслеживаемые стримеры: ${process.env.STREAMERS}`);
  
  // Запускаем первоначальную загрузку клипов
  twitchService.startClipsUpdater();
});

module.exports = app;