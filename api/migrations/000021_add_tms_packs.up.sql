DO $$
DECLARE
    pack_id UUID;
BEGIN

    -- Pack: DEP 242 Pack 1 - Assurance of Salvation
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 1 - Assurance of Salvation', 'DEP-242-Pack-1-Assurance-of-Salvation', 'DEP 242 Pack 1 - Assurance of Salvation', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 13:5', 'ESV', '["\u2014 Unless", "See Whether", "Christ Jesus", "Test", "Realize", "Faith"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:11-12', 'ESV', '["Whoever", "Testimony", "Son", "Life", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:13', 'ESV', '["May Know", "Eternal Life", "Write", "Things", "Son", "Name"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 6:47', 'ESV', '["Eternal Life", "Truly", "Tell", "One", "Believes"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:7', 'ESV', '["Sins", "Riches", "Redemption", "Grace", "God", "Forgiveness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:1', 'ESV', '["Christ Jesus", "Therefore", "Condemnation"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:24', 'ESV', '["Justified Freely", "Christ Jesus", "Redemption", "Grace", "Came"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:1', 'ESV', '["Lord Jesus Christ", "Therefore", "Since", "Peace", "Justified", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:3', 'ESV', '["Lord Jesus Christ", "Jesus Christ", "Living Hope", "Great Mercy", "Resurrection", "Praise"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 3:5', 'ESV', '["Saved Us", "Righteous Things", "Holy Spirit", "Washing", "Renewal", "Rebirth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 3:26', 'ESV', '["Christ Jesus", "God", "Faith", "Children"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:14', 'ESV', '["Spirit", "Led", "God", "Children"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:9', 'ESV', '["God Lives", "Spirit", "Realm", "Indeed", "However", "Flesh"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:16-17', 'ESV', '["World Cannot Accept", "Neither Sees", "Forever \u2014", "Another Advocate", "Truth", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 10:28-29', 'ESV', '["Shall Never Perish", "Eternal Life", "Snatch", "One", "Hand", "Greater"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:39', 'ESV', '["Separate Us", "Neither Height", "Christ Jesus", "Anything Else", "Love", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:23', 'ESV', '["Perishable Seed", "Enduring Word", "Living", "Imperishable", "God", "Born"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:13', 'ESV', '["Promised Holy Spirit", "Truth", "Seal", "Salvation", "Message", "Marked"]');

    -- Pack: DEP 242 Pack 2 - Quiet Time
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 2 - Quiet Time', 'DEP-242-Pack-2-Quiet-Time', 'DEP 242 Pack 2 - Quiet Time', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:9', 'ESV', '["Jesus Christ", "Son", "Lord", "God", "Fellowship", "Faithful"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 30:18', 'ESV', '["Lord Longs", "Lord", "Yet", "Wait", "Therefore", "Show"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:6', 'ESV', '["Seek", "Near", "May", "Lord", "Found", "Call"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 27:8', 'ESV', '["Heart Says", "Face ", "Face", "Seek", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:5', 'ESV', '["Bear Much Fruit", "Vine", "Remain", "Nothing", "Branches", "Apart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 34:10', 'ESV', '["Lord Lack", "Good Thing", "Seek", "Hungry"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Habakkuk 2:1', 'ESV', '["Watch", "Station", "Stand", "See", "Say", "Ramparts"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 143:8, 10', 'ESV', '["Good Spirit Lead", "Unfailing Love", "Morning Bring", "Level Ground", "Word", "Way"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:2', 'ESV', '["Right Hand", "Joy Set", "Throne", "Shame", "Scorning", "Sat"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 42:11', 'ESV', '["Yet Praise", "Disturbed Within", "Soul", "Savior", "Put", "Hope"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 42:1', 'ESV', '["Soul Pants", "Deer Pants", "Water", "Streams", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 130:5-6', 'ESV', '["Watchmen Wait", "Wait", "Word", "Whole", "Waits", "Put"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 55:22', 'ESV', '["Never Let", "Sustain", "Shaken", "Righteous", "Lord", "Cast"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 68:19', 'ESV', '["Daily Bears", "Savior", "Praise", "Lord", "God", "Burdens"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 91:9-10', 'ESV', '["Refuge ", "Come Near", "Tent", "Say", "Overtake", "Make"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 11:28-29', 'ESV', '["Yoke Upon", "Find Rest", "Rest", "Weary", "Take", "Souls"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:147-148', 'ESV', '["Eyes Stay Open", "May Meditate", "Word", "Watches", "Rise", "Put"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 5:3', 'ESV', '["Wait Expectantly", "Voice", "Requests", "Morning", "Lord", "Lay"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 95:6', 'ESV', '["Let Us Kneel", "Let Us Bow", "Worship", "Maker", "Lord", "Come"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 13:15', 'ESV', '["Praise \u2014", "Openly Profess", "Therefore", "Sacrifice", "Name", "Lips"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 90:14', 'ESV', '["Unfailing Love", "Satisfy Us", "May Sing", "Morning", "Joy", "Glad"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 107:9', 'ESV', '["Good Things", "Thirsty", "Satisfies", "Hungry", "Fills"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 1:35', 'ESV', '["Still Dark", "Solitary Place", "Jesus Got", "Went", "Prayed", "Morning"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Exodus 33:11', 'ESV', '["Lord Would Speak", "Moses Would Return", "One Speaks", "Moses Face", "Face", "Tent"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Daniel 6:10', 'ESV', '["Went Home", "Upstairs Room", "Three Times", "Giving Thanks", "Daniel Learned", "Published"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 4:34', 'ESV', '["Without Using", "Say Anything", "Explained Everything", "Parable", "Disciples", "Alone"]');

    -- Pack: DEP 242 Pack 3 - The Word
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 3 - The Word', 'DEP-242-Pack-3-The-Word', 'DEP 242 Pack 3 - The Word', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 3:16', 'ESV', '["Useful", "Training", "Teaching", "Scripture", "Righteousness", "Rebuking"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:21', 'ESV', '["Prophecy Never", "Holy Spirit", "Carried Along", "Though Human", "Human", "Spoke"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 24:35', 'ESV', '["Never Pass Away", "Pass Away", "Words", "Heaven", "Earth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:24-25', 'ESV', '["Grass Withers", "Like Grass", "Flowers Fall", "Like", "Flowers", "Word"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:17', 'ESV', '["Word", "Truth", "Sanctify"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Samuel 7:28', 'ESV', '["Sovereign Lord", "Good Things", "Trustworthy", "Servant", "Promised", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 23:29', 'ESV', '["Rock", "Pieces", "Lord", "Like", "Hammer", "Breaks"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:9', 'ESV', '["Suffering Even", "Chained Like", "Chained", "Word", "Point", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:4', 'ESV', '["Man Shall", "Jesus Answered", "God ", "Every Word", "Bread Alone", "Written"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 24:27', 'ESV', '["Scriptures Concerning", "Said", "Prophets", "Moses", "Explained", "Beginning"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:23', 'ESV', '["Perishable Seed", "Enduring Word", "Living", "Imperishable", "God", "Born"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:18', 'ESV', '["Give Us Birth", "Word", "Truth", "Might", "Kind", "Firstfruits"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:2', 'ESV', '["Like Newborn Babies", "May Grow", "Salvation"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 20:32', 'ESV', '["Inheritance Among", "Word", "Sanctified", "Grace", "God", "Give"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:105', 'ESV', '["Word", "Path", "Light", "Lamp", "Feet"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 6:22-23', 'ESV', '["Way", "Watch", "Walk", "Teaching", "Speak", "Sleep"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 107:20', 'ESV', '["Word", "Sent", "Rescued", "Healed", "Grave"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 8:8', 'ESV', '["Centurion Replied", "Word", "Servant", "Say", "Roof", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 15:16', 'ESV', '["Lord God Almighty", "Words Came", "Name", "Joy", "Heart", "Delight"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:111', 'ESV', '["Heritage Forever", "Statutes", "Joy", "Heart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:17', 'ESV', '["Word", "Take", "Sword", "Spirit", "Salvation", "Helmet"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '["Penetrates Even", "Edged Sword", "Dividing Soul", "Word", "Thoughts", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 17:11', 'ESV', '["Scriptures Every Day", "Paul Said", "Noble Character", "Great Eagerness", "Berean Jews", "True"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:97', 'ESV', '["Day Long", "Meditate", "Love", "Law"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ezra 7:10', 'ESV', '["Teaching", "Study", "Observance", "Lord", "Laws", "Law"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 23:12', 'ESV', '["Daily Bread", "Words", "Treasured", "Mouth", "Lips", "Departed"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 5:5-6', 'ESV', '["Worked Hard", "Simon Answered", "Nets Began", "Nets ", "Large Number", "Caught Anything"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:17', 'ESV', '["Faith Comes", "Word", "Message", "Hearing", "Heard", "Consequently"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 11:28', 'ESV', '["Blessed Rather", "Word", "Replied", "Obey", "Hear", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 1:3', 'ESV', '["Reads Aloud", "Written", "Words", "Time", "Take", "Prophecy"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Deuteronomy 17:19', 'ESV', '["May Learn", "Follow Carefully", "Words", "Revere", "Read", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 17:11', 'ESV', '["Scriptures Every Day", "Paul Said", "Noble Character", "Great Eagerness", "Berean Jews", "True"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:15', 'ESV', '["One Approved", "Correctly Handles", "Worker", "Word", "Truth", "Present"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Deuteronomy 6:6', 'ESV', '["Today", "Hearts", "Give", "Commandments"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 7:1-3', 'ESV', '["Commands Within", "Commands", "Write", "Words", "Teachings", "Tablet"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 1:1-2', 'ESV', '["Whose Delight", "Sinners Take", "Law Day", "Law", "Wicked", "Way"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Joshua 1:8', 'ESV', '["Law Always", "Everything Written", "Successful", "Prosperous", "Night", "Meditate"]');

    -- Pack: 180 Series 2 - Growing in Love
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 2 - Growing in Love', '180-Series-2-Growing-in-Love', '180 Series 2 - Growing in Love', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:15', 'ESV', '["Mature Body", "Every Respect", "Truth", "Speaking", "Love", "Instead"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:9', 'ESV', '["Old Self", "Taken", "Since", "Practices", "Lie"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 17:9', 'ESV', '["Whoever Repeats", "Offense", "Whoever", "Would", "Foster", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 11:13', 'ESV', '["Trustworthy Person Keeps", "Gossip Betrays", "Secret", "Confidence"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 4:4-6', 'ESV', '["Act Toward Outsiders", "May Proclaim", "May Know", "Every Opportunity", "Answer Everyone", "Always Full"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 15:1', 'ESV', '["Harsh Word Stirs", "Anger", "Gentle", "Answer", "Turns", "Away"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 5:16', 'ESV', '["Therefore Confess", "Righteous Person", "Sins", "Prayer", "Pray", "Powerful"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:23-24', 'ESV', '["First Go", "Therefore", "Something", "Sister", "Remember", "Reconciled"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:19', 'ESV', '["Take Note", "Dear Brothers", "Become Angry", "Speak", "Slow", "Sisters"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 18:13', 'ESV', '["Listening \u2014", "Shame", "Folly", "Answer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 9:8-9', 'ESV', '["Wiser Still", "Rebuke Mockers", "Rebuke", "Wise", "Teach", "Righteous"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 18:15', 'ESV', '["Sister Sins", "Two", "Point", "Listen", "Fault", "Brother"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:32', 'ESV', '["Christ God Forgave", "One Another", "Kind", "Forgiving", "Compassionate"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:13', 'ESV', '["Forgive One Another", "Lord Forgave", "Forgive", "Someone", "Grievance", "Bear"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:2', 'ESV', '["One Another", "Completely Humble", "Patient", "Love", "Gentle", "Bearing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:24-25', 'ESV', '["Repentance Leading", "Gently Instructed", "Servant Must", "Opponents Must", "Must", "Truth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:26', 'ESV', '["Sun Go", "Still Angry", "Sin ", "Let", "Anger"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:8', 'ESV', '["Must Also Rid", "Filthy Language", "Things", "Slander", "Rage", "Malice"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:15', 'ESV', '["One Falls Short", "Bitter Root Grows", "Defile Many", "Cause Trouble", "See", "Grace"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:31', 'ESV', '["Get Rid", "Every Form", "Slander", "Rage", "Malice", "Brawling"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:20-21', 'ESV', '["Christ Suffered", "Wrong", "Suffer", "Steps", "Receive", "Leaving"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:19', 'ESV', '["Repay ,\" Says", "Take Revenge", "Leave Room", "Dear Friends", "Written", "Wrath"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:4', 'ESV', '["Fury Overwhelming", "Stand", "Jealousy", "Cruel", "Anger"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 3:16', 'ESV', '["Every Evil Practice", "Selfish Ambition", "Find Disorder", "Envy"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:5-6', 'ESV', '["Lord Jesus Christ", "Christ Jesus", "One Voice", "One Mind", "Mind Toward", "Gives Endurance"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:10', 'ESV', '["Lord Jesus Christ", "Perfectly United", "One Another", "Divisions Among", "Thought", "Sisters"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 20:26-27', 'ESV', '["Become Great Among", "Whoever Wants", "Slave \u2014", "First Must", "Must", "Servant"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:13', 'ESV', '["Use", "Sisters", "Rather", "Love", "Indulge", "Freedom"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:2', 'ESV', '["Please", "Neighbors", "Good", "Build"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:3-4', 'ESV', '["Humility Value Others", "Vain Conceit", "Selfish Ambition", "Others", "Rather", "Nothing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:11', 'ESV', '["Fact", "Build", "Therefore", "Encourage", "Another", "Each"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ecclesiastes 4:9-10', 'ESV', '["Pity Anyone", "Good Return", "Two", "One", "Labor", "Help"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 9:36', 'ESV', '["Like Sheep Without", "Shepherd", "Saw", "Helpless", "Harassed", "Crowds"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:15', 'ESV', '["Rejoice", "Mourn", "With", "Those"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 3:17', 'ESV', '["Good Fruit", "Wisdom", "Submissive", "Sincere", "Pure", "Peace"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 6:1', 'ESV', '["Person Gently", "Also May", "Watch", "Tempted", "Spirit", "Someone"]');

    -- Pack: 180 Series 1 - Getting to Know God
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 1 - Getting to Know God', '180-Series-1-Getting-to-Know-God', '180 Series 1 - Getting to Know God', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:1,14', 'ESV', '["Dwelling Among Us", "Word Became Flesh", "Word", "Truth", "Son", "Seen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 1:8', 'ESV', '["Throne", "Son", "Scepter", "Says", "Last", "Kingdom"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:15', 'ESV', '["\u2014 Yet", "High Priest", "Every Way", "Weaknesses", "Unable", "Tempted"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 2:52', 'ESV', '["Jesus Grew", "Wisdom", "Stature", "Man", "God", "Favor"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:3-4', 'ESV', '["Third Day According", "Sins According", "First Importance", "Christ Died", "Scriptures", "Received"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:20', 'ESV', '["Fallen Asleep", "Raised", "Indeed", "Firstfruits", "Dead", "Christ"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:18', 'ESV', '["Ever Seen God", "Closest Relationship", "God", "Son", "One", "Made"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 1:3', 'ESV', '["Right Hand", "Provided Purification", "Powerful Word", "Exact Representation", "Things", "Sustaining"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 19:10', 'ESV', '["Man Came", "Lost ", "Son", "Seek", "Save"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:18-19', 'ESV', '["Lamb Without Blemish", "Precious Blood", "Perishable Things", "Life Handed", "Empty Way", "Silver"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 4:16-17', 'ESV', '["Trumpet Call", "Still Alive", "Rise First", "Loud Command", "Lord Forever", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 3:2-3', 'ESV', '["Made Known", "Dear Friends", "Christ Appears", "Shall See", "Shall", "Yet"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:13-14', 'ESV', '["Make Known", "Yet", "Truth", "Tell", "Spirit", "Speak"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 12:3', 'ESV', '["Lord ,\" Except", "Cursed ", "God Says", "Holy Spirit", "Spirit", "Want"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:9', 'ESV', '["God Lives", "Spirit", "Realm", "Indeed", "However", "Flesh"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 4:6', 'ESV', '["God Sent", "Father ", "Spirit", "Sons", "Son", "Hearts"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:18', 'ESV', '["Get Drunk", "Wine", "Spirit", "Leads", "Instead", "Filled"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:16', 'ESV', '["Walk", "Spirit", "Say", "Gratify", "Flesh", "Desires"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:9-10', 'ESV', '["Human Mind", "Conceived \"\u2014", "Deep Things", "Spirit Searches", "Things God", "Things"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:26', 'ESV', '["Holy Spirit", "Things", "Teach", "Send", "Said", "Remind"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:4-5', 'ESV', '["Persuasive Words", "Human Wisdom", "Faith Might", "Wise", "Spirit", "Rest"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 1:5', 'ESV', '["Lived Among", "Holy Spirit", "Gospel Came", "Deep Conviction", "Words", "Simply"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 12:11', 'ESV', '["Work", "Spirit", "One", "Distributes", "Determines"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 12:4-6', 'ESV', '["Spirit Distributes", "Different Kinds", "Working", "Work", "Service", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 32:17', 'ESV', '["Sovereign Lord", "Outstretched Arm", "Great Power", "Nothing", "Made", "Heavens"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 11:33', 'ESV', '["Paths Beyond Tracing", "Wisdom", "Unsearchable", "Riches", "Knowledge", "Judgments"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 23:24', 'ESV', '["Secret Places", "Fill Heaven", "Cannot See", "Earth ", "Lord", "Hide"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 9:8', 'ESV', '["Every Good Work", "Times", "Things", "Need", "God", "Bless"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Chronicles 29:11-13', 'ESV', '["Honor Come", "Glorious Name", "Give Strength", "Strength", "Give", "Wealth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 3:3', 'ESV', '["Evil One", "Strengthen", "Protect", "Lord", "Faithful"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 4:24', 'ESV', '["Worshipers Must Worship", "Truth ", "Spirit", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:15-16', 'ESV', '["Holy ", "Holy", "Written", "Called"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 145:3', 'ESV', '["Worthy", "Praise", "One", "Lord", "Greatness", "Great"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:10', 'ESV', '["Loved Us", "Loved God", "Atoning Sacrifice", "Son", "Sins", "Sent"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 86:15', 'ESV', '["Gracious God", "Slow", "Love", "Lord", "Faithfulness", "Compassionate"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:28', 'ESV', '["Things God Works", "Called According", "Purpose", "Love", "Know", "Good"]');

    -- Pack: 180 Series 3 - Growing in Faith
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 3 - Growing in Faith', '180-Series-3-Growing-in-Faith', '180 Series 3 - Growing in Faith', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:3-4', 'ESV', '["Given Us Everything", "Given Us", "Called Us", "World Caused", "Precious Promises", "May Participate"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 1:20', 'ESV', '["Many Promises God", "God", "Yes", "Spoken", "Matter", "Made"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 1:3-4', 'ESV', '["Lord Jesus Christ", "Comforts Us", "Troubles", "Trouble", "Receive", "Praise"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 3:7-8', 'ESV', '["Makes Things Grow", "Rewarded According", "One Purpose", "One", "Waters", "Plants"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 1:2-3', 'ESV', '["Wither \u2014 Whatever", "Whose Leaf", "Whose Delight", "Tree Planted", "Law Day", "Law"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:8', 'ESV', '["Lord Jesus Christ", "Increasing Measure", "Unproductive", "Qualities", "Possess", "Knowledge"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 5:10', 'ESV', '["Eternal Glory", "Suffered", "Strong", "Steadfast", "Restore", "Make"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 12:9', 'ESV', '["Power May Rest", "Weakness ", "Made Perfect", "Power", "Weaknesses", "Therefore"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:3', 'ESV', '["Every Spiritual Blessing", "Lord Jesus Christ", "Heavenly Realms", "Blessed Us", "Christ", "Praise"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 37:4-5', 'ESV', '["Take Delight", "Way", "Trust", "Lord", "Heart", "Give"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 2:1-2', 'ESV', '["Whole World", "Righteous One", "Dear Children", "Atoning Sacrifice", "Write", "Sins"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 103:12', 'ESV', '["West", "Transgressions", "Removed", "Far", "East"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '["Penetrates Even", "Edged Sword", "Dividing Soul", "Word", "Thoughts", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:10-11', 'ESV', '["Yields Seed", "Without Watering", "Snow Come", "Word", "Sower", "Sent"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:20-21', 'ESV', '["Scripture Came", "Must Understand", "Holy Spirit", "Carried Along", "Though Human", "Prophecy Never"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:13', 'ESV', '["Human Word", "Word", "Work", "Received", "Indeed", "Heard"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 15:16', 'ESV', '["Lord God Almighty", "Words Came", "Name", "Joy", "Heart", "Delight"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 23:12', 'ESV', '["Daily Bread", "Words", "Treasured", "Mouth", "Lips", "Departed"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 17:11', 'ESV', '["Scriptures Every Day", "Paul Said", "Noble Character", "Great Eagerness", "Berean Jews", "True"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:39-40', 'ESV', '["Scriptures Diligently", "Eternal Life", "Scriptures", "Life", "Yet", "Think"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 8:31-32', 'ESV', '["Jesus Said", "Free ", "Truth", "Teaching", "Set", "Really"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:4', 'ESV', '["Man Shall", "Jesus Answered", "God ", "Every Word", "Bread Alone", "Written"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 24:35', 'ESV', '["Never Pass Away", "Pass Away", "Words", "Heaven", "Earth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:17', 'ESV', '["Word", "Truth", "Sanctify"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:6-7', 'ESV', '["Suffer Grief", "Proven Genuineness", "Jesus Christ", "Greatly Rejoice", "Greater Worth", "Faith \u2014"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:2-4', 'ESV', '["Let Perseverance Finish", "Faith Produces Perseverance", "Pure Joy", "Many Kinds", "Lacking Anything", "Face Trials"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:2', 'ESV', '["Good News Proclaimed", "Value", "Share", "Obeyed", "Message", "Heard"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:38', 'ESV', '["Shrinks Back ", "Righteous One", "One", "Take", "Pleasure", "Live"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:16', 'ESV', '["Flaming Arrows", "Evil One", "Take", "Shield", "Faith", "Extinguish"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 6:11-12', 'ESV', '["Take Hold", "Pursue Righteousness", "Many Witnesses", "Good Confession", "Eternal Life", "Good Fight"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:1', 'ESV', '["See", "Hope", "Faith", "Confidence", "Assurance"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:17', 'ESV', '["Faith Comes", "Word", "Message", "Hearing", "Heard", "Consequently"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 6:12', 'ESV', '["Patience Inherit", "Become Lazy", "Want", "Promised", "Imitate", "Faith"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 2:17', 'ESV', '["Way", "Faith", "Dead", "Action", "Accompanied"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 2:16', 'ESV', '["Jesus Christ", "Christ Jesus", "Christ", "Works", "Put", "Person"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:1', 'ESV', '["Lord Jesus Christ", "Therefore", "Since", "Peace", "Justified", "God"]');

    -- Pack: DEP 242 - Pack 4
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 - Pack 4', 'DEP-242-Pack-4', 'DEP 242 - Pack 4', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:17 ', 'ESV', '["Pray Continually", "Pray", "Continually"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 4:2', 'ESV', '["Watchful", "Thankful", "Prayer", "Devote"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 4:7', 'ESV', '["Sober Mind", "May Pray", "Things", "Therefore", "Near", "End"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 2:1-2', 'ESV', '["May Live Peaceful", "Quiet Lives", "People \u2014", "Urge", "Thanksgiving", "Prayers"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:13-14', 'ESV', '["Father May", "May Ask", "Ask", "Whatever", "Son", "Name"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:20', 'ESV', '["Work Within Us", "Power", "Immeasurably", "Imagine", "Ask", "According"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 33:3', 'ESV', '["Unsearchable Things", "Know ", "Tell", "Great", "Call", "Answer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:5', 'ESV', '["Without Finding Fault", "Lacks Wisdom", "Gives Generously", "Ask God", "Given"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 34:4', 'ESV', '["Sought", "Lord", "Fears", "Delivered", "Answered"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 50:15', 'ESV', '["Trouble", "Honor", "Deliver", "Day", "Call"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 4:31', 'ESV', '["Holy Spirit", "God Boldly", "Word", "Spoke", "Shaken", "Prayed"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 4:3', 'ESV', '["God May Open", "May Proclaim", "Pray", "Mystery", "Message", "Door"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:24', 'ESV', '["Receive", "Name", "Joy", "Complete", "Asked", "Ask"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 21:22', 'ESV', '["Receive Whatever", "Prayer ", "Believe", "Ask"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:18', 'ESV', '["Always Keep", "Spirit", "Requests", "Praying", "Prayers", "Pray"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:14-15', 'ESV', '["Ask Anything According", "Ask \u2014", "Hears Us", "Approaching God", "Know", "Confidence"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 66:18', 'ESV', '["Lord Would", "Cherished Sin", "Listened", "Heart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 3:22', 'ESV', '["Receive", "Pleases", "Keep", "Commands", "Ask", "Anything"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 18:19', 'ESV', '["Earth Agree", "Two", "Truly", "Tell", "Heaven", "Father"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:7-8', 'ESV', '["Seeks Finds", "Asks Receives", "Seek", "Opened", "One", "Knocks"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 29:12-13', 'ESV', '["Seek", "Pray", "Listen", "Heart", "Find", "Come"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Nehemiah 1:8-9', 'ESV', '["Servant Moses", "Name ", "Farthest Horizon", "Exiled People", "Unfaithful", "Scatter"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Daniel 6:10', 'ESV', '["Went Home", "Upstairs Room", "Three Times", "Giving Thanks", "Daniel Learned", "Published"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 5:15-16', 'ESV', '["Jesus Often Withdrew", "People Came", "Lonely Places", "Yet", "Spread", "Sicknesses"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 6:12', 'ESV', '["Days Jesus Went", "Night Praying", "Spent", "Pray", "One", "Mountainside"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 5:17-18', 'ESV', '["Heavens Gave Rain", "Half Years", "Earth Produced", "Prayed Earnestly", "Rain", "Prayed"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 16:25', 'ESV', '["Singing Hymns", "Midnight Paul", "Silas", "Prisoners", "Praying", "Listening"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Chronicles 29:11-13', 'ESV', '["Honor Come", "Glorious Name", "Give Strength", "Strength", "Give", "Wealth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 2:1-2', 'ESV', '["May Live Peaceful", "Quiet Lives", "People \u2014", "Urge", "Thanksgiving", "Prayers"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:18', 'ESV', '["Give Thanks", "Christ Jesus", "God", "Circumstances"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:9', 'ESV', '["Purify Us", "Forgive Us", "Unrighteousness", "Sins", "Faithful", "Confess"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:6-7', 'ESV', '["Every Situation", "Christ Jesus", "Understanding", "Transcends", "Thanksgiving", "Requests"]');

    -- Pack: DEP 242 Pack 5 - Fellowship
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 5 - Fellowship', 'DEP-242-Pack-5-Fellowship', 'DEP 242 Pack 5 - Fellowship', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 2:13', 'ESV', '["Far Away", "Brought Near", "Christ Jesus", "Christ", "Blood"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:20', 'ESV', '["Making Peace", "Whether Things", "Things", "Shed", "Reconcile", "Heaven"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:3', 'ESV', '["Jesus Christ", "Also May", "Son", "Seen", "Proclaim", "Heard"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 13:14', 'ESV', '["Lord Jesus Christ", "Holy Spirit", "May", "Love", "Grace", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 18:20', 'ESV', '["Three Gather", "Two", "Name"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 133:1-3', 'ESV', '["People Live Together", "Even Life Forevermore", "Mount Zion", "Lord Bestows", "Unity", "Running"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 3:13', 'ESV', '["Today ", "Sin", "None", "May", "Long", "Hardened"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ecclesiastes 4:9-10', 'ESV', '["Pity Anyone", "Good Return", "Two", "One", "Labor", "Help"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:22', 'ESV', '["Pursue Righteousness", "Pure Heart", "Evil Desires", "Youth", "Peace", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:17,19', 'ESV', '["Iron Sharpens Iron", "Water Reflects", "Life Reflects", "One", "Heart", "Face"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 13:20', 'ESV', '["Fools Suffers Harm", "Become Wise", "Wise", "Walk", "Companion"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:13', 'ESV', '["Whole Measure", "Reach Unity", "Become Mature", "Son", "Knowledge", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 2:42,47', 'ESV', '["Praising God", "Number Daily", "Lord Added", "Teaching", "Saved", "Prayer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 1:5,27', 'ESV', '["Whatever Happens", "Striving Together", "Stand Firm", "Manner Worthy", "First Day", "One Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:24-25', 'ESV', '["Let Us Consider", "Toward Love", "Meeting Together", "Good Deeds", "Day Approaching", "See"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 8:3-4', 'ESV', '["Urgently Pleaded", "Even Beyond", "Testify", "Sharing", "Service", "Privilege"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 4:13', 'ESV', '["Rejoice Inasmuch", "Sufferings", "Revealed", "Participate", "Overjoyed", "May"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 6:2', 'ESV', '["Way", "Law", "Fulfill", "Christ", "Carry", "Burdens"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:1-2', 'ESV', '["Joy Complete", "Common Sharing", "One Mind", "One", "United", "Therefore"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:3-4', 'ESV', '["Humility Value Others", "Vain Conceit", "Selfish Ambition", "Others", "Rather", "Nothing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 6:12-13', 'ESV', '["Fair Exchange \u2014", "Hearts Also", "Withholding", "Speak", "Affection"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:21', 'ESV', '["One Another", "Submit", "Reverence", "Christ"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 9:34-35', 'ESV', '["Kept Quiet", "Jesus Called", "First Must", "Way", "Wants", "Twelve"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:15', 'ESV', '["One Falls Short", "Bitter Root Grows", "Defile Many", "Cause Trouble", "See", "Grace"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:11', 'ESV', '["Rather Expose", "Fruitless Deeds", "Nothing", "Darkness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:4', 'ESV', '["Soldier Gets Entangled", "Rather Tries", "One Serving", "Commanding Officer", "Civilian Affairs", "Please"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:23-24', 'ESV', '["First Go", "Therefore", "Something", "Sister", "Remember", "Reconciled"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 18:15', 'ESV', '["Sister Sins", "Two", "Point", "Listen", "Fault", "Brother"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:9', 'ESV', '["Purify Us", "Forgive Us", "Unrighteousness", "Sins", "Faithful", "Confess"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 28:13', 'ESV', '["Whoever Conceals", "Finds Mercy", "Sins", "Renounces", "Prosper", "One"]');

    -- Pack: DEP 242 Pack 6 - Witnessing
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 6 - Witnessing', 'DEP-242-Pack-6-Witnessing', 'DEP 242 Pack 6 - Witnessing', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 1:8', 'ESV', '["Holy Spirit Comes", "Receive Power", "Earth ", "Witnesses", "Samaria", "Judea"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:18-19', 'ESV', '["Counting People", "Reconciled Us", "Gave Us", "World", "Sins", "Reconciling"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 4:1-2', 'ESV', '["Great Patience", "Encourage \u2014", "Christ Jesus", "Careful Instruction", "Word", "View"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 19:10', 'ESV', '["Man Came", "Lost ", "Son", "Seek", "Save"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:4', 'ESV', '["Please People", "Trying", "Tests", "Speak", "Hearts", "Gospel"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 9:16', 'ESV', '["Cannot Boast", "Woe", "Since", "Preach", "Gospel", "Compelled"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:14', 'ESV', '["One", "Heard", "Call", "Believed", "Believe"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 2:3-4', 'ESV', '["Pleases God", "Wants", "Truth", "Savior", "Saved", "People"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jonah 4:10-11', 'ESV', '["Twenty Thousand People", "Right Hand", "Lord Said", "Left \u2014", "Great City", "Cannot Tell"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 15:7', 'ESV', '["Nine Righteous Persons", "One Sinner", "Way", "Tell", "Repents", "Repent"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:3', 'ESV', '["Jesus Christ", "Also May", "Son", "Seen", "Proclaim", "Heard"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Daniel 12:3', 'ESV', '["Lead Many", "Shine Like", "Like", "Wise", "Stars", "Righteousness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:19', 'ESV', '[" Jesus Said", "People ", "Send", "Follow", "Fish", "Come"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:15-16', 'ESV', '["May Become Blameless", "God Without Fault", "Crooked Generation ", "Shine Among", "Like Stars", "Hold Firmly"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:19', 'ESV', '["Fearlessly Make Known", "Words May", "Pray Also", "Whenever", "Speak", "Mystery"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:15', 'ESV', '["Hearts Revere Christ", "Respect", "Reason", "Prepared", "Lord", "Hope"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 4:39', 'ESV', '["Town Believed", "Woman", "Told", "Testimony", "Samaritans", "Many"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '["Penetrates Even", "Edged Sword", "Dividing Soul", "Word", "Thoughts", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:11', 'ESV', '["Word", "Sent", "Return", "Purpose", "Mouth", "Goes"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:23-24', 'ESV', '["Preach Christ Crucified", "Stumbling Block", "Christ", "Wisdom", "Power", "Jews"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 1:5', 'ESV', '["Lived Among", "Holy Spirit", "Gospel Came", "Deep Conviction", "Words", "Simply"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:9-10', 'ESV', '["Lord ", "God Raised", "Saved", "Profess", "Mouth", "Justified"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 6:2', 'ESV', '["Time", "Tell", "Says", "Salvation", "Helped", "Heard"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 13:5', 'ESV', '["\u2014 Unless", "See Whether", "Christ Jesus", "Test", "Realize", "Faith"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 17:2-3', 'ESV', '["Three Sabbath Days", "Paul Went", "Messiah ", "Messiah", "Synagogue", "Suffer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 20:24', 'ESV', '["Life Worth Nothing", "Lord Jesus", "Good News", "Testifying", "Task", "Race"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 8:29-30', 'ESV', '["Man Reading Isaiah", "Spirit Told Philip", "Reading ", "Philip Ran", "Philip Asked", "Stay Near"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 4:20', 'ESV', '["Cannot Help Speaking", "Heard ", "Seen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 1:27', 'ESV', '["God Created Mankind", "God", "Created", "Male", "Image", "Female"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 59:1-2', 'ESV', '["Surely", "Sins", "Short", "Separated", "Save", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:12', 'ESV', '["Way Death Came", "Sinned \u2014", "One Man", "Sin Entered", "Death", "Sin"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:23', 'ESV', '["Fall Short", "Sinned", "God", "Glory"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 9:27', 'ESV', '["Face Judgment", "People", "Die", "Destined"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 1:8-9', 'ESV', '["Know God", "Everlasting Destruction", "Lord Jesus", "Lord", "Shut", "Punished"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:23', 'ESV', '["Eternal Life", "Christ Jesus", "Wages", "Sin", "Lord", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 21:8', 'ESV', '["Second Death ", "Practice Magic Arts", "Sexually Immoral", "Liars \u2014", "Fiery Lake", "Burning Sulfur"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 2:8-9', 'ESV', '["God \u2014", "Faith \u2014", "Works", "Saved", "One", "Grace"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 2:21', 'ESV', '["Set Aside", "Righteousness Could", "Nothing ", "Christ Died", "Law", "Grace"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 4:12', 'ESV', '["Saved ", "One Else", "Heaven Given", "Salvation", "Name", "Must"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:18', 'ESV', '["Perishable Things", "Life Handed", "Empty Way", "Silver", "Redeemed", "Know"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:22-23', 'ESV', '["Preach Christ Crucified", "Jews Demand Signs", "Stumbling Block", "Greeks Look", "Jews", "Wisdom"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:21', 'ESV', '["World", "Wisdom", "Since", "Save", "Preached", "Pleased"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 2:8', 'ESV', '["Elemental Spiritual Forces", "World Rather", "One Takes", "Human Tradition", "Deceptive Philosophy", "See"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:13', 'ESV', '["Natural Descent", "Human Decision", "Children Born", "Born", "Husband", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 3:6-7', 'ESV', '["Spirit Gives Birth", "Flesh Gives Birth", "Spirit", "Flesh", "Surprised", "Saying"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 6:63', 'ESV', '["Spirit Gives Life", "Flesh Counts", "Spirit", "Life", "Words", "Spoken"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:18', 'ESV', '["Christ Also Suffered", "Made Alive", "Unrighteous", "Spirit", "Sins", "Righteous"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 3:16', 'ESV', '["Whoever Believes", "Eternal Life", "World", "Son", "Shall", "Perish"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:8', 'ESV', '["Still Sinners", "God Demonstrates", "Christ Died", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:3-4', 'ESV', '["Third Day According", "Sins According", "First Importance", "Christ Died", "Scriptures", "Received"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 2:13', 'ESV', '["God Made", "Forgave Us", "Uncircumcision", "Sins", "Flesh", "Dead"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:24', 'ESV', '["Whoever Hears", "Eternal Life", "Life", "Word", "Truly", "Tell"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:12', 'ESV', '["God \u2014", "Become Children", "Yet", "Right", "Receive", "Name"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 3:20', 'ESV', '["Anyone Hears", "Voice", "Stand", "Person", "Opens", "Knock"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:9-10', 'ESV', '["Lord ", "God Raised", "Saved", "Profess", "Mouth", "Justified"]');

    -- Pack: 180 Series 4 - Walking in Victory
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 4 - Walking in Victory', '180-Series-4-Walking-in-Victory', '180 Series 4 - Walking in Victory', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:57', 'ESV', '["Lord Jesus Christ", "Gives Us", "Victory", "Thanks", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 2:14', 'ESV', '["Always Leads Us", "Uses Us", "Triumphal Procession", "Thanks", "Spread", "Knowledge"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:4-5', 'ESV', '["Every Pretension", "Divine Power", "Demolish Strongholds", "Demolish Arguments", "World", "Weapons"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:10-11', 'ESV', '["Mighty Power", "Full Armor", "Take", "Strong", "Stand", "Schemes"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 12:11', 'ESV', '["Word", "Triumphed", "Testimony", "Shrink", "Much", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 4:7-8', 'ESV', '["Come Near", "Wash", "Submit", "Sinners", "Resist", "Purify"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:5-6', 'ESV', '["Minds Set", "Mind Governed", "Spirit Desires", "Live According", "Flesh Desires", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 13:14', 'ESV', '["Lord Jesus Christ", "Think", "Rather", "Gratify", "Flesh", "Desires"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:4', 'ESV', '["Dear Children", "World", "Overcome", "One", "Greater", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:4-5', 'ESV', '["Everyone Born", "God Overcomes", "Overcomes", "God", "World", "Victory"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 37:31', 'ESV', '["Slip", "Law", "Hearts", "God", "Feet"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:12-13', 'ESV', '["Let Sin Reign", "Offer Every Part", "Rather Offer", "Mortal Body", "Evil Desires", "Sin"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:8', 'ESV', '["Praiseworthy \u2014 Think", "Admirable \u2014", "Whatever", "True", "Things", "Sisters"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 1:15', 'ESV', '["Things", "Pure", "Nothing", "Minds", "Fact", "Corrupted"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 6:45', 'ESV', '["Mouth Speaks", "Good Stored", "Evil Stored", "Heart", "Full"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 4:23', 'ESV', '["Heart", "Guard", "Flows", "Everything", "Else"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:22', 'ESV', '["Whole Body", "Body", "Light", "Lamp", "Healthy", "Full"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:28', 'ESV', '["Already Committed Adultery", "Woman Lustfully", "Tell", "Looks", "Heart", "Anyone"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 4:3', 'ESV', '["Avoid Sexual Immorality", "Sanctified", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:13', 'ESV', '["Sexual Immorality", "Stomach", "Say", "Meant", "Lord", "However"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:29', 'ESV', '["Unwholesome Talk Come", "May Benefit", "Building Others", "Needs", "Mouths", "Listen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 12:36-37', 'ESV', '["Every Empty Word", "Give Account", "Condemned ", "Words", "Tell", "Spoken"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:22', 'ESV', '["Reject Every Kind", "Evil", "Reject", "Every", "Kind"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 5:1-2', 'ESV', '["Treat Younger Men", "Older Man Harshly", "Younger Women", "Older Women", "Absolute Purity", "Sisters"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 21:22', 'ESV', '["Receive Whatever", "Prayer ", "Believe", "Ask"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:7-8', 'ESV', '["Seeks Finds", "Asks Receives", "Seek", "Opened", "One", "Knocks"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:6', 'ESV', '["Unseen", "Sees", "Secret", "Room", "Reward", "Pray"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 1:35', 'ESV', '["Still Dark", "Solitary Place", "Jesus Got", "Went", "Prayed", "Morning"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 33:3', 'ESV', '["Unsearchable Things", "Know ", "Tell", "Great", "Call", "Answer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:20', 'ESV', '["Work Within Us", "Power", "Immeasurably", "Imagine", "Ask", "According"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:14-15', 'ESV', '["Ask Anything According", "Ask \u2014", "Hears Us", "Approaching God", "Know", "Confidence"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:17-18', 'ESV', '["Pray Continually", "Give Thanks", "Christ Jesus", "God", "Circumstances"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Samuel 12:23', 'ESV', '["Way", "Teach", "Sin", "Right", "Pray", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 9:37-38', 'ESV', '["Harvest Field ", "Harvest", "Workers", "Therefore", "Send", "Said"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 13:15', 'ESV', '["Praise \u2014", "Openly Profess", "Therefore", "Sacrifice", "Name", "Lips"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 146:1-2', 'ESV', '["Sing Praise", "Praise", "Soul", "Lord", "Long", "Live"]');

    -- Pack: 180 Series 5 - Sharing Christ with Others
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), '180 Series 5 - Sharing Christ with Others', '180-Series-5-Sharing-Christ-with-Others', '180 Series 5 - Sharing Christ with Others', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:27-28', 'ESV', '["Make Known Among", "Teaching Everyone", "Glorious Riches", "Wisdom", "Proclaim", "One"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:19-20', 'ESV', '["Counting People", "Though God", "Therefore Christ", "God", "Christ", "World"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:4', 'ESV', '["Please People", "Trying", "Tests", "Speak", "Hearts", "Gospel"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 9:19', 'ESV', '["Win", "Though", "Slave", "Possible", "One", "Many"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 4:35', 'ESV', '["Still Four Months", "Harvest ", "Harvest", "Tell", "Saying", "Ripe"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:8', 'ESV', '["Well", "Share", "Much", "Loved", "Lives", "Gospel"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:10-12', 'ESV', '["Together Become Worthless", "Even One ", "Even One", "Turned Away", "Seeks God", "One Righteous"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 1:8-9', 'ESV', '["Know God", "Everlasting Destruction", "Lord Jesus", "Lord", "Shut", "Punished"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:24', 'ESV', '["Might Die", "Healed ", "Wounds", "Sins", "Righteousness", "Live"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 1:9', 'ESV', '["Holy Life \u2014", "Saved Us", "Given Us", "Christ Jesus", "Called Us", "Time"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 10:9-10', 'ESV', '["Lord ", "God Raised", "Saved", "Profess", "Mouth", "Justified"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 10:28-29', 'ESV', '["Shall Never Perish", "Eternal Life", "Snatch", "One", "Hand", "Greater"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 21:2', 'ESV', '["Person May Think", "Lord Weighs", "Ways", "Right", "Heart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 8:36', 'ESV', '["Yet Forfeit", "Whole World", "Soul", "Someone", "Good", "Gain"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 7:17', 'ESV', '["Teaching Comes", "Whether", "Speak", "God", "Find", "Chooses"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 5:31-32', 'ESV', '["Repentance ", "Jesus Answered", "Sinners", "Sick", "Righteous", "Need"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:44', 'ESV', '["One Another", "Believe Since", "Accept Glory", "Glory", "Seek", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 7:25', 'ESV', '["Save Completely", "Always Lives", "Therefore", "Intercede", "God", "Come"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:1', 'ESV', '["Day May Bring", "Tomorrow", "Know", "Boast"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 14:12', 'ESV', '["God", "Give", "Account"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:39', 'ESV', '["Eternal Life", "Scriptures Diligently", "Scriptures", "Think", "Testify", "Study"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 14:12', 'ESV', '["Way", "Right", "Leads", "End", "Death", "Appears"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 25:41', 'ESV', '["Eternal Fire Prepared", "Say", "Left", "Devil", "Depart", "Cursed"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 1:20', 'ESV', '["Invisible Qualities \u2014", "Divine Nature \u2014", "World God", "Without Excuse", "Eternal Power", "Clearly Seen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:18-19', 'ESV', '["Lamb Without Blemish", "Precious Blood", "Perishable Things", "Life Handed", "Empty Way", "Silver"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:18', 'ESV', '["Reconciled Us", "Gave Us", "Reconciliation", "Ministry", "God", "Christ"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:7', 'ESV', '["Sins", "Riches", "Redemption", "Grace", "God", "Forgiveness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:14-15', 'ESV', '["Sin Shall", "Sin", "Shall", "Means", "Master", "Longer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 3:26', 'ESV', '["Christ Jesus", "God", "Faith", "Children"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:9', 'ESV', '["Wonderful Light", "Special Possession", "Royal Priesthood", "May Declare", "Holy Nation", "Chosen People"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:21', 'ESV', '["Might Become", "God Made", "God", "Sin", "Righteousness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:30', 'ESV', '["Us Wisdom", "God \u2014", "Christ Jesus", "Righteousness", "Redemption", "Holiness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:14', 'ESV', '["Made Perfect Forever", "Made Holy", "One Sacrifice"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:24', 'ESV', '["Justified Freely", "Christ Jesus", "Redemption", "Grace", "Came"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 3:20', 'ESV', '["Lord Jesus Christ", "Eagerly Await", "Savior", "Heaven", "Citizenship"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 2:9-10', 'ESV', '["Every Power", "Deity Lives", "Bodily Form", "Head", "Fullness", "Christ"]');

    -- Pack: DEP 242 Pack 7 - The Lordship of Christ
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 7 - The Lordship of Christ', 'DEP-242-Pack-7-The-Lordship-of-Christ', 'DEP 242 Pack 7 - The Lordship of Christ', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:2-3', 'ESV', '["Without", "Things", "Nothing", "Made", "God", "Beginning"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:16-17', 'ESV', '["Things Hold Together", "Whether Thrones", "Things", "Visible", "Rulers", "Powers"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:18', 'ESV', '["Supremacy", "Might", "Head", "Firstborn", "Everything", "Dead"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 1:21', 'ESV', '["Present Age", "Every Name", "Rule", "Power", "One", "Invoked"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:10-11', 'ESV', '["Every Tongue Acknowledge", "Jesus Every Knee", "Jesus Christ", "Name", "Lord", "Heaven"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 14:9', 'ESV', '["Christ Died", "Returned", "Reason", "Might", "Lord", "Living"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:19-20', 'ESV', '["Therefore Honor God", "Holy Spirit", "God", "Temples", "Received", "Price"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:15', 'ESV', '["Longer Live", "Live", "Raised", "Died"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:33', 'ESV', '["Seek First", "Well", "Things", "Righteousness", "Kingdom", "Given"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 10:29-30', 'ESV', '["Come Eternal Life", " Jesus Replied", "Fields \u2014 Along", "Persecutions \u2014", "Left Home", "Hundred Times"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 22:16-17', 'ESV', '["Take Possession", "Surely Bless", "Withheld", "Swear", "Stars", "Son"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Chronicles 28:9', 'ESV', '["Understands Every Desire", "Every Thought", "Willing Mind", "Wholehearted Devotion", "Son Solomon", "Serve"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Chronicles 16:9', 'ESV', '["Lord Range Throughout", "Whose Hearts", "War ", "Fully Committed", "Foolish Thing", "Strengthen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 91:14', 'ESV', '[" Says", "Rescue", "Protect", "Name", "Loves", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Samuel 2:30', 'ESV', '["Family Would Minister", "Forever ", "Lord Declares", "Lord", "Declares", "Therefore"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 9:23', 'ESV', '["Disciple Must Deny", "Whoever Wants", "Cross Daily", "Take", "Said", "Follow"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 16:7', 'ESV', '["Lord Takes Pleasure", "Make Peace", "Way", "Enemies", "Causes", "Anyone"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 3:5-6', 'ESV', '["Ways Submit", "Paths Straight", "Understanding", "Trust", "Make", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:22-24', 'ESV', '["True Righteousness", "Old Self", "New Self", "Made New", "Like God", "Former Way"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:22', 'ESV', '["Pursue Righteousness", "Pure Heart", "Evil Desires", "Youth", "Peace", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 6:7-9', 'ESV', '["Get Rich Fall", "Take Nothing", "Plunge People", "Many Foolish", "Harmful Desires", "Brought Nothing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:15-16', 'ESV', '["Live \u2014", "Every Opportunity", "Wise", "Unwise", "Making", "Evil"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 1:20', 'ESV', '["Without Delay", "Hired Men", "Father Zebedee", "Left", "Followed", "Called"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 1:20', 'ESV', '["Sufficient Courage", "Eagerly Expect", "Always Christ", "Whether", "Way", "Life"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:21-22', 'ESV', '["Jesus Christ", "Everyone Looks", "Work", "Timothy", "Son", "Served"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:36-38', 'ESV', '["Mistreated \u2014", "Faced Jeers", "Even Chains", "Worthy", "World", "Went"]');

    -- Pack: DEP 242 Pack 8 - World Vision
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'DEP 242 Pack 8 - World Vision', 'DEP-242-Pack-8-World-Vision', 'DEP 242 Pack 8 - World Vision', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 12:2-3', 'ESV', '["Whoever Curses", "Name Great", "Great Nation", "Peoples", "Make", "Earth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 49:6', 'ESV', '["Salvation May Reach", "Earth ", "Bring Back", "Also Make", "Tribes", "Thing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 28:19-20', 'ESV', '["Therefore Go", "Obey Everything", "Make Disciples", "Holy Spirit", "Age ", "Teaching"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 1:8', 'ESV', '["Holy Spirit Comes", "Receive Power", "Earth ", "Witnesses", "Samaria", "Judea"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 6:8', 'ESV', '["Us ", "Lord Saying", "Voice", "Shall", "Send", "Said"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:16', 'ESV', '["Gentiles Might Become", "Priestly Duty", "Offering Acceptable", "Holy Spirit", "Christ Jesus", "Gentiles"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 2:8', 'ESV', '["Possession", "Nations", "Make", "Inheritance", "Ends", "Earth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 57:5', 'ESV', '["Let", "Heavens", "God", "Glory", "Exalted", "Earth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 1:2-3', 'ESV', '["Preaching Entrusted", "Eternal Life", "Appointed Season", "Time", "Savior", "Promised"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 4:17', 'ESV', '["Gentiles Might Hear", "Message Might", "Lord Stood", "Fully Proclaimed", "Strength", "Side"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:28-29', 'ESV', '["Teaching Everyone", "Strenuously Contend", "Powerfully Works", "Energy Christ", "Christ", "Wisdom"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 1:9-11', 'ESV', '["Love May Abound", "Jesus Christ \u2014", "May", "Christ", "Righteousness", "Pure"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:2', 'ESV', '["Many Witnesses Entrust", "Teach Others", "Reliable People", "Things", "Say", "Qualified"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 60:22', 'ESV', '["Swiftly ", "Mighty Nation", "Time", "Thousand", "Smallest", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:1', 'ESV', '["Let Us Throw", "Let Us Run", "Race Marked", "Great Cloud", "Easily Entangles", "Witnesses"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Chronicles 16:9', 'ESV', '["Lord Range Throughout", "Whose Hearts", "War ", "Fully Committed", "Foolish Thing", "Strengthen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Habakkuk 2:14', 'ESV', '["Waters Cover", "Sea", "Lord", "Knowledge", "Glory", "Filled"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Malachi 1:11', 'ESV', '["Every Place Incense", "Nations ,\" Says", "Sun Rises", "Pure Offerings", "Lord Almighty", "Great Among"]');

    -- Pack: Promises Pack
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Promises Pack', 'Promises-Pack', 'Promises Pack', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 33:22', 'ESV', '["Save Us", "Lord", "Lawgiver", "King", "Judge"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 9:28', 'ESV', '["Take Away", "Second Time", "Bring Salvation", "Bear Sin", "Waiting", "Sins"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:37', 'ESV', '["Delay ", "Little", "Coming", "Come"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 22:12', 'ESV', '["Person According", "Coming Soon", "Reward", "Look", "Give", "Done"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 11:13', 'ESV', '["Give Good Gifts", "Heaven Give", "Holy Spirit", "Though", "Much", "Know"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:26', 'ESV', '["Father \u2014", "Advocate Comes", "Truth", "Testify", "Spirit", "Send"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 32:8', 'ESV', '["Loving Eye", "Way", "Teach", "Instruct", "Counsel"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 48:17', 'ESV', '["Lord Says \u2014", "Holy One", "Lord", "Way", "Teaches", "Redeemer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:7', 'ESV', '["Seek", "Opened", "Knock", "Given", "Find", "Door"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:13-14', 'ESV', '["Father May", "May Ask", "Ask", "Whatever", "Son", "Name"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 103:2-3', 'ESV', '["Benefits \u2014", "Soul", "Sins", "Praise", "Lord", "Heals"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:14', 'ESV', '["Heavenly Father", "Also Forgive", "Forgive", "Sin", "People"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 21:22', 'ESV', '["Receive Whatever", "Prayer ", "Believe", "Ask"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 17:5-6', 'ESV', '["Sea ", "Mustard Seed", "Mulberry Tree", "Apostles Said", "Faith ", "Faith"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:23', 'ESV', '["Jesus Replied", "Teaching", "Obey", "Make", "Loves", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:10', 'ESV', '["Remain", "Love", "Kept", "Keep", "Father", "Commands"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 112:5', 'ESV', '["Lend Freely", "Justice", "Good", "Generous", "Conduct", "Come"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 6:10', 'ESV', '["Work", "Unjust", "Shown", "People", "Love", "Helped"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 121:7', 'ESV', '["Harm \u2014", "Watch", "Lord", "Life", "Keep"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 46:4', 'ESV', '["Old Age", "Gray Hairs", "Sustain", "Rescue", "Made", "Even"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 91:1', 'ESV', '["Whoever Dwells", "Shelter", "Shadow", "Rest", "High", "Almighty"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:38-39', 'ESV', '["Separate Us", "Christ Jesus", "Anything Else", "Neither Height", "Neither Death", "Neither Angels"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 40:31', 'ESV', '["Wings Like Eagles", "Grow Weary", "Walk", "Strength", "Soar", "Run"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 3:3', 'ESV', '["Evil One", "Strengthen", "Protect", "Lord", "Faithful"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:27', 'ESV', '["World Gives", "Troubled", "Peace", "Let", "Leave", "Hearts"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:33', 'ESV', '["Take Heart", "World ", "World", "Trouble", "Told", "Things"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 16:25', 'ESV', '["Whoever Wants", "Whoever Loses", "Save", "Lose", "Life", "Find"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 4:10', 'ESV', '["Lord", "Lift", "Humble"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 138:8', 'ESV', '["Endures Forever \u2014", "Works", "Vindicate", "Love", "Lord", "Hands"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Zechariah 14:9', 'ESV', '["Whole Earth", "One Lord", "Lord", "Name", "King", "Day"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 8:31-32', 'ESV', '["Jesus Said", "Free ", "Truth", "Teaching", "Set", "Really"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:13', 'ESV', '["Yet", "Truth", "Tell", "Spirit", "Speak", "Hears"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 17:7-8', 'ESV', '["Bear Fruit ", "Whose Confidence", "Tree Planted", "Never Fails", "Heat Comes", "Always Green"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:8', 'ESV', '["Lord Jesus Christ", "Increasing Measure", "Unproductive", "Qualities", "Possess", "Knowledge"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 11:28', 'ESV', '["Weary", "Rest", "Give", "Come", "Burdened"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 21:3-4', 'ESV', '["Wipe Every Tear", "Passed Away ", "Throne Saying", "Old Order", "Loud Voice", "Dwelling Place"]');

    -- Pack: Self-Respect Pack
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Self-Respect Pack', 'Self-Respect-Pack', 'Self-Respect Pack', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 149:4', 'ESV', '["Lord Takes Delight", "Victory", "People", "Humble", "Crowns"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:6-7', 'ESV', '["Five Sparrows Sold", "Many Sparrows", "Two Pennies", "Yet", "Worth", "One"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:24', 'ESV', '["Yet God Feeds", "Valuable", "Storeroom", "Sow", "Reap", "Ravens"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:16', 'ESV', '["Testifies", "Spirit", "God", "Children"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:19-20', 'ESV', '["Therefore Honor God", "Holy Spirit", "God", "Temples", "Received", "Price"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 2:10', 'ESV', '["Good Works", "Christ Jesus", "God Prepared", "God", "Handiwork", "Created"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 10:14', 'ESV', '["Good Shepherd", "Sheep Know", "Sheep", "Know"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 14:8', 'ESV', '["Whether", "Lord", "Live", "Die", "Belong"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:15', 'ESV', '["Longer Live", "Live", "Raised", "Died"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:17-18', 'ESV', '["Lord ", "Boasts Boast", "Lord Commends", "Commends", "One", "Let"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 1:21', 'ESV', '["Live", "Gain", "Die", "Christ"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 3:12', 'ESV', '["Godly Life", "Christ Jesus", "Wants", "Persecuted", "Live", "Fact"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 16:8', 'ESV', '["Right Hand", "Eyes Always", "Shaken", "Lord", "Keep"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 42:8', 'ESV', '["Lord Directs", "Song", "Prayer", "Night", "Love", "Life"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 63:7-8', 'ESV', '["Right Hand Upholds", "Wings", "Sing", "Shadow", "Help", "Cling"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 73:25', 'ESV', '["Desire Besides", "Nothing", "Heaven", "Earth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 94:18-19', 'ESV', '["Unfailing Love", "Slipping ", "Great Within", "Consolation Brought", "Supported", "Said"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 116:1-2', 'ESV', '["Voice", "Turned", "Mercy", "Love", "Lord", "Long"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:22', 'ESV', '["One \u2014", "One", "May", "Glory", "Given", "Gave"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:4-5', 'ESV', '["Form One Body", "One Body", "Though Many", "Member Belongs", "Many Members", "Members"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:5-6', 'ESV', '["Lord Jesus Christ", "Christ Jesus", "One Voice", "One Mind", "Mind Toward", "Gives Endurance"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 12:26', 'ESV', '["Every Part Suffers", "Every Part Rejoices", "One Part Suffers", "One Part", "Honored"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 3:28', 'ESV', '["Neither Slave", "Neither Jew", "Christ Jesus", "One", "Male", "Gentile"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 2:19', 'ESV', '["Longer Foreigners", "Fellow Citizens", "Also Members", "Strangers", "People", "Household"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 18:35-36', 'ESV', '["Right Hand Sustains", "Give Way", "Broad Path", "Saving Help", "Help", "Shield"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 56:4', 'ESV', '["Whose Word", "Praise \u2014", "Mere Mortals", "Trust", "God", "Afraid"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 57:2', 'ESV', '["Vindicates", "High", "God", "Cry"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Habakkuk 3:19', 'ESV', '["Stringed Instruments", "Sovereign Lord", "Feet Like", "Feet", "Tread", "Strength"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:5', 'ESV', '["Competence Comes", "Claim Anything", "God", "Competent"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 1:7', 'ESV', '["Make Us Timid", "Gives Us Power", "Self", "Love", "Discipline"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:3', 'ESV', '["Sober Judgment", "Grace Given", "Faith God", "Every One", "Rather Think", "Think"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 15:7', 'ESV', '["Accept One Another", "Christ Accepted", "Bring Praise", "Order", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 6:2', 'ESV', '["Way", "Law", "Fulfill", "Christ", "Carry", "Burdens"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:25', 'ESV', '["Speak Truthfully", "One Body", "Must Put", "Therefore", "Neighbor", "Members"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:4', 'ESV', '["Lord Always", "Say", "Rejoice"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:15', 'ESV', '["One Body", "Christ Rule", "Thankful", "Since", "Peace", "Members"]');

    -- Pack: God's Word
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'God''s Word', 'God-s-Word', 'God''s Word', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 24:35', 'ESV', '["Never Pass Away", "Pass Away", "Words", "Heaven", "Earth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:24-25', 'ESV', '["Grass Withers", "Like Grass", "Flowers Fall", "Like", "Flowers", "Word"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 138:2', 'ESV', '["Unfailing Love", "Solemn Decree", "Holy Temple", "Toward", "Surpasses", "Praise"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:17', 'ESV', '["Word", "Truth", "Sanctify"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 23:29', 'ESV', '["Rock", "Pieces", "Lord", "Like", "Hammer", "Breaks"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '["Penetrates Even", "Edged Sword", "Dividing Soul", "Word", "Thoughts", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 24:27', 'ESV', '["Scriptures Concerning", "Said", "Prophets", "Moses", "Explained", "Beginning"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:46-47', 'ESV', '["Say ", "Believed Moses", "Would Believe", "Believe", "Wrote", "Since"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:4', 'ESV', '["Man Shall", "Jesus Answered", "God ", "Every Word", "Bread Alone", "Written"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 20:32', 'ESV', '["Inheritance Among", "Word", "Sanctified", "Grace", "God", "Give"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 33:4', 'ESV', '["Word", "True", "Right", "Lord", "Faithful"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:21', 'ESV', '["Word Planted", "Moral Filth", "Humbly Accept", "Get Rid", "Therefore", "Save"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ezra 7:10', 'ESV', '["Teaching", "Study", "Observance", "Lord", "Laws", "Law"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 15:16', 'ESV', '["Lord God Almighty", "Words Came", "Name", "Joy", "Heart", "Delight"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:7', 'ESV', '["Upright Heart", "Righteous Laws", "Praise", "Learn"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:14', 'ESV', '["One Rejoices", "Great Riches", "Statutes", "Rejoice", "Following"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:18', 'ESV', '["Open", "Law", "Eyes"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:30', 'ESV', '["Way", "Set", "Laws", "Heart", "Faithfulness", "Chosen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:34', 'ESV', '["May Keep", "Understanding", "Obey", "Law", "Heart", "Give"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:45', 'ESV', '["Walk", "Sought", "Precepts", "Freedom"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:56', 'ESV', '["Precepts", "Practice", "Obey"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:60', 'ESV', '["Obey", "Hasten", "Delay", "Commands"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:71', 'ESV', '["Might Learn", "Good", "Decrees", "Afflicted"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:80', 'ESV', '["Wholeheartedly Follow", "Shame", "Put", "May", "Decrees"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:81', 'ESV', '["Soul Faints", "Word", "Salvation", "Put", "Longing", "Hope"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:93', 'ESV', '["Never Forget", "Preserved", "Precepts", "Life"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:97', 'ESV', '["Day Long", "Meditate", "Love", "Law"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:105', 'ESV', '["Word", "Path", "Light", "Lamp", "Feet"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:114', 'ESV', '["Word", "Shield", "Refuge", "Put", "Hope"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:125', 'ESV', '["May Understand", "Statutes", "Servant", "Give", "Discernment"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:130', 'ESV', '["Words Gives Light", "Gives Understanding", "Unfolding", "Simple"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:143', 'ESV', '["Commands Give", "Come Upon", "Trouble", "Distress", "Delight"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:148', 'ESV', '["Eyes Stay Open", "May Meditate", "Watches", "Promises", "Night"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:160', 'ESV', '["Righteous Laws", "Words", "True", "Eternal"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:165', 'ESV', '["Great Peace", "Stumble", "Nothing", "Make", "Love", "Law"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:173', 'ESV', '["Ready", "Precepts", "May", "Help", "Hand", "Chosen"]');

    -- Pack: Established For Growth
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Established For Growth', 'Established-For-Growth', 'Established For Growth', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:23', 'ESV', '["Eternal Life", "Christ Jesus", "Wages", "Sin", "Lord", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:6', 'ESV', '["Without Faith", "Please God", "Must Believe", "Earnestly Seek", "Rewards", "Impossible"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 3:16-17', 'ESV', '["Every Good Work", "Thoroughly Equipped", "God May", "God", "Useful", "Training"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 10:28-29', 'ESV', '["Shall Never Perish", "Eternal Life", "Snatch", "One", "Hand", "Greater"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 2:20', 'ESV', '["Longer Live", "Christ Lives", "Live", "Christ", "Son", "Loved"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:25', 'ESV', '["Let Us Keep", "Step", "Spirit", "Since", "Live"]');

    -- Pack: Walking in Victory
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Walking in Victory', 'Walking-in-Victory', 'Walking in Victory', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:57', 'ESV', '["Lord Jesus Christ", "Gives Us", "Victory", "Thanks", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 2:14', 'ESV', '["Always Leads Us", "Uses Us", "Triumphal Procession", "Thanks", "Spread", "Knowledge"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:4-5', 'ESV', '["Every Pretension", "Divine Power", "Demolish Strongholds", "Demolish Arguments", "World", "Weapons"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:10,11', 'ESV', '["Mighty Power", "Full Armor", "Take", "Strong", "Stand", "Schemes"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Revelation 12:11', 'ESV', '["Word", "Triumphed", "Testimony", "Shrink", "Much", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 4:7-8', 'ESV', '["Come Near", "Wash", "Submit", "Sinners", "Resist", "Purify"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:5-6', 'ESV', '["Minds Set", "Mind Governed", "Spirit Desires", "Live According", "Flesh Desires", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 13:14', 'ESV', '["Lord Jesus Christ", "Think", "Rather", "Gratify", "Flesh", "Desires"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:4', 'ESV', '["Dear Children", "World", "Overcome", "One", "Greater", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:4-5', 'ESV', '["Everyone Born", "God Overcomes", "Overcomes", "God", "World", "Victory"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 37:31', 'ESV', '["Slip", "Law", "Hearts", "God", "Feet"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:12-13', 'ESV', '["Let Sin Reign", "Offer Every Part", "Rather Offer", "Mortal Body", "Evil Desires", "Sin"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:8', 'ESV', '["Praiseworthy \u2014 Think", "Admirable \u2014", "Whatever", "True", "Things", "Sisters"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 1:15', 'ESV', '["Things", "Pure", "Nothing", "Minds", "Fact", "Corrupted"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 6:45', 'ESV', '["Mouth Speaks", "Good Stored", "Evil Stored", "Heart", "Full"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 4:23', 'ESV', '["Heart", "Guard", "Flows", "Everything", "Else"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:22', 'ESV', '["Whole Body", "Body", "Light", "Lamp", "Healthy", "Full"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:28', 'ESV', '["Already Committed Adultery", "Woman Lustfully", "Tell", "Looks", "Heart", "Anyone"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 4:3', 'ESV', '["Avoid Sexual Immorality", "Sanctified", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:13', 'ESV', '["Sexual Immorality", "Stomach", "Say", "Meant", "Lord", "However"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:29', 'ESV', '["Unwholesome Talk Come", "May Benefit", "Building Others", "Needs", "Mouths", "Listen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 12:36-37', 'ESV', '["Every Empty Word", "Give Account", "Condemned ", "Words", "Tell", "Spoken"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:22', 'ESV', '["Reject Every Kind", "Evil", "Reject", "Every", "Kind"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 5:1-2', 'ESV', '["Treat Younger Men", "Older Man Harshly", "Younger Women", "Older Women", "Absolute Purity", "Sisters"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 21:22', 'ESV', '["Receive Whatever", "Prayer ", "Believe", "Ask"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:7-8', 'ESV', '["Seeks Finds", "Asks Receives", "Seek", "Opened", "One", "Knocks"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:6', 'ESV', '["Unseen", "Sees", "Secret", "Room", "Reward", "Pray"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 1:35', 'ESV', '["Still Dark", "Solitary Place", "Jesus Got", "Went", "Prayed", "Morning"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 33:3', 'ESV', '["Unsearchable Things", "Know ", "Tell", "Great", "Call", "Answer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:20', 'ESV', '["Work Within Us", "Power", "Immeasurably", "Imagine", "Ask", "According"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 5:14-15', 'ESV', '["Ask Anything According", "Ask \u2014", "Hears Us", "Approaching God", "Know", "Confidence"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:17,18', 'ESV', '["Pray Continually", "Give Thanks", "Christ Jesus", "God", "Circumstances"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Samuel 12:23', 'ESV', '["Way", "Teach", "Sin", "Right", "Pray", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 9:37-38', 'ESV', '["Harvest Field ", "Harvest", "Workers", "Therefore", "Send", "Said"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 13:15', 'ESV', '["Praise \u2014", "Openly Profess", "Therefore", "Sacrifice", "Name", "Lips"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 146:1-2', 'ESV', '["Sing Praise", "Praise", "Soul", "Lord", "Long", "Live"]');

    -- Pack: Empowered For Multiplication
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Empowered For Multiplication', 'Empowered-For-Multiplication', 'Empowered For Multiplication', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:15', 'ESV', '["Hearts Revere Christ", "Respect", "Reason", "Prepared", "Lord", "Hope"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 4:10', 'ESV', '["Use Whatever Gift", "Various Forms", "Serve Others", "Faithful Stewards", "Received", "Grace"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 3:5-6', 'ESV', '["Ways Submit", "Paths Straight", "Understanding", "Trust", "Make", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 24:1', 'ESV', '["World", "Lord", "Live", "Everything", "Earth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 4:7-8', 'ESV', '["Righteous Judge", "Good Fight", "Day \u2014", "Store", "Righteousness", "Race"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:2-3', 'ESV', '["Right Hand", "Lose Heart", "Joy Set", "Grow Weary", "Throne", "Sinners"]');

    -- Pack: Equipped For Ministry
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Equipped For Ministry', 'Equipped-For-Ministry', 'Equipped For Ministry', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 145:3-4', 'ESV', '["One Generation Commends", "Mighty Acts", "One", "Worthy", "Works", "Tell"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:10', 'ESV', '["Made Holy", "Jesus Christ", "Sacrifice", "Body"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '["Penetrates Even", "Edged Sword", "Dividing Soul", "Word", "Thoughts", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 4:6-7', 'ESV', '["Every Situation", "Christ Jesus", "Understanding", "Transcends", "Thanksgiving", "Requests"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:32', 'ESV', '["Christ God Forgave", "One Another", "Kind", "Forgiving", "Compassionate"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:2', 'ESV', '["Many Witnesses Entrust", "Teach Others", "Reliable People", "Things", "Say", "Qualified"]');

    -- Pack: Hide His word in your heart
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Hide His word in your heart', 'Hide-His-word-in-your-heart', 'Hide His word in your heart', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:37', 'ESV', '["Loved Us", "Things", "Conquerors"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:28', 'ESV', '["Things God Works", "Called According", "Purpose", "Love", "Know", "Good"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:32', 'ESV', '["Graciously Give Us", "Things", "Spare", "Son", "Gave", "Also"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:26', 'ESV', '["Spirit Helps Us", "Wordless Groans", "Spirit", "Weakness", "Way", "Pray"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:15', 'ESV', '["Father ", "Received Brought", "Received", "Spirit", "Sonship", "Slaves"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:1-2', 'ESV', '["Gives Life", "Christ Jesus", "Therefore", "Spirit", "Sin", "Set"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:23', 'ESV', '["Fall Short", "Sinned", "God", "Glory"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 6:23', 'ESV', '["Eternal Life", "Christ Jesus", "Wages", "Sin", "Lord", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 1:16', 'ESV', '["Brings Salvation", "Power", "Jew", "Gospel", "God", "Gentile"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:3-5', 'ESV', '["Suffering Produces Perseverance", "Holy Spirit", "Also Glory", "Put Us", "Perseverance", "Sufferings"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:8', 'ESV', '["Still Sinners", "God Demonstrates", "Christ Died", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 3:12-14', 'ESV', '["Taken Hold", "Take Hold", "One Thing", "Christ Jesus", "Already Obtained", "Already Arrived"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:14-16', 'ESV', '["May Receive Mercy", "Great High Priest", "High Priest", "\u2014 Yet", "Let Us", "Help Us"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:1-3', 'ESV', '["Let Us Throw", "Let Us Run", "Right Hand", "Race Marked", "Lose Heart", "Joy Set"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:10', 'ESV', '["Loved Us", "Loved God", "Atoning Sacrifice", "Son", "Sins", "Sent"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:5', 'ESV', '["Bear Much Fruit", "Vine", "Remain", "Nothing", "Branches", "Apart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 16:33', 'ESV', '["Take Heart", "World ", "World", "Trouble", "Told", "Things"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:33', 'ESV', '["Seek First", "Well", "Things", "Righteousness", "Kingdom", "Given"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 8:12', 'ESV', '["Whoever Follows", "Never Walk", "Life ", "Jesus Spoke", "World", "Said"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 4:23', 'ESV', '["True Worshipers", "Father Seeks", "Worshipers", "Father", "Yet", "Worship"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:27', 'ESV', '["World Gives", "Troubled", "Peace", "Let", "Leave", "Hearts"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:14', 'ESV', '["Word Became Flesh", "Dwelling Among Us", "Truth", "Son", "Seen", "One"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 1:27', 'ESV', '["God Created Mankind", "God", "Created", "Male", "Image", "Female"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 19:10', 'ESV', '["Man Came", "Lost ", "Son", "Seek", "Save"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 6:2-3', 'ESV', '["May Go Well", "Promise \u2014", "Mother \"\u2014", "First Commandment", "Earth ", "Honor"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 1:7', 'ESV', '["Fools Despise Wisdom", "Lord", "Knowledge", "Instruction", "Fear", "Beginning"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Lamentations 3:22-23', 'ESV', '["New Every Morning", "Compassions Never Fail", "Great Love", "Great", "Lord", "Faithfulness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:6-7', 'ESV', '["Five Sparrows Sold", "Many Sparrows", "Two Pennies", "Yet", "Worth", "One"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:34', 'ESV', '["Treasure", "Heart", "Also"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 90:14', 'ESV', '["Unfailing Love", "Satisfy Us", "May Sing", "Morning", "Joy", "Glad"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 13:4', 'ESV', '["Proud", "Patient", "Love", "Kind", "Envy", "Boast"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:2-3', 'ESV', '["Faith Produces Perseverance", "Pure Joy", "Many Kinds", "Face Trials", "Whenever", "Testing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:16', 'ESV', '["May See", "Light Shine", "Good Deeds", "Way", "Others", "Let"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 28:19-20', 'ESV', '["Therefore Go", "Obey Everything", "Make Disciples", "Holy Spirit", "Age ", "Teaching"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 3:18', 'ESV', '["Let Us", "Dear Children", "Words", "Truth", "Speech", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 4:23', 'ESV', '["Heart", "Guard", "Flows", "Everything", "Else"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:16', 'ESV', '["Holy ", "Holy", "Written"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 5:7', 'ESV', '["Cast", "Cares", "Anxiety"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 1:4', 'ESV', '["Mankind", "Light", "Life"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 46:10', 'ESV', '["Earth ", "Exalted Among", "Exalted", "Still", "Says", "Nations"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 16:26', 'ESV', '["Yet Forfeit", "Whole World", "Anyone Give", "Soul", "Someone", "Good"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:18', 'ESV', '["Open", "Law", "Eyes"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 139:14', 'ESV', '["Wonderfully Made", "Full Well", "Works", "Wonderful", "Praise", "Know"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:13', 'ESV', '["Remains Faithful", "Cannot Disown", "Faithless"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 19:21', 'ESV', '["Purpose", "Prevails", "Plans", "Person", "Many", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 46:1', 'ESV', '["Present Help", "Trouble", "Strength", "Refuge", "God", "Ever"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 100:1-5', 'ESV', '["Love Endures Forever", "Made Us", "Joyful Songs", "Give Thanks", "Faithfulness Continues", "Worship"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Habakkuk 3:17-19', 'ESV', '["Olive Crop Fails", "Stringed Instruments", "Fig Tree", "Fields Produce", "Sovereign Lord", "Feet Like"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Nehemiah 9:6', 'ESV', '["Starry Host", "Heaven Worship", "Give Life", "Highest Heavens", "Heavens", "Seas"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 16:2', 'ESV', '["Good Thing ", "Say", "Lord", "Apart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:20-21', 'ESV', '["Work Within Us", "Christ Jesus Throughout", "Power", "Immeasurably", "Imagine", "Glory"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 42:11', 'ESV', '["Yet Praise", "Disturbed Within", "Soul", "Savior", "Put", "Hope"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 139:1-18', 'ESV', '["Would Outnumber", "Unformed Body", "Shine Like", "Secret Place", "Sand \u2014", "Full Well"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:9', 'ESV', '["Wonderful Light", "Special Possession", "Royal Priesthood", "May Declare", "Holy Nation", "Chosen People"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:3-4', 'ESV', '["Quiet Spirit", "Outward Adornment", "Inner Self", "Great Worth", "Gold Jewelry", "Fine Clothes"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 2:6-7', 'ESV', '["Received Christ Jesus", "Thankfulness", "Taught", "Strengthened", "Rooted", "Overflowing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:22', 'ESV', '["Sincere Love", "Truth", "Purified", "Obeying", "Heart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 3:13', 'ESV', '["Forgive One Another", "Lord Forgave", "Forgive", "Someone", "Grievance", "Bear"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 3:15', 'ESV', '["Hearts Revere Christ", "Respect", "Reason", "Prepared", "Lord", "Hope"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 4:10-11', 'ESV', '["Use Whatever Gift", "Things God May", "Strength God Provides", "Various Forms", "Serve Others", "Jesus Christ"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 1:1-3', 'ESV', '["Wither \u2014 Whatever", "Whose Leaf", "Whose Delight", "Tree Planted", "Sinners Take", "Law Day"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 112:7', 'ESV', '["Bad News", "Trusting", "Steadfast", "Lord", "Hearts", "Fear"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:23', 'ESV', '["Promised", "Profess", "Hope", "Faithful"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 1:17', 'ESV', '["Perfect Gift", "Heavenly Lights", "Every Good", "Father", "Coming"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 5:10', 'ESV', '["Eternal Glory", "Suffered", "Strong", "Steadfast", "Restore", "Make"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 3:7-8', 'ESV', '["Shun Evil", "Bring Health", "Wise", "Nourishment", "Lord", "Fear"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:28', 'ESV', '["Things God Works", "Called According", "Purpose", "Love", "Know", "Good"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 43:4', 'ESV', '["Give People", "Since", "Sight", "Precious", "Nations", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 27:13', 'ESV', '["Remain Confident", "See", "Lord", "Living", "Land", "Goodness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Jeremiah 2:13', 'ESV', '["Committed Two Sins", "Cannot Hold Water", "Living Water", "Broken Cisterns", "Cisterns", "Spring"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 1:15', 'ESV', '["Save Sinners \u2014", "Deserves Full Acceptance", "Christ Jesus Came", "Trustworthy Saying", "Worst", "World"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 143:8', 'ESV', '["Unfailing Love", "Morning Bring", "Word", "Way", "Trust", "Show"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:5', 'ESV', '["Every Pretension", "Demolish Arguments", "Sets", "Obedient", "Make", "Knowledge"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:18', 'ESV', '["Open", "Law", "Eyes"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 4:12', 'ESV', '["Penetrates Even", "Edged Sword", "Dividing Soul", "Word", "Thoughts", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:105', 'ESV', '["Word", "Path", "Light", "Lamp", "Feet"]');

    -- Pack: Come Follow Me - Answering the Call
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Come Follow Me - Answering the Call', 'Come-Follow-Me-Answering-the-Call', 'Come Follow Me - Answering the Call', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 4:19', 'ESV', '[" Jesus Said", "People ", "Send", "Follow", "Fish", "Come"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 11:28-30', 'ESV', '["Light ", "Yoke Upon", "Find Rest", "Yoke", "Rest", "Weary"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:18', 'ESV', '["Unveiled Faces Contemplate", "Increasing Glory", "Glory", "Transformed", "Spirit", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:6', 'ESV', '["One Comes", "Jesus Answered", "Father Except", "Way", "Truth", "Life"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:5', 'ESV', '["Bear Much Fruit", "Vine", "Remain", "Nothing", "Branches", "Apart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 6:68-69', 'ESV', '["Simon Peter Answered", "Holy One", "God ", "Eternal Life", "Words", "Shall"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 10:37-38', 'ESV', '["Worthy", "Whoever", "Take", "Son", "Mother", "Loves"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 9:23', 'ESV', '["Disciple Must Deny", "Whoever Wants", "Cross Daily", "Take", "Said", "Follow"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 14:33', 'ESV', '["Way", "Give", "Everything", "Disciples", "Cannot"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 9:62', 'ESV', '["Looks Back", "Jesus Replied", "God ", "Service", "Puts", "Plow"]');

    -- Pack: Disciple-making
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Disciple-making', 'Disciple-making', 'Disciple-making', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:16', 'ESV', '["Might Go", "Last \u2014", "Whatever", "Name", "Give", "Father"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:14-15', 'ESV', '["Love Compels Us", "One Died", "Longer Live", "Live", "Died", "Therefore"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:12-13', 'ESV', '["Whole Measure", "Reach Unity", "Become Mature", "Christ May", "Christ", "Works"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:28', 'ESV', '["Teaching Everyone", "Wisdom", "Proclaim", "One", "Christ", "Admonishing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:2', 'ESV', '["Many Witnesses Entrust", "Teach Others", "Reliable People", "Things", "Say", "Qualified"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 2:8', 'ESV', '["Well", "Share", "Much", "Loved", "Lives", "Gospel"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Timothy 4:12', 'ESV', '["Let Anyone Look", "Young", "Speech", "Set", "Purity", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 78:72', 'ESV', '["Skillful Hands", "David Shepherded", "Led", "Integrity", "Heart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 3:10', 'ESV', '["May See", "Supply", "Pray", "Night", "Lacking", "Faith"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ezra 7:10', 'ESV', '["Teaching", "Study", "Observance", "Lord", "Laws", "Law"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 2:20-21', 'ESV', '["Show Genuine Concern", "One Else Like", "Jesus Christ", "Everyone Looks", "Welfare", "Interests"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Colossians 1:29', 'ESV', '["Strenuously Contend", "Powerfully Works", "Energy Christ", "End"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:23', 'ESV', '["Give Careful Attention", "Sure", "Know", "Herds", "Flocks", "Condition"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 4:2', 'ESV', '["Great Patience", "Encourage \u2014", "Careful Instruction", "Word", "Season", "Rebuke"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 14:26', 'ESV', '["Sisters \u2014 Yes", "Life \u2014", "Person Cannot", "Hate Father", "Anyone Comes", "Wife"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 14:27', 'ESV', '["Whoever", "Follow", "Disciple", "Cross", "Carry", "Cannot"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 14:33', 'ESV', '["Way", "Give", "Everything", "Disciples", "Cannot"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 8:31-32', 'ESV', '["Jesus Said", "Free ", "Truth", "Teaching", "Set", "Really"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:8', 'ESV', '["Bear Much Fruit", "Showing", "Glory", "Father", "Disciples"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 13:34-35', 'ESV', '["Love One Another", "New Command", "Loved", "Know", "Give", "Everyone"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 3:10-11', 'ESV', '["Becoming Like", "Want", "Sufferings", "Somehow", "Resurrection", "Power"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 20:24', 'ESV', '["Life Worth Nothing", "Lord Jesus", "Good News", "Testifying", "Task", "Race"]');

    -- Pack: Sufferings Pack
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Sufferings Pack', 'Sufferings-Pack', 'Sufferings Pack', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 6:22-23', 'ESV', '["People Hate", "Ancestors Treated", "Son", "Reward", "Rejoice", "Reject"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 2:20-21', 'ESV', '["Christ Suffered", "Wrong", "Suffer", "Steps", "Receive", "Leaving"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:18-19', 'ESV', '["Would Love", "World Hates", "World", "Mind", "Keep", "Hated"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:6', 'ESV', '["Son ", "Lord Disciplines", "Chastens Everyone", "One", "Loves", "Accepts"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:2', 'ESV', '["Every Branch", "Bear Fruit", "Fruit", "Prunes", "Fruitful", "Even"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:6-7', 'ESV', '["Suffer Grief", "Proven Genuineness", "Jesus Christ", "Greatly Rejoice", "Greater Worth", "Faith \u2014"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 1:4-5', 'ESV', '["Counted Worthy", "Among God", "God", "Trials", "Therefore", "Suffering"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:25 ', 'ESV', '["Mistreated Along", "God Rather", "Fleeting Pleasures", "Sin", "People", "Enjoy"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 1:20-22 ', 'ESV', '["Taken Away", "Praised ", "Charging God", "Lord Gave", "Job Got", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 5:40-41', 'ESV', '["Suffering Disgrace", "Speech Persuaded", "Counted Worthy", "Apostles Left", "Apostles", "Speak"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 7:59-60', 'ESV', '["Stephen Prayed", "Spirit ", "Lord Jesus", "Fell Asleep", "Lord", "Fell"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:18', 'ESV', '["Worth Comparing", "Present Sufferings", "Revealed", "Glory", "Consider"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:11', 'ESV', '["Discipline Seems Pleasant", "Trained", "Time", "Righteousness", "Produces", "Peace"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 1:3-4', 'ESV', '["Lord Jesus Christ", "Comforts Us", "Troubles", "Trouble", "Receive", "Praise"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:67,71', 'ESV', '["Went Astray", "Might Learn", "Word", "Obey", "Good", "Decrees"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 4:17', 'ESV', '["Momentary Troubles", "Far Outweighs", "Eternal Glory", "Light", "Achieving"]');

    -- Pack: The Pursuit of Holiness
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'The Pursuit of Holiness', 'The-Pursuit-of-Holiness', 'The Pursuit of Holiness', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:19-20', 'ESV', '["Therefore Honor God", "Holy Spirit", "God", "Temples", "Received", "Price"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 2:11-12', 'ESV', '["Worldly Passions", "Teaches Us", "Present Age", "Offers Salvation", "Live Self", "Godly Lives"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 39:8-9', 'ESV', '["Withheld Nothing", "Wicked Thing", "God ", "Charge ", "Wife", "Told"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 4:3', 'ESV', '["Avoid Sexual Immorality", "Sanctified", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:13-14', 'ESV', '["Live According", "Live", "Spirit", "Put", "Misdeeds", "Led"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 6:7-8', 'ESV', '["Reap Eternal Life", "Reap Destruction", "Man Reaps", "God Cannot", "Whoever Sows", "Sows"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 24:3-6', 'ESV', '["Receive Blessing", "Pure Heart", "May Stand", "May Ascend", "Holy Place", "Clean Hands"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 12:1', 'ESV', '["Proper Worship", "Living Sacrifice", "God \u2014", "God", "View", "Urge"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Philippians 3:8', 'ESV', '["May Gain Christ", "Knowing Christ Jesus", "Whose Sake", "Surpassing Worth", "Consider Everything", "Consider"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:22-24', 'ESV', '["True Righteousness", "Old Self", "New Self", "Made New", "Like God", "Former Way"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:5-6', 'ESV', '["Minds Set", "Mind Governed", "Spirit Desires", "Live According", "Flesh Desires", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 4:7', 'ESV', '["Must Rule", "Sin", "Right", "Door", "Desires", "Crouching"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:19', 'ESV', '["Sets Aside One", "Teaches Others Accordingly", "Whoever Practices", "Therefore Anyone", "Called Great", "Called Least"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:1', 'ESV', '["Let Us Throw", "Let Us Run", "Race Marked", "Great Cloud", "Easily Entangles", "Witnesses"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 51:10-13', 'ESV', '["Steadfast Spirit Within", "Willing Spirit", "Holy Spirit", "Turn Back", "Teach Transgressors", "Pure Heart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:1-2', 'ESV', '["Milk Without Money", "Without Cost", "Spend Money", "Buy Wine", "Money", "Buy"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 2:15-16', 'ESV', '["Life \u2014 Comes", "World \u2014", "Anyone Loves", "World", "Pride", "Lust"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:2', 'ESV', '["Every Branch", "Bear Fruit", "Fruit", "Prunes", "Fruitful", "Even"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:20-21', 'ESV', '["Special Purposes", "Made Holy", "Large House", "Good Work", "Common Use", "Wood"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 17:19', 'ESV', '["Truly Sanctified", "Sanctify", "May"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 3:21', 'ESV', '["Whoever Lives", "Truth Comes", "Seen Plainly", "Sight", "May", "Light"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 13:4', 'ESV', '["Sexually Immoral", "Marriage", "Judge", "Honored", "God", "Adulterer"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:29-33', 'ESV', '["Also Must Love", "Profound Mystery \u2014", "Wife Must Respect", "One Ever Hated", "Church \u2014", "Wife"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Genesis 1:27', 'ESV', '["God Created Mankind", "God", "Created", "Male", "Image", "Female"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:9', 'ESV', '["Purify Us", "Forgive Us", "Unrighteousness", "Sins", "Faithful", "Confess"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Lamentations 3:22-23', 'ESV', '["New Every Morning", "Compassions Never Fail", "Great Love", "Great", "Lord", "Faithfulness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 2:11-12', 'ESV', '["Worldly Passions", "Teaches Us", "Present Age", "Offers Salvation", "Live Self", "Godly Lives"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 5:14-15', 'ESV', '["Love Compels Us", "One Died", "Longer Live", "Live", "Died", "Therefore"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:13-14', 'ESV', '["Live According", "Live", "Spirit", "Put", "Misdeeds", "Led"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:16', 'ESV', '["Walk", "Spirit", "Say", "Gratify", "Flesh", "Desires"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 11:25-26', 'ESV', '["Regarded Disgrace", "Mistreated Along", "Looking Ahead", "Greater Value", "God Rather", "Fleeting Pleasures"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 12:2', 'ESV', '["Right Hand", "Joy Set", "Throne", "Shame", "Scorning", "Sat"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 7:1-5', 'ESV', '["Wayward Woman", "Sister ", "Relative ", "Adulterous Woman", "Seductive Words", "Commands Within"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:9, 11', 'ESV', '["Young Person Stay", "Living According", "Word", "Sin", "Purity", "Path"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 15:7-8', 'ESV', '["Bear Much Fruit", "Ask Whatever", "Words Remain", "Remain", "Wish", "Showing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 55:22', 'ESV', '["Never Let", "Sustain", "Shaken", "Righteous", "Lord", "Cast"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 3:13', 'ESV', '["Today ", "Sin", "None", "May", "Long", "Hardened"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 2:22', 'ESV', '["Pursue Righteousness", "Pure Heart", "Evil Desires", "Youth", "Peace", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:3-4', 'ESV', '["Given Us Everything", "Given Us", "Called Us", "World Caused", "Precious Promises", "May Participate"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 10:4-5', 'ESV', '["Every Pretension", "Divine Power", "Demolish Strongholds", "Demolish Arguments", "World", "Weapons"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 97:10', 'ESV', '["Lord Hate Evil", "Faithful Ones", "Wicked", "Love", "Lives", "Let"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:29', 'ESV', '["Right Eye Causes", "Lose One Part", "Whole Body", "Body", "Thrown", "Throw"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:44', 'ESV', '["One Another", "Believe Since", "Accept Glory", "Glory", "Seek", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 28:28', 'ESV', '["Understanding ", "Shun Evil", "Lord \u2014", "Human Race", "Wisdom", "Said"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 2:1-2', 'ESV', '["Whole World", "Righteous One", "Dear Children", "Atoning Sacrifice", "Write", "Sins"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Micah 7:8', 'ESV', '["Though", "Sit", "Rise", "Lord", "Light", "Gloat"]');

    -- Pack: Excuses Pack
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'Excuses Pack', 'Excuses-Pack', 'Excuses Pack', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 1:10', 'ESV', '["Word", "Sinned", "Make", "Liar", "Claim"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:10', 'ESV', '["One Righteous", "Even One", "Written"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 14:12', 'ESV', '["Way", "Right", "Leads", "End", "Death", "Appears"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 10:31', 'ESV', '["Living God", "Dreadful Thing", "Hands", "Fall"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Mark 8:36', 'ESV', '["Yet Forfeit", "Whole World", "Soul", "Someone", "Good", "Gain"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 6:33', 'ESV', '["Seek First", "Well", "Things", "Righteousness", "Kingdom", "Given"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 119:130', 'ESV', '["Words Gives Light", "Gives Understanding", "Unfolding", "Simple"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:14', 'ESV', '["Person Without", "Cannot Understand", "Things", "Spirit", "God", "Foolishness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 4:4', 'ESV', '["Cannot See", "Unbelievers", "Minds", "Light", "Image", "Gospel"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 3:19-20', 'ESV', '["Evil Hates", "Evil", "World", "Verdict", "Light", "Fear"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 5:32', 'ESV', '["Repentance ", "Sinners", "Righteous", "Come", "Call"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 4:4-5', 'ESV', '["Trusts God", "Works", "Work", "Wages", "Ungodly", "Righteousness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 19:10', 'ESV', '["Man Came", "Lost ", "Son", "Seek", "Save"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 1:18', 'ESV', '["Matter ,\" Says", "Let Us Settle", "Like Wool", "Like Scarlet", "White", "Though"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 29:25', 'ESV', '["Whoever Trusts", "Kept Safe", "Snare", "Prove", "Man", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 5:11-12', 'ESV', '["People Insult", "Falsely Say", "Way", "Reward", "Rejoice", "Prophets"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 3:3', 'ESV', '["Evil One", "Strengthen", "Protect", "Lord", "Faithful"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Peter 1:5', 'ESV', '["Last Time", "Shielded", "Salvation", "Revealed", "Ready", "Power"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Proverbs 27:1', 'ESV', '["Day May Bring", "Tomorrow", "Know", "Boast"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:6', 'ESV', '["Seek", "Near", "May", "Lord", "Found", "Call"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ecclesiastes 11:9', 'ESV', '["Things God", "Eyes See", "Heart Give", "Heart", "Youth", "Young"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 12:19-20', 'ESV', '["Take Life Easy", "Merry ", "Many Years", "Grain Laid", "God Said", "Life"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 14:12', 'ESV', '["God", "Give", "Account"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 13:16', 'ESV', '["Turn", "Indeed", "Deliverance"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 1:18', 'ESV', '["Saved", "Power", "Perishing", "Message", "God", "Foolishness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 55:8,9', 'ESV', '["Ways ,\" Declares", "Ways Higher", "Ways", "Higher", "Thoughts", "Neither"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 3:3', 'ESV', '["Unfaithfulness Nullify God", "Unfaithful", "Faithfulness"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 5:44', 'ESV', '["One Another", "Believe Since", "Accept Glory", "Glory", "Seek", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 8:20', 'ESV', '["Speak According", "Consult God", "Word", "Warning", "Testimony", "Light"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 43:11', 'ESV', '["Savior", "Lord", "Even", "Apart"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 25:41', 'ESV', '["Eternal Fire Prepared", "Say", "Left", "Devil", "Depart", "Cursed"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Matthew 7:13-14', 'ESV', '["Many Enter", "Narrow Gate", "Narrow", "Gate", "Enter", "Wide"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 14:1', 'ESV', '["God ", "Fool Says", "Vile", "One", "Heart", "Good"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 1:20', 'ESV', '["Invisible Qualities \u2014", "Divine Nature \u2014", "World God", "Without Excuse", "Eternal Power", "Clearly Seen"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Timothy 3:16', 'ESV', '["Useful", "Training", "Teaching", "Scripture", "Righteousness", "Rebuking"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Peter 1:21', 'ESV', '["Prophecy Never", "Holy Spirit", "Carried Along", "Though Human", "Human", "Spoke"]');

    -- Pack: The Spirit-Filled Life
    INSERT INTO verse_packs (id, title, identifier, description, is_public)
    VALUES (gen_random_uuid(), 'The Spirit-Filled Life', 'The-Spirit-Filled-Life', 'The Spirit-Filled Life', true)
    RETURNING id INTO pack_id;
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Isaiah 44:3', 'ESV', '["Thirsty Land", "Dry Ground", "Pour Water", "Pour", "Streams", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 7:37-38', 'ESV', '["Whoever Believes", "Thirsty Come", "Loud Voice", "Living Water", "Let Anyone", "Jesus Stood"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 15:45', 'ESV', '["Last Adam", "Giving Spirit", "Written", "Living", "Life"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Job 33:4', 'ESV', '["Almighty Gives", "Spirit", "Made", "Life", "God", "Breath"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:9', 'ESV', '["God Lives", "Spirit", "Realm", "Indeed", "However", "Flesh"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:16', 'ESV', '["Walk", "Spirit", "Say", "Gratify", "Flesh", "Desires"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:15-18', 'ESV', '["Live \u2014", "Get Drunk", "Every Opportunity", "Wise", "Wine", "Unwise"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:19-20', 'ESV', '["Always Giving Thanks", "Lord Jesus Christ", "One Another", "Make Music", "Lord", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ezekiel 36:26-27', 'ESV', '["New Spirit", "New Heart", "Spirit", "Heart", "Stone", "Remove"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 2:29', 'ESV', '["Written Code", "One Inwardly", "Spirit", "Praise", "Person", "People"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:22-23', 'ESV', '["Things", "Spirit", "Self", "Peace", "Love", "Law"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:13', 'ESV', '["Given Us", "Spirit", "Live", "Know"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:17', 'ESV', '["Whoever", "United", "Spirit", "One", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 3:16-19', 'ESV', '["Surpasses Knowledge \u2014", "Christ May Dwell", "May Strengthen", "Holy People", "Glorious Riches", "May"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 5:5', 'ESV', '["Holy Spirit", "Put Us", "Shame", "Poured", "Love", "Hope"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'John 14:23', 'ESV', '["Jesus Replied", "Teaching", "Obey", "Make", "Loves", "Love"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Luke 15:31', 'ESV', '["Son ", "Father Said", "Everything", "Always"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Thessalonians 2:13', 'ESV', '["Thank God", "Sisters Loved", "Sanctifying Work", "Ought Always", "God Chose", "Truth"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 6:11', 'ESV', '["Lord Jesus Christ", "Washed", "Spirit", "Sanctified", "Name", "Justified"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:17', 'ESV', '["Spirit", "Lord", "Freedom"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:18', 'ESV', '["Unveiled Faces Contemplate", "Increasing Glory", "Glory", "Transformed", "Spirit", "Lord"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Titus 3:5-7', 'ESV', '["Might Become Heirs", "Us Generously", "Saved Us", "Righteous Things", "Jesus Christ", "Holy Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:11-12', 'ESV', '["Freely Given Us", "May Understand", "Thoughts Except", "Spirit Within", "One Knows", "God Except"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Corinthians 2:16', 'ESV', '["Mind", "Lord", "Known", "Instruct", "Christ"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 4:31', 'ESV', '["Holy Spirit", "God Boldly", "Word", "Spoke", "Shaken", "Prayed"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 John 4:4', 'ESV', '["Dear Children", "World", "Overcome", "One", "Greater", "God"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 4:6', 'ESV', '["God Sent", "Father ", "Spirit", "Sons", "Son", "Hearts"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:15-17', 'ESV', '["May Also Share", "Heirs \u2014 Heirs", "Father ", "Received Brought", "Heirs", "Share"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:2-3', 'ESV', '["Living God", "Human Hearts", "Hearts", "Written", "Tablets", "Stone"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 3:4-6', 'ESV', '["New Covenant \u2014", "Spirit Gives Life", "Made Us Competent", "Competence Comes", "Claim Anything", "Letter Kills"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 4:30', 'ESV', '["Holy Spirit", "Sealed", "Redemption", "Grieve", "God", "Day"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '1 Thessalonians 5:16-19', 'ESV', '["Rejoice Always", "Pray Continually", "Give Thanks", "Christ Jesus", "Spirit", "Quench"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Acts 7:51', 'ESV', '["Still Uncircumcised", "Necked People", "Holy Spirit", "Always Resist", "Stiff", "Like"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'James 4:4-5', 'ESV', '["World Means Enmity", "World Becomes", "Jealously Longs", "Adulterous People", "Therefore", "Spirit"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Hebrews 3:7-8', 'ESV', '["Holy Spirit Says", "Wilderness", "Voice", "Today", "Time", "Testing"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Galatians 5:25', 'ESV', '["Let Us Keep", "Step", "Spirit", "Since", "Live"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Romans 8:13-14', 'ESV', '["Live According", "Live", "Spirit", "Put", "Misdeeds", "Led"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, '2 Corinthians 4:7-10', 'ESV', '["Always Carry Around", "Jesus May Also", "Surpassing Power", "Hard Pressed", "Every Side", "Jesus"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Psalm 51:10-12', 'ESV', '["Steadfast Spirit Within", "Willing Spirit", "Holy Spirit", "Pure Heart", "Take", "Sustain"]');
    INSERT INTO memory_verses (verse_pack_id, reference, version, tags)
    VALUES (pack_id, 'Ephesians 5:18-21', 'ESV', '["Always Giving Thanks", "Lord Jesus Christ", "One Another", "Make Music", "Get Drunk", "Lord"]');
END; $$;