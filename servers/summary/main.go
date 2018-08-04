package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// MetaTag represents an html meta tag with easier
// to access fields
type MetaTag struct {
	property string
	content  string
	name     string
}

// LinkTag represents and html link tag
type LinkTag struct {
	relIsIcon bool
	href      string
	sizes     string
	linkType  string
}

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	requestedURL := r.URL.Query().Get("url")
	// Added in this dummy URL to test quickly if https is working.
	// I like it, so I think I should keep it for now...
	if requestedURL == "test" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello world"))
		return
	} else if requestedURL == "" {
		log.Println(http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - url value not set"))
		return
	}

	pageBody, err := fetchHTML(requestedURL)
	if err != nil {
		log.Printf("Error in fetchHTML: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - provided url could not be fetched. see log for more details"))
	}

	iPageSummary, err := extractSummary(requestedURL, pageBody)
	if err != nil {
		log.Printf("Error in fetchHTML: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - error evaluating page tokens"))
	}

	//Close page body
	pageBody.Close()

	// Adding proper headers and returning the json.
	w.Header().Add("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(iPageSummary)

}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	// Fetch the body
	pageURL = makeValidURL("", pageURL)
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get experienced an error. response code returned was %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("response content-type is not html. it is %s", contentType)
	}
	return resp.Body, nil
}

// Takes in a htmlTag and returns a MetaTag that can
// be better accessed.
func makeMetaTag(tag html.Token) (MetaTag, error) {
	if tag.Data != "meta" {
		return MetaTag{}, fmt.Errorf("error tag is not meta. cannot convert")
	}
	metaTag := MetaTag{}
	// Sets the appropriate fields of a meta tag
	for _, attrib := range tag.Attr {
		if attrib.Key == "property" {
			metaTag.property = attrib.Val
		} else if attrib.Key == "content" {
			metaTag.content = attrib.Val
		} else if attrib.Key == "name" {
			metaTag.name = attrib.Val
		}
	}
	return metaTag, nil
}

// Returns a lnik tag that we can access more easily.
func makeLinkTag(tag html.Token) (LinkTag, error) {
	if tag.Data != "link" {
		return LinkTag{}, fmt.Errorf("error tag is not link. cannot convert")
	}
	linkTag := LinkTag{}
	// Sets the approiate attributes of the link tag.
	for _, attrib := range tag.Attr {
		if attrib.Key == "rel" && attrib.Val == "icon" {
			linkTag.relIsIcon = true
		} else if attrib.Key == "href" {
			linkTag.href = attrib.Val
		} else if attrib.Key == "sizes" {
			linkTag.sizes = attrib.Val
		} else if attrib.Key == "type" {
			linkTag.linkType = attrib.Val
		}
	}
	return linkTag, nil
}

// This adds "https://pageURL" to the beginning of URL
// so that it is an actual URL if it is incomplete.
func makeValidURL(pageURL string, input string) string {
	if strings.HasPrefix(input, "http") {
		return input
	} else {
		u, err := url.Parse(pageURL)
		if err != nil {
			log.Fatal("cannot parse page url")
		}
		return "http://" + u.Host + input
	}
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	// Create our containers
	page := &PageSummary{}
	openGraphImage := &PreviewImage{}
	tokenizer := html.NewTokenizer(htmlStream)
	backupTitle := ""
	backupDescription := ""

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			// need to return an empty PageSummary when error is found
			return &PageSummary{}, nil
		} else if tokenType == html.EndTagToken { // We are done when we see </head>
			token := tokenizer.Token()
			if token.Data == "head" {
				break
			}
		} else if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			// If we see a starting tag or a self-closing tag, we want to check it
			token := tokenizer.Token()
			// Handling case of title tag
			if token.Data == "title" {
				tokenType = tokenizer.Next()
				backupTitle = tokenizer.Token().Data
			}
			// Checking to see if we are dealing with a meta tag
			// and then looking through its attributes.
			if token.Data == "meta" {
				metaTag, err := makeMetaTag(token)
				if err == nil {
					switch metaTag.property {
					case "og:type":
						page.Type = metaTag.content
					case "og:url":
						page.URL = metaTag.content
					case "og:title":
						page.Title = metaTag.content
					case "og:site_name":
						page.SiteName = metaTag.content
					case "og:description":
						page.Description = metaTag.content
					case "og:image":
						fullURL := makeValidURL(pageURL, metaTag.content)
						// Check for relative pathing in image url
						openGraphImage = &PreviewImage{
							URL: fullURL,
						}
						// Because there coud be multiple images,
						// we need to append the image we have had if we see a new
						// og:image meta tag. Then we can edit its field Because
						// we have a pointer to it.
						page.Images = append(page.Images, openGraphImage)
					case "og:image:secure_url":
						openGraphImage.SecureURL = metaTag.content
					case "og:image:type":
						openGraphImage.Type = metaTag.content
					case "og:image:width":
						openGraphImage.Width, err = strconv.Atoi(metaTag.content)
					case "og:image:height":
						openGraphImage.Height, err = strconv.Atoi(metaTag.content)
					case "og:image:alt":
						openGraphImage.Alt = metaTag.content
					}
					// Switching on the name attribute
					switch metaTag.name {
					case "author":
						page.Author = metaTag.content
					case "keywords":
						page.Keywords = strings.Split(strings.Replace(metaTag.content, " ", "", -1), ",")
					case "description":
						backupDescription = metaTag.content
					}
				}
			}
			// Handle the link tag
			preview := &PreviewImage{}
			if token.Data == "link" {
				linkTag, err := makeLinkTag(token)
				// Make sure rel="icon"
				if err == nil && linkTag.relIsIcon {
					preview.Type = linkTag.linkType
					preview.URL = makeValidURL(pageURL, linkTag.href)
					// Make sure size is not blank/specififed as "any"
					if linkTag.sizes != "" && linkTag.sizes != "any" {
						widthHeight := strings.Split(strings.ToLower(linkTag.sizes), "x")
						width, widthErr := strconv.Atoi(widthHeight[1])
						height, heightErr := strconv.Atoi(widthHeight[0])
						// Make sure conversions didn't mess up
						if widthErr == nil && heightErr == nil {
							preview.Width = width
							preview.Height = height
						}
					}
				}
				page.Icon = preview
			}
		}
	}
	// Backup settings for when the og:properties are missing.
	if page.Title == "" {
		page.Title = backupTitle
	}
	if page.Description == "" {
		page.Description = backupDescription
	}

	return page, nil
}

func main() {
	summaryAddr := os.Getenv("SUMMARYADDR")
	if len(summaryAddr) == 0 {
		summaryAddr = "localhost:4001"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", SummaryHandler)
	log.Printf("listening on %s...", summaryAddr)
	log.Fatal(http.ListenAndServe(summaryAddr, mux))
}
