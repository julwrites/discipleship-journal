CREATE TABLE study_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    structure JSONB DEFAULT '{}', -- Sections like Passage, Reflection, etc.
    prompts JSONB DEFAULT '{}', -- AI System Prompts
    fields JSONB DEFAULT '[]', -- User input fields definition
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_study_templates_creator_id ON study_templates(creator_id);
CREATE INDEX idx_study_templates_is_public ON study_templates(is_public);
