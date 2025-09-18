-- Удалим пустые значения, если вдруг есть
UPDATE users
SET nickname = 'user_' || id
WHERE nickname IS NULL OR nickname = '';

-- Сделаем колонку обязательной и уникальной
ALTER TABLE users
ALTER COLUMN nickname SET NOT NULL;

ALTER TABLE users
ADD CONSTRAINT users_nickname_unique UNIQUE (nickname);
