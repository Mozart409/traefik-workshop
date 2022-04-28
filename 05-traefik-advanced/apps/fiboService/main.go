package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Fibonacci(n int) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	var n2, n1 uint64 = 0, 1
	for i := int(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}

func main() {

	// jaeger.docker.localhost:16686

	app := fiber.New()

	// app.Use(fibertracing.New())

	app.Get("/", func(c *fiber.Ctx) error {

		return c.SendString("Hello, World from Fiber!")
	})

	app.Get("/fibonacci/:n", func(c *fiber.Ctx) error {

		n, err := c.ParamsInt("n")
		if err != nil {
			return c.Status(400).SendString("Invalid parameter")
		}

		fib, err := Fibonacci(n)
		if err != nil {
			return c.Status(500).SendString("Internal server error")
		}

		return c.SendString(fmt.Sprintf("%d", fib))
	})

	app.Listen(":3333")
}
