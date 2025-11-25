-- Drop message_deduplication table and related objects
DROP INDEX IF EXISTS idx_message_deduplication_active;
DROP INDEX IF EXISTS idx_message_deduplication_event_type;
DROP INDEX IF EXISTS idx_message_deduplication_expires_at;
DROP INDEX IF EXISTS idx_message_deduplication_message_consumer;

DROP TABLE IF EXISTS message_deduplication;