# boccob
An external website auditing CLI tool written in Go.

### What It Does
Currently this tool will crawl a page and find any links (using the `href=` attribute), ignoring all external links. For any internal links it has found, it will make a request to that page and repeat the process, ignoring any links to pages it has already visited. If there are pages that have no links on the site, they will not show up in the results. Once site crawling is complete, then the name, URL, and weight of each page is output in a comma-separated format. 

### Building & Running
To run: `go run boccob.go -s https://site.to.crawl.com`

To build: `go build boccob.go -o boccob`, and then to run: `./boccob -s https://site.to.crawl.com`

### What's Next?
- [ ] Improve gathered data for each page (SEO, image count, etc.)
- [ ] Integrate with Cobra for CLI (as tooling increases)
- [ ] Add different reporting output options
