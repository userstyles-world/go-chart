package chart

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
	"github.com/userstyles-world/go-chart/v2/drawing"
	"github.com/valyala/bytebufferpool"
)

// SVG returns a new png/raster renderer.
func SVG(width, height int) (Renderer, error) {
	buffer := bytebufferpool.Get()
	buffer2 := bytebufferpool.Get()

	canvas := newCanvas(buffer)
	canvas.Start(width, height)
	return &vectorRenderer{
		b:   buffer,
		c:   canvas,
		s:   &Style{},
		p:   buffer2,
		dpi: DefaultDPI,
	}, nil
}

// SVGWithCSS returns a new png/raster renderer with attached custom CSS
// The optional nonce argument sets a CSP nonce.
func SVGWithCSS(css string, nonce string) func(width, height int) (Renderer, error) {
	return func(width, height int) (Renderer, error) {
		buffer := bytebufferpool.Get()
		buffer2 := bytebufferpool.Get()

		canvas := newCanvas(buffer)
		canvas.css = css
		canvas.nonce = nonce
		canvas.Start(width, height)
		return &vectorRenderer{
			b:   buffer,
			c:   canvas,
			s:   &Style{},
			p:   buffer2,
			dpi: DefaultDPI,
		}, nil
	}
}

// vectorRenderer renders chart commands to a bitmap.
type vectorRenderer struct {
	dpi float64
	b   *bytebufferpool.ByteBuffer
	c   *canvas
	s   *Style
	p   *bytebufferpool.ByteBuffer
	fc  font.Face
}

func (vr *vectorRenderer) ResetStyle() {
	vr.s = &Style{Font: vr.s.Font}
	vr.fc = nil
}

// GetDPI returns the dpi.
func (vr *vectorRenderer) GetDPI() float64 {
	return vr.dpi
}

// SetDPI implements the interface method.
func (vr *vectorRenderer) SetDPI(dpi float64) {
	vr.dpi = dpi
	vr.c.dpi = dpi
}

// SetClassName implements the interface method.
func (vr *vectorRenderer) SetClassName(classname string) {
	vr.s.ClassName = classname
}

// SetStrokeColor implements the interface method.
func (vr *vectorRenderer) SetStrokeColor(c drawing.Color) {
	vr.s.StrokeColor = c
}

// SetFillColor implements the interface method.
func (vr *vectorRenderer) SetFillColor(c drawing.Color) {
	vr.s.FillColor = c
}

// SetLineWidth implements the interface method.
func (vr *vectorRenderer) SetStrokeWidth(width float64) {
	vr.s.StrokeWidth = width
}

// StrokeDashArray sets the stroke dash array.
func (vr *vectorRenderer) SetStrokeDashArray(dashArray []float64) {
	vr.s.StrokeDashArray = dashArray
}

var (
	mStart = []byte("M ")
	space  = []byte(" ")
)

// MoveTo implements the interface method.
func (vr *vectorRenderer) MoveTo(x, y int) {
	_, _ = vr.p.Write(mStart)
	_, _ = vr.p.WriteString(itoa(x))
	_, _ = vr.p.Write(space)
	_, _ = vr.p.WriteString(itoa(y))
}

var (
	lStart = []byte("L ")
)

// LineTo implements the interface method.
func (vr *vectorRenderer) LineTo(x, y int) {
	_, _ = vr.p.Write(lStart)
	_, _ = vr.p.WriteString(itoa(x))
	_, _ = vr.p.Write(space)
	_, _ = vr.p.WriteString(itoa(y))
}

var (
	qStart = []byte("Q")
	comma  = []byte(",")
)

// QuadCurveTo draws a quad curve.
func (vr *vectorRenderer) QuadCurveTo(cx, cy, x, y int) {
	_, _ = vr.p.Write(qStart)
	_, _ = vr.p.WriteString(itoa(cx))
	_, _ = vr.p.Write(comma)
	_, _ = vr.p.WriteString(itoa(cy))
	_, _ = vr.p.Write(space)
	_, _ = vr.p.WriteString(itoa(x))
	_, _ = vr.p.Write(comma)
	_, _ = vr.p.WriteString(itoa(y))
}

var (
	zClose = []byte("Z")
)

// Close closes a shape.
func (vr *vectorRenderer) Close() {
	_, _ = vr.p.Write(zClose)
}

// Stroke draws the path with no fill.
func (vr *vectorRenderer) Stroke() {
	vr.drawPath(vr.s.GetStrokeOptions())
}

// Fill draws the path with no stroke.
func (vr *vectorRenderer) Fill() {
	vr.drawPath(vr.s.GetFillOptions())
}

// FillStroke draws the path with both fill and stroke.
func (vr *vectorRenderer) FillStroke() {
	vr.drawPath(vr.s.GetFillAndStrokeOptions())
}

// drawPath draws a path.
func (vr *vectorRenderer) drawPath(s Style) {
	vr.c.Path(vr.p.String(), s.GetFillAndStrokeOptions())
	vr.p.Reset()
}

// Circle implements the interface method.
func (vr *vectorRenderer) Circle(radius float64, x, y int) {
	vr.c.Circle(x, y, int(radius), vr.s.GetFillAndStrokeOptions())
}

// SetFont implements the interface method.
func (vr *vectorRenderer) SetFont(f *truetype.Font) {
	vr.s.Font = f
}

// SetFontColor implements the interface method.
func (vr *vectorRenderer) SetFontColor(c drawing.Color) {
	vr.s.FontColor = c
}

// SetFontSize implements the interface method.
func (vr *vectorRenderer) SetFontSize(size float64) {
	vr.s.FontSize = size
}

// Text draws a text blob.
func (vr *vectorRenderer) Text(body string, x, y int) {
	vr.c.Text(x, y, body, vr.s.GetTextOptions())
}

type runeCache map[rune]fixed.Int26_6

var (
	fontCache  = map[string]font.Face{}
	fontToRune = map[string]runeCache{}
)

// Implement a caching measureString
func measureString(str, fontName string, font font.Face) fixed.Int26_6 {
	cache := fontToRune[fontName]
	if cache == nil {
		cache = make(runeCache)
	}
	prevC := rune(-1)
	var advance fixed.Int26_6
	for _, c := range str {
		if prevC >= 0 {
			advance += font.Kern(prevC, c)
		}
		var a fixed.Int26_6
		if res, ok := cache[c]; ok {
			a = res
		} else {
			a, ok = font.GlyphAdvance(c)
			if !ok {
				// TODO: is falling back on the U+FFFD glyph the responsibility of
				// the Drawer or the Face?
				// TODO: set prevC = '\ufffd'?
				continue
			}
			cache[c] = a
		}
		advance += a
		prevC = c
	}
	fontToRune[fontName] = cache
	return advance
}

// MeasureText uses the truetype font drawer to measure the width of text.
func (vr *vectorRenderer) MeasureText(body string) (box Box) {
	if vr.s.GetFont() != nil {
		fontName := vr.s.GetFont().Name(truetype.NameIDFontFamily)
		cacheName := fontName + strconv.FormatInt(int64(vr.dpi), 10) + strconv.FormatInt(int64(vr.s.FontSize), 10)
		if vr.fc == nil {
			if f, ok := fontCache[cacheName]; ok {
				vr.fc = f
			} else {
				newFace := truetype.NewFace(vr.s.GetFont(), &truetype.Options{
					DPI:  vr.dpi,
					Size: vr.s.FontSize,
				})
				vr.fc = newFace
				fontCache[cacheName] = newFace
			}

		}
		w := measureString(body, fontName, vr.fc).Ceil()

		box.Right = w
		box.Bottom = int(drawing.PointsToPixels(vr.dpi, vr.s.FontSize))
		if vr.c.textTheta == 0.0 {
			return
		}
		box = box.Corners().Rotate(RadiansToDegrees(vr.c.textTheta)).Box()
	}
	return
}

// SetTextRotation sets the text rotation.
func (vr *vectorRenderer) SetTextRotation(radians float64) {
	vr.c.textTheta = radians
}

// ClearTextRotation clears the text rotation.
func (vr *vectorRenderer) ClearTextRotation() {
	vr.c.textTheta = 0.0
}

// Save saves the renderer's contents to a writer.
func (vr *vectorRenderer) Save(w io.Writer) error {
	vr.c.End()
	_, err := w.Write(vr.b.Bytes())
	bytebufferpool.Put(vr.b)
	return err
}

func newCanvas(w *bytebufferpool.ByteBuffer) *canvas {
	return &canvas{
		w:   w,
		dpi: DefaultDPI,
	}
}

type canvas struct {
	w         *bytebufferpool.ByteBuffer
	dpi       float64
	textTheta float64
	width     int
	height    int
	css       string
	nonce     string
}

var (
	canvasStart  = []byte(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewbox="0 0 `)
	canvasWidth  = []byte(`" width="`)
	canvasHeight = []byte(`" height="`)
	canvasEnd    = []byte("\">\n")
)

func (c *canvas) Start(width, height int) {
	c.width = width
	c.height = height
	sWidth := itoa(c.width)
	sHeight := itoa(c.height)

	_, _ = c.w.Write(canvasStart)
	_, _ = c.w.WriteString(sWidth)
	_, _ = c.w.Write(space)
	_, _ = c.w.WriteString(sHeight)
	_, _ = c.w.Write(canvasWidth)
	_, _ = c.w.WriteString(sWidth)
	_, _ = c.w.Write(canvasHeight)
	_, _ = c.w.WriteString(sHeight)
	_, _ = c.w.Write(canvasEnd)
	if c.css != "" {
		_, _ = c.w.Write([]byte(`<style type="text/css"`))
		if c.nonce != "" {
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy
			_, _ = c.w.Write([]byte(fmt.Sprintf(` nonce="%s"`, c.nonce)))
		}
		// To avoid compatibility issues between XML and CSS (f.e. with child selectors) we should encapsulate the CSS with CDATA.
		_, _ = c.w.Write([]byte(fmt.Sprintf(`><![CDATA[%s]]></style>`, c.css)))
	}
}

var (
	pathStart = []byte("<path ")
	pathD     = []byte("d=\"")
	pathMark  = []byte("\" ")
	pathEnd   = []byte("/>")
)

func (c *canvas) Path(d string, style Style) {
	_, _ = c.w.Write(pathStart)
	if len(style.StrokeDashArray) > 0 {
		_, _ = c.w.WriteString(c.getStrokeDashArray(style))
	}
	_, _ = c.w.Write(pathD)
	_, _ = c.w.WriteString(d)
	_, _ = c.w.Write(pathMark)
	_, _ = c.w.WriteString(c.styleAsSVG(style))
	_, _ = c.w.Write(pathEnd)
}

var (
	textStart = []byte(`<text x="`)
	textY     = []byte(`" y="`)
	textMark  = []byte(`" `)
	textMark2 = []byte(`>`)
	textEnd   = []byte(`</text>`)

	transformStarts = []byte(`transform="rotate(`)
	transformCoords = []byte(`,`)
	transformEnds   = []byte(`)"`)
)

func (c *canvas) Text(x, y int, body string, style Style) {
	sX := itoa(x)
	sY := itoa(y)

	_, _ = c.w.Write(textStart)
	_, _ = c.w.WriteString(sX)
	_, _ = c.w.Write(textY)
	_, _ = c.w.WriteString(sY)
	_, _ = c.w.Write(textMark)
	_, _ = c.w.WriteString(c.styleAsSVG(style))

	if c.textTheta != 0.0 {
		_, _ = c.w.Write(transformStarts)
		_, _ = c.w.WriteString(ftoa2(RadiansToDegrees(c.textTheta)))
		_, _ = c.w.Write(transformCoords)
		_, _ = c.w.WriteString(sX)
		_, _ = c.w.Write(transformCoords)
		_, _ = c.w.WriteString(sY)
		_, _ = c.w.Write(transformEnds)
	}
	_, _ = c.w.Write(textMark2)
	_, _ = c.w.WriteString(body)
	_, _ = c.w.Write(textEnd)
}

var (
	circleStart  = []byte(`<circle cx="`)
	circleCY     = []byte(`" cy="`)
	circleRadius = []byte(`" r="`)
	circleMark   = []byte(`" `)
	circleEnd    = []byte(`/>`)
)

func (c *canvas) Circle(x, y, r int, style Style) {
	_, _ = c.w.Write(circleStart)
	_, _ = c.w.WriteString(itoa(x))
	_, _ = c.w.Write(circleCY)
	_, _ = c.w.WriteString(itoa(y))
	_, _ = c.w.Write(circleRadius)
	_, _ = c.w.WriteString(itoa(r))
	_, _ = c.w.Write(circleMark)
	_, _ = c.w.WriteString(c.styleAsSVG(style))
	_, _ = c.w.Write(circleEnd)
}

var (
	svgEnd = []byte("</svg>")
)

func (c *canvas) End() {
	_, _ = c.w.Write(svgEnd)
}

// getStrokeDashArray returns the stroke-dasharray property of a style.
func (*canvas) getStrokeDashArray(s Style) string {
	if len(s.StrokeDashArray) > 0 {
		values := make([]string, 0, len(s.StrokeDashArray))
		for _, v := range s.StrokeDashArray {
			values = append(values, ftoa1(v))
		}
		return "stroke-dasharray=\"" + strings.Join(values, ", ") + "\" "
	}
	return ""
}

// GetFontFace returns the font face for the style.
func (*canvas) getFontFace(s Style) string {
	family := "sans-serif"
	if s.GetFont() != nil {
		name := s.GetFont().Name(truetype.NameIDFontFamily)
		if len(name) != 0 {
			family = `'` + name + `', ` + family
		}
	}
	return "font-family:" + family
}

// styleAsSVG returns the style as a svg style or class string.
func (c *canvas) styleAsSVG(s Style) string {
	sc := s.StrokeColor
	fc := s.FillColor
	fs := s.FontSize
	fnc := s.FontColor

	if s.ClassName != "" {
		var classes []string
		classes = append(classes, s.ClassName)
		if !sc.IsZero() {
			classes = append(classes, "stroke")
		}
		if !fc.IsZero() {
			classes = append(classes, "fill")
		}
		if fs != 0 || s.Font != nil {
			classes = append(classes, "text")
		}

		return fmt.Sprintf("class=\"%s\"", strings.Join(classes, " "))
	}

	var pieces []string

	if sw := s.StrokeWidth; sw != 0 {
		pieces = append(pieces, "stroke-width:"+fmt.Sprintf("%d", int(sw)))
	} else {
		pieces = append(pieces, "stroke-width:0")
	}

	if !sc.IsZero() {
		pieces = append(pieces, "stroke:"+sc.String())
	} else {
		pieces = append(pieces, "stroke:none")
	}

	switch {
	case !fnc.IsZero():
		pieces = append(pieces, "fill:"+fnc.String())
	case !fc.IsZero():
		pieces = append(pieces, "fill:"+fc.String())
	default:
		pieces = append(pieces, "fill:none")
	}

	if fs != 0 {
		pieces = append(pieces, "font-size:"+ftoa1(drawing.PointsToPixels(c.dpi, fs))+"px")
	}

	if s.Font != nil {
		pieces = append(pieces, c.getFontFace(s))
	}
	return "style=\"" + strings.Join(pieces, ";") + "\""
}
