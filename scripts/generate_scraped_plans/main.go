package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReadingPlanConfig struct {
	Title       string
	Slug        string // URL slug e.g. "chronological"
	Description string
	Type        string
	MaxDays     int
	MigrationID string // e.g. "000024_add_chronological_plan"
}

const (
	baseURL = "https://www.biblegateway.com/reading-plans"
	year    = 2026
	delay   = 200 * time.Millisecond // Reduced delay slightly, but still respectful
)

var configs = []ReadingPlanConfig{
	{
		Title:       "Chronological",
		Slug:        "chronological",
		Description: "Read the Bible in the chronological order in which its stories and events occurred.",
		Type:        "calendar",
		MaxDays:     365,
		MigrationID: "000024_add_chronological_plan",
	},
	{
		Title:       "Bible in 90 Days",
		Slug:        "bible-in-90-days",
		Description: "An intensive Bible reading plan that walks through the entire Bible in 90 days.",
		Type:        "sequential",
		MaxDays:     100, // It's around 88-90 days, we'll stop when we see repetition
		MigrationID: "000025_add_bible_in_90_days",
	},
	{
		Title:       "New Testament in a Year",
		Slug:        "new-testament-in-a-year",
		Description: "Read through the entire New Testament in one year.",
		Type:        "calendar",
		MaxDays:     365,
		MigrationID: "000026_add_new_testament_in_a_year",
	},
}

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, config := range configs {
		fmt.Printf("Generating plan: %s...\n", config.Title)
		generatePlan(client, config)
	}

	fmt.Println("All plans generated successfully!")
}

func generatePlan(client *http.Client, config ReadingPlanConfig) {
	upSQL := strings.Builder{}
	downSQL := strings.Builder{}

	planID := uuid.New()

	// Header for UP migration
	upSQL.WriteString(fmt.Sprintf("-- %s Reading Plan Seed\n", config.Title))
	upSQL.WriteString(fmt.Sprintf("INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at) VALUES ('%s', '%s', '%s', %d, '%s', NOW(), NOW());\n", planID, escapeSQL(config.Title), escapeSQL(config.Description), 0, config.Type)) // Days 0 placeholder, will update later? No, usually days is fixed or derived. I'll use a placeholder or count.
	// Actually, I should probably count the days and update the record, or just insert the count if I know it.
	// For variable length plans, I'll count as I go and might need to update the migration file or use a variable.
	// But SQL is sequential. I can't update the previous INSERT easily in a streaming write.
	// However, for "Bible in 90 Days", I don't know the exact count yet.
	// I'll store the days in memory and write the SQL at the end of the function.

	// Better approach: Buffer day inserts, then write plan insert with correct count, then day inserts.

	dayInserts := strings.Builder{}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	firstPassage := ""
	daysCount := 0

	for day := 1; day <= config.MaxDays; day++ {
		date := startDate.AddDate(0, 0, day-1)
		url := fmt.Sprintf("%s/%s/%d/%02d/%02d?version=ESV", baseURL, config.Slug, date.Year(), date.Month(), date.Day())

		fmt.Printf("  Fetching day %d: %s... ", day, url)

		passage, err := fetchPassage(client, url)
		if err != nil {
			fmt.Printf("Error: %v. Stopping.\n", err)
			break
		}

		// Check for repetition (End of cycle)
		if day == 1 {
			firstPassage = passage
		} else if passage == firstPassage {
			fmt.Printf("Passage repeated (%s). End of plan.\n", passage)
			break
		}

		fmt.Printf("Found: %s\n", passage)

		// Clean up passage
		passage = cleanPassage(passage)

		dayID := uuid.New()
		dayInserts.WriteString(fmt.Sprintf("INSERT INTO reading_plan_days (id, reading_plan_id, day_number, passage, created_at) VALUES ('%s', '%s', %d, '%s', NOW());\n", dayID, planID, day, escapeSQL(passage)))
		daysCount++

		time.Sleep(delay)
	}

	// Re-write the plan insert with correct days count
	finalUpSQL := strings.Builder{}
	finalUpSQL.WriteString(fmt.Sprintf("-- %s Reading Plan Seed\n", config.Title))
	finalUpSQL.WriteString(fmt.Sprintf("INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at) VALUES ('%s', '%s', '%s', %d, '%s', NOW(), NOW());\n", planID, escapeSQL(config.Title), escapeSQL(config.Description), daysCount, config.Type))
	finalUpSQL.WriteString(dayInserts.String())

	// Down migration
	downSQL.WriteString(fmt.Sprintf("-- Remove %s Reading Plan\n", config.Title))
	downSQL.WriteString(fmt.Sprintf("DELETE FROM reading_plan_days WHERE reading_plan_id = '%s';\n", planID))
	downSQL.WriteString(fmt.Sprintf("DELETE FROM reading_plans WHERE id = '%s';\n", planID))

	// Write files
	upFile := fmt.Sprintf("../../api/migrations/%s.up.sql", config.MigrationID)
	downFile := fmt.Sprintf("../../api/migrations/%s.down.sql", config.MigrationID)

	if err := os.WriteFile(upFile, []byte(finalUpSQL.String()), 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(downFile, []byte(downSQL.String()), 0644); err != nil {
		panic(err)
	}
}

func fetchPassage(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	html := string(body)

	// Look for title="Listen to {Reference}"
	// <a class="audio-link" ... title="Listen to Genesis 1:1 - Genesis 16:16" >

	// Find class="audio-link"
	idx := strings.Index(html, "class=\"audio-link\"")
	if idx == -1 {
		// Try without quotes? Browsers might normalize, but source usually has them.
		// Or try single quotes?
		idx = strings.Index(html, "class='audio-link'")
	}

	if idx == -1 {
		return "", fmt.Errorf("audio-link class not found")
	}

	// Look for title attribute after audio-link
	// It should be within the <a> tag.
	// Let's find the closing > of the tag
	tagEnd := strings.Index(html[idx:], ">")
	if tagEnd == -1 {
		return "", fmt.Errorf("tag end not found")
	}

	tagContent := html[idx : idx+tagEnd]

	// Look for title="
	titleIdx := strings.Index(tagContent, "title=\"")
	if titleIdx == -1 {
		return "", fmt.Errorf("title attribute not found in audio-link")
	}

	// Extract content
	start := titleIdx + 7 // len("title=\"")
	end := strings.Index(tagContent[start:], "\"")
	if end == -1 {
		return "", fmt.Errorf("title attribute closing quote not found")
	}

	titleValue := tagContent[start : start+end]

	// Remove "Listen to "
	if strings.HasPrefix(titleValue, "Listen to ") {
		return titleValue[10:], nil
	}

	return titleValue, nil
}

func cleanPassage(p string) string {
	p = strings.TrimSpace(p)
	// Normalize separators
	p = strings.ReplaceAll(p, "\n", "; ")
	p = strings.ReplaceAll(p, " + ", "; ")
	p = strings.ReplaceAll(p, " & ", "; ")

	// Enforce semicolon rule: if we see ", " followed by a Book Name (heuristic: capital letter),
	// and NOT a digit (chapter/verse), we warn or replace.
	// But automating this is risky without a book list.
	// For now, we trust the source mostly but ensure standard delimiters are semi-colons.

	// If the source uses commas to separate books (e.g. "Gen 1, Ex 1"), we might need to fix it manually or improve this.
	// But BibleGateway usually formats clearly.

	return p
}

func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
