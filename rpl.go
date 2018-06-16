package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Args struct {
	IgnoreCase bool
	RegExp     bool
	WordWise   bool
}

var options = Args{
	IgnoreCase: false,
	RegExp:     false,
	WordWise:   false,
}

var from = ""
var to = ""
var filenames = make([]string, 0, 100)

var re *regexp.Regexp
var escapedTo string

var replacedFileCount = 0
var noChangeFileCount = 0
var ignoredFileCount = 0
var totalMatchCount = 0

func main() {
	parseArgs()

	buildRegexp()

	fmt.Printf("Search: %v\n", re)

	for _, filename := range filenames {
		isProcessed, matchCount := rpl(filename)
		if isProcessed {
			if matchCount > 0 {
				replacedFileCount += 1
			} else {
				noChangeFileCount += 1
			}
		} else {
			ignoredFileCount += 1
		}
		totalMatchCount += matchCount
	}

	fmt.Print("\n")
	fmt.Printf("%d files (replaced: %d / no change: %d / ignored: %d) Total %d matches\n", len(filenames), replacedFileCount, noChangeFileCount, ignoredFileCount, totalMatchCount)
}

func parseArgs() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <from> <to> files...\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVar(&options.IgnoreCase, "i", false, "ignore case")
	flag.BoolVar(&options.RegExp, "r", false, "regular expression search")
	flag.BoolVar(&options.WordWise, "w", false, "match whole word")

	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		//fmt.Println("Usage: rpl [-i] [-r] [-w] <from> <to> [files...]")
		os.Exit(1)
	}

	from = os.Args[flag.NFlag()+1]
	to = os.Args[flag.NFlag()+2]
	filenames = os.Args[flag.NFlag()+3:]
}

func buildRegexp() {
	var strRe string
	if options.RegExp {
		strRe = from
	} else {
		strRe = regexp.QuoteMeta(from)
	}
	if options.WordWise {
		strRe = `\b` + strRe + `\b`
	}
	if options.IgnoreCase {
		strRe = `(?i)` + strRe
	}
	re = regexp.MustCompile(strRe)

	escapedTo = regexp.MustCompile(`\$`).ReplaceAllString(to, "$$$$")
}

func message(process string, filename string, detail string, color string) {
	fmt.Printf("\x1b[%sm%s\x1b[0m    [%s] %s\n", color, process, filename, detail)
}

func rpl(filename string) (bool, int) {
	// 通常ファイルかどうか判定する
	fileInfo, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		message("Ignore", filename, "(not a regular file)", "1;33")
		return false, 0
	}

	message("Replace", filename, "n matches", "1;32")

	content := readFile(filename)

	replaced := re.ReplaceAllString(content, escapedTo)

	writeFile(filename, replaced)

	return true, 0
}

func readFile(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func writeFile(filename string, content string) {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}
