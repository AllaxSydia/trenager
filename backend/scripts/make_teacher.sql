-- Скрипт для назначения пользователя преподавателем
-- Использование: 
-- 1. Замените 'user@example.com' на email пользователя, которого хотите сделать преподавателем
-- 2. Выполните этот скрипт в вашей базе данных

-- Пример 1: Назначить преподавателем по email
UPDATE users 
SET role = 'teacher' 
WHERE email = 'user@example.com';

-- Пример 2: Назначить преподавателем по username
-- UPDATE users 
-- SET role = 'teacher' 
-- WHERE username = 'teacher_username';

-- Пример 3: Назначить преподавателем по ID
-- UPDATE users 
-- SET role = 'teacher' 
-- WHERE id = 1;

-- Проверка результата
SELECT id, username, email, role 
FROM users 
WHERE role = 'teacher';

