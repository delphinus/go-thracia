package thracia

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
	"gopkg.in/cheggaaa/pb.v1"
)

type toDownload struct {
	filename, URL string
	pb            chan *pb.ProgressBar
}

func download(ctx context.Context, toDL []*toDownload) error {
	if len(toDL) == 0 {
		return nil
	}
	eg, ctx := errgroup.WithContext(ctx)
	filenameLen := 0
	for _, dl := range toDL {
		dl := dl
		if len(dl.filename) > filenameLen {
			filenameLen = len(dl.filename)
		}
		if dl.pb == nil {
			dl.pb = make(chan *pb.ProgressBar)
		}
		eg.Go(func() error { return fetch(ctx, dl) })
	}

	eg.Go(func() error { return showBar(ctx, toDL, filenameLen) })

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error in Wait: %v", err)
	}
	return nil
}

func showBar(ctx context.Context, toDL []*toDownload, filenameLen int) error {
	pbs := make([]*pb.ProgressBar, len(toDL))
	format := fmt.Sprintf("%%%ds", filenameLen)
	for i, dl := range toDL {
		select {
		case <-ctx.Done():
			return nil
		case b := <-dl.pb:
			prefix := fmt.Sprintf(format, dl.filename)
			pbs[i] = b.Prefix(prefix)
		}
	}
	if _, err := pb.StartPool(pbs...); err != nil {
		return fmt.Errorf("error in StartPool: %v", err)
	}
	return nil
}

func fetch(ctx context.Context, dl *toDownload) (err error) {
	req, err := http.NewRequest("GET", dl.URL, nil)
	if err != nil {
		return fmt.Errorf("error in NewRequest: %v", err)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error in Do: %v", err)
	}
	defer checkClose(resp.Body, &err)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("invalid status: %d %s", resp.StatusCode, resp.Status)
	}

	total, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	b := pb.New64(total).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10).SetWidth(80)
	b.ShowSpeed = true
	dl.pb <- b

	file, err := os.Create(pathInTempDir(ctx, dl.filename))
	if err != nil {
		return fmt.Errorf("error in Create: %v", err)
	}
	defer checkClose(file, &err)

	go func() {
		<-ctx.Done()
		b.Finish()
	}()
	_, err = io.Copy(file, b.NewProxyReader(resp.Body))
	return err
}
