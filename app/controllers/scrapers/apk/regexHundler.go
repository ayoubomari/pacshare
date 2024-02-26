package apk

import (
	"errors"
	"regexp"
	"strings"
)

func RegexHundlerMessage(sender_psid string, subStrings []string) error {
	if len(subStrings) < 2 {
		return errors.New("subStrings length are low then 2")
	}

	return scrapeApk(sender_psid, subStrings[1])
}

func RegexHundlerPostback(sender_psid string, subStrings []string) error {
	// DOWNLOAD -> ^(DOWNLOAD_::_)(.+)$
	DOWNLOADRegex := regexp.MustCompile("^(DOWNLOAD_::_)(.+)$")
	if match := DOWNLOADRegex.FindStringSubmatch(subStrings[1]); match != nil {
		return downloadApk(sender_psid, strings.Split(match[2], "_::_"))
	}
	// PHOTO -> ^(PHOTO_::_)(.+)$
	PHOTORegex := regexp.MustCompile("^(PHOTO_::_)(.+)$")
	if match := PHOTORegex.FindStringSubmatch(subStrings[1]); match != nil {
		return photoApk(sender_psid, strings.Split(match[2], "_::_"))
	}
	// DESCRIPTION -> ^(DESCRIPTION_::_)(.+)$
	DESCRIPTIONRegex := regexp.MustCompile("^(DESCRIPTION_::_)(.+)$")
	if match := DESCRIPTIONRegex.FindStringSubmatch(subStrings[1]); match != nil {
		return descriptionApk(sender_psid, strings.Split(match[2], "_::_"))
	}
	return errors.New("the postback isn't compatible with any available postback")
}
