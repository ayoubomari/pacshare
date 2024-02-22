package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Your HTML string
	htmlString := `
		<html>
			<head>
				<title>Sample HTML</title>
			</head>
			<body>
				<div class="content">
					<p>Hello, goquery!</p>
				</div>
			</body>
		</html>
	`

	// Create a new reader for the HTML string
	reader := strings.NewReader(htmlString)

	// Use goquery to parse the HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Example: Accessing and printing the text inside the <p> tag
	paragraphText := doc.Find("p").Text()
	fmt.Println("Text inside <p> tag:", paragraphText)

	// Example: Accessing and printing the text inside the .content class
	contentText := doc.Find(".content").Text()
	fmt.Println("Text inside .content class:", contentText)
}
