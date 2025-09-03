import React, { useState, useEffect } from 'react';
import {
  Paper,
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Button,
  Grid,
  Typography,
  Chip,
  Stack
} from '@mui/material';
import { Refresh, FilterList } from '@mui/icons-material';
import { clipsAPI } from '../services/api';

function FilterPanel({ filters, onFilterChange, onRefresh, loading }) {
  const [streamers, setStreamers] = useState([]);

  // Загрузка списка стримеров
  useEffect(() => {
    const loadStreamers = async () => {
      try {
        const streamersData = await clipsAPI.getStreamers();
        setStreamers(streamersData);
      } catch (error) {
        console.error('Ошибка загрузки стримеров:', error);
        // Fallback - используем стандартный список стримеров
        setStreamers([
          { streamer_login: 'RavshanN', streamer_display_name: 'RavshanN' },
          { streamer_login: 'steel', streamer_display_name: 'steel' },
          { streamer_login: 'renatko', streamer_display_name: 'renatko' }
        ]);
      }
    };
    loadStreamers();
  }, []);

  const handleFilterChange = (field, value) => {
    onFilterChange({ [field]: value });
  };

  const resetFilters = () => {
    onFilterChange({
      daysBack: 7,
      streamer: '',
      minViews: 0,
      sortBy: 'created_at',
      sortOrder: 'DESC',
      limit: 50
    });
  };

  return (
    <Paper sx={{ p: 3, mb: 3 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <FilterList sx={{ mr: 1 }} />
        <Typography variant="h6">Фильтры</Typography>
      </Box>

      <Grid container spacing={2} alignItems="center">
        {/* Период времени */}
        <Grid item xs={12} sm={6} md={2}>
          <FormControl fullWidth size="small">
            <InputLabel>Период</InputLabel>
            <Select
              value={filters.daysBack}
              label="Период"
              onChange={(e) => handleFilterChange('daysBack', e.target.value)}
            >
              <MenuItem value={1}>1 день</MenuItem>
              <MenuItem value={3}>3 дня</MenuItem>
              <MenuItem value={7}>1 неделя</MenuItem>
              <MenuItem value={14}>2 недели</MenuItem>
              <MenuItem value={30}>1 месяц</MenuItem>
            </Select>
          </FormControl>
        </Grid>

        {/* Стример */}
        <Grid item xs={12} sm={6} md={2}>
          <FormControl fullWidth size="small">
            <InputLabel>Стример</InputLabel>
            <Select
              value={filters.streamer}
              label="Стример"
              onChange={(e) => handleFilterChange('streamer', e.target.value)}
            >
              <MenuItem value="">Все стримеры</MenuItem>
              {streamers.map((streamer) => (
                <MenuItem key={streamer.streamer_login} value={streamer.streamer_login}>
                  {streamer.streamer_display_name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>

        {/* Минимальные просмотры */}
        <Grid item xs={12} sm={6} md={2}>
          <TextField
            fullWidth
            size="small"
            label="Мин. просмотров"
            type="number"
            value={filters.minViews}
            onChange={(e) => handleFilterChange('minViews', parseInt(e.target.value) || 0)}
          />
        </Grid>

        {/* Сортировка */}
        <Grid item xs={12} sm={6} md={2}>
          <FormControl fullWidth size="small">
            <InputLabel>Сортировка</InputLabel>
            <Select
              value={filters.sortBy}
              label="Сортировка"
              onChange={(e) => handleFilterChange('sortBy', e.target.value)}
            >
              <MenuItem value="created_at">По времени</MenuItem>
              <MenuItem value="view_count">По просмотрам</MenuItem>
              <MenuItem value="title">По названию</MenuItem>
            </Select>
          </FormControl>
        </Grid>

        {/* Порядок сортировки */}
        <Grid item xs={12} sm={6} md={1}>
          <FormControl fullWidth size="small">
            <InputLabel>Порядок</InputLabel>
            <Select
              value={filters.sortOrder}
              label="Порядок"
              onChange={(e) => handleFilterChange('sortOrder', e.target.value)}
            >
              <MenuItem value="DESC">↓</MenuItem>
              <MenuItem value="ASC">↑</MenuItem>
            </Select>
          </FormControl>
        </Grid>

        {/* Кнопки */}
        <Grid item xs={12} md={3}>
          <Stack direction="row" spacing={1}>
            <Button
              variant="contained"
              startIcon={<Refresh />}
              onClick={onRefresh}
              disabled={loading}
              size="small"
            >
              Обновить
            </Button>
            <Button
              variant="outlined"
              onClick={resetFilters}
              size="small"
            >
              Сбросить
            </Button>
          </Stack>
        </Grid>
      </Grid>

      {/* Активные фильтры */}
      <Box sx={{ mt: 2 }}>
        <Stack direction="row" spacing={1} flexWrap="wrap">
          <Chip 
            label={`Период: ${filters.daysBack} дн.`} 
            size="small" 
            color="primary" 
            variant="outlined" 
          />
          {filters.streamer && (
            <Chip 
              label={`Стример: ${filters.streamer}`} 
              size="small" 
              color="secondary" 
              variant="outlined" 
            />
          )}
          {filters.minViews > 0 && (
            <Chip 
              label={`Мин. просмотров: ${filters.minViews}`} 
              size="small" 
              variant="outlined" 
            />
          )}
        </Stack>
      </Box>
    </Paper>
  );
}

export default FilterPanel;