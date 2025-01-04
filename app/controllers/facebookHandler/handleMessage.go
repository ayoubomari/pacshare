package facebookHandler

import (
	"regexp"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/apk"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/pdf"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/wiki"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/yt"
	"github.com/ayoubomari/pacshare/app/models/facebook"
)

// "📺 .yt - Search and download YouTube videos\n" +
var MenuResponse = facebook.ResponseMessage{
	Text: "🚀 Welcome to PacShare! Here are the available commands:\n\n" +
		"📚 .wiki - Search Wikipedia articles\n" +
		"📄 .pdf - Search and download PDF documents\n\n",
	// "💡 Tip: You can also send a YouTube link directly!\n" +
	// "🔭 For more information, please visit https://fb.com/pacshare1/",
}

// handle message come from facebook (text message)
func handleMessage(sender_psid string, message facebook.Message) error {
	// work with lower case from now one
	trimmedMessage := strings.Trim(message.Text, " ")

	// redirect to specific regex handler

	// menu -> (?i)^(.menu)$
	menuRegex := regexp.MustCompile("(?i)^(.menu)$")
	if match := menuRegex.FindStringSubmatch(trimmedMessage); match != nil {
		//send message to sender with menu of all commands
		return facebookSender.CallSendAPI(sender_psid, MenuResponse)
	}

	// yt -> (?i)^(.yt) (.+)$
	ytRegex := regexp.MustCompile("(?i)^(.yt) (.+)$")
	if match := ytRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return yt.RegexHundlerMessage(sender_psid, match[1:])
	}

	// wiki -> (?i)^(.wiki)( |-)([a-z]{2}) (.+)$
	wikiRegex := regexp.MustCompile("(?i)^(.wiki)( |-)([a-z]{2}) (.+)$")
	if match := wikiRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return wiki.RegexHundlerMessage(sender_psid, match[1:])
	}
	wikiNoLangRegex := regexp.MustCompile("(?i)^(.wiki) (.+)$")
	if match := wikiNoLangRegex.FindStringSubmatch(trimmedMessage); match != nil {
		newMatch := make([]string, 5)
		newMatch[0] = match[0]
		newMatch[1] = ".wiki"
		newMatch[2] = "-"
		newMatch[3] = "en"
		newMatch[4] = match[2]
		return wiki.RegexHundlerMessage(sender_psid, newMatch[1:])
	}

	// apk -> (?i)^(.apk) (.+)$
	apkRegex := regexp.MustCompile("(?i)^(.apk) (.+)$")
	if match := apkRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return apk.RegexHundlerMessage(sender_psid, match[1:])
	}

	// pdf -> (?i)^(.pdf) (.+)$
	pdfRegex := regexp.MustCompile("(?i)^(.pdf) (.+)$")
	if match := pdfRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return pdf.RegexHundlerMessage(sender_psid, match[1:])
	}

	// // gemeni -> (?i)^(.(ask)) (.+)$
	// geminiRegex := regexp.MustCompile("(?i)^(.(gemini|ask)) (.+)$")
	// if match := geminiRegex.FindStringSubmatch(trimmedMessage); match != nil {
	// 	return gemini.RegexHundlerMessage(sender_psid, match[1:], message.Reply_to.MID)
	// }

	// default
	// return yt.RegexHundlerMessage(
	// 	sender_psid,
	// 	[]string{
	// 		".yt",
	// 		trimmedMessage,
	// 	},
	// )
	return facebookSender.CallSendAPI(sender_psid, MenuResponse)
}
