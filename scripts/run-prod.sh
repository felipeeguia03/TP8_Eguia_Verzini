#!/bin/bash

# Script para ejecutar el proyecto en modo PROD localmente
# Uso: ./scripts/run-prod.sh

set -e

echo "🚀 Iniciando proyecto en modo PROD..."

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Verificar que docker-compose esté instalado
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose no está instalado"
    exit 1
fi

echo -e "${BLUE}📦 Configurando variables de entorno para PROD...${NC}"

# Crear archivo .env.prod si no existe
if [ ! -f .env.prod ]; then
    cat > .env.prod << EOF
# PROD Environment Variables
NODE_ENV=production
NEXT_PUBLIC_API_URL=http://localhost:8080
PORT=3000

# Backend
BACKEND_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=
DB_NAME=final_clj4
DB_SSLMODE=disable
EOF
    echo -e "${GREEN}✅ Archivo .env.prod creado${NC}"
fi

echo -e "${BLUE}🐳 Iniciando servicios con Docker Compose...${NC}"

# Ejecutar docker-compose con el archivo de configuración
docker-compose --env-file .env.prod up --build

echo -e "${GREEN}✅ Servicios PROD iniciados${NC}"
echo -e "${BLUE}📍 Frontend: http://localhost:3000${NC}"
echo -e "${BLUE}📍 Backend: http://localhost:8080${NC}"

