package thracia

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const (
	migu1mFile = "migu1m.zip"
	migu1mURL  = "https://osdn.jp/frs/redir.php?m=gigenet&f=%2Fmix-mplus-ipa%2F63545%2Fmigu-1m-20150712.zip"
	glyphsDir  = "src/glyphs/"
)

func files(ctx context.Context) ([]*toDownload, error) {
	dir := filepath.Join(tempDir(ctx), glyphsDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error in MkdirAll: %v", err)
	}
	return []*toDownload{
		{
			filename: migu1mFile,
			URL:      migu1mURL,
		},
		{
			filename: glyphsDir + "devicons.ttf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/devicons.ttf",
		},
		{
			filename: glyphsDir + "font-awesome-extension.ttf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/font-awesome-extension.ttf",
		},
		{
			filename: glyphsDir + "font-linux.ttf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/font-linux.ttf",
		},
		{
			filename: glyphsDir + "FontAwesome.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/FontAwesome.otf",
		},
		{
			filename: glyphsDir + "octicons.ttf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/octicons.ttf",
		},
		{
			filename: glyphsDir + "original-source.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/original-source.otf",
		},
		{
			filename: glyphsDir + "Pomicons.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/Pomicons.otf",
		},
		{
			filename: glyphsDir + "PowerlineExtraSymbols.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/PowerlineExtraSymbols.otf",
		},
		{
			filename: glyphsDir + "PowerlineSymbols.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/PowerlineSymbols.otf",
		},
		{
			filename: glyphsDir + "Unicode_IEC_symbol_font.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/Unicode_IEC_symbol_font.otf",
		},
	}, nil
}
