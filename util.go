package thracia

import (
	"context"

	"github.com/urfave/cli"
)

type cliContextKey struct{}

var ccKey = cliContextKey{}

func contextWithCLI(ctx context.Context, c *cli.Context) context.Context {
	return context.WithValue(ctx, ccKey, c)
}

func cliContext(ctx context.Context) *cli.Context {
	if v, ok := ctx.Value(ccKey).(*cli.Context); ok {
		return v
	}
	return nil
}
