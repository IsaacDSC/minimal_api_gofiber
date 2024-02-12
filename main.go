package main

import (
	"context"
	"errors"
	"log"

	o11yfiber "github.com/IsaacDSC/O11Y_fiber"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/monitor"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const serviceName = "minimal-api"

func main() {
	tp := o11yfiber.StartTracing(o11yfiber.TracingConfig{
		EndpointCollector: "http://localhost:14268/api/traces",
		ServiceNameKey:    serviceName,
	})

	log.Fatal(o11yfiber.StartServerHttp(o11yfiber.SettingsHttp{
		TracerProvider:     tp,
		ServiceNameMetrics: serviceName,
		Handlers: []o11yfiber.Handler{
			{HandlerFunc: handleUser, Path: "/users/:id", Method: o11yfiber.GET},
			{HandlerFunc: handleError, Path: "/error", Method: o11yfiber.GET},
			// {HandlerFunc: monitor.New(), Path: "/fiber/metrics", Method: o11yfiber.GET},
		},
		Middleware: []func(c *fiber.Ctx) error{
			o11yfiber.MiddlewareIO(),
		},
		ServerPort: 3333,
	}))

}

func handleError(ctx *fiber.Ctx) error {
	return errors.New("abc")
}

func handleUser(c *fiber.Ctx) error {
	id := c.Params("id")
	name := getUser(c.UserContext(), id)
	return c.JSON(fiber.Map{"id": id, name: name})
}

func getUser(ctx context.Context, id string) string {
	_, span := o11yfiber.Span().Start(ctx, "getUser", oteltrace.WithAttributes(attribute.String("id", id)))
	defer span.End()
	if id == "123" {
		return "otelfiber tester"
	}
	return "unknown"
}
