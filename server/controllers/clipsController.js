const twitchService = require('../services/twitchService');
const { getClips, getClipsStats, getStreamersFromDB } = require('../database/db');

class ClipsController {
  // Получение клипов с фильтрацией
  async getClips(req, res) {
    try {
      const filters = {
        streamer: req.query.streamer,
        daysBack: parseInt(req.query.daysBack) || 7,
        minViews: parseInt(req.query.minViews) || 0,
        sortBy: req.query.sortBy || 'created_at',
        sortOrder: req.query.sortOrder || 'DESC',
        limit: parseInt(req.query.limit) || 50
      };

      const clips = await getClips(filters);
      const stats = await getClipsStats();

      res.json({
        clips,
        stats,
        filters: filters,
        total: clips.length,
        timestamp: new Date().toISOString()
      });
    } catch (error) {
      console.error('❌ Ошибка получения клипов:', error);
      res.status(500).json({
        error: 'Ошибка получения клипов',
        message: error.message
      });
    }
  }

  // Принудительное обновление клипов
  async refreshClips(req, res) {
    try {
      const daysBack = parseInt(req.query.daysBack) || 7;
      console.log(`🔄 Запущено принудительное обновление клипов за ${daysBack} дней`);
      
      const clips = await twitchService.getAllClips(daysBack);
      
      res.json({
        message: 'Клипы успешно обновлены',
        count: clips.length,
        timestamp: new Date().toISOString()
      });
    } catch (error) {
      console.error('❌ Ошибка обновления клипов:', error);
      res.status(500).json({
        error: 'Ошибка обновления клипов',
        message: error.message
      });
    }
  }

  // Получение списка стримеров
  async getStreamers(req, res) {
    try {
      const streamers = await getStreamersFromDB();
      res.json(streamers);
    } catch (error) {
      console.error('❌ Ошибка получения стримеров:', error);
      res.status(500).json({
        error: 'Ошибка получения стримеров',
        message: error.message
      });
    }
  }
}

module.exports = new ClipsController();