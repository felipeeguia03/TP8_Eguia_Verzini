# Configuración con GitHub Actions y Render (100% Gratis)

Esta guía explica cómo configurar el proyecto para usar **GitHub Actions** para construir y subir imágenes Docker a **GitHub Container Registry (GHCR)**, y luego desplegar esas imágenes en **Render** - todo completamente gratis.

## 🎯 Stack Tecnológico

- ✅ **GitHub Actions**: CI/CD gratuito
- ✅ **GitHub Container Registry (GHCR)**: Almacenamiento de imágenes Docker gratuito
- ✅ **Render**: Hosting gratuito con tier free

## 📋 Prerrequisitos

1. **Cuenta de GitHub** con un repositorio del proyecto
2. **Cuenta de Render** (gratis)
3. **Docker** instalado localmente (solo para pruebas, opcional)

## 🔐 Paso 1: Configurar Secrets en GitHub

> 📖 **¿Necesitas ayuda paso a paso?** Lee la [Guía Detallada de Secrets](./GUIA-SECRETS-GITHUB.md)

### Ruta Rápida

1. Ve a tu repositorio en GitHub
2. Click en **Settings** (en la parte superior del repositorio)
3. En el menú lateral, busca **Secrets and variables** → **Actions**
4. Click en **New repository secret**

### Secrets a Configurar

#### 1. `NEXT_PUBLIC_API_URL` (Opcional pero Recomendado)

- **Name**: `NEXT_PUBLIC_API_URL`
- **Value**: URL de tu backend (ej: `https://tu-backend.onrender.com`)
- **Cuándo**: Puedes agregarlo ahora con `http://localhost:8080` y cambiarlo después

#### 2. `RENDER_BACKEND_WEBHOOK_URL` (Opcional)

- **Name**: `RENDER_BACKEND_WEBHOOK_URL`
- **Value**: Webhook URL del backend en Render
- **Cuándo**: Después de crear el servicio backend en Render
- **Dónde obtenerlo**: Render → Tu servicio → Settings → Manual Deploy → Webhook URL

#### 3. `RENDER_FRONTEND_WEBHOOK_URL` (Opcional)

- **Name**: `RENDER_FRONTEND_WEBHOOK_URL`
- **Value**: Webhook URL del frontend en Render
- **Cuándo**: Después de crear el servicio frontend en Render
- **Dónde obtenerlo**: Render → Tu servicio → Settings → Manual Deploy → Webhook URL

### ⚠️ Importante

- **`GITHUB_TOKEN`**: NO necesitas crearlo. GitHub lo proporciona automáticamente.
- **URL Directa**: `https://github.com/[tu-usuario]/[tu-repositorio]/settings/secrets/actions`
- **Ver guía completa**: [GUIA-SECRETS-GITHUB.md](./GUIA-SECRETS-GITHUB.md)

## 🚀 Paso 2: Configurar Render

> 📖 **¿Necesitas ayuda paso a paso detallada?** Lee la [Guía Completa de Render](./GUIA-CONFIGURAR-RENDER.md)

### Crear Servicio Backend en Render

1. Ve a [Render Dashboard](https://dashboard.render.com)
2. Click en **New** → **Web Service**
3. Selecciona **Deploy an existing image from a registry**
4. Configura:

   - **Registry**: `Docker Registry`
   - **Image URL**: `ghcr.io/[tu-usuario-github]/backend:latest`
     - Ejemplo: `ghcr.io/felipeeguia03/backend:latest`
   - **Registry Credentials**:
     - **Username**: Tu usuario de GitHub
     - **Password**: Un Personal Access Token de GitHub con scopes:
       - `read:packages`
       - `write:packages` (si quieres que Render pueda hacer pull)
   - **Name**: `tp8-backend`
   - **Environment**: `Docker`
   - **Region**: Elige la más cercana (Oregon, Frankfurt, etc.)
   - **Plan**: `Free` (o el que prefieras)

5. **Environment Variables**:

   ```
   PORT=8080
   MYSQL_HOST=tu-mysql-host
   MYSQL_PORT=3306
   MYSQL_USER=tu-usuario
   MYSQL_PASSWORD=tu-password
   MYSQL_DB=tu-database
   MYSQL_SSL=false
   ```

6. **Health Check Path**: `/healthz`

### Crear Servicio Frontend en Render

1. Ve a [Render Dashboard](https://dashboard.render.com)
2. Click en **New** → **Web Service**
3. Selecciona **Deploy an existing image from a registry**
4. Configura:

   - **Registry**: `Docker Registry`
   - **Image URL**: `ghcr.io/[tu-usuario-github]/frontend:latest`
     - Ejemplo: `ghcr.io/felipeeguia03/frontend:latest`
   - **Registry Credentials**:
     - **Username**: Tu usuario de GitHub
     - **Password**: El mismo Personal Access Token de GitHub
   - **Name**: `tp8-frontend`
   - **Environment**: `Docker`
   - **Region**: La misma que el backend
   - **Plan**: `Free`

5. **Environment Variables**:
   ```
   PORT=8080
   NEXT_PUBLIC_API_URL=https://tu-backend.onrender.com
   ```
   (Reemplaza `tu-backend.onrender.com` con la URL real de tu backend en Render)

### Crear Personal Access Token de GitHub (para Render)

Si necesitas un token para que Render acceda a GHCR:

1. Ve a GitHub → **Settings** → **Developer settings** → **Personal access tokens** → **Tokens (classic)**
2. Click en **Generate new token (classic)**
3. Configura:
   - **Note**: `Render GHCR Access`
   - **Expiration**: Elige según prefieras
   - **Scopes**: Marca solo:
     - ✅ `read:packages`
4. Click en **Generate token**
5. **Copia el token inmediatamente** (no podrás verlo después)
6. Úsalo como password en Render cuando configures el registry

## 🔄 Paso 3: Auto-deploy en Render

Para que Render actualice automáticamente cuando GitHub Actions suba nuevas imágenes:

### Opción 1: Webhook de Render (Recomendado)

1. En Render, ve a tu servicio → **Settings** → **Manual Deploy**
2. Copia el **Webhook URL** (debería verse algo como: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`)
3. Agrega los secrets en GitHub:
   - `RENDER_BACKEND_WEBHOOK_URL`: Webhook del servicio backend
   - `RENDER_FRONTEND_WEBHOOK_URL`: Webhook del servicio frontend

El workflow `.github/workflows/notify-render.yml` se ejecutará automáticamente después de que las imágenes se suban a GHCR.

### Opción 2: Render API

Si prefieres usar la API directamente:

1. Obtén tu **API Key** de Render: [Account Settings](https://dashboard.render.com/account/api-keys)
2. Agrega `RENDER_API_KEY` como secret en GitHub
3. Obtén el **Service ID** de cada servicio (está en la URL del servicio)
4. Actualiza el workflow para usar la API

### Opción 3: Auto-deploy desde GHCR (Sin webhook)

Render puede configurarse para hacer pull automático de imágenes. En **Settings** → **Auto-Deploy**, configura:

- **Pull latest image**: Habilitado
- **Pull schedule**: Cada X minutos (o manual)

## 📝 Estructura de Imágenes en GHCR

Las imágenes se suben a GHCR con el siguiente formato:

- **Backend**: `ghcr.io/[usuario-github]/backend:latest`
- **Frontend**: `ghcr.io/[usuario-github]/frontend:latest`

También se crean tags con el SHA del commit:

- `ghcr.io/[usuario-github]/backend:[sha]`
- `ghcr.io/[usuario-github]/frontend:[sha]`

### Ver tus imágenes en GHCR

1. Ve a tu repositorio en GitHub
2. En la barra lateral derecha, busca **Packages**
3. Ahí verás todas las imágenes Docker que has subido

## 🔍 Verificar el Despliegue

### Verificar imágenes en GHCR

1. Ve a tu repositorio en GitHub
2. Click en **Packages** (en la barra lateral)
3. Deberías ver `backend` y `frontend` como packages

O desde la línea de comandos:

```bash
# Listar imágenes (requiere autenticación)
gh auth login
gh api user/packages?package_type=container
```

### Verificar en Render

1. Ve al dashboard de Render
2. Revisa los logs del servicio
3. Verifica que el health check esté funcionando
4. Prueba acceder a las URLs de tus servicios

## 🐛 Troubleshooting

### Error: "unauthorized: authentication required" en GitHub Actions

- Verifica que el workflow tenga los permisos correctos:
  ```yaml
  permissions:
    contents: read
    packages: write
  ```
- Asegúrate de que `GITHUB_TOKEN` esté disponible (se proporciona automáticamente)

### Error: "image pull failed" en Render

- Verifica que la imagen exista en GHCR (ve a Packages en GitHub)
- Asegúrate de que el Personal Access Token tenga el scope `read:packages`
- Verifica que la URL de la imagen sea correcta: `ghcr.io/[usuario]/[imagen]:latest`
- **Importante**: Si el repositorio es privado, el token debe tener acceso al repositorio

### Error: "build failed" en GitHub Actions

- Revisa los logs del workflow en GitHub
- Verifica que los Dockerfiles estén correctos
- Asegúrate de que todas las dependencias estén disponibles
- Verifica que los paths en el workflow sean correctos (`./final/backend`, `./final/frontend`)

### Las imágenes no aparecen en Packages

- Espera unos minutos (puede haber un delay)
- Verifica que el workflow se haya ejecutado correctamente
- Asegúrate de que el push haya sido exitoso (revisa los logs)

### Render no puede hacer pull de imágenes privadas

Si tu repositorio es privado:

1. Asegúrate de que el Personal Access Token tenga acceso al repositorio
2. O haz el package público en GitHub:
   - Ve a Packages → Selecciona el package → Package settings → Change visibility → Public

## 💡 Tips y Mejores Prácticas

1. **Usa tags específicos para producción**: En lugar de `latest`, considera usar tags como `v1.0.0` o el SHA del commit
2. **Configura webhooks**: Para deploy automático cuando se suban nuevas imágenes
3. **Monitorea los logs**: Tanto en GitHub Actions como en Render
4. **Usa variables de entorno**: No hardcodees URLs ni credenciales
5. **Health checks**: Asegúrate de que tus servicios tengan endpoints de health check (`/healthz`)

## 📚 Recursos

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Render Documentation](https://render.com/docs)
- [Docker Build Push Action](https://github.com/docker/build-push-action)

## 🎉 ¡Listo!

Ahora tienes un pipeline CI/CD completamente gratuito:

- ✅ GitHub Actions construye y prueba tu código
- ✅ GHCR almacena tus imágenes Docker
- ✅ Render despliega tus servicios

Cada vez que hagas push a `main` o `master`, todo se actualiza automáticamente.
