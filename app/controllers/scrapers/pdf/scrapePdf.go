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

func scrapePdf(sender_psid string, searchKeyWords string) error {
	searchKeyWords = url.QueryEscape(searchKeyWords)

	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://www.pdfdrive.com/search?q=%s&pagecount=&pubyear=&searchin=&em=&more=true", searchKeyWords),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("GetPdfInfo: %w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("GetPdfInfo: %w", err)
	}

	if doc.Find(".col-sm").Length() == 0 {
		response := facebook.ResponseMessage{
			Text: "No result found. 🔭\n" +
				"Try different keywords.",
		}
		facebookSender.CallSendAPI(sender_psid, response)

		return nil
	}

	comul := 0
	n := config.MaxReturnedPdf
	for i := 0; i < doc.Find(".col-sm").Length(); i++ {
		// get app size
		sizeFloat, err := formats.ParseFloat(doc.Find(".file-info").Eq(i).Children().Eq(4).Text())
		if err != nil {
			continue
		}
		size := formats.NumberToMegaBytes(sizeFloat)
		fmt.Printf("%d\n", size)
		if size <= config.ApkMaxAppSize {
			n--
			comul++
			if n+1 == 0 {
				break
			}

			imgURL, _ := doc.Find(".file-left").Eq(i).Find("img").First().Attr("src")
			pdfLink, _ := doc.Find(".file-left").Eq(i).Find("a").First().Attr("href")
			fmt.Println("imgURL:", imgURL)
			fmt.Println("pdfLink:", pdfLink)

			response := facebook.ResponseTemplateAttachment{
				Type: "template",
				Payload: facebook.TemplateAttachmentPayload{
					TemplateType: "generic",
					Elements: []facebook.TemplateAttachmentElement{
						{
							Title:    fmt.Sprintf("%d# %s", comul, strings.Trim(doc.Find(".file-right").Eq(i).Find("h2").First().Text(), " ")),
							Subtitle: doc.Find(".file-info").Eq(i).Children().Eq(0).Text(),
							ImageURL: "", // "https://pacshare.omzor.com/img/forContent/pdfLogo.png"
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

			go facebookSender.CallSendAPI(sender_psid, response)
		}
	}
	if comul == 0 {
		response := facebook.ResponseMessage{
			Text: "All these pdfs are large 😥.\n" +
				"Try different keywords.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
	}

	return nil
}
