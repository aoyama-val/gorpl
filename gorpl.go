package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

var (
	options struct {
		IgnoreCase bool `short:"i" long:"ignore-case" description:"Ignore case"`
		RegExp     bool `short:"r" long:"regexp" description:"Regular expression search. \\1 \\2 ... \\9 are replaced to corresponding submatch."`
		WordWise   bool `short:"w" long:"word" description:"Match whole word"`
	}

	from      = ""
	to        = ""
	filenames = []string{}

	re        *regexp.Regexp
	escapedTo = ""

	replacedFileCount = 0
	noChangeFileCount = 0
	ignoredFileCount  = 0
	totalMatchCount   = 0
)

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
	parser := flags.NewParser(&options, flags.PrintErrors|flags.PassDoubleDash)
	parser.Usage = "[options...] <from> <to> [filenames...]"

	args, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if len(args) < 2 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	from = args[0]
	to = args[1]
	filenames = args[2:]
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

	replaced := rplString(content)

	writeFile(filename, replaced)

	return true, 0
}

func rplString(str string) string {
	cb := func(s string) string {
		// Treat \1, \2, ..., \9
		if options.RegExp {
			submatches := re.FindStringSubmatch(s)
			replaceMap := []string{}
			for i := 1; i <= 9; i++ {
				replaceMap = append(replaceMap, "\\"+strconv.Itoa(i))
				if i < len(submatches) {
					replaceMap = append(replaceMap, submatches[i])
				} else {
					replaceMap = append(replaceMap, "")
				}
			}
			replacer := strings.NewReplacer(replaceMap...)
			return replacer.Replace(to)
		} else {
			return to
		}
	}

	replaced := re.ReplaceAllStringFunc(str, cb)
	return replaced
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
