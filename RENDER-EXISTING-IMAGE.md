# 🎯 Render - Configurar con "Existing Image"

Perfecto, veo que tienes la opción **"Existing image"**. Esa es la correcta. Sigue estos pasos:

## 📍 Paso 1: Seleccionar "Existing image"

1. En Render, click en **"New +"** → **"Web Service"**
2. Verás las opciones:
   - Git provider
   - Public git repository
   - **Existing image** ← **SELECCIONA ESTA**

## 🔧 Paso 2: Configurar el Servicio Backend

### 2.1 Después de seleccionar "Existing image"

Deberías ver una pantalla con campos. Completa:

### **Image URL** o **Docker Image URL**

```
ghcr.io/[tu-usuario-github]/backend:latest
```

**Ejemplo real**:

```
ghcr.io/felipeeguia03/backend:latest
```

**⚠️ IMPORTANTE**: Reemplaza `[tu-usuario-github]` con tu usuario real de GitHub.

### **Registry** (si te lo pide)

Si hay un campo "Registry" o "Image Source":

- Selecciona: **"Docker Registry"** o **"Other"** o **"Custom"**

### **Authentication** o **Registry Credentials**

Busca una sección que diga:

- "Authentication"
- "Registry Credentials"
- "Docker Registry Credentials"
- O un botón "Add credentials"

Ahí deberías ver:

**Username**:

```
[tu-usuario-github]
```

(Tu usuario de GitHub, el mismo que en la URL de la imagen)

**Password**:

```
[tu-token-ghp_...]
```

(El Personal Access Token que creaste, el que empieza con `ghp_`)

### **Name** o **Service Name**

```
tp8-backend
```

### **Environment** o **Runtime**

- Selecciona: **"Docker"**

### **Region**

- Elige la más cercana (Oregon, Frankfurt, Singapore, etc.)

### **Plan**

- Selecciona: **"Free"**

## 🔧 Paso 3: Variables de Entorno

Busca una sección que diga:

- "Environment Variables"
- "Env Vars"
- "Environment"

Agrega estas variables (una por una):

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

## 🔧 Paso 4: Health Check

Busca:

- "Health Check Path"
- "Health Check"
- "Health"

Escribe:

```
/healthz
```

## ✅ Paso 5: Crear el Servicio

1. Revisa toda la configuración
2. Click en **"Create Web Service"** o **"Create"**
3. Render comenzará a descargar la imagen y desplegarla
4. Espera 2-5 minutos

## 🔗 Paso 6: Obtener Webhook URL

Una vez que el servicio esté creado:

1. Ve a tu servicio (click en el nombre `tp8-backend`)
2. En el menú lateral, click en **"Settings"**
3. Busca la sección **"Manual Deploy"** o **"Deploy Hook"**
4. Verás un campo **"Webhook URL"** o un botón **"Copy webhook URL"**
5. Copia la URL completa (algo como: `https://api.render.com/deploy/srv-xxxxx?key=xxxxx`)
6. **Guarda esta URL** - la necesitarás para GitHub Secrets

---

## 🎨 Paso 7: Crear Servicio Frontend (Mismo Proceso)

Repite los mismos pasos pero para el frontend:

### 7.1 Seleccionar "Existing image"

1. **New +** → **Web Service** → **Existing image**

### 7.2 Configurar

**Image URL**:

```
ghcr.io/[tu-usuario-github]/frontend:latest
```

**Authentication** (mismo usuario y token):

- **Username**: Tu usuario de GitHub
- **Password**: El mismo token (`ghp_...`)

**Name**:

```
tp8-frontend
```

**Environment**: `Docker`

**Region**: La misma que el backend

**Plan**: `Free`

### 7.3 Variables de Entorno

| Key                   | Value                             |
| --------------------- | --------------------------------- |
| `PORT`                | `8080`                            |
| `NEXT_PUBLIC_API_URL` | `https://tu-backend.onrender.com` |

**⚠️ IMPORTANTE**:

- Reemplaza `tu-backend.onrender.com` con la URL real de tu backend
- La encontrarás en el dashboard de Render, en la página del servicio backend
- Debería verse algo como: `https://tp8-backend-xxxxx.onrender.com`

### 7.4 Crear y Obtener Webhook

1. **Create Web Service**
2. Espera a que se despliegue
3. **Settings** → **Manual Deploy** → Copia el **Webhook URL**

---

## 🔐 Paso 8: Agregar Webhooks a GitHub

1. Ve a tu repositorio → **Settings** → **Secrets and variables** → **Actions**
2. Agrega:
   - `RENDER_BACKEND_WEBHOOK_URL` → URL del webhook del backend
   - `RENDER_FRONTEND_WEBHOOK_URL` → URL del webhook del frontend

---

## 🐛 Si No Ves Algún Campo

### No veo "Authentication" o "Registry Credentials"

Algunas veces Render no muestra estos campos hasta que:

1. Escribes la URL de la imagen
2. O haces click en "Advanced" o "Show more options"

**Intenta**:

- Escribir primero la Image URL
- Buscar un botón "Advanced" o "More options"
- Hacer scroll hacia abajo en la página

### No veo dónde poner el Username/Password

Puede que Render use un formato diferente:

1. Busca un botón **"Add credentials"** o **"Configure registry"**
2. O puede que te pida las credenciales después de hacer click en "Create"
3. O puede estar en **Settings** después de crear el servicio

**Si no encuentras dónde poner las credenciales**:

- Crea el servicio primero
- Luego ve a **Settings** → Busca **"Docker"** o **"Registry"**
- Ahí deberías poder configurar las credenciales

---

## ✅ Resumen Rápido

```
1. New + → Web Service → Existing image
2. Image URL: ghcr.io/[usuario]/backend:latest
3. Username: [tu-usuario-github]
4. Password: [tu-token-ghp_...]
5. Name: tp8-backend
6. Environment: Docker
7. Variables de entorno (PORT, MYSQL_*, etc.)
8. Health Check: /healthz
9. Create
10. Settings → Manual Deploy → Copiar Webhook URL
```

---

## 🆘 Si Aún Tienes Problemas

Dime exactamente:

1. **¿Qué campos ves después de seleccionar "Existing image"?**

   - Lista todos los campos que aparecen

2. **¿Ves algún campo relacionado con "Authentication" o "Credentials"?**

   - O un botón "Add credentials"

3. **¿Puedes hacer una captura de pantalla?**
   - O describe qué ves en la pantalla

Con esa información puedo ayudarte mejor.
