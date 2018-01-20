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
	if script, err := generateScripts(ctx, modifyMigu1mTmpl, h{
		"FontForge":  fontforge,
		"SrcRegular": migu1mTTFs[0],
		"SrcBold":    migu1mTTFs[1],
		"DstRegular": modifiedMigu1mTTFs[0],
		"DstBold":    modifiedMigu1mTTFs[1],
		"ScaleDown":  cliContext(ctx).BoolT("scale-down"),
		"DstDir":     tempDir(ctx),
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

func execScripts(ctx context.Context, script string) error {
	c := cliContext(ctx)
	cmd := exec.CommandContext(ctx, script)
	cmd.Dir = tempDir(ctx)
	cmd.Stdout = c.App.Writer
	cmd.Stderr = c.App.ErrWriter
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error in Run: %v", err)
	}
	return nil
}
