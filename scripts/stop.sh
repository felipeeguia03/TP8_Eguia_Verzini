#!/bin/bash

# Script para detener todos los servicios
# Uso: ./scripts/stop.sh

set -e

echo "🛑 Deteniendo servicios..."

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Verificar que docker-compose esté instalado
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose no está instalado"
    exit 1
fi

echo -e "${BLUE}🐳 Deteniendo contenedores...${NC}"
docker-compose down

echo -e "${GREEN}✅ Servicios detenidos${NC}"

