package apk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/apkModels"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

type scrapApkResponseBody struct {
	Datalist struct {
		List []apkModels.ApkInfo `json:"list,omitempty"`
	} `json:"datalist,omitempty"`
}

func scrapeApk(sender_psid string, searchKeyWords string) error {
	searchKeyWords = url.QueryEscape(searchKeyWords)

	// request search page by search key words
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://web-api-cache.aptoide.com/search?query=%s&country=MR&mature=false", searchKeyWords),
		nil,
		map[string]string{
			"Host":                      "web-api-cache.aptoide.com",
			"User-Agent":                "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0",
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
			"Accept-Language":           "en-US,en;q=0.5",
			"Connection":                "keep-alive",
			"Cookie":                    `_gcl_au=1.1.951052803.1708254414; _ga_WVQ7GSYQDV=GS1.1.1708806968.5.1.1708807316.0.0.0; _ga=GA1.2.1433557566.1708254414; Indicative_305bdd41-271f-4618-a1ea-0793da9e04ef="%7B%22defaultUniqueID%22%3A%22fb0f7178-1a58-471c-8a88-db4e4d1a185d%22%2C%22props%22%3A%7B%22subdomain%22%3A%22en%22%2C%22countryCode%22%3A%22MA%22%2C%22aptoide_package%22%3A%22aptoide.com%22%7D%2C%22lastSessionTime%22%3A1708807316710%7D"; cookie_settings=marketing%3Dtrue%26analytics%3Dtrue%26saved_settings%3Dtrue; _gid=GA1.2.1709577985.1708789560; searchHistory=%5B%22facebook%20lite%22%2C%22dama%22%2C%22karta%22%5D`,
			"Upgrade-Insecure-Requests": "1",
			"Sec-Fetch-Dest":            "document",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-Site":            "cross-site",
			"Pragma":                    "no-cache",
			"Cache-Control":             "no-cache",
		},
		false,
	)
	if err != nil {
		return fmt.Errorf("scrapeApk: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("scrapeApk: %w", err)
	}
	// fmt.Println(string(bodyBytes))
	var bodyJson scrapApkResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return fmt.Errorf("scrapeApk: %w", err)
	}

	//if there is no application found
	if len(bodyJson.Datalist.List) == 0 {
		response := facebook.ResponseMessage{
			Text: "No result found. 🔭\n" +
				"Try different keywords.",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	comul := 0
	n := config.MaxReturnedApk
	for i := 0; i < len(bodyJson.Datalist.List); i++ {
		// get app size => size = (apk size + obb size)
		// fmt.Printf("%+v\n", bodyJson.Datalist.List[i])
		size := bodyJson.Datalist.List[i].File.Filesize
		if bodyJson.Datalist.List[i].Obb != nil {
			// fmt.Println("bodyJson.Datalist.List[i].Obb.Main.Filesize:", bodyJson.Datalist.List[i].Obb.Main.Filesize)
			size += bodyJson.Datalist.List[i].Obb.Main.Filesize
		}

		if size <= config.ApkMaxAppSize {
			n--
			comul++
			if n+1 == 0 {
				break
			}

			response := facebook.ResponseTemplateAttachment{
				Type: "template",
				Payload: facebook.TemplateAttachmentPayload{
					TemplateType: "generic",
					Elements: []facebook.TemplateAttachmentElement{
						{
							Title:    fmt.Sprintf("%d# %s", comul, bodyJson.Datalist.List[i].Name),
							Subtitle: fmt.Sprintf("%.2f MB", formats.ByteToMegaByte(size)),
							ImageURL: bodyJson.Datalist.List[i].Icon,
							Buttons: []facebook.TemplateAttachmentButton{
								{
									Type:    "postback",
									Title:   "Download Now",
									Payload: fmt.Sprintf("APK_::_DOWNLOAD_::_%d", bodyJson.Datalist.List[i].ID),
								},
								{
									Type:    "postback",
									Title:   "See screenshots",
									Payload: fmt.Sprintf("APK_::_PHOTO_::_%d", bodyJson.Datalist.List[i].ID),
								},
								{
									Type:    "postback",
									Title:   "Description",
									Payload: fmt.Sprintf("APK_::_DESCRIPTION_::_%d", bodyJson.Datalist.List[i].ID),
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
			Text: "All these apps are large 😥.\n" +
				"Try different keywords.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
	}

	return nil
}
