package thracia

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	ascent  = 1638
	descent = 410
)

func scripts(ctx context.Context) error {
	c := cliContext(ctx)
	fontforge, err := exec.LookPath("fontforge")
	if err != nil {
		return errors.New("cannot find `fontforge` executable")
	}
	if err := generateOblique(ctx, fontforge); err != nil {
		return fmt.Errorf("error in generateOblique: %v", err)
	}
	if err := modifyMigu1m(ctx, fontforge); err != nil {
		return fmt.Errorf("error in modifyMigu1m: %v", err)
	}
	if c.Bool("square") {
		if err := modifySFMono(ctx, fontforge); err != nil {
			return fmt.Errorf("error in modifySFMono: %v", err)
		}
	}
	if err := generateSFMonoMod(ctx, fontforge); err != nil {
		return fmt.Errorf("error in generateSFMonoMod: %v", err)
	}
	python2, err := exec.LookPath("python2")
	if err != nil {
		return errors.New("cannot find `python2` executable")
	}
	if err := execFontPatcher(ctx, python2); err != nil {
		return fmt.Errorf("error in execFontPatcher: %v", err)
	}
	if err := copyBuildDir(ctx); err != nil {
		return fmt.Errorf("error in copyBuildDir: %v", err)
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
	// Migu1M has this:
	// 860 + 140 = 1000
	var square string
	var padding, scale float64
	var width, hankakuWidth, halfHankakuWidth, hankakuPadding int
	if c.Bool("square") {
		// So zenkaku glyphs should move:
		// 140 - 2048 * (140 / 1000) = -146.72
		square = "true"
		padding = -146.72
		scale = 82.0
		width = 2048
		hankakuWidth = 1024
		halfHankakuWidth = 512
	} else {
		// So zenkaku glyphs should have padding on left and right:
		// (1266 * 2 - 2048) / 2 = 242
		padding = 242.0
		hankakuWidth = 1266
		halfHankakuWidth = 1266
		hankakuPadding = int(padding) / 2
	}
	if script, err := generateScripts(ctx, modifyMigu1mTmpl, h{
		"FontForge":        fontforge,
		"Square":           square,
		"Ascent":           ascent,
		"Descent":          descent,
		"Padding":          padding,
		"Scale":            scale,
		"Width":            width,
		"HankakuWidth":     hankakuWidth,
		"HalfHankakuWidth": halfHankakuWidth,
		"HankakuPadding":   hankakuPadding,
		"Inputs":           inputs,
		"Outputs":          outputs,
	}); err != nil {
		return fmt.Errorf("error in generateScripts: %v", err)
	} else if err := execScripts(ctx, script); err != nil {
		return fmt.Errorf("error in execScripts: %v", err)
	}
	return nil
}

func modifySFMono(ctx context.Context, fontforge string) error {
	c := cliContext(ctx)
	var inputs, outputs string
	if c.BoolT("bold") && c.BoolT("italic") {
		inputs = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
			SFMonoTTFs[0], SFMonoTTFs[1], SFMonoTTFs[2], SFMonoTTFs[3])
		outputs = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
			modifiedSFMonoTTFs[0], modifiedSFMonoTTFs[1], modifiedSFMonoTTFs[2], modifiedSFMonoTTFs[3])
	} else if c.BoolT("bold") {
		inputs = fmt.Sprintf(`"%s", "%s"`, SFMonoTTFs[0], SFMonoTTFs[1])
		outputs = fmt.Sprintf(`"%s", "%s"`, modifiedSFMonoTTFs[0], modifiedSFMonoTTFs[1])
	} else if c.BoolT("italic") {
		inputs = fmt.Sprintf(`"%s", "%s"`, SFMonoTTFs[0], SFMonoTTFs[2])
		outputs = fmt.Sprintf(`"%s", "%s"`, modifiedSFMonoTTFs[0], modifiedSFMonoTTFs[2])
	} else {
		inputs = fmt.Sprintf(`"%s"`, SFMonoTTFs[0])
		outputs = fmt.Sprintf(`"%s"`, modifiedSFMonoTTFs[0])
	}
	// New parameters
	// Scale: 1024 / 1266 * 100 = 80.8847
	if script, err := generateScripts(ctx, modifySFMonoTmpl, h{
		"FontForge": fontforge,
		"Scale":     80.8847,
		"CenterX":   0,
		"CenterY":   0,
		"Width":     1024,
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
		if c.Bool("square") {
			hankakus = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
				modifiedSFMonoTTFs[0], modifiedSFMonoTTFs[1], modifiedSFMonoTTFs[2], modifiedSFMonoTTFs[3])
		} else {
			hankakus = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
				SFMonoTTFs[0], SFMonoTTFs[1], SFMonoTTFs[2], SFMonoTTFs[3])
		}
		zenkakus = fmt.Sprintf(`"%s", "%s", "%s", "%s"`,
			modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[1], modifiedMigu1mTTFs[2], modifiedMigu1mTTFs[3])
	} else if c.BoolT("bold") {
		if c.Bool("square") {
			hankakus = fmt.Sprintf(`"%s", "%s"`, modifiedSFMonoTTFs[0], modifiedSFMonoTTFs[1])
		} else {
			hankakus = fmt.Sprintf(`"%s", "%s"`, SFMonoTTFs[0], SFMonoTTFs[1])
		}
		zenkakus = fmt.Sprintf(`"%s", "%s"`, modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[1])
	} else if c.BoolT("italic") {
		if c.Bool("square") {
			hankakus = fmt.Sprintf(`"%s", "%s"`, modifiedSFMonoTTFs[0], modifiedSFMonoTTFs[2])
		} else {
			hankakus = fmt.Sprintf(`"%s", "%s"`, SFMonoTTFs[0], SFMonoTTFs[2])
		}
		zenkakus = fmt.Sprintf(`"%s", "%s"`, modifiedMigu1mTTFs[0], modifiedMigu1mTTFs[2])
	} else {
		if c.Bool("square") {
			hankakus = fmt.Sprintf(`"%s"`, modifiedSFMonoTTFs[0])
		} else {
			hankakus = fmt.Sprintf(`"%s"`, SFMonoTTFs[0])
		}
		zenkakus = fmt.Sprintf(`"%s"`, modifiedMigu1mTTFs[0])
	}
	var square string
	var winAscent, winDescent, hankakuWidth, zenkakuWidth, padding int
	if c.Bool("square") {
		winAscent = ascent
		winDescent = descent
		hankakuWidth = 1024
		zenkakuWidth = 1024 * 2
		padding = 1024 / 2
	} else {
		winAscent = 1950
		winDescent = 494
		hankakuWidth = 1266
		zenkakuWidth = 1266 * 2
		padding = 1266 / 2
	}
	if script, err := generateScripts(ctx, generateSFMonoModTmpl, h{
		"FontForge":         fontforge,
		"Square":            square,
		"Hankakus":          hankakus,
		"Zenkakus":          zenkakus,
		"FamilyName":        familyName,
		"FamilyNameSuffix":  c.String("suffix"),
		"Version":           version,
		"WinAscent":         winAscent,
		"WinDescent":        winDescent,
		"Ascent":            ascent,
		"Descent":           descent,
		"ZenkakuSpaceGlyph": "",
		"HankakuWidth":      hankakuWidth,
		"ZenkakuWidth":      zenkakuWidth,
		"Padding":           padding,
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

func execFontPatcher(ctx context.Context, python2 string) error {
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

	fp, err := generateScripts(ctx, fontPatcherTmpl, h{
		"Python2": python2,
	})
	if err != nil {
		return fmt.Errorf("error in generateScripts: %v", err)
	}
	if _, err := generateScripts(ctx, changelogTmpl, nil); err != nil {
		return fmt.Errorf("error in generateScripts: %v", err)
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

func copyBuildDir(ctx context.Context) error {
	buildDir := pathInTempDir(ctx, "build")
	if st, err := os.Stat(buildDir); os.IsNotExist(err) || !st.IsDir() {
		return nil
	}
	files, err := ioutil.ReadDir(buildDir)
	if err != nil {
		return fmt.Errorf("error in ReadDir: %v", err)
	}
	if err := os.MkdirAll("build", 0755); err != nil {
		return fmt.Errorf("error in MkdirAll: %v", err)
	}
	for _, f := range files {
		in, err := os.Open(filepath.Join(buildDir, f.Name()))
		if err != nil {
			return fmt.Errorf("error in Open: %s: %v", f.Name(), err)
		}
		outFile := "build/" + f.Name()
		out, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("error in Create: %s: %v", outFile, err)
		}
		if _, err := io.Copy(out, in); err != nil {
			return fmt.Errorf("error in Copy: %v", err)
		}
		if err := in.Close(); err != nil {
			return fmt.Errorf("error in Close: %v", err)
		}
		if err := out.Close(); err != nil {
			return fmt.Errorf("error in Close: %v", err)
		}
		fmt.Printf("copied: %s\n", outFile)
	}
	return nil
}
