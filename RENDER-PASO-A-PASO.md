# 🎯 Render - Configuración Paso a Paso Actualizada

Esta guía te muestra exactamente qué verás en Render y dónde encontrar cada opción.

## 📍 Paso 1: Llegar a la Pantalla Correcta

### 1.1 Ir al Dashboard

1. Ve a [dashboard.render.com](https://dashboard.render.com)
2. Inicia sesión si es necesario

### 1.2 Crear Nuevo Servicio

1. Click en el botón **"New +"** (arriba a la derecha, puede ser verde o azul)
2. Verás un menú desplegable con opciones

### 1.3 Seleccionar el Tipo Correcto

**IMPORTANTE**: Debes seleccionar la opción correcta:

- ❌ **NO selecciones**: "Web Service" directamente
- ❌ **NO selecciones**: "Deploy from Git"
- ✅ **SÍ selecciona**: Busca la opción que diga algo como:
  - **"Deploy an existing image from a registry"**
  - **"Docker Registry"**
  - **"Container Registry"**
  - O en algunos casos: **"Web Service"** → Luego busca la opción de registry

**Si no ves estas opciones**, puede que Render haya cambiado la interfaz. Sigue estos pasos alternativos:

---

## 🔄 Opción Alternativa: Si No Ves "Deploy Existing Image"

### Método 1: Buscar en Web Service

1. Click en **"New +"** → **"Web Service"**
2. En la pantalla que aparece, busca:
   - Una pestaña o sección que diga **"Docker"** o **"Container"**
   - O un botón que diga **"Deploy from registry"**
   - O un selector que te permita elegir entre "Git" y "Docker Registry"

### Método 2: Usar Blueprint (render.yaml)

Si Render no te muestra la opción de registry directamente:

1. Asegúrate de tener el archivo `render.yaml` en tu repositorio
2. En Render, selecciona **"New Blueprint"** o **"Deploy from Git"**
3. Conecta tu repositorio de GitHub
4. Render detectará el `render.yaml` y usará esa configuración

**Nota**: El `render.yaml` que creamos está configurado para construir desde Dockerfile, no desde imagen. Si quieres usar imágenes de GHCR, necesitas configurarlo manualmente.

---

## 🎯 Lo Que Deberías Ver (Pantalla Correcta)

Cuando estés en la pantalla correcta, deberías ver campos como:

### Campos que Render DEBERÍA Mostrar:

1. **"Registry"** o **"Image Source"** o **"Deploy from"**

   - Con opciones como: "Docker Hub", "Docker Registry", "GitHub Container Registry", etc.

2. **"Image URL"** o **"Docker Image"** o **"Container Image"**

   - Un campo de texto donde pegar la URL de la imagen

3. **"Registry Credentials"** o **"Authentication"**

   - Username
   - Password

4. **"Name"** o **"Service Name"**

   - Nombre del servicio

5. **"Environment"** o **"Runtime"**

   - Opciones como: "Docker", "Node", "Python", etc.

6. **"Plan"**
   - Free, Starter, Standard, etc.

---

## 🔍 Si No Ves Esos Campos

### Posibles Razones:

1. **Estás en la opción incorrecta**

   - Vuelve atrás y busca "Deploy existing image" o "Docker Registry"

2. **Render ha cambiado su interfaz**

   - Busca en la documentación de Render: [render.com/docs](https://render.com/docs)
   - O busca "deploy docker image render" en Google

3. **Necesitas una cuenta Pro** (poco probable)
   - El tier Free debería permitir desplegar desde registry

---

## 📸 Qué Buscar Exactamente

### En el Menú "New +":

Busca opciones como:

- ✅ "Deploy an existing image"
- ✅ "Docker Registry"
- ✅ "Container Registry"
- ✅ "Web Service" (y luego busca opción de Docker/Registry)

### En la Pantalla de Configuración:

Busca secciones o pestañas como:

- ✅ "Docker"
- ✅ "Container"
- ✅ "Registry"
- ✅ "Image"

---

## 🆘 Si Aún No Lo Encuentras

### Opción 1: Contactar Soporte de Render

1. Ve a [render.com/support](https://render.com/support)
2. Pregunta: "How do I deploy a Docker image from GitHub Container Registry?"

### Opción 2: Usar Render CLI

Si tienes problemas con la interfaz web, puedes usar la CLI:

```bash
# Instalar Render CLI
npm install -g render-cli

# Login
render login

# Crear servicio desde imagen
render services:create \
  --name tp8-backend \
  --type web \
  --docker-image ghcr.io/[tu-usuario]/backend:latest \
  --docker-registry ghcr.io \
  --docker-username [tu-usuario] \
  --docker-password [tu-token]
```

### Opción 3: Configurar Manualmente Después

1. Crea el servicio con cualquier método (puede ser desde Git)
2. Una vez creado, ve a **Settings**
3. Busca la sección **"Docker"** o **"Container"**
4. Ahí deberías poder cambiar la configuración para usar una imagen

---

## 💡 Método Alternativo: Construir desde Dockerfile

Si Render no te permite desplegar desde registry directamente, puedes:

1. **Conectar tu repositorio de GitHub a Render**
2. **Seleccionar "Build from Dockerfile"**
3. Render construirá la imagen desde el Dockerfile en tu repo
4. Esto funciona, pero las imágenes se construirán en Render, no usarán las de GHCR

**Para usar este método**:

- Render → New → Web Service
- Conecta tu repositorio de GitHub
- Selecciona la rama (main/master)
- Render detectará el Dockerfile automáticamente
- Configura las variables de entorno
- Create

---

## ❓ Preguntas para Ayudarte Mejor

Para ayudarte mejor, dime:

1. **¿Qué opciones ves cuando haces click en "New +"?**

   - Lista todas las opciones que aparecen

2. **¿Estás en la pantalla de crear "Web Service"?**

   - ¿Qué campos o secciones ves?

3. **¿Ves alguna opción relacionada con "Docker", "Container" o "Registry"?**

   - Aunque sea en otro lugar

4. **¿Qué versión de Render estás usando?**
   - Free, Starter, etc.

Con esta información puedo darte instrucciones más específicas para tu caso.

---

## 📚 Recursos

- [Render Docs - Docker](https://render.com/docs/docker)
- [Render Docs - Private Images](https://render.com/docs/private-images)
- [Render Support](https://render.com/support)
