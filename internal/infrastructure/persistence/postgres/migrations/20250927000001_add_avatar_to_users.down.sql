ALTER TABLE users 
DROP COLUMN IF EXISTS avatar_file_key,
DROP COLUMN IF EXISTS avatar_cdn_url;