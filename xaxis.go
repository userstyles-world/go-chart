package chart

import (
	"math"

	"github.com/userstyles-world/go-chart/v2/drawing"
)

// HideXAxis hides the x-axis.
func HideXAxis() XAxis {
	return XAxis{
		Style: Hidden(),
	}
}

// XAxis represents the horizontal axis.
type XAxis struct {
	Name      string
	NameStyle Style

	Style          Style
	ValueFormatter ValueFormatter
	Range          Range

	TickStyle    Style
	Ticks        []Tick
	TickPosition TickPosition

	GridLines      []GridLine
	GridMajorStyle Style
	GridMinorStyle Style
}

// GetName returns the name.
func (xa XAxis) GetName() string {
	return xa.Name
}

// GetStyle returns the style.
func (xa XAxis) GetStyle() Style {
	return xa.Style
}

// GetValueFormatter returns the value formatter for the axis.
func (xa XAxis) GetValueFormatter() ValueFormatter {
	if xa.ValueFormatter != nil {
		return xa.ValueFormatter
	}
	return FloatValueFormatter
}

// GetTickPosition returns the tick position option for the axis.
func (xa XAxis) GetTickPosition(defaults ...TickPosition) TickPosition {
	if xa.TickPosition == TickPositionUnset {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return TickPositionUnderTick
	}
	return xa.TickPosition
}

// GetTicks returns the ticks for a series.
// The coalesce priority is:
// 	- User Supplied Ticks (i.e. Ticks array on the axis itself).
// 	- Range ticks (i.e. if the range provides ticks).
//	- Generating continuous ticks based on minimum spacing and canvas width.
func (xa XAxis) GetTicks(r Renderer, ra Range, defaults Style, vf ValueFormatter) []Tick {
	if len(xa.Ticks) > 0 {
		return xa.Ticks
	}
	if tp, isTickProvider := ra.(TicksProvider); isTickProvider {
		return tp.GetTicks(r, defaults, vf)
	}
	tickStyle := xa.Style.InheritFrom(defaults)
	return GenerateContinuousTicks(r, ra, false, tickStyle, vf)
}

// GetGridLines returns the gridlines for the axis.
func (xa XAxis) GetGridLines(ticks []Tick) []GridLine {
	if len(xa.GridLines) > 0 {
		return xa.GridLines
	}
	return GenerateGridLines(ticks, xa.GridMajorStyle, xa.GridMinorStyle)
}

// Measure returns the bounds of the axis.
func (xa XAxis) Measure(r Renderer, canvasBox Box, ra Range, defaults Style, ticks []Tick) Box {
	tickStyle := xa.TickStyle.InheritFrom(xa.Style.InheritFrom(defaults))

	tp := xa.GetTickPosition()

	var ltx, rtx int
	var tx, ty int
	var left, right, bottom = math.MaxInt32, 0, 0
	for index, t := range ticks {
		v := t.Value
		tb := Draw.MeasureText(r, t.Label, tickStyle.GetTextOptions())

		tx = canvasBox.Left + ra.Translate(v)
		ty = canvasBox.Bottom + DefaultXAxisMargin + tb.Height()
		switch tp {
		case TickPositionUnderTick, TickPositionUnset:
			ltx = tx - tb.Width()>>1
			rtx = tx + tb.Width()>>1
		case TickPositionBetweenTicks:
			if index > 0 {
				ltx = ra.Translate(ticks[index-1].Value)
				rtx = tx
			}
		}

		left = MinInt(left, ltx)
		right = MaxInt(right, rtx)
		bottom = MaxInt(bottom, ty)
	}

	if !xa.NameStyle.Hidden && len(xa.Name) > 0 {
		tb := Draw.MeasureText(r, xa.Name, xa.NameStyle.InheritFrom(defaults))
		bottom += DefaultXAxisMargin + tb.Height()
	}

	return Box{
		Top:    canvasBox.Bottom,
		Left:   left,
		Right:  right,
		Bottom: bottom,
	}
}

// Render renders the axis
func (xa XAxis) Render(r Renderer, canvasBox Box, ra Range, defaults Style, ticks []Tick) {
	tickStyle := xa.TickStyle.InheritFrom(xa.Style.InheritFrom(defaults))

	tickStyle.GetStrokeOptions().WriteToRenderer(r)
	r.MoveTo(canvasBox.Left, canvasBox.Bottom)
	r.LineTo(canvasBox.Right, canvasBox.Bottom)
	r.Stroke()

	tp := xa.GetTickPosition()

	var tx, ty int
	var maxTextHeight int
	for index, t := range ticks {
		v := t.Value
		lx := ra.Translate(v)

		tx = canvasBox.Left + lx

		tickStyle.GetStrokeOptions().WriteToRenderer(r)
		r.MoveTo(tx, canvasBox.Bottom)
		r.LineTo(tx, canvasBox.Bottom+DefaultVerticalTickHeight)
		r.Stroke()

		tickWithAxisStyle := xa.TickStyle.InheritFrom(xa.Style.InheritFrom(defaults))
		tb := Draw.MeasureText(r, t.Label, tickWithAxisStyle)

		switch tp {
		case TickPositionUnderTick, TickPositionUnset:
			if tickStyle.TextRotationDegrees == 0 {
				tx -= tb.Width() >> 1
				ty = canvasBox.Bottom + DefaultXAxisMargin + tb.Height()
			} else {
				ty = canvasBox.Bottom + (2 * DefaultXAxisMargin)
			}
			Draw.Text(r, t.Label, tx, ty, tickWithAxisStyle)
			maxTextHeight = MaxInt(maxTextHeight, tb.Height())
		case TickPositionBetweenTicks:
			if index > 0 {
				llx := ra.Translate(ticks[index-1].Value)
				ltx := canvasBox.Left + llx
				finalTickStyle := tickWithAxisStyle.InheritFrom(Style{TextHorizontalAlign: TextHorizontalAlignCenter})

				Draw.TextWithin(r, t.Label, Box{
					Left:   ltx,
					Right:  tx,
					Top:    canvasBox.Bottom + DefaultXAxisMargin,
					Bottom: canvasBox.Bottom + DefaultXAxisMargin,
				}, finalTickStyle)

				ftb := Text.MeasureLines(r, Text.WrapFit(r, t.Label, tx-ltx, finalTickStyle), finalTickStyle)
				maxTextHeight = MaxInt(maxTextHeight, ftb.Height())
			}
		}
	}

	nameStyle := xa.NameStyle.InheritFrom(defaults)
	if !xa.NameStyle.Hidden && len(xa.Name) > 0 {
		tb := Draw.MeasureText(r, xa.Name, nameStyle)
		tx := canvasBox.Right - (canvasBox.Width()>>1 + tb.Width()>>1)
		ty := canvasBox.Bottom + DefaultXAxisMargin + maxTextHeight + DefaultXAxisMargin + tb.Height()
		Draw.Text(r, xa.Name, tx, ty, nameStyle)
	}

	if !xa.GridMajorStyle.Hidden || !xa.GridMinorStyle.Hidden {
		isMinorDefault := xa.GridMinorStyle.StrokeColor == drawing.Color{}
		isMajorDefault := xa.GridMajorStyle.StrokeColor == drawing.Color{}
		// If the grid style is not set, then skip it.
		if isMinorDefault && isMajorDefault {
			return
		}

		for _, gl := range xa.GetGridLines(ticks) {
			if (gl.IsMinor && !xa.GridMinorStyle.Hidden) || (!gl.IsMinor && !xa.GridMajorStyle.Hidden) {
				var defaults Style

				if gl.IsMinor {
					if isMinorDefault {
						break
					}
					defaults = xa.GridMinorStyle
				} else {
					if isMajorDefault {
						break
					}
					defaults = xa.GridMajorStyle
				}

				gl.Render(r, canvasBox, ra, true, gl.Style.InheritFrom(defaults))
			}
		}
	}
}
