package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	net "wfmon/pkg/network/mac"
)

var manufLineRegex = regexp.MustCompile(
	`^((?:[A-Fa-f0-9]{2}[\.:\-]){2,5}[A-Fa-f0-9]{2})(\/\d{2})?(?:\t(.*)\t(.*)\t#.*|\t(.*)\t#.*|\t(.*)\t(.*)|\t(.*))$`,
)

const (
	numGroups9    = 9
	manufSize     = 50000
	manufFile     = "manuf"
	manufTemplate = "manuf.tmpl"
)

// First arg is program name.
// Second --
// Third is base input directory.
// Forth is output file.
func main() {
	//nolint:gomnd // ignore
	if len(os.Args) < 4 {
		log.Fatal("input or output options not provided")
	}
	baseDir := os.Args[2]
	outputDest := os.Args[3]

	processTemplate(parseManuf(filepath.Join(baseDir, manufFile)), filepath.Join(baseDir, manufTemplate), outputDest)
}

func parseManuf(filePath string) map[string][]string {
	m := make(map[string][]string, manufSize)

	var err error
	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, short, long := parseLine(scanner.Text())
		if len(key) != 0 {
			m[key] = []string{short, long}
		}
	}

	if err = scanner.Err(); err != nil {
		log.Panic(err)
	}

	return m
}

func parseLine(line string) (string, string, string) {
	line = strings.TrimSpace(line)
	// skip comments
	if strings.HasPrefix(line, "#") || len(line) == 0 {
		return "", "", ""
	}

	// strings will be stored as raw - replace `
	line = strings.ReplaceAll(line, "`", "'")

	manufA := manufLineRegex.FindStringSubmatch(line)
	if len(manufA) != numGroups9 {
		log.Panicf("failed to parse %s", line)
	}

	mac, mask := manufA[1], manufA[2]
	short, long := manufA[3], manufA[4]
	if short == "" && manufA[5] != "" {
		short = manufA[5]
	}
	if short == "" && manufA[6] != "" {
		short = manufA[6]
	}
	if long == "" && manufA[7] != "" {
		long = manufA[7]
	}
	if short == "" && manufA[8] != "" {
		short = manufA[8]
	}

	//nolint:forbidigo // ignore
	fmt.Printf("MAC: %s Mask: %s Short: %s Long: %s\n", mac, mask, short, long)

	addr := new(net.HardwareAddr).WithAddr(mac).WithPrefix(mask)
	key := net.WildcardDotBigInt(*addr).String()

	return key, short, long
}

func processTemplate(m map[string][]string, tmplPath, outputPath string) {
	var err error

	tmpl := template.New(filepath.Base(tmplPath))
	if tmpl, err = tmpl.ParseFiles(tmplPath); err != nil {
		log.Panic(err)
	}

	var output bytes.Buffer
	if err = tmpl.Execute(&output, m); err != nil {
		log.Panic(err)
	}

	var formatted []byte
	if formatted, err = format.Source(output.Bytes()); err != nil {
		log.Panic(err)
	}
	// formatted = output.Bytes()

	var file *os.File
	if file, err = os.Create(outputPath); err != nil {
		log.Panic(err)
	}
	w := bufio.NewWriter(file)
	if _, err = w.Write(formatted); err != nil {
		log.Panic(err)
	}
	w.Flush()
}
