package wiki

import (
	"errors"

	"github.com/ayoubomari/pacshare/config"
)

// wiki -> (?i)^(.wiki)( |-)([a-z]{2}) (.+)$
func RegexHundlerMessage(sender_psid string, subStrings []string) error {
	if len(subStrings) < 4 {
		return errors.New("subStrings length are low then 2")
	}

	// Check if language present in WikiSepportedLanguages array
	found := false
	for _, lang := range config.WikiSepportedLanguages {
		if lang == subStrings[2] {
			found = true
			break
		}
	}

	// if it's not found
	if !found {
		subStrings[2] = "en"
	}

	return scrapeWiki(sender_psid, subStrings[2:])
}

// subStrings: [WIKI_::_, subCommandStringPayload]
func RegexHundlerPostback(subStrings []string) error {
	return nil
}
