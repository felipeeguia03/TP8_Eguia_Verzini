# 🏗️ Arquitectura CI/CD: GitHub Actions + Render

## 📊 Diagrama del Flujo Completo

```
┌─────────────────────────────────────────────────────────────────┐
│                    DESARROLLADOR                                │
│              (git push origin main)                             │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│              GITHUB REPOSITORY                                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Código Fuente (Go Backend + Next.js Frontend)          │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ Trigger: push a main
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│           GITHUB ACTIONS (PIPELINE CI/CD)                       │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  JOB 1: build-backend                                    │  │
│  │  ├─ Checkout código                                      │  │
│  │  ├─ Setup Go 1.24                                       │  │
│  │  ├─ Run tests (go test)                                  │  │
│  │  ├─ Build Docker image (backend)                         │  │
│  │  └─ Push a GHCR: ghcr.io/user/backend:latest            │  │
│  │                                                           │  │
│  │  JOB 2: build-frontend                                   │  │
│  │  ├─ Checkout código                                      │  │
│  │  ├─ Setup Node.js 20                                     │  │
│  │  ├─ Run tests (npm test)                                 │  │
│  │  ├─ Build Docker image (frontend)                        │  │
│  │  └─ Push a GHCR: ghcr.io/user/frontend:latest           │  │
│  └──────────────────────────────────────────────────────────┘  │
│                         │                                        │
│                         │ Después de push exitoso               │
│                         ▼                                        │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Trigger Webhooks de Render                              │  │
│  │  ├─ curl POST → RENDER_BACKEND_WEBHOOK_URL               │  │
│  │  └─ curl POST → RENDER_FRONTEND_WEBHOOK_URL              │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ Webhook HTTP POST
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    RENDER (HOSTING)                             │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Servicio Backend (tp8-backend)                          │  │
│  │  ├─ Recibe webhook                                       │  │
│  │  ├─ Pull imagen: ghcr.io/user/backend:latest             │  │
│  │  ├─ Reinicia contenedor con nueva imagen                 │  │
│  │  └─ Deploy en: https://tp8-backend.onrender.com          │  │
│  │                                                           │  │
│  │  Servicio Frontend (tp8-frontend)                        │  │
│  │  ├─ Recibe webhook                                       │  │
│  │  ├─ Pull imagen: ghcr.io/user/frontend:latest            │  │
│  │  ├─ Reinicia contenedor con nueva imagen                 │  │
│  │  └─ Deploy en: https://tp8-frontend.onrender.com         │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ Conexión a BD
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│         POSTGRESQL (Render Database)                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Base de datos: final_clj4                              │  │
│  │  Host: dpg-d48vp9be5dus73cf9250-a.oregon-postgres...    │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## 🔄 Flujo Paso a Paso

### 1️⃣ **CI (Continuous Integration) - GitHub Actions**

**¿Qué es?** Pipeline automatizado que se ejecuta cuando haces `git push`.

**¿Qué hace?**

- ✅ **Build**: Construye las imágenes Docker
- ✅ **Test**: Ejecuta tests (Go y Node.js)
- ✅ **Push**: Sube imágenes a GitHub Container Registry (GHCR)
- ✅ **Notify**: Notifica a Render vía webhook

**Archivo**: `.github/workflows/docker-build-push.yml`

### 2️⃣ **CD (Continuous Deployment) - Render**

**¿Qué es?** Plataforma de hosting que despliega tu aplicación.

**¿Qué hace?**

- ✅ **Pull**: Descarga la imagen Docker desde GHCR
- ✅ **Deploy**: Despliega el contenedor en producción
- ✅ **Host**: Sirve tu aplicación en URLs públicas

**Configuración**: Render Dashboard → Services

## 🎯 Roles de Cada Componente

### GitHub Actions (Pipeline CI/CD)

```
┌─────────────────────────────────────┐
│  GITHUB ACTIONS                     │
│  ─────────────────────────────────  │
│  • Construye imágenes Docker        │
│  • Ejecuta tests                     │
│  • Sube imágenes a GHCR             │
│  • Notifica a Render (webhook)      │
│                                     │
│  ⚙️  Automatización                 │
│  📦 Build & Test                    │
│  🚀 Push a Registry                 │
└─────────────────────────────────────┘
```

### GitHub Container Registry (GHCR)

```
┌─────────────────────────────────────┐
│  GHCR (ghcr.io)                     │
│  ─────────────────────────────────  │
│  • Almacena imágenes Docker         │
│  • Versiones: latest, :sha          │
│  • Privado/Gratis                   │
│                                     │
│  📦 Registry de Imágenes            │
│  🔐 Autenticación GitHub            │
└─────────────────────────────────────┘
```

### Render (Hosting/Deployment)

```
┌─────────────────────────────────────┐
│  RENDER                              │
│  ─────────────────────────────────  │
│  • Pull imagen desde GHCR            │
│  • Despliega contenedor              │
│  • Hosting público                   │
│  • Base de datos PostgreSQL          │
│                                     │
│  🌐 Hosting                          │
│  🐳 Docker Containers                │
│  💾 PostgreSQL Database              │
└─────────────────────────────────────┘
```

## 🔗 Conexión Entre Componentes

### GitHub Actions → GHCR

```yaml
# En el workflow
- name: Build and push backend image
  uses: docker/build-push-action@v5
  with:
    tags: ghcr.io/user/backend:latest
```

### GitHub Actions → Render

```yaml
# Webhook trigger
- name: Trigger Render Backend Deploy
  run: |
    curl -X POST "${{ secrets.RENDER_BACKEND_WEBHOOK_URL }}"
```

### Render → GHCR

```
Render Service Config:
  Registry: Docker Registry
  Image URL: ghcr.io/user/backend:latest
  Credentials: GitHub username + PAT token
```

### Render → PostgreSQL

```
Environment Variables en Render:
  DB_HOST: dpg-d48vp9be5dus73cf9250-a.oregon-postgres...
  DB_PORT: 5432
  DB_NAME: final_clj4
  DB_USER: admin
  DB_PASSWORD: ***
```

## 📋 Resumen

| Componente         | Rol            | Responsabilidad               |
| ------------------ | -------------- | ----------------------------- |
| **GitHub Actions** | Pipeline CI/CD | Build, Test, Push imágenes    |
| **GHCR**           | Registry       | Almacenar imágenes Docker     |
| **Render**         | Hosting        | Desplegar y servir aplicación |
| **PostgreSQL**     | Base de Datos  | Almacenar datos persistentes  |

## 🚀 Flujo Completo en Acción

1. **Desarrollador** → `git push origin main`
2. **GitHub** → Detecta push, ejecuta workflow
3. **GitHub Actions** → Build + Test + Push a GHCR
4. **GitHub Actions** → Llama webhook de Render
5. **Render** → Pull nueva imagen desde GHCR
6. **Render** → Reinicia servicio con nueva imagen
7. **Render** → Aplicación actualizada en producción ✅

## 💡 Analogía Simple

- **GitHub Actions** = Fábrica que construye y empaqueta
- **GHCR** = Almacén donde guardas los paquetes
- **Render** = Tienda donde vendes/expones tus productos
- **PostgreSQL** = Bodega donde guardas tus inventarios
