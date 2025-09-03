import React, { useState } from 'react';
import {
  Card,
  CardContent,
  CardMedia,
  Typography,
  Box,
  Chip,
  Avatar,
  Stack,
  IconButton,
  Dialog,
  DialogContent,
  DialogTitle,
  Tooltip,
  Link
} from '@mui/material';
import {
  PlayArrow,
  Visibility,
  Schedule,
  OpenInNew,
  Close
} from '@mui/icons-material';
import { formatDistanceToNow, format } from 'date-fns';
import { ru } from 'date-fns/locale';
import ReactPlayer from 'react-player';

function ClipCard({ clip }) {
  const [dialogOpen, setDialogOpen] = useState(false);

  const formatViews = (views) => {
    if (views >= 1000000) {
      return (views / 1000000).toFixed(1) + 'M';
    } else if (views >= 1000) {
      return (views / 1000).toFixed(1) + 'K';
    }
    return views.toString();
  };

  const formatDuration = (duration) => {
    const minutes = Math.floor(duration / 60);
    const seconds = Math.floor(duration % 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  const createdAt = new Date(clip.created_at);

  return (
    <>
      <Card 
        sx={{ 
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          transition: 'transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out',
          '&:hover': {
            transform: 'translateY(-4px)',
            boxShadow: 6
          }
        }}
      >
        {/* Thumbnail */}
        <Box sx={{ position: 'relative', cursor: 'pointer' }} onClick={() => setDialogOpen(true)}>
          <CardMedia
            component="img"
            height="200"
            image={clip.thumbnail_url}
            alt={clip.title}
            sx={{ objectFit: 'cover' }}
          />
          <Box
            sx={{
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              background: 'rgba(0,0,0,0.3)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              opacity: 0,
              transition: 'opacity 0.2s ease-in-out',
              '&:hover': { opacity: 1 }
            }}
          >
            <PlayArrow sx={{ fontSize: 60, color: 'white' }} />
          </Box>
          <Chip
            label={formatDuration(clip.duration)}
            size="small"
            sx={{
              position: 'absolute',
              bottom: 8,
              right: 8,
              backgroundColor: 'rgba(0,0,0,0.8)',
              color: 'white'
            }}
          />
        </Box>

        <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
          {/* Заголовок */}
          <Typography 
            variant="h6" 
            component="h3" 
            sx={{ 
              mb: 1,
              fontSize: '1rem',
              lineHeight: 1.3,
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical',
              overflow: 'hidden'
            }}
          >
            {clip.title}
          </Typography>

          {/* Информация о стримере */}
          <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
            <Avatar 
              src={clip.streamer_profile_image} 
              sx={{ width: 24, height: 24 }}
            />
            <Typography variant="body2" color="primary" fontWeight="bold">
              {clip.streamer_display_name}
            </Typography>
          </Stack>

          {/* Метрики */}
          <Stack direction="row" spacing={1} sx={{ mb: 1, flexWrap: 'wrap', gap: 0.5 }}>
            <Chip
              icon={<Visibility />}
              label={formatViews(clip.view_count)}
              size="small"
              variant="outlined"
              color="secondary"
            />
            <Chip
              icon={<Schedule />}
              label={formatDistanceToNow(createdAt, { addSuffix: true, locale: ru })}
              size="small"
              variant="outlined"
            />
          </Stack>

          {/* Создатель клипа */}
          {clip.creator_name && (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 'auto' }}>
              Создал: {clip.creator_name}
            </Typography>
          )}

          {/* Дата создания */}
          <Typography variant="caption" color="text.secondary">
            {format(createdAt, 'dd.MM.yyyy HH:mm', { locale: ru })}
          </Typography>
        </CardContent>
      </Card>

      {/* Диалог с видео */}
      <Dialog
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        maxWidth="md"
        fullWidth
        PaperProps={{
          sx: {
            backgroundColor: 'background.paper',
            backgroundImage: 'none'
          }
        }}
      >
        <DialogTitle sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Avatar src={clip.streamer_profile_image} sx={{ width: 32, height: 32 }} />
            <Box>
              <Typography variant="h6" component="div">
                {clip.title}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {clip.streamer_display_name}
              </Typography>
            </Box>
          </Box>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Tooltip title="Открыть на Twitch">
              <IconButton 
                component={Link} 
                href={clip.url} 
                target="_blank" 
                rel="noopener noreferrer"
                color="primary"
              >
                <OpenInNew />
              </IconButton>
            </Tooltip>
            <IconButton onClick={() => setDialogOpen(false)}>
              <Close />
            </IconButton>
          </Box>
        </DialogTitle>
        <DialogContent>
          <Box sx={{ position: 'relative', paddingTop: '56.25%' /* 16:9 aspect ratio */ }}>
            <ReactPlayer
              url={clip.url}
              width="100%"
              height="100%"
              style={{ position: 'absolute', top: 0, left: 0 }}
              controls
              playing={dialogOpen}
            />
          </Box>
          <Box sx={{ mt: 2 }}>
            <Stack direction="row" spacing={1} sx={{ mb: 1 }}>
              <Chip
                icon={<Visibility />}
                label={`${clip.view_count.toLocaleString()} просмотров`}
                color="secondary"
                variant="outlined"
              />
              <Chip
                icon={<Schedule />}
                label={format(new Date(clip.created_at), 'dd.MM.yyyy HH:mm', { locale: ru })}
                variant="outlined"
              />
            </Stack>
            {clip.creator_name && (
              <Typography variant="body2" color="text.secondary">
                Создатель клипа: {clip.creator_name}
              </Typography>
            )}
          </Box>
        </DialogContent>
      </Dialog>
    </>
  );
}

export default ClipCard;