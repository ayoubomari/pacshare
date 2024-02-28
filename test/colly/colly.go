package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	agent := fiber.Get("https://z45st682fx.zlib-cdn.com/dl2.php?id=36551844&h=0d482032919c1a052d06d96ed8df7bb7&u=cache&ext=pdf&n=Living%20in%20the%20light%20a%20guide%20to%20personal%20transformation")
	statusCode, _, errs := agent.Bytes()
	if len(errs) > 0 {
		fmt.Println(errs)
		return
	}
	fmt.Println("statusCode:", statusCode)

	resp := fiber.AcquireResponse()
	agent.SetResponse(resp)

	// Visit and print all the headers in the response
	resp.Header.VisitAll(func(key, value []byte) {
		fmt.Println("Header", string(key), "value", string(value))
	})

	contentLength := resp.Header.ContentLength()
	fmt.Println("contentLength:", contentLength)
	// Release the response to free up resources
	defer fiber.ReleaseResponse(resp)
}
