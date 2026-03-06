ALTER TABLE study_templates
ADD COLUMN bible_references JSON,
ADD COLUMN allow_user_passages BOOLEAN DEFAULT false,
ADD COLUMN template_body TEXT,
ADD COLUMN required_version VARCHAR(50) DEFAULT '';
