DROP TABLE IF EXISTS users;
ALTER TABLE IF EXISTS bookmarks DROP COLUMN user_fk IF EXISTS CASCADE;