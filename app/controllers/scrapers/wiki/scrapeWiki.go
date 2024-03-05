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

var ErrReferalPage = errors.New("wiki peferal page")
var SomethingWasWrong = facebook.ResponseMessage{
	Text: "No result found. 🔭\n" +
		"Try different keywords.",
}

// wiki search response body json struct
type WikiSearchResponseBody struct {
	Query struct {
		Search []struct {
			Title string `json:"title,omitempty"`
		} `json:"search"`
	} `josn:"query"`
}

// response body json struct
type scrapeWikiResponseBody struct {
	Error struct {
		Code string `json:"code,omitempty"`
	} `json:"error,omitempty"`
	Parse *struct {
		Text string `json:"text"`
	} `json:"parse,omitempty"`
}

// for scrape the wiki content and images
// @params: sender_psid="facebook sender ID"
// @params: arguments=[]string{"language","wiki title"}
func scrapeWiki(sender_psid string, arguments []string) error {
	lang := arguments[0]
	userTitle := arguments[1]

	pagesTitles, err := getwikiPagesTitles(lang, userTitle)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("scrapeWiki: %w", err)
	}
	fmt.Println("pagesTitles:", pagesTitles)

	_, text, err := getWikiArticleTextHtmlByPageTitle(lang, pagesTitles[0], true)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("scrapeWiki: %w", err)
	}

	//send the article as chunks
	//write the article content in a file
	userTitle = strings.ReplaceAll(pagesTitles[0], "_", "-")
	randomNumber := rand.Intn(1000)
	fileName := formats.ToFileNameString(userTitle) + "_" + fmt.Sprint(randomNumber) + ".1_1_pac.md"
	err = fs.WriteFile(fmt.Sprintf("./public/src/wikis/%s", fileName), text)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("scrapeWiki: %w", err)
	}

	response := facebook.ResponseMediaAttachment{
		Type: "file",
		Payload: facebook.WebhookBodyAttachmentPayload{
			URL: fmt.Sprintf("https://pacshare.omzor.com/src/wikis/%s", fileName),
		},
	}
	go facebookSender.CallSendAPIWithCallback(sender_psid, response, func(err error) error {
		return fs.DeleteFile(fmt.Sprintf("./public/src/wikis/%s", fileName))
	})

	// if the page and the content are exist
	// send the article photos
	// Find all images within the element with class "thumb" and get the src attribute
	// Replace digits followed by "px-" with "1080px-"
	// re2 := regexp.MustCompile("[0-9]{1,4}px-")
	// doc.Find(".thumb img").Each(func(i int, s *goquery.Selection) {
	// 	if i < 10 { // only the firt 10 photos
	// 		src, exists := s.Attr("src")
	// 		fmt.Println("src:", src)
	// 		finalResult := re2.ReplaceAllString(src, "1080px-")

	// 		if exists {
	// 			fmt.Printf("Image %d: %s\n", i+1, fmt.Sprintf("https:%s", finalResult))
	// 			response := facebook.ResponseMediaAttachment{
	// 				Type: "file",
	// 				Payload: facebook.WebhookBodyAttachmentPayload{
	// 					URL:         fmt.Sprintf("https:%s", finalResult),
	// 					Is_reusable: false,
	// 				},
	// 			}
	// 			go facebookSender.CallSendAPI(sender_psid, response)
	// 		}
	// 	}
	// })

	return nil
}

func getwikiPagesTitles(lang string, userTitle string) ([]string, error) {
	pagesTitles := make([]string, 0)

	// Specify the Wikipedia API endpoint
	apiEndpoint := fmt.Sprintf("https://%s.wikipedia.org/w/api.php", lang)
	// Specify the search parameters
	searchParams := url.Values{}
	searchParams.Set("action", "query")
	searchParams.Set("format", "json")
	searchParams.Set("list", "search")
	searchParams.Set("srsearch", userTitle)
	// Build the full URL with parameters
	fullURL := fmt.Sprintf("%s?%s", apiEndpoint, searchParams.Encode())
	fmt.Println("fullURL:", fullURL)

	res, err := request.JSONReqest(
		"GET",
		fullURL,
		nil,
		nil,
	)
	if err != nil {
		return pagesTitles, fmt.Errorf("scrapeWiki: %w", err)
	}
	defer res.Body.Close()

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

	for _, search := range wikiSearchBodyJson.Query.Search {
		pagesTitles = append(pagesTitles, search.Title)
	}

	return pagesTitles, nil
}

func getWikiArticleTextHtmlByPageTitle(lang string, pageTitle string, isRecursive bool) (*goquery.Document, *string, error) {
	var text string
	pageTitle = strings.ReplaceAll(pageTitle, " ", "_")
	fmt.Printf("https://%s.wikipedia.org/w/api.php?action=parse&prop=text&formatversion=2&format=json&page=%s\n", lang, url.QueryEscape(pageTitle))
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://%s.wikipedia.org/w/api.php?action=parse&prop=text&formatversion=2&format=json&page=%s", lang, url.QueryEscape(pageTitle)),
		nil,
		nil,
	)
	if err != nil {
		return nil, &text, fmt.Errorf("scrapeWiki: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &text, fmt.Errorf("scrapeWiki: %w", err)
	}
	var bodyJson scrapeWikiResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		fmt.Println(string(bodyBytes))
		return nil, &text, fmt.Errorf("scrapeWiki: json unmarshale error %w", err)
	}

	if bodyJson.Error.Code != "" {
		return nil, &text, fmt.Errorf("wiki search: %s", bodyJson.Error.Code)
	}

	// Create a new reader for the HTML string
	reader := strings.NewReader(bodyJson.Parse.Text)

	// Use goquery to parse the HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, &text, fmt.Errorf("error parsing HTML: %w", err)
	}

	if doc.Find(".mw-parser-output") == nil {
		return nil, &text, errors.New("no result found")
	}

	text = doc.Find(".mw-parser-output").Find("p").Text()
	if len(text) < 50 {
		if !isRecursive {
			return nil, &text, ErrReferalPage
		}

		a := doc.Find("ul").First().Find("li").First().Find("a").First()
		// Get the href attribute value
		href, exists := a.Attr("href")
		if !exists {
			return nil, &text, errors.New("no href attribute found")
		}

		wikiTitle := strings.ReplaceAll(href, "/wiki/", "")
		fmt.Println("new wiki title:", wikiTitle)
		return getWikiArticleTextHtmlByPageTitle(lang, wikiTitle, false)
	}

	text, err = formats.Utf8ToBase64(text)
	if err != nil {
		return nil, &text, fmt.Errorf("error converting utf8 to base64: %w", err)
	}

	return doc, &text, nil
}
