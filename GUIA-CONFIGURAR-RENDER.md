# 🚀 Guía Paso a Paso: Configurar Render

Esta guía te muestra exactamente cómo configurar tus servicios backend y frontend en Render para usar las imágenes Docker de GitHub Container Registry (GHCR).

## 📋 Prerrequisitos

Antes de empezar, asegúrate de tener:

1. ✅ Cuenta de Render creada ([dashboard.render.com](https://dashboard.render.com))
2. ✅ Repositorio en GitHub con el workflow funcionando
3. ✅ Imágenes Docker ya subidas a GHCR (o las subirás después)
4. ✅ Personal Access Token de GitHub con scope `read:packages` (para que Render acceda a GHCR)

---

## 🔑 Paso 0: Crear Personal Access Token de GitHub (para Render)

Render necesita un token de GitHub para acceder a las imágenes en GHCR.

### Cómo Crear el Token

1. Ve a GitHub → Click en tu avatar (arriba derecha) → **Settings**
2. En el menú lateral izquierdo, busca **Developer settings**
3. Click en **Personal access tokens** → **Tokens (classic)**
4. Click en **Generate new token (classic)**
5. Configura el token:
   - **Note**: `Render GHCR Access` (o el nombre que prefieras)
   - **Expiration**: Elige según prefieras (90 días, 1 año, etc.)
   - **Scopes**: Marca SOLO:
     - ✅ `read:packages` (para leer imágenes de GHCR)
6. Click en **Generate token** (abajo)
7. **⚠️ IMPORTANTE**: Copia el token inmediatamente (algo como `ghp_xxxxxxxxxxxxx`)
   - No podrás verlo después
   - Guárdalo en un lugar seguro temporalmente

**Este token lo usarás como "Password" en Render cuando configures el registry.**

---

## 🎯 Paso 1: Crear Servicio Backend en Render

### 1.1 Ir al Dashboard de Render

1. Ve a [dashboard.render.com](https://dashboard.render.com)
2. Si no has iniciado sesión, haz login con tu cuenta

### 1.2 Crear Nuevo Servicio

1. Click en el botón **New +** (arriba a la derecha)
2. Selecciona **Web Service**

### 1.3 Configurar el Servicio

Verás varias opciones. Selecciona:

**"Deploy an existing image from a registry"**

### 1.4 Configurar Registry e Imagen

Completa los siguientes campos:

#### Registry

- Selecciona: **Docker Registry**

#### Image URL

- Escribe: `ghcr.io/[tu-usuario-github]/backend:latest`
  - Ejemplo: `ghcr.io/felipeeguia03/backend:latest`
  - Reemplaza `[tu-usuario-github]` con tu usuario de GitHub

#### Registry Credentials

**Username**:

- Escribe tu usuario de GitHub (el mismo que en la URL de la imagen)

**Password**:

- Pega el Personal Access Token que creaste en el Paso 0
- Es el token que empieza con `ghp_`

### 1.5 Configurar Nombre y Región

#### Name

- Escribe: `tp8-backend` (o el nombre que prefieras)

#### Environment

- Selecciona: **Docker**

#### Region

- Elige la región más cercana a ti:
  - `Oregon` (US West)
  - `Frankfurt` (EU)
  - `Singapore` (Asia)
  - etc.

#### Plan

- Selecciona: **Free** (o el plan que prefieras)

### 1.6 Configurar Variables de Entorno

En la sección **Environment Variables**, agrega:

| Key              | Value           | Descripción                            |
| ---------------- | --------------- | -------------------------------------- |
| `PORT`           | `8080`          | Puerto donde corre la aplicación       |
| `MYSQL_HOST`     | `tu-mysql-host` | Host de tu base de datos MySQL         |
| `MYSQL_PORT`     | `3306`          | Puerto de MySQL                        |
| `MYSQL_USER`     | `tu-usuario`    | Usuario de MySQL                       |
| `MYSQL_PASSWORD` | `tu-password`   | Contraseña de MySQL                    |
| `MYSQL_DB`       | `tu-database`   | Nombre de la base de datos             |
| `MYSQL_SSL`      | `false`         | SSL para MySQL (false para desarrollo) |

**Nota**: Reemplaza los valores con tus datos reales de MySQL.

### 1.7 Configurar Health Check

En la sección **Health Check Path**:

- Escribe: `/healthz`

### 1.8 Crear el Servicio

1. Revisa toda la configuración
2. Click en **Create Web Service**
3. Render comenzará a hacer pull de la imagen y desplegarla
4. Espera unos minutos mientras Render descarga y despliega la imagen

### 1.9 Obtener Webhook URL (para Auto-deploy)

Una vez que el servicio esté creado:

1. Ve a tu servicio → **Settings** (en el menú lateral)
2. Busca la sección **Manual Deploy**
3. Busca **Webhook URL** o **Deploy Hook**
4. Copia la URL completa (algo como: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`)
5. **Guarda esta URL** - la necesitarás para el secret `RENDER_BACKEND_WEBHOOK_URL` en GitHub

---

## 🎨 Paso 2: Crear Servicio Frontend en Render

### 2.1 Crear Nuevo Servicio

1. En el dashboard de Render, click en **New +**
2. Selecciona **Web Service**
3. Selecciona **"Deploy an existing image from a registry"**

### 2.2 Configurar Registry e Imagen

#### Registry

- Selecciona: **Docker Registry**

#### Image URL

- Escribe: `ghcr.io/[tu-usuario-github]/frontend:latest`
  - Ejemplo: `ghcr.io/felipeeguia03/frontend:latest`
  - Reemplaza `[tu-usuario-github]` con tu usuario de GitHub

#### Registry Credentials

**Username**:

- Tu usuario de GitHub (el mismo de antes)

**Password**:

- El mismo Personal Access Token que usaste para el backend

### 2.3 Configurar Nombre y Región

#### Name

- Escribe: `tp8-frontend` (o el nombre que prefieras)

#### Environment

- Selecciona: **Docker**

#### Region

- Elige la **misma región** que el backend (recomendado)

#### Plan

- Selecciona: **Free**

### 2.4 Configurar Variables de Entorno

En la sección **Environment Variables**, agrega:

| Key                   | Value                             | Descripción                      |
| --------------------- | --------------------------------- | -------------------------------- |
| `PORT`                | `8080`                            | Puerto donde corre la aplicación |
| `NEXT_PUBLIC_API_URL` | `https://tu-backend.onrender.com` | URL del backend                  |

- ⚠️ **Importante**: Reemplaza `tu-backend.onrender.com` con la URL real de tu backend
- La URL la encontrarás en el dashboard de Render, en la página del servicio backend
- Debería verse algo como: `https://tp8-backend-xxxxx.onrender.com`

### 2.5 Crear el Servicio

1. Revisa toda la configuración
2. Click en **Create Web Service**
3. Espera a que Render despliegue la imagen

### 2.6 Obtener Webhook URL

1. Ve a tu servicio frontend → **Settings**
2. Busca **Manual Deploy** → **Webhook URL**
3. Copia la URL
4. **Guarda esta URL** - la necesitarás para el secret `RENDER_FRONTEND_WEBHOOK_URL` en GitHub

---

## ✅ Paso 3: Agregar Webhooks a GitHub Secrets

Ahora que tienes los webhooks de Render, agrégalos a GitHub:

1. Ve a tu repositorio en GitHub
2. **Settings** → **Secrets and variables** → **Actions**
3. Agrega:
   - `RENDER_BACKEND_WEBHOOK_URL` → URL del webhook del backend
   - `RENDER_FRONTEND_WEBHOOK_URL` → URL del webhook del frontend

**Ver guía completa**: [GUIA-SECRETS-GITHUB.md](./GUIA-SECRETS-GITHUB.md)

---

## 🔍 Paso 4: Verificar el Despliegue

### Verificar Backend

1. En Render, ve a tu servicio backend
2. Click en **Logs** (en el menú lateral)
3. Deberías ver logs de la aplicación iniciándose
4. Si hay errores, revísalos y corrígelos

### Verificar Frontend

1. En Render, ve a tu servicio frontend
2. Click en **Logs**
3. Revisa que no haya errores

### Probar los Servicios

1. En Render, cada servicio tiene una URL (algo como `https://tp8-backend-xxxxx.onrender.com`)
2. Prueba acceder a:
   - Backend: `https://tu-backend.onrender.com/healthz` (debería responder `{"status":"ok"}`)
   - Frontend: `https://tu-frontend.onrender.com` (debería mostrar tu aplicación)

---

## 🔄 Paso 5: Configurar Auto-Deploy (Opcional)

### Opción 1: Usar Webhooks (Recomendado)

Si ya agregaste los webhooks a GitHub Secrets, el workflow `.github/workflows/notify-render.yml` se ejecutará automáticamente después de cada push y notificará a Render para hacer deploy.

### Opción 2: Auto-Pull en Render

1. En Render, ve a tu servicio → **Settings**
2. Busca **Auto-Deploy**
3. Habilita **"Pull latest image"**
4. Configura la frecuencia (cada X minutos, o manual)

---

## 🐛 Troubleshooting

### Error: "Failed to pull image"

**Causa**: Render no puede acceder a la imagen en GHCR

**Solución**:

- Verifica que el Personal Access Token tenga el scope `read:packages`
- Verifica que la URL de la imagen sea correcta: `ghcr.io/[usuario]/backend:latest`
- Verifica que las credenciales (username y password) sean correctas
- Si el repositorio es privado, asegúrate de que el token tenga acceso al repositorio

### Error: "Image not found"

**Causa**: La imagen no existe en GHCR o el nombre es incorrecto

**Solución**:

- Verifica que el workflow de GitHub Actions haya ejecutado correctamente
- Ve a tu repositorio en GitHub → **Packages** (barra lateral) → Verifica que las imágenes existan
- Verifica que el nombre de la imagen sea exactamente: `ghcr.io/[usuario]/backend:latest`

### Error: "Connection refused" o "Cannot connect to database"

**Causa**: Variables de entorno incorrectas o base de datos no accesible

**Solución**:

- Verifica todas las variables de entorno en Render
- Asegúrate de que `MYSQL_HOST` sea accesible desde Render
- Verifica que `MYSQL_USER`, `MYSQL_PASSWORD` y `MYSQL_DB` sean correctos

### El servicio no inicia

**Causa**: Varias posibles (puerto incorrecto, variables faltantes, etc.)

**Solución**:

1. Ve a **Logs** en Render
2. Revisa los errores
3. Verifica que `PORT=8080` esté configurado
4. Verifica que todas las variables de entorno necesarias estén configuradas

### Frontend no puede conectar con Backend

**Causa**: `NEXT_PUBLIC_API_URL` incorrecto o CORS mal configurado

**Solución**:

- Verifica que `NEXT_PUBLIC_API_URL` en el frontend sea la URL correcta del backend
- Verifica que el backend tenga configurado CORS para permitir el origen del frontend
- Revisa los logs del frontend para ver errores de conexión

---

## 📝 Resumen de URLs y Credenciales

Guarda esta información:

### GitHub Container Registry (GHCR)

- **Registry**: `ghcr.io`
- **Backend Image**: `ghcr.io/[tu-usuario]/backend:latest`
- **Frontend Image**: `ghcr.io/[tu-usuario]/frontend:latest`
- **Username**: Tu usuario de GitHub
- **Password**: Personal Access Token con scope `read:packages`

### Render Services

- **Backend URL**: `https://tu-backend.onrender.com`
- **Frontend URL**: `https://tu-frontend.onrender.com`
- **Backend Webhook**: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`
- **Frontend Webhook**: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`

---

## 🎉 ¡Listo!

Una vez configurado todo:

1. ✅ GitHub Actions construye y sube imágenes a GHCR
2. ✅ Render despliega las imágenes automáticamente
3. ✅ Los servicios están disponibles en las URLs de Render
4. ✅ Auto-deploy funciona con webhooks

**Siguiente paso**: Prueba hacer un cambio en tu código, haz push a GitHub, y verifica que todo se actualice automáticamente! 🚀

---

## 📚 Recursos Adicionales

- [Render Documentation](https://render.com/docs)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Render Webhooks](https://render.com/docs/webhooks)
