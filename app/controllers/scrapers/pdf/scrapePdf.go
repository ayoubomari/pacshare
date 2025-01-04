package pdf

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

// scrapePdf searches for PDF files based on given keywords and sends results to the user via Facebook Messenger
func scrapePdf(sender_psid string, searchKeyWords string) error {
	// URL encode the search keywords
	searchKeyWords = url.QueryEscape(searchKeyWords)

	// Make a GET request to PDFDrive search page
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://www.pdfdrive.com/search?q=%s&pagecount=&pubyear=&searchin=&em=&more=true", searchKeyWords),
		nil,
		nil,
		false,
	)
	if err != nil {
		return fmt.Errorf("GetPdfInfo: %w", err)
	}
	defer res.Body.Close()

	// Parse the HTML response
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("GetPdfInfo: %w", err)
	}

	// Check if any results were found
	if doc.Find(".col-sm").Length() == 0 {
		response := facebook.ResponseMessage{
			Text: "No result found. 🔭\nTry different keywords.",
		}
		go facebookSender.CallSendAPI(sender_psid, response)
		return nil
	}

	comul := 0 // Counter for found PDFs
	n := config.MaxReturnedPdf // Maximum number of PDFs to return

	// Iterate through search results
	for i := 0; i < doc.Find(".col-sm").Length() && n > 0; i++ {
		fileInfo := doc.Find(".file-info").Eq(i)
		
		// Ensure we have enough children elements
		if fileInfo.Children().Length() < 5 {
			continue
		}

		// Parse the PDF file size
		sizeText := fileInfo.Children().Eq(4).Text()
		sizeFloat, err := formats.ParseFloat(sizeText)
		if err != nil {
			continue
		}

		size := formats.NumberToMegaBytes(sizeFloat)
		
		// Skip if file is too large
		if size > config.PdfMaxSize {
			continue
		}

		comul++
		n--

		// Extract PDF link
		fileLeft := doc.Find(".file-left").Eq(i)
		pdfLink, exists := fileLeft.Find("a").First().Attr("href")
		if !exists {
			continue
		}

		// Extract title and subtitle
		title := strings.Trim(doc.Find(".file-right").Eq(i).Find("h2").First().Text(), " ")
		subtitle := fileInfo.Children().Eq(0).Text() // return the size of the pdf in MB

		// Prepare response template
		response := facebook.ResponseTemplateAttachment{
			Type: "template",
			Payload: facebook.TemplateAttachmentPayload{
				TemplateType: "generic",
				Elements: []facebook.TemplateAttachmentElement{
					{
						Title:    fmt.Sprintf("%d# %s", comul, title),
						Subtitle: subtitle,
						ImageURL: "",
						Buttons: []facebook.TemplateAttachmentButton{
							{
								Type:    "postback",
								Title:   "Download Now",
								Payload: fmt.Sprintf("PDF_::_DOWNLOAD_::_%s", pdfLink),
							},
							{
								Type:    "postback",
								Title:   "See the cover",
								Payload: fmt.Sprintf("PDF_::_PHOTO_::_%s", pdfLink),
							},
							{
								Type:    "postback",
								Title:   "Description",
								Payload: fmt.Sprintf("PDF_::_DESCRIPTION_::_%s", pdfLink),
							},
						},
					},
				},
			},
		}

		// Send the response asynchronously
		go facebookSender.CallSendAPI(sender_psid, response)
	}

	// If no suitable PDFs were found, send a message to the user
	if comul == 0 {
		response := facebook.ResponseMessage{
			Text: "All these pdfs are large 😥.\nTry different keywords.",
		}
		go facebookSender.CallSendAPI(sender_psid, response)
	}

	return nil
}