# 🚀 Configuración QA y PROD

## 📋 Resumen

Tu pipeline ahora soporta **2 ambientes**:
- **QA**: Branch `qa` → Despliega a servicios QA en Render
- **PROD**: Branch `main` → Despliega a servicios PROD en Render

## 🔄 Flujo

```
git push origin qa    →  Build imagen :qa  →  Deploy a Render QA
git push origin main  →  Build imagen :prod →  Deploy a Render PROD
```

## 📦 Imágenes Docker

Las imágenes se etiquetan según el ambiente:
- **QA**: `ghcr.io/usuario/backend:qa` y `ghcr.io/usuario/frontend:qa`
- **PROD**: `ghcr.io/usuario/backend:prod` y `ghcr.io/usuario/frontend:prod`

## 🛠️ Paso 1: Crear Branch QA

```bash
# Crear branch qa desde main
git checkout -b qa
git push origin qa
```

## 🛠️ Paso 2: Crear Servicios en Render

### 2.1 Backend QA

1. **Render Dashboard** → **New** → **Web Service**
2. **Configuración**:
   - **Name**: `tp8-backend-qa`
   - **Environment**: `Docker`
   - **Registry**: `Docker Registry`
   - **Image URL**: `ghcr.io/[tu-usuario-github]/backend:qa`
   - **Credentials**:
     - Username: Tu usuario de GitHub
     - Password: Personal Access Token (con scope `read:packages`)
3. **Environment Variables**:
   ```
   PORT=8080
   DB_HOST=dpg-d48vp9be5dus73cf9250-a.oregon-postgres.render.com
   DB_PORT=5432
   DB_USER=admin
   DB_PASSWORD=gImS5z9Riw8683kfI5uMY5NSNYTGXclu
   DB_NAME=final_clj4
   DB_SSLMODE=require
   ```
4. **Health Check Path**: `/healthz`
5. **Plan**: Free

### 2.2 Frontend QA

1. **Render Dashboard** → **New** → **Web Service**
2. **Configuración**:
   - **Name**: `tp8-frontend-qa`
   - **Environment**: `Docker`
   - **Registry**: `Docker Registry`
   - **Image URL**: `ghcr.io/[tu-usuario-github]/frontend:qa`
   - **Credentials**: Mismo token que backend
3. **Environment Variables**:
   ```
   PORT=8080
   NEXT_PUBLIC_API_URL=https://tp8-backend-qa.onrender.com
   ```
4. **Plan**: Free

### 2.3 Backend PROD

1. **Render Dashboard** → **New** → **Web Service** (o usar el existente)
2. **Configuración**:
   - **Name**: `tp8-backend` (o `tp8-backend-prod`)
   - **Environment**: `Docker`
   - **Registry**: `Docker Registry`
   - **Image URL**: `ghcr.io/[tu-usuario-github]/backend:prod`
   - **Credentials**: Mismo token
3. **Environment Variables**: Mismas que QA (o base de datos PROD si tienes una separada)
4. **Health Check Path**: `/healthz`

### 2.4 Frontend PROD

1. **Render Dashboard** → **New** → **Web Service** (o usar el existente)
2. **Configuración**:
   - **Name**: `tp8-frontend` (o `tp8-frontend-prod`)
   - **Environment**: `Docker`
   - **Registry**: `Docker Registry`
   - **Image URL**: `ghcr.io/[tu-usuario-github]/frontend:prod`
   - **Credentials**: Mismo token
3. **Environment Variables**:
   ```
   PORT=8080
   NEXT_PUBLIC_API_URL=https://tp8-backend.onrender.com
   ```

## 🛠️ Paso 3: Obtener Webhooks de Render

Para cada servicio en Render:

1. Ve al servicio → **Settings**
2. Busca **Manual Deploy** → **Webhook URL**
3. Copia la URL del webhook

## 🛠️ Paso 4: Configurar GitHub Secrets

Ve a tu repositorio → **Settings** → **Secrets and variables** → **Actions**

### Secrets para QA:

- `RENDER_BACKEND_WEBHOOK_URL_QA` → Webhook del backend QA
- `RENDER_FRONTEND_WEBHOOK_URL_QA` → Webhook del frontend QA
- `NEXT_PUBLIC_API_URL_QA` → `https://tp8-backend-qa.onrender.com` (opcional, tiene default)

### Secrets para PROD:

- `RENDER_BACKEND_WEBHOOK_URL_PROD` → Webhook del backend PROD
- `RENDER_FRONTEND_WEBHOOK_URL_PROD` → Webhook del frontend PROD
- `NEXT_PUBLIC_API_URL_PROD` → `https://tp8-backend.onrender.com` (opcional, tiene default)

## ✅ Paso 5: Probar

### Probar QA:

```bash
# Hacer un cambio
echo "# QA Test" >> README.md
git add .
git commit -m "Test: QA deployment"
git push origin qa
```

Verifica en GitHub Actions que se ejecute el workflow y despliegue a QA.

### Probar PROD:

```bash
# Hacer un cambio
echo "# PROD Test" >> README.md
git add .
git commit -m "Test: PROD deployment"
git push origin main
```

Verifica en GitHub Actions que se ejecute el workflow y despliegue a PROD.

## 📊 Estructura Final

```
Render Services:
├── tp8-backend-qa      → ghcr.io/user/backend:qa
├── tp8-frontend-qa     → ghcr.io/user/frontend:qa
├── tp8-backend         → ghcr.io/user/backend:prod
└── tp8-frontend        → ghcr.io/user/frontend:prod

GitHub Branches:
├── qa    → Deploy a QA
└── main  → Deploy a PROD
```

## 🔍 Verificar Despliegues

1. **GitHub Actions**: Ve a la pestaña Actions y verifica que los workflows se ejecuten correctamente
2. **Render Dashboard**: Verifica que los servicios se actualicen después del webhook
3. **URLs**:
   - QA Frontend: `https://tp8-frontend-qa.onrender.com`
   - QA Backend: `https://tp8-backend-qa.onrender.com`
   - PROD Frontend: `https://tp8-frontend.onrender.com`
   - PROD Backend: `https://tp8-backend.onrender.com`

## 🎯 Resumen de Secrets Necesarios

| Secret | Valor | Ambiente |
|--------|-------|----------|
| `RENDER_BACKEND_WEBHOOK_URL_QA` | Webhook de backend QA | QA |
| `RENDER_FRONTEND_WEBHOOK_URL_QA` | Webhook de frontend QA | QA |
| `RENDER_BACKEND_WEBHOOK_URL_PROD` | Webhook de backend PROD | PROD |
| `RENDER_FRONTEND_WEBHOOK_URL_PROD` | Webhook de frontend PROD | PROD |
| `NEXT_PUBLIC_API_URL_QA` | `https://tp8-backend-qa.onrender.com` | QA (opcional) |
| `NEXT_PUBLIC_API_URL_PROD` | `https://tp8-backend.onrender.com` | PROD (opcional) |

