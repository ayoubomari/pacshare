package facebookHandler

import (
	"regexp"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/scrapers/apk"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/pdf"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/wiki"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/yt"
	"github.com/ayoubomari/pacshare/app/models/facebook"
)

// handle message come from facebook (text message)
func handleMessage(sender_psid string, message facebook.Message) error {
	// work with lower case from now one
	trimmedMessage := strings.Trim(message.Text, " ")

	// redirect to specific regex handler

	// yt -> (?i)^(.yt) (.+)$
	ytRegex := regexp.MustCompile("(?i)^(.yt) (.+)$")
	if match := ytRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return yt.RegexHundlerMessage(sender_psid, match[1:])
	}

	// wiki -> (?i)^(.wiki)( |-)([a-z]{2}) (.+)$
	wikiRegex := regexp.MustCompile("(?i)^(.wiki)( |-)([a-z]{2}) (.+)$")
	if match := wikiRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return wiki.RegexHundlerMessage(match[1:])
	}

	// apk -> (?i)^(.apk) (.+)$
	apkRegex := regexp.MustCompile("(?i)^(.apk) (.+)$")
	if match := apkRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return apk.RegexHundlerMessage(match[1:])
	}

	// pdf -> (?i)^(.pdf) (.+)$
	pdfRegex := regexp.MustCompile("(?i)^(.pdf) (.+)$")
	if match := pdfRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return pdf.RegexHundlerMessage(match[1:])
	}

	//default
	return yt.RegexHundlerMessage(
		sender_psid,
		[]string{
			".yt",
			trimmedMessage,
		},
	)
}
