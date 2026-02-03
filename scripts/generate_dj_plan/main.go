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

var (
	Genesis      = Book{"Genesis", 50}
	Exodus       = Book{"Exodus", 40}
	Leviticus    = Book{"Leviticus", 27}
	Numbers      = Book{"Numbers", 36}
	Deuteronomy  = Book{"Deuteronomy", 34}
	Joshua       = Book{"Joshua", 24}
	Judges       = Book{"Judges", 21}
	Ruth         = Book{"Ruth", 4}
	ISamuel      = Book{"1 Samuel", 31}
	IISamuel     = Book{"2 Samuel", 24}
	IKings       = Book{"1 Kings", 22}
	IIKings      = Book{"2 Kings", 25}
	IChronicles  = Book{"1 Chronicles", 29}
	IIChronicles = Book{"2 Chronicles", 36}
	Ezra         = Book{"Ezra", 10}
	Nehemiah     = Book{"Nehemiah", 13}
	Esther       = Book{"Esther", 10}

	Job           = Book{"Job", 42}
	Psalms        = Book{"Psalm", 150}
	Proverbs      = Book{"Proverbs", 31}
	Ecclesiastes  = Book{"Ecclesiastes", 12}
	SongOfSolomon = Book{"Song of Solomon", 8}
	Isaiah        = Book{"Isaiah", 66}
	Jeremiah      = Book{"Jeremiah", 52}
	Lamentations  = Book{"Lamentations", 5}
	Ezekiel       = Book{"Ezekiel", 48}
	Daniel        = Book{"Daniel", 12}
	Hosea         = Book{"Hosea", 14}
	Joel          = Book{"Joel", 3}
	Amos          = Book{"Amos", 9}
	Obadiah       = Book{"Obadiah", 1}
	Jonah         = Book{"Jonah", 4}
	Micah         = Book{"Micah", 7}
	Nahum         = Book{"Nahum", 3}
	Habakkuk      = Book{"Habakkuk", 3}
	Zephaniah     = Book{"Zephaniah", 3}
	Haggai        = Book{"Haggai", 2}
	Zechariah     = Book{"Zechariah", 14}
	Malachi       = Book{"Malachi", 4}

	Matthew = Book{"Matthew", 28}
	Mark    = Book{"Mark", 16}
	Luke    = Book{"Luke", 24}
	John    = Book{"John", 21}

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

func main() {
	// 4 Streams
	stream1 := []Book{Genesis, Exodus, Leviticus, Numbers, Deuteronomy, Joshua, Judges, Ruth, ISamuel, IISamuel, IKings, IIKings, IChronicles, IIChronicles, Ezra, Nehemiah, Esther}
	stream2 := []Book{Job, Psalms, Proverbs, Ecclesiastes, SongOfSolomon, Isaiah, Jeremiah, Lamentations, Ezekiel, Daniel, Hosea, Joel, Amos, Obadiah, Jonah, Micah, Nahum, Habakkuk, Zephaniah, Haggai, Zechariah, Malachi}
	stream3 := []Book{Matthew, Mark, Luke, John}
	stream4 := []Book{Acts, Romans, ICorinthians, IICorinthians, Galatians, Ephesians, Philippians, Colossians, IThessalonians, IIThessalonians, ITimothy, IITimothy, Titus, Philemon, Hebrews, James, IPeter, IIPeter, IJohn, IIJohn, IIIJohn, Jude, Revelation}

	// 12 months * 25 days = 300 reading days
	targetDays := 300

	readings1 := distribute(expand(stream1), targetDays)
	readings2 := distribute(expand(stream2), targetDays)
	readings3 := distribute(expand(stream3), targetDays)
	readings4 := distribute(expand(stream4), targetDays)

	// Generate SQL
	planID := uuid.New()
	title := "Discipleship Journal Bible Reading Plan"
	description := "A balanced approach to reading the entire Bible in a year. Four daily readings (Law/History, Psalms/Prophets, Gospels, Epistles). 25 reading days per month allow for catch-up."
	planType := "calendar"
	totalDays := 365 // It is a calendar plan

	var upSQL strings.Builder
	var downSQL strings.Builder

	// Header
	upSQL.WriteString(fmt.Sprintf("-- %s\n", title))
	upSQL.WriteString(fmt.Sprintf("INSERT INTO reading_plans (id, title, description, days, plan_type, created_at, updated_at) VALUES ('%s', '%s', '%s', %d, '%s', NOW(), NOW());\n", planID, escapeSQL(title), escapeSQL(description), totalDays, planType))

	downSQL.WriteString(fmt.Sprintf("DELETE FROM reading_plan_days WHERE reading_plan_id = '%s';\n", planID))
	downSQL.WriteString(fmt.Sprintf("DELETE FROM reading_plans WHERE id = '%s';\n", planID))

	// Iterate calendar year (use 2024 as leap year reference or just generic 365? Plan says 365.)
	// We'll just loop 1 to 365.
	// We map Day Number to Month/Day logic to apply the "25 days" rule.
	// Assume standard year (Feb has 28).

	currentReadingDay := 0

	// Days in months
	monthDays := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	globalDay := 0
	for _, daysInMonth := range monthDays {
		for d := 1; d <= daysInMonth; d++ {
			globalDay++

			// Check if this is a reading day (<= 25)
			if d <= 25 {
				if currentReadingDay < targetDays {
					// Get readings
					r1 := readings1[currentReadingDay]
					r2 := readings2[currentReadingDay]
					r3 := readings3[currentReadingDay]
					r4 := readings4[currentReadingDay]

					passage := fmt.Sprintf("%s; %s; %s; %s", r1, r2, r3, r4)

					dayID := uuid.New()
					upSQL.WriteString(fmt.Sprintf("INSERT INTO reading_plan_days (id, reading_plan_id, day_number, passage, created_at) VALUES ('%s', '%s', %d, '%s', NOW());\n", dayID, planID, globalDay, escapeSQL(passage)))

					currentReadingDay++
				} else {
					// Ran out of readings? Should not happen if targetDays is correct.
					// But if it does, leave empty.
				}
			} else {
				// Catch-up day. No reading assigned.
				// We don't insert a row, or we insert a row with empty passage?
				// Usually empty or "Catch-up".
				// Frontend might expect a row for every day for a "calendar" plan?
				// "Calendar" plans usually imply specific dates. If we omit day 26, the UI might show it as "rest" or nothing.
				// Let's omit it to save DB space, assuming UI handles gaps.
				// OR, insert with NULL/Empty passage to be explicit.
				// Let's Insert "Catch Up / Reflection" text?
				// The system expects a passage ref. "Catch Up" is not a ref.
				// So we skip insert.
			}
		}
	}

	// Write files
	// Migration ID: Use a new timestamp.
	migrationID := "000028_add_discipleship_journal_plan"
	upFile := fmt.Sprintf("../../api/migrations/%s.up.sql", migrationID)
	downFile := fmt.Sprintf("../../api/migrations/%s.down.sql", migrationID)

	if err := os.WriteFile(upFile, []byte(upSQL.String()), 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(downFile, []byte(downSQL.String()), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %s\n", upFile)
}

func expand(books []Book) []string {
	var chapters []string
	for _, b := range books {
		for c := 1; c <= b.Chapters; c++ {
			chapters = append(chapters, fmt.Sprintf("%s %d", b.Name, c))
		}
	}
	return chapters
}

func distribute(chapters []string, days int) []string {
	var schedule []string

	if len(chapters) == 0 {
		return make([]string, days)
	}

	if len(chapters) >= days {
		// Compress: We have more chapters than days.
		// Use accumulator
		ratio := float64(len(chapters)) / float64(days)
		acc := 0.0
		chkIdx := 0

		for i := 0; i < days; i++ {
			acc += ratio
			count := int(math.Floor(acc))
			// floating point precision might make acc 1.999999 -> 1.
			// Ideally we consume 'count' chapters.
			// Adjust acc
			acc -= float64(count)

			// If we are at the last day, take all remaining
			if i == days-1 {
				count = len(chapters) - chkIdx
			}

			// Collect 'count' chapters
			var daysReadings []string
			for k := 0; k < count; k++ {
				if chkIdx < len(chapters) {
					daysReadings = append(daysReadings, chapters[chkIdx])
					chkIdx++
				}
			}

			if len(daysReadings) > 0 {
				// Compress ranges if possible? e.g. "Gen 1", "Gen 2" -> "Gen 1-2"
				// Simple join for now.
				schedule = append(schedule, compress(daysReadings))
			} else {
				schedule = append(schedule, "") // Should not happen often if ratio >= 1
			}
		}
	} else {
		// Expand: We have fewer chapters than days. Loop them.
		// Loop logic: just repeat the chapters until we fill 'days'.
		for i := 0; i < days; i++ {
			chapter := chapters[i%len(chapters)]
			schedule = append(schedule, chapter)
		}
	}

	return schedule
}

// compress takes "Genesis 1", "Genesis 2" and returns "Genesis 1-2"
func compress(readings []string) string {
	if len(readings) == 0 {
		return ""
	}
	if len(readings) == 1 {
		return readings[0]
	}

	// Check if all same book
	first := readings[0]
	parts := strings.Split(first, " ")
	if len(parts) < 2 {
		return strings.Join(readings, ", ")
	}
	book := strings.Join(parts[:len(parts)-1], " ")

	// Verify all same book
	for _, r := range readings {
		if !strings.HasPrefix(r, book) {
			return strings.Join(readings, ", ") // Different books, just join
		}
	}

	// Same book, range chums
	start := parts[len(parts)-1]
	last := readings[len(readings)-1]
	lastParts := strings.Split(last, " ")
	end := lastParts[len(lastParts)-1]

	return fmt.Sprintf("%s %s-%s", book, start, end)
}

func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
