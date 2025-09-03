const twitchService = require('../services/twitchService');
const { getClips, getClipsStats, getStreamersFromDB } = require('../database/db');

// Тестовые данные для демонстрации
function getTestClips() {
  const now = new Date();
  return [
    {
      id: 'test-clip-1',
      url: 'https://www.twitch.tv/ravshanN/clip/TestClip1',
      title: 'Невероятный камбэк в CS:GO!',
      view_count: 15243,
      created_at: new Date(now - 2 * 60 * 60 * 1000).toISOString(), // 2 часа назад
      thumbnail_url: 'https://via.placeholder.com/320x180/9146ff/ffffff?text=RavshanN+Clip',
      duration: 45.5,
      streamer_login: 'RavshanN',
      streamer_display_name: 'RavshanN',
      streamer_profile_image: 'https://via.placeholder.com/50x50/9146ff/ffffff?text=R',
      creator_name: 'viewer123'
    },
    {
      id: 'test-clip-2', 
      url: 'https://www.twitch.tv/steel/clip/TestClip2',
      title: 'Мастерский играм steel',
      view_count: 8967,
      created_at: new Date(now - 5 * 60 * 60 * 1000).toISOString(), // 5 часов назад
      thumbnail_url: 'https://via.placeholder.com/320x180/f50057/ffffff?text=steel+Clip',
      duration: 32.2,
      streamer_login: 'steel',
      streamer_display_name: 'steel',
      streamer_profile_image: 'https://via.placeholder.com/50x50/f50057/ffffff?text=S',
      creator_name: 'fan456'
    },
    {
      id: 'test-clip-3',
      url: 'https://www.twitch.tv/renatko/clip/TestClip3',
      title: 'Смешной момент на стриме',
      view_count: 12456,
      created_at: new Date(now - 8 * 60 * 60 * 1000).toISOString(), // 8 часов назад
      thumbnail_url: 'https://via.placeholder.com/320x180/00ff88/ffffff?text=renatko+Clip',
      duration: 28.7,
      streamer_login: 'renatko',
      streamer_display_name: 'renatko', 
      streamer_profile_image: 'https://via.placeholder.com/50x50/00ff88/ffffff?text=R',
      creator_name: 'moderator789'
    }
  ];
}

// Тестовая статистика
function getTestStats() {
  return [
    { streamer_login: 'RavshanN', clips_count: 1, total_views: 15243, avg_views: 15243 },
    { streamer_login: 'steel', clips_count: 1, total_views: 8967, avg_views: 8967 },
    { streamer_login: 'renatko', clips_count: 1, total_views: 12456, avg_views: 12456 }
  ];
}

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

      console.log('🔍 Получен запрос клипов с фильтрами:', filters);

      const clips = await getClips(filters);
      const stats = await getClipsStats();

      console.log(`📦 Найдено клипов: ${clips.length}`);

      // Если клипов нет и API ключи не настроены, возвращаем тестовые данные
      if (clips.length === 0 && process.env.TWITCH_CLIENT_ID === 'your_client_id_here') {
        const testClips = getTestClips();
        const testStats = getTestStats();
        
        // Применяем фильтр по стримеру к тестовым данным
        let filteredTestClips = testClips;
        if (filters.streamer) {
          filteredTestClips = testClips.filter(clip => clip.streamer_login === filters.streamer);
        }
        
        res.json({
          clips: filteredTestClips,
          stats: testStats,
          filters: filters,
          total: filteredTestClips.length,
          timestamp: new Date().toISOString(),
          isTestData: true
        });
        return;
      }

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
      // Если API ключи не настроены, возвращаем тестовых стримеров
      if (process.env.TWITCH_CLIENT_ID === 'your_client_id_here') {
        const testStreamers = [
          { streamer_login: 'RavshanN', streamer_display_name: 'RavshanN' },
          { streamer_login: 'steel', streamer_display_name: 'steel' },
          { streamer_login: 'renatko', streamer_display_name: 'renatko' }
        ];
        res.json(testStreamers);
        return;
      }

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