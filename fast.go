package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")

	for i, line := range lines {
		user := make(map[string]interface{})
		// fmt.Printf("%v %v\n", err, line)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}
			if ok := strings.Contains(browser, "Android"); ok {
				isAndroid = true
				seenBrowsers = addIfNotSeenBefore(seenBrowsers, browser, uniqueBrowsers)
			}

			if ok := strings.Contains(browser, "MSIE"); ok {
				isMSIE = true
				seenBrowsers = addIfNotSeenBefore(seenBrowsers, browser, uniqueBrowsers)
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

func addIfNotSeenBefore(seenBrowsers []string, browser string, uniqueBrowsers int) []string {
	for _, item := range seenBrowsers {
		if item == browser {
			return seenBrowsers
		}
	}
	// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
	seenBrowsers = append(seenBrowsers, browser)
	uniqueBrowsers++
	return seenBrowsers
}
