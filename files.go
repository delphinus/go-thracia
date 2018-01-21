package thracia

const (
	migu1mFile = "migu1m.zip"
	migu1mURL  = "https://osdn.jp/frs/redir.php?m=gigenet&f=%2Fmix-mplus-ipa%2F63545%2Fmigu-1m-20150712.zip"
)

func files() []*toDownload {
	return []*toDownload{
		{
			filename: migu1mFile,
			URL:      migu1mURL,
		},
		{
			filename: "devicons.ttf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/devicons.ttf",
		},
		{
			filename: "font-awesome-extension.ttf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/font-awesome-extension.ttf",
		},
		{
			filename: "font-linux.ttf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/font-linux.ttf",
		},
		{
			filename: "FontAwesome.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/FontAwesome.otf",
		},
		{
			filename: "octicons.ttf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/octicons.ttf",
		},
		{
			filename: "original-source.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/original-source.otf",
		},
		{
			filename: "Pomicons.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/Pomicons.otf",
		},
		{
			filename: "PowerlineExtraSymbols.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/PowerlineExtraSymbols.otf",
		},
		{
			filename: "PowerlineSymbols.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/PowerlineSymbols.otf",
		},
		{
			filename: "Unicode_IEC_symbol_font.otf",
			URL:      "https://github.com/delphinus/nerd-fonts-simple/raw/master/src/glyphs/Unicode_IEC_symbol_font.otf",
		},
	}
}
