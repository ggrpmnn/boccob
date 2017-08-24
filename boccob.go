package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Page contains the name and URL of a specific page
type Page struct {
	// the name of the page, determined by its <title /> element
	Name string
	// the full URL of the page
	PageURL string
	// how many times the link was encountered (relative importance)
	Weight int
}

var (
	baseURL string
)

func main() {
	// parse command line flags
	flag.StringVar(&baseURL, "s", "", "the base URL of the site to crawl (include http/https and subdomain)")
	flag.Parse()

	// check the command line params
	if baseURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	// create list and add site index to it
	pages := make(map[string]Page, 0)

	// create HTTP client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := &http.Client{Transport: tr}

	// recursively audit all the pages on the site
	Audit(hc, baseURL, pages)

	// print each entry
	c := 1
	fmt.Println("Entry, Name, URL, Weight")
	for _, p := range pages {
		fmt.Printf("%d, %s, %s, %d\n", c, p.Name, p.PageURL, p.Weight)
		c++
	}

	return
}

// Audit GETs a page, parses the HTML, and adds new URLs to the slice param
func Audit(hc *http.Client, urlStr string, pages map[string]Page) {
	// make request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("error creating new request to '%s'", urlStr)
		return
	}
	res, err := hc.Do(req)
	if err != nil {
		log.Printf("error making request to '%s': %s", urlStr, err.Error())
		return
	}
	if res.Body == nil {
		log.Printf("error response body '%s' is nil", urlStr)
		return
	}
	defer res.Body.Close()

	// read page response
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading response from '%s': %s", urlStr, err.Error())
		return
	}

	// set the name of the current page
	bodyStr := string(bodyBytes)
	rx := regexp.MustCompile(`<title>(.*)<\/title>`)
	title := rx.FindStringSubmatch(bodyStr)
	pages[urlStr] = Page{Name: title[1], PageURL: urlStr, Weight: 1}

	// find links on page using regexp
	rx = regexp.MustCompile(`href=["'](.*?)['"]`)
	matches := rx.FindAllStringSubmatch(bodyStr, -1)
	for _, match := range matches {
		linkURL := match[1]
		if ignoreURL(linkURL) {
			continue
		}
		// check whether the URL references the current site; criteria:
		// 1. full URL contains the site's base URL OR
		// 2. is a relative link (e.g. begins with '/')
		if !strings.HasPrefix(linkURL, baseURL) && !strings.HasPrefix(linkURL, "/") {
			continue
		}

		linkURL = validateAndFormURL(linkURL)
		// check if the URL exists in the map; add it ONLY if not
		if _, ok := pages[linkURL]; !ok {
			Audit(hc, linkURL, pages)
		} else {
			page := pages[linkURL]
			page.Weight = page.Weight + 1
			pages[linkURL] = page
		}
	}

	return
}

// ignoreURL checks against certain rules to see if the URL should be ignored
func ignoreURL(URLStr string) bool {
	// any URLs matching the rules will be ignored
	ignoreRules := []string{"_css", "css", "_js", "_html", "#", "events"}
	for _, r := range ignoreRules {
		if strings.Contains(URLStr, r) {
			return true
		}
	}
	return false
}

func validateAndFormURL(URLStr string) string {
	if strings.HasPrefix(URLStr, "/") {
		return baseURL + strings.TrimSuffix(URLStr, "/")
	}

	return URLStr
}
