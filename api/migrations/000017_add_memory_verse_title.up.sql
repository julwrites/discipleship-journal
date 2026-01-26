-- Up Migration
ALTER TABLE memory_verses ADD COLUMN title VARCHAR(255);
DO $$
BEGIN

    UPDATE memory_verses mv
    SET title = 'Christ the Center'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 5 : 17', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Christ the Center'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Galatians 2 : 20', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Obedience to Christ'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 12 : 1', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Obedience to Christ'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('John 14 : 21', ' ', '');

    UPDATE memory_verses mv
    SET title = 'The Word'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Timothy 3 : 16', ' ', '');

    UPDATE memory_verses mv
    SET title = 'The Word'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Joshua 1 : 8', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Prayer'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('John 15 : 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Prayer'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Philippians 4 : 6 - 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Fellowship'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 18 : 20', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Fellowship'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Hebrews 10 : 24 - 25', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Witnessing'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 4 : 19', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Witnessing'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'A'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 1 : 16', ' ', '');

    UPDATE memory_verses mv
    SET title = 'All have Sinned'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 3 : 23', ' ', '');

    UPDATE memory_verses mv
    SET title = 'All have Sinned'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Isaiah 53 : 6', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin''s Penalty'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 6 : 23', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin''s Penalty'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Hebrews 9 : 27', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Christ paid the Penalty'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 5 : 8', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Christ paid the Penalty'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Peter 3 : 18', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Salvation not by Works'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Ephesians 2 : 8 - 9', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Salvation not by Works'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Titus 3 : 5', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Must receive Christ'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('John 1 : 12', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Must receive Christ'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Revelation 3 : 20', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Assurance of Salvation'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 5 : 13', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Assurance of Salvation'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'B'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('John 5 : 24', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Spirit'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Corinthians 3 : 16', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Spirit'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Corinthians 2 : 12', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Strength'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Isaiah 41 : 10', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Strength'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Philippians 4 : 13', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Faithfulness'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Lamentations 3 : 22 - 23', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Faithfulness'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Numbers 23 : 19', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Peace'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Isaiah 26 : 3', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Peace'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Peter 5 : 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Provision'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 8 : 32', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Provision'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Philippians 4 : 19', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Help in Temptation'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Hebrews 2 : 18', ' ', '');

    UPDATE memory_verses mv
    SET title = 'His Help in Temptation'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'C'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalms 119 : 9, 11', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Put Christ First'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 6 : 33', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Put Christ First'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Luke 9 : 23', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Separate from the World'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 2 : 15 - 16', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Separate from the World'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 12 : 2', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Be Steadfast'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Corinthians 15 : 58', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Be Steadfast'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Hebrews 12 : 3', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Serve Others'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Mark 10 : 45', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Serve Others'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 4 : 5', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Give Generously'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Proverbs 3 : 9 - 10', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Give Generously'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 9 : 6 - 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Develop World Vision'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Acts 1 : 8', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Develop World Vision'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'D'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 28 : 19 - 20', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Love'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('John 13 : 34 - 35', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Love'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 3 : 18', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Humility'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Philippians 2 : 3 - 4', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Humility'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Peter 5 : 5 - 6', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Purity'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Ephesians 5 : 3', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Purity'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Peter 2 : 11', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Honesty'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Leviticus 19 : 11', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Honesty'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Acts 24 : 16', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Faith'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Hebrews 11 : 6', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Faith'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 4 : 20 - 21', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Good Works'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Galatians 6 : 9 - 10', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Good Works'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'E'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 5 : 16', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Assurance of Salvation'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'LOA'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 5 : 11 - 12', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Assurance of Answered Prayer'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'LOA'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('John 16 : 24', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Assurance of Victory'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'LOA'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Corinthians 10 : 13', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Assurance of Forgiveness'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'LOA'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 1 : 9', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Assurance of Guidance'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'LOA'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Proverbs 3 : 5 - 6', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sin'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 6 : 11 - 13', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sin'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Corinthians 10 : 13', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sin'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Galatians 6 : 1 - 2', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sin'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Ephesians 6 : 10 - 12', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sin'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('James 4 : 7 - 8', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sin'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 1 : 8', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sin'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sin'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 1 : 9', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Depression'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Depression'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Isaiah 43 : 1 - 3', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Depression'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Depression'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 4 : 7 - 10', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Depression'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Depression'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalm 42 : 5', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Depression'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Depression'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalm 34 : 17 - 18', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Depression'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Depression'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Lamentations 3 : 19 - 23', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Depression'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Depression'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 1 : 8 - 9', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Guilt'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Guilt'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalm 32 : 1 - 2', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Guilt'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Guilt'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalm 51 : 9 - 10', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Guilt'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Guilt'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Proverbs 28 : 13', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Guilt'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Guilt'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 8 : 1 - 2', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Guilt'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Guilt'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 7 : 10', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Guilt'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Guilt'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('James 5 : 16', ' ', '');

    UPDATE memory_verses mv
    SET title = 'God''s Will'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'God''s Will'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Proverbs 3 : 5 - 6', ' ', '');

    UPDATE memory_verses mv
    SET title = 'God''s Will'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'God''s Will'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Proverbs 3 : 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'God''s Will'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'God''s Will'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Proverbs 16 : 9', ' ', '');

    UPDATE memory_verses mv
    SET title = 'God''s Will'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'God''s Will'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Isaiah 30 : 21', ' ', '');

    UPDATE memory_verses mv
    SET title = 'God''s Will'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'God''s Will'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Jeremiah 29 : 11 - 13', ' ', '');

    UPDATE memory_verses mv
    SET title = 'God''s Will'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'God''s Will'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 12 : 1 - 2', ' ', '');

    UPDATE memory_verses mv
    SET title = 'God''s Will'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'God''s Will'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 5 : 14 - 15', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Love'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Love'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 22 : 37 - 40', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Love'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Love'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('John 13 : 34 - 35', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Love'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Love'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 8 : 38 - 39', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Love'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Love'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Corinthians 13 : 1 - 3', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Love'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Love'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Corinthians 13 : 4 - 8', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Love'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Love'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 John 4 : 20', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Money'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Money'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Deuteronomy 8 : 17 - 18', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Money'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Money'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Proverbs 3 : 9 - 10', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Money'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Money'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 6 : 19 - 21', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Money'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Money'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 6 : 24', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Money'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Money'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Philippians 4 : 11 - 13', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Money'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Money'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Timothy 6 : 9 - 10', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Perfectionism'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Perfectionism'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalm 127 : 1 - 2', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Perfectionism'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Perfectionism'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Ecclesiastes 2 : 10 - 11', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Perfectionism'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Perfectionism'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Luke 10 : 40 - 42', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Perfectionism'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Perfectionism'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 12 : 9', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Perfectionism'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Perfectionism'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Galatians 3 : 3', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Perfectionism'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Perfectionism'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Ephesians 2 : 8 - 9', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Self-Image'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Self-Image'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Samuel 16 : 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Self-Image'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Self-Image'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalm 139 : 13 - 14', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Self-Image'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Self-Image'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Jeremiah 9 : 23 - 24', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Self-Image'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Self-Image'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 10 : 29 - 31', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Self-Image'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Self-Image'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Philippians 2 : 3 - 11', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Self-Image'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Self-Image'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Peter 3 : 3 - 4', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sex'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sex'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 5 : 27 - 28', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sex'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sex'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 13 : 13 - 14', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sex'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sex'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Corinthians 6 : 18 - 20', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sex'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sex'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Ephesians 5 : 3', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sex'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sex'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Thessalonians 4 : 3 - 5', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Sex'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Sex'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Hebrews 13 : 4', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Stress'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Stress'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalm 73 : 26', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Stress'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Stress'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Psalm 118 : 5 - 6', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Stress'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Stress'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Matthew 11 : 28 - 30', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Stress'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Stress'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 4 : 16 - 18', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Stress'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Stress'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Philippians 4 : 6 - 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Stress'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Stress'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Peter 5 : 5 - 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Suffering'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Suffering'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('Romans 5 : 2 - 5', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Suffering'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Suffering'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('2 Corinthians 1 : 3 - 4', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Suffering'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Suffering'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('James 1 : 2 - 4', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Suffering'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Suffering'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('James 1 : 12', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Suffering'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Suffering'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Peter 1 : 6 - 7', ' ', '');

    UPDATE memory_verses mv
    SET title = 'Suffering'
    FROM verse_packs vp
    WHERE mv.verse_pack_id = vp.id
      AND vp.identifier = 'Suffering'
      AND REPLACE(mv.reference, ' ', '') = REPLACE('1 Peter 4 : 12 - 13', ' ', '');

END $$;
