ALTER TABLE study_templates
ADD COLUMN bible_references JSONB DEFAULT '[]',
ADD COLUMN allow_user_passages BOOLEAN DEFAULT FALSE,
ADD COLUMN template_body TEXT DEFAULT '',
ADD COLUMN required_version TEXT DEFAULT '';
