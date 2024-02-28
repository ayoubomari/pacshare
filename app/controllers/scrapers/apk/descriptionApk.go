package apk

import (
	"errors"
	"fmt"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
)

func descriptionApk(sender_psid string, arguments []string) error {
	if len(arguments) < 1 {
		return errors.New("arguments length is lower then 1")
	}
	apkId := arguments[0]

	appInfo, err := GetApkInfoWS2(apkId)
	if err != nil {
		response := facebook.ResponseMessage{
			Text: "Something wrong try another time 🙁.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
		return fmt.Errorf("descriptionApk: %w", err)
	}

	// send description text
	if appInfo.Nodes.Meta.Data.Media.Description == "" {
		response := facebook.ResponseMessage{
			Text: "No description was found. 🤷‍♂️",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	facebookSender.SendMessageByChunks(sender_psid, appInfo.Nodes.Meta.Data.Media.Description)

	return nil
}
