package thracia

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

const (
	migu1mURL  = "https://osdn.jp/frs/redir.php?m=gigenet&f=%2Fmix-mplus-ipa%2F63545%2Fmigu-1m-20150712.zip"
	migu1mFile = "migu1m.zip"
)

var migu1mTTFs = []string{"migu-1m-regular.ttf", "migu-1m-bold.ttf"}

// New returns the App instance.
func New() *cli.App {
	app := cli.NewApp()
	app.ErrWriter = cli.ErrWriter
	app.Name = "thracia"
	app.Usage = "Make fonts with SFMono + others"
	app.Version = "v0.1.0"
	app.Flags = flags()
	app.Action = action
	return app
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "clear, c",
			Usage: "Clear all downloads, and run",
		},
	}
}

func action(c *cli.Context) error {
	ctx := contextWithCLI(context.Background(), c)
	toDL := []*toDownload{}
	toDL = append(toDL, &toDownload{
		filename: migu1mFile,
		URL:      migu1mURL,
	})
	if err := clear(ctx, &toDL); err != nil {
		return fmt.Errorf("error in clear: %v", err)
	}
	if err := download(ctx, toDL); err != nil {
		return fmt.Errorf("error in download: %v", err)
	}
	if err := extract(ctx, migu1mFile, migu1mTTFs); err != nil {
		return fmt.Errorf("error in extract: %v", err)
	}
	return nil
}

func clear(ctx context.Context, toDL *[]*toDownload) error {
	cleared := make([]*toDownload, 0, len(*toDL))
	for _, dl := range *toDL {
		if _, err := os.Stat(dl.filename); os.IsNotExist(err) {
			cleared = append(cleared, dl)
		} else if err != nil {
			return fmt.Errorf("error in Stat: %v", err)
		} else if cliContext(ctx).Bool("clear") {
			if err := os.Remove(dl.filename); err != nil {
				return fmt.Errorf("error in Remove: %v", err)
			}
			cleared = append(cleared, dl)
		}
	}
	*toDL = cleared
	return nil
}

func extract(ctx context.Context, zipFile string, members []string) (err error) {
	z, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("error in OpenReader: %v", err)
	}
	defer checkClose(z, &err)
	for _, f := range z.File {
		for _, m := range members {
			if m != filepath.Base(f.Name) {
				continue
			}
			src, err := f.Open()
			if err != nil {
				return fmt.Errorf("error in Open: %v", err)
			}
			defer checkClose(src, &err)
			dst, err := os.Create(m)
			if err != nil {
				return fmt.Errorf("error in Create: %v", err)
			}
			defer checkClose(dst, &err)
			if _, err := io.Copy(dst, src); err != nil {
				return fmt.Errorf("error in Copy: %v", err)
			}
		}
	}
	return nil
}
