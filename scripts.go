package thracia

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func scripts(ctx context.Context) error {
	c := cliContext(ctx)
	fontforge, err := exec.LookPath("fontforge")
	if err != nil {
		return fmt.Errorf("cannot find `fontforge` executable")
	}

	if script, err := generateScripts(ctx, generateObliqueTmpl, h{
		"FontForge":      fontforge,
		"SrcRegular":     migu1mTTFs[0],
		"SrcBold":        migu1mTTFs[1],
		"FilenameFamily": "migu-1m",
		"FontFamily":     "Migu 1M",
	}); err != nil {
		return fmt.Errorf("error in script: %v", err)
	} else if err := execScripts(ctx, script); err != nil {
		return fmt.Errorf("error in execScripts: %v", err)
	}

	// Glyphs of SFMono has this metrics:
	// h: 1638 + 410 = 2048
	// w: 1266
	// So zenkaku glyphs should have padding on left and right:
	// (1266 * 2 - 2048) / 2 = 242
	if script, err := generateScripts(ctx, modifyMigu1mTmpl, h{
		"FontForge":  fontforge,
		"Ascent":     1638,
		"Descent":    410,
		"Padding":    242,
		"SrcRegular": migu1mTTFs[0],
		"SrcBold":    migu1mTTFs[1],
		"DstRegular": modifiedMigu1mTTFs[0],
		"DstBold":    modifiedMigu1mTTFs[1],
		"DstDir":     tempDir(ctx),
	}); err != nil {
		return fmt.Errorf("error in generateScripts: %v", err)
	} else if err := execScripts(ctx, script); err != nil {
		return fmt.Errorf("error in execScripts: %v", err)
	}

	if script, err := generateScripts(ctx, generateSFMonoModTmpl, h{
		"FontForge":         fontforge,
		"SFMonoRegular":     SFMonoTTFs[0],
		"SFMonoBold":        SFMonoTTFs[1],
		"Migu1mRegular":     modifiedMigu1mTTFs[0],
		"Migu1mBold":        modifiedMigu1mTTFs[1],
		"FamilyName":        familyName,
		"FamilyNameSuffix":  c.String("suffix"),
		"Version":           version,
		"WinAscent":         1950,
		"WinDescent":        494,
		"Ascent":            1638,
		"Descent":           410,
		"ZenkakuSpaceGlyph": "",
	}); err != nil {
		return fmt.Errorf("error in generateScripts: %v", err)
	} else if err := execScripts(ctx, script); err != nil {
		return fmt.Errorf("error in execScripts: %v", err)
	}

	if err := execFontPatcher(ctx); err != nil {
		return fmt.Errorf("error in execFontPatcher: %v", err)
	}
	return nil
}

func scriptFilename(ctx context.Context, tmpl string) string {
	filename := strings.TrimSuffix(filepath.Base(tmpl), ".tmpl")
	return pathInTempDir(ctx, filename)
}

func generateScripts(ctx context.Context, tmpl string, data interface{}) (script string, err error) {
	f, err := Assets.Open(tmpl)
	if err != nil {
		return "", fmt.Errorf("error in Open: %v", err)
	}
	defer checkClose(f, &err)
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("error in ReadAll: %v", err)
	}
	t, err := template.New("").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("error in Parse: %v", err)
	}
	script = scriptFilename(ctx, tmpl)
	dst, err := os.OpenFile(script, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return "", fmt.Errorf("error in Create: %v", err)
	}
	defer checkClose(dst, &err)
	if err := t.Execute(dst, data); err != nil {
		return "", fmt.Errorf("error in Execute: %v", err)
	}
	return
}

func execScripts(ctx context.Context, script string, args ...string) error {
	c := cliContext(ctx)
	if c.Bool("verbose") {
		fmt.Fprintf(c.App.Writer, "%s %v\n", script, args)
	}
	cmd := exec.CommandContext(ctx, script, args...)
	cmd.Dir = tempDir(ctx)
	cmd.Stdout = c.App.Writer
	cmd.Stderr = c.App.ErrWriter
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error in Run: %v", err)
	}
	return nil
}

func execFontPatcher(ctx context.Context) error {
	fp := pathInTempDir(ctx, fontPatcher)
	if err := os.Chmod(fp, 0755); err != nil {
		return fmt.Errorf("error in Chmod: %v", err)
	}
	mod := modified(SFMonoTTFs[0], cliContext(ctx).String("suffix"))
	if err := execScripts(ctx, fp, "-c", "-q", "-out", "build", mod); err != nil {
		return fmt.Errorf("error in execScripts: %v", err)
	}
	return nil
}

func modified(orig, suffix string) string {
	return strings.Replace(orig, familyName, familyName+suffix, 1)
}
