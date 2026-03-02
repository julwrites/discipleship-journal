CREATE TABLE memory_verses (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE, -- Nullable for system packs
    pack_name VARCHAR(100), -- E.g., "TMS: Living the New Life" or "My Verses"
    reference VARCHAR(100) NOT NULL,
    text TEXT NOT NULL,
    version VARCHAR(20) NOT NULL DEFAULT 'ESV',
    tags JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_memory_verses_user_id ON memory_verses(user_id);
CREATE INDEX idx_memory_verses_pack_name ON memory_verses(pack_name);


-- Seed TMS Data (Topical Memory System) - A - Live the New Life
INSERT INTO memory_verses (pack_name, reference, text, version, tags) VALUES
('TMS: Live the New Life', '2 Corinthians 5:17', 'Therefore, if anyone is in Christ, he is a new creation. The old has passed away; behold, the new has come.', 'ESV', '["Christ the Center", "Assurance"]'),
('TMS: Live the New Life', 'Galatians 2:20', 'I have been crucified with Christ. It is no longer I who live, but Christ who lives in me. And the life I now live in the flesh I live by faith in the Son of God, who loved me and gave himself for me.', 'ESV', '["Christ the Center", "Identity"]'),
('TMS: Live the New Life', 'Romans 12:1', 'I appeal to you therefore, brothers, by the mercies of God, to present your bodies as a living sacrifice, holy and acceptable to God, which is your spiritual worship.', 'ESV', '["Obedience to Christ", "Worship"]'),
('TMS: Live the New Life', 'John 14:21', 'Whoever has my commandments and keeps them, he it is who loves me. And he who loves me will be loved by my Father, and I will love him and manifest myself to him.', 'ESV', '["Obedience to Christ", "Love"]'),
('TMS: Live the New Life', '2 Timothy 3:16', 'All Scripture is breathed out by God and profitable for teaching, for reproof, for correction, and for training in righteousness,', 'ESV', '["The Word", "Scripture"]'),
('TMS: Live the New Life', 'Joshua 1:8', 'This Book of the Law shall not depart from your mouth, but you shall meditate on it day and night, so that you may be careful to do according to all that is written in it. For then you will make your way prosperous, and then you will have good success.', 'ESV', '["The Word", "Success"]'),
('TMS: Live the New Life', 'John 15:7', 'If you abide in me, and my words abide in you, ask whatever you wish, and it will be done for you.', 'ESV', '["Prayer", "Abiding"]'),
('TMS: Live the New Life', 'Philippians 4:6-7', 'Do not be anxious about anything, but in everything by prayer and supplication with thanksgiving let your requests be made known to God. And the peace of God, which surpasses all understanding, will guard your hearts and your minds in Christ Jesus.', 'ESV', '["Prayer", "Peace"]'),
('TMS: Live the New Life', 'Matthew 18:20', 'For where two or three are gathered in my name, there am I among them.', 'ESV', '["Fellowship", "Presence"]'),
('TMS: Live the New Life', 'Hebrews 10:24-25', 'And let us consider how to stir up one another to love and good works, not neglecting to meet together, as is the habit of some, but encouraging one another, and all the more as you see the Day drawing near.', 'ESV', '["Fellowship", "Community"]'),
('TMS: Live the New Life', 'Matthew 4:19', 'And he said to them, "Follow me, and I will make you fishers of men."', 'ESV', '["Witnessing", "Evangelism"]'),
('TMS: Live the New Life', 'Romans 1:16', 'For I am not ashamed of the gospel, for it is the power of God for salvation to everyone who believes, to the Jew first and also to the Greek.', 'ESV', '["Witnessing", "Gospel"]');
