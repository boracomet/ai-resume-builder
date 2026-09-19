package pdf

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	allocOnce   sync.Once
	allocCtx    context.Context
	allocCancel context.CancelFunc
	allocErr    error
	pdfMu       sync.Mutex
)

func browserAllocator() (context.Context, error) {
	allocOnce.Do(func() {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("hide-scrollbars", true),
		)
		if chromePath := strings.TrimSpace(os.Getenv("CHROME_PATH")); chromePath != "" {
			opts = append(opts, chromedp.ExecPath(chromePath))
		}

		allocCtx, allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	})
	return allocCtx, allocErr
}

func Generate(html string) ([]byte, error) {
	pdfMu.Lock()
	defer pdfMu.Unlock()

	parent, err := browserAllocator()
	if err != nil {
		return nil, err
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		_, _ = w.Write([]byte(html))
	}))
	defer server.Close()

	ctx, cancelCtx := chromedp.NewContext(parent)
	defer cancelCtx()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	url := fmt.Sprintf("%s/?t=%d", server.URL, time.Now().UnixNano())

	var pdfBuf []byte
	err = chromedp.Run(ctx,
		network.Enable(),
		network.SetCacheDisabled(true),
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(250*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	)
	if err != nil {
		// Reset shared browser so the next request starts clean.
		resetBrowserLocked()
		return nil, err
	}
	return pdfBuf, nil
}

func resetBrowserLocked() {
	if allocCancel != nil {
		allocCancel()
	}
	allocOnce = sync.Once{}
	allocCtx = nil
	allocCancel = nil
	allocErr = nil
}
