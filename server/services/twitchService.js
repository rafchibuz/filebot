const axios = require('axios');
const cron = require('node-cron');
const { saveClips, getStreamersFromDB } = require('../database/db');

class TwitchService {
  constructor() {
    this.accessToken = null;
    this.tokenExpiry = null;
    this.streamers = (process.env.STREAMERS || 'RavshanN,steel,renatko').split(',');
  }

  // Получение токена доступа к Twitch API
  async getAccessToken() {
    if (this.accessToken && this.tokenExpiry && Date.now() < this.tokenExpiry) {
      return this.accessToken;
    }

    // Проверяем наличие API ключей
    if (!process.env.TWITCH_CLIENT_ID || !process.env.TWITCH_CLIENT_SECRET) {
      throw new Error('❌ TWITCH_CLIENT_ID и TWITCH_CLIENT_SECRET должны быть настроены в .env файле');
    }

    if (process.env.TWITCH_CLIENT_ID === 'your_client_id_here') {
      throw new Error('❌ Пожалуйста, замените your_client_id_here на ваш реальный Client ID в .env файле');
    }

    try {
      const response = await axios.post('https://id.twitch.tv/oauth2/token', {
        client_id: process.env.TWITCH_CLIENT_ID,
        client_secret: process.env.TWITCH_CLIENT_SECRET,
        grant_type: 'client_credentials'
      });

      this.accessToken = response.data.access_token;
      this.tokenExpiry = Date.now() + (response.data.expires_in * 1000) - 60000; // -1 минута для безопасности

      console.log('✅ Получен новый токен доступа Twitch API');
      return this.accessToken;
    } catch (error) {
      console.error('❌ Ошибка получения токена:', error.response?.data || error.message);
      throw error;
    }
  }

  // Получение информации о пользователе по имени
  async getUserByLogin(login) {
    const token = await this.getAccessToken();
    
    try {
      const response = await axios.get('https://api.twitch.tv/helix/users', {
        headers: {
          'Client-ID': process.env.TWITCH_CLIENT_ID,
          'Authorization': `Bearer ${token}`
        },
        params: {
          login: login
        }
      });

      return response.data.data[0];
    } catch (error) {
      console.error(`❌ Ошибка получения пользователя ${login}:`, error.response?.data || error.message);
      return null;
    }
  }

  // Получение клипов стримера за период
  async getClipsForStreamer(broadcasterId, startDate, endDate) {
    const token = await this.getAccessToken();
    
    try {
      const response = await axios.get('https://api.twitch.tv/helix/clips', {
        headers: {
          'Client-ID': process.env.TWITCH_CLIENT_ID,
          'Authorization': `Bearer ${token}`
        },
        params: {
          broadcaster_id: broadcasterId,
          started_at: startDate.toISOString(),
          ended_at: endDate.toISOString(),
          first: 100 // Максимум 100 клипов за запрос
        }
      });

      return response.data.data;
    } catch (error) {
      console.error(`❌ Ошибка получения клипов для стримера ${broadcasterId}:`, error.response?.data || error.message);
      return [];
    }
  }

  // Получение всех клипов за период
  async getAllClips(daysBack = 7) {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - daysBack);

    console.log(`🔍 Загружаем клипы за период с ${startDate.toLocaleDateString()} по ${endDate.toLocaleDateString()}`);

    const allClips = [];
    
    for (const streamerLogin of this.streamers) {
      try {
        const user = await this.getUserByLogin(streamerLogin);
        if (!user) {
          console.warn(`⚠️ Стример ${streamerLogin} не найден`);
          continue;
        }

        const clips = await this.getClipsForStreamer(user.id, startDate, endDate);
        
        // Добавляем информацию о стримере к каждому клипу
        const enrichedClips = clips.map(clip => ({
          ...clip,
          streamer_login: streamerLogin,
          streamer_display_name: user.display_name,
          streamer_profile_image: user.profile_image_url
        }));

        allClips.push(...enrichedClips);
        console.log(`✅ Загружено ${clips.length} клипов для ${streamerLogin}`);
        
        // Небольшая пауза между запросами
        await new Promise(resolve => setTimeout(resolve, 100));
      } catch (error) {
        console.error(`❌ Ошибка обработки стримера ${streamerLogin}:`, error.message);
      }
    }

    // Сортируем по времени создания (новые сначала)
    allClips.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

    console.log(`🎬 Всего загружено ${allClips.length} клипов`);
    
    // Сохраняем в базу данных
    await saveClips(allClips);
    
    return allClips;
  }

  // Запуск автоматического обновления клипов
  startClipsUpdater() {
    // Загружаем клипы при старте
    this.getAllClips().catch(console.error);

    // Настраиваем автоматическое обновление каждые 30 минут
    const updateInterval = process.env.UPDATE_INTERVAL || 30;
    cron.schedule(`*/${updateInterval} * * * *`, () => {
      console.log('🔄 Автоматическое обновление клипов...');
      this.getAllClips().catch(console.error);
    });

    console.log(`⏰ Автоматическое обновление настроено на каждые ${updateInterval} минут`);
  }
}

module.exports = new TwitchService();