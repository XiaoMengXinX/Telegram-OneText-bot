package utils

import (
	_ "embed"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unsafe"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

var symbolsReg = regexp.MustCompile("^[a-zA-Z|{P}| ]$")

func setFontFace(gc *gg.Context, f *truetype.Font, point int) {
	gc.SetFontFace(truetype.NewFace(f, &truetype.Options{
		Size: float64(point),
	}))
	v := reflect.ValueOf(gc).Elem().FieldByName("fontHeight")
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(reflect.ValueOf(float64(point * 72 / 96)))
}

func strWrapper(dc *gg.Context, str string, maxTextWidth float64) (warpStr string) {
	if str == "" {
		return ""
	}
	warpStr = walkStrSlice(dc, splitWords(str), maxTextWidth)
	symbols := "？！，。、；：”’）》〉】』」〕…—～﹏" + `]})>!?:;,.~\|/`
	for _, r := range symbols {
		warpStr = strings.Replace(warpStr, fmt.Sprintf("\n%s", string(r)), fmt.Sprintf("%s\n", string(r)), -1)
	}
	warpStr = strings.ReplaceAll(warpStr, "\n\n", "\n")
	if warpStr[len(warpStr)-1] == '\n' {
		warpStr = warpStr[:len(warpStr)-1]
	}
	return
}

func walkStrSlice(dc *gg.Context, s []string, maxTextWidth float64) string {
	var result string
	for i := 0; i < len(s); {
		tmp := truncateText(dc, s, i, maxTextWidth)
		if tmp != nil {
			result = result + strings.Join(tmp, "") + "\n"
			i = i + len(tmp)
		} else {
			i++
		}
	}
	return result
}

func truncateText(dc *gg.Context, textSlice []string, count int, maxTextWidth float64) []string {
	tmpStr := ""
	var result []string
	for _, r := range textSlice[count:] {
		if r == "\n" {
			break
		}
		tmpStr = tmpStr + r
		w, _ := dc.MeasureString(tmpStr)
		if w > maxTextWidth {
			if len(tmpStr) <= 1 {
				return nil
			} else {
				break
			}
		} else {
			result = append(result, r)
		}
	}
	return result
}

func splitWords(str string) []string {
	var result []string
	var tmpStr string
	for _, r := range str {
		if !symbolsReg.MatchString(string(r)) {
			if tmpStr != "" {
				result = append(result, tmpStr)
				tmpStr = ""
			}
			result = append(result, string(r))
		} else {
			if unicode.IsSpace(r) {
				result = append(result, tmpStr+string(r))
				tmpStr = ""
			} else {
				tmpStr = tmpStr + string(r)
			}
		}
	}
	if tmpStr != "" {
		result = append(result, tmpStr)
	}
	return result
}
