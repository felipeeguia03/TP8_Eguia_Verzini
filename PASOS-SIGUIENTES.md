# 🎯 Pasos Siguientes - Tu Workflow Funcionó ✅

¡Felicitaciones! Tu workflow se ejecutó correctamente en 1 minuto 27 segundos. Ahora sigue estos pasos para completar el despliegue.

## ✅ Paso 1: Verificar que las Imágenes Están en GHCR

### Verificar en GitHub

1. Ve a tu perfil de GitHub: `https://github.com/[tu-usuario]?tab=packages`
2. O desde tu repositorio, busca **"Packages"** en la barra lateral derecha
3. Deberías ver:
   - ✅ `backend` (package)
   - ✅ `frontend` (package)

**Si no los ves inmediatamente**: Espera 1-2 minutos, a veces hay un pequeño delay.

---

## 🔐 Paso 2: Configurar Secrets en GitHub (Si Aún No Lo Hiciste)

Ve a tu repositorio → **Settings** → **Secrets and variables** → **Actions**

Agrega estos secrets (si aún no los tienes):

### 1. `NEXT_PUBLIC_API_URL` (Opcional pero Recomendado)

- **Value**: Por ahora puedes usar `http://localhost:8080`
- Lo cambiarás después con la URL real del backend en Render

### 2. `RENDER_BACKEND_WEBHOOK_URL` (Opcional - Lo agregarás después)

- Lo obtendrás cuando crees el servicio backend en Render

### 3. `RENDER_FRONTEND_WEBHOOK_URL` (Opcional - Lo agregarás después)

- Lo obtendrás cuando crees el servicio frontend en Render

**Guía detallada**: [GUIA-SECRETS-GITHUB.md](./GUIA-SECRETS-GITHUB.md)

---

## 🚀 Paso 3: Crear Personal Access Token de GitHub (Para Render)

Render necesita un token para acceder a las imágenes en GHCR.

### Cómo Crearlo

1. GitHub → Tu avatar (arriba derecha) → **Settings**
2. **Developer settings** → **Personal access tokens** → **Tokens (classic)**
3. **Generate new token (classic)**
4. Configura:
   - **Note**: `Render GHCR Access`
   - **Expiration**: Elige según prefieras
   - **Scopes**: Marca SOLO ✅ `read:packages`
5. **Generate token**
6. **⚠️ COPIA EL TOKEN INMEDIATAMENTE** (empieza con `ghp_`)
   - No podrás verlo después
   - Guárdalo temporalmente en un lugar seguro

**Este token lo usarás como "Password" en Render.**

---

## 🎯 Paso 4: Configurar Render - Backend

### 4.1 Crear Servicio Backend

1. Ve a [dashboard.render.com](https://dashboard.render.com)
2. **New +** → **Web Service**
3. Selecciona: **"Deploy an existing image from a registry"**

### 4.2 Configurar

**Registry**: `Docker Registry`

**Image URL**:

```
ghcr.io/[tu-usuario-github]/backend:latest
```

Reemplaza `[tu-usuario-github]` con tu usuario real.

**Registry Credentials**:

- **Username**: Tu usuario de GitHub
- **Password**: El token que creaste en el Paso 3 (el que empieza con `ghp_`)

**Name**: `tp8-backend`

**Environment**: `Docker`

**Region**: Elige la más cercana (Oregon, Frankfurt, etc.)

**Plan**: `Free`

### 4.3 Variables de Entorno

Agrega estas variables:

| Key              | Value               |
| ---------------- | ------------------- |
| `PORT`           | `8080`              |
| `MYSQL_HOST`     | `tu-mysql-host`     |
| `MYSQL_PORT`     | `3306`              |
| `MYSQL_USER`     | `tu-usuario-mysql`  |
| `MYSQL_PASSWORD` | `tu-password-mysql` |
| `MYSQL_DB`       | `tu-database`       |
| `MYSQL_SSL`      | `false`             |

**Reemplaza los valores con tus datos reales de MySQL.**

### 4.4 Health Check

**Health Check Path**: `/healthz`

### 4.5 Crear y Obtener Webhook

1. Click en **Create Web Service**
2. Espera a que Render despliegue (puede tardar 2-5 minutos)
3. Una vez creado, ve a **Settings** → **Manual Deploy**
4. Copia el **Webhook URL** (algo como: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`)
5. **Guarda esta URL** - la necesitarás para GitHub Secrets

---

## 🎨 Paso 5: Configurar Render - Frontend

### 5.1 Crear Servicio Frontend

1. En Render, **New +** → **Web Service**
2. **"Deploy an existing image from a registry"**

### 5.2 Configurar

**Registry**: `Docker Registry`

**Image URL**:

```
ghcr.io/[tu-usuario-github]/frontend:latest
```

**Registry Credentials**:

- **Username**: Tu usuario de GitHub (el mismo)
- **Password**: El mismo token del Paso 3

**Name**: `tp8-frontend`

**Environment**: `Docker`

**Region**: La misma que el backend

**Plan**: `Free`

### 5.3 Variables de Entorno

**IMPORTANTE**: Necesitas la URL del backend que creaste en el Paso 4.

| Key                   | Value                             |
| --------------------- | --------------------------------- |
| `PORT`                | `8080`                            |
| `NEXT_PUBLIC_API_URL` | `https://tu-backend.onrender.com` |

**⚠️ Reemplaza `tu-backend.onrender.com` con la URL real de tu backend en Render.**

- La encontrarás en el dashboard de Render, en la página del servicio backend
- Debería verse algo como: `https://tp8-backend-xxxxx.onrender.com`

### 5.4 Crear y Obtener Webhook

1. Click en **Create Web Service**
2. Espera a que Render despliegue
3. Ve a **Settings** → **Manual Deploy**
4. Copia el **Webhook URL**
5. **Guarda esta URL**

---

## 🔗 Paso 6: Agregar Webhooks a GitHub Secrets

Ahora que tienes los webhooks de Render:

1. Ve a tu repositorio → **Settings** → **Secrets and variables** → **Actions**
2. Agrega:
   - `RENDER_BACKEND_WEBHOOK_URL` → URL del webhook del backend
   - `RENDER_FRONTEND_WEBHOOK_URL` → URL del webhook del frontend

**Esto habilitará el auto-deploy automático cuando hagas push a GitHub.**

---

## ✅ Paso 7: Verificar que Todo Funciona

### Verificar Backend

1. En Render, ve a tu servicio backend
2. Click en **Logs** (menú lateral)
3. Deberías ver logs de la aplicación iniciándose
4. Prueba acceder a: `https://tu-backend.onrender.com/healthz`
   - Debería responder: `{"status":"ok"}`

### Verificar Frontend

1. En Render, ve a tu servicio frontend
2. Click en **Logs**
3. Revisa que no haya errores
4. Prueba acceder a: `https://tu-frontend.onrender.com`
   - Debería mostrar tu aplicación

---

## 🎉 ¡Listo! Flujo Completo

Ahora tienes configurado:

```
Push a GitHub
    ↓
GitHub Actions construye imágenes (1 min 27 seg) ✅
    ↓
GitHub Actions sube a GHCR ✅
    ↓
GitHub Actions notifica a Render (webhook) ✅
    ↓
Render hace pull de GHCR ✅
    ↓
Render despliega servicios ✅
    ↓
¡Tu app está en producción! 🚀
```

---

## 📚 Guías de Referencia

- **Configurar Secrets**: [GUIA-SECRETS-GITHUB.md](./GUIA-SECRETS-GITHUB.md)
- **Configurar Render**: [GUIA-CONFIGURAR-RENDER.md](./GUIA-CONFIGURAR-RENDER.md)
- **Verificar Workflow**: [COMO-VERIFICAR-WORKFLOW.md](./COMO-VERIFICAR-WORKFLOW.md)

---

## 🐛 Si Algo No Funciona

### Error: "Failed to pull image" en Render

**Solución**:

- Verifica que el token de GitHub tenga scope `read:packages`
- Verifica que la URL de la imagen sea correcta: `ghcr.io/[usuario]/backend:latest`
- Verifica que las credenciales sean correctas

### Error: "Image not found"

**Solución**:

- Verifica que las imágenes existan en GHCR (ve a Packages en GitHub)
- Espera unos minutos si acabas de subir las imágenes

### Frontend no puede conectar con Backend

**Solución**:

- Verifica que `NEXT_PUBLIC_API_URL` sea la URL correcta del backend
- Verifica que el backend tenga CORS configurado correctamente

---

## 🎯 Resumen de URLs que Necesitarás

Guarda esta información:

- **Backend en Render**: `https://tu-backend.onrender.com`
- **Frontend en Render**: `https://tu-frontend.onrender.com`
- **Backend Webhook**: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`
- **Frontend Webhook**: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`
- **Imagen Backend**: `ghcr.io/[tu-usuario]/backend:latest`
- **Imagen Frontend**: `ghcr.io/[tu-usuario]/frontend:latest`

---

¡Sigue estos pasos y tendrás tu aplicación desplegada en Render! 🚀
