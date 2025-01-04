package pdf

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ayoubomari/pacshare/app/models/pdfModels"
	"github.com/ayoubomari/pacshare/util/request"
)

// GetPdfInfo retrieves PDF information by scraping the pdfLink HTML page
func GetPdfInfo(pdfLink string) (pdfModels.ApkInfo, error) {
	var apkInfo pdfModels.ApkInfo

	// Ensure pdfLink is not empty
	if pdfLink == "" {
		return apkInfo, fmt.Errorf("GetPdfInfo: empty pdfLink")
	}

	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://www.pdfdrive.com%s", pdfLink),
		nil,
		nil,
		false,
	)
	if err != nil {
		return apkInfo, fmt.Errorf("GetPdfInfo: %w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return apkInfo, fmt.Errorf("GetPdfInfo: %w", err)
	}

	cover, _ := doc.Find(".ebook-left").First().Find("img").First().Attr("src")
	name := strings.TrimSpace(doc.Find(".ebook-right-inner").First().Find("h1").Text())
	author := strings.TrimSpace(doc.Find(".card-author").First().Text())
	fileSize := strings.TrimSpace(doc.Find(".ebook-file-info").Children().Eq(4).Text())
	pagesNum := strings.TrimSpace(doc.Find(".ebook-file-info").Children().Eq(0).Text())
	language := strings.TrimSpace(doc.Find(".ebook-file-info").Children().Eq(6).Text())

	// Get session
	dataPreview, _ := doc.Find(".ebook-buttons").First().Find("button").First().Attr("data-preview")
	sessionID := ""
	if dataPreview != "" {
		dataPreviewSlices := regexp.MustCompile(`=|&`).Split(dataPreview, -1)
		if len(dataPreviewSlices) >= 4 {
			sessionID = fmt.Sprintf("%s_%s", dataPreviewSlices[1], dataPreviewSlices[3])
		}
	}

	// Build the return apkInfo
	apkInfo = pdfModels.ApkInfo{
		Cover:     cover,
		Name:      name,
		Author:    author,
		FileSize:  fileSize,
		PagesNum:  pagesNum,
		Language:  language,
		SessionID: sessionID,
	}

	return apkInfo, nil
}