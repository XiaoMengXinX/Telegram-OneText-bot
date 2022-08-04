package utils

import (
	_ "embed"
	"reflect"
	"strings"
	"unsafe"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

func setFontFace(gc *gg.Context, f *truetype.Font, point int) {
	gc.SetFontFace(truetype.NewFace(f, &truetype.Options{
		Size: float64(point),
	}))
	v := reflect.ValueOf(gc).Elem().FieldByName("fontHeight")
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(reflect.ValueOf(float64(point * 72 / 96)))
}

func strWrapper(dc *gg.Context, str string, maxTextWidth float64) (warpStr string) {
	for i := 0; i < len(str); {
		if str[i] == '\n' {
			i++
		}
		str := truncateText(dc, str[i:], maxTextWidth)
		i = i + len(str)
		warpStr = warpStr + str + "\n"
	}
	punctuationReplacer := strings.NewReplacer("\n。", "。\n", "\n，", "，\n", "\n！", "！\n", "\n？", "？\n", "\n\"", "\"\n", "\n“", "“\n", "\n”", "”\n")
	warpStr = strings.ReplaceAll(punctuationReplacer.Replace(warpStr), "\n\n", "\n")
	if warpStr[len(warpStr)-1] == '\n' {
		warpStr = warpStr[:len(warpStr)-1]
	}
	return
}

func truncateText(dc *gg.Context, originalText string, maxTextWidth float64) string {
	tmpStr := ""
	result := make([]rune, 0)
	for _, r := range originalText {
		if r == '\n' {
			break
		}
		tmpStr = tmpStr + string(r)
		w, _ := dc.MeasureString(tmpStr)
		if w > maxTextWidth {
			if len(tmpStr) <= 1 {
				return ""
			} else {
				break
			}
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
