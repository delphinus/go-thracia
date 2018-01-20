package thracia

import (
	"context"
	"fmt"
	"os"
)

const (
	migu1mURL  = "https://osdn.jp/frs/redir.php?m=gigenet&f=%2Fmix-mplus-ipa%2F63545%2Fmigu-1m-20150712.zip"
	migu1mFile = "migu1m.zip"
	migu1mDir  = "migu-1m-20150712"
	// migu1mTTFs = []string{"migu-1m-regular.ttf", "migu-1m-bold.ttf"}
)

// Run is the main routine
func Run(ctx context.Context) error {
	toDL := []*toDownload{}
	if _, err := os.Stat(migu1mFile); os.IsNotExist(err) {
		toDL = append(toDL, &toDownload{
			filename: migu1mFile,
			URL:      migu1mURL,
		})
	} else if err != nil {
		return fmt.Errorf("error in Stat: %v", err)
	}
	if err := download(ctx, toDL); err != nil {
		return fmt.Errorf("error in download: %v", err)
	}
	return nil
}
