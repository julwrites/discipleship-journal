-- Create verse_packs table
CREATE TABLE verse_packs (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE, -- NULL for system packs
    title VARCHAR(100) NOT NULL,
    identifier VARCHAR(50), -- Increased to 50 for safety
    description TEXT,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_verse_packs_user_id ON verse_packs(user_id);
-- Modify memory_verses table
DELETE FROM memory_verses;
-- Remove old incompatible data
ALTER TABLE memory_verses ADD COLUMN verse_pack_id VARCHAR(36) REFERENCES verse_packs(id) ON DELETE CASCADE;
ALTER TABLE memory_verses DROP COLUMN pack_name;
ALTER TABLE memory_verses DROP COLUMN text;
ALTER TABLE memory_verses DROP COLUMN user_id;
ALTER TABLE memory_verses MODIFY verse_pack_id VARCHAR(36) NOT NULL;
-- Modify group_shares table
ALTER TABLE group_shares ADD COLUMN verse_pack_id VARCHAR(36) REFERENCES verse_packs(id) ON DELETE CASCADE;
ALTER TABLE group_shares MODIFY note_id VARCHAR(36) NULL;
ALTER TABLE group_shares ADD CONSTRAINT chk_share_target CHECK ( (note_id IS NOT NULL AND verse_pack_id IS NULL) OR (note_id IS NULL AND verse_pack_id IS NOT NULL) );
-- Seed System Packs and Verses

    -- Pack: Living the New Life
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Living the New Life', 'A', 'TMS 60', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', '2 Corinthians 5:17', 'ESV', '["Christ Jesus", "New Creation", "Old", "New"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Galatians 2:20', 'ESV', '["Christ Jesus", "Christ Crucified", "Faith", "Love"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Romans 12:1', 'ESV', '["Mercy", "Lord God", "Living Sacrifice", "Worship"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'John 14:21', 'ESV', '["Obey", "Love", "Lord God", "Jesus Christ"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', '2 Timothy 3:16', 'ESV', '["Lord God", "Scripture", "Teaching", "Righteousness"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Joshua 1:8', 'ESV', '["Law", "Meditate", "Do", "Prosperous"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'John 15:7', 'ESV', '["Remain", "Ask", "Given", "Pray"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Philippians 4:6-7', 'ESV', '["Pray", "Petition", "Request", "Peace"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Matthew 18:20', 'ESV', '["Come", "Together", "Jesus Christ", "Name"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Hebrews 10:24-25', 'ESV', '["One Another", "Love", "Good Deeds", "Encourage"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Matthew 4:19', 'ESV', '["Follow", "Fishers", "Men", "Jesus Christ"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('fa45fcbd-4cdc-4f65-9d6c-c618a00c7f1c', 'Romans 1:16', 'ESV', '["Gospel", "Power", "Lord God", "Salvation"]');
-- Pack: Proclaiming Christ
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Proclaiming Christ', 'B', 'TMS 60', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Romans 3:23', 'ESV', '["Sin", "Fall", "Glory", "Lord God"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Isaiah 53:6', 'ESV', '["Astray", "Own Way", "Iniquity", "Jesus Christ"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Romans 6:23', 'ESV', '["Sin", "Death", "Gift", "Life"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Hebrews 9:27', 'ESV', '["Death", "Judgement", "Man", "Destiny"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Romans 5:8', 'ESV', '["Lord God", "Love", "Sinner", "Christ Jesus"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', '1 Peter 3:18', 'ESV', '["Christ Jesus", "Died", "Sins", "Reconcile"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Ephesians 2:8-9', 'ESV', '["Grace", "Salvation", "Faith", "Works"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Titus 3:5', 'ESV', '["Salvation", "Righteous", "Mercy", "Holy Spirit"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'John 1:12', 'ESV', '["Received", "Believe", "Children", "Lord God"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'Revelation 3:20', 'ESV', '["Here", "Knock", "Hears", "Opens"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', '1 John 5:13', 'ESV', '["Believe", "Christ Jesus", "Eternal Life", "Know"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('414a8407-c785-486d-ab34-3de3458be753', 'John 5:24', 'ESV', '["Hears", "Believe", "Eternal Life", "Life"]');
-- Pack: Reliance on God's Resources
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Reliance on God''s Resources', 'C', 'TMS 60', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', '1 Corinthians 3:16', 'ESV', '["Temple", "Lord God", "Spirit", "Within"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', '1 Corinthians 2:12', 'ESV', '["World", "Spirit", "Lord God", "Freely"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Isaiah 41:10', 'ESV', '["Fear", "Dismay", "Strengthen", "Help"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Philippians 4:13', 'ESV', '["Do", "Everything", "Strengthen", "Christ"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Lamentations 3:22-23', 'ESV', '["Faithfulness", "Love", "Compassions", "Change"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Numbers 23:19', 'ESV', '["Lie", "Change", "Say", "Act"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Isaiah 26:3', 'ESV', '["Peace", "Lord God", "Steadfast", "Trust"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', '1 Peter 5:7', 'ESV', '["Lord God", "Anxiety", "Care", "Worry"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Romans 8:32', 'ESV', '["Lord God", "Son", "Gracious", "Give"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Philippians 4:19', 'ESV', '["Need", "Riches", "Christ Jesus", "Lord God"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Hebrews 2:18', 'ESV', '["Suffered", "Tempted", "Help", "Temptation"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('f4e176a2-26f9-491c-bcf5-56e5a17b9a6d', 'Psalms 119:9,11', 'ESV', '["Pure", "Word", "Heart", "Sin"]');
-- Pack: Being Christ's Disciple
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Being Christ''s Disciple', 'D', 'TMS 60', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Matthew 6:33', 'ESV', '["Christ Jesus", "First", "Kingdom", "Righteousness"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Luke 9:23', 'ESV', '["Christ Jesus", "Deny", "Cross", "Follow"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', '1 John 2:15-16', 'ESV', '["World", "Lust", "Pride", "Craving"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Romans 12:2', 'ESV', '["World", "Conform", "Transforme", "Renewing"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', '1 Corinthians 15:58', 'ESV', '["Steadfast", "Labor", "Work", "Stand Firm"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Hebrews 12:3', 'ESV', '["Endure", "Opposition", "Weary", "Lose Heart"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Mark 10:45', 'ESV', '["Serve", "Ransom", "Give", "Christ Jesus"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', '2 Corinthians 4:5', 'ESV', '["Serve", "Servants", "Christ Jesus", "Lord"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Proverbs 3:9-10', 'ESV', '["Tithe", "First", "Filled", "Overflow"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', '2 Corinthians 9:6-7', 'ESV', '["Tithe", "Give", "Generous", "Wealth"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Acts 1:8', 'ESV', '["Power", "Holy Spirit", "Witness", "World"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b3f842c8-044b-4f29-9e44-a233fb0288e3', 'Matthew 28:19-20', 'ESV', '["Disciples", "Baptise", "End", "Teaching"]');
-- Pack: Growth in Christlikeness
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Growth in Christlikeness', 'E', 'TMS 60', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'John 13:34-35', 'ESV', '["Jesus Christ", "Love", "Disciples", "One Another"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', '1 John 3:18', 'ESV', '["Love", "Words", "Actions", "Truth"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Philippians 2:3-4', 'ESV', '["Selfish", "Conceit", "Humility", "Interests"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', '1 Peter 5:5-6', 'ESV', '["Submission", "Humility", "Proud", "Humble"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Ephesians 5:3', 'ESV', '["Sexual Immorality", "Impurity", "Greed", "Improper", "Holy"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', '1 Peter 2:11', 'ESV', '["Aliens", "Strangers", "Abstain", "Sin"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Leviticus 19:11', 'ESV', '["Lie", "Cheat", "Steal", "Deception"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Acts 24:16', 'ESV', '["Conscience", "Honesty", "Integrity", "Lord God"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Hebrews 11:6', 'ESV', '["Lord God", "Believe", "Seek", "Faith"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Romans 4:20-21', 'ESV', '["Unbelief", "Faith", "Lord God", "Promise"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Galatians 6:9-10', 'ESV', '["Good", "Give Up", "Believers", "Do Good"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('67028be9-7466-4542-be27-2702cd36418e', 'Matthew 5:16', 'ESV', '["Light", "Shine", "Good", "Praise"]');
-- Pack: Lessons of Assurance
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('877dd6ad-3fe2-47c5-b529-f30986ceb2c4', 'Lessons of Assurance', 'LOA', 'Lessons of Assurance', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('877dd6ad-3fe2-47c5-b529-f30986ceb2c4', '1 John 5:11-12', 'ESV', '["Testimony", "Eternal Life", "Son", "Lord God"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('877dd6ad-3fe2-47c5-b529-f30986ceb2c4', 'John 16:24', 'ESV', '["Jesus''s Name", "Ask", "Joy", "Receive"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('877dd6ad-3fe2-47c5-b529-f30986ceb2c4', '1 Corinthians 10:13', 'ESV', '["Temptation", "Common", "Lord God", "Tempted"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('877dd6ad-3fe2-47c5-b529-f30986ceb2c4', '1 John 1:9', 'ESV', '["Sin", "Confess", "Faithful", "Forgive"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('877dd6ad-3fe2-47c5-b529-f30986ceb2c4', 'Proverbs 3:5-6', 'ESV', '["Trust", "Lord God", "Way", "Path"]');
-- Pack: Life Issues: Sin
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('3afef658-b47b-4eb2-a165-487da1c73b93', 'Life Issues: Sin', 'Sin', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('3afef658-b47b-4eb2-a165-487da1c73b93', 'Romans 6:11-13', 'ESV', '["Dead", "Alive", "Reign", "Obey", "Righteousness"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('3afef658-b47b-4eb2-a165-487da1c73b93', '1 Corinthians 10:13', 'ESV', '["Temptation", "Faithful", "Tempted", "Victory"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('3afef658-b47b-4eb2-a165-487da1c73b93', 'Galatians 6:1-2', 'ESV', '["Transgression", "Gentleness", "Tempted", "Burdens"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('3afef658-b47b-4eb2-a165-487da1c73b93', 'Ephesians 6:10-12', 'ESV', '["Strong", "Armor", "Evil", "Spiritual"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('3afef658-b47b-4eb2-a165-487da1c73b93', 'James 4:7-8', 'ESV', '["Submit", "Resist", "Devil", "Flee", "Cleanse"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('3afef658-b47b-4eb2-a165-487da1c73b93', '1 John 1:8', 'ESV', '["Deceive", "Truth", "Blameless", "Say"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('3afef658-b47b-4eb2-a165-487da1c73b93', '1 John 1:9', 'ESV', '["Confess", "Faithful", "Just", "Forgive", "Cleanse"]');
-- Pack: Life Issues: Depression
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('b2dd451c-9aac-44f9-b3be-f00fa7fb0d9e', 'Life Issues: Depression', 'Depression', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b2dd451c-9aac-44f9-b3be-f00fa7fb0d9e', 'Isaiah 43:1-3', 'ESV', '["Fear", "With", "Overwhelm", "Consume"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b2dd451c-9aac-44f9-b3be-f00fa7fb0d9e', '2 Corinthians 4:7-10', 'ESV', '["Treasure", "Power", "Affliction", "Crushed", "Life"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b2dd451c-9aac-44f9-b3be-f00fa7fb0d9e', 'Psalm 42:5', 'ESV', '["Turmoil", "Soul", "Hope", "Salvation"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b2dd451c-9aac-44f9-b3be-f00fa7fb0d9e', 'Psalm 34:17-18', 'ESV', '["Righteous", "Lord God", "Delivers", "Brokenhearted", "Save"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b2dd451c-9aac-44f9-b3be-f00fa7fb0d9e', 'Lamentations 3:19-23', 'ESV', '["Affliction", "Soul", "Steadfast", "Love", "Lord God"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b2dd451c-9aac-44f9-b3be-f00fa7fb0d9e', '2 Corinthians 1:8-9', 'ESV', '["Affliction", "Burden", "Rely", "Lord God"]');
-- Pack: Life Issues: Guilt
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('1c2374a7-d98a-4217-98b7-63549b2db186', 'Life Issues: Guilt', 'Guilt', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('1c2374a7-d98a-4217-98b7-63549b2db186', 'Psalm 32:1-2', 'ESV', '["Blessing", "Sin", "Cover", "Spirit"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('1c2374a7-d98a-4217-98b7-63549b2db186', 'Psalm 51:9-10', 'ESV', '["Hide", "Sin", "Clean", "Heart", "Spirit"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('1c2374a7-d98a-4217-98b7-63549b2db186', 'Proverbs 28:13', 'ESV', '["Conceal", "Transgression", "Confess", "Forsake"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('1c2374a7-d98a-4217-98b7-63549b2db186', 'Romans 8:1-2', 'ESV', '["Condemn", "Christ Jesus", "Law", "Free"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('1c2374a7-d98a-4217-98b7-63549b2db186', '2 Corinthians 7:10', 'ESV', '["Grief", "Repentance", "Salvation", "Death"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('1c2374a7-d98a-4217-98b7-63549b2db186', 'James 5:16', 'ESV', '["Confess", "Sin", "Prayer", "Righteous"]');
-- Pack: Life Issues: God's Will
                INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('14b071d1-f95d-4bff-bda3-0052dfd066b2', 'Life Issues: God''s Will', 'God''s Will', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('14b071d1-f95d-4bff-bda3-0052dfd066b2', 'Proverbs 3:5-6', 'ESV', '["Trust", "Heart", "Understanding", "Path"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('14b071d1-f95d-4bff-bda3-0052dfd066b2', 'Proverbs 3:7', 'ESV', '["Wisdom", "Own", "Fear", "Turn", "Evil"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('14b071d1-f95d-4bff-bda3-0052dfd066b2', 'Proverbs 16:9', 'ESV', '["Heart", "Man", "Plans", "Step"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('14b071d1-f95d-4bff-bda3-0052dfd066b2', 'Isaiah 30:21', 'ESV', '["Hear", "Word", "Walk", "Turn"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('14b071d1-f95d-4bff-bda3-0052dfd066b2', 'Jeremiah 29:11-13', 'ESV', '["Plans", "Prosper", "Harm", "Call", "Hear"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('14b071d1-f95d-4bff-bda3-0052dfd066b2', 'Romans 12:1-2', 'ESV', '["Mercy", "Offer", "Sacrifice", "Transform"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('14b071d1-f95d-4bff-bda3-0052dfd066b2', '1 John 5:14-15', 'ESV', '["Confidence", "Ask", "Hear", "Know"]');
-- Pack: Life Issues: Love
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('b4bd6b61-4316-47ec-9141-110f9cd3bafa', 'Life Issues: Love', 'Love', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b4bd6b61-4316-47ec-9141-110f9cd3bafa', 'Matthew 22:37-40', 'ESV', '["Lord God", "Heart", "Soul", "Mind", "Yourself"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b4bd6b61-4316-47ec-9141-110f9cd3bafa', 'John 13:34-35', 'ESV', '["Command", "One Another", "Disciples", "Jesus Christ"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b4bd6b61-4316-47ec-9141-110f9cd3bafa', 'Romans 8:38-39', 'ESV', '["Death", "Life", "Separate", "Jesus Christ"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b4bd6b61-4316-47ec-9141-110f9cd3bafa', '1 Corinthians 13:1-3', 'ESV', '["Tongues", "Prophecy", "Knowledge", "Faith"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b4bd6b61-4316-47ec-9141-110f9cd3bafa', '1 Corinthians 13:4-8', 'ESV', '["Patient", "Kind", "Envy", "Boasting", "Arrogance", "Rude", "Truth", "Irritation", "Resentment", "Hope", "Endures"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('b4bd6b61-4316-47ec-9141-110f9cd3bafa', '1 John 4:20', 'ESV', '["Hate", "Brother", "Liar", "Lord God"]');
-- Pack: Life Issues: Money
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('eef60076-3c8c-497d-be26-55a97b945dcc', 'Life Issues: Money', 'Money', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('eef60076-3c8c-497d-be26-55a97b945dcc', 'Deuteronomy 8:17-18', 'ESV', '["Heart", "Power", "Lord God", "Give"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('eef60076-3c8c-497d-be26-55a97b945dcc', 'Proverbs 3:9-10', 'ESV', '["Lord God", "First", "Produce", "Plenty"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('eef60076-3c8c-497d-be26-55a97b945dcc', 'Matthew 6:19-21', 'ESV', '["Hoard", "Earth", "Heaven", "Heart"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('eef60076-3c8c-497d-be26-55a97b945dcc', 'Matthew 6:24', 'ESV', '["Master", "Hate", "Love", "Lord God"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('eef60076-3c8c-497d-be26-55a97b945dcc', 'Philippians 4:11-13', 'ESV', '["Need", "Content", "Plenty", "Jesus Christ"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('eef60076-3c8c-497d-be26-55a97b945dcc', '1 Timothy 6:9-10', 'ESV', '["Desire", "Temptation", "Destruction", "Evil"]');
-- Pack: Life Issues: Perfectionism
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('4dec3464-29ee-479b-985e-b7b818f08641', 'Life Issues: Perfectionism', 'Perfectionism', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('4dec3464-29ee-479b-985e-b7b818f08641', 'Psalm 127:1-2', 'ESV', '["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('4dec3464-29ee-479b-985e-b7b818f08641', 'Ecclesiastes 2:10-11', 'ESV', '["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('4dec3464-29ee-479b-985e-b7b818f08641', 'Luke 10:40-42', 'ESV', '["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('4dec3464-29ee-479b-985e-b7b818f08641', '2 Corinthians 12:9', 'ESV', '["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('4dec3464-29ee-479b-985e-b7b818f08641', 'Galatians 3:3', 'ESV', '["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('4dec3464-29ee-479b-985e-b7b818f08641', 'Ephesians 2:8-9', 'ESV', '["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]');
-- Pack: Life Issues: Self-Image
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('89b54275-8d04-48ec-84b4-34ea89a31e38', 'Life Issues: Self-Image', 'Self-Image', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('89b54275-8d04-48ec-84b4-34ea89a31e38', '1 Samuel 16:7', 'ESV', '["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('89b54275-8d04-48ec-84b4-34ea89a31e38', 'Psalm 139:13-14', 'ESV', '["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('89b54275-8d04-48ec-84b4-34ea89a31e38', 'Jeremiah 9:23-24', 'ESV', '["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('89b54275-8d04-48ec-84b4-34ea89a31e38', 'Matthew 10:29-31', 'ESV', '["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('89b54275-8d04-48ec-84b4-34ea89a31e38', 'Philippians 2:3-11', 'ESV', '["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('89b54275-8d04-48ec-84b4-34ea89a31e38', '1 Peter 3:3-4', 'ESV', '["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]');
-- Pack: Life Issues: Sex
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('8a4c1bab-16db-40ca-ab33-2e1a9aa5b670', 'Life Issues: Sex', 'Sex', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('8a4c1bab-16db-40ca-ab33-2e1a9aa5b670', 'Matthew 5:27-28', 'ESV', '["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('8a4c1bab-16db-40ca-ab33-2e1a9aa5b670', 'Romans 13:13-14', 'ESV', '["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('8a4c1bab-16db-40ca-ab33-2e1a9aa5b670', '1 Corinthians 6:18-20', 'ESV', '["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('8a4c1bab-16db-40ca-ab33-2e1a9aa5b670', 'Ephesians 5:3', 'ESV', '["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('8a4c1bab-16db-40ca-ab33-2e1a9aa5b670', '1 Thessalonians 4:3-5', 'ESV', '["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('8a4c1bab-16db-40ca-ab33-2e1a9aa5b670', 'Hebrews 13:4', 'ESV', '["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]');
-- Pack: Life Issues: Stress
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('dc3a2152-2d71-4b2e-926a-ad64914c9dfd', 'Life Issues: Stress', 'Stress', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('dc3a2152-2d71-4b2e-926a-ad64914c9dfd', 'Psalm 73:26', 'ESV', '["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('dc3a2152-2d71-4b2e-926a-ad64914c9dfd', 'Psalm 118:5-6', 'ESV', '["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('dc3a2152-2d71-4b2e-926a-ad64914c9dfd', 'Matthew 11:28-30', 'ESV', '["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('dc3a2152-2d71-4b2e-926a-ad64914c9dfd', '2 Corinthians 4:16-18', 'ESV', '["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('dc3a2152-2d71-4b2e-926a-ad64914c9dfd', 'Philippians 4:6-7', 'ESV', '["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('dc3a2152-2d71-4b2e-926a-ad64914c9dfd', '1 Peter 5:5-7', 'ESV', '["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]');
-- Pack: Life Issues: Suffering
    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES ('0de97ba4-93f7-4743-a7c7-b4f7c2e73148', 'Life Issues: Suffering', 'Suffering', 'Life Issues', true) ;
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('0de97ba4-93f7-4743-a7c7-b4f7c2e73148', 'Romans 5:2-5', 'ESV', '["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('0de97ba4-93f7-4743-a7c7-b4f7c2e73148', '2 Corinthians 1:3-4', 'ESV', '["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('0de97ba4-93f7-4743-a7c7-b4f7c2e73148', 'James 1:2-4', 'ESV', '["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('0de97ba4-93f7-4743-a7c7-b4f7c2e73148', 'James 1:12', 'ESV', '["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('0de97ba4-93f7-4743-a7c7-b4f7c2e73148', '1 Peter 1:6-7', 'ESV', '["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]');
INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES ('0de97ba4-93f7-4743-a7c7-b4f7c2e73148', '1 Peter 4:12-13', 'ESV', '["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]');