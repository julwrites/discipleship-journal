CREATE TABLE IF NOT EXISTS reading_plans (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    days INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reading_plan_days (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    reading_plan_id VARCHAR(36) NOT NULL REFERENCES reading_plans(id) ON DELETE CASCADE,
    day_number INT NOT NULL,
    passage TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(reading_plan_id, day_number)
);

CREATE TABLE IF NOT EXISTS user_reading_plans (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reading_plan_id VARCHAR(36) NOT NULL REFERENCES reading_plans(id) ON DELETE CASCADE,
    start_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, completed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_reading_plan_progress (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_reading_plan_id VARCHAR(36) NOT NULL REFERENCES user_reading_plans(id) ON DELETE CASCADE,
    day_number INT NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_reading_plan_id, day_number)
);

CREATE INDEX idx_reading_plans_created_at ON reading_plans(created_at);
CREATE INDEX idx_user_reading_plans_user_id ON user_reading_plans(user_id);
CREATE INDEX idx_user_reading_plan_progress_plan_id ON user_reading_plan_progress(user_reading_plan_id);
