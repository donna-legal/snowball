package snowball

import "testing"

// Stuct for holding tests where a word is transformed
// into another by the function to be tested.
//
type simpleStringTestCase struct {
	wordIn  string
	wordOut string
}

// A type representing all functions that take one
// string and return another.
type simpleStringFunction func(string) string

// Runs a series of test cases for functions that just
// transform one string into another.
//
func runSimpleStringTests(t *testing.T, f simpleStringFunction, tcs []simpleStringTestCase) {
	for _, testCase := range tcs {
		output := f(testCase.wordIn)
		if output != testCase.wordOut {
			t.Errorf("Expected \"%v\", but got \"%v\"", testCase.wordOut, output)
		}
	}
}

func Test_constants(t *testing.T) {
	expectedVowels := "aeiouy"
	if vowels != "aeiouy" {
		t.Errorf("Expected %v, got %v", expectedVowels, vowels)
	}
}

// Test stopWords for things we know should be true
// or false.
//
func Test_stopWords(t *testing.T) {

	// Test true
	knownTrueStopwords := [...]string{
		"a",
		"for",
		"be",
		"was",
	}
	for _, word := range knownTrueStopwords {
		if stopWords[word] == false {
			t.Errorf("Expected %v, to be in stopWords", word)
		}
	}

	// Test false
	knownFalseStopwords := [...]string{
		"truck",
		"deoxyribonucleic",
		"farse",
		"bullschnizzle",
	}
	for _, word := range knownFalseStopwords {
		if stopWords[word] == true {
			t.Errorf("Expected %v, to be in stopWords", word)
		}
	}
}

// Test specialWords for things we know should be present
// and not present.
//
func Test_specialWords(t *testing.T) {

	// Test true
	knownTrueSpecialwords := [...]string{
		"exceeding",
		"early",
		"outing",
	}
	for _, word := range knownTrueSpecialwords {
		if _, ok := specialWords[word]; !ok {
			t.Errorf("Expected %v, to be in specialWords", word)
		}
	}

	// Test false
	knownFalseSpecialwords := [...]string{
		"truck",
		"deoxyribonucleic",
		"farse",
		"bullschnizzle",
	}
	for _, word := range knownFalseSpecialwords {
		if _, ok := specialWords[word]; ok {
			t.Errorf("Expected %v, to NOT be in specialWords", word)
		}
	}
}

func Test_normalizeApostrophes(t *testing.T) {
	variants := [...]string{
		"\u2019xxx\u2019",
		"\u2018xxx\u2018",
		"\u201Bxxx\u201B",
		"’xxx’",
		"‘xxx‘",
		"‛xxx‛",
	}
	for _, v := range variants {
		normalizedVersion := normalizeApostrophes(v)
		if normalizedVersion != "'xxx'" {
			t.Errorf("Expected \"'xxx'\", not \"%v\"", normalizedVersion)
		}
	}
}

func Test_isLowerVowel(t *testing.T) {
	for _, r := range vowels {
		if isLowerVowel(r) == false {
			t.Errorf("Expected \"%v\" to be a vowel", r)
		}
	}

	consonant := "bcdfghjklmnpqrstvwxz"
	for _, r := range consonant {
		if isLowerVowel(r) == true {
			t.Errorf("Expected \"%v\" to NOT be a vowel", r)
		}
	}
}

func Test_capitalizeYs(t *testing.T) {
	var wordTests = []simpleStringTestCase{
		{"ysdcsdeysdfsysdfsdiyoyyyxyxayxey", "YsdcsdeYsdfsysdfsdiYoYyYxyxaYxeY"},
	}
	runSimpleStringTests(t, capitalizeYs, wordTests)
}

func Test_preprocessWord(t *testing.T) {
	var wordTests = []simpleStringTestCase{
		{"arguing", "arguing"},
		{"Arguing", "arguing"},
		{"'catty", "catty"},
		{"Kyle’s", "kyle's"},
		{"toy", "toY"},
	}
	runSimpleStringTests(t, preprocessWord, wordTests)
}

func Test_vnvSuffix(t *testing.T) {
	var wordTests = []simpleStringTestCase{
		{"crepuscular", "uscular"},
		{"uscular", "cular"},
	}
	runSimpleStringTests(t, vnvSuffix, wordTests)
}

func Test_r1r2(t *testing.T) {
	var wordTests = []struct {
		word string
		r1   string
		r2   string
	}{
		{"crepuscular", "uscular", "cular"},
		{"beautiful", "iful", "ul"},
		{"beauty", "y", ""},
		{"eucharist", "harist", "ist"},
		{"animadversion", "imadversion", "adversion"},
		{"communism", "ism", "m"},
		{"arsenal", "al", ""},
		{"generalities", "alities", "ities"},
	}
	for _, testCase := range wordTests {
		r1, r2 := r1r2(testCase.word)
		if r1 != testCase.r1 || r2 != testCase.r2 {
			t.Errorf("Expected \"{%v, %v}\", but got \"{%v, %v}\"", testCase.r1, testCase.r2, r1, r2)
		}
	}
}

func Test_step0(t *testing.T) {
	var wordTests = []struct {
		wordIn  string
		r1in    string
		r2in    string
		wordOut string
		r1out   string
		r2out   string
	}{
		{"general's", "al's", "", "general", "al", ""},
		{"general's'", "al's'", "", "general", "al", ""},
		{"spices'", "es'", "", "spices", "es", ""},
	}
	for _, testCase := range wordTests {
		wordOut, r1out, r2out := step0(testCase.wordIn, testCase.r1in, testCase.r2in)
		if wordOut != testCase.wordOut || r1out != testCase.r1out || r2out != testCase.r2out {
			t.Errorf("Expected \"{%v, %v, %v}\", but got \"{%v, %v, %v}\"", testCase.wordOut, testCase.r1out, testCase.r2out, wordOut, r1out, r2out)
		}
	}
}