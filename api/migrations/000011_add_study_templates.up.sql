CREATE TABLE study_templates (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    creator_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    structure JSON, -- Sections like Passage, Reflection, etc.
    prompts JSON, -- AI System Prompts
    fields JSON, -- User input fields definition
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_study_templates_creator_id ON study_templates(creator_id);
CREATE INDEX idx_study_templates_is_public ON study_templates(is_public);
