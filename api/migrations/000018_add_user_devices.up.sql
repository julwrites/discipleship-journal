CREATE TABLE IF NOT EXISTS user_devices (
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fcm_token VARCHAR(255) NOT NULL,
    device_type TEXT,
    last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, fcm_token)
);

CREATE INDEX idx_user_devices_user_id ON user_devices(user_id);
