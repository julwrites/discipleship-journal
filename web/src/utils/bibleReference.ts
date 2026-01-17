

const BibleBooks: Record<string, string> = {
    // Old Testament
    "genesis": "Genesis", "gen": "Genesis", "gn": "Genesis",
    "exodus": "Exodus", "ex": "Exodus", "exo": "Exodus", "exod": "Exodus",
    "leviticus": "Leviticus", "lev": "Leviticus", "le": "Leviticus", "lv": "Leviticus",
    "numbers": "Numbers", "num": "Numbers", "nu": "Numbers", "nm": "Numbers",
    "deuteronomy": "Deuteronomy", "deut": "Deuteronomy", "de": "Deuteronomy", "dt": "Deuteronomy",
    "joshua": "Joshua", "josh": "Joshua", "jos": "Joshua",
    "judges": "Judges", "judg": "Judges", "jdg": "Judges",
    "ruth": "Ruth", "ru": "Ruth", "rth": "Ruth",
    "1 samuel": "1 Samuel", "1 sam": "1 Samuel", "1 sa": "1 Samuel", "1sam": "1 Samuel", "1sa": "1 Samuel", "i sam": "1 Samuel", "i samuel": "1 Samuel", "1st samuel": "1 Samuel",
    "2 samuel": "2 Samuel", "2 sam": "2 Samuel", "2 sa": "2 Samuel", "2sam": "2 Samuel", "2sa": "2 Samuel", "ii sam": "2 Samuel", "ii samuel": "2 Samuel", "2nd samuel": "2 Samuel",
    "1 kings": "1 Kings", "1 kgs": "1 Kings", "1 ki": "1 Kings", "1kings": "1 Kings", "i kgs": "1 Kings", "i kings": "1 Kings", "1st kings": "1 Kings",
    "2 kings": "2 Kings", "2 kgs": "2 Kings", "2 ki": "2 Kings", "2kings": "2 Kings", "ii kgs": "2 Kings", "ii kings": "2 Kings", "2nd kings": "2 Kings",
    "1 chronicles": "1 Chronicles", "1 chron": "1 Chronicles", "1 chr": "1 Chronicles", "1 ch": "1 Chronicles", "1chronicles": "1 Chronicles", "i chron": "1 Chronicles", "i chronicles": "1 Chronicles", "1st chronicles": "1 Chronicles",
    "2 chronicles": "2 Chronicles", "2 chron": "2 Chronicles", "2 chr": "2 Chronicles", "2 ch": "2 Chronicles", "2chronicles": "2 Chronicles", "ii chron": "2 Chronicles", "ii chronicles": "2 Chronicles", "2nd chronicles": "2 Chronicles",
    "ezra": "Ezra", "ezr": "Ezra",
    "nehemiah": "Nehemiah", "neh": "Nehemiah", "ne": "Nehemiah",
    "esther": "Esther", "est": "Esther", "esth": "Esther",
    "job": "Job", "jb": "Job",
    "psalms": "Psalms", "psalm": "Psalms", "ps": "Psalms", "psa": "Psalms",
    "proverbs": "Proverbs", "prov": "Proverbs", "pro": "Proverbs", "pr": "Proverbs",
    "ecclesiastes": "Ecclesiastes", "eccl": "Ecclesiastes", "ecc": "Ecclesiastes",
    "song of solomon": "Song of Solomon", "song": "Song of Solomon", "songs": "Song of Solomon", "sos": "Song of Solomon", "song of songs": "Song of Solomon",
    "isaiah": "Isaiah", "isa": "Isaiah",
    "jeremiah": "Jeremiah", "jer": "Jeremiah", "je": "Jeremiah",
    "lamentations": "Lamentations", "lam": "Lamentations",
    "ezekiel": "Ezekiel", "ezek": "Ezekiel", "eze": "Ezekiel", "ezk": "Ezekiel",
    "daniel": "Daniel", "dan": "Daniel", "dn": "Daniel",
    "hosea": "Hosea", "hos": "Hosea",
    "joel": "Joel", "jl": "Joel",
    "amos": "Amos", "am": "Amos",
    "obadiah": "Obadiah", "obad": "Obadiah", "ob": "Obadiah",
    "jonah": "Jonah", "jon": "Jonah", "jnh": "Jonah",
    "micah": "Micah", "mic": "Micah",
    "nahum": "Nahum", "nah": "Nahum",
    "habakkuk": "Habakkuk", "hab": "Habakkuk",
    "zephaniah": "Zephaniah", "zeph": "Zephaniah", "zep": "Zephaniah",
    "haggai": "Haggai", "hag": "Haggai", "hg": "Haggai",
    "zechariah": "Zechariah", "zech": "Zechariah", "zec": "Zechariah",
    "malachi": "Malachi", "mal": "Malachi",

    // New Testament
    "matthew": "Matthew", "matt": "Matthew", "mat": "Matthew", "mt": "Matthew",
    "mark": "Mark", "mrk": "Mark", "mk": "Mark", "mr": "Mark",
    "luke": "Luke", "luk": "Luke", "lk": "Luke",
    "john": "John", "jn": "John", "jhn": "John", "joh": "John",
    "acts": "Acts", "ac": "Acts", "act": "Acts",
    "romans": "Romans", "rom": "Romans", "ro": "Romans", "rm": "Romans",
    "1 corinthians": "1 Corinthians", "1 cor": "1 Corinthians", "1 co": "1 Corinthians", "1cor": "1 Corinthians", "i cor": "1 Corinthians", "i corinthians": "1 Corinthians", "1st corinthians": "1 Corinthians",
    "2 corinthians": "2 Corinthians", "2 cor": "2 Corinthians", "2 co": "2 Corinthians", "2cor": "2 Corinthians", "ii cor": "2 Corinthians", "ii corinthians": "2 Corinthians", "2nd corinthians": "2 Corinthians",
    "galatians": "Galatians", "gal": "Galatians", "ga": "Galatians",
    "ephesians": "Ephesians", "eph": "Ephesians", "ep": "Ephesians",
    "philippians": "Philippians", "phil": "Philippians", "php": "Philippians",
    "colossians": "Colossians", "col": "Colossians",
    "1 thessalonians": "1 Thessalonians", "1 thess": "1 Thessalonians", "1 th": "1 Thessalonians", "1thess": "1 Thessalonians", "i thess": "1 Thessalonians", "i thessalonians": "1 Thessalonians", "1st thessalonians": "1 Thessalonians",
    "2 thessalonians": "2 Thessalonians", "2 thess": "2 Thessalonians", "2 th": "2 Thessalonians", "2thess": "2 Thessalonians", "ii thess": "2 Thessalonians", "ii thessalonians": "2 Thessalonians", "2nd thessalonians": "2 Thessalonians",
    "1 timothy": "1 Timothy", "1 tim": "1 Timothy", "1 ti": "1 Timothy", "1tim": "1 Timothy", "i tim": "1 Timothy", "i timothy": "1 Timothy", "1st timothy": "1 Timothy",
    "2 timothy": "2 Timothy", "2 tim": "2 Timothy", "2 ti": "2 Timothy", "2tim": "2 Timothy", "ii tim": "2 Timothy", "ii timothy": "2 Timothy", "2nd timothy": "2 Timothy",
    "titus": "Titus", "tit": "Titus", "ti": "Titus",
    "philemon": "Philemon", "philem": "Philemon", "phlm": "Philemon", "phm": "Philemon",
    "hebrews": "Hebrews", "heb": "Hebrews",
    "james": "James", "jas": "James", "jm": "James",
    "1 peter": "1 Peter", "1 pet": "1 Peter", "1 pe": "1 Peter", "1 pt": "1 Peter", "1peter": "1 Peter", "i pet": "1 Peter", "i peter": "1 Peter", "1st peter": "1 Peter",
    "2 peter": "2 Peter", "2 pet": "2 Peter", "2 pe": "2 Peter", "2 pt": "2 Peter", "2peter": "2 Peter", "ii pet": "2 Peter", "ii peter": "2 Peter", "2nd peter": "2 Peter",
    "1 john": "1 John", "1 jn": "1 John", "1jn": "1 John", "1john": "1 John", "i jn": "1 John", "i john": "1 John", "1st john": "1 John",
    "2 john": "2 John", "2 jn": "2 John", "2jn": "2 John", "2john": "2 John", "ii jn": "2 John", "ii john": "2 John", "2nd john": "2 John",
    "3 john": "3 John", "3 jn": "3 John", "3jn": "3 John", "3john": "3 John", "iii jn": "3 John", "iii john": "3 John", "3rd john": "3 John",
    "jude": "Jude", "jud": "Jude", "jd": "Jude",
    "revelation": "Revelation", "rev": "Revelation",
};

const SingleChapterBooks: Record<string, boolean> = {
    "Obadiah": true,
    "Philemon": true,
    "2 John": true,
    "3 John": true,
    "Jude": true,
};

// Initialize sorted keys sorted by length (descending)
const sortedBookKeys = Object.keys(BibleBooks).sort((a, b) => b.length - a.length);

function isLetter(c: string) {
    return c.toLowerCase() !== c.toUpperCase();
}

function hasDigit(s: string) {
    return /\d/.test(s);
}

function consumeReferenceSyntax(s: string): [string, number] {
    let lastDigit = -1;

    for (let i = 0; i < s.length; i++) {
        const r = s[i];
        if (r >= '0' && r <= '9') {
            lastDigit = i;
        } else if (r === ':' || r === '-' || r === '.' || r === ',' || r === ' ' || r === '\t') {
            continue;
        } else {
            break;
        }
    }

    if (lastDigit === -1) {
        return ["", 0];
    }

    return [s.substring(0, lastDigit + 1), lastDigit + 1];
}

export function parseBibleReference(input: string): string | null {
    if (!input || !input.trim()) return null;

    // Skip leading whitespace
    const trimmedInput = input.trimStart();
    const lowerInput = trimmedInput.toLowerCase();

    let foundBook = "";
    let bookName = "";
    let matchLen = 0;

    // 1. Try exact match (Greedy)
    for (const key of sortedBookKeys) {
        if (lowerInput.startsWith(key)) {
            const mLen = key.length;
            // Ensure whole word match
            if (lowerInput.length > mLen) {
                const nextChar = lowerInput[mLen];
                if (isLetter(nextChar)) {
                    continue;
                }
            }

            foundBook = key;
            bookName = BibleBooks[key];
            matchLen = mLen;
            break;
        }
    }

    if (!foundBook) {
        return null;
    }

    // Parse numbers
    const remainderStart = matchLen;
    const currentInput = trimmedInput;
    let remainder = currentInput.substring(remainderStart);
    remainder = remainder.trimStart(); // consume spaces after book

    // Consume valid reference syntax
    const [syntax] = consumeReferenceSyntax(remainder);

    if (!syntax) {
        if (SingleChapterBooks[bookName]) {
            return bookName;
        }
        return `${bookName} 1`;
    }

    if (hasDigit(syntax)) {
        // Standardize syntax: remove spaces around separators? 
        // The user example "Philippians 2:3-4" implies preserving user's numbers but fixing book name.
        // We can just append the syntax as is (trimmed).
        return `${bookName} ${syntax.trim()}`;
    }

    if (SingleChapterBooks[bookName]) {
        return bookName;
    }
    return `${bookName} 1`;
}
