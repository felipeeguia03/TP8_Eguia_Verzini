# ✅ Cómo Verificar si tu Workflow está Funcionando

Esta guía te muestra cómo verificar si tu repositorio tiene el workflow configurado y si está funcionando correctamente.

## 🔍 Paso 1: Verificar que el Workflow Existe

### Opción A: Desde el Repositorio en GitHub

1. Ve a tu repositorio en GitHub
2. Click en la pestaña **Actions** (arriba del repositorio, junto a Code, Issues, etc.)
3. Deberías ver:
   - Si el workflow existe: Verás "Build and Push Docker Images to GHCR" en la lista
   - Si no existe: Verás un mensaje como "Get started with GitHub Actions"

### Opción B: Verificar el Archivo Directamente

1. Ve a tu repositorio en GitHub
2. Busca la carpeta `.github/workflows/`
3. Deberías ver el archivo: `docker-build-push.yml`

**Si no ves la carpeta `.github`**:

- Puede estar oculta (archivos que empiezan con punto)
- Click en "Show hidden files" o busca directamente: `https://github.com/[tu-usuario]/[tu-repo]/tree/main/.github/workflows`

---

## 🚀 Paso 2: Verificar que el Workflow se ha Ejecutado

### Ver Ejecuciones del Workflow

1. Ve a tu repositorio → **Actions** (pestaña superior)
2. En el menú lateral izquierdo, busca **"Build and Push Docker Images to GHCR"**
3. Click en el workflow
4. Verás una lista de ejecuciones (runs)

### ¿Qué Buscar?

#### ✅ Workflow Funcionando Correctamente

Deberías ver:

- **Estado verde** (✓) con checkmark
- **Título**: "Build and Push Docker Images to GHCR"
- **Última ejecución**: Fecha reciente
- **Jobs**:
  - `build-backend` (verde)
  - `build-frontend` (verde)

#### ⚠️ Workflow con Problemas

Si ves:

- **Estado rojo** (✗) o amarillo (⚠)
- **Error** en algún job
- **Failed** en el estado

---

## 📊 Paso 3: Ver Detalles de una Ejecución

### Ver Logs de una Ejecución

1. En la página de **Actions**, click en una ejecución específica
2. Verás los jobs: `build-backend` y `build-frontend`
3. Click en un job para ver los logs detallados

### ¿Qué Deberías Ver en los Logs?

#### Job `build-backend`:

```
✓ Checkout code
✓ Set up Go
✓ Run tests
✓ Login to GitHub Container Registry
✓ Build and push backend image
```

#### Job `build-frontend`:

```
✓ Checkout code
✓ Set up Node.js
✓ Install dependencies and test
✓ Login to GitHub Container Registry
✓ Build and push frontend image
```

### Errores Comunes

Si ves errores, revisa:

- **"unauthorized"**: Problema con `GITHUB_TOKEN` (raro, debería funcionar automáticamente)
- **"Dockerfile not found"**: El path del Dockerfile es incorrecto
- **"test failed"**: Los tests están fallando
- **"push failed"**: Problema con permisos de GHCR

---

## 📦 Paso 4: Verificar que las Imágenes se Subieron a GHCR

### Método 1: Desde GitHub (Más Fácil)

1. Ve a tu repositorio en GitHub
2. En la barra lateral derecha, busca **"Packages"** (debajo de "About")
3. Click en **Packages**
4. Deberías ver:
   - `backend` (package)
   - `frontend` (package)

**Si no ves "Packages"**:

- Puede estar en la parte inferior de la barra lateral
- O busca directamente: `https://github.com/[tu-usuario]?tab=packages`

### Método 2: Ver Imágenes Específicas

1. Click en el package `backend` o `frontend`
2. Verás:
   - **Versions**: Lista de tags (latest, SHA del commit, etc.)
   - **Last updated**: Fecha de última actualización
   - **Visibility**: Público o Privado

### Método 3: Desde la Línea de Comandos

```bash
# Si tienes GitHub CLI instalado
gh auth login
gh api user/packages?package_type=container
```

---

## 🎯 Paso 5: Verificar que el Workflow se Ejecuta Automáticamente

### Probar el Workflow

1. Haz un cambio pequeño en tu código (ej: agregar un comentario)
2. Haz commit y push:
   ```bash
   git add .
   git commit -m "test workflow"
   git push origin main
   ```
3. Ve a **Actions** en GitHub
4. Deberías ver una nueva ejecución iniciándose automáticamente

### Verificar el Trigger

El workflow se ejecuta cuando:

- ✅ Haces push a `main` o `master`
- ✅ Haces un Pull Request a `main` o `master`
- ✅ Haces push manual (si está configurado)

---

## 📋 Checklist de Verificación

Usa este checklist para verificar que todo está funcionando:

### ✅ Workflow Configurado

- [ ] Existe el archivo `.github/workflows/docker-build-push.yml`
- [ ] El workflow aparece en la pestaña **Actions**

### ✅ Workflow Ejecutado

- [ ] Hay al menos una ejecución en **Actions**
- [ ] La última ejecución tiene estado verde (✓)
- [ ] Ambos jobs (`build-backend` y `build-frontend`) completaron exitosamente

### ✅ Imágenes en GHCR

- [ ] El package `backend` existe en GitHub Packages
- [ ] El package `frontend` existe en GitHub Packages
- [ ] Ambos tienen el tag `latest`
- [ ] Ambos tienen tags con SHA del commit

### ✅ Auto-ejecución

- [ ] Al hacer push, el workflow se ejecuta automáticamente
- [ ] Los logs muestran que las imágenes se construyeron y subieron correctamente

---

## 🐛 Problemas Comunes y Soluciones

### Problema: "No veo la pestaña Actions"

**Causa**: El repositorio no tiene workflows configurados o no tienes permisos

**Solución**:

- Verifica que exista `.github/workflows/docker-build-push.yml`
- Asegúrate de tener permisos de escritura en el repositorio

---

### Problema: "El workflow no se ejecuta"

**Causa**: Puede ser que no hayas hecho push a `main` o `master`

**Solución**:

1. Verifica en qué rama estás: `git branch`
2. Haz push a `main` o `master`:
   ```bash
   git push origin main
   ```
3. O verifica el trigger en el workflow (debe incluir `main` y `master`)

---

### Problema: "Workflow falla con error de permisos"

**Causa**: Falta el permiso `packages: write` en el workflow

**Solución**:
Verifica que el workflow tenga:

```yaml
permissions:
  contents: read
  packages: write
```

---

### Problema: "No veo los packages en GitHub"

**Causa**: Las imágenes no se subieron o el repositorio es privado

**Solución**:

1. Revisa los logs del workflow para ver si hubo errores en el push
2. Si el repositorio es privado, los packages también serán privados
3. Espera unos minutos (puede haber un delay)

---

### Problema: "Los tests fallan"

**Causa**: Los tests en tu código están fallando

**Solución**:

1. Revisa los logs del job que falla
2. Ejecuta los tests localmente:

   ```bash
   # Backend
   cd final/backend
   go test ./...

   # Frontend
   cd final/frontend
   npm test
   ```

3. Corrige los tests y haz push nuevamente

---

## 🔗 URLs Útiles

### Ver Workflows

```
https://github.com/[tu-usuario]/[tu-repositorio]/actions
```

### Ver Packages

```
https://github.com/[tu-usuario]?tab=packages
```

### Ver una Ejecución Específica

```
https://github.com/[tu-usuario]/[tu-repositorio]/actions/runs/[run-id]
```

---

## 📸 Qué Deberías Ver (Ejemplo Visual)

### En la Pestaña Actions:

```
Actions
├── Build and Push Docker Images to GHCR
    ├── ✓ Run #5 (hace 2 horas) - main
    ├── ✓ Run #4 (hace 1 día) - main
    └── ✓ Run #3 (hace 3 días) - main
```

### En Packages:

```
Packages
├── backend
│   ├── Latest: latest
│   └── Updated: hace 2 horas
└── frontend
    ├── Latest: latest
    └── Updated: hace 2 horas
```

---

## ✅ Resumen Rápido

**Para verificar rápidamente si tu workflow funciona**:

1. Ve a: `https://github.com/[tu-usuario]/[tu-repositorio]/actions`
2. Deberías ver ejecuciones con estado verde (✓)
3. Ve a: `https://github.com/[tu-usuario]?tab=packages`
4. Deberías ver los packages `backend` y `frontend`

**Si todo esto está correcto, ¡tu workflow está funcionando!** 🎉

---

## 🚀 Siguiente Paso

Una vez que verifiques que el workflow funciona:

- ✅ Configura Render (ver [GUIA-CONFIGURAR-RENDER.md](./GUIA-CONFIGURAR-RENDER.md))
- ✅ Agrega los webhooks de Render a GitHub Secrets
- ✅ ¡Disfruta del auto-deploy! 🎊
