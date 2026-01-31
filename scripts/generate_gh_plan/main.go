package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Book struct {
	Name     string
	Chapters int
}

var (
	Matthew         = Book{"Matthew", 28}
	Mark            = Book{"Mark", 16}
	Luke            = Book{"Luke", 24}
	John            = Book{"John", 21}
	Genesis         = Book{"Genesis", 50}
	Exodus          = Book{"Exodus", 40}
	Leviticus       = Book{"Leviticus", 27}
	Numbers         = Book{"Numbers", 36}
	Deuteronomy     = Book{"Deuteronomy", 34}
	Romans          = Book{"Romans", 16}
	ICorinthians    = Book{"1 Corinthians", 16}
	IICorinthians   = Book{"2 Corinthians", 13}
	Galatians       = Book{"Galatians", 6}
	Ephesians       = Book{"Ephesians", 6}
	Philippians     = Book{"Philippians", 4}
	Colossians      = Book{"Colossians", 4}
	IThessalonians  = Book{"1 Thessalonians", 5}
	IIThessalonians = Book{"2 Thessalonians", 3}
	ITimothy        = Book{"1 Timothy", 6}
	IITimothy       = Book{"2 Timothy", 4}
	Titus           = Book{"Titus", 3}
	Philemon        = Book{"Philemon", 1}
	Hebrews         = Book{"Hebrews", 13}
	James           = Book{"James", 5}
	IPeter          = Book{"1 Peter", 5}
	IIPeter         = Book{"2 Peter", 3}
	IJohn           = Book{"1 John", 5}
	IIJohn          = Book{"2 John", 1}
	IIIJohn         = Book{"3 John", 1}
	Jude            = Book{"Jude", 1}
	Revelation      = Book{"Revelation", 22}
	Job             = Book{"Job", 42}
	Ecclesiastes    = Book{"Ecclesiastes", 12}
	SongOfSolomon   = Book{"Song of Solomon", 8}
	Psalms          = Book{"Psalm", 150} // Using "Psalm" singular for reference usually, but often "Psalms" is used for the book. Checking M'Cheyne: "Psalm 1-2".
	Proverbs        = Book{"Proverbs", 31}
	Joshua          = Book{"Joshua", 24}
	Judges          = Book{"Judges", 21}
	Ruth            = Book{"Ruth", 4}
	ISamuel         = Book{"1 Samuel", 31}
	IISamuel        = Book{"2 Samuel", 24}
	IKings          = Book{"1 Kings", 22}
	IIKings         = Book{"2 Kings", 25}
	IChronicles     = Book{"1 Chronicles", 29}
	IIChronicles    = Book{"2 Chronicles", 36}
	Ezra            = Book{"Ezra", 10}
	Nehemiah        = Book{"Nehemiah", 13}
	Esther          = Book{"Esther", 10}
	Isaiah          = Book{"Isaiah", 66}
	Jeremiah        = Book{"Jeremiah", 52}
	Lamentations    = Book{"Lamentations", 5}
	Ezekiel         = Book{"Ezekiel", 48}
	Daniel          = Book{"Daniel", 12}
	Hosea           = Book{"Hosea", 14}
	Joel            = Book{"Joel", 3}
	Amos            = Book{"Amos", 9}
	Obadiah         = Book{"Obadiah", 1}
	Jonah           = Book{"Jonah", 4}
	Micah           = Book{"Micah", 7}
	Nahum           = Book{"Nahum", 3}
	Habakkuk        = Book{"Habakkuk", 3}
	Zephaniah       = Book{"Zephaniah", 3}
	Haggai          = Book{"Haggai", 2}
	Zechariah       = Book{"Zechariah", 14}
	Malachi         = Book{"Malachi", 4}
	Acts            = Book{"Acts", 28}
)

func main() {
	lists := [][]Book{
		{Matthew, Mark, Luke, John},
		{Genesis, Exodus, Leviticus, Numbers, Deuteronomy},
		{Romans, ICorinthians, IICorinthians, Galatians, Ephesians, Philippians, Colossians},
		{IThessalonians, IIThessalonians, ITimothy, IITimothy, Titus, Philemon, Hebrews, James, IPeter, IIPeter, IJohn, IIJohn, IIIJohn, Jude, Revelation},
		{Job, Ecclesiastes, SongOfSolomon},
		{Psalms},
		{Proverbs},
		{Joshua, Judges, Ruth, ISamuel, IISamuel, IKings, IIKings, IChronicles, IIChronicles, Ezra, Nehemiah, Esther},
		{Isaiah, Jeremiah, Lamentations, Ezekiel, Daniel, Hosea, Joel, Amos, Obadiah, Jonah, Micah, Nahum, Habakkuk, Zephaniah, Haggai, Zechariah, Malachi},
		{Acts},
	}

	// Expand lists into flattened chapters
	// each list becomes a slice of strings e.g. "Matthew 1", "Matthew 2", ...
	expandedLists := make([][]string, len(lists))
	for i, list := range lists {
		for _, book := range list {
			for c := 1; c <= book.Chapters; c++ {
				expandedLists[i] = append(expandedLists[i], fmt.Sprintf("%s %d", book.Name, c))
			}
		}
	}

	planID := uuid.New() // Generate a new UUID for the plan
	days := 365
	title := "Professor Grant Horner's System (1 Year Sample)"
	description := "A system consisting of 10 lists of books. You read one chapter from each list every day."

	var upSQL strings.Builder
	var downSQL strings.Builder

	// Header
	upSQL.WriteString(fmt.Sprintf("-- %s\n", title))
	upSQL.WriteString(fmt.Sprintf("INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at) VALUES ('%s', '%s', '%s', %d, 'sequential', NOW(), NOW());\n", planID, escapeSQL(title), escapeSQL(description), days))

	downSQL.WriteString(fmt.Sprintf("DELETE FROM reading_plan_days WHERE reading_plan_id = '%s';\n", planID)) // Delete children first
	downSQL.WriteString(fmt.Sprintf("DELETE FROM reading_plans WHERE id = '%s';\n", planID))

	for d := 1; d <= days; d++ {
		var readings []string
		for i, list := range expandedLists {
			// Get reading for this day from this list (looping)
			idx := (d - 1) % len(list)
			readings = append(readings, list[idx])
			_ = i // unused
		}
		passage := strings.Join(readings, "; ")

		dayID := uuid.New()
		upSQL.WriteString(fmt.Sprintf("INSERT INTO reading_plan_days (id, reading_plan_id, day_number, passage, created_at) VALUES ('%s', '%s', %d, '%s', NOW());\n", dayID, planID, d, escapeSQL(passage)))
	}

	if err := os.WriteFile("../../api/migrations/000022_add_grant_horner_plan.up.sql", []byte(upSQL.String()), 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("../../api/migrations/000022_add_grant_horner_plan.down.sql", []byte(downSQL.String()), 0644); err != nil {
		panic(err)
	}

	fmt.Println("Migration files generated successfully.")
}

func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
