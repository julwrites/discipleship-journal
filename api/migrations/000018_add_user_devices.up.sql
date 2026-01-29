CREATE TABLE IF NOT EXISTS user_devices (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fcm_token TEXT NOT NULL,
    device_type TEXT,
    last_used_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, fcm_token)
);

CREATE INDEX idx_user_devices_user_id ON user_devices(user_id);
