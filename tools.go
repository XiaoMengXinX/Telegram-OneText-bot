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

func drawString(dc *gg.Context, s string, x, y, width, lineSpacing float64, align gg.Align) {
	lines := strings.Split(s, "\n")
	var ax, ay float64

	h := float64(len(lines)) * dc.FontHeight() * lineSpacing
	h -= (lineSpacing - 1) * dc.FontHeight()

	switch align {
	case gg.AlignLeft:
		ax = 0
	case gg.AlignCenter:
		ax = 0.5
		x += width / 2
	case gg.AlignRight:
		ax = 1
		x += width
	}
	ay = 1
	for _, line := range lines {
		dc.DrawStringAnchored(line, x, y, ax, ay)
		y += dc.FontHeight() * lineSpacing
	}
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
