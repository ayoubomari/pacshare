package yt

import (
	"errors"
	"regexp"
	"strings"
)

// yt -> (?i)^(.yt) (.+)$
func RegexHundlerMessage(sender_psid string, subStrings []string) error {
	if len(subStrings) < 2 {
		return errors.New("subStrings length are low then 2")
	}

	// check if subStrings is a Youtube link
	if strings.Contains(subStrings[1], "watch?v=") || strings.Contains(subStrings[1], "youtu.be/") || strings.Contains(subStrings[1], "shorts/") {
		return scrapeYtByVideoUrl(sender_psid, subStrings[1])
	}

	// if the keywords are not a Youtube link
	return scrapeYt(sender_psid, subStrings[1])
}

// subStrings: [YT_::_, subCommandStringPayload]
func RegexHundlerPostback(sender_psid string, subStrings []string) error {
	// WATCH -> ^(WATCH_::_)(.+)$
	WATCHRegex := regexp.MustCompile("^(WATCH_::_)(.+)$")
	if match := WATCHRegex.FindStringSubmatch(subStrings[1]); match != nil {
		return watchYt(sender_psid, strings.Split(match[2], "_::_"))
	}
	// LISTEN -> ^(LISTEN_::_)(.+)$
	LISTENRegex := regexp.MustCompile("^(LISTEN_::_)(.+)$")
	if match := LISTENRegex.FindStringSubmatch(subStrings[1]); match != nil {
		return listenYt(sender_psid, strings.Split(match[2], "_::_"))
	}
	// RECOMMENDED_VIDEO -> ^(RECOMMENDED_VIDEO_::_)(.+)$
	RECOMMENDED_VIDEORegex := regexp.MustCompile("^(RECOMMENDED_VIDEO_::_)(.+)$")
	if match := RECOMMENDED_VIDEORegex.FindStringSubmatch(subStrings[1]); match != nil {
		return recommendedVideos(sender_psid, strings.Split(match[2], "_::_"))
	}
	// DESCRIPTION -> ^(DESCRIPTION_::_)(.+)$
	DESCRIPTIONRegex := regexp.MustCompile("^(DESCRIPTION_::_)(.+)$")
	if match := DESCRIPTIONRegex.FindStringSubmatch(subStrings[1]); match != nil {
		return VideoDescription(sender_psid, strings.Split(match[2], "_::_"))
	}
	return errors.New("the postback isn't compatible with any available postback")
}
