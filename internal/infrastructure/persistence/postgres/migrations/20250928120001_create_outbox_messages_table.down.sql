-- Drop outbox_messages table and related objects
DROP INDEX IF EXISTS idx_outbox_messages_processed_at;
DROP INDEX IF EXISTS idx_outbox_messages_pending_failed;
DROP INDEX IF EXISTS idx_outbox_messages_aggregate_id;
DROP INDEX IF EXISTS idx_outbox_messages_event_type;
DROP INDEX IF EXISTS idx_outbox_messages_message_id;
DROP INDEX IF EXISTS idx_outbox_messages_status_created;

DROP TABLE IF EXISTS outbox_messages;