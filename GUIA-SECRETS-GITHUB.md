# 🔐 Guía Paso a Paso: Configurar Secrets en GitHub

Esta guía te muestra exactamente dónde y cómo configurar los secrets necesarios para que GitHub Actions funcione correctamente.

## 📍 Paso 1: Ir a la Configuración de Secrets

### Opción A: Desde el Repositorio

1. Ve a tu repositorio en GitHub (ej: `https://github.com/tu-usuario/tu-repositorio`)
2. En la parte superior del repositorio, haz click en **Settings** (Configuración)
   - Si no ves "Settings", asegúrate de que tengas permisos de administrador en el repositorio
3. En el menú lateral izquierdo, busca la sección **Secrets and variables**
4. Haz click en **Actions**

### Opción B: URL Directa

Puedes ir directamente a:

```
https://github.com/[tu-usuario]/[tu-repositorio]/settings/secrets/actions
```

Reemplaza:

- `[tu-usuario]` con tu usuario de GitHub
- `[tu-repositorio]` con el nombre de tu repositorio

## 🔑 Paso 2: Agregar Secrets

Una vez que estés en la página de Secrets, verás:

### Pantalla Inicial

- Si no tienes secrets, verás un mensaje como "There are no secrets"
- Si ya tienes secrets, verás una lista
- En la parte superior derecha, verás un botón **New repository secret**

### Agregar cada Secret

Haz click en **New repository secret** para cada uno:

---

### Secret 1: `NEXT_PUBLIC_API_URL` (Opcional pero Recomendado)

**Cuándo agregarlo**: Si ya sabes la URL de tu backend en Render

**Pasos**:

1. Click en **New repository secret**
2. **Name**: Escribe exactamente: `NEXT_PUBLIC_API_URL`
3. **Secret**: Escribe la URL de tu backend
   - Ejemplo: `https://tu-backend.onrender.com`
   - O si aún no lo tienes: `http://localhost:8080` (puedes cambiarlo después)
4. Click en **Add secret**

**Nota**: Si no lo agregas, el workflow usará `http://localhost:8080` por defecto.

---

### Secret 2: `RENDER_BACKEND_WEBHOOK_URL` (Opcional)

**Cuándo agregarlo**: Después de crear el servicio backend en Render

**Pasos**:

1. Primero, crea el servicio backend en Render (ver guía de Render)
2. En Render, ve a tu servicio backend → **Settings** → **Manual Deploy**
3. Busca **Webhook URL** o **Deploy Hook**
4. Copia la URL completa (algo como: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`)
5. En GitHub, click en **New repository secret**
6. **Name**: Escribe exactamente: `RENDER_BACKEND_WEBHOOK_URL`
7. **Secret**: Pega la URL que copiaste
8. Click en **Add secret**

**Nota**: Si no lo agregas, el auto-deploy no funcionará automáticamente (tendrás que hacer deploy manual en Render).

---

### Secret 3: `RENDER_FRONTEND_WEBHOOK_URL` (Opcional)

**Cuándo agregarlo**: Después de crear el servicio frontend en Render

**Pasos**:

1. Primero, crea el servicio frontend en Render (ver guía de Render)
2. En Render, ve a tu servicio frontend → **Settings** → **Manual Deploy**
3. Busca **Webhook URL** o **Deploy Hook**
4. Copia la URL completa
5. En GitHub, click en **New repository secret**
6. **Name**: Escribe exactamente: `RENDER_FRONTEND_WEBHOOK_URL`
7. **Secret**: Pega la URL que copiaste
8. Click en **Add secret**

---

## ✅ Verificar que los Secrets Están Configurados

Después de agregar los secrets, deberías ver una lista como esta:

```
Secrets (3)
├── NEXT_PUBLIC_API_URL          [Update] [Delete]
├── RENDER_BACKEND_WEBHOOK_URL    [Update] [Delete]
└── RENDER_FRONTEND_WEBHOOK_URL   [Update] [Delete]
```

**Importante**:

- Los valores de los secrets están ocultos (solo ves `••••••••`)
- Puedes actualizarlos o eliminarlos en cualquier momento
- Una vez que guardas un secret, **no puedes ver su valor** (solo cambiarlo)

---

## 🚫 Lo que NO Necesitas Hacer

### ❌ NO necesitas crear `GITHUB_TOKEN`

- GitHub lo proporciona automáticamente
- Ya tiene los permisos necesarios para GHCR
- No aparece en la lista de secrets (es automático)

### ❌ NO necesitas crear secrets para ACR

- Ya no usamos Azure Container Registry
- Todo funciona con GHCR que usa `GITHUB_TOKEN` automáticamente

---

## 📝 Resumen Rápido

| Secret                        | Requerido | Cuándo Agregarlo          | Dónde Obtenerlo                 |
| ----------------------------- | --------- | ------------------------- | ------------------------------- |
| `NEXT_PUBLIC_API_URL`         | Opcional  | Antes o después           | URL de tu backend en Render     |
| `RENDER_BACKEND_WEBHOOK_URL`  | Opcional  | Después de crear backend  | Settings del servicio en Render |
| `RENDER_FRONTEND_WEBHOOK_URL` | Opcional  | Después de crear frontend | Settings del servicio en Render |

---

## 🎯 Orden Recomendado

1. **Primero**: Agrega `NEXT_PUBLIC_API_URL` (puedes usar `http://localhost:8080` temporalmente)
2. **Después**: Crea los servicios en Render
3. **Finalmente**: Agrega los webhooks de Render

---

## 🔍 Dónde Ver los Secrets en GitHub

**Ruta completa**:

```
Repositorio → Settings → Secrets and variables → Actions
```

**URL directa**:

```
https://github.com/[usuario]/[repositorio]/settings/secrets/actions
```

---

## ❓ Preguntas Frecuentes

### ¿Puedo ver el valor de un secret después de guardarlo?

No, GitHub no muestra los valores por seguridad. Solo puedes actualizarlos o eliminarlos.

### ¿Qué pasa si escribo mal el nombre del secret?

El workflow fallará porque no encontrará el secret. Asegúrate de escribir el nombre exactamente como se muestra (case-sensitive).

### ¿Puedo usar los mismos secrets en múltiples repositorios?

No, los secrets son específicos de cada repositorio. Si tienes múltiples repos, debes configurarlos en cada uno.

### ¿Los secrets funcionan en forks?

No, los secrets no se copian a forks por seguridad. Cada fork necesita configurar sus propios secrets.

---

## 🎉 ¡Listo!

Una vez que hayas configurado los secrets, el workflow de GitHub Actions podrá:

- ✅ Construir las imágenes Docker
- ✅ Subirlas a GHCR
- ✅ Notificar a Render para hacer deploy automático

¡Siguiente paso: Configurar Render! 🚀
