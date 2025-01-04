package yt

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	goaway "github.com/TwiN/go-away"
	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

type scrapeYtResponseBody struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer struct {
							Contents []struct {
								VideoRenderer struct {
									VideoID    string `json:"videoId"`
									LengthText struct {
										SimpleText string `json:"simpleText"`
									} `json:"lengthText"`
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
								} `json:"videoRenderer"`
							} `json:"contents"`
						} `json:"itemSectionRenderer"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

func (s *scrapeYtResponseBody) UnmarshalJSON(data []byte) error {
	type Alias scrapeYtResponseBody
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Ensure that all slices are at least empty slices, not nil
	if s.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents == nil {
		s.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents = []struct {
			ItemSectionRenderer struct {
				Contents []struct {
					VideoRenderer struct {
						VideoID    string `json:"videoId"`
						LengthText struct {
							SimpleText string `json:"simpleText"`
						} `json:"lengthText"`
						Title struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"title"`
					} `json:"videoRenderer"`
				} `json:"contents"`
			} `json:"itemSectionRenderer"`
		}{}
	}
	return nil
}

func scrapeYt(sender_psid string, searchKeyWords string) error {
	if goaway.IsProfane(searchKeyWords) {
		response := facebook.ResponseMessage{
			Text: "No result found. 🔭\nTry different keywords.",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	jsonBytes := []byte(`{"context":{"client":{"hl":"en","gl":"MA","remoteHost":"196.64.162.134","deviceMake":"","deviceModel":"","visitorData":"CgtGWGVqTEdWWEJRNCiM6-q4BjIKCgJNQRIEGgAgLQ%3D%3D","userAgent":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36,gzip(gfe)","clientName":"WEB","clientVersion":"2.20241024.01.00","osName":"X11","osVersion":"","originalUrl":"https://www.youtube.com/results?search_query=` + searchKeyWords + `","platform":"DESKTOP","clientFormFactor":"UNKNOWN_FORM_FACTOR","configInfo":{"appInstallData":"CIzr6rgGEPOisQUQydewBRDE2LEFELfq_hIQgdaxBRDalM4cEKLUsQUQppOxBRD2hrEFEN6tsQUQvYqwBRCHw7EFEPDHsQUQz9H_EhDW3bAFEOHssAUQsO6wBRD4ubEFEI7QsQUQndCwBRD2q7AFELfvrwUQqLH_EhCBw7EFEInorgUQ5s-xBRCJp7EFEMrYsQUQ0ZTOHBDvzbAFEL2ZsAUQxfWwBRCIh7AFEIzQsQUQxqSxBRComrAFENuvrwUQmY2xBRCdprAFEKaasAUQhaexBRCKobEFEMjYsQUQyeawBRD7lc4cEMvRsQUQiOOvBRDxnLAFEKaSsQUQjcywBRCQzLEFENXWsQUQytSxBRDgzbEFEMzfrgUQksuxBRDqkM4cEJeUzhwQi9SxBRDgjf8SEJT-sAUQmoG4IhDzn84cEOLUrgUQ276xBRDUwbEFEPirsQUQop2xBRCN1LEFEOuZsQUQhLawBRCNlLEFEOrDrwUQ5bmxBRCq2LAFENPhrwUQrZ7OHBDGv7EFEI_DsQUQ3ej-EhCWlbAFEKPN_xIQ18GxBRDQjbAFEOvo_hIQwc2xBRDJ968FENfprwUQusSxBRDtubEFEKDZ_xIQzdGxBRCazrEFEIzU_xIQhcOxBRCkmc4cEM3XsAUQtv-3IhC9tq4FEJmYsQUQis7_EhCYnM4cEKTUsQUQl9exBRDb17EFELDU_xIQ_9ixBRD0tLEFKixDQU1TSEJVZm9MMndETkhrQnUySW9BQ08zTWdDal9RT3FBemhjdFo4SFFjPQ%3D%3D"},"userInterfaceTheme":"USER_INTERFACE_THEME_DARK","timeZone":"Africa/Casablanca","browserName":"Chrome","browserVersion":"130.0.0.0","acceptHeader":"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7","deviceExperimentId":"ChxOelF5T1RRMU1ERTBPRFk0TURFME16UXhNdz09EIzr6rgGGIzr6rgG","screenWidthPoints":1141,"screenHeightPoints":781,"screenPixelDensity":1,"screenDensityFloat":1,"utcOffsetMinutes":60,"memoryTotalKbytes":"8000000","mainAppWebInfo":{"graftUrl":"/results?search_query=` + searchKeyWords + `","pwaInstallabilityStatus":"PWA_INSTALLABILITY_STATUS_CAN_BE_INSTALLED","webDisplayMode":"WEB_DISPLAY_MODE_BROWSER","isWebNativeShareAvailable":false}},"user":{"lockedSafetyMode":false},"request":{"useSsl":true,"consistencyTokenJars":[{"encryptedTokenJarContents":"AKreu9ttE1I3wjoxGirEb7mjUxe65Tlmg8OdAwy70Cszz_61zja1qQdyvBZ4sxAYfGrG7v89ILzdZ-VHHJL6-agrgRmhjPeSJ5VCJKtLo_EFvFNANljoSbuLURl2EXv32O1LRvChdCNJbg3-8AjmqUrJ"}],"internalExperimentFlags":[]},"clickTracking":{"clickTrackingParams":"CBMQ7VAiEwjCwdXp9KeJAxWqPwYAHWfcHpY="},"adSignalsInfo":{"params":[{"key":"dt","value":"1729803661360"},{"key":"flash","value":"0"},{"key":"frm","value":"0"},{"key":"u_tz","value":"60"},{"key":"u_his","value":"6"},{"key":"u_h","value":"900"},{"key":"u_w","value":"1600"},{"key":"u_ah","value":"868"},{"key":"u_aw","value":"1600"},{"key":"u_cd","value":"24"},{"key":"bc","value":"31"},{"key":"bih","value":"781"},{"key":"biw","value":"1126"},{"key":"brdim","value":"0,32,0,32,1600,32,1600,868,1141,781"},{"key":"vis","value":"1"},{"key":"wgl","value":"true"},{"key":"ca_type","value":"image"}],"bid":"ANyPxKom_1ac7pBPZGN0VmY2l3aXP6NXpHm1lIRyLft7OCw3lr_yd9byb_yjnKTkClp58ThBkVCQrKZGjppndkQn7e82WMx_rg"}},"query":"` + searchKeyWords + `","webSearchboxStatsUrl":"/search?oq=` + searchKeyWords + `&gs_l=youtube.3..0i512i433i67i131i650k1j0..."}`)
	headersKeysValuesPairs := map[string]string{
		"accept":                        "*/*",
		"accept-language":               "en,fr-FR;q=0.9,fr;q=0.8,en-US;q=0.7",
		"authorization":                 "SAPISIDHASH 1729803691_a2943e8b303ae09307585e06617e27be81acd4e8_u",
		"cache-control":                 "no-cache",
		"content-type":                  "application/json",
		"cookie":                        "YSC=BAUuYRMkRWg; VISITOR_PRIVACY_METADATA=CgJNQRIEGgAgLQ%3D%3D; VISITOR_INFO1_LIVE=FXejLGVXBQ4; PREF=f7=4100&tz=Africa.Casablanca&f5=30000&f6=40000000; SID=g.a000ogijVWnOsznebc_Q3XtfOWfD9lncd7LzLQMVtd5VlT7iNtPwNmgaRkM_7Yb9n3BsBQDphAACgYKARkSARMSFQHGX2MiKHsgkP-pUnVkpd1vIAClBRoVAUF8yKp6o1zR8tSfUBU-vl5S9fOA0076; __Secure-1PSID=g.a000ogijVWnOsznebc_Q3XtfOWfD9lncd7LzLQMVtd5VlT7iNtPwfjqlb4UMcj8BaQXHOUreygACgYKAc4SARMSFQHGX2MiL-xX6lCqmXoLg5HgK4OHNhoVAUF8yKrSbx4rtWSpf7nzaPokpgxO0076; __Secure-3PSID=g.a000ogijVWnOsznebc_Q3XtfOWfD9lncd7LzLQMVtd5VlT7iNtPwKPZu0dUa7ADO_6QnE5tpRQACgYKATMSARMSFQHGX2MipatA1aeIS-_D6RtCQhdEbhoVAUF8yKrGv3lwX7beY8QFwAo-n7kR0076; HSID=ALxVH5MXEYIzPzy_u; SSID=A7d2MuXs7cIJBqyO0; APISID=0UvM2c36pCUL299b/ADdRrCRFEYGgi2Nde; SAPISID=FIi-86L8A2TI_WD_/AlU3rjC0laAEZAIDE; __Secure-1PAPISID=FIi-86L8A2TI_WD_/AlU3rjC0laAEZAIDE; __Secure-3PAPISID=FIi-86L8A2TI_WD_/AlU3rjC0laAEZAIDE; LOGIN_INFO=AFmmF2swRQIgAYpzfY6Q6vCi-dCFs8MJuW2VkQvfMdYLs4UR7IB1hwECIQD2Zf6-G2QV7h3BgMqoTJy9IRUkxjJ7pilQlh7APa_bHg:QUQ3MjNmeC1EQnpfeDFUS2NVTVk5OW8yUUlUSG5fSzRjQ1E5NVN2ajU5VkJHU2NjRERSdVNMWVNGMF9JS0NRVGVHWDRfYzhDTUJpOHlEVGdzblVWS0QzT2RBQktwb0NLN3VBa3I2UEJ3UHFjQmppT3pSRkduZkhsOFV6Sl90aFN4MXFjUmN0Zmk0RGdfSWcxbVAwQXFBcFlJWGFWS2RycTZB; __Secure-1PSIDTS=sidts-CjEBQT4rX28cVx_kTaTU86w8uiJVdgeCPt7OeBj6kiw6VkQbepj3NesQWdyekxjqHDruEAA; __Secure-3PSIDTS=sidts-CjEBQT4rX28cVx_kTaTU86w8uiJVdgeCPt7OeBj6kiw6VkQbepj3NesQWdyekxjqHDruEAA; CONSISTENCY=AKreu9ttE1I3wjoxGirEb7mjUxe65Tlmg8OdAwy70Cszz_61zja1qQdyvBZ4sxAYfGrG7v89ILzdZ-VHHJL6-agrgRmhjPeSJ5VCJKtLo_EFvFNANljoSbuLURl2EXv32O1LRvChdCNJbg3-8AjmqUrJ; SIDCC=AKEyXzUM4MhVoAG0VtYx4S9TSISY0d77cdACNFlqoEP1MOKYCUx0nhbBxD7edjq517LwQJeP; __Secure-1PSIDCC=AKEyXzUU5GB_qY696ze_llin8SLUrbCqRdOOuG9z2R04-EyqiZXDlTdlLWJVU3alvwwvMWFWtw; __Secure-3PSIDCC=AKEyXzU3HfUBgIhVRczpVL5KQI7spY4WvvSif81p0ay3Y4jcyGzckts6bVnoLqAQXpwjGNYyPA; ST-2bzi54=oq=marrakech&gs_l=youtube.3..0i512i433i67i131i650k1j0i512i67i650k1j0i512i433k1j0i512i433i131k1j0i512i67i650k1j0i512i433i131k1j0i512i67i650k1j0i512i433i131k1l2j0i512k1j0i512i433i131k1j0i512i433k1j0i512i433i131i10k1j0i512i433k1.933.13928.0.17989.11.9.1.1.1.0.189.1074.0j8.9.0.ytsznadd10mni2%2Cytpo-bo-me%3D1%2Cytposo-bo-me%3D1%2Cytpo-bo-zo-mndr%3D15%2Cytposo-bo-zo-mndr%3D15%2Cytpo-bo-zo-epn%3D1%2Cytposo-bo-zo-epn%3D1%2Cytpo-bo-cwo-zrm%3D2%2Cytposo-bo-cwo-zrm%3D2%2Cytpo-bo-ndr-e%3D1%2Cytposo-bo-ndr-e%3D1%2Cytpo-bo-ndr-vsz%3D1%2Cytposo-bo-ndr-vsz%3D1%2Cytpo-bo-ndr-mi%3D51231707%2Cytposo-bo-ndr-mi%3D51231707%2Cytpo-bo-ro-mndr%3D1000%2Cytposo-bo-ro-mndr%3D1000%2Ccfro%3D1%2Cytpo-bo-me%3D0%2Cytposo-bo-me%3D0...0...1ac.1.64.youtube..1.10.1087.15..35i39i362k1j35i39k1j0i433i131i637k1j0i471k1j0i512i433i131i650k1j0i3k1j0i512i67i650i10k1j0i512i433i67i650i10k1j0i512i433i131i650i10k1.184.A5_O6rBGaQQ&itct=CBMQ7VAiEwjCwdXp9KeJAxWqPwYAHWfcHpY%3D&csn=vyFlaHOAQJv5Tm11&session_logininfo=AFmmF2swRQIgAYpzfY6Q6vCi-dCFs8MJuW2VkQvfMdYLs4UR7IB1hwECIQD2Zf6-G2QV7h3BgMqoTJy9IRUkxjJ7pilQlh7APa_bHg%3AQUQ3MjNmeC1EQnpfeDFUS2NVTVk5OW8yUUlUSG5fSzRjQ1E5NVN2ajU5VkJHU2NjRERSdVNMWVNGMF9JS0NRVGVHWDRfYzhDTUJpOHlEVGdzblVWS0QzT2RBQktwb0NLN3VBa3I2UEJ3UHFjQmppT3pSRkduZkhsOFV6Sl90aFN4MXFjUmN0Zmk0RGdfSWcxbVAwQXFBcFlJWGFWS2RycTZB&endpoint=%7B%22clickTrackingParams%22%3A%22CBMQ7VAiEwjCwdXp9KeJAxWqPwYAHWfcHpY%3D%22%2C%22commandMetadata%22%3A%7B%22webCommandMetadata%22%3A%7B%22url%22%3A%22%2Fresults%3Fsearch_query%3Dmarrakech%22%2C%22webPageType%22%3A%22WEB_PAGE_TYPE_SEARCH%22%2C%22rootVe%22%3A4724%7D%7D%2C%22searchEndpoint%22%3A%7B%22query%22%3A%22marrakech%22%7D%7D; ST-1k06sw0=session_logininfo=AFmmF2swRQIgAYpzfY6Q6vCi-dCFs8MJuW2VkQvfMdYLs4UR7IB1hwECIQD2Zf6-G2QV7h3BgMqoTJy9IRUkxjJ7pilQlh7APa_bHg%3AQUQ3MjNmeC1EQnpfeDFUS2NVTVk5OW8yUUlUSG5fSzRjQ1E5NVN2ajU5VkJHU2NjRERSdVNMWVNGMF9JS0NRVGVHWDRfYzhDTUJpOHlEVGdzblVWS0QzT2RBQktwb0NLN3VBa3I2UEJ3UHFjQmppT3pSRkduZkhsOFV6Sl90aFN4MXFjUmN0Zmk0RGdfSWcxbVAwQXFBcFlJWGFWS2RycTZB",
		"origin":                        "https://www.youtube.com",
		"pragma":                        "no-cache",
		"priority":                      "u=1, i",
		"referer":                       "https://www.youtube.com/results?search_query=marrakech",
		"sec-ch-ua":                     "\"Chromium\";v=\"130\", \"Google Chrome\";v=\"130\", \"Not?A_Brand\";v=\"99\"",
		"sec-ch-ua-arch":                "\"x86\"",
		"sec-ch-ua-bitness":             "\"64\"",
		"sec-ch-ua-form-factors":        "\"Desktop\"",
		"sec-ch-ua-full-version":        "\"130.0.6723.58\"",
		"sec-ch-ua-full-version-list":   "\"Chromium\";v=\"130.0.6723.58\", \"Google Chrome\";v=\"130.0.6723.58\", \"Not?A_Brand\";v=\"99.0.0.0\"",
		"sec-ch-ua-mobile":              "?0",
		"sec-ch-ua-model":               "\"\"",
		"sec-ch-ua-platform":            "\"Linux\"",
		"sec-ch-ua-platform-version":    "\"6.8.0\"",
		"sec-ch-ua-wow64":               "?0",
		"sec-fetch-dest":                "empty",
		"sec-fetch-mode":                "same-origin",
		"sec-fetch-site":                "same-origin",
		"user-agent":                    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		"x-client-data":                 "CIS2yQEIprbJAQipncoBCICXywEIk6HLAQiHoM0BCLKezgEI/aXOAQj7vs4BCIHDzgEIo8bOAQioyM4BCInJzgEI/srOAQiYy84BCMbMzgEY9cnNARicsc4BGIDKzgE=",
		"x-goog-authuser":               "3",
		"x-goog-visitor-id":             "CgtGWGVqTEdWWEJRNCiM6-q4BjIKCgJNQRIEGgAgLQ%3D%3D",
		"x-origin":                      "https://www.youtube.com",
		"x-youtube-bootstrap-logged-in": "true",
		"x-youtube-client-name":         "1",
		"x-youtube-client-version":      "2.20241024.01.00",
	}

	res, err := request.JSONReqest(
		"POST",
		"https://www.youtube.com/youtubei/v1/search?prettyPrint=false",
		jsonBytes,
		headersKeysValuesPairs,
		false,
	)
	if err != nil {
		return fmt.Errorf("scrapeYt: JSON request error: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("scrapeYt: failed to read response body: %w", err)
	}

	var bodyJson scrapeYtResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return fmt.Errorf("scrapeYt: failed to unmarshal the body: %w", err)
	}

	contents := bodyJson.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents
	fmt.Printf("contents: %+v\n", contents)
	if len(contents) == 0 {
		response := facebook.ResponseMessage{
			Text: "No result found. 🔭\nTry different keywords.",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	// lastIndex := len(contents) - 2 // VideoRenderer
	lastIndex := 0 // VideoRenderer
	things := contents[lastIndex].ItemSectionRenderer.Contents
	if len(things) <= 1 {
		response := facebook.ResponseMessage{
			Text: "No result found. 🔭\nTry different keywords.",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	comul := 0
	n := config.MaxReturnedVideo

	var wg sync.WaitGroup
	for _, thing := range things {
		if thing.VideoRenderer.VideoID == "" {
			continue
		}

		durationInSeconds := formats.DurationStrToSeconds(thing.VideoRenderer.LengthText.SimpleText)
		if durationInSeconds > config.MaxDurationInSeconds || durationInSeconds <= 0 {
			continue
		}

		n--
		comul++
		if n < 0 {
			break
		}

		title := "No title"
		if len(thing.VideoRenderer.Title.Runs) > 0 {
			title = thing.VideoRenderer.Title.Runs[0].Text
		}

		response := facebook.ResponseTemplateAttachment{
			Type: "template",
			Payload: facebook.TemplateAttachmentPayload{
				TemplateType: "generic",
				Elements: []facebook.TemplateAttachmentElement{
					{
						Title:    fmt.Sprintf("%d# %s", comul, title),
						Subtitle: thing.VideoRenderer.LengthText.SimpleText,
						ImageURL: fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", thing.VideoRenderer.VideoID),
						Buttons: []facebook.TemplateAttachmentButton{
							{
								Type:    "postback",
								Title:   "Watch now",
								Payload: fmt.Sprintf("YT_::_WATCH_::_%s_::_%d", thing.VideoRenderer.VideoID, durationInSeconds),
							},
							{
								Type:    "postback",
								Title:   "Listen now",
								Payload: fmt.Sprintf("YT_::_LISTEN_::_%s_::_%d", thing.VideoRenderer.VideoID, durationInSeconds),
							},
							{
								Type:    "postback",
								Title:   "Description",
								Payload: fmt.Sprintf("YT_::_DESCRIPTION_::_%s", thing.VideoRenderer.VideoID),
							},
						},
					},
				},
			},
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := facebookSender.CallSendAPI(sender_psid, response)
			if err != nil {
				fmt.Printf("Error sending Facebook message: %v\n", err)
			}
		}()
	}

	wg.Wait()

	if comul == 0 {
		response := facebook.ResponseMessage{
			Text: "All these videos are long 😥.\nTry different keywords.",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	return nil
}
