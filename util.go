package thracia

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

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
	verboseLog(ctx, "tempDir: %s\n", name)
	return context.WithValue(ctx, tdKey, name)
}

func tempDir(ctx context.Context) string {
	if v, ok := ctx.Value(tdKey).(string); ok {
		return v
	}
	return ""
}

func verboseLog(ctx context.Context, format string, args ...interface{}) {
	if c := cliContext(ctx); c.Bool("verbose") {
		fmt.Fprintf(c.App.Writer, format, args...)
	}
}

func pathInTempDir(ctx context.Context, path string) string {
	p := filepath.Join(tempDir(ctx), path)
	verboseLog(ctx, "path: %s\n", p)
	return p
}

func checkClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
