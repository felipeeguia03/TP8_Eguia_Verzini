# 🎓 Explicación Completa: Pipeline, Workflows, Docker y CI/CD

## 📚 Conceptos Básicos

### ¿Qué es un Pipeline (Tubería)?

Un **pipeline** es una secuencia automatizada de pasos que se ejecutan cuando haces cambios en tu código. Es como una línea de producción en una fábrica:

```
Código → Tests → Build → Deploy → Producción
```

### ¿Qué es un Workflow?

Un **workflow** en GitHub Actions es un archivo YAML que define qué hacer cuando ocurre un evento (como un `git push`).

---

## 🔄 TU PIPELINE COMPLETO

### Flujo Visual Simplificado

```
TÚ (Desarrollador)
    ↓
git push origin main/qa
    ↓
GitHub detecta el push
    ↓
GitHub Actions ejecuta workflow
    ↓
┌─────────────────────────────────────┐
│  WORKFLOW: docker-build-push.yml    │
│  ─────────────────────────────────  │
│  1. Descarga tu código              │
│  2. Ejecuta tests                   │
│  3. Construye imágenes Docker       │
│  4. Sube imágenes a GHCR            │
│  5. Llama webhooks de Render        │
└─────────────────────────────────────┘
    ↓
Render recibe webhook
    ↓
Render descarga imagen de GHCR
    ↓
Render despliega el contenedor
    ↓
TU APLICACIÓN ESTÁ EN PRODUCCIÓN ✅
```

---

## 📄 ARCHIVO 1: `.github/workflows/docker-build-push.yml`

### ¿Qué hace este archivo?

Este es el **workflow principal**. Se ejecuta cada vez que haces `git push`.

### Estructura del Archivo

#### 1. Cuándo se ejecuta (Triggers)

```yaml
on:
  push:
    branches:
      - main    # Cuando haces push a main → PROD
      - qa      # Cuando haces push a qa → QA
  pull_request:
    branches:
      - main
      - qa
```

**Traducción**: "Ejecuta este workflow cuando alguien hace push a `main` o `qa`"

#### 2. Variables Globales

```yaml
env:
  REGISTRY: ghcr.io
  BACKEND_IMAGE: backend
  FRONTEND_IMAGE: frontend
```

**Traducción**: "Define nombres que usaré en todo el workflow"

#### 3. Jobs (Trabajos)

El workflow tiene **2 jobs que se ejecutan en paralelo**:

---

### JOB 1: `build-backend` (Construir Backend)

#### Paso 1: Checkout (Descargar código)
```yaml
- name: Checkout code
  uses: actions/checkout@v4
```
**¿Qué hace?** Descarga tu código del repositorio a la máquina virtual de GitHub Actions.

#### Paso 2: Setup Go
```yaml
- name: Set up Go
  uses: actions/setup-go@v4
  with:
    go-version: "1.24"
```
**¿Qué hace?** Instala Go 1.24 en la máquina virtual.

#### Paso 3: Run Tests
```yaml
- name: Run tests
  working-directory: ./final/backend
  run: |
    go mod tidy
    go mod download
    go test ./... -v
```
**¿Qué hace?**
- `go mod tidy`: Limpia dependencias
- `go mod download`: Descarga dependencias
- `go test ./... -v`: Ejecuta todos los tests

**Si los tests fallan → El workflow se detiene aquí**

#### Paso 4: Login to GHCR
```yaml
- name: Login to GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
```
**¿Qué hace?** Se autentica en GitHub Container Registry para poder subir imágenes.

**Explicación**:
- `ghcr.io`: Es donde GitHub guarda las imágenes Docker
- `github.actor`: Tu usuario de GitHub
- `GITHUB_TOKEN`: Token automático que GitHub genera

#### Paso 5: Determine Environment (Determinar Ambiente)
```yaml
- name: Determine environment
  id: env
  run: |
    if [[ "${{ github.ref }}" == "refs/heads/main" ]]; then
      echo "environment=prod" >> $GITHUB_OUTPUT
      echo "tag=prod" >> $GITHUB_OUTPUT
    elif [[ "${{ github.ref }}" == "refs/heads/qa" ]]; then
      echo "environment=qa" >> $GITHUB_OUTPUT
      echo "tag=qa" >> $GITHUB_OUTPUT
    fi
```
**¿Qué hace?** Detecta en qué branch estás:
- Si es `main` → ambiente = `prod`, tag = `prod`
- Si es `qa` → ambiente = `qa`, tag = `qa`

**Resultado**: Guarda `tag=prod` o `tag=qa` para usarlo después.

#### Paso 6: Build and Push (Construir y Subir Imagen)
```yaml
- name: Build and push backend image
  uses: docker/build-push-action@v5
  with:
    context: ./final/backend
    file: ./final/backend/Dockerfile
    push: true
    tags: |
      ghcr.io/felipeeguia03/backend:prod
      ghcr.io/felipeeguia03/backend:prod-abc123
    cache-from: type=registry,ref=ghcr.io/felipeeguia03/backend:prod
    cache-to: type=inline,mode=max
```

**¿Qué hace?**
1. **Build**: Construye la imagen Docker usando `final/backend/Dockerfile`
2. **Tags**: Etiqueta la imagen con:
   - `:prod` (o `:qa`) → Versión "latest" de ese ambiente
   - `:prod-abc123` → Versión específica con el hash del commit
3. **Push**: Sube la imagen a GHCR
4. **Cache**: Usa caché para acelerar builds futuros

**Resultado**: Imagen disponible en `ghcr.io/felipeeguia03/backend:prod`

#### Paso 7: Trigger Render Deploy (Disparar Deploy en Render)

**Para PROD:**
```yaml
- name: Trigger Render Backend Deploy (PROD)
  if: steps.env.outputs.environment == 'prod'
  env:
    WEBHOOK_URL: ${{ secrets.RENDER_BACKEND_WEBHOOK_URL_PROD }}
  run: |
    curl -X POST "$WEBHOOK_URL"
```

**¿Qué hace?**
- Solo se ejecuta si `environment == 'prod'`
- Obtiene el webhook URL de GitHub Secrets
- Hace un `POST` a ese URL
- Render recibe el webhook y despliega automáticamente

**Para QA:**
```yaml
- name: Trigger Render Backend Deploy (QA)
  if: steps.env.outputs.environment == 'qa'
  env:
    WEBHOOK_URL: ${{ secrets.RENDER_BACKEND_WEBHOOK_URL_QA }}
  run: |
    curl -X POST "$WEBHOOK_URL"
```

**Mismo proceso, pero para QA**

---

### JOB 2: `build-frontend` (Construir Frontend)

Es **igual que el backend**, pero:

#### Diferencias Clave:

1. **Setup Node.js** en lugar de Go
2. **Build Args**: Pasa `NEXT_PUBLIC_API_URL` al Dockerfile
   ```yaml
   build-args: |
     NEXT_PUBLIC_API_URL=https://backend-prod-2hc0.onrender.com
   ```
   **¿Por qué?** Next.js necesita esta variable en **build time** (cuando construye), no en runtime.

3. **Tags diferentes**: `ghcr.io/felipeeguia03/frontend:prod`

---

## 📄 ARCHIVO 2: `.github/workflows/notify-render.yml`

### ¿Qué hace este archivo?

Este es un **workflow secundario** que se ejecuta **después** de que el workflow principal termine.

### Estructura

```yaml
on:
  workflow_run:
    workflows: ["Build and Push Docker Images to GHCR"]
    types:
      - completed
```

**Traducción**: "Ejecuta este workflow cuando el workflow 'Build and Push Docker Images to GHCR' termine"

### ¿Por qué existe?

**Razón**: Es un **backup** por si los webhooks en el workflow principal fallan.

**Nota**: Actualmente **NO se está usando** porque los webhooks están directamente en `docker-build-push.yml`. Podrías eliminarlo si quieres.

---

## 🐳 DOCKERFILES: Cómo se Construyen las Imágenes

### Backend Dockerfile (`final/backend/Dockerfile`)

#### Etapa 1: Builder (Construcción)
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .
```

**¿Qué hace?**
1. Usa imagen con Go instalado
2. Copia archivos de dependencias
3. Descarga dependencias
4. Copia código fuente
5. Compila el binario `server`

**Resultado**: Tienes un binario compilado

#### Etapa 2: Runtime (Ejecución)
```dockerfile
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

**¿Qué hace?**
1. Usa imagen Alpine (muy pequeña, ~5MB)
2. Copia **solo el binario** desde la etapa anterior
3. Expone puerto 8080
4. Ejecuta el binario

**¿Por qué dos etapas?**
- **Etapa 1**: Tiene Go, compiladores, etc. (pesada)
- **Etapa 2**: Solo tiene el binario (ligera)
- **Resultado**: Imagen final de ~10MB en lugar de ~300MB

**Espacio para screenshot:**
<!-- [Screenshot: Dockerfile del backend] -->

---

### Frontend Dockerfile (`final/frontend/Dockerfile`)

#### Etapa 1: Dependencies
```dockerfile
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci
```

**¿Qué hace?** Descarga todas las dependencias de Node.js

#### Etapa 2: Builder
```dockerfile
FROM node:20-alpine AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
RUN npm run build
```

**¿Qué hace?**
1. Copia dependencias de la etapa anterior
2. Copia código fuente
3. **Recibe** `NEXT_PUBLIC_API_URL` como argumento
4. Construye la aplicación Next.js

**Importante**: `NEXT_PUBLIC_API_URL` se pasa aquí porque Next.js lo necesita en **build time**.

#### Etapa 3: Runner
```dockerfile
FROM node:20-alpine AS runner
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
CMD ["node", "server.js"]
```

**¿Qué hace?**
1. Copia solo los archivos necesarios para ejecutar
2. Ejecuta el servidor Next.js

**¿Por qué 3 etapas?**
- **Etapa 1**: Solo dependencias (se puede cachear)
- **Etapa 2**: Build (necesita dependencias + código)
- **Etapa 3**: Runtime (solo archivos compilados)

**Espacio para screenshot:**
<!-- [Screenshot: Dockerfile del frontend] -->

---

## 🔗 CÓMO SE CONECTA TODO

### Flujo Completo Paso a Paso

#### 1. Tú haces cambios
```bash
git add .
git commit -m "Nuevo feature"
git push origin main
```

#### 2. GitHub detecta el push
- GitHub ve que hay un push a `main`
- Busca workflows en `.github/workflows/`
- Encuentra `docker-build-push.yml`
- **Inicia el workflow**

#### 3. GitHub Actions ejecuta el workflow

**Job 1: build-backend**
```
1. Descarga código ✅
2. Instala Go ✅
3. Ejecuta tests ✅
4. Se autentica en GHCR ✅
5. Detecta que es PROD (branch main) ✅
6. Construye imagen Docker ✅
   → Usa: final/backend/Dockerfile
   → Resultado: imagen con tu código compilado
7. Sube imagen a GHCR ✅
   → ghcr.io/felipeeguia03/backend:prod
8. Llama webhook de Render ✅
   → curl -X POST https://api.render.com/deploy/...
```

**Job 2: build-frontend** (en paralelo)
```
1. Descarga código ✅
2. Instala Node.js ✅
3. Ejecuta tests ✅
4. Se autentica en GHCR ✅
5. Detecta que es PROD ✅
6. Determina API URL ✅
   → NEXT_PUBLIC_API_URL=https://backend-prod-2hc0.onrender.com
7. Construye imagen Docker ✅
   → Usa: final/frontend/Dockerfile
   → Pasa: NEXT_PUBLIC_API_URL como build arg
   → Resultado: imagen con frontend compilado
8. Sube imagen a GHCR ✅
   → ghcr.io/felipeeguia03/frontend:prod
9. Llama webhook de Render ✅
```

#### 4. Render recibe los webhooks

**Backend:**
1. Render recibe: `POST /deploy/webhook-xxx`
2. Render sabe: "Necesito actualizar backend-prod-2hc0"
3. Render hace: `docker pull ghcr.io/felipeeguia03/backend:prod`
4. Render reinicia el contenedor con la nueva imagen
5. ✅ Backend actualizado

**Frontend:**
1. Render recibe: `POST /deploy/webhook-yyy`
2. Render sabe: "Necesito actualizar frontend-prod-2czt"
3. Render hace: `docker pull ghcr.io/felipeeguia03/frontend:prod`
4. Render reinicia el contenedor con la nueva imagen
5. ✅ Frontend actualizado

#### 5. Tu aplicación está actualizada
- Backend: Nueva versión corriendo
- Frontend: Nueva versión corriendo
- Todo automático ✅

---

## 🎯 CONCEPTOS CLAVE

### 1. ¿Qué es Docker?

**Docker** es una herramienta que empaqueta tu aplicación y todas sus dependencias en una "caja" llamada **contenedor**.

**Ventajas**:
- Funciona igual en cualquier máquina
- Incluye todo lo necesario
- Fácil de desplegar

**Ejemplo**:
```
Sin Docker:
- Necesitas instalar Go
- Necesitas instalar Node.js
- Necesitas configurar variables
- Puede fallar en otra máquina

Con Docker:
- Todo está en la imagen
- Funciona igual en todas partes
- Solo necesitas Docker instalado
```

### 2. ¿Qué es GitHub Container Registry (GHCR)?

Es donde GitHub guarda tus imágenes Docker. Es como un "almacén" de imágenes.

**URLs de tus imágenes**:
- `ghcr.io/felipeeguia03/backend:prod`
- `ghcr.io/felipeeguia03/frontend:qa`

### 3. ¿Qué es un Webhook?

Un **webhook** es una URL que Render te da. Cuando llamas a esa URL (con `curl -X POST`), Render sabe que debe hacer algo (en este caso, desplegar).

**Ejemplo**:
```
Webhook URL: https://api.render.com/deploy/srv-abc123?key=xyz

Cuando GitHub Actions hace:
curl -X POST https://api.render.com/deploy/srv-abc123?key=xyz

Render dice: "Ah, necesito actualizar el servicio srv-abc123"
```

### 4. ¿Qué es Build Time vs Runtime?

**Build Time** (Tiempo de construcción):
- Cuando se **construye** la imagen Docker
- Es cuando Next.js compila tu código
- Las variables `NEXT_PUBLIC_*` se "queman" en el código

**Runtime** (Tiempo de ejecución):
- Cuando el contenedor **está corriendo**
- Es cuando tu aplicación está funcionando
- Las variables normales se leen del ambiente

**Por eso**:
- `NEXT_PUBLIC_API_URL` se pasa en **build time** (en el Dockerfile)
- `DATABASE_URL` se lee en **runtime** (cuando el backend corre)

---

## 📊 RESUMEN VISUAL

### Tu Pipeline Completo

```
┌─────────────────────────────────────────────────────────┐
│  TÚ: git push origin main                               │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  GITHUB ACTIONS                                          │
│  ─────────────────────────────────────────────────────  │
│                                                          │
│  JOB 1: build-backend                                    │
│  ├─ Checkout código                                      │
│  ├─ Setup Go                                             │
│  ├─ Run tests                                            │
│  ├─ Login GHCR                                           │
│  ├─ Determine: tag=prod                                   │
│  ├─ Build Docker image                                   │
│  │   └─ Usa: final/backend/Dockerfile                    │
│  ├─ Push: ghcr.io/user/backend:prod                      │
│  └─ Webhook: curl POST → Render                         │
│                                                          │
│  JOB 2: build-frontend                                   │
│  ├─ Checkout código                                      │
│  ├─ Setup Node.js                                        │
│  ├─ Run tests                                            │
│  ├─ Login GHCR                                           │
│  ├─ Determine: tag=prod                                  │
│  ├─ Determine: API_URL=backend-prod-2hc0.onrender.com   │
│  ├─ Build Docker image                                   │
│  │   └─ Usa: final/frontend/Dockerfile                  │
│  │   └─ Build arg: NEXT_PUBLIC_API_URL=...              │
│  ├─ Push: ghcr.io/user/frontend:prod                     │
│  └─ Webhook: curl POST → Render                          │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ Webhooks HTTP POST
                     ▼
┌─────────────────────────────────────────────────────────┐
│  RENDER                                                  │
│  ─────────────────────────────────────────────────────  │
│                                                          │
│  Backend PROD:                                           │
│  ├─ Recibe webhook                                       │
│  ├─ docker pull ghcr.io/user/backend:prod                │
│  ├─ Reinicia contenedor                                 │
│  └─ ✅ Nuevo backend corriendo                           │
│                                                          │
│  Frontend PROD:                                          │
│  ├─ Recibe webhook                                       │
│  ├─ docker pull ghcr.io/user/frontend:prod               │
│  ├─ Reinicia contenedor                                 │
│  └─ ✅ Nuevo frontend corriendo                          │
└─────────────────────────────────────────────────────────┘
```

**Espacio para screenshot:**
<!-- [Screenshot: Diagrama completo del pipeline] -->

---

## 🔍 DETALLES TÉCNICOS

### ¿Por qué Multi-stage Build?

**Backend sin multi-stage**:
```
Imagen final: ~300MB
- Incluye Go, compiladores, herramientas
- Solo necesitas el binario
```

**Backend con multi-stage**:
```
Etapa 1 (builder): ~300MB
Etapa 2 (runtime): ~10MB ✅
- Solo el binario compilado
```

**Resultado**: Imagen 30x más pequeña, despliega más rápido.

### ¿Por qué Tags con SHA?

```yaml
tags: |
  ghcr.io/user/backend:prod
  ghcr.io/user/backend:prod-abc123
```

**Razón**:
- `:prod` → Siempre apunta a la última versión
- `:prod-abc123` → Versión específica del commit
- **Ventaja**: Puedes volver a una versión anterior si algo falla

### ¿Por qué Cache?

```yaml
cache-from: type=registry,ref=ghcr.io/user/backend:prod
cache-to: type=inline,mode=max
```

**Razón**:
- Docker puede reutilizar capas de builds anteriores
- **Ejemplo**: Si no cambias `go.mod`, no necesita descargar dependencias de nuevo
- **Resultado**: Builds más rápidos

---

## 🎓 RESUMEN EN PALABRAS SIMPLES

### ¿Qué hace tu pipeline?

1. **Tú haces push** → GitHub lo detecta
2. **GitHub Actions** → Construye imágenes Docker
3. **GitHub Actions** → Sube imágenes a GHCR
4. **GitHub Actions** → Llama a Render (webhook)
5. **Render** → Descarga la nueva imagen
6. **Render** → Reinicia el servicio
7. **✅ Tu app está actualizada**

### ¿Por qué es útil?

- **Automático**: No necesitas hacer nada manual
- **Rápido**: Todo en minutos
- **Confiable**: Tests antes de desplegar
- **Reproducible**: Mismo proceso siempre

---

## 📸 Espacios para Screenshots

**Espacio para screenshot del workflow ejecutándose:**
<!-- [Screenshot: GitHub Actions workflow en ejecución] -->

**Espacio para screenshot de imágenes en GHCR:**
<!-- [Screenshot: Imágenes en GitHub Container Registry] -->

**Espacio para screenshot de Render recibiendo webhook:**
<!-- [Screenshot: Logs de Render mostrando deploy automático] -->

---

**¿Tienes preguntas sobre alguna parte específica?** Puedo explicar más detalles de cualquier sección.

