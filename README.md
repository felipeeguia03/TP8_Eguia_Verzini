# TP8 - Plataforma de Cursos Online

## 📋 Descripción del Proyecto

Plataforma web para gestión de cursos online desarrollada con **Next.js** (Frontend) y **Go** (Backend), desplegada en **Render** con bases de datos **PostgreSQL** en **Railway**.

---

## 🏗️ Arquitectura del Sistema

### Stack Tecnológico

- **Frontend**: Next.js 14, React 18, TypeScript, TailwindCSS
- **Backend**: Go 1.24, Gin Framework, GORM
- **Base de Datos**: PostgreSQL 17
- **Containerización**: Docker
- **CI/CD**: GitHub Actions
- **Registry**: GitHub Container Registry (GHCR)
- **Hosting**: Render (Frontend y Backend)
- **Base de Datos**: Railway (PostgreSQL)

### Diagrama de Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                    DESARROLLADOR                             │
│              (git push origin main/qa)                       │
└────────────────────────┬──────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              GITHUB REPOSITORY                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Código Fuente (Go Backend + Next.js Frontend)      │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬──────────────────────────────────────┘
                         │
                         │ Trigger: push a main/qa
                         ▼
┌─────────────────────────────────────────────────────────────┐
│           GITHUB ACTIONS (PIPELINE CI/CD)                   │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  JOB 1: build-backend                                │  │
│  │  ├─ Checkout código                                  │  │
│  │  ├─ Setup Go 1.24                                   │  │
│  │  ├─ Run tests (go test)                             │  │
│  │  ├─ Build Docker image (backend)                    │  │
│  │  └─ Push a GHCR: ghcr.io/user/backend:qa/prod       │  │
│  │                                                      │  │
│  │  JOB 2: build-frontend                               │  │
│  │  ├─ Checkout código                                  │  │
│  │  ├─ Setup Node.js 20                                 │  │
│  │  ├─ Run tests (npm test)                            │  │
│  │  ├─ Build Docker image (frontend)                   │  │
│  │  └─ Push a GHCR: ghcr.io/user/frontend:qa/prod     │  │
│  └──────────────────────────────────────────────────────┘  │
│                         │                                    │
│                         │ Después de push exitoso            │
│                         ▼                                    │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Trigger Webhooks de Render                          │  │
│  │  ├─ curl POST → RENDER_BACKEND_WEBHOOK_URL_QA       │  │
│  │  ├─ curl POST → RENDER_FRONTEND_WEBHOOK_URL_QA       │  │
│  │  ├─ curl POST → RENDER_BACKEND_WEBHOOK_URL_PROD     │  │
│  │  └─ curl POST → RENDER_FRONTEND_WEBHOOK_URL_PROD    │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬──────────────────────────────────────┘
                         │
                         │ Webhook HTTP POST
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    RENDER (HOSTING)                          │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Servicio Backend QA (backend-qa-1)                  │  │
│  │  ├─ Pull imagen: ghcr.io/user/backend:qa             │  │
│  │  ├─ Deploy contenedor                                 │  │
│  │  └─ URL: https://backend-qa-1.onrender.com            │  │
│  │                                                       │  │
│  │  Servicio Frontend QA                                 │  │
│  │  ├─ Pull imagen: ghcr.io/user/frontend:qa             │  │
│  │  ├─ Deploy contenedor                                 │  │
│  │  └─ URL: https://frontend-qa-eop0.onrender.com        │  │
│  │                                                       │  │
│  │  Servicio Backend PROD (backend-prod-2hc0)           │  │
│  │  ├─ Pull imagen: ghcr.io/user/backend:prod           │  │
│  │  ├─ Deploy contenedor                                 │  │
│  │  └─ URL: https://backend-prod-2hc0.onrender.com      │  │
│  │                                                       │  │
│  │  Servicio Frontend PROD (frontend-prod-2czt)         │  │
│  │  ├─ Pull imagen: ghcr.io/user/frontend:prod          │  │
│  │  ├─ Deploy contenedor                                 │  │
│  │  └─ URL: https://frontend-prod-2czt.onrender.com     │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬──────────────────────────────────────┘
                         │
                         │ Conexión a BD
                         ▼
┌─────────────────────────────────────────────────────────────┐
│         POSTGRESQL (Railway Database)                         │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Base de datos QA: final_clj4 QA                    │  │
│  │  Host: turntable.proxy.rlwy.net                      │  │
│  │                                                       │  │
│  │  Base de datos PROD: final_clj4                      │  │
│  │  Host: centerbeam.proxy.rlwy.net                     │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

**Espacio para screenshot del diagrama de arquitectura:**

<!-- [Screenshot: Diagrama de arquitectura completo] -->

---

## 🔄 Flujo CI/CD Completo

### 1. Desarrollo Local

- Desarrollo en máquina local
- Testing con Docker Compose
- Commit y push a GitHub

### 2. GitHub Actions (CI)

- **Trigger**: Push a `main` (PROD) o `qa` (QA)
- **Build**: Construcción de imágenes Docker
- **Test**: Ejecución de tests automatizados
- **Push**: Subida de imágenes a GHCR con tags `:qa` o `:prod`

### 3. Render (CD)

- **Webhook**: GitHub Actions notifica a Render
- **Pull**: Render descarga la imagen desde GHCR
- **Deploy**: Render despliega el contenedor
- **Health Check**: Verificación automática de salud

**Espacio para screenshot del workflow de GitHub Actions:**

<!-- [Screenshot: GitHub Actions workflow ejecutándose] -->

**Espacio para screenshot de Render Dashboard:**

<!-- [Screenshot: Servicios desplegados en Render] -->

---

## 🐳 Dockerización

### Decisión: Multi-stage Build

**¿Por qué?**

- Reduce el tamaño final de la imagen
- Mejora la seguridad (no incluye herramientas de build)
- Acelera el despliegue

### Backend Dockerfile

```dockerfile
# Etapa 1: Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# Etapa 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
ENV PORT=8080
CMD ["./server"]
```

**Decisiones:**

- **Alpine Linux**: Imagen base ligera (~5MB)
- **CGO_ENABLED=0**: Binario estático, no requiere librerías C
- **Multi-stage**: Solo incluye el binario compilado

**Espacio para screenshot del Dockerfile:**

<!-- [Screenshot: Dockerfile del backend] -->

### Frontend Dockerfile

```dockerfile
# Etapa 1: Dependencies
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci

# Etapa 2: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
RUN npm run build

# Etapa 3: Production
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=8080
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
RUN chown -R nextjs:nodejs /app
USER nextjs
EXPOSE 8080
CMD ["node", "server.js"]
```

**Decisiones:**

- **Standalone output**: Next.js genera un servidor independiente
- **Usuario no-root**: Mejora la seguridad
- **Build args**: `NEXT_PUBLIC_API_URL` se pasa en build time

**Espacio para screenshot del Dockerfile del frontend:**

<!-- [Screenshot: Dockerfile del frontend] -->

---

## 🗄️ Base de Datos: Migración de MySQL a PostgreSQL

### Decisión: Migrar de MySQL a PostgreSQL

**¿Por qué?**

1. **Railway**: Ofrece PostgreSQL gratuito con mejor rendimiento
2. **Compatibilidad**: PostgreSQL es más estándar y compatible
3. **Herramientas**: Mejor soporte en el ecosistema Go
4. **Escalabilidad**: PostgreSQL maneja mejor cargas altas

### Cambios Realizados

#### 1. Driver de Base de Datos

```go
// Antes (MySQL)
import "gorm.io/driver/mysql"

// Después (PostgreSQL)
import "gorm.io/driver/postgres"
```

#### 2. Variables de Entorno

```bash
# Antes (MySQL)
MYSQL_HOST=...
MYSQL_PORT=3306
MYSQL_USER=...
MYSQL_PASSWORD=...
MYSQL_DB=...

# Después (PostgreSQL)
DB_HOST=...
DB_PORT=5432
DB_USER=...
DB_PASSWORD=...
DB_NAME=...
DB_SSLMODE=require
```

#### 3. Connection String

```go
// Formato PostgreSQL
dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
    dbHost, dbUser, dbPassword, dbName, dbPort, sslMode)
```

**Espacio para screenshot de la migración:**

<!-- [Screenshot: Código de migración MySQL a PostgreSQL] -->

---

## 🚀 CI/CD: GitHub Actions + Render

### Decisión: Stack 100% Gratuito

**¿Por qué?**

- **GitHub Actions**: 2000 minutos/mes gratis
- **GitHub Container Registry**: Ilimitado y gratuito
- **Render**: Plan free disponible
- **Railway**: PostgreSQL gratuito

### Workflow de GitHub Actions

**Archivo**: `.github/workflows/docker-build-push.yml`

**Características:**

- Build paralelo de backend y frontend
- Tests automatizados antes del build
- Tagging inteligente según branch (`:qa` o `:prod`)
- Webhooks automáticos a Render

**Espacio para screenshot del workflow:**

<!-- [Screenshot: Workflow de GitHub Actions] -->

### Separación de Ambientes

**Branch Strategy:**

- `main` → PROD (despliegue a producción)
- `qa` → QA (despliegue a ambiente de pruebas)

**Imágenes Docker:**

- `ghcr.io/felipeeguia03/backend:qa` → Backend QA
- `ghcr.io/felipeeguia03/backend:prod` → Backend PROD
- `ghcr.io/felipeeguia03/frontend:qa` → Frontend QA
- `ghcr.io/felipeeguia03/frontend:prod` → Frontend PROD

**Espacio para screenshot de branches:**

<!-- [Screenshot: Branches en GitHub] -->

---

## 🔐 Seguridad y Configuración

### Variables de Entorno

#### Backend

```bash
PORT=8080
DATABASE_URL=postgresql://user:password@host:port/database
# O variables individuales:
DB_HOST=...
DB_PORT=5432
DB_USER=...
DB_PASSWORD=...
DB_NAME=...
DB_SSLMODE=require
```

#### Frontend

```bash
PORT=8080
NEXT_PUBLIC_API_URL=https://backend-prod-2hc0.onrender.com
```

### CORS Configuration

**Decisión**: Whitelist de orígenes permitidos

**¿Por qué?**

- Seguridad: Solo permite requests desde dominios conocidos
- Flexibilidad: Fácil agregar nuevos orígenes
- Control: Evita ataques CSRF

**Orígenes configurados:**

- `http://localhost:3000` (desarrollo local)
- `https://frontend-prod-2czt.onrender.com` (Frontend PROD)
- `https://backend-qa-1.onrender.com` (Backend QA)
- URLs de Azure (legacy)

**Espacio para screenshot de CORS config:**

<!-- [Screenshot: Configuración de CORS en main.go] -->

---

## 📊 Estructura de Base de Datos

### Tablas Principales

1. **users**: Usuarios del sistema (estudiantes e instructores)
2. **courses**: Cursos disponibles
3. **subscriptions**: Suscripciones de usuarios a cursos
4. **comments**: Comentarios en cursos
5. **files**: Archivos subidos

**Espacio para screenshot del esquema de BD:**

<!-- [Screenshot: Diagrama ER de la base de datos] -->

**Espacio para screenshot de Railway:**

<!-- [Screenshot: Bases de datos en Railway] -->

---

## 🎯 Decisiones de Diseño

### 1. Separación de Bases de Datos

**Decisión**: Dos bases de datos separadas (QA y PROD)

**¿Por qué?**

- Aislamiento total entre ambientes
- Puedes resetear QA sin afectar PROD
- Diferentes planes/recursos si es necesario
- Más simple que usar schemas

**Espacio para screenshot:**

<!-- [Screenshot: Configuración de bases de datos en Railway] -->

### 2. Imágenes Docker en GHCR vs Build en Render

**Decisión**: Usar imágenes de GHCR

**¿Por qué?**

- GitHub Actions ya construye y prueba
- Render solo despliega (más rápido)
- Separación de responsabilidades (CI construye, CD despliega)
- Mejor control de versiones

**Espacio para screenshot:**

<!-- [Screenshot: Imágenes en GitHub Container Registry] -->

### 3. Health Checks

**Decisión**: Endpoint `/healthz` en el backend

**¿Por qué?**

- Render puede verificar que el servicio está vivo
- Útil para monitoreo
- Estándar en la industria

**Espacio para screenshot:**

<!-- [Screenshot: Health check configurado en Render] -->

---

## 🔧 Configuración de Servicios

### Backend QA

- **Nombre**: `backend-qa-1`
- **URL**: `https://backend-qa-1.onrender.com`
- **Imagen**: `ghcr.io/felipeeguia03/backend:qa`
- **Base de datos**: Railway QA
- **Plan**: Free

**Espacio para screenshot:**

<!-- [Screenshot: Configuración del backend QA en Render] -->

### Backend PROD

- **Nombre**: `backend-prod-2hc0`
- **URL**: `https://backend-prod-2hc0.onrender.com`
- **Imagen**: `ghcr.io/felipeeguia03/backend:prod`
- **Base de datos**: Railway PROD
- **Plan**: Starter/Standard (para más recursos)

**Espacio para screenshot:**

<!-- [Screenshot: Configuración del backend PROD en Render] -->

### Frontend QA

- **Nombre**: `frontend-qa-eop0`
- **URL**: `https://frontend-qa-eop0.onrender.com`
- **Imagen**: `ghcr.io/felipeeguia03/frontend:qa`
- **API URL**: `https://backend-qa-1.onrender.com`
- **Plan**: Free

**Espacio para screenshot:**

<!-- [Screenshot: Configuración del frontend QA en Render] -->

### Frontend PROD

- **Nombre**: `frontend-prod-2czt`
- **URL**: `https://frontend-prod-2czt.onrender.com`
- **Imagen**: `ghcr.io/felipeeguia03/frontend:prod`
- **API URL**: `https://backend-prod-2hc0.onrender.com`
- **Plan**: Free

**Espacio para screenshot:**

<!-- [Screenshot: Configuración del frontend PROD en Render] -->

---

## 📝 GitHub Secrets Configurados

### Secrets para Webhooks

- `RENDER_BACKEND_WEBHOOK_URL_QA`
- `RENDER_FRONTEND_WEBHOOK_URL_QA`
- `RENDER_BACKEND_WEBHOOK_URL_PROD`
- `RENDER_FRONTEND_WEBHOOK_URL_PROD`

### Secrets para Frontend Build

- `NEXT_PUBLIC_API_URL_QA` → `https://backend-qa-1.onrender.com`
- `NEXT_PUBLIC_API_URL_PROD` → `https://backend-prod-2hc0.onrender.com`

**Espacio para screenshot:**

<!-- [Screenshot: GitHub Secrets configurados] -->

---

## 🚦 Flujo de Despliegue

### Para QA

```bash
git checkout qa
# Hacer cambios
git add .
git commit -m "Cambios para QA"
git push origin qa
```

→ GitHub Actions construye imagen `:qa`
→ Render despliega automáticamente

**Espacio para screenshot:**

<!-- [Screenshot: Deploy a QA] -->

### Para PROD

```bash
git checkout main
# Hacer cambios
git add .
git commit -m "Cambios para PROD"
git push origin main
```

→ GitHub Actions construye imagen `:prod`
→ Render despliega automáticamente

**Espacio para screenshot:**

<!-- [Screenshot: Deploy a PROD] -->

---

## 🐛 Troubleshooting

### Problemas Comunes

#### 1. Error de conexión a base de datos

**Solución**: Verificar que `DATABASE_URL` use `DATABASE_PUBLIC_URL` de Railway (no la interna)

#### 2. Error de CORS

**Solución**: Agregar la URL del frontend a `allowedOrigins` en `main.go`

#### 3. Imagen no se actualiza

**Solución**: Forzar nuevo deploy en Render o verificar que GitHub Actions haya construido la imagen

#### 4. Variables de entorno no se leen

**Solución**: Verificar que estén configuradas en Render Environment y que el servicio se haya reiniciado

**Espacio para screenshot:**

<!-- [Screenshot: Logs de troubleshooting] -->

---

## 📈 Métricas y Monitoreo

### Health Checks

- Backend: `GET /healthz` → `{"status":"ok"}`
- Frontend: Verificación automática por Render

### Logs

- **Render**: Logs en tiempo real de cada servicio
- **GitHub Actions**: Logs de build y tests
- **Railway**: Logs de base de datos

**Espacio para screenshot:**

<!-- [Screenshot: Logs de Render] -->

---

## 🎓 Aprendizajes y Mejores Prácticas

### 1. Docker Multi-stage Build

- Reduce tamaño de imágenes
- Mejora seguridad
- Acelera despliegues

### 2. Separación de Ambientes

- QA para pruebas
- PROD para producción
- Bases de datos separadas

### 3. CI/CD Automatizado

- Tests antes de deploy
- Builds automáticos
- Deploys automáticos

### 4. Variables de Entorno

- Configuración por ambiente
- Secrets seguros
- Fácil cambio entre ambientes

---

## 📚 Tecnologías y Versiones

- **Go**: 1.24
- **Node.js**: 20
- **Next.js**: 14.2.3
- **PostgreSQL**: 17.6
- **Docker**: Latest
- **GitHub Actions**: v4/v5

---

## 🔗 URLs de Producción

### QA

- Frontend: `https://frontend-qa-eop0.onrender.com`
- Backend: `https://backend-qa-1.onrender.com`

### PROD

- Frontend: `https://frontend-prod-2czt.onrender.com`
- Backend: `https://backend-prod-2hc0.onrender.com`

**Espacio para screenshot:**

<!-- [Screenshot: Aplicación funcionando en producción] -->

---

## 👥 Autores

- Felipe Eguia
- [Tu compañero]

---

## 📄 Licencia

[Especificar licencia]

---

## 🙏 Agradecimientos

- Railway por el servicio gratuito de PostgreSQL
- Render por el hosting gratuito
- GitHub por GitHub Actions y GHCR

---

**Última actualización**: Noviembre 2025


scripts 
./scripts/run-qa.sh
./scripts/run-prod.sh
./scripts/deploy-qa.sh

./scripts/stop.sh