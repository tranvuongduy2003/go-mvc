-- Create inbox_messages table for the inbox pattern
CREATE TABLE IF NOT EXISTS inbox_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    consumer_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'received',
    received_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT inbox_messages_status_check 
        CHECK (status IN ('received', 'processed', 'ignored')),
        
    -- Ensure uniqueness of message_id per consumer
    CONSTRAINT inbox_messages_message_consumer_unique 
        UNIQUE (message_id, consumer_id)
);

-- Create indexes for performance
CREATE INDEX idx_inbox_messages_message_consumer 
    ON inbox_messages (message_id, consumer_id);
    
CREATE INDEX idx_inbox_messages_consumer_status 
    ON inbox_messages (consumer_id, status);
    
CREATE INDEX idx_inbox_messages_event_type 
    ON inbox_messages (event_type);
    
CREATE INDEX idx_inbox_messages_received_at 
    ON inbox_messages (received_at);

-- Create index for cleanup operations
CREATE INDEX idx_inbox_messages_created_at 
    ON inbox_messages (created_at);

COMMENT ON TABLE inbox_messages IS 'Stores received messages for deduplication using the inbox pattern';
COMMENT ON COLUMN inbox_messages.message_id IS 'Unique identifier of the received message';
COMMENT ON COLUMN inbox_messages.event_type IS 'Type of event/message received';
COMMENT ON COLUMN inbox_messages.consumer_id IS 'Identifier of the consumer that received the message';
COMMENT ON COLUMN inbox_messages.status IS 'Processing status: received, processed, or ignored (duplicate)';
COMMENT ON COLUMN inbox_messages.received_at IS 'When the message was first received';
COMMENT ON COLUMN inbox_messages.processed_at IS 'When the message was successfully processed';