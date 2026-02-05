import { describe, it, expect } from 'vitest';
import { parseBibleReference } from './bibleReference';

describe('parseBibleReference', () => {
    it('returns null for empty or null input', () => {
        expect(parseBibleReference('')).toBeNull();
        expect(parseBibleReference('   ')).toBeNull();
        // @ts-ignore
        expect(parseBibleReference(null)).toBeNull();
    });

    it('parses standard references correctly', () => {
        expect(parseBibleReference('John 3:16')).toBe('John 3:16');
        expect(parseBibleReference('Genesis 1:1')).toBe('Genesis 1:1');
        expect(parseBibleReference('Romans 8:28')).toBe('Romans 8:28');
    });

    it('handles abbreviations', () => {
        expect(parseBibleReference('Gen 1:1')).toBe('Genesis 1:1');
        expect(parseBibleReference('Ex 20:1')).toBe('Exodus 20:1');
        expect(parseBibleReference('Matt 5:3')).toBe('Matthew 5:3');
        expect(parseBibleReference('Rom 8:1')).toBe('Romans 8:1');
        expect(parseBibleReference('Rev 21:1')).toBe('Revelation 21:1');
    });

    it('handles numbered books', () => {
        expect(parseBibleReference('1 John 1:9')).toBe('1 John 1:9');
        expect(parseBibleReference('1Jn 1:9')).toBe('1 John 1:9');
        expect(parseBibleReference('2 Cor 5:17')).toBe('2 Corinthians 5:17');
        expect(parseBibleReference('ii cor 5:17')).toBe('2 Corinthians 5:17');
        expect(parseBibleReference('1 Kings 18')).toBe('1 Kings 18');
    });

    it('handles single chapter books', () => {
        expect(parseBibleReference('Jude')).toBe('Jude');
        expect(parseBibleReference('Jude 24')).toBe('Jude 24');
        expect(parseBibleReference('Philemon')).toBe('Philemon');
        expect(parseBibleReference('Obadiah')).toBe('Obadiah');
        // If user types "Jude 1:24", it should probably be normalized or kept?
        // Current implementation appends syntax if digits found.
        expect(parseBibleReference('Jude 1:24')).toBe('Jude 1:24');
    });

    it('handles implicit chapter 1 for books without numbers', () => {
        expect(parseBibleReference('Genesis')).toBe('Genesis 1');
        expect(parseBibleReference('John')).toBe('John 1');
        // Unless it's a single chapter book
        expect(parseBibleReference('Jude')).toBe('Jude');
    });

    it('normalizes casing', () => {
        expect(parseBibleReference('john 3:16')).toBe('John 3:16');
        expect(parseBibleReference('GENESIS 1')).toBe('Genesis 1');
    });

    it('handles whitespace', () => {
        expect(parseBibleReference('  John   3:16  ')).toBe('John 3:16');
    });

    it('returns null for unknown books', () => {
        expect(parseBibleReference('Unknown 3:16')).toBeNull();
        expect(parseBibleReference('Harry Potter 1')).toBeNull();
    });

    it('handles various separators', () => {
        expect(parseBibleReference('John 3.16')).toBe('John 3.16');
        expect(parseBibleReference('John 3,16')).toBe('John 3,16'); // European style?
        expect(parseBibleReference('John 3 16')).toBe('John 3 16');
    });

    it('handles complex ranges', () => {
        expect(parseBibleReference('Philippians 2:3-4')).toBe('Philippians 2:3-4');
        expect(parseBibleReference('Gen 1:1-2:3')).toBe('Genesis 1:1-2:3');
    });
});
