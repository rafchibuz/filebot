const sqlite3 = require('sqlite3').verbose();
const path = require('path');
const fs = require('fs');

// Создаем папку для базы данных если её нет
const dataDir = path.join(__dirname, '../../data');
if (!fs.existsSync(dataDir)) {
  fs.mkdirSync(dataDir, { recursive: true });
}

const dbPath = process.env.DB_PATH || path.join(dataDir, 'clips.db');
const db = new sqlite3.Database(dbPath);

// Инициализация базы данных
function initDatabase() {
  return new Promise((resolve, reject) => {
    db.serialize(() => {
      // Таблица для клипов
      db.run(`
        CREATE TABLE IF NOT EXISTS clips (
          id TEXT PRIMARY KEY,
          url TEXT NOT NULL,
          embed_url TEXT,
          broadcaster_id TEXT,
          broadcaster_name TEXT,
          creator_id TEXT,
          creator_name TEXT,
          video_id TEXT,
          game_id TEXT,
          language TEXT,
          title TEXT,
          view_count INTEGER,
          created_at TEXT,
          thumbnail_url TEXT,
          duration REAL,
          vod_offset INTEGER,
          is_featured INTEGER DEFAULT 0,
          streamer_login TEXT,
          streamer_display_name TEXT,
          streamer_profile_image TEXT,
          updated_at TEXT DEFAULT CURRENT_TIMESTAMP
        )
      `);

      // Таблица для стримеров
      db.run(`
        CREATE TABLE IF NOT EXISTS streamers (
          id TEXT PRIMARY KEY,
          login TEXT UNIQUE,
          display_name TEXT,
          profile_image_url TEXT,
          last_updated TEXT DEFAULT CURRENT_TIMESTAMP
        )
      `);

      // Индексы для оптимизации запросов
      db.run(`CREATE INDEX IF NOT EXISTS idx_clips_created_at ON clips(created_at)`);
      db.run(`CREATE INDEX IF NOT EXISTS idx_clips_streamer_login ON clips(streamer_login)`);
      db.run(`CREATE INDEX IF NOT EXISTS idx_clips_view_count ON clips(view_count)`);

      console.log('✅ База данных инициализирована');
      resolve();
    });
  });
}

// Сохранение клипов в базу данных
function saveClips(clips) {
  return new Promise((resolve, reject) => {
    const stmt = db.prepare(`
      INSERT OR REPLACE INTO clips (
        id, url, embed_url, broadcaster_id, broadcaster_name, creator_id, 
        creator_name, video_id, game_id, language, title, view_count, 
        created_at, thumbnail_url, duration, vod_offset, streamer_login,
        streamer_display_name, streamer_profile_image, updated_at
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `);

    db.serialize(() => {
      db.run('BEGIN TRANSACTION');
      
      clips.forEach(clip => {
        stmt.run([
          clip.id,
          clip.url,
          clip.embed_url,
          clip.broadcaster_id,
          clip.broadcaster_name,
          clip.creator_id,
          clip.creator_name,
          clip.video_id,
          clip.game_id,
          clip.language,
          clip.title,
          clip.view_count,
          clip.created_at,
          clip.thumbnail_url,
          clip.duration,
          clip.vod_offset,
          clip.streamer_login,
          clip.streamer_display_name,
          clip.streamer_profile_image,
          new Date().toISOString()
        ]);
      });
      
      db.run('COMMIT', (err) => {
        if (err) {
          console.error('❌ Ошибка сохранения клипов:', err);
          reject(err);
        } else {
          console.log(`✅ Сохранено ${clips.length} клипов в базу данных`);
          resolve();
        }
      });
    });

    stmt.finalize();
  });
}

// Получение клипов из базы данных
function getClips(filters = {}) {
  return new Promise((resolve, reject) => {
    let query = `
      SELECT * FROM clips 
      WHERE 1=1
    `;
    const params = [];

    // Фильтр по стримеру
    if (filters.streamer) {
      query += ` AND streamer_login = ?`;
      params.push(filters.streamer);
    }

    // Фильтр по периоду
    if (filters.daysBack) {
      const startDate = new Date();
      startDate.setDate(startDate.getDate() - filters.daysBack);
      query += ` AND created_at >= ?`;
      params.push(startDate.toISOString());
    }

    // Фильтр по минимальному количеству просмотров
    if (filters.minViews) {
      query += ` AND view_count >= ?`;
      params.push(filters.minViews);
    }

    // Сортировка
    const sortBy = filters.sortBy || 'created_at';
    const sortOrder = filters.sortOrder || 'DESC';
    query += ` ORDER BY ${sortBy} ${sortOrder}`;

    // Лимит
    if (filters.limit) {
      query += ` LIMIT ?`;
      params.push(filters.limit);
    }

    db.all(query, params, (err, rows) => {
      if (err) {
        console.error('❌ Ошибка получения клипов:', err);
        reject(err);
      } else {
        resolve(rows);
      }
    });
  });
}

// Получение статистики
function getClipsStats() {
  return new Promise((resolve, reject) => {
    const query = `
      SELECT 
        streamer_login,
        COUNT(*) as clips_count,
        SUM(view_count) as total_views,
        AVG(view_count) as avg_views,
        MAX(created_at) as last_clip_date
      FROM clips 
      GROUP BY streamer_login
      ORDER BY total_views DESC
    `;

    db.all(query, (err, rows) => {
      if (err) {
        reject(err);
      } else {
        resolve(rows);
      }
    });
  });
}

// Получение списка стримеров
function getStreamersFromDB() {
  return new Promise((resolve, reject) => {
    const query = `
      SELECT DISTINCT 
        streamer_login,
        streamer_display_name,
        streamer_profile_image
      FROM clips 
      ORDER BY streamer_display_name
    `;

    db.all(query, (err, rows) => {
      if (err) {
        reject(err);
      } else {
        resolve(rows);
      }
    });
  });
}

module.exports = {
  initDatabase,
  saveClips,
  getClips,
  getClipsStats,
  getStreamersFromDB
};