package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Book struct {
	Name     string
	Chapters int
}

// Protestant Canon
var (
	Genesis         = Book{"Genesis", 50}
	Exodus          = Book{"Exodus", 40}
	Leviticus       = Book{"Leviticus", 27}
	Numbers         = Book{"Numbers", 36}
	Deuteronomy     = Book{"Deuteronomy", 34}
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
	Job             = Book{"Job", 42}
	Psalms          = Book{"Psalm", 150} // Using "Psalm" singular for reference building
	Proverbs        = Book{"Proverbs", 31}
	Ecclesiastes    = Book{"Ecclesiastes", 12}
	SongOfSolomon   = Book{"Song of Solomon", 8}
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
	Matthew         = Book{"Matthew", 28}
	Mark            = Book{"Mark", 16}
	Luke            = Book{"Luke", 24}
	John            = Book{"John", 21}
	Acts            = Book{"Acts", 28}
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
)

var AllBooks = []Book{
	Genesis, Exodus, Leviticus, Numbers, Deuteronomy,
	Joshua, Judges, Ruth, ISamuel, IISamuel, IKings, IIKings, IChronicles, IIChronicles, Ezra, Nehemiah, Esther,
	Job, Psalms, Proverbs, Ecclesiastes, SongOfSolomon,
	Isaiah, Jeremiah, Lamentations, Ezekiel, Daniel,
	Hosea, Joel, Amos, Obadiah, Jonah, Micah, Nahum, Habakkuk, Zephaniah, Haggai, Zechariah, Malachi,
	Matthew, Mark, Luke, John, Acts,
	Romans, ICorinthians, IICorinthians, Galatians, Ephesians, Philippians, Colossians,
	IThessalonians, IIThessalonians, ITimothy, IITimothy, Titus, Philemon,
	Hebrews, James, IPeter, IIPeter, IJohn, IIJohn, IIIJohn, Jude, Revelation,
}

var OTBooks = []Book{
	Genesis, Exodus, Leviticus, Numbers, Deuteronomy,
	Joshua, Judges, Ruth, ISamuel, IISamuel, IKings, IIKings, IChronicles, IIChronicles, Ezra, Nehemiah, Esther,
	Job, Psalms, Proverbs, Ecclesiastes, SongOfSolomon,
	Isaiah, Jeremiah, Lamentations, Ezekiel, Daniel,
	Hosea, Joel, Amos, Obadiah, Jonah, Micah, Nahum, Habakkuk, Zephaniah, Haggai, Zechariah, Malachi,
}

func main() {
	var upSQL strings.Builder
	var downSQL strings.Builder

	// 1. Whole Bible in 3 Years
	generateWholeBiblePlan(&upSQL, &downSQL)

	// 2. Old Testament in a Year
	generateOTPlan(&upSQL, &downSQL)

	// 3. Psalms & Proverbs
	generatePsalmsProverbsPlan(&upSQL, &downSQL)

	if err := os.WriteFile("../../api/migrations/000023_add_algorithmic_plans.up.sql", []byte(upSQL.String()), 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("../../api/migrations/000023_add_algorithmic_plans.down.sql", []byte(downSQL.String()), 0644); err != nil {
		panic(err)
	}

	fmt.Println("Migration files generated successfully.")
}

func generateWholeBiblePlan(up *strings.Builder, down *strings.Builder) {
	planID := uuid.New()
	title := "Whole Bible (3 Years)"
	description := "Read through the entire Bible at a pace of one chapter per day. Takes about 3 years and 3 months."
	planType := "sequential"

	allChapters := []string{}
	for _, book := range AllBooks {
		for c := 1; c <= book.Chapters; c++ {
			allChapters = append(allChapters, fmt.Sprintf("%s %d", book.Name, c))
		}
	}
	days := len(allChapters)

	up.WriteString(fmt.Sprintf("-- %s\n", title))
	up.WriteString(fmt.Sprintf("INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at) VALUES ('%s', '%s', '%s', %d, '%s', NOW(), NOW());\n", planID, escapeSQL(title), escapeSQL(description), days, planType))
	down.WriteString(fmt.Sprintf("DELETE FROM reading_plans WHERE id = '%s';\n", planID))

	for i, passage := range allChapters {
		dayNum := i + 1
		dayID := uuid.New()
		up.WriteString(fmt.Sprintf("INSERT INTO reading_plan_days (id, reading_plan_id, day_number, passage, created_at) VALUES ('%s', '%s', %d, '%s', NOW());\n", dayID, planID, dayNum, escapeSQL(passage)))
	}
	up.WriteString("\n")
}

func generateOTPlan(up *strings.Builder, down *strings.Builder) {
	planID := uuid.New()
	title := "Old Testament in a Year"
	description := "Read through the Old Testament in one year."
	planType := "calendar"
	days := 365

	allOTChapters := []string{}
	for _, book := range OTBooks {
		for c := 1; c <= book.Chapters; c++ {
			allOTChapters = append(allOTChapters, fmt.Sprintf("%s %d", book.Name, c))
		}
	}
	totalChapters := len(allOTChapters) // Should be 929

	up.WriteString(fmt.Sprintf("-- %s\n", title))
	up.WriteString(fmt.Sprintf("INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at) VALUES ('%s', '%s', '%s', %d, '%s', NOW(), NOW());\n", planID, escapeSQL(title), escapeSQL(description), days, planType))
	down.WriteString(fmt.Sprintf("DELETE FROM reading_plans WHERE id = '%s';\n", planID))

	chaptersAssigned := 0
	for d := 1; d <= days; d++ {
		// Calculate how many chapters we should have covered by end of this day
		targetTotal := int(math.Round(float64(d) * float64(totalChapters) / float64(days)))
		chaptersTodayCount := targetTotal - chaptersAssigned

		var todaysReadings []string
		for k := 0; k < chaptersTodayCount; k++ {
			if chaptersAssigned < totalChapters {
				todaysReadings = append(todaysReadings, allOTChapters[chaptersAssigned])
				chaptersAssigned++
			}
		}

		passage := strings.Join(todaysReadings, "; ")
		dayID := uuid.New()
		up.WriteString(fmt.Sprintf("INSERT INTO reading_plan_days (id, reading_plan_id, day_number, passage, created_at) VALUES ('%s', '%s', %d, '%s', NOW());\n", dayID, planID, d, escapeSQL(passage)))
	}
	up.WriteString("\n")
}

func generatePsalmsProverbsPlan(up *strings.Builder, down *strings.Builder) {
	planID := uuid.New()
	title := "Psalms & Proverbs"
	description := "Read through Psalms and Proverbs twice in a year."
	planType := "calendar"
	days := 365

	up.WriteString(fmt.Sprintf("-- %s\n", title))
	up.WriteString(fmt.Sprintf("INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at) VALUES ('%s', '%s', '%s', %d, '%s', NOW(), NOW());\n", planID, escapeSQL(title), escapeSQL(description), days, planType))
	down.WriteString(fmt.Sprintf("DELETE FROM reading_plans WHERE id = '%s';\n", planID))

	for d := 1; d <= days; d++ {
		psalmIdx := (d - 1) % 150
		proverbIdx := (d - 1) % 31

		psalm := fmt.Sprintf("Psalm %d", psalmIdx+1)
		proverb := fmt.Sprintf("Proverbs %d", proverbIdx+1)

		passage := fmt.Sprintf("%s; %s", psalm, proverb)

		dayID := uuid.New()
		up.WriteString(fmt.Sprintf("INSERT INTO reading_plan_days (id, reading_plan_id, day_number, passage, created_at) VALUES ('%s', '%s', %d, '%s', NOW());\n", dayID, planID, d, escapeSQL(passage)))
	}
	up.WriteString("\n")
}

func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
