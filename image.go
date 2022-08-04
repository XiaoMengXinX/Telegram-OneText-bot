package utils

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"

	onetext "github.com/XiaoMengXinX/OneTextAPI-Go"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

func CreateOnetextImage(s onetext.Sentence) ([]byte, error) {
	weight := 1080
	height := 0

	text := s.Text
	by := s.By
	from := s.From
	recordTime := ""
	if len(s.Time) > 0 {
		recordTime = s.Time[0]
	}
	createTime := ""
	if len(s.Time) > 1 {
		createTime = s.Time[1]
	}

	f, err := truetype.Parse(fontFile)
	if err != nil {
		return nil, err
	}

	textContent := gg.NewContext(1080, 4000)
	textContent.SetHexColor("#FFFFFF")
	setFontFace(textContent, f, 59)
	textContent.SetHexColor("#000000")

	warpStr := strWrapper(textContent, text, 780)

	_, oneLineHeight := textContent.MeasureString("字")
	newLineCount := float64(strings.Count(warpStr, "\n"))
	imgTextHeight := (newLineCount + 1) * (oneLineHeight * 1.8)
	textContent.DrawStringWrapped(strings.ReplaceAll(warpStr, " ", "\b"), 0, 0, 0, 0, 59, 1.8, gg.AlignLeft)
	height = int(imgTextHeight + oneLineHeight*1.8 + 220)

	byContent := gg.NewContext(weight, 200)
	byContent.SetHexColor("#FFFFFF")
	if by != "" {
		height = height + 70
		setFontFace(byContent, f, 48)
		byContent.SetHexColor("#313131")
		byContent.DrawStringWrapped(fmt.Sprintf("——\b%s", strings.ReplaceAll(by, " ", "\b")), 930, 0, 0, 0, 48, 1.8, gg.AlignRight)
	}

	timeContent := gg.NewContext(weight, 200)
	timeContent.SetHexColor("#FFFFFF")
	if recordTime != "" {
		height = height + 110
		setFontFace(timeContent, f, 40)
		timeContent.SetHexColor("#313131")
		timeStr := ""
		if createTime != "" {
			timeStr = fmt.Sprintf("记录于：%s\b创作于：%s", recordTime, createTime)
		} else {
			timeStr = fmt.Sprintf("记录于：%s", recordTime)
		}
		timeContent.DrawStringWrapped(timeStr, 935, 0, 0, 0, 40, 1.8, gg.AlignRight)
	}

	fromContent := gg.NewContext(weight, 200)
	fromContent.SetHexColor("#FFFFFF")
	if from != "" {
		setFontFace(fromContent, f, 38)
		fromContent.SetHexColor("#313131")
		fromStr := strWrapper(fromContent, from, 860)
		_, fromOnelineHeight := fromContent.MeasureString("字")
		height = height + strings.Count(fromStr, "\n")*int(fromOnelineHeight*1.8) + 110
		fromContent.DrawStringWrapped(strings.ReplaceAll(fromStr, " ", "\b"), 0, 0, 0, 0, 38, 1.8, gg.AlignLeft)
	}

	height = height + 150

	fw := gg.NewContext(weight, height)
	fw.SetHexColor("#FFFFFF")
	fw.Clear()
	fw.DrawRoundedRectangle(55, 55, float64(weight-55*2), float64(height-55*2), 10)
	fw.SetLineWidth(4)
	fw.SetHexColor("#e3e3e3")
	fw.StrokePreserve()
	setFontFace(fw, f, 55)
	fw.SetRGB(0, 0, 0)
	fw.DrawString("“", 110, 165)
	fw.DrawImage(textContent.Image(), 160, 220)
	lastHeight := imgTextHeight + oneLineHeight*1.8 + 220
	fw.DrawString("”", 940, lastHeight)
	if by != "" {
		fw.DrawImage(byContent.Image(), 0, int(lastHeight+70))
		lastHeight = lastHeight + 70
	}
	if recordTime != "" {
		fw.DrawImage(timeContent.Image(), 0, int(lastHeight+110))
		lastHeight = lastHeight + 110
	}
	if from != "" {
		fw.DrawImage(fromContent.Image(), 110, int(lastHeight+110))
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, fw.Image())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
