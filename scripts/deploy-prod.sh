#!/bin/bash

# Script para desplegar a PROD
# Uso: ./scripts/deploy-prod.sh

set -e

echo "🚀 Desplegando a PROD..."

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar que estamos en la branch main
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${YELLOW}⚠️  Estás en la branch: $CURRENT_BRANCH${NC}"
    echo -e "${YELLOW}⚠️  Cambiando a branch main...${NC}"
    git checkout main
fi

echo -e "${RED}⚠️  ATENCIÓN: Estás desplegando a PRODUCCIÓN${NC}"
read -p "¿Estás seguro? (escribe 'yes' para continuar): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    echo "❌ Deploy cancelado"
    exit 1
fi

echo -e "${BLUE}📝 Verificando cambios sin commitear...${NC}"

# Verificar si hay cambios sin commitear
if ! git diff-index --quiet HEAD --; then
    echo -e "${YELLOW}⚠️  Hay cambios sin commitear${NC}"
    read -p "¿Quieres hacer commit? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Mensaje del commit: " COMMIT_MSG
        git add .
        git commit -m "$COMMIT_MSG"
    else
        echo "❌ Abortado. Haz commit de tus cambios primero."
        exit 1
    fi
fi

echo -e "${BLUE}📤 Haciendo push a branch main...${NC}"
git push origin main

echo -e "${GREEN}✅ Push completado${NC}"
echo -e "${BLUE}⏳ Esperando a que GitHub Actions construya las imágenes...${NC}"
echo -e "${BLUE}🔗 Verifica el progreso en: https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/actions${NC}"

echo -e "${GREEN}✅ Deploy a PROD iniciado${NC}"
echo -e "${BLUE}📍 El deploy se completará automáticamente en Render${NC}"

