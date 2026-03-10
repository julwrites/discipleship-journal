# Task: Implement Additional Well-Established Reading Plans

## Objective
Implement scripts or logic to seed the database with additional well-established Bible reading plans. This expands the options available to users beyond the basic algorithmic plans.

## 1. Chronological Plan (1 Year)
*   **Description**: Read the Bible in the chronological order in which its stories and events occurred.
*   **Type**: Calendar (365 Days)
*   **Method**: Scraping (Implemented in `scripts/generate_scraped_plans/main.go`).
*   **Status**: Completed. Migration `000024_add_chronological_plan` added.
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
*   **Method**: Scraping (Implemented in `scripts/generate_scraped_plans/main.go`).
*   **Status**: Completed. Migration `000025_add_bible_in_90_days` added.
*   **Source**: BibleGateway "Bible in 90 Days"
*   **URL Pattern**: `https://www.biblegateway.com/reading-plans/bible-in-90-days/{year}/{month}/{day}?version=ESV`
    *   Example: `https://www.biblegateway.com/reading-plans/bible-in-90-days/2026/01/01?version=ESV`
*   **Scraping Logic**:
    *   Iterate dates from `2026-01-01` to `2026-03-31` (approx 88-90 days).
    *   Extract passage reference from breadcrumb (e.g., `Bible in 90 Days / Genesis 1:1 - Genesis 16:16`).
    *   Note: The duration listed on the site says "88 days", so the loop should check for completion or 404s.

## 3. Professor Grant Horner's Bible Reading System (Implemented)
*   **Description**: A system consisting of 10 lists of books. You read one chapter from each list every day (10 chapters/day). The lists vary in length, causing the readings to rotate and interweave constantly.
*   **Type**: Sequential (365 Days Generated)
*   **Method**: Algorithmic Generation (Implemented in `scripts/generate_gh_plan/main.go`).
*   **Status**: Completed. Migration `000022_add_grant_horner_plan` added.

## 4. Discipleship Journal Bible Reading Plan
*   **Description**: Four daily readings starting in: Genesis, Psalms, Matthew, and Acts. 25 readings per month to allow catch-up days.
*   **Type**: Calendar (365 days, with rest days)
*   **Method**: Algorithmic Generation (Implemented in `scripts/generate_dj_plan/main.go`).
*   **Status**: Completed. Migration `000028_add_discipleship_journal_plan` added.
*   **Source**: The Navigators / Discipleship Journal (Adapted).
*   **Strategy**:
    *   Implemented a "Balanced 4-Stream" approach that distributes biblical books into 4 streams over 300 active days (25 days * 12 months).
    *   Structure:
        *   Stream 1: Law/History (Genesis - Esther)
        *   Stream 2: Psalms/Prophets (Job - Malachi)
        *   Stream 3: Gospels (Matthew - John, looped)
        *   Stream 4: Epistles (Acts - Revelation, looped)
