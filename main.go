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
	migu1mURL        = "https://osdn.jp/frs/redir.php?m=gigenet&f=%2Fmix-mplus-ipa%2F63545%2Fmigu-1m-20150712.zip"
	migu1mFile       = "migu1m.zip"
	modifyMigu1mTmpl = "/assets/modify-migu1m.pe.tmpl"
	// SFMonoDir is a dir to store SFMono fonts
	SFMonoDir = "/Applications/Utilities/Terminal.app/Contents/Resources/Fonts"
)

var migu1mTTFs = []string{"migu-1m-regular.ttf", "migu-1m-bold.ttf"}
var modifiedMigu1mTTFs = []string{
	"modified-migu-1m-regular.ttf",
	"modified-migu-1m-bold.ttf",
}

// SFMonoTTFs is SFMono themselves.
var SFMonoTTFs = []string{
	"SFMono-Regular.otf",
	"SFMono-RegularItalic.otf",
	"SFMono-Bold.otf",
	"SFMono-BoldItalic.otf",
}

type h map[string]interface{}

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
			Name:  "verbose, vv",
			Usage: "Print logs verbosely",
		},
		cli.BoolTFlag{
			Name:  "scale-down, s",
			Usage: "Scale down Migu 1M (default: true)",
		},
	}
}

func action(c *cli.Context) error {
	ctx := contextWithCLI(context.Background(), c)
	ctx = contextWithTempDir(ctx)
	toDL := []*toDownload{}
	toDL = append(toDL, &toDownload{
		filename: migu1mFile,
		URL:      migu1mURL,
	})
	if err := download(ctx, toDL); err != nil {
		return fmt.Errorf("error in download: %v", err)
	}
	if err := extract(ctx, migu1mFile, migu1mTTFs); err != nil {
		return fmt.Errorf("error in extract: %v", err)
	}
	if err := copySFMono(ctx); err != nil {
		return fmt.Errorf("error in copySFMono: %v", err)
	}
	if err := scripts(ctx); err != nil {
		return fmt.Errorf("error in scripts: %v", err)
	}
	return nil
}

func extract(ctx context.Context, zipFile string, members []string) (err error) {
	z, err := zip.OpenReader(pathInTempDir(ctx, zipFile))
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

func copySFMono(ctx context.Context) (err error) {
	for _, f := range SFMonoTTFs {
		font := filepath.Join(SFMonoDir, f)
		if _, err := os.Stat(font); os.IsNotExist(err) {
			return fmt.Errorf("cannot find font: %s", font)
		} else if err != nil {
			return fmt.Errorf("error in Stat: %v", err)
		}
		src, err := os.Open(font)
		if err != nil {
			return fmt.Errorf("error in Open: %v", err)
		}
		dst, err := os.Create(pathInTempDir(ctx, f))
		if err != nil {
			return fmt.Errorf("error in Create: %v", err)
		}
		defer checkClose(dst, &err)
		if _, err := io.Copy(dst, src); err != nil {
			return fmt.Errorf("error in Copy: %v", err)
		}
	}
	return nil
}