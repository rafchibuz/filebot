import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Box,
  AppBar,
  Toolbar,
  Grid,
  CircularProgress,
  Alert,
  IconButton,
  Tooltip
} from '@mui/material';
import {
  Refresh,
  PlayArrow
} from '@mui/icons-material';
import ClipCard from './components/ClipCard';
import FilterPanel from './components/FilterPanel';
import StatsPanel from './components/StatsPanel';
import { clipsAPI } from './services/api';
import './App.css';

function App() {
  const [clips, setClips] = useState([]);
  const [stats, setStats] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    daysBack: 7,
    streamer: '',
    minViews: 0,
    sortBy: 'created_at',
    sortOrder: 'DESC',
    limit: 50
  });

  // Загрузка клипов
  const loadClips = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await clipsAPI.getClips(filters);
      setClips(data.clips);
      setStats(data.stats);
    } catch (err) {
      setError('Ошибка загрузки клипов: ' + err.message);
      console.error('Ошибка загрузки клипов:', err);
    } finally {
      setLoading(false);
    }
  };

  // Обновление клипов
  const refreshClips = async () => {
    try {
      setLoading(true);
      await clipsAPI.refreshClips(filters.daysBack);
      await loadClips();
    } catch (err) {
      setError('Ошибка обновления клипов: ' + err.message);
    }
  };

  // Загрузка при монтировании и изменении фильтров
  useEffect(() => {
    loadClips();
  }, [filters]);

  // Обработка изменения фильтров
  const handleFilterChange = (newFilters) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      {/* Header */}
      <AppBar position="static" sx={{ background: 'linear-gradient(45deg, #9146ff 30%, #f50057 90%)' }}>
        <Toolbar>
          <PlayArrow sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Twitch Clips Aggregator
          </Typography>
          <Tooltip title="Обновить клипы">
            <IconButton color="inherit" onClick={refreshClips} disabled={loading}>
              <Refresh />
            </IconButton>
          </Tooltip>
        </Toolbar>
      </AppBar>

      <Container maxWidth="xl" sx={{ mt: 3, mb: 3 }}>
        {/* Панель фильтров */}
        <FilterPanel 
          filters={filters}
          onFilterChange={handleFilterChange}
          onRefresh={refreshClips}
          loading={loading}
        />

        {/* Панель статистики */}
        {stats.length > 0 && (
          <Box sx={{ mb: 3 }}>
            <StatsPanel stats={stats} />
          </Box>
        )}

        {/* Ошибки */}
        {error && (
          <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {/* Загрузка */}
        {loading && (
          <Box display="flex" justifyContent="center" my={4}>
            <CircularProgress size={60} />
          </Box>
        )}

        {/* Список клипов */}
        {!loading && clips.length === 0 && (
          <Alert severity="info">
            Клипы не найдены. Попробуйте изменить фильтры или обновить данные.
          </Alert>
        )}

        {!loading && clips.length > 0 && (
          <Grid container spacing={3}>
            {clips.map((clip) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={clip.id}>
                <ClipCard clip={clip} />
              </Grid>
            ))}
          </Grid>
        )}
      </Container>
    </Box>
  );
}

export default App;