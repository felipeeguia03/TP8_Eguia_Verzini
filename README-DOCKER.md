# Dockerización del Proyecto TP8

Este proyecto ha sido dockerizado para facilitar el despliegue y desarrollo local.

## 📋 Prerrequisitos

- Docker Desktop instalado
- Docker Compose instalado (incluido en Docker Desktop)

## 🚀 Desarrollo Local con Docker

### Opción 1: Docker Compose (Recomendado)

Ejecuta todos los servicios con un solo comando:

```bash
docker-compose up --build
```

Esto iniciará:

- **MySQL** en el puerto 3306
- **Backend** (Go) en el puerto 8080
- **Frontend** (Next.js) en el puerto 3000

### Opción 2: Docker Individual

#### Backend

```bash
cd final/backend
docker build -t tp8-backend .
docker run -p 8080:8080 \
  -e MYSQL_HOST=localhost \
  -e MYSQL_PORT=3306 \
  -e MYSQL_USER=root \
  -e MYSQL_PASSWORD="" \
  -e MYSQL_DB=final \
  tp8-backend
```

#### Frontend

```bash
cd final/frontend
docker build -t tp8-frontend --build-arg NEXT_PUBLIC_API_URL=http://localhost:8080 .
docker run -p 3000:8080 \
  -e NEXT_PUBLIC_API_URL=http://localhost:8080 \
  tp8-frontend
```

## 🏗️ Estructura de Dockerfiles

### Backend Dockerfile

- **Multi-stage build**: Compila en una imagen y ejecuta en otra más pequeña
- **Base**: `golang:1.22-alpine` para build, `alpine:latest` para runtime
- **Puerto**: 8080

### Frontend Dockerfile

- **Multi-stage build**: Dependencies → Build → Production
- **Base**: `node:20-alpine`
- **Standalone mode**: Next.js en modo standalone para optimizar el tamaño
- **Puerto**: 8080

## ☁️ Despliegue en Azure

El pipeline de Azure DevOps (`azure-pipelines.yml`) ha sido actualizado para:

1. **Construir imágenes Docker** durante el build
2. **Subir imágenes a Azure Container Registry (ACR)**
3. **Desplegar contenedores** en Azure App Service

### Configuración Requerida

1. **Crear Azure Container Registry**:

   ```bash
   az acr create --resource-group rg-tp05-ingsoft3-2025 \
     --name acrtp05verzinieguia \
     --sku Basic
   ```

2. **Crear conexión de Docker Registry en Azure DevOps**:

   - Ve a **Project Settings** → **Service connections** → **New service connection**
   - Selecciona **Docker Registry**
   - Elige **Azure Container Registry**
   - Selecciona tu suscripción y ACR
   - Nombre la conexión: `acr-connection` (o el que prefieras)
   - **Actualiza** `azure-pipelines.yml` para usar esta conexión en lugar de `$(azureServiceConnection)` en las tareas Docker

3. **Actualizar variables en Azure DevOps**:

   - `acrName`: Nombre de tu ACR (debe coincidir con el creado)
   - Actualiza la variable en el pipeline o directamente en el YAML

4. **Configurar permisos**:

   - El service principal debe tener rol "AcrPush" en el ACR
   - O usa **Admin user** habilitado en ACR (Settings → Access keys)

5. **Habilitar Admin User en ACR** (alternativa más simple):
   ```bash
   az acr update -n acrtp05verzinieguia --admin-enabled true
   ```

## 📝 Notas Importantes

- Las contraseñas de MySQL están en texto plano en el pipeline. **Deberías moverlas a Azure Key Vault**.
- El frontend necesita `NEXT_PUBLIC_API_URL` en build time, por eso se pasa como `--build-arg`.
- Los contenedores se reinician automáticamente después de cada despliegue.

## 🔧 Troubleshooting

### Error: "Cannot connect to Docker daemon"

- Asegúrate de que Docker Desktop esté ejecutándose

### Error: "Port already in use"

- Cambia los puertos en `docker-compose.yml` o detén los servicios que los usan

### Error: "ACR login failed"

- Verifica que el service connection tenga los permisos correctos
- Revisa que el nombre del ACR sea correcto en las variables

## 📚 Recursos

- [Docker Documentation](https://docs.docker.com/)
- [Azure Container Registry](https://docs.microsoft.com/azure/container-registry/)
- [Azure App Service con Contenedores](https://docs.microsoft.com/azure/app-service/containers/)
