package pdf

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/ayoubomari/pacshare/app/models/pdfModels"
	"github.com/ayoubomari/pacshare/util/request"
)

func GetPdfInfo(pdfLink string) (pdfModels.ApkInfo, error) {
	var apkInfo pdfModels.ApkInfo

	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://www.pdfdrive.com%s", pdfLink),
		nil,
		nil,
	)
	if err != nil {
		return apkInfo, fmt.Errorf("GetPdfInfo: %w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return apkInfo, fmt.Errorf("GetPdfInfo: %w", err)
	}

	conver, _ := doc.Find(".ebook-left").First().Find("img").First().Attr("src")
	name := doc.Find(".ebook-right-inner").First().Find("h1").Text()
	author := doc.Find(".card-author").First().Text()
	fileSize := doc.Find(".ebook-file-info").Children().Eq(4).Text()
	pagesNum := doc.Find(".ebook-file-info").Children().Eq(0).Text()
	language := doc.Find(".ebook-file-info").Children().Eq(6).Text()

	// get session
	dataPreview, _ := doc.Find(".ebook-buttons").First().Find("button").First().Attr("data-preview")
	sessionID := ""
	dataPreviewSlices := regexp.MustCompile("=|&").Split(dataPreview, -1)
	if len(dataPreviewSlices) >= 4 {
		sessionID = fmt.Sprintf("%s_%s", dataPreviewSlices[1], dataPreviewSlices[3])
	}
	fmt.Println("sessionID:", sessionID)

	// build the return
	apkInfo = pdfModels.ApkInfo{
		Cover:     conver,
		Name:      name,
		Author:    author,
		FileSize:  fileSize,
		PagesNum:  pagesNum,
		Language:  language,
		SessionID: sessionID,
	}

	return apkInfo, nil
}
