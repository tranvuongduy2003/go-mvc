-- Create message_deduplication table for lightweight deduplication
CREATE TABLE IF NOT EXISTS message_deduplication (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    consumer_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Ensure uniqueness of message_id per consumer
    CONSTRAINT message_deduplication_message_consumer_unique 
        UNIQUE (message_id, consumer_id)
);

-- Create indexes for performance
CREATE INDEX idx_message_deduplication_message_consumer 
    ON message_deduplication (message_id, consumer_id);
    
CREATE INDEX idx_message_deduplication_expires_at 
    ON message_deduplication (expires_at);
    
CREATE INDEX idx_message_deduplication_event_type 
    ON message_deduplication (event_type);

-- Create partial index for non-expired records (most commonly queried)
CREATE INDEX idx_message_deduplication_active 
    ON message_deduplication (message_id, consumer_id, expires_at) 
    WHERE expires_at > NOW();

COMMENT ON TABLE message_deduplication IS 'Lightweight message deduplication table with TTL-based expiry';
COMMENT ON COLUMN message_deduplication.message_id IS 'Unique identifier of the processed message';
COMMENT ON COLUMN message_deduplication.consumer_id IS 'Identifier of the consumer that processed the message';
COMMENT ON COLUMN message_deduplication.event_type IS 'Type of event/message processed';
COMMENT ON COLUMN message_deduplication.processed_at IS 'When the message was processed';
COMMENT ON COLUMN message_deduplication.expires_at IS 'When this deduplication record expires and can be cleaned up';