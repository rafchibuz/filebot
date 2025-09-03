# Multi-stage build for production
FROM node:18-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./
COPY client/package*.json ./client/

# Install dependencies
RUN npm install
RUN cd client && npm install

# Copy source code
COPY . .

# Build client
RUN cd client && npm run build

# Production stage
FROM node:18-alpine AS production

WORKDIR /app

# Copy package files and install only production dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy built application
COPY --from=builder /app/server ./server
COPY --from=builder /app/client/build ./client/build
COPY --from=builder /app/.env.example ./.env.example

# Create data directory
RUN mkdir -p data

# Expose ports
EXPOSE 3001

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:3001/api/health', (res) => process.exit(res.statusCode === 200 ? 0 : 1))"

# Start application
CMD ["node", "server/server.js"]