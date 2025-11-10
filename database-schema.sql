-- ============================================
-- ESQUEMA DE BASE DE DATOS - TP8
-- Base de datos: final_clj4 (PostgreSQL)
-- ============================================

-- ============================================
-- TABLA: users
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    nickname VARCHAR(50) NOT NULL,
    email VARCHAR(150) NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    type BOOLEAN NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para búsquedas rápidas
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_nickname ON users(nickname);

-- ============================================
-- TABLA: courses
-- ============================================
CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    description VARCHAR(300) NOT NULL,
    category VARCHAR(50) NOT NULL,
    instructor VARCHAR(100) NOT NULL,
    duration BIGINT NOT NULL,
    requirement VARCHAR(150) NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para búsquedas
CREATE INDEX IF NOT EXISTS idx_courses_title ON courses(title);
CREATE INDEX IF NOT EXISTS idx_courses_category ON courses(category);

-- ============================================
-- TABLA: subscriptions
-- ============================================
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, course_id)
);

-- Índices y foreign keys (relaciones)
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_course_id ON subscriptions(course_id);

-- ============================================
-- TABLA: comments
-- ============================================
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    text VARCHAR(500) NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_course_id ON comments(course_id);

-- ============================================
-- TABLA: files
-- ============================================
CREATE TABLE IF NOT EXISTS files (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    url VARCHAR(500) NOT NULL,
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_files_user_id ON files(user_id);
CREATE INDEX IF NOT EXISTS idx_files_course_id ON files(course_id);

-- ============================================
-- EJEMPLOS DE INSERCIONES
-- ============================================

-- ============================================
-- INSERTAR USUARIOS
-- ============================================
-- Usuario Estudiante (type = false)
INSERT INTO users (nickname, email, password_hash, type) 
VALUES 
    ('juan_perez', 'juan.perez@example.com', '$2a$10$hashedpassword123', false),
    ('maria_garcia', 'maria.garcia@example.com', '$2a$10$hashedpassword456', false),
    ('carlos_lopez', 'carlos.lopez@example.com', '$2a$10$hashedpassword789', false);

-- Usuario Instructor/Admin (type = true)
INSERT INTO users (nickname, email, password_hash, type) 
VALUES 
    ('prof_smith', 'prof.smith@example.com', '$2a$10$hashedpasswordadmin', true),
    ('admin_user', 'admin@example.com', '$2a$10$hashedpasswordadmin2', true);

-- ============================================
-- INSERTAR CURSOS
-- ============================================
INSERT INTO courses (title, description, category, instructor, duration, requirement) 
VALUES 
    ('Introducción a Go', 
     'Curso completo de programación en Go desde cero', 
     'Programación', 
     'prof_smith', 
     40, 
     'Conocimientos básicos de programación'),
    
    ('Next.js Avanzado', 
     'Aprende Next.js con SSR, SSG y API Routes', 
     'Desarrollo Web', 
     'prof_smith', 
     60, 
     'Conocimientos de React y JavaScript'),
    
    ('PostgreSQL para Desarrolladores', 
     'Curso completo de bases de datos PostgreSQL', 
     'Bases de Datos', 
     'admin_user', 
     50, 
     'Conocimientos básicos de SQL'),
    
    ('Docker y Contenedores', 
     'Aprende a containerizar aplicaciones con Docker', 
     'DevOps', 
     'admin_user', 
     30, 
     'Conocimientos básicos de Linux'),
    
    ('CI/CD con GitHub Actions', 
     'Automatiza tus despliegues con GitHub Actions', 
     'DevOps', 
     'prof_smith', 
     25, 
     'Conocimientos de Git y GitHub');

-- ============================================
-- INSERTAR SUSCRIPCIONES
-- ============================================
-- Usuario 1 se suscribe a cursos 1, 2, 3
INSERT INTO subscriptions (user_id, course_id) 
VALUES 
    (1, 1),
    (1, 2),
    (1, 3);

-- Usuario 2 se suscribe a cursos 2, 4
INSERT INTO subscriptions (user_id, course_id) 
VALUES 
    (2, 2),
    (2, 4);

-- Usuario 3 se suscribe a todos los cursos
INSERT INTO subscriptions (user_id, course_id) 
VALUES 
    (3, 1),
    (3, 2),
    (3, 3),
    (3, 4),
    (3, 5);

-- ============================================
-- INSERTAR COMENTARIOS
-- ============================================
INSERT INTO comments (user_id, course_id, text) 
VALUES 
    (1, 1, 'Excelente curso, muy bien explicado'),
    (1, 2, 'Me ayudó mucho a entender Next.js'),
    (2, 2, 'Buen contenido pero podría tener más ejemplos'),
    (2, 4, 'Docker es genial, gracias por el curso'),
    (3, 1, 'Perfecto para empezar con Go'),
    (3, 3, 'PostgreSQL es muy potente, buen curso'),
    (3, 5, 'GitHub Actions simplifica mucho el CI/CD');

-- ============================================
-- INSERTAR ARCHIVOS
-- ============================================
INSERT INTO files (user_id, course_id, name, url) 
VALUES 
    (1, 1, 'ejemplo_go.go', 'https://storage.example.com/files/ejemplo_go.go'),
    (1, 2, 'componente_nextjs.jsx', 'https://storage.example.com/files/componente_nextjs.jsx'),
    (2, 4, 'docker-compose.yml', 'https://storage.example.com/files/docker-compose.yml'),
    (3, 3, 'schema.sql', 'https://storage.example.com/files/schema.sql'),
    (3, 5, 'workflow.yml', 'https://storage.example.com/files/workflow.yml');

-- ============================================
-- CONSULTAS ÚTILES
-- ============================================

-- Ver todos los usuarios
SELECT * FROM users;

-- Ver todos los cursos
SELECT * FROM courses;

-- Ver suscripciones con nombres de usuario y curso
SELECT 
    u.nickname AS usuario,
    c.title AS curso,
    s.creation_date AS fecha_suscripcion
FROM subscriptions s
JOIN users u ON s.user_id = u.id
JOIN courses c ON s.course_id = c.id
ORDER BY s.creation_date DESC;

-- Ver comentarios con información del usuario y curso
SELECT 
    u.nickname AS usuario,
    c.title AS curso,
    cm.text AS comentario,
    cm.creation_date AS fecha
FROM comments cm
JOIN users u ON cm.user_id = u.id
JOIN courses c ON cm.course_id = c.id
ORDER BY cm.creation_date DESC;

-- Contar cursos por categoría
SELECT 
    category,
    COUNT(*) AS total_cursos
FROM courses
GROUP BY category
ORDER BY total_cursos DESC;

-- Usuarios con más suscripciones
SELECT 
    u.nickname,
    COUNT(s.id) AS total_suscripciones
FROM users u
LEFT JOIN subscriptions s ON u.id = s.user_id
GROUP BY u.id, u.nickname
ORDER BY total_suscripciones DESC;

-- ============================================
-- RELACIONES ENTRE TABLAS
-- ============================================
-- 
-- users (1) ──< subscriptions (N) ──> courses (1)
-- users (1) ──< comments (N) ──> courses (1)
-- users (1) ──< files (N) ──> courses (1)
--
-- Nota: Las relaciones se manejan a nivel de aplicación (Go/GORM)
-- No hay foreign keys explícitas en la base de datos por diseño
-- ============================================

