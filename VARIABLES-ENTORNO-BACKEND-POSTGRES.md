# 🔧 Variables de Entorno - Backend (PostgreSQL)

Estas son las variables de entorno **exactas** que necesita tu backend con PostgreSQL.

## ✅ Variables Requeridas

Agrega estas variables de entorno en Render (una por una):

| Key           | Value                                                   | Descripción                         | Requerido                           |
| ------------- | ------------------------------------------------------- | ----------------------------------- | ----------------------------------- |
| `PORT`        | `8080`                                                  | Puerto donde corre el servidor Go   | ✅ Sí                               |
| `DB_HOST`     | `dpg-d48vp9be5dus73cf9250-a.oregon-postgres.render.com` | Host de tu base de datos PostgreSQL | ✅ Sí                               |
| `DB_PORT`     | `5432`                                                  | Puerto de PostgreSQL                | ⚠️ Opcional (default: 5432)         |
| `DB_USER`     | `admin`                                                 | Usuario de PostgreSQL               | ✅ Sí                               |
| `DB_PASSWORD` | `gImS5z9Riw8683kfI5uMY5NSNYTGXclu`                      | Contraseña de PostgreSQL            | ✅ Sí                               |
| `DB_NAME`     | `final_clj4`                                            | Nombre de la base de datos          | ⚠️ Opcional (default: "final_clj4") |
| `DB_SSLMODE`  | `require`                                               | Modo SSL para PostgreSQL            | ⚠️ Opcional (default: "require")    |

---

## 📝 Valores Exactos para Render

Copia y pega estos valores exactos en Render:

```
PORT=8080
DB_HOST=dpg-d48vp9be5dus73cf9250-a.oregon-postgres.render.com
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=gImS5z9Riw8683kfI5uMY5NSNYTGXclu
DB_NAME=final_clj4
DB_SSLMODE=require
```

---

## 🔧 Cómo Agregarlas en Render

### Paso a Paso:

1. En Render, cuando estés configurando el servicio backend
2. Busca la sección **"Environment Variables"** o **"Env Vars"**
3. Click en **"Add Environment Variable"** o **"+"**
4. Agrega cada variable una por una:

### Variable 1: PORT

- **Key**: `PORT`
- **Value**: `8080`

### Variable 2: DB_HOST

- **Key**: `DB_HOST`
- **Value**: `dpg-d48vp9be5dus73cf9250-a.oregon-postgres.render.com`

### Variable 3: DB_PORT

- **Key**: `DB_PORT`
- **Value**: `5432`

### Variable 4: DB_USER

- **Key**: `DB_USER`
- **Value**: `admin`

### Variable 5: DB_PASSWORD

- **Key**: `DB_PASSWORD`
- **Value**: `gImS5z9Riw8683kfI5uMY5NSNYTGXclu`
- ⚠️ **Importante**: Copia exactamente, sin espacios

### Variable 6: DB_NAME

- **Key**: `DB_NAME`
- **Value**: `final_clj4`

### Variable 7: DB_SSLMODE

- **Key**: `DB_SSLMODE`
- **Value**: `require`

---

## ⚠️ Valores por Defecto (Si No Los Configuras)

Si no configuras algunas variables, el backend usará estos valores por defecto:

| Variable      | Valor por Defecto |
| ------------- | ----------------- |
| `PORT`        | `8080`            |
| `DB_HOST`     | `127.0.0.1`       |
| `DB_PORT`     | `5432`            |
| `DB_USER`     | `admin`           |
| `DB_PASSWORD` | `""` (vacío)      |
| `DB_NAME`     | `final_clj4`      |
| `DB_SSLMODE`  | `require`         |

**⚠️ IMPORTANTE**: En producción, **DEBES** configurar todas las variables, especialmente:

- `DB_HOST` (no puede ser 127.0.0.1 en Render)
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`

---

## 🔍 Cómo Verificar que Funcionan

Después de configurar las variables:

1. Ve a tu servicio en Render
2. Click en **"Logs"**
3. Deberías ver algo como:
   ```
   DB connected host=dpg-d48vp9be5dus73cf9250-a.oregon-postgres.render.com db=final_clj4 attempt=1
   ```
4. Si ves errores de conexión, verifica:
   - Que `DB_HOST` sea correcto
   - Que `DB_USER` y `DB_PASSWORD` sean correctos
   - Que `DB_NAME` exista en tu PostgreSQL
   - Que `DB_SSLMODE=require` esté configurado

---

## 📋 Checklist Rápido

Antes de crear el servicio, asegúrate de tener:

- [ ] `PORT=8080`
- [ ] `DB_HOST=dpg-d48vp9be5dus73cf9250-a.oregon-postgres.render.com`
- [ ] `DB_PORT=5432`
- [ ] `DB_USER=admin`
- [ ] `DB_PASSWORD=gImS5z9Riw8683kfI5uMY5NSNYTGXclu`
- [ ] `DB_NAME=final_clj4`
- [ ] `DB_SSLMODE=require`

---

## 🆘 Troubleshooting

### Error: "Connection refused" o "Cannot connect to database"

**Causa**: `DB_HOST` incorrecto o PostgreSQL no accesible

**Solución**:

- Verifica que `DB_HOST` sea exactamente: `dpg-d48vp9be5dus73cf9250-a.oregon-postgres.render.com`
- Verifica que PostgreSQL permita conexiones desde fuera
- Verifica que el firewall de Render permita conexiones

### Error: "Access denied for user" o "password authentication failed"

**Causa**: `DB_USER` o `DB_PASSWORD` incorrectos

**Solución**:

- Verifica que `DB_USER=admin`
- Verifica que `DB_PASSWORD=gImS5z9Riw8683kfI5uMY5NSNYTGXclu` (sin espacios)
- Asegúrate de no tener espacios extra en los valores

### Error: "database does not exist"

**Causa**: `DB_NAME` no existe

**Solución**:

- Verifica que `DB_NAME=final_clj4`
- La base de datos debería existir en tu PostgreSQL de Render

### Error: "SSL connection required"

**Causa**: `DB_SSLMODE` no está configurado o es incorrecto

**Solución**:

- Asegúrate de que `DB_SSLMODE=require`
- Render PostgreSQL requiere SSL

---

## ✅ Listo

Una vez que agregues todas estas variables en Render, tu backend debería conectarse correctamente a PostgreSQL.
