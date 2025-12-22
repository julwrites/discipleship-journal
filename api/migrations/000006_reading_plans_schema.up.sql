CREATE TABLE IF NOT EXISTS reading_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    days INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reading_plan_days (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reading_plan_id UUID NOT NULL REFERENCES reading_plans(id) ON DELETE CASCADE,
    day_number INT NOT NULL,
    passage TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(reading_plan_id, day_number)
);

CREATE TABLE IF NOT EXISTS user_reading_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reading_plan_id UUID NOT NULL REFERENCES reading_plans(id) ON DELETE CASCADE,
    start_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, completed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_reading_plan_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_reading_plan_id UUID NOT NULL REFERENCES user_reading_plans(id) ON DELETE CASCADE,
    day_number INT NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_reading_plan_id, day_number)
);

CREATE INDEX idx_reading_plans_created_at ON reading_plans(created_at);
CREATE INDEX idx_user_reading_plans_user_id ON user_reading_plans(user_id);
CREATE INDEX idx_user_reading_plan_progress_plan_id ON user_reading_plan_progress(user_reading_plan_id);
