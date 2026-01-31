DO $$
DECLARE
    pack_id UUID;
BEGIN

    -- Pack: DEP 242 Pack 1 - Assurance of Salvation
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 1 - Assurance of Salvation', 'DEP-242-Pack-1-Assurance-of-Salvation', 'DEP 242 Pack 1 - Assurance of Salvation', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 13:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:11-12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 6:47', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 3:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 3:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:16-17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 10:28-29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:39', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:13', 'ESV', '[]');

    -- Pack: DEP 242 Pack 2 - Quiet Time
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 2 - Quiet Time', 'DEP-242-Pack-2-Quiet-Time', 'DEP 242 Pack 2 - Quiet Time', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 30:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 27:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 34:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Habakkuk 2:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 143:8, 10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 42:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 42:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 130:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 55:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 68:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 91:9-10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 11:28-29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:147-148', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 5:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 95:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 13:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 90:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 107:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 1:35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Exodus 33:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Daniel 6:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 4:34', 'ESV', '[]');

    -- Pack: DEP 242 Pack 3 - The Word
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 3 - The Word', 'DEP-242-Pack-3-The-Word', 'DEP 242 Pack 3 - The Word', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 3:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 24:35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:24-25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Samuel 7:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 23:29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 24:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 20:32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:105', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 6:22-23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 107:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 8:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 15:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:111', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 17:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:97', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ezra 7:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 23:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 5:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 11:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 1:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Deuteronomy 17:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 17:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Deuteronomy 6:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 7:1-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 1:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Joshua 1:8', 'ESV', '[]');

    -- Pack: 180 Series 2 - Growing in Love
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 2 - Growing in Love', '180-Series-2-Growing-in-Love', '180 Series 2 - Growing in Love', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 17:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 11:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 4:4-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 15:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 5:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:23-24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 18:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 9:8-9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 18:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:24-25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:31', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:20-21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 3:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 20:26-27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ecclesiastes 4:9-10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 9:36', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 3:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 6:1', 'ESV', '[]');

    -- Pack: 180 Series 1 - Getting to Know God
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 1 - Getting to Know God', '180-Series-1-Getting-to-Know-God', '180 Series 1 - Getting to Know God', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:1,14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 1:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 2:52', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 1:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 19:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:18-19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 4:16-17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 3:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:13-14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 12:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 4:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:9-10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 1:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 12:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 12:4-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 32:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 11:33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 23:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 9:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Chronicles 29:11-13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 3:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 4:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:15-16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 145:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 86:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:28', 'ESV', '[]');

    -- Pack: 180 Series 3 - Growing in Faith
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 3 - Growing in Faith', '180-Series-3-Growing-in-Faith', '180 Series 3 - Growing in Faith', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 1:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 1:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 3:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 1:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 5:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 12:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 37:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 2:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 103:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:10-11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:20-21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 15:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 23:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 17:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:39-40', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 8:31-32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 24:35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:6-7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:2-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:38', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 6:11-12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 6:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 2:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 2:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:1', 'ESV', '[]');

    -- Pack: DEP 242 - Pack 4
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 - Pack 4', 'DEP-242-Pack-4', 'DEP 242 - Pack 4', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:17 ', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 4:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 4:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 2:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:13-14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 33:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 34:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 50:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 4:31', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 4:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 21:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:14-15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 66:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 3:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 18:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 29:12-13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Nehemiah 1:8-9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Daniel 6:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 5:15-16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 6:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 5:17-18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 16:25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Chronicles 29:11-13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 2:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:6-7', 'ESV', '[]');

    -- Pack: DEP 242 Pack 5 - Fellowship
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 5 - Fellowship', 'DEP-242-Pack-5-Fellowship', 'DEP 242 Pack 5 - Fellowship', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 2:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 13:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 18:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 133:1-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 3:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ecclesiastes 4:9-10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:17,19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 13:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 2:42,47', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 1:5,27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:24-25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 8:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 4:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 6:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 6:12-13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 9:34-35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:23-24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 18:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 28:13', 'ESV', '[]');

    -- Pack: DEP 242 Pack 6 - Witnessing
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 6 - Witnessing', 'DEP-242-Pack-6-Witnessing', 'DEP 242 Pack 6 - Witnessing', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 1:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:18-19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 4:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 19:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 9:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 2:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jonah 4:10-11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 15:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Daniel 12:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:15-16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 4:39', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:23-24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 1:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:9-10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 6:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 13:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 17:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 20:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 8:29-30', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 4:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 1:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 59:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 9:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 1:8-9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 21:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 2:8-9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 2:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 4:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:22-23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 2:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 3:6-7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 6:63', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 3:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 2:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 3:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:9-10', 'ESV', '[]');

    -- Pack: 180 Series 4 - Walking in Victory
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 4 - Walking in Victory', '180-Series-4-Walking-in-Victory', '180 Series 4 - Walking in Victory', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:57', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 2:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:10-11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 12:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 4:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 13:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 37:31', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:12-13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 1:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 6:45', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 4:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 4:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 12:36-37', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 5:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 21:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 1:35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 33:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:14-15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:17-18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Samuel 12:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 9:37-38', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 13:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 146:1-2', 'ESV', '[]');

    -- Pack: 180 Series 5 - Sharing Christ with Others
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 5 - Sharing Christ with Others', '180-Series-5-Sharing-Christ-with-Others', '180 Series 5 - Sharing Christ with Others', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:27-28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 9:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 4:35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:10-12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 1:8-9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 1:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:9-10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 10:28-29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 21:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 8:36', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 7:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 5:31-32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:44', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 7:25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 14:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:39', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 14:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 25:41', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 1:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:18-19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:14-15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 3:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:30', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 3:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 2:9-10', 'ESV', '[]');

    -- Pack: DEP 242 Pack 7 - The Lordship of Christ
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 7 - The Lordship of Christ', 'DEP-242-Pack-7-The-Lordship-of-Christ', 'DEP 242 Pack 7 - The Lordship of Christ', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:16-17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:10-11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 14:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 10:29-30', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 22:16-17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Chronicles 28:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Chronicles 16:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 91:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Samuel 2:30', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 9:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 16:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 3:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:22-24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 6:7-9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:15-16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 1:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 1:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:21-22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:36-38', 'ESV', '[]');

    -- Pack: DEP 242 Pack 8 - World Vision
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 8 - World Vision', 'DEP-242-Pack-8-World-Vision', 'DEP 242 Pack 8 - World Vision', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 12:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 49:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 28:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 1:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 6:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 2:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 57:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 1:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 4:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:28-29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 1:9-11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 60:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Chronicles 16:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Habakkuk 2:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Malachi 1:11', 'ESV', '[]');

    -- Pack: Promises Pack
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Promises Pack', 'Promises-Pack', 'Promises Pack', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 33:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 9:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:37', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 22:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 11:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 32:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 48:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:13-14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 103:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 21:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 17:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 112:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 6:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 121:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 46:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 91:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:38-39', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 40:31', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 3:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 16:25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 4:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 138:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Zechariah 14:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 8:31-32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 17:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 11:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 21:3-4', 'ESV', '[]');

    -- Pack: Self-Respect Pack
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Self-Respect Pack', 'Self-Respect-Pack', 'Self-Respect Pack', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 149:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:6-7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 2:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 10:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 14:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:17-18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 1:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 3:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 16:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 42:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 63:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 73:25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 94:18-19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 116:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 12:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 3:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 2:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 18:35-36', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 56:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 57:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Habakkuk 3:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 1:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 6:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:15', 'ESV', '[]');

    -- Pack: God's Word
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'God''s Word', 'God-s-Word', 'God''s Word', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 24:35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:24-25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 138:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 23:29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 24:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:46-47', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 20:32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 33:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ezra 7:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 15:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:30', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:34', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:45', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:56', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:60', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:71', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:80', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:81', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:93', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:97', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:105', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:114', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:125', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:130', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:143', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:148', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:160', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:165', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:173', 'ESV', '[]');

    -- Pack: Established For Growth
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Established For Growth', 'Established-For-Growth', 'Established For Growth', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 3:16-17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 10:28-29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 2:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:25', 'ESV', '[]');

    -- Pack: Walking in Victory
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Walking in Victory', 'Walking-in-Victory', 'Walking in Victory', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:57', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 2:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:10,11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 12:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 4:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 13:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 37:31', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:12-13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 1:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 6:45', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 4:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 4:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 12:36-37', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 5:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 21:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 1:35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 33:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:14-15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:17,18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Samuel 12:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 9:37-38', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 13:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 146:1-2', 'ESV', '[]');

    -- Pack: Empowered For Multiplication
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Empowered For Multiplication', 'Empowered-For-Multiplication', 'Empowered For Multiplication', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 4:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 3:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 24:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 4:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:2-3', 'ESV', '[]');

    -- Pack: Equipped For Ministry
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Equipped For Ministry', 'Equipped-For-Ministry', 'Equipped For Ministry', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 145:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:6-7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:2', 'ESV', '[]');

    -- Pack: Hide His word in your heart
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Hide His word in your heart', 'Hide-His-word-in-your-heart', 'Hide His word in your heart', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:37', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 1:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:3-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 3:12-14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:14-16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:1-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 8:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 4:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 1:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 19:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 1:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Lamentations 3:22-23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:6-7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:34', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 90:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 13:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 28:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 3:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 4:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 5:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 46:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 16:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 139:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 19:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 46:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 100:1-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Habakkuk 3:17-19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Nehemiah 9:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 16:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:20-21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 42:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 139:1-18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 2:6-7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 4:10-11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 1:1-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 112:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 5:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 3:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 43:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 27:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 2:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 1:15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 143:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:105', 'ESV', '[]');

    -- Pack: Come Follow Me - Answering the Call
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Come Follow Me - Answering the Call', 'Come-Follow-Me-Answering-the-Call', 'Come Follow Me - Answering the Call', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 11:28-30', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 6:68-69', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 10:37-38', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 9:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 14:33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 9:62', 'ESV', '[]');

    -- Pack: Disciple-making
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Disciple-making', 'Disciple-making', 'Disciple-making', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:14-15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:12-13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 4:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 78:72', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 3:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ezra 7:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:20-21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 4:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 14:26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 14:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 14:33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 8:31-32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 13:34-35', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 3:10-11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 20:24', 'ESV', '[]');

    -- Pack: Sufferings Pack
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Sufferings Pack', 'Sufferings-Pack', 'Sufferings Pack', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 6:22-23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:20-21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:18-19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:6-7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 1:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:25 ', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 1:20-22 ', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 5:40-41', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 7:59-60', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 1:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:67,71', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 4:17', 'ESV', '[]');

    -- Pack: The Pursuit of Holiness
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'The Pursuit of Holiness', 'The-Pursuit-of-Holiness', 'The Pursuit of Holiness', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 2:11-12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 39:8-9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 4:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:13-14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 6:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 24:3-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 3:8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:22-24', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:5-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 4:7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 51:10-13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 2:15-16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:20-21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 3:21', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 13:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:29-33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 1:27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Lamentations 3:22-23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 2:11-12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:14-15', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:13-14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:25-26', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 7:1-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:9, 11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 55:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 3:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:22', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:3-4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 97:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:44', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 28:28', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 2:1-2', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Micah 7:8', 'ESV', '[]');

    -- Pack: Excuses Pack
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Excuses Pack', 'Excuses-Pack', 'Excuses Pack', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 14:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:31', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 8:36', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:33', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:130', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 4:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 3:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 5:32', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 4:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 19:10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 1:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 29:25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:11-12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 3:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ecclesiastes 11:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 14:12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 13:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:8,9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:44', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 8:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 43:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 25:41', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:13-14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 14:1', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 1:20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 3:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:21', 'ESV', '[]');

    -- Pack: The Spirit-Filled Life
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'The Spirit-Filled Life', 'The-Spirit-Filled-Life', 'The Spirit-Filled Life', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 44:3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 7:37-38', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:45', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 33:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:9', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:15-18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:19-20', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ezekiel 36:26-27', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 2:29', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:22-23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:16-19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:23', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 15:31', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 2:13', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:11', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:18', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 3:5-7', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:11-12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:16', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 4:31', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:4', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 4:6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:15-17', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:2-3', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:4-6', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:30', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:16-19', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 7:51', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 4:4-5', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 3:7-8', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:25', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:13-14', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 4:7-10', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 51:10-12', 'ESV', '[]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:18-21', 'ESV', '[]');
END; $$;