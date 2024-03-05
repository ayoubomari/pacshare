package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	agent := fiber.Get("https://upload.wikimedia.org/wikipedia/commons/thumb/c/c1/Gosforth_Cross_Loki_and_Sigyn.jpg/1080px-Gosforth_Cross_Loki_and_Sigyn.jpeg")
	statusCode, body, errs := agent.Bytes()
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

	fmt.Println("body:", string(body))
	// Release the response to free up resources
	defer fiber.ReleaseResponse(resp)
}
