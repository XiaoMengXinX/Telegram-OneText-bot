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

func CreateOnetextImage(s onetext.Sentence, font FontConfig) ([]byte, error) {
	weight := 1080
	height := 0

	// default font size is for canger.ttf
	var textFontSize = int(59 * font.FontScale)
	var byFontSize = int(48 * font.FontScale)
	var fromFontSize = int(38 * font.FontScale)
	var timeFontSize = int(40 * font.FontScale)

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

	f, err := truetype.Parse(font.FontFile)
	if err != nil {
		return nil, err
	}

	textContent := gg.NewContext(1080, 3000)
	textContent.SetHexColor("#FFFFFF")
	setFontFace(textContent, f, textFontSize)
	textContent.SetHexColor("#000000")

	warpStr := strWrapper(textContent, text, 780)

	_, oneLineHeight := textContent.MeasureString("字")
	newLineCount := float64(strings.Count(warpStr, "\n"))
	imgTextHeight := (newLineCount + 1) * (oneLineHeight * 1.8)
	drawString(textContent, warpStr, 0, 20, float64(textFontSize), 1.8, gg.AlignLeft)
	height = int(imgTextHeight + oneLineHeight*1.8 + 165*font.FontScale + 55)

	byContent := gg.NewContext(weight, 200)
	byContent.SetHexColor("#FFFFFF")
	var byHeight float64
	if by != "" {
		height = height + int(70*font.FontScale)
		setFontFace(byContent, f, byFontSize)
		byContent.SetHexColor("#313131")
		byStr := strWrapper(byContent, fmt.Sprintf("—— %s", by), 860)
		_, byOnelineHeight := byContent.MeasureString("字")
		byHeight = float64(strings.Count(byStr, "\n"))*byOnelineHeight*1.8 + 70*font.FontScale
		height = height + int(byHeight)
		drawString(byContent, byStr, 930, 10, float64(byFontSize), 1.8, gg.AlignRight)
	}

	timeContent := gg.NewContext(weight, 200)
	timeContent.SetHexColor("#FFFFFF")
	if recordTime != "" {
		height = height + int(110*font.FontScale)
		setFontFace(timeContent, f, timeFontSize)
		timeContent.SetHexColor("#313131")
		timeStr := ""
		if createTime != "" {
			timeStr = fmt.Sprintf("记录于：%s 创作于：%s", recordTime, createTime)
		} else {
			timeStr = fmt.Sprintf("记录于：%s", recordTime)
		}
		drawString(timeContent, timeStr, 935, 10, float64(timeFontSize), 1.8, gg.AlignRight)
	}

	fromContent := gg.NewContext(weight, 200)
	fromContent.SetHexColor("#FFFFFF")
	if from != "" {
		setFontFace(fromContent, f, fromFontSize)
		fromContent.SetHexColor("#313131")
		fromStr := strWrapper(fromContent, from, 860)
		_, fromOnelineHeight := fromContent.MeasureString("字")
		height = height + int(float64(strings.Count(fromStr, "\n"))*fromOnelineHeight*1.8) + int(25*font.FontScale)
		drawString(fromContent, fromStr, 0, 10, float64(fromFontSize), 1.8, gg.AlignLeft)
	}

	height = height + 55 + int(110*font.FontScale)

	fw := gg.NewContext(weight, height)
	fw.SetHexColor("#FFFFFF")
	fw.Clear()
	fw.DrawRoundedRectangle(55, 55, float64(weight-55*2), float64(height-55*2), 10)
	fw.SetLineWidth(4)
	fw.SetHexColor("#e3e3e3")
	fw.StrokePreserve()
	setFontFace(fw, f, textFontSize)
	fw.SetRGB(0, 0, 0)
	lastY := 55 + 110*font.FontScale
	fw.DrawString("“", 110, lastY)
	fw.DrawImage(textContent.Image(), 160, int(lastY+55*font.FontScale)-20)
	lastY = imgTextHeight + oneLineHeight*1.8 + lastY + 55*font.FontScale
	fw.DrawString("”", 940, lastY)
	if by != "" {
		lastY = lastY + 70*font.FontScale
		fw.DrawImage(byContent.Image(), 0, int(lastY)-10)
		lastY = lastY + byHeight
	}
	if recordTime != "" {
		lastY = lastY + 40*font.FontScale
		fw.DrawImage(timeContent.Image(), 0, int(lastY)-10)
		lastY = lastY + 70*font.FontScale
	}
	if from != "" {
		lastY = lastY + 40*font.FontScale
		fw.DrawImage(fromContent.Image(), 110, int(lastY)-10)
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, fw.Image())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
