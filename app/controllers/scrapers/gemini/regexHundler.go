package gemini

import (
	"errors"
)

// gemini -> (?i)^(.(gemini|bard)) (.+)$
func RegexHundlerMessage(sender_psid string, subStrings []string, mid string) error {
	if len(subStrings) < 2 {
		return errors.New("subStrings length are low then 2")
	}

	return scrapeGemini(sender_psid, subStrings[2:], mid)
}

// subStrings: [GEMINI_::_, subCommandStringPayload]
func RegexHundlerPostback(subStrings []string) error {
	return nil
}
