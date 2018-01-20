package thracia

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/urfave/cli"
)

type cliContextKey struct{}

var ccKey = cliContextKey{}

type tempDirKey struct{}

var tdKey = tempDirKey{}

func contextWithCLI(ctx context.Context, c *cli.Context) context.Context {
	return context.WithValue(ctx, ccKey, c)
}

func cliContext(ctx context.Context) *cli.Context {
	if v, ok := ctx.Value(ccKey).(*cli.Context); ok {
		return v
	}
	return nil
}

func contextWithTempDir(ctx context.Context) context.Context {
	name, err := ioutil.TempDir("", "")
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, tdKey, name)
}

func tempDir(ctx context.Context) string {
	if v, ok := ctx.Value(tdKey).(string); ok {
		return v
	}
	return ""
}

func checkClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
