-- Drop inbox_messages table and related objects
DROP INDEX IF EXISTS idx_inbox_messages_created_at;
DROP INDEX IF EXISTS idx_inbox_messages_received_at;
DROP INDEX IF EXISTS idx_inbox_messages_event_type;
DROP INDEX IF EXISTS idx_inbox_messages_consumer_status;
DROP INDEX IF EXISTS idx_inbox_messages_message_consumer;

DROP TABLE IF EXISTS inbox_messages;