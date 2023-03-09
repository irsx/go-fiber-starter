package middlewares

import (
	"bufio"
	"go-fiber-starter/utils/sse"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func SseMiddleware(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("Transfer-Encoding", "chunked")

	ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		sse.BrokerList.Subscribe(w)
		sse.BrokerList.Listen()
	}))

	return nil
}
