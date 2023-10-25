package _import

import (
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func AddBarPersistent(bars *mpb.Progress, size int64, prefix string, suffix any) *mpb.Bar {
	return bars.AddBar(size,
		//mpb.BarRemoveOnComplete(),
		//mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("%s: %v ", prefix, suffix), decor.WCSyncSpaceR),
			decor.CountersNoUnit("%5d/%5d ", decor.WCSyncWidth),
			decor.Percentage(decor.WC{W: 5}),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), " done",
			),
		),
	)
}
func AddBarMomentary(bars *mpb.Progress, size int64, prefix string, suffix any) *mpb.Bar {
	return bars.AddBar(size,
		mpb.BarRemoveOnComplete(),
		//mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("%s: %v ", prefix, suffix), decor.WCSyncSpaceR),
			decor.CountersNoUnit("%5d/%5d ", decor.WCSyncWidth),
			decor.Percentage(decor.WC{W: 5}),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), " done",
			),
		),
	)
}
