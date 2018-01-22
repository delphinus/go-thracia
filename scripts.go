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
	fontforge, err := exec.LookPath("fontforge")
	if err != nil {
		return fmt.Errorf("cannot find `fontforge` executable")
	}
	if err := generateOblique(ctx, fontforge); err != nil {
		return fmt.Errorf("error in generateOblique: %v", err)
	}
	if err := modifyMigu1m(ctx, fontforge); err != nil {
		return fmt.Errorf("error in modifyMigu1m: %v", err)
	}
	if err := generateSFMonoMod(ctx, fontforge); err != nil {
		return fmt.Errorf("error in generateSFMonoMod: %v", err)
	}
	if err := execFontPatcher(ctx); err != nil {
		return fmt.Errorf("error in execFontPatcher: %v", err)
	}
	return nil
}

func generateOblique(ctx context.Context, fontforge string) error {
	c := cliContext(ctx)
	if !c.BoolT("italic") {
		return nil
	}
	filenames := ""
	if c.BoolT("bold") {
		filenames = fmt.Sprintf(`"%s", "%s"`, migu1mTTFs[0], migu1mTTFs[1])
	} else {
		filenames = fmt.Sprintf(`"%s"`, migu1mTTFs[0])
	}
	if script, err := generateScripts(ctx, generateObliqueTmpl, h{
		"FontForge":      fontforge,
		"Filenames":      filenames,
		"FilenameFamily": "migu-1m",
		"FontFamily":     "Migu 1M",
	}); err != nil {
		return fmt.Errorf("error in script: %v", err)
	} else if err := execScripts(ctx, script); err != nil {
		return fmt.Errorf("error in execScripts: %v", err)
	}
	return nil
}

func modifyMigu1m(ctx context.Context, fontforge string) error {
	c := cliContext(ctx)
	var inputs, outputs string
	if c.BoolT("bold") && c.BoolT("italic") {
		inputs = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
			migu1mTTFs[0], migu1mTTFs[1], migu1mTTFs[2], migu1mTTFs[3])
		outputs = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
			modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[1], modifiedMigu1mTTFs[2], modifiedMigu1mTTFs[3])
	} else if c.BoolT("bold") {
		inputs = fmt.Sprintf(`"%s", "%s"`, migu1mTTFs[0], migu1mTTFs[1])
		outputs = fmt.Sprintf(`"%s", "%s"`, modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[1])
	} else if c.BoolT("italic") {
		inputs = fmt.Sprintf(`"%s", "%s"`, migu1mTTFs[0], migu1mTTFs[2])
		outputs = fmt.Sprintf(`"%s", "%s"`, modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[2])
	} else {
		inputs = fmt.Sprintf(`"%s"`, migu1mTTFs[0])
		outputs = fmt.Sprintf(`"%s"`, modifiedMigu1mTTFs[0])
	}
	// Glyphs of SFMono has this metrics:
	// h: 1638 + 410 = 2048
	// w: 1266
	// So zenkaku glyphs should have padding on left and right:
	// (1266 * 2 - 2048) / 2 = 242
	if script, err := generateScripts(ctx, modifyMigu1mTmpl, h{
		"FontForge": fontforge,
		"Ascent":    1638,
		"Descent":   410,
		"Padding":   242,
		"Inputs":    inputs,
		"Outputs":   outputs,
	}); err != nil {
		return fmt.Errorf("error in generateScripts: %v", err)
	} else if err := execScripts(ctx, script); err != nil {
		return fmt.Errorf("error in execScripts: %v", err)
	}
	return nil
}

func generateSFMonoMod(ctx context.Context, fontforge string) error {
	c := cliContext(ctx)
	var hankakus, zenkakus string
	if c.BoolT("bold") && c.BoolT("italic") {
		hankakus = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
			SFMonoTTFs[0], SFMonoTTFs[1], SFMonoTTFs[2], SFMonoTTFs[3])
		zenkakus = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
			modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[1], modifiedMigu1mTTFs[2], modifiedMigu1mTTFs[3])
	} else if c.BoolT("bold") {
		hankakus = fmt.Sprintf(`"%s", "%s"`, SFMonoTTFs[0], SFMonoTTFs[1])
		zenkakus = fmt.Sprintf(`"%s", "%s"`, modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[1])
	} else if c.BoolT("italic") {
		hankakus = fmt.Sprintf(`"%s", "%s"`, SFMonoTTFs[0], SFMonoTTFs[2])
		zenkakus = fmt.Sprintf(`"%s", "%s"`, modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[2])
	} else {
		hankakus = fmt.Sprintf(`"%s"`, SFMonoTTFs[0])
		zenkakus = fmt.Sprintf(`"%s"`, modifiedMigu1mTTFs[0])
	}
	if script, err := generateScripts(ctx, generateSFMonoModTmpl, h{
		"FontForge":         fontforge,
		"Hankakus":          hankakus,
		"Zenkakus":          zenkakus,
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
	c := cliContext(ctx)
	if !c.BoolT("nerd-fonts") {
		return nil
	}
	var SFMonos []string
	if c.BoolT("bold") && c.BoolT("italic") {
		SFMonos = SFMonoTTFs
	} else if c.BoolT("bold") {
		SFMonos = []string{SFMonoTTFs[0], SFMonoTTFs[1]}
	} else if c.BoolT("italic") {
		SFMonos = []string{SFMonoTTFs[0], SFMonoTTFs[2]}
	} else {
		SFMonos = []string{SFMonoTTFs[0]}
	}
	fp := pathInTempDir(ctx, fontPatcher)
	if err := os.Chmod(fp, 0755); err != nil {
		return fmt.Errorf("error in Chmod: %v", err)
	}
	for _, f := range SFMonos {
		mod := modified(f, c.String("suffix"))
		if err := execScripts(ctx, fp, "-c", "-q", "-out", "build", mod); err != nil {
			return fmt.Errorf("error in execScripts: %v", err)
		}
	}
	return nil
}

func modified(orig, suffix string) string {
	return strings.Replace(orig, familyName, familyName+suffix, 1)
}
