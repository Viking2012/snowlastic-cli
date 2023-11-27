package _import

import (
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"reflect"
	"strconv"
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
func SegmentIsGiven(segment any, givenSegments []string) bool {
	switch segment.(type) {
	case bool:
		for i := range givenSegments {
			seg, err := strconv.ParseBool(givenSegments[i])
			if err == nil && reflect.ValueOf(segment).Bool() == seg {
				return true
			}
		}
	case string:
		for i := range givenSegments {
			seg := givenSegments[i]
			if segment == seg {
				return true
			}
		}
	case int, int8, int16, int32, int64:
		for i := range givenSegments {
			seg, err := strconv.ParseInt(givenSegments[i], 10, 64)
			if err == nil && reflect.ValueOf(segment).Int() == seg {
				return true
			}
		}
	case uint, uint8, uint16, uint32, uint64, uintptr:
		for i := range givenSegments {
			seg, err := strconv.ParseUint(givenSegments[i], 10, 64)
			if err == nil && reflect.ValueOf(segment).Uint() == seg {
				return true
			}
		}
	case float32, float64:
		for i := range givenSegments {
			seg, err := strconv.ParseFloat(givenSegments[i], 64)
			if err == nil && reflect.ValueOf(segment).Float() == seg {
				return true
			}
		}
	case complex64, complex128:
		for i := range givenSegments {
			seg, err := strconv.ParseComplex(givenSegments[i], 128)
			if err == nil && reflect.ValueOf(segment).Complex() == seg {
				return true
			}
		}
	}
	return false
}
