package yt

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

type scrapeYtResponseBody struct {
	Contents *struct {
		TwoColumnSearchResultsRenderer *struct {
			PrimaryContents *struct {
				SectionListRenderer *struct {
					Contents []struct {
						ItemSectionRenderer *struct {
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
								}
							} `json:"contents"`
						} `json:"itemSectionRenderer"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

func scrapeYt(sender_psid string, searchKeyWords string) error {
	jsonBytes := []byte(`{"context":{"client":{"hl":"en","gl":"MA","remoteHost":"41.141.106.74","deviceMake":"","deviceModel":"","visitorData":"CgtYUEFKMFBFX05HdyjY_fWLBg%3D%3D","userAgent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36,gzip(gfe)","clientName":"WEB","clientVersion":"2.20211029.00.00","osName":"Windows","osVersion":"10.0","originalUrl":"https://www.youtube.com/resultssearch_query=youtube&sp=EgQQAVgD","platform":"DESKTOP","clientFormFactor":"UNKNOWN_FORM_FACTOR","configInfo":{"appInstallData":"CNj99YsGEODWrQUQktWtBRC3y60FELDUrQUQ47v9EhDU0K0FENi-rQUQkfj8Eg%3D%3D"},"timeZone":"Africa/Casablanca","browserName":"Chrome","browserVersion":"95.0.4638.54","screenWidthPoints":1440,"screenHeightPoints":241,"screenPixelDensity":1,"screenDensityFloat":1,"utcOffsetMinutes":60,"userInterfaceTheme":"USER_INTERFACE_THEME_LIGHT","mainAppWebInfo":{"graftUrl":"/resultssearch_query=hello+world&sp=EgQQAVgD","webDisplayMode":"WEB_DISPLAY_MODE_BROWSER","isWebNativeShareAvailable":true}},"user":{"enableSafetyMode":true,"lockedSafetyMode":false},"request":{"useSsl":true,"internalExperimentFlags":[],"consistencyTokenJars":[]},"clickTracking":{"clickTrackingParams":"CDYQk3UYACITCJScqPnT8vMCFYsQBgAdoiQKCQ=="},"adSignalsInfo":{"params":[{"key":"dt","value":"1635614426721"},{"key":"flash","value":"0"},{"key":"frm","value":"0"},{"key":"u_tz","value":"60"},{"key":"u_his","value":"7"},{"key":"u_h","value":"900"},{"key":"u_w","value":"1440"},{"key":"u_ah","value":"860"},{"key":"u_aw","value":"1440"},{"key":"u_cd","value":"24"},{"key":"bc","value":"31"},{"key":"bih","value":"241"},{"key":"biw","value":"1424"},{"key":"brdim","value":"0,0,0,0,1440,0,1440,860,1440,241"},{"key":"vis","value":"1"},{"key":"wgl","value":"true"},{"key":"ca_type","value":"image"}]}},"query": "` + searchKeyWords + `","params":"EgQQAVgD"}`)

	res, err := request.JSONReqest(
		"POST",
		"https://www.youtube.com/youtubei/v1/search?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false",
		jsonBytes,
		make(map[string]string),
	)
	if err != nil {
		return fmt.Errorf("scrapeYt: json requeset err: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("scrapeYt: fail to read res.body %w", err)
	}
	var bodyJson scrapeYtResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return fmt.Errorf("scrapeYt: failt to unmarshale the body %w", err)
	}

	if len(bodyJson.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents)-2 < 0 {
		response := facebook.ResponseMessage{
			Text: "No result found. 🔭\n" +
				"Try different keywords.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
		return nil
	}
	things := bodyJson.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents[len(bodyJson.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents)-2].ItemSectionRenderer.Contents
	if things == nil || len(things) <= 1 {
		response := facebook.ResponseMessage{
			Text: "No result found. 🔭\n" +
				"Try different keywords.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
		return nil
	}

	comul := 0
	n := config.MaxReturnedVideo
	for i := 0; i < len(things); i++ {
		if things[i].VideoRenderer.VideoID != "" {
			durationInSeconds := formats.DurationStrToSeconds(things[i].VideoRenderer.LengthText.SimpleText)
			if durationInSeconds <= config.MaxDurationInSeconds && durationInSeconds > 0 {
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
								Title:    fmt.Sprintf("%d# %s", comul, things[i].VideoRenderer.Title.Runs[0].Text),
								Subtitle: things[i].VideoRenderer.LengthText.SimpleText,
								ImageURL: fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", things[i].VideoRenderer.VideoID),
								Buttons: []facebook.TemplateAttachmentButton{
									{
										Type:    "postback",
										Title:   "Watch now",
										Payload: fmt.Sprintf("YT_::_WATCH_::_%s_::_%d", things[i].VideoRenderer.VideoID, durationInSeconds),
									},
									{
										Type:    "postback",
										Title:   "Listen now",
										Payload: fmt.Sprintf("YT_::_LISTEN_::_%s_::_%d", things[i].VideoRenderer.VideoID, durationInSeconds),
									},
									{
										Type:    "postback",
										Title:   "Description",
										Payload: fmt.Sprintf("YT_::_DESCRIPTION_::_%s", things[i].VideoRenderer.VideoID),
									},
								},
							},
						},
					},
				}

				go facebookSender.CallSendAPI(sender_psid, response)
			}
		}
	}
	if comul == 0 {
		response := facebook.ResponseMessage{
			Text: "All these videos are long 😥.\n" +
				"Try different keywords.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
	}

	return nil
}
