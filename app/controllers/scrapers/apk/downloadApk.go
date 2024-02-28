package apk

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/fileDownloader"
	"github.com/ayoubomari/pacshare/util/formats"
)

func downloadApk(sender_psid string, arguments []string) error {
	fmt.Println("from downloadAPK")
	if len(arguments) < 1 {
		return errors.New("arguments length is lower then 1")
	}
	apkId := arguments[0]
	fmt.Println("apkId:", apkId)

	appInfo, err := GetApkInfoWS2(apkId)
	if err != nil {
		response := facebook.ResponseMessage{
			Text: "Something wrong try another time 🙁.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
		return fmt.Errorf("downloadApk: %w", err)
	}

	fmt.Println("appInfo.Nodes.Meta.Data.File.Path:", appInfo.Nodes.Meta.Data.File.Path)
	fmt.Println(strings.ReplaceAll(appInfo.Nodes.Meta.Data.File.Path, "pool.", "premium."))
	// download and send apk
	fileDownloader.DownloadAndSendFile(sender_psid, strings.ReplaceAll(appInfo.Nodes.Meta.Data.File.Path, "pool.", "premium."), "./public/src/apks/", formats.ToFileNameString(appInfo.Nodes.Meta.Data.Uname), "_pac.apk", appInfo.Nodes.Meta.Data.File.Filesize, config.ApkChunksMaxSize, "file")
	//send apk complition response message
	response := facebook.ResponseMessage{
		Text: "All apk files have been sent. ✅",
	}
	facebookSender.CallSendAPI(sender_psid, response)

	// if there is an obb file
	if appInfo.Nodes.Meta.Data.Obb != nil {
		// download and send obb
		fileDownloader.DownloadAndSendFile(sender_psid, strings.ReplaceAll(appInfo.Nodes.Meta.Data.Obb.Main.Path, "pool", "premium"), "./public/src/apks/", formats.ToFileNameString(strings.ReplaceAll(appInfo.Nodes.Meta.Data.Obb.Main.Filename, ".obb", "")), "_pac.obb", appInfo.Nodes.Meta.Data.File.Filesize, config.ApkChunksMaxSize, "file")

		//send obb complition response message
		response = facebook.ResponseMessage{
			Text: "All obb files have been sent. ✅\n" +
				fmt.Sprintf("obb folder name: %s\n", appInfo.Nodes.Meta.Data.Obb.Main.Filename[strings.Index(appInfo.Nodes.Meta.Data.Obb.Main.Filename, "com."):strings.Index(appInfo.Nodes.Meta.Data.Obb.Main.Filename, ".obb")]) +
				fmt.Sprintf("obb file name: %s", appInfo.Nodes.Meta.Data.Obb.Main.Filename),
		}
		facebookSender.CallSendAPI(sender_psid, response)
	}

	return nil
}
