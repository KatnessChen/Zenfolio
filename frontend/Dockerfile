# Development stage
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Install dependencies for development
RUN apk add --no-cache \
    git \
    curl

# Install pnpm globally
RUN npm install -g pnpm

# Copy package files first for better caching
COPY package*.json ./
COPY pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install

# Copy source code
COPY . .

# Expose port
EXPOSE 3000

# Set environment variables for development
ENV NODE_ENV=development
ENV CHOKIDAR_USEPOLLING=true

# Start development server with hot reload
CMD ["pnpm", "run", "dev"]
