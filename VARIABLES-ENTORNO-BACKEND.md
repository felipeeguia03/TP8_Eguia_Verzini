# 🔧 Variables de Entorno - Backend

Estas son las variables de entorno **exactas** que necesita tu backend según el código.

## ✅ Variables Requeridas

Agrega estas variables de entorno en Render (una por una):

| Key              | Value               | Descripción                       | Requerido                      |
| ---------------- | ------------------- | --------------------------------- | ------------------------------ |
| `PORT`           | `8080`              | Puerto donde corre el servidor Go | ✅ Sí                          |
| `MYSQL_HOST`     | `tu-mysql-host`     | Host de tu base de datos MySQL    | ✅ Sí                          |
| `MYSQL_PORT`     | `3306`              | Puerto de MySQL                   | ⚠️ Opcional (default: 3306)    |
| `MYSQL_USER`     | `tu-usuario-mysql`  | Usuario de MySQL                  | ✅ Sí                          |
| `MYSQL_PASSWORD` | `tu-password-mysql` | Contraseña de MySQL               | ✅ Sí                          |
| `MYSQL_DB`       | `tu-database`       | Nombre de la base de datos        | ⚠️ Opcional (default: "final") |
| `MYSQL_SSL`      | `false`             | Usar SSL para MySQL               | ⚠️ Opcional (default: "false") |

---

## 📝 Cómo Agregarlas en Render

### Paso a Paso:

1. En Render, cuando estés configurando el servicio backend
2. Busca la sección **"Environment Variables"** o **"Env Vars"**
3. Click en **"Add Environment Variable"** o **"+"**
4. Agrega cada variable una por una:

### Variable 1: PORT

- **Key**: `PORT`
- **Value**: `8080`

### Variable 2: MYSQL_HOST

- **Key**: `MYSQL_HOST`
- **Value**: `[tu-host-mysql]`
  - Ejemplo si usas Render MySQL: `dpg-xxxxx-a.oregon-postgres.render.com`
  - Ejemplo si usas otro servicio: `tu-servidor-mysql.com`

### Variable 3: MYSQL_PORT

- **Key**: `MYSQL_PORT`
- **Value**: `3306`
  - (Solo si tu MySQL no usa el puerto 3306, cámbialo)

### Variable 4: MYSQL_USER

- **Key**: `MYSQL_USER`
- **Value**: `[tu-usuario-mysql]`
  - Ejemplo: `root`, `admin`, `usuario_db`, etc.

### Variable 5: MYSQL_PASSWORD

- **Key**: `MYSQL_PASSWORD`
- **Value**: `[tu-password-mysql]`
  - ⚠️ **Importante**: Esta es sensible, asegúrate de escribirla correctamente

### Variable 6: MYSQL_DB

- **Key**: `MYSQL_DB`
- **Value**: `[nombre-de-tu-database]`
  - Ejemplo: `final`, `cursos`, `tp8_db`, etc.
  - Si no la agregas, usará `"final"` por defecto

### Variable 7: MYSQL_SSL

- **Key**: `MYSQL_SSL`
- **Value**: `false`
  - Usa `true` solo si tu MySQL requiere SSL
  - Para desarrollo/local usa `false`

---

## 🎯 Ejemplo Completo

Si tu MySQL está en Render o en otro servicio, aquí tienes un ejemplo:

```
PORT=8080
MYSQL_HOST=dpg-xxxxx-a.oregon-postgres.render.com
MYSQL_PORT=3306
MYSQL_USER=admin
MYSQL_PASSWORD=MiPassword123!
MYSQL_DB=final
MYSQL_SSL=false
```

---

## ⚠️ Valores por Defecto (Si No Los Configuras)

Si no configuras algunas variables, el backend usará estos valores por defecto:

| Variable         | Valor por Defecto |
| ---------------- | ----------------- |
| `PORT`           | `8080`            |
| `MYSQL_HOST`     | `127.0.0.1`       |
| `MYSQL_PORT`     | `3306`            |
| `MYSQL_USER`     | `root`            |
| `MYSQL_PASSWORD` | `""` (vacío)      |
| `MYSQL_DB`       | `final`           |
| `MYSQL_SSL`      | `false`           |

**⚠️ IMPORTANTE**: En producción, **DEBES** configurar todas las variables, especialmente:

- `MYSQL_HOST` (no puede ser 127.0.0.1 en Render)
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DB`

---

## 🔍 Cómo Verificar que Funcionan

Después de configurar las variables:

1. Ve a tu servicio en Render
2. Click en **"Logs"**
3. Deberías ver algo como:
   ```
   DB connected host=tu-mysql-host db=tu-database attempt=1
   ```
4. Si ves errores de conexión, verifica:
   - Que `MYSQL_HOST` sea accesible desde Render
   - Que `MYSQL_USER` y `MYSQL_PASSWORD` sean correctos
   - Que `MYSQL_DB` exista en tu MySQL

---

## 📋 Checklist Rápido

Antes de crear el servicio, asegúrate de tener:

- [ ] `PORT=8080`
- [ ] `MYSQL_HOST` (host real de tu MySQL)
- [ ] `MYSQL_USER` (usuario real)
- [ ] `MYSQL_PASSWORD` (contraseña real)
- [ ] `MYSQL_DB` (nombre de la base de datos)
- [ ] `MYSQL_PORT=3306` (si no es el puerto estándar)
- [ ] `MYSQL_SSL=false` (o `true` si necesitas SSL)

---

## 🆘 Troubleshooting

### Error: "Connection refused" o "Cannot connect to database"

**Causa**: `MYSQL_HOST` incorrecto o MySQL no accesible

**Solución**:

- Verifica que `MYSQL_HOST` sea el host correcto (no `127.0.0.1` o `localhost`)
- Si usas Render MySQL, usa el host que te da Render
- Verifica que MySQL permita conexiones desde fuera

### Error: "Access denied for user"

**Causa**: `MYSQL_USER` o `MYSQL_PASSWORD` incorrectos

**Solución**:

- Verifica que el usuario y contraseña sean correctos
- Asegúrate de no tener espacios extra en los valores
- Verifica que el usuario tenga permisos en la base de datos

### Error: "Unknown database"

**Causa**: `MYSQL_DB` no existe

**Solución**:

- Verifica que la base de datos exista en tu MySQL
- Crea la base de datos si no existe:
  ```sql
  CREATE DATABASE tu_database;
  ```

---

## ✅ Listo

Una vez que agregues todas estas variables en Render, tu backend debería conectarse correctamente a MySQL.
