package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Get("/:key/*", ProxyHandler)
	//app.Delete("cache/:key/*", EvictCacheHandler)
	//app.Delete("limit/:key/*", ResetLimitHandler)

	err := app.Listen(":3000")
	if err != nil {
		fmt.Println(err)
	}
}
