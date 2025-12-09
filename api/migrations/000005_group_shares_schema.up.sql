CREATE TABLE IF NOT EXISTS group_shares (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    note_id UUID REFERENCES notes(id) ON DELETE CASCADE,
    shared_by UUID REFERENCES users(id) ON DELETE CASCADE,
    shared_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    comment TEXT,
    UNIQUE(group_id, note_id)
);

CREATE INDEX idx_group_shares_group_id ON group_shares(group_id);
