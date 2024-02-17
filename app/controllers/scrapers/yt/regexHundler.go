package yt

import "errors"

func RegexHundlerMessage(sender_psid string, subStrings []string) error {
	if len(subStrings) < 2 {
		return errors.New("subStrings length are low then 2")
	}
	return scrapeYt(sender_psid, subStrings[1])
}

func RegexHundlerPostback(subStrings []string) error {
	return nil
}
