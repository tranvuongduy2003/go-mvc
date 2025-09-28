-- Create outbox_messages table for the outbox pattern
CREATE TABLE IF NOT EXISTS outbox_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL UNIQUE,
    event_type VARCHAR(255) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    retries INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    
    CONSTRAINT outbox_messages_status_check 
        CHECK (status IN ('pending', 'processed', 'failed'))
);

-- Create indexes for performance
CREATE INDEX idx_outbox_messages_status_created 
    ON outbox_messages (status, created_at);
    
CREATE INDEX idx_outbox_messages_message_id 
    ON outbox_messages (message_id);
    
CREATE INDEX idx_outbox_messages_event_type 
    ON outbox_messages (event_type);
    
CREATE INDEX idx_outbox_messages_aggregate_id 
    ON outbox_messages (aggregate_id);

-- Create partial index for pending and failed messages (most commonly queried)
CREATE INDEX idx_outbox_messages_pending_failed 
    ON outbox_messages (created_at) 
    WHERE status IN ('pending', 'failed');

-- Create partial index for processed messages for cleanup queries
CREATE INDEX idx_outbox_messages_processed_at 
    ON outbox_messages (processed_at) 
    WHERE status = 'processed';

COMMENT ON TABLE outbox_messages IS 'Stores messages for reliable publishing using the outbox pattern';
COMMENT ON COLUMN outbox_messages.message_id IS 'Unique identifier for the message, used for deduplication';
COMMENT ON COLUMN outbox_messages.event_type IS 'Type of event/message for routing purposes';
COMMENT ON COLUMN outbox_messages.aggregate_id IS 'ID of the aggregate that generated this message';
COMMENT ON COLUMN outbox_messages.payload IS 'JSON payload of the message';
COMMENT ON COLUMN outbox_messages.status IS 'Processing status: pending, processed, or failed';
COMMENT ON COLUMN outbox_messages.retries IS 'Number of retry attempts made';
COMMENT ON COLUMN outbox_messages.max_retries IS 'Maximum number of retry attempts allowed';