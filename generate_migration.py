import yaml
import uuid

# The YAML content provided in the tool output
yaml_content = """
Series:
  - ID: "TMS 60"
    Title: "TMS 60"
    Packs:
      - ID: "A"
        Title: "Living the New Life"
        Verses:
          - ID: "A1"
            Title: "Christ the Center"
            Reference: "2 Corinthians 5:17"
            Tags: ["Christ Jesus", "New Creation", "Old", "New"]
          - ID: "A2"
            Title: "Christ the Center"
            Reference: "Galatians 2:20"
            Tags: ["Christ Jesus", "Christ Crucified", "Faith", "Love"]
          - ID: "A3"
            Title: "Obedience to Christ"
            Reference: "Romans 12:1"
            Tags: ["Mercy", "Lord God", "Living Sacrifice", "Worship"]
          - ID: "A4"
            Title: "Obedience to Christ"
            Reference: "John 14:21"
            Tags: ["Obey", "Love", "Lord God", "Jesus Christ"]
          - ID: "A5"
            Title: "The Word"
            Reference: "2 Timothy 3:16"
            Tags: ["Lord God", "Scripture", "Teaching", "Righteousness"]
          - ID: "A6"
            Title: "The Word"
            Reference: "Joshua 1:8"
            Tags: ["Law", "Meditate", "Do", "Prosperous"]
          - ID: "A7"
            Title: "Prayer"
            Reference: "John 15:7"
            Tags: ["Remain", "Ask", "Given", "Pray"]
          - ID: "A8"
            Title: "Prayer"
            Reference: "Philippians 4:6-7"
            Tags: ["Pray", "Petition", "Request", "Peace"]
          - ID: "A9"
            Title: "Fellowship"
            Reference: "Matthew 18:20"
            Tags: ["Come", "Together", "Jesus Christ", "Name"]
          - ID: "A10"
            Title: "Fellowship"
            Reference: "Hebrews 10:24-25"
            Tags: ["One Another", "Love", "Good Deeds", "Encourage"]
          - ID: "A11"
            Title: "Witnessing"
            Reference: "Matthew 4:19"
            Tags: ["Follow", "Fishers", "Men", "Jesus Christ"]
          - ID: "A12"
            Title: "Witnessing"
            Reference: "Romans 1:16"
            Tags: ["Gospel", "Power", "Lord God", "Salvation"]
      - ID: "B"
        Title: "Proclaiming Christ"
        Verses:
          - ID: "B1"
            Title: "All have Sinned"
            Reference: "Romans 3:23"
            Tags: ["Sin", "Fall", "Glory", "Lord God"]
          - ID: "B2"
            Title: "All have Sinned"
            Reference: "Isaiah 53:6"
            Tags: ["Astray", "Own Way", "Iniquity", "Jesus Christ"]
          - ID: "B3"
            Title: "Sin's Penalty"
            Reference: "Romans 6:23"
            Tags: ["Sin", "Death", "Gift", "Life"]
          - ID: "B4"
            Title: "Sin's Penalty"
            Reference: "Hebrews 9:27"
            Tags: ["Death", "Judgement", "Man", "Destiny"]
          - ID: "B5"
            Title: "Christ paid the Penalty"
            Reference: "Romans 5:8"
            Tags: ["Lord God", "Love", "Sinner", "Christ Jesus"]
          - ID: "B6"
            Title: "Christ paid the Penalty"
            Reference: "1 Peter 3:18"
            Tags: ["Christ Jesus", "Died", "Sins", "Reconcile"]
          - ID: "B7"
            Title: "Salvation not by Works"
            Reference: "Ephesians 2:8-9"
            Tags: ["Grace", "Salvation", "Faith", "Works"]
          - ID: "B8"
            Title: "Salvation not by Works"
            Reference: "Titus 3:5"
            Tags: ["Salvation", "Righteous", "Mercy", "Holy Spirit"]
          - ID: "B9"
            Title: "Must receive Christ"
            Reference: "John 1:12"
            Tags: ["Received", "Believe", "Children", "Lord God"]
          - ID: "B10"
            Title: "Must receive Christ"
            Reference: "Revelation 3:20"
            Tags: ["Here", "Knock", "Hears", "Opens"]
          - ID: "B11"
            Title: "Assurance of Salvation"
            Reference: "1 John 5:13"
            Tags: ["Believe", "Christ Jesus", "Eternal Life", "Know"]
          - ID: "B12"
            Title: "Assurance of Salvation"
            Reference: "John 5:24"
            Tags: ["Hears", "Believe", "Eternal Life", "Life"]
      - ID: "C"
        Title: "Reliance on God's Resources"
        Verses:
          - ID: "C1"
            Title: "His Spirit"
            Reference: "1 Corinthians 3:16"
            Tags: ["Temple", "Lord God", "Spirit", "Within"]
          - ID: "C2"
            Title: "His Spirit"
            Reference: "1 Corinthians 2:12"
            Tags: ["World", "Spirit", "Lord God", "Freely"]
          - ID: "C3"
            Title: "His Strength"
            Reference: "Isaiah 41:10"
            Tags: ["Fear", "Dismay", "Strengthen", "Help"]
          - ID: "C4"
            Title: "His Strength"
            Reference: "Philippians 4:13"
            Tags: ["Do", "Everything", "Strengthen", "Christ"]
          - ID: "C5"
            Title: "His Faithfulness"
            Reference: "Lamentations 3:22-23"
            Tags: ["Faithfulness", "Love", "Compassions", "Change"]
          - ID: "C6"
            Title: "His Faithfulness"
            Reference: "Numbers 23:19"
            Tags: ["Lie", "Change", "Say", "Act"]
          - ID: "C7"
            Title: "His Peace"
            Reference: "Isaiah 26:3"
            Tags: ["Peace", "Lord God", "Steadfast", "Trust"]
          - ID: "C8"
            Title: "His Peace"
            Reference: "1 Peter 5:7"
            Tags: ["Lord God", "Anxiety", "Care", "Worry"]
          - ID: "C9"
            Title: "His Provision"
            Reference: "Romans 8:32"
            Tags: ["Lord God", "Son", "Gracious", "Give"]
          - ID: "C10"
            Title: "His Provision"
            Reference: "Philippians 4:19"
            Tags: ["Need", "Riches", "Christ Jesus", "Lord God"]
          - ID: "C11"
            Title: "His Help in Temptation"
            Reference: "Hebrews 2:18"
            Tags: ["Suffered", "Tempted", "Help", "Temptation"]
          - ID: "C12"
            Title: "His Help in Temptation"
            Reference: "Psalms 119:9,11"
            Tags: ["Pure", "Word", "Heart", "Sin"]
      - ID: "D"
        Title: "Being Christ's Disciple"
        Verses:
          - ID: "D1"
            Title: "Put Christ First"
            Reference: "Matthew 6:33"
            Tags: ["Christ Jesus", "First", "Kingdom", "Righteousness"]
          - ID: "D2"
            Title: "Put Christ First"
            Reference: "Luke 9:23"
            Tags: ["Christ Jesus", "Deny", "Cross", "Follow"]
          - ID: "D3"
            Title: "Separate from the World"
            Reference: "1 John 2:15-16"
            Tags: ["World", "Lust", "Pride", "Craving"]
          - ID: "D4"
            Title: "Separate from the World"
            Reference: "Romans 12:2"
            Tags: ["World", "Conform", "Transforme", "Renewing"]
          - ID: "D5"
            Title: "Be Steadfast"
            Reference: "1 Corinthians 15:58"
            Tags: ["Steadfast", "Labor", "Work", "Stand Firm"]
          - ID: "D6"
            Title: "Be Steadfast"
            Reference: "Hebrews 12:3"
            Tags: ["Endure", "Opposition", "Weary", "Lose Heart"]
          - ID: "D7"
            Title: "Serve Others"
            Reference: "Mark 10:45"
            Tags: ["Serve", "Ransom", "Give", "Christ Jesus"]
          - ID: "D8"
            Title: "Serve Others"
            Reference: "2 Corinthians 4:5"
            Tags: ["Serve", "Servants", "Christ Jesus", "Lord"]
          - ID: "D9"
            Title: "Give Generously"
            Reference: "Proverbs 3:9-10"
            Tags: ["Tithe", "First", "Filled", "Overflow"]
          - ID: "D10"
            Title: "Give Generously"
            Reference: "2 Corinthians 9:6-7"
            Tags: ["Tithe", "Give", "Generous", "Wealth"]
          - ID: "D11"
            Title: "Develop World Vision"
            Reference: "Acts 1:8"
            Tags: ["Power", "Holy Spirit", "Witness", "World"]
          - ID: "D12"
            Title: "Develop World Vision"
            Reference: "Matthew 28:19-20"
            Tags: ["Disciples", "Baptise", "End", "Teaching"]
      - ID: "E"
        Title: "Growth in Christlikeness"
        Verses:
          - ID: "E1"
            Title: "Love"
            Reference: "John 13:34-35"
            Tags: ["Jesus Christ", "Love", "Disciples", "One Another"]
          - ID: "E2"
            Title: "Love"
            Reference: "1 John 3:18"
            Tags: ["Love", "Words", "Actions", "Truth"]
          - ID: "E3"
            Title: "Humility"
            Reference: "Philippians 2:3-4"
            Tags: ["Selfish", "Conceit", "Humility", "Interests"]
          - ID: "E4"
            Title: "Humility"
            Reference: "1 Peter 5:5-6"
            Tags: ["Submission", "Humility", "Proud", "Humble"]
          - ID: "E5"
            Title: "Purity"
            Reference: "Ephesians 5:3"
            Tags: ["Sexual Immorality", "Impurity", "Greed", "Improper", "Holy"]
          - ID: "E6"
            Title: "Purity"
            Reference: "1 Peter 2:11"
            Tags: ["Aliens", "Strangers", "Abstain", "Sin"]
          - ID: "E7"
            Title: "Honesty"
            Reference: "Leviticus 19:11"
            Tags: ["Lie", "Cheat", "Steal", "Deception"]
          - ID: "E8"
            Title: "Honesty"
            Reference: "Acts 24:16"
            Tags: ["Conscience", "Honesty", "Integrity", "Lord God"]
          - ID: "E9"
            Title: "Faith"
            Reference: "Hebrews 11:6"
            Tags: ["Lord God", "Believe", "Seek", "Faith"]
          - ID: "E10"
            Title: "Faith"
            Reference: "Romans 4:20-21"
            Tags: ["Unbelief", "Faith", "Lord God", "Promise"]
          - ID: "E11"
            Title: "Good Works"
            Reference: "Galatians 6:9-10"
            Tags: ["Good", "Give Up", "Believers", "Do Good"]
          - ID: "E12"
            Title: "Good Works"
            Reference: "Matthew 5:16"
            Tags: ["Light", "Shine", "Good", "Praise"]
  - ID: "LOA"
    Title: "Lessons of Assurance"
    Packs:
      - ID: "LOA"
        Title: "Lessons of Assurance"
        Verses:
          - ID: "LOA1"
            Title: "Assurance of Salvation"
            Reference: "1 John 5:11-12"
            Tags: ["Testimony", "Eternal Life", "Son", "Lord God"]
          - ID: "LOA2"
            Title: "Assurance of Answered Prayer"
            Reference: "John 16:24"
            Tags: ["Jesus's Name", "Ask", "Joy", "Receive"]
          - ID: "LOA3"
            Title: "Assurance of Victory"
            Reference: "1 Corinthians 10:13"
            Tags: ["Temptation", "Common", "Lord God", "Tempted"]
          - ID: "LOA4"
            Title: "Assurance of Forgiveness"
            Reference: "1 John 1:9"
            Tags: ["Sin", "Confess", "Faithful", "Forgive"]
          - ID: "LOA5"
            Title: "Assurance of Guidance"
            Reference: "Proverbs 3:5-6"
            Tags: ["Trust", "Lord God", "Way", "Path"]
  - ID: "Life Issues"
    Title: "Life Issues"
    Packs:
      - ID: "Sin"
        Title: "Life Issues: Sin"
        Verses:
          - ID: "Sin1"
            Title: "Sin"
            Reference: "Romans 6:11-13"
            Tags: ["Dead", "Alive", "Reign", "Obey", "Righteousness"]
          - ID: "Sin2"
            Title: "Sin"
            Reference: "1 Corinthians 10:13"
            Tags: ["Temptation", "Faithful", "Tempted", "Victory"]
          - ID: "Sin3"
            Title: "Sin"
            Reference: "Galatians 6:1-2"
            Tags: ["Transgression", "Gentleness", "Tempted", "Burdens"]
          - ID: "Sin4"
            Title: "Sin"
            Reference: "Ephesians 6:10-12"
            Tags: ["Strong", "Armor", "Evil", "Spiritual"]
          - ID: "Sin5"
            Title: "Sin"
            Reference: "James 4:7-8"
            Tags: ["Submit", "Resist", "Devil", "Flee", "Cleanse"]
          - ID: "Sin6"
            Title: "Sin"
            Reference: "1 John 1:8"
            Tags: ["Deceive", "Truth", "Blameless", "Say"]
          - ID: "Sin7"
            Title: "Sin"
            Reference: "1 John 1:9"
            Tags: ["Confess", "Faithful", "Just", "Forgive", "Cleanse"]
      - ID: "Depression"
        Title: "Life Issues: Depression"
        Verses:
          - ID: "Depression1"
            Title: "Depression"
            Reference: "Isaiah 43:1-3"
            Tags: ["Fear", "With", "Overwhelm", "Consume"]
          - ID: "Depression2"
            Title: "Depression"
            Reference: "2 Corinthians 4:7-10"
            Tags: ["Treasure", "Power", "Affliction", "Crushed", "Life"]
          - ID: "Depression3"
            Title: "Depression"
            Reference: "Psalm 42:5"
            Tags: ["Turmoil", "Soul", "Hope", "Salvation"]
          - ID: "Depression4"
            Title: "Depression"
            Reference: "Psalm 34:17-18"
            Tags: ["Righteous", "Lord God", "Delivers", "Brokenhearted", "Save"]
          - ID: "Depression5"
            Title: "Depression"
            Reference: "Lamentations 3:19-23"
            Tags: ["Affliction", "Soul", "Steadfast", "Love", "Lord God"]
          - ID: "Depression6"
            Title: "Depression"
            Reference: "2 Corinthians 1:8-9"
            Tags: ["Affliction", "Burden", "Rely", "Lord God"]
      - ID: "Guilt"
        Title: "Life Issues: Guilt"
        Verses:
          - ID: "Guilt1"
            Title: "Guilt"
            Reference: "Psalm 32:1-2"
            Tags: ["Blessing", "Sin", "Cover", "Spirit"]
          - ID: "Guilt2"
            Title: "Guilt"
            Reference: "Psalm 51:9-10"
            Tags: ["Hide", "Sin", "Clean", "Heart", "Spirit"]
          - ID: "Guilt3"
            Title: "Guilt"
            Reference: "Proverbs 28:13"
            Tags: ["Conceal", "Transgression", "Confess", "Forsake"]
          - ID: "Guilt4"
            Title: "Guilt"
            Reference: "Romans 8:1-2"
            Tags: ["Condemn", "Christ Jesus", "Law", "Free"]
          - ID: "Guilt5"
            Title: "Guilt"
            Reference: "2 Corinthians 7:10"
            Tags: ["Grief", "Repentance", "Salvation", "Death"]
          - ID: "Guilt6"
            Title: "Guilt"
            Reference: "James 5:16"
            Tags: ["Confess", "Sin", "Prayer", "Righteous"]
      - ID: "God's Will"
        Title: "Life Issues: God's Will"
        Verses:
          - ID: "Will1"
            Title: "God's Will"
            Reference: "Proverbs 3:5-6"
            Tags: ["Trust", "Heart", "Understanding", "Path"]
          - ID: "Will2"
            Title: "God's Will"
            Reference: "Proverbs 3:7"
            Tags: ["Wisdom", "Own", "Fear", "Turn", "Evil"]
          - ID: "Will3"
            Title: "God's Will"
            Reference: "Proverbs 16:9"
            Tags: ["Heart", "Man", "Plans", "Step"]
          - ID: "Will4"
            Title: "God's Will"
            Reference: "Isaiah 30:21"
            Tags: ["Hear", "Word", "Walk", "Turn"]
          - ID: "Will5"
            Title: "God's Will"
            Reference: "Jeremiah 29:11-13"
            Tags: ["Plans", "Prosper", "Harm", "Call", "Hear"]
          - ID: "Will6"
            Title: "God's Will"
            Reference: "Romans 12:1-2"
            Tags: ["Mercy", "Offer", "Sacrifice", "Transform"]
          - ID: "Will7"
            Title: "God's Will"
            Reference: "1 John 5:14-15"
            Tags: ["Confidence", "Ask", "Hear", "Know"]
      - ID: "Love"
        Title: "Life Issues: Love"
        Verses:
          - ID: "Love1"
            Title: "Love"
            Reference: "Matthew 22:37-40"
            Tags: ["Lord God", "Heart", "Soul", "Mind", "Yourself"]
          - ID: "Love2"
            Title: "Love"
            Reference: "John 13:34-35"
            Tags: ["Command", "One Another", "Disciples", "Jesus Christ"]
          - ID: "Love3"
            Title: "Love"
            Reference: "Romans 8:38-39"
            Tags: ["Death", "Life", "Separate", "Jesus Christ"]
          - ID: "Love4"
            Title: "Love"
            Reference: "1 Corinthians 13:1-3"
            Tags: ["Tongues", "Prophecy", "Knowledge", "Faith"]
          - ID: "Love5"
            Title: "Love"
            Reference: "1 Corinthians 13:4-8"
            Tags: ["Patient", "Kind", "Envy", "Boasting", "Arrogance", "Rude", "Truth", "Irritation", "Resentment", "Hope", "Endures"]
          - ID: "Love6"
            Title: "Love"
            Reference: "1 John 4:20"
            Tags: ["Hate", "Brother", "Liar", "Lord God"]
      - ID: "Money"
        Title: "Life Issues: Money"
        Verses:
          - ID: "Money1"
            Title: "Money"
            Reference: "Deuteronomy 8:17-18"
            Tags: ["Heart", "Power", "Lord God", "Give"]
          - ID: "Money2"
            Title: "Money"
            Reference: "Proverbs 3:9-10"
            Tags: ["Lord God", "First", "Produce", "Plenty"]
          - ID: "Money3"
            Title: "Money"
            Reference: "Matthew 6:19-21"
            Tags: ["Hoard", "Earth", "Heaven", "Heart"]
          - ID: "Money4"
            Title: "Money"
            Reference: "Matthew 6:24"
            Tags: ["Master", "Hate", "Love", "Lord God"]
          - ID: "Money5"
            Title: "Money"
            Reference: "Philippians 4:11-13"
            Tags: ["Need", "Content", "Plenty", "Jesus Christ"]
          - ID: "Money6"
            Title: "Money"
            Reference: "1 Timothy 6:9-10"
            Tags: ["Desire", "Temptation", "Destruction", "Evil"]
      - ID: "Perfectionism"
        Title: "Life Issues: Perfectionism"
        Verses:
          - ID: "Perfectionism1"
            Title: "Perfectionism"
            Reference: "Psalm 127:1-2"
            Tags: ["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]
          - ID: "Perfectionism2"
            Title: "Perfectionism"
            Reference: "Ecclesiastes 2:10-11"
            Tags: ["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]
          - ID: "Perfectionism3"
            Title: "Perfectionism"
            Reference: "Luke 10:40-42"
            Tags: ["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]
          - ID: "Perfectionism4"
            Title: "Perfectionism"
            Reference: "2 Corinthians 12:9"
            Tags: ["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]
          - ID: "Perfectionism5"
            Title: "Perfectionism"
            Reference: "Galatians 3:3"
            Tags: ["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]
          - ID: "Perfectionism6"
            Title: "Perfectionism"
            Reference: "Ephesians 2:8-9"
            Tags: ["Perfectionism", "Perfect", "Perfection", "Salvation", "Achievements", "Grace"]
      - ID: "Self-Image"
        Title: "Life Issues: Self-Image"
        Verses:
          - ID: "Self-Image1"
            Title: "Self-Image"
            Reference: "1 Samuel 16:7"
            Tags: ["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]
          - ID: "Self-Image2"
            Title: "Self-Image"
            Reference: "Psalm 139:13-14"
            Tags: ["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]
          - ID: "Self-Image3"
            Title: "Self-Image"
            Reference: "Jeremiah 9:23-24"
            Tags: ["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]
          - ID: "Self-Image4"
            Title: "Self-Image"
            Reference: "Matthew 10:29-31"
            Tags: ["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]
          - ID: "Self-Image5"
            Title: "Self-Image"
            Reference: "Philippians 2:3-11"
            Tags: ["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]
          - ID: "Self-Image6"
            Title: "Self-Image"
            Reference: "1 Peter 3:3-4"
            Tags: ["Self Image", "Self", "Self Worth", "Self Esteem", "Confidence", "Self Pity"]
      - ID: "Sex"
        Title: "Life Issues: Sex"
        Verses:
          - ID: "Sex1"
            Title: "Sex"
            Reference: "Matthew 5:27-28"
            Tags: ["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]
          - ID: "Sex2"
            Title: "Sex"
            Reference: "Romans 13:13-14"
            Tags: ["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]
          - ID: "Sex3"
            Title: "Sex"
            Reference: "1 Corinthians 6:18-20"
            Tags: ["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]
          - ID: "Sex4"
            Title: "Sex"
            Reference: "Ephesians 5:3"
            Tags: ["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]
          - ID: "Sex5"
            Title: "Sex"
            Reference: "1 Thessalonians 4:3-5"
            Tags: ["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]
          - ID: "Sex6"
            Title: "Sex"
            Reference: "Hebrews 13:4"
            Tags: ["Sex", "Sexuality", "Body", "Intimacy", "Physical Intimacy", "Intimate"]
      - ID: "Stress"
        Title: "Life Issues: Stress"
        Verses:
          - ID: "Stress1"
            Title: "Stress"
            Reference: "Psalm 73:26"
            Tags: ["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]
          - ID: "Stress2"
            Title: "Stress"
            Reference: "Psalm 118:5-6"
            Tags: ["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]
          - ID: "Stress3"
            Title: "Stress"
            Reference: "Matthew 11:28-30"
            Tags: ["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]
          - ID: "Stress4"
            Title: "Stress"
            Reference: "2 Corinthians 4:16-18"
            Tags: ["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]
          - ID: "Stress5"
            Title: "Stress"
            Reference: "Philippians 4:6-7"
            Tags: ["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]
          - ID: "Stress6"
            Title: "Stress"
            Reference: "1 Peter 5:5-7"
            Tags: ["Stress", "Stressful", "Pressure", "Burden", "Worry", "Anxiety"]
      - ID: "Suffering"
        Title: "Life Issues: Suffering"
        Verses:
          - ID: "Suffering1"
            Title: "Suffering"
            Reference: "Romans 5:2-5"
            Tags: ["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]
          - ID: "Suffering2"
            Title: "Suffering"
            Reference: "2 Corinthians 1:3-4"
            Tags: ["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]
          - ID: "Suffering3"
            Title: "Suffering"
            Reference: "James 1:2-4"
            Tags: ["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]
          - ID: "Suffering4"
            Title: "Suffering"
            Reference: "James 1:12"
            Tags: ["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]
          - ID: "Suffering5"
            Title: "Suffering"
            Reference: "1 Peter 1:6-7"
            Tags: ["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]
          - ID: "Suffering6"
            Title: "Suffering"
            Reference: "1 Peter 4:12-13"
            Tags: ["Suffering", "Suffer", "Pain", "Enduring", "Agony", "Comfort"]
"""

data = yaml.safe_load(yaml_content)

print("-- Create verse_packs table")
print("""CREATE TABLE verse_packs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- NULL for system packs
    title VARCHAR(100) NOT NULL,
    identifier VARCHAR(50), -- Increased to 50 for safety
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);""")

print("CREATE INDEX idx_verse_packs_user_id ON verse_packs(user_id);")

print("-- Modify memory_verses table")
print("DELETE FROM memory_verses; -- Remove old incompatible data")
print("ALTER TABLE memory_verses ADD COLUMN verse_pack_id UUID REFERENCES verse_packs(id) ON DELETE CASCADE;")
print("ALTER TABLE memory_verses DROP COLUMN pack_name;")
print("ALTER TABLE memory_verses DROP COLUMN text;")
print("ALTER TABLE memory_verses DROP COLUMN user_id;")
print("ALTER TABLE memory_verses ALTER COLUMN verse_pack_id SET NOT NULL;")

print("-- Modify group_shares table")
print("ALTER TABLE group_shares ADD COLUMN verse_pack_id UUID REFERENCES verse_packs(id) ON DELETE CASCADE;")
print("ALTER TABLE group_shares ALTER COLUMN note_id DROP NOT NULL;")
print("ALTER TABLE group_shares ADD CONSTRAINT chk_share_target CHECK ( (note_id IS NOT NULL AND verse_pack_id IS NULL) OR (note_id IS NULL AND verse_pack_id IS NOT NULL) );")

print("-- Seed System Packs and Verses")
print("DO 1321")
print("DECLARE")
print("    pack_id UUID;")
print("BEGIN")

for series in data['Series']:
    series_title = series['Title']
    for pack in series['Packs']:
        pack_title = pack['Title']
        pack_id_str = pack['ID']
        # Escape quotes in title
        safe_pack_title = pack_title.replace("'", "''")
        safe_series_title = series_title.replace("'", "''")

        # Use description to store Series info
        print(f"    -- Pack: {pack_title}")
        print(f"    INSERT INTO verse_packs (id, title, identifier, description, is_public) VALUES (gen_random_uuid(), '{safe_pack_title}', '{pack_id_str}', '{safe_series_title}', true) RETURNING id INTO pack_id;")

        for verse in pack['Verses']:
            ref = verse['Reference']
            tags = verse['Tags']
            # Convert list to PG array/json string
            tags_json = str(tags).replace("'", '"')

            print(f"    INSERT INTO memory_verses (verse_pack_id, reference, version, tags) VALUES (pack_id, '{ref}', 'ESV', '{tags_json}');")

print("END 1321;")
