# Task: Implement Additional Bible Reading Plans

## Objective
Implement scripts or logic to seed the database with four new Bible reading plans:
1.  **Whole Bible in 3 Years** (1 chapter/day)
2.  **Old Testament in a Year**
3.  **New Testament in a Year**
4.  **Psalms & Proverbs**

## Strategy & Sources

### 1. Whole Bible in 3 Years (1 Chapter/Day)
*   **Method**: Algorithmic Generation.
*   **Logic**:
    *   Iterate through the Protestant canon (Genesis to Revelation).
    *   One chapter per day.
    *   Total chapters: 1,189.
    *   Duration: ~3 years and 3 months.
*   **Implementation**: Create a Go script/function that utilizes a map/list of Bible book chapter counts.

### 2. Old Testament in a Year
*   **Method**: Algorithmic Generation.
*   **Logic**:
    *   Iterate through Old Testament books (Genesis to Malachi).
    *   Total chapters: 929.
    *   Target duration: 365 days.
    *   Daily rate: `math.Ceil(929 / 365)`. Some days will have 2 chapters, some 3.
    *   Alternatively, simple sequential chunks.

### 3. New Testament in a Year
*   **Method**: Scraping (Primary) or Algorithmic (Secondary).
*   **Source**: BibleGateway "New Testament in a Year".
*   **URL Pattern**: `https://www.biblegateway.com/reading-plans/new-testament-in-a-year/{year}/{month}/{day}?version=ESV`
    *   Example: `https://www.biblegateway.com/reading-plans/new-testament-in-a-year/2026/01/01?version=ESV`
*   **Scraping Logic**:
    *   Loop dates from `2026-01-01` to `2026-12-31`.
    *   Fetch HTML.
    *   Extract reference from the breadcrumb (e.g., `New Testament in a Year / Matthew 20:17-34`) or the header text.
    *   Clean the string to get just the reference (e.g., "Matthew 20:17-34").
*   **Note**: This plan is 5 days/week in some variations, or 365 days. BibleGateway's version appears to be 365 days.

### 4. Psalms & Proverbs
*   **Method**: Algorithmic Generation.
*   **Logic**:
    *   **Psalms**: 150 chapters.
    *   **Proverbs**: 31 chapters.
    *   **Daily Reading**:
        *   Cycle through Psalms (e.g., Psalm 1 on Day 1, Psalm 150 on Day 150, repeat).
        *   Cycle through Proverbs (e.g., Proverbs 1 on Day 1, Proverbs 31 on Day 31, repeat).
        *   Combine: "Psalm {X}; Proverbs {Y}".
    *   **Duration**: Indefinite (or set to 365 days).

## Output Format
Generate a SQL migration file (e.g., `api/migrations/YYYYMMDD_add_new_plans.up.sql`) containing:

```sql
-- Insert Plan Metadata
INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at)
VALUES (gen_random_uuid(), 'Title', 'Description', 365, 'calendar', NOW(), NOW());

-- Insert Days
INSERT INTO reading_plan_days (id, reading_plan_id, day_number, passage, created_at)
VALUES (gen_random_uuid(), 'PLAN_UUID', 1, 'Passage String', NOW());
```
