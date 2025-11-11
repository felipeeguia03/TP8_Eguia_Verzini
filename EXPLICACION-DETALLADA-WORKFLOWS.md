# 📚 Explicación Detallada de los Workflows de GitHub Actions

## 📄 Archivo 1: `.github/workflows/docker-build-push.yml`

Este es el **workflow principal** que construye y sube las imágenes Docker a GHCR, y luego notifica a Render.

---

### 🔹 SECCIÓN 1: Metadatos del Workflow

```yaml
name: Build and Push Docker Images to GHCR
```

**¿Qué hace?**

- Define el **nombre** del workflow que aparece en GitHub Actions
- Es solo para identificación, no afecta la funcionalidad

---

### 🔹 SECCIÓN 2: Triggers (Cuándo se ejecuta)

```yaml
on:
  push:
    branches:
      - main # PROD
      - master # PROD
      - qa # QA
  pull_request:
    branches:
      - main
      - master
      - qa
```

**¿Qué hace?**

- **`on: push`**: Se ejecuta cuando haces `git push` a estas ramas:
  - `main` o `master` → Despliega a PROD
  - `qa` → Despliega a QA
- **`on: pull_request`**: También se ejecuta cuando creas un Pull Request a estas ramas
  - Útil para probar antes de mergear

**Ejemplo:**

```bash
git push origin qa  # → Ejecuta el workflow para QA
git push origin main # → Ejecuta el workflow para PROD
```

---

### 🔹 SECCIÓN 3: Variables de Entorno Globales

```yaml
env:
  REGISTRY: ghcr.io
  BACKEND_IMAGE: backend
  FRONTEND_IMAGE: frontend
```

**¿Qué hace?**

- Define variables que se usan en **todos los jobs**
- **`REGISTRY`**: Dónde guardar las imágenes (GitHub Container Registry)
- **`BACKEND_IMAGE`**: Nombre base de la imagen del backend
- **`FRONTEND_IMAGE`**: Nombre base de la imagen del frontend

**¿Por qué?**

- Evita repetir `ghcr.io` en todos lados
- Si cambias de registry, solo cambias aquí
- Más fácil de mantener

**Resultado:**

- Backend: `ghcr.io/felipeeguia03/backend:qa`
- Frontend: `ghcr.io/felipeeguia03/frontend:qa`

---

### 🔹 SECCIÓN 4: Job 1 - Build Backend

#### 4.1. Configuración del Job

```yaml
build-backend:
  name: Build and Push Backend Image
  runs-on: ubuntu-latest
  permissions:
    contents: read
    packages: write
```

**¿Qué hace?**

- **`name`**: Nombre que aparece en los logs
- **`runs-on: ubuntu-latest`**: Ejecuta en una máquina virtual Ubuntu
- **`permissions`**:
  - `contents: read`: Puede leer el código del repo
  - `packages: write`: Puede escribir (subir) imágenes a GHCR

---

#### 4.2. Step 1: Checkout Code

```yaml
- name: Checkout code
  uses: actions/checkout@v4
```

**¿Qué hace?**

- Descarga el código del repositorio a la máquina virtual
- Sin esto, no tendrías acceso a tus archivos

**Equivalente a:**

```bash
git clone https://github.com/felipeeguia03/TP8_Eguia_Verzini.git
```

---

#### 4.3. Step 2: Set up Go

```yaml
- name: Set up Go
  uses: actions/setup-go@v4
  with:
    go-version: "1.24"
```

**¿Qué hace?**

- Instala Go versión 1.24 en la máquina virtual
- Necesario para compilar y testear el backend

**Equivalente a:**

```bash
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

---

#### 4.4. Step 3: Run Tests

```yaml
- name: Run tests
  working-directory: ./final/backend
  run: |
    go mod tidy
    go mod download
    go test ./... -v
```

**¿Qué hace?**

- **`working-directory`**: Cambia al directorio del backend
- **`go mod tidy`**: Limpia y actualiza `go.mod`
- **`go mod download`**: Descarga dependencias
- **`go test ./... -v`**: Ejecuta todos los tests

**¿Por qué?**

- Si los tests fallan, el workflow se detiene
- Evita construir imágenes con código roto

---

#### 4.5. Step 4: Login to GHCR

```yaml
- name: Login to GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ${{ env.REGISTRY }}
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
```

**¿Qué hace?**

- Se autentica en GitHub Container Registry (GHCR)
- **`${{ env.REGISTRY }}`**: `ghcr.io` (de la sección `env`)
- **`${{ github.actor }}`**: Tu usuario de GitHub (ej: `felipeeguia03`)
- **`${{ secrets.GITHUB_TOKEN }}`**: Token automático de GitHub (no necesitas crearlo)

**Equivalente a:**

```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u felipeeguia03 --password-stdin
```

**¿Por qué?**

- Sin autenticación, no puedes subir imágenes
- GitHub genera el token automáticamente

---

#### 4.6. Step 5: Determine Environment

```yaml
- name: Determine environment
  id: env
  run: |
    if [[ "${{ github.ref }}" == "refs/heads/main" ]] || [[ "${{ github.ref }}" == "refs/heads/master" ]]; then
      echo "environment=prod" >> $GITHUB_OUTPUT
      echo "tag=prod" >> $GITHUB_OUTPUT
    elif [[ "${{ github.ref }}" == "refs/heads/qa" ]]; then
      echo "environment=qa" >> $GITHUB_OUTPUT
      echo "tag=qa" >> $GITHUB_OUTPUT
    else
      echo "environment=dev" >> $GITHUB_OUTPUT
      echo "tag=dev" >> $GITHUB_OUTPUT
    fi
```

**¿Qué hace?**

- Detecta en qué ambiente estás (QA o PROD)
- **`${{ github.ref }}`**: La rama que activó el workflow
  - `refs/heads/main` → PROD
  - `refs/heads/qa` → QA
- **`id: env`**: Guarda los resultados para usar después
- **`$GITHUB_OUTPUT`**: Archivo donde GitHub Actions guarda variables

**Resultado:**

- Si pusheas a `qa`: `environment=qa`, `tag=qa`
- Si pusheas a `main`: `environment=prod`, `tag=prod`

**Uso después:**

- `${{ steps.env.outputs.environment }}` → `qa` o `prod`
- `${{ steps.env.outputs.tag }}` → `qa` o `prod`

---

#### 4.7. Step 6: Build and Push Backend Image

```yaml
- name: Build and push backend image
  uses: docker/build-push-action@v5
  with:
    context: ./final/backend
    file: ./final/backend/Dockerfile
    push: true
    tags: |
      ${{ env.REGISTRY }}/${{ github.repository_owner }}/${{ env.BACKEND_IMAGE }}:${{ steps.env.outputs.tag }}
      ${{ env.REGISTRY }}/${{ github.repository_owner }}/${{ env.BACKEND_IMAGE }}:${{ steps.env.outputs.tag }}-${{ github.sha }}
    cache-from: type=registry,ref=${{ env.REGISTRY }}/${{ github.repository_owner }}/${{ env.BACKEND_IMAGE }}:${{ steps.env.outputs.tag }}
    cache-to: type=inline,mode=max
```

**¿Qué hace cada parte?**

**`context: ./final/backend`**

- Directorio base para construir la imagen
- Docker busca archivos desde aquí

**`file: ./final/backend/Dockerfile`**

- Qué Dockerfile usar
- Instrucciones para construir la imagen

**`push: true`**

- Sube la imagen a GHCR después de construirla
- Si es `false`, solo construye localmente

**`tags:`**

- Etiquetas (nombres) para la imagen
- **Tag 1**: `ghcr.io/felipeeguia03/backend:qa`
  - Tag fijo, siempre apunta a la última versión
- **Tag 2**: `ghcr.io/felipeeguia03/backend:qa-abc123def`
  - Tag con SHA del commit, versión específica

**`cache-from:`**

- Usa caché de la imagen anterior para acelerar builds
- Si no cambió nada, reutiliza capas

**`cache-to: type=inline,mode=max`**

- Guarda caché dentro de la imagen para próximos builds

**Equivalente a:**

```bash
docker build -t ghcr.io/felipeeguia03/backend:qa ./final/backend
docker push ghcr.io/felipeeguia03/backend:qa
```

---

#### 4.8. Step 7: Trigger Render Backend Deploy (PROD)

```yaml
- name: Trigger Render Backend Deploy (PROD)
  if: steps.env.outputs.environment == 'prod'
  continue-on-error: true
  env:
    WEBHOOK_URL: ${{ secrets.RENDER_BACKEND_WEBHOOK_URL_PROD }}
  run: |
    if [ -n "$WEBHOOK_URL" ]; then
      curl -X POST "$WEBHOOK_URL" || echo "Webhook failed but continuing"
      echo "✅ Triggered PROD backend deploy"
    else
      echo "⚠️ RENDER_BACKEND_WEBHOOK_URL_PROD not configured, skipping"
    fi
```

**¿Qué hace cada parte?**

**`if: steps.env.outputs.environment == 'prod'`**

- Solo se ejecuta si el ambiente es PROD
- Si es QA, salta este step

**`continue-on-error: true`**

- Si el webhook falla, el workflow continúa
- No quieres que un error de Render rompa todo

**`env: WEBHOOK_URL`**

- Define variable de entorno con el webhook de Render
- Viene de GitHub Secrets

**`if [ -n "$WEBHOOK_URL" ]`**

- Verifica que el webhook esté configurado
- Si no está, muestra warning pero no falla

**`curl -X POST "$WEBHOOK_URL"`**

- Llama al webhook de Render
- Render recibe la notificación y despliega

**Equivalente a:**

```bash
curl -X POST https://api.render.com/deploy/srv-xxxxx?key=yyyyy
```

---

#### 4.9. Step 8: Trigger Render Backend Deploy (QA)

```yaml
- name: Trigger Render Backend Deploy (QA)
  if: steps.env.outputs.environment == 'qa'
  continue-on-error: true
  env:
    WEBHOOK_URL: ${{ secrets.RENDER_BACKEND_WEBHOOK_URL_QA }}
  run: |
    if [ -n "$WEBHOOK_URL" ]; then
      curl -X POST "$WEBHOOK_URL" || echo "Webhook failed but continuing"
      echo "✅ Triggered QA backend deploy"
    else
      echo "⚠️ RENDER_BACKEND_WEBHOOK_URL_QA not configured, skipping"
    fi
```

**¿Qué hace?**

- Igual que el step anterior, pero para QA
- Solo se ejecuta si `environment == 'qa'`

---

### 🔹 SECCIÓN 5: Job 2 - Build Frontend

Este job es **similar al backend**, pero con diferencias clave:

#### 5.1. Step 1-3: Setup (igual que backend)

```yaml
- name: Checkout code
  uses: actions/checkout@v4

- name: Set up Node.js
  uses: actions/setup-node@v4
  with:
    node-version: "20"

- name: Install dependencies and test
  working-directory: ./final/frontend
  run: |
    npm ci
    npm test
```

**Diferencias:**

- Usa Node.js en vez de Go
- `npm ci` en vez de `go mod download`
- `npm test` en vez de `go test`

---

#### 5.2. Step 4: Determine API URL (NUEVO)

```yaml
- name: Determine API URL
  id: api_url
  run: |
    if [[ "${{ steps.env.outputs.environment }}" == "prod" ]]; then
      echo "url=${{ secrets.NEXT_PUBLIC_API_URL_PROD || 'https://tp8-backend.onrender.com' }}" >> $GITHUB_OUTPUT
    elif [[ "${{ steps.env.outputs.environment }}" == "qa" ]]; then
      echo "url=${{ secrets.NEXT_PUBLIC_API_URL_QA || 'https://tp8-backend-qa.onrender.com' }}" >> $GITHUB_OUTPUT
    else
      echo "url=http://localhost:8080" >> $GITHUB_OUTPUT
    fi
```

**¿Qué hace?**

- Determina qué URL del backend usar
- **PROD**: `https://backend-prod-2hc0.onrender.com`
- **QA**: `https://backend-qa-1.onrender.com`
- **Dev**: `http://localhost:8080`

**¿Por qué?**

- El frontend necesita saber a qué backend conectarse
- Se inyecta durante el build de Docker (no en runtime)

**`||` significa:**

- Si el secret no existe, usa el valor por defecto

---

#### 5.3. Step 5: Build and Push Frontend Image

```yaml
- name: Build and push frontend image
  uses: docker/build-push-action@v5
  with:
    context: ./final/frontend
    file: ./final/frontend/Dockerfile
    push: true
    tags: |
      ${{ env.REGISTRY }}/${{ github.repository_owner }}/${{ env.FRONTEND_IMAGE }}:${{ steps.env.outputs.tag }}
      ${{ env.REGISTRY }}/${{ github.repository_owner }}/${{ env.FRONTEND_IMAGE }}:${{ steps.env.outputs.tag }}-${{ github.sha }}
    build-args: |
      NEXT_PUBLIC_API_URL=${{ steps.api_url.outputs.url }}
    cache-from: type=registry,ref=${{ env.REGISTRY }}/${{ github.repository_owner }}/${{ env.FRONTEND_IMAGE }}:${{ steps.env.outputs.tag }}
    cache-to: type=inline,mode=max
```

**Diferencias con backend:**

**`build-args:`**

- Pasa variables al Dockerfile durante el build
- **`NEXT_PUBLIC_API_URL`**: URL del backend
- Next.js necesita esto en **build time** (no runtime)

**En el Dockerfile del frontend:**

```dockerfile
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
```

**Resultado:**

- Frontend QA se conecta a `backend-qa-1.onrender.com`
- Frontend PROD se conecta a `backend-prod-2hc0.onrender.com`

---

#### 5.4. Steps 6-7: Trigger Render Frontend Deploy

Igual que los steps del backend, pero para frontend:

- Step 6: PROD
- Step 7: QA

---

## 📄 Archivo 2: `.github/workflows/notify-render.yml`

Este archivo es un **workflow secundario** que se ejecuta **después** del principal.

---

### 🔹 SECCIÓN 1: Metadatos

```yaml
name: Notify Render on Image Update
```

**¿Qué hace?**

- Nombre del workflow secundario

---

### 🔹 SECCIÓN 2: Trigger (Cuándo se ejecuta)

```yaml
on:
  workflow_run:
    workflows: ["Build and Push Docker Images to GHCR"]
    types:
      - completed
    branches:
      - main
      - master
```

**¿Qué hace?**

- **`workflow_run`**: Se ejecuta cuando **otro workflow** termina
- **`workflows: ["Build and Push Docker Images to GHCR"]`**: Espera a que termine el workflow principal
- **`types: completed`**: Solo cuando el workflow principal **termina** (exitoso o fallido)
- **`branches: main, master`**: Solo para PROD (no QA)

**Flujo:**

```
1. docker-build-push.yml se ejecuta
2. docker-build-push.yml termina (completed)
3. notify-render.yml se ejecuta automáticamente
```

---

### 🔹 SECCIÓN 3: Job: Notify Render

```yaml
jobs:
  notify-render:
    name: Trigger Render Deploy
    runs-on: ubuntu-latest
    if: ${{ github.event.workflow_run.conclusion == 'success' }}
    permissions:
      contents: read
```

**¿Qué hace cada parte?**

**`if: github.event.workflow_run.conclusion == 'success'`**

- Solo se ejecuta si el workflow principal **tuvo éxito**
- Si falló, no notifica a Render

**`permissions: contents: read`**

- Solo necesita leer el repo (no escribir)

---

### 🔹 SECCIÓN 4: Step 1: Trigger Backend Deploy

```yaml
- name: Trigger Backend Deploy
  continue-on-error: true
  env:
    WEBHOOK_URL: ${{ secrets.RENDER_BACKEND_WEBHOOK_URL }}
  run: |
    if [ -n "$WEBHOOK_URL" ] && [ "$WEBHOOK_URL" != "" ]; then
      curl -X POST "$WEBHOOK_URL" || true
      echo "Backend webhook triggered"
    else
      echo "RENDER_BACKEND_WEBHOOK_URL not configured, skipping"
    fi
```

**¿Qué hace?**

- Llama al webhook del backend
- **Nota**: Usa `RENDER_BACKEND_WEBHOOK_URL` (sin `_PROD` o `_QA`)
- Este archivo es más simple, solo para PROD

**¿Por qué existe este archivo?**

- Era un método alternativo de notificación
- **Ahora está obsoleto** porque `docker-build-push.yml` ya llama los webhooks directamente

---

### 🔹 SECCIÓN 5: Step 2: Trigger Frontend Deploy

```yaml
- name: Trigger Frontend Deploy
  continue-on-error: true
  env:
    WEBHOOK_URL: ${{ secrets.RENDER_FRONTEND_WEBHOOK_URL }}
  run: |
    if [ -n "$WEBHOOK_URL" ] && [ "$WEBHOOK_URL" != "" ]; then
      curl -X POST "$WEBHOOK_URL" || true
      echo "Frontend webhook triggered"
    else
      echo "RENDER_FRONTEND_WEBHOOK_URL not configured, skipping"
    fi
```

**¿Qué hace?**

- Igual que el step anterior, pero para frontend

---

## 🔄 Comparación: ¿Cuál usar?

### `docker-build-push.yml` (Principal)

✅ **Ventajas:**

- Llama webhooks **inmediatamente** después de construir
- Soporta QA y PROD
- Más directo y eficiente

❌ **Desventajas:**

- Ninguna

### `notify-render.yml` (Secundario)

✅ **Ventajas:**

- Se ejecuta solo si el workflow principal tuvo éxito
- Separación de responsabilidades

❌ **Desventajas:**

- Solo para PROD (no QA)
- Duplica funcionalidad
- Más lento (espera a que termine el workflow principal)

---

## 🎯 Recomendación

**Usa solo `docker-build-push.yml`** porque:

1. Ya llama los webhooks directamente
2. Soporta QA y PROD
3. Es más rápido
4. `notify-render.yml` es redundante

**Puedes eliminar `notify-render.yml`** si quieres simplificar.

---

## 📊 Flujo Completo Visual

```
1. git push origin qa
   ↓
2. GitHub detecta push
   ↓
3. docker-build-push.yml se ejecuta
   ├─ Job 1: build-backend
   │  ├─ Checkout code
   │  ├─ Setup Go
   │  ├─ Run tests
   │  ├─ Login GHCR
   │  ├─ Determine: environment=qa
   │  ├─ Build: backend:qa
   │  ├─ Push: ghcr.io/user/backend:qa
   │  └─ Webhook: curl POST → Render QA ✅
   │
   └─ Job 2: build-frontend (en paralelo)
      ├─ Checkout code
      ├─ Setup Node.js
      ├─ Run tests
      ├─ Login GHCR
      ├─ Determine: environment=qa
      ├─ Determine: API_URL=backend-qa-1.onrender.com
      ├─ Build: frontend:qa (con API_URL)
      ├─ Push: ghcr.io/user/frontend:qa
      └─ Webhook: curl POST → Render QA ✅
   ↓
4. Render recibe webhooks
   ├─ Backend: docker pull → Reinicia ✅
   └─ Frontend: docker pull → Reinicia ✅
   ↓
5. ✅ Aplicación actualizada
```

---

## 🔑 Conceptos Clave

### 1. Variables de GitHub Actions

**`${{ github.ref }}`**

- La rama que activó el workflow
- Ej: `refs/heads/qa`

**`${{ github.actor }}`**

- Usuario que hizo el push
- Ej: `felipeeguia03`

**`${{ github.sha }}`**

- Hash del commit
- Ej: `abc123def456`

**`${{ secrets.NOMBRE }}`**

- Secretos configurados en GitHub
- No se muestran en logs

**`${{ steps.ID.outputs.VARIABLE }}`**

- Variables de un step anterior
- Ej: `${{ steps.env.outputs.tag }}`

### 2. Permisos

**`contents: read`**

- Puede leer el código

**`packages: write`**

- Puede subir imágenes a GHCR

### 3. Condiciones

**`if: condition`**

- Solo ejecuta si la condición es verdadera

**`continue-on-error: true`**

- Si falla, continúa el workflow

---

## ❓ Preguntas Frecuentes

### ¿Por qué dos tags por imagen?

- **Tag fijo** (`backend:qa`): Siempre apunta a la última versión
- **Tag con SHA** (`backend:qa-abc123`): Versión específica para rollback

### ¿Qué pasa si el webhook falla?

- El workflow continúa (por `continue-on-error: true`)
- Puedes hacer deploy manual en Render

### ¿Por qué `notify-render.yml` solo para PROD?

- Fue diseñado antes de tener QA
- Ahora es redundante

### ¿Puedo eliminar `notify-render.yml`?

- Sí, `docker-build-push.yml` ya hace todo

---

¡Eso es todo! ¿Alguna parte específica que quieras que profundice más?
