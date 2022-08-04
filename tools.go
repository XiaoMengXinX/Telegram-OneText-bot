package utils

import (
	_ "embed"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unsafe"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

var symbolsReg = "^[\u4e00-\u9fa5|。|？|！|，|、|；|：|“|”|‘|’|（|）|《|》|〈|〉|【|】|『|』|「|」|﹃|﹄|〔|〕|…|—|～|﹏|￥]$"

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
	warpStr = walkStrSlice(dc, splitHanAndASCII(str), maxTextWidth)
	punctuationReplacer := strings.NewReplacer("\n。", "。\n", "\n，", "，\n", "\n！", "！\n", "\n？", "？\n", "\n\"", "\"\n", "\n“", "“\n", "\n”", "”\n")
	warpStr = strings.ReplaceAll(punctuationReplacer.Replace(warpStr), "\n\n", "\n")
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
		if regexp.MustCompile(symbolsReg).MatchString(r) {
			tmpStr = tmpStr + r
		} else {
			tmpStr = tmpStr + " " + r
		}
		w, _ := dc.MeasureString(tmpStr)
		if w > maxTextWidth {
			if len(tmpStr) <= 1 {
				return nil
			} else {
				break
			}
		} else {
			if regexp.MustCompile(symbolsReg).MatchString(r) {
				result = append(result, r)
			} else {
				result = append(result, r+" ")
			}
		}
	}
	return result
}

func splitHanAndASCII(str string) []string {
	var result []string
	var tmpStr string
	for _, r := range str {
		if regexp.MustCompile(symbolsReg).MatchString(string(r)) {
			if tmpStr != "" {
				result = append(result, tmpStr)
				tmpStr = ""
			}
			result = append(result, string(r))
		} else {
			if r == '\n' {
				result = append(result, tmpStr, "\n")
				tmpStr = ""
				continue
			}
			if unicode.IsSpace(r) {
				result = append(result, tmpStr)
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
