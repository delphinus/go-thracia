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
	version               = "v0.1.0"
	modifyMigu1mTmpl      = "/assets/modify-migu1m.pe.tmpl"
	generateSFMonoModTmpl = "/assets/generate-sfmono-mod.pe.tmpl"
	generateObliqueTmpl   = "/assets/generate-oblique.pe.tmpl"
	fontPatcherTmpl       = "/assets/font-patcher.tmpl"
	changelogTmpl         = "/assets/changelog.md.tmpl"
	// SFMonoDir is a dir to store SFMono fonts
	SFMonoDir  = "/Applications/Utilities/Terminal.app/Contents/Resources/Fonts"
	familyName = "SFMono"
)

var migu1mTTFs = []string{
	"migu-1m-regular.ttf",
	"migu-1m-bold.ttf",
	"migu-1m-oblique.ttf",
	"migu-1m-bold-oblique.ttf",
}
var modifiedMigu1mTTFs = []string{
	"modified-migu-1m-regular.ttf",
	"modified-migu-1m-bold.ttf",
	"modified-migu-1m-oblique.ttf",
	"modified-migu-1m-bold-oblique.ttf",
}

// SFMonoTTFs is SFMono themselves.
var SFMonoTTFs = []string{
	"SFMono-Regular.otf",
	"SFMono-Bold.otf",
	"SFMono-RegularItalic.otf",
	"SFMono-BoldItalic.otf",
}

type h map[string]interface{}

// New returns the App instance.
func New() *cli.App {
	app := cli.NewApp()
	app.ErrWriter = cli.ErrWriter
	app.Name = "thracia"
	app.Usage = "Make fonts with SFMono + others"
	app.Version = version
	app.Flags = flags()
	app.Action = action
	return app
}

func flags() []cli.Flag {
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print logs verbosely",
		},
		cli.StringFlag{
			Name:  "suffix, s",
			Usage: "set fontfamily suffix",
		},
		cli.BoolFlag{
			Name:  "download-only",
			Usage: "it prepares files only, and does not generate (default: false)",
		},
		cli.BoolTFlag{
			Name:  "nerd-fonts, n",
			Usage: "generate fonts with nerd-fonts (default: true)",
		},
		cli.BoolTFlag{
			Name:  "bold, b",
			Usage: "generate bold fonts (default: true)",
		},
		cli.BoolTFlag{
			Name:  "italic, i",
			Usage: "generate italic fonts (default: true)",
		},
	}
}

func action(c *cli.Context) error {
	ctx := contextWithCLI(context.Background(), c)
	ctx = contextWithTempDir(ctx)
	if err := copySFMono(ctx); err != nil {
		return fmt.Errorf("error in copySFMono: %v", err)
	}
	toDL, err := files(ctx)
	if err != nil {
		return fmt.Errorf("error in files: %v", err)
	}
	if err := download(ctx, toDL); err != nil {
		return fmt.Errorf("error in download: %v", err)
	}
	if err := extract(ctx, migu1mFile, migu1mTTFs); err != nil {
		return fmt.Errorf("error in extract: %v", err)
	}
	if c.Bool("download-only") {
		return nil
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
			dst, err := os.Create(pathInTempDir(ctx, m))
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
