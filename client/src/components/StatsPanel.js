import React from 'react';
import {
  Paper,
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Avatar,
  Stack,
  Chip
} from '@mui/material';
import {
  TrendingUp,
  Visibility,
  VideoLibrary
} from '@mui/icons-material';

function StatsPanel({ stats }) {
  const formatViews = (views) => {
    if (views >= 1000000) {
      return (views / 1000000).toFixed(1) + 'M';
    } else if (views >= 1000) {
      return (views / 1000).toFixed(1) + 'K';
    }
    return views?.toString() || '0';
  };

  const totalClips = stats.reduce((sum, stat) => sum + stat.clips_count, 0);
  const totalViews = stats.reduce((sum, stat) => sum + (stat.total_views || 0), 0);

  return (
    <Paper sx={{ p: 3 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <TrendingUp sx={{ mr: 1 }} />
        <Typography variant="h6">Статистика</Typography>
      </Box>

      {/* Общая статистика */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        <Grid item xs={6} md={3}>
          <Card variant="outlined">
            <CardContent sx={{ textAlign: 'center', py: 2 }}>
              <Typography variant="h4" color="primary">
                {totalClips}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Всего клипов
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6} md={3}>
          <Card variant="outlined">
            <CardContent sx={{ textAlign: 'center', py: 2 }}>
              <Typography variant="h4" color="secondary">
                {formatViews(totalViews)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Всего просмотров
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6} md={3}>
          <Card variant="outlined">
            <CardContent sx={{ textAlign: 'center', py: 2 }}>
              <Typography variant="h4" color="success.main">
                {stats.length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Стримеров
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6} md={3}>
          <Card variant="outlined">
            <CardContent sx={{ textAlign: 'center', py: 2 }}>
              <Typography variant="h4" color="warning.main">
                {totalClips > 0 ? Math.round(totalViews / totalClips) : 0}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Средние просмотры
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Статистика по стримерам */}
      <Typography variant="subtitle1" sx={{ mb: 2, fontWeight: 'bold' }}>
        По стримерам:
      </Typography>
      <Grid container spacing={2}>
        {stats.map((stat) => (
          <Grid item xs={12} md={4} key={stat.streamer_login}>
            <Card variant="outlined">
              <CardContent>
                <Stack direction="row" spacing={2} alignItems="center">
                  <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h6" color="primary">
                      {stat.clips_count}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      клипов
                    </Typography>
                  </Box>
                  <Box sx={{ flexGrow: 1 }}>
                    <Typography variant="subtitle2" fontWeight="bold">
                      {stat.streamer_login}
                    </Typography>
                    <Stack direction="row" spacing={1} sx={{ mt: 0.5 }}>
                      <Chip
                        icon={<Visibility />}
                        label={formatViews(stat.total_views || 0)}
                        size="small"
                        variant="outlined"
                        color="secondary"
                      />
                      <Chip
                        label={`~${formatViews(stat.avg_views || 0)}`}
                        size="small"
                        variant="outlined"
                      />
                    </Stack>
                  </Box>
                </Stack>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Paper>
  );
}

export default StatsPanel;