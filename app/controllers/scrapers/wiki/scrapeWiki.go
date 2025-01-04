package wiki

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/fs"
	"github.com/ayoubomari/pacshare/util/request"
)

// Error for referral pages
var ErrReferalPage = errors.New("wiki peferal page")

// Default error message for users
var SomethingWasWrong = facebook.ResponseMessage{
	Text: "No result found. 🔭\n" +
		"Try different keywords.",
}

// WikiSearchResponseBody represents the structure of the Wikipedia search API response
type WikiSearchResponseBody struct {
	Query struct {
		Search []struct {
			Title string `json:"title,omitempty"`
		} `json:"search"`
	} `json:"query"`
}

// scrapeWikiResponseBody represents the structure of the Wikipedia content API response
type scrapeWikiResponseBody struct {
	Error struct {
		Code string `json:"code,omitempty"`
	} `json:"error,omitempty"`
	Parse struct {
		Text string `json:"text"`
	} `json:"parse"`
}

// scrapeWiki fetches and processes Wikipedia content
// @params: sender_psid - Facebook sender ID
// @params: arguments - slice containing language and wiki title
func scrapeWiki(sender_psid string, arguments []string) error {
	// Check if we have sufficient arguments
	if len(arguments) < 2 {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return errors.New("insufficient arguments")
	}

	lang := arguments[0]
	userTitle := arguments[1]

	// Get Wikipedia page titles based on user query
	pagesTitles, err := getwikiPagesTitles(lang, userTitle)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("scrapeWiki: %w", err)
	}

	// Check if we found any pages
	if len(pagesTitles) == 0 {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return errors.New("no pages found")
	}

	// Get the content of the first page found
	_, text, err := getWikiArticleTextHtmlByPageTitle(lang, pagesTitles[0], true)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("scrapeWiki: %w", err)
	}

	// Prepare file name for storing the article content
	userTitle = strings.ReplaceAll(pagesTitles[0], "_", "-")
	randomNumber := rand.Intn(1000)
	fileName := formats.ToFileNameString(userTitle) + "_" + fmt.Sprint(randomNumber) + ".1_1_pac.txt"

	// Write the article content to a file
	err = fs.WriteFile(fmt.Sprintf("./public/src/wikis/%s", fileName), text)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("scrapeWiki: %w", err)
	}

	// Prepare response with file URL
	response := facebook.ResponseMediaAttachment{
		Type: "file",
		Payload: facebook.WebhookBodyAttachmentPayload{
			URL: fmt.Sprintf("https://pacshare.omzor.com/src/wikis/%s", fileName),
		},
	}

	// Send response and delete file after sending
	go facebookSender.CallSendAPIWithCallback(sender_psid, response, func(err error) error {
		return fs.DeleteFile(fmt.Sprintf("./public/src/wikis/%s", fileName))
	})

	return nil
}

// getwikiPagesTitles fetches Wikipedia page titles based on user query
func getwikiPagesTitles(lang string, userTitle string) ([]string, error) {
	pagesTitles := make([]string, 0)

	// Prepare Wikipedia API endpoint
	apiEndpoint := fmt.Sprintf("https://%s.wikipedia.org/w/api.php", lang)
	searchParams := url.Values{}
	searchParams.Set("action", "query")
	searchParams.Set("format", "json")
	searchParams.Set("list", "search")
	searchParams.Set("srsearch", userTitle)
	fullURL := fmt.Sprintf("%s?%s", apiEndpoint, searchParams.Encode())

	// Make API request
	res, err := request.JSONReqest("GET", fullURL, nil, nil, false)
	if err != nil {
		return pagesTitles, fmt.Errorf("scrapeWiki: %w", err)
	}
	defer res.Body.Close()

	// Read and parse response
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return pagesTitles, fmt.Errorf("scrapeWiki: %w", err)
	}
	var wikiSearchBodyJson WikiSearchResponseBody
	err = json.Unmarshal(bodyBytes, &wikiSearchBodyJson)
	if err != nil {
		return pagesTitles, fmt.Errorf("scrapeWiki: json unmarshale wikiSearch error %w", err)
	}
	if len(wikiSearchBodyJson.Query.Search) == 0 {
		return pagesTitles, errors.New("no result found")
	}

	// Extract page titles from response
	for _, search := range wikiSearchBodyJson.Query.Search {
		pagesTitles = append(pagesTitles, search.Title)
	}

	return pagesTitles, nil
}

// getWikiArticleTextHtmlByPageTitle fetches and processes Wikipedia article content
func getWikiArticleTextHtmlByPageTitle(lang string, pageTitle string, isRecursive bool) (*goquery.Document, *string, error) {
	var text string
	pageTitle = strings.ReplaceAll(pageTitle, " ", "_")

	// Make API request for article content
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://%s.wikipedia.org/w/api.php?action=parse&prop=text&formatversion=2&format=json&page=%s", lang, url.QueryEscape(pageTitle)),
		nil,
		nil,
		false,
	)
	if err != nil {
		return nil, &text, fmt.Errorf("scrapeWiki: %w", err)
	}
	defer res.Body.Close()

	// Read and parse response
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &text, fmt.Errorf("scrapeWiki: %w", err)
	}
	var bodyJson scrapeWikiResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return nil, &text, fmt.Errorf("scrapeWiki: json unmarshale error %w", err)
	}

	if bodyJson.Error.Code != "" {
		return nil, &text, fmt.Errorf("wiki search: %s", bodyJson.Error.Code)
	}

	// Parse HTML content
	reader := strings.NewReader(bodyJson.Parse.Text)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, &text, fmt.Errorf("error parsing HTML: %w", err)
	}

	// Extract main content
	mwParserOutput := doc.Find(".mw-parser-output")
	if mwParserOutput.Length() == 0 {
		return nil, &text, errors.New("no result found")
	}

	text = mwParserOutput.Find("p").Text()

	// Handle referral pages
	if len(text) < 50 {
		if !isRecursive {
			return nil, &text, ErrReferalPage
		}

		// Try to find a link to the actual content page
		firstUl := doc.Find("ul").First()
		if firstUl.Length() == 0 {
			return nil, &text, errors.New("no ul element found")
		}

		firstLi := firstUl.Find("li").First()
		if firstLi.Length() == 0 {
			return nil, &text, errors.New("no li element found")
		}

		a := firstLi.Find("a").First()
		if a.Length() == 0 {
			return nil, &text, errors.New("no a element found")
		}

		href, exists := a.Attr("href")
		if !exists {
			return nil, &text, errors.New("no href attribute found")
		}

		wikiTitle := strings.ReplaceAll(href, "/wiki/", "")
		return getWikiArticleTextHtmlByPageTitle(lang, wikiTitle, false)
	}

	return doc, &text, nil
}