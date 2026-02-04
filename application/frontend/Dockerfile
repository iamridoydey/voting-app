# Stage 1: Build dependencies and production assets
FROM node:20-alpine3.22 AS builder

WORKDIR /app

# Copy dependency manifests
COPY package*.json ./
COPY tsconfig*.json ./
COPY vite.config.ts ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build production assets
RUN npm run build

# Stage 2: Serve with nginx
FROM nginx:1.27-alpine AS production

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Optional: custom nginx config
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
