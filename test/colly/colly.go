package main

import (
	"fmt"
	"strings"
)

func main() {
	urls := []string{
		"https://www.youtube.com/shorts/UUb68mu2G1M?hellno=wtf&ok=yes",
		"https://www.youtube.com/watch?v=KNHEXOoV-H4&hellno=wtf&ok=yes",
		"https://youtu.be/KNHEXOoV-H4?si=Jzqw82ONyOKk8i9e/hellno?wtf&ok",
	}

	for _, url := range urls {
		var videoID string
		if strings.Contains(url, "watch?v=") {
			url2 := strings.ReplaceAll(url, "&", "watch?v=")
			videoID = strings.Split(url2, "watch?v=")[1]
		} else if strings.Contains(url, "youtu.be/") {
			videoID = strings.Split(url, "youtu.be/")[1]
			videoID = strings.Split(videoID, "?")[0]
		} else if strings.Contains(url, "shorts/") {
			videoID = strings.Split(url, "shorts/")[1]
			videoID = strings.Split(videoID, "?")[0]
		}
		fmt.Println("videoID:", videoID)
	}
}
