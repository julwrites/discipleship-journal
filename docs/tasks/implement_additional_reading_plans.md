# Task: Implement Additional Well-Established Reading Plans

## Objective
Implement scripts or logic to seed the database with additional well-established Bible reading plans. This expands the options available to users beyond the basic algorithmic plans.

## 1. Chronological Plan (1 Year)
*   **Description**: Read the Bible in the chronological order in which its stories and events occurred.
*   **Type**: Calendar (365 Days)
*   **Method**: Scraping.
*   **Source**: BibleGateway "Chronological"
*   **URL Pattern**: `https://www.biblegateway.com/reading-plans/chronological/{year}/{month}/{day}?version=ESV`
    *   Example: `https://www.biblegateway.com/reading-plans/chronological/2026/01/01?version=ESV`
*   **Scraping Logic**:
    *   Iterate dates from `2026-01-01` to `2026-12-31`.
    *   Fetch HTML.
    *   Extract the passage reference. It is typically found in the breadcrumb (e.g., `Chronological / Genesis 1-3`) or the header.
    *   Parse and clean the reference string.

## 2. Bible in 90 Days
*   **Description**: An intensive Bible reading plan that walks through the entire Bible in 90 days.
*   **Type**: Sequential (88-90 Days)
*   **Method**: Scraping.
*   **Source**: BibleGateway "Bible in 90 Days"
*   **URL Pattern**: `https://www.biblegateway.com/reading-plans/bible-in-90-days/{year}/{month}/{day}?version=ESV`
    *   Example: `https://www.biblegateway.com/reading-plans/bible-in-90-days/2026/01/01?version=ESV`
*   **Scraping Logic**:
    *   Iterate dates from `2026-01-01` to `2026-03-31` (approx 88-90 days).
    *   Extract passage reference from breadcrumb (e.g., `Bible in 90 Days / Genesis 1:1 - Genesis 16:16`).
    *   Note: The duration listed on the site says "88 days", so the loop should check for completion or 404s.

## 3. Professor Grant Horner's Bible Reading System
*   **Description**: A system consisting of 10 lists of books. You read one chapter from each list every day (10 chapters/day). The lists vary in length, causing the readings to rotate and interweave constantly.
*   **Type**: Algorithmic / Generated Sequence (e.g., Generate 365 or 730 days of the sequence).
*   **Method**: Algorithmic Generation.
*   **Logic**:
    *   Define 10 lists of books/chapters.
    *   Day `n` reading = `(List1[n % len1])` + `(List2[n % len2])` + ...
    *   Generate a `reading_plan` entry (e.g., "Professor Grant Horner's System (1 Year Sample)") and insert `reading_plan_days`.
*   **Lists** (Based on common variation):
    1.  **Gospels**: Matthew, Mark, Luke, John
    2.  **Pentateuch**: Genesis, Exodus, Leviticus, Numbers, Deuteronomy
    3.  **Romans & Pauline Epistles (Part 1)**: Romans, 1 Corinthians, 2 Corinthians, Galatians, Ephesians, Philippians, Colossians
    4.  **Pauline (Part 2) & General Epistles**: 1 Thessalonians, 2 Thessalonians, 1 Timothy, 2 Timothy, Titus, Philemon, Hebrews, James, 1 Peter, 2 Peter, 1 John, 2 John, 3 John, Jude, Revelation
    5.  **Wisdom**: Job, Ecclesiastes, Song of Solomon
    6.  **Psalms**: Psalms
    7.  **Proverbs**: Proverbs
    8.  **History**: Joshua, Judges, Ruth, 1 Samuel, 2 Samuel, 1 Kings, 2 Kings, 1 Chronicles, 2 Chronicles, Ezra, Nehemiah, Esther
    9.  **Prophets**: Isaiah, Jeremiah, Lamentations, Ezekiel, Daniel, Hosea, Joel, Amos, Obadiah, Jonah, Micah, Nahum, Habakkuk, Zephaniah, Haggai, Zechariah, Malachi
    10. **Acts**: Acts

## 4. Discipleship Journal Bible Reading Plan
*   **Description**: Four daily readings starting in: Genesis, Psalms, Matthew, and Acts. 25 readings per month to allow catch-up days.
*   **Type**: Calendar (365 days, with rest days)
*   **Method**: Manual Entry or Scraping if source found.
*   **Source**: The Navigators / Discipleship Journal.
*   **Strategy**:
    *   Since this plan uses "Day 1" to "Day 25" for each month, it maps well to a 365-day plan where days 26-End of month might be empty or skipped in a sequential implementation.
    *   Alternatively, populate a JSON file manually from the PDF source (`navigators.org`) if no clean digital list is available.
    *   Structure:
        *   Stream 1: Law/History
        *   Stream 2: Psalms/Poetry
        *   Stream 3: Gospels
        *   Stream 4: Epistles
