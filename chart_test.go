package chart

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/userstyles-world/go-chart/v2/drawing"
	"github.com/userstyles-world/go-chart/v2/testutil"
)

func TestChartGetDPI(t *testing.T) {
	// replaced new assertions helper

	unset := Chart{}
	testutil.AssertEqual(t, DefaultDPI, unset.GetDPI())
	testutil.AssertEqual(t, 192, unset.GetDPI(192))

	set := Chart{DPI: 128}
	testutil.AssertEqual(t, 128, set.GetDPI())
	testutil.AssertEqual(t, 128, set.GetDPI(192))
}

func TestChartGetFont(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	unset := Chart{}
	testutil.AssertNil(t, unset.GetFont())

	set := Chart{Font: f}
	testutil.AssertNotNil(t, set.GetFont())
}

func TestChartGetWidth(t *testing.T) {
	// replaced new assertions helper

	unset := Chart{}
	testutil.AssertEqual(t, DefaultChartWidth, unset.GetWidth())

	set := Chart{Width: DefaultChartWidth + 10}
	testutil.AssertEqual(t, DefaultChartWidth+10, set.GetWidth())
}

func TestChartGetHeight(t *testing.T) {
	// replaced new assertions helper

	unset := Chart{}
	testutil.AssertEqual(t, DefaultChartHeight, unset.GetHeight())

	set := Chart{Height: DefaultChartHeight + 10}
	testutil.AssertEqual(t, DefaultChartHeight+10, set.GetHeight())
}

func TestChartGetRanges(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
			ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{-2.1, -1.0, 0, 1.0, 2.0},
			},
			ContinuousSeries{
				YAxis:   YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{10.0, 11.0, 12.0, 13.0, 14.0},
			},
		},
	}

	xrange, yrange, yrangeAlt := c.getRanges()
	testutil.AssertEqual(t, -2.0, xrange.GetMin())
	testutil.AssertEqual(t, 5.0, xrange.GetMax())

	testutil.AssertEqual(t, -2.1, yrange.GetMin())
	testutil.AssertEqual(t, 4.5, yrange.GetMax())

	testutil.AssertEqual(t, 10.0, yrangeAlt.GetMin())
	testutil.AssertEqual(t, 14.0, yrangeAlt.GetMax())

	cSet := Chart{
		XAxis: XAxis{
			Range: &ContinuousRange{Min: 9.8, Max: 19.8},
		},
		YAxis: YAxis{
			Range: &ContinuousRange{Min: 9.9, Max: 19.9},
		},
		YAxisSecondary: YAxis{
			Range: &ContinuousRange{Min: 9.7, Max: 19.7},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
			ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{-2.1, -1.0, 0, 1.0, 2.0},
			},
			ContinuousSeries{
				YAxis:   YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{10.0, 11.0, 12.0, 13.0, 14.0},
			},
		},
	}

	xr2, yr2, yra2 := cSet.getRanges()
	testutil.AssertEqual(t, 9.8, xr2.GetMin())
	testutil.AssertEqual(t, 19.8, xr2.GetMax())

	testutil.AssertEqual(t, 9.9, yr2.GetMin())
	testutil.AssertEqual(t, 19.9, yr2.GetMax())

	testutil.AssertEqual(t, 9.7, yra2.GetMin())
	testutil.AssertEqual(t, 19.7, yra2.GetMax())
}

func TestChartGetRangesUseTicks(t *testing.T) {
	// replaced new assertions helper

	// this test asserts that ticks should supercede manual ranges when generating the overall ranges.

	c := Chart{
		YAxis: YAxis{
			Ticks: []Tick{
				{0.0, "Zero"},
				{1.0, "1.0"},
				{2.0, "2.0"},
				{3.0, "3.0"},
				{4.0, "4.0"},
				{5.0, "Five"},
			},
			Range: &ContinuousRange{
				Min: -5.0,
				Max: 5.0,
			},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
		},
	}

	xr, yr, yar := c.getRanges()
	testutil.AssertEqual(t, -2.0, xr.GetMin())
	testutil.AssertEqual(t, 2.0, xr.GetMax())
	testutil.AssertEqual(t, 0.0, yr.GetMin())
	testutil.AssertEqual(t, 5.0, yr.GetMax())
	testutil.AssertTrue(t, yar.IsZero(), yar.String())
}

func TestChartGetRangesUseUserRanges(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: -5.0,
				Max: 5.0,
			},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
		},
	}

	xr, yr, yar := c.getRanges()
	testutil.AssertEqual(t, -2.0, xr.GetMin())
	testutil.AssertEqual(t, 2.0, xr.GetMax())
	testutil.AssertEqual(t, -5.0, yr.GetMin())
	testutil.AssertEqual(t, 5.0, yr.GetMax())
	testutil.AssertTrue(t, yar.IsZero(), yar.String())
}

func TestChartGetBackgroundStyle(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		Background: Style{
			FillColor: drawing.ColorBlack,
		},
	}

	bs := c.getBackgroundStyle()
	testutil.AssertEqual(t, bs.FillColor.String(), drawing.ColorBlack.String())
}

func TestChartGetCanvasStyle(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		Canvas: Style{
			FillColor: drawing.ColorBlack,
		},
	}

	bs := c.getCanvasStyle()
	testutil.AssertEqual(t, bs.FillColor.String(), drawing.ColorBlack.String())
}

func TestChartGetDefaultCanvasBox(t *testing.T) {
	// replaced new assertions helper

	c := Chart{}
	canvasBoxDefault := c.getDefaultCanvasBox()
	testutil.AssertFalse(t, canvasBoxDefault.IsZero())
	testutil.AssertEqual(t, DefaultBackgroundPadding.Top, canvasBoxDefault.Top)
	testutil.AssertEqual(t, DefaultBackgroundPadding.Left, canvasBoxDefault.Left)
	testutil.AssertEqual(t, c.GetWidth()-DefaultBackgroundPadding.Right, canvasBoxDefault.Right)
	testutil.AssertEqual(t, c.GetHeight()-DefaultBackgroundPadding.Bottom, canvasBoxDefault.Bottom)

	custom := Chart{
		Background: Style{
			Padding: Box{
				Top:    DefaultBackgroundPadding.Top + 1,
				Left:   DefaultBackgroundPadding.Left + 1,
				Right:  DefaultBackgroundPadding.Right + 1,
				Bottom: DefaultBackgroundPadding.Bottom + 1,
			},
		},
	}
	canvasBoxCustom := custom.getDefaultCanvasBox()
	testutil.AssertFalse(t, canvasBoxCustom.IsZero())
	testutil.AssertEqual(t, DefaultBackgroundPadding.Top+1, canvasBoxCustom.Top)
	testutil.AssertEqual(t, DefaultBackgroundPadding.Left+1, canvasBoxCustom.Left)
	testutil.AssertEqual(t, c.GetWidth()-(DefaultBackgroundPadding.Right+1), canvasBoxCustom.Right)
	testutil.AssertEqual(t, c.GetHeight()-(DefaultBackgroundPadding.Bottom+1), canvasBoxCustom.Bottom)
}

func TestChartGetValueFormatters(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
			ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{-2.1, -1.0, 0, 1.0, 2.0},
			},
			ContinuousSeries{
				YAxis:   YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{10.0, 11.0, 12.0, 13.0, 14.0},
			},
		},
	}

	dxf, dyf, dyaf := c.getValueFormatters()
	testutil.AssertNotNil(t, dxf)
	testutil.AssertNotNil(t, dyf)
	testutil.AssertNotNil(t, dyaf)
}

func TestChartHasAxes(t *testing.T) {
	// replaced new assertions helper

	testutil.AssertTrue(t, Chart{}.hasAxes())
	testutil.AssertFalse(t, Chart{XAxis: XAxis{Style: Hidden()}, YAxis: YAxis{Style: Hidden()}, YAxisSecondary: YAxis{Style: Hidden()}}.hasAxes())

	x := Chart{
		XAxis: XAxis{
			Style: Hidden(),
		},
		YAxis: YAxis{
			Style: Shown(),
		},
		YAxisSecondary: YAxis{
			Style: Hidden(),
		},
	}
	testutil.AssertTrue(t, x.hasAxes())

	y := Chart{
		XAxis: XAxis{
			Style: Shown(),
		},
		YAxis: YAxis{
			Style: Hidden(),
		},
		YAxisSecondary: YAxis{
			Style: Hidden(),
		},
	}
	testutil.AssertTrue(t, y.hasAxes())

	ya := Chart{
		XAxis: XAxis{
			Style: Hidden(),
		},
		YAxis: YAxis{
			Style: Hidden(),
		},
		YAxisSecondary: YAxis{
			Style: Shown(),
		},
	}
	testutil.AssertTrue(t, ya.hasAxes())
}

func TestChartGetAxesTicks(t *testing.T) {
	// replaced new assertions helper

	r, err := PNG(1024, 1024)
	testutil.AssertNil(t, err)

	c := Chart{
		XAxis: XAxis{
			Range: &ContinuousRange{Min: 9.8, Max: 19.8},
		},
		YAxis: YAxis{
			Range: &ContinuousRange{Min: 9.9, Max: 19.9},
		},
		YAxisSecondary: YAxis{
			Range: &ContinuousRange{Min: 9.7, Max: 19.7},
		},
	}
	xr, yr, yar := c.getRanges()

	xt, yt, yat := c.getAxesTicks(r, xr, yr, yar, FloatValueFormatter, FloatValueFormatter, FloatValueFormatter)
	testutil.AssertNotEmpty(t, xt)
	testutil.AssertNotEmpty(t, yt)
	testutil.AssertNotEmpty(t, yat)
}

func TestChartSingleSeries(t *testing.T) {
	// replaced new assertions helper
	now := time.Now()
	c := Chart{
		Title:  "Hello!",
		Width:  1024,
		Height: 400,
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: 0.0,
				Max: 4.0,
			},
		},
		Series: []Series{
			TimeSeries{
				Name:    "goog",
				XValues: []time.Time{now.AddDate(0, 0, -3), now.AddDate(0, 0, -2), now.AddDate(0, 0, -1)},
				YValues: []float64{1.0, 2.0, 3.0},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := c.Render(PNG, buffer)
	testutil.AssertNil(t, err)
	testutil.AssertNotEmpty(t, buffer.Bytes())
}

func TestChartRegressionBadRanges(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1)},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
		},
	}
	buffer := bytes.NewBuffer([]byte{})
	err := c.Render(PNG, buffer)
	testutil.AssertNotNil(t, err)
	testutil.AssertEqual(t, "infinite x-range delta", err.Error(), "Should error about infinite x-delta")
	testutil.AssertTrue(t, true, "Render needs to finish.")
}

func TestChartRegressionBadRangesByUser(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: math.Inf(-1),
				Max: math.Inf(1), // this could really happen? eh.
			},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: LinearRange(1.0, 10.0),
				YValues: LinearRange(1.0, 10.0),
			},
		},
	}
	buffer := bytes.NewBuffer([]byte{})
	err := c.Render(PNG, buffer)
	testutil.AssertNotNil(t, err)
	testutil.AssertEqual(t, "infinite y-range delta", err.Error(), "Should error about infinite y-delta")
	testutil.AssertTrue(t, true, "Render needs to finish.")
}

func TestChartValidatesSeries(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: LinearRange(1.0, 10.0),
				YValues: LinearRange(1.0, 10.0),
			},
		},
	}

	testutil.AssertNil(t, c.validateSeries())

	c = Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: LinearRange(1.0, 10.0),
			},
		},
	}

	testutil.AssertNotNil(t, c.validateSeries())
}

func TestChartCheckRanges(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{1.0, 2.0},
				YValues: []float64{3.10, 3.14},
			},
		},
	}

	xr, yr, yra := c.getRanges()
	testutil.AssertNil(t, c.checkRanges(xr, yr, yra))
}

func TestChartCheckRangesWithRanges(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		XAxis: XAxis{
			Range: &ContinuousRange{
				Min: 0,
				Max: 10,
			},
		},
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: 0,
				Max: 5,
			},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{1.0, 2.0},
				YValues: []float64{3.14, 3.14},
			},
		},
	}

	xr, yr, yra := c.getRanges()
	testutil.AssertNil(t, c.checkRanges(xr, yr, yra))
}

func at(i image.Image, x, y int) drawing.Color {
	return drawing.ColorFromAlphaMixedRGBA(i.At(x, y).RGBA())
}

func TestChartE2ELine(t *testing.T) {
	// replaced new assertions helper

	c := Chart{
		Height:         50,
		Width:          50,
		TitleStyle:     Hidden(),
		XAxis:          HideXAxis(),
		YAxis:          HideYAxis(),
		YAxisSecondary: HideYAxis(),
		Canvas: Style{
			Padding: BoxZero,
		},
		Background: Style{
			Padding: BoxZero,
		},
		Series: []Series{
			ContinuousSeries{
				XValues: LinearRangeWithStep(0, 4, 1),
				YValues: LinearRangeWithStep(0, 4, 1),
			},
		},
	}

	var buffer = &bytes.Buffer{}
	err := c.Render(PNG, buffer)
	testutil.AssertNil(t, err)

	// do color tests ...

	i, err := png.Decode(buffer)
	testutil.AssertNil(t, err)

	// test the bottom and top of the line
	testutil.AssertEqual(t, drawing.ColorWhite, at(i, 0, 0))
	testutil.AssertEqual(t, drawing.ColorWhite, at(i, 49, 49))

	// test a line mid point
	defaultSeriesColor := GetDefaultColor(0)
	testutil.AssertEqual(t, defaultSeriesColor, at(i, 0, 49))
	testutil.AssertEqual(t, defaultSeriesColor, at(i, 49, 0))
	testutil.AssertEqual(t, drawing.ColorFromHex("bddbf6"), at(i, 24, 24))
}

func TestChartE2ELineWithFill(t *testing.T) {
	// replaced new assertions helper

	logBuffer := new(bytes.Buffer)

	c := Chart{
		Height: 50,
		Width:  50,
		Canvas: Style{
			Padding: BoxZero,
		},
		Background: Style{
			Padding: BoxZero,
		},
		TitleStyle:     Hidden(),
		XAxis:          HideXAxis(),
		YAxis:          HideYAxis(),
		YAxisSecondary: HideYAxis(),
		Series: []Series{
			ContinuousSeries{
				Style: Style{
					StrokeColor: drawing.ColorBlue,
					FillColor:   drawing.ColorRed,
				},
				XValues: LinearRangeWithStep(0, 4, 1),
				YValues: LinearRangeWithStep(0, 4, 1),
			},
		},
		Log: NewLogger(OptLoggerStdout(logBuffer), OptLoggerStderr(logBuffer)),
	}

	testutil.AssertEqual(t, 5, len(c.Series[0].(ContinuousSeries).XValues))
	testutil.AssertEqual(t, 5, len(c.Series[0].(ContinuousSeries).YValues))

	var buffer = &bytes.Buffer{}
	err := c.Render(PNG, buffer)
	testutil.AssertNil(t, err)

	i, err := png.Decode(buffer)
	testutil.AssertNil(t, err)

	// test the bottom and top of the line
	testutil.AssertEqual(t, drawing.ColorWhite, at(i, 0, 0))
	testutil.AssertEqual(t, drawing.ColorRed, at(i, 49, 49))

	// test a line mid point
	defaultSeriesColor := drawing.ColorBlue
	testutil.AssertEqual(t, defaultSeriesColor, at(i, 0, 49))
	testutil.AssertEqual(t, defaultSeriesColor, at(i, 49, 0))
}

// Regression: check when the XAsis has no majorGrid or minorGrid.
// It should not draw the grid.
func TestChartEmptyGrid(t *testing.T) {
	now := time.Now()

	c := Chart{
		Width:      1248,
		Canvas:     Style{ClassName: "bg inner"},
		Background: Style{ClassName: "bg outer"},
		XAxis:      XAxis{Name: "Date"},
		YAxis:      YAxis{Name: "Daily count"},
		Series: []Series{
			TimeSeries{
				Name:    "goog",
				XValues: []time.Time{now.AddDate(0, 0, -3), now.AddDate(0, 0, -2), now.AddDate(0, 0, -1)},
				YValues: []float64{1.0, 2.0, 3.0},
			},
		},
	}

	var buffer = &bytes.Buffer{}
	err := c.Render(SVG, buffer)
	testutil.AssertNil(t, err)

	// check the grid is not drawn
	// Should match: <path d="M 38 -9223372036854775457 L 1185 -9223372036854775457" style="stroke-width:0;stroke:none;fill:none"/>
	re := regexp.MustCompile(`<path d=\"[^"]*\" ?style="stroke-width:0;stroke:none;fill:none" ?(\/>| ?<\/path>)`)
	matches := re.FindAllString(buffer.String(), -1)
	testutil.AssertEqual(t, 0, len(matches))
}

// Regression: go-charts shouldn't render paths with multiple spaces in the path.
func TestChartDoubleSpace(t *testing.T) {
	t.Parallel()

	c := Chart{
		Width:      1248,
		Canvas:     Style{ClassName: "bg inner"},
		Background: Style{ClassName: "bg outer"},
		XAxis:      XAxis{Name: "Date"},
		YAxis:      YAxis{Name: "Daily count"},
		Series: []Series{
			TimeSeries{
				Name:    "goog",
				XValues: []time.Time{time.Now(), time.Now().AddDate(0, 0, -1)},
				YValues: []float64{1.0, 2.0},
			},
		},
	}
	var buffer = &bytes.Buffer{}
	err := c.Render(SVG, buffer)
	testutil.AssertNil(t, err)

	re := regexp.MustCompile(` {2,}`)
	matches := re.FindAllString(buffer.String(), -1)
	testutil.AssertEqual(t, 0, len(matches))
}

// Regression: check if their are no <text> with negative x or y values
func TestChartTextValues(t *testing.T) {
	now := time.Now()

	c := Chart{
		Width:      1248,
		Canvas:     Style{ClassName: "bg inner"},
		Background: Style{ClassName: "bg outer"},
		XAxis:      XAxis{Name: "Date"},
		YAxis:      YAxis{Name: "Daily count"},
		Series: []Series{
			TimeSeries{
				Name:    "goog",
				XValues: []time.Time{now.AddDate(0, 0, -3), now.AddDate(0, 0, -2), now.AddDate(0, 0, -1)},
				YValues: []float64{1.0, 2.0, 3.0},
			},
		},
	}

	var buffer = &bytes.Buffer{}
	err := c.Render(SVG, buffer)
	testutil.AssertNil(t, err)

	// check there are no negative values in the text
	// x="-9" y="-9"
	re := regexp.MustCompile(`(x="|y=")-\d+"`)
	matches := re.FindAllString(buffer.String(), -1)
	testutil.AssertEqual(t, 0, len(matches))
}

func BenchmarkBarChartLegend(b *testing.B) {
	// Exported data [6 4 9 3 11 3 8 4 5 11 8 8 8 9 9 9 8 5 6 13 10 6 8 7 11 5 16 6 9 8 7 8 8 17 11 11 11 10 14 6 19 14 9 11 10 13 6 9 12]
	dailyViews := []float64{6, 4, 9, 3, 11, 3, 8, 4, 5, 11, 8, 8, 8, 9, 9, 9, 8, 5, 6, 13, 10, 6, 8, 7, 11, 5, 16, 6, 9, 8, 7, 8, 8, 17, 11, 11, 11, 10, 14, 6, 19, 14, 9, 11, 10, 13, 6, 9, 12}

	// Exported data [25 22 26 23 23 26 28 30 34 30 29 26 28 31 31 35 32 30 30 35 43 40 40 37 34 36 38 39 35 40 39 31 44 43 49 44 54 48 49 59 58 56 58 63 58 62 59 71 72]
	dailyUpdates := []float64{25, 22, 26, 23, 23, 26, 28, 30, 34, 30, 29, 26, 28, 31, 31, 35, 32, 30, 30, 35, 43, 40, 40, 37, 34, 36, 38, 39, 35, 40, 39, 31, 44, 43, 49, 44, 54, 48, 49, 59, 58, 56, 58, 63, 58, 62, 59, 71, 72}

	// Exported data [25 7 8 5 5 4 9 9 10 7 4 5 10 3 7 13 9 7 7 11 14 10 13 12 10 11 12 12 13 6 11 11 12 13 12 11 12 14 11 12 15 14 13 18 19 25 13 19 13]
	dailyInstalls := []float64{25, 7, 8, 5, 5, 4, 9, 9, 10, 7, 4, 5, 10, 3, 7, 13, 9, 7, 7, 11, 14, 10, 13, 12, 10, 11, 12, 12, 13, 6, 11, 11, 12, 13, 12, 11, 12, 14, 11, 12, 15, 14, 13, 18, 19, 25, 13, 19, 13}

	// Exported data [2021-07-07 23:59:00.284008573 +0000 UTC 2021-07-08 23:59:00.448171544 +0000 UTC 2021-07-09 23:59:00.430622428 +0000 UTC 2021-07-10 23:59:00.518588335 +0000 UTC 2021-07-11 23:59:00.598871418 +0000 UTC 2021-07-12 23:59:00.762924884 +0000 UTC 2021-07-13 23:59:00.894566364 +0000 UTC 2021-07-14 23:59:01.183382398 +0000 UTC 2021-07-15 23:59:02.110516404 +0000 UTC 2021-07-16 23:59:03.01251904 +0000 UTC 2021-07-17 23:59:03.364826174 +0000 UTC 2021-07-18 23:59:07.053386001 +0000 UTC 2021-07-19 23:59:04.613313208 +0000 UTC 2021-07-20 23:59:05.157727215 +0000 UTC 2021-07-21 23:59:06.256763794 +0000 UTC 2021-07-22 23:59:08.234239808 +0000 UTC 2021-07-23 23:59:07.321565106 +0000 UTC 2021-07-24 23:59:08.298329788 +0000 UTC 2021-07-25 23:59:08.633682537 +0000 UTC 2021-07-26 23:59:10.320324132 +0000 UTC 2021-07-27 23:59:11.124146932 +0000 UTC 2021-07-28 23:59:12.724467815 +0000 UTC 2021-07-29 23:59:26.634052314 +0000 UTC 2021-07-30 23:59:15.329860625 +0000 UTC 2021-07-31 23:59:17.940346159 +0000 UTC 2021-08-01 23:59:18.490248966 +0000 UTC 2021-08-02 23:59:21.941662503 +0000 UTC 2021-08-03 23:59:35.401362307 +0000 UTC 2021-08-04 23:59:00.678379832 +0000 UTC 2021-08-05 23:59:00.8755122 +0000 UTC 2021-08-06 23:59:00.934198575 +0000 UTC 2021-08-07 23:59:00.964950689 +0000 UTC 2021-08-08 23:59:00.829227921 +0000 UTC 2021-08-09 23:59:00.971994166 +0000 UTC 2021-08-10 23:59:00.852921338 +0000 UTC 2021-08-11 23:59:00.933554721 +0000 UTC 2021-08-12 23:59:00.989925044 +0000 UTC 2021-08-13 23:59:01.078751216 +0000 UTC 2021-08-14 23:59:01.113226627 +0000 UTC 2021-08-15 23:59:01.133484412 +0000 UTC 2021-08-16 23:59:01.104150383 +0000 UTC 2021-08-17 23:59:01.148971349 +0000 UTC 2021-08-18 23:59:01.357223006 +0000 UTC 2021-08-19 23:59:01.131506323 +0000 UTC 2021-08-20 23:59:01.272777903 +0000 UTC 2021-08-21 23:59:01.173679573 +0000 UTC 2021-08-22 23:59:01.270419124 +0000 UTC 2021-08-23 23:59:01.266545934 +0000 UTC 2021-08-24 23:59:01.455933761 +0000 UTC]
	dates := []time.Time{
		time.Date(2021, time.July, 7, 23, 59, 0, 0, time.UTC),
		time.Date(2021, time.July, 8, 23, 59, 0, 0, time.UTC),
		time.Date(2021, time.July, 9, 23, 59, 0, 0, time.UTC),
		time.Date(2021, time.July, 10, 23, 59, 0, 518588335, time.UTC),
		time.Date(2021, time.July, 11, 23, 59, 0, 598871418, time.UTC),
		time.Date(2021, time.July, 12, 23, 59, 0, 762924884, time.UTC),
		time.Date(2021, time.July, 13, 23, 59, 1, 894566364, time.UTC),
		time.Date(2021, time.July, 14, 23, 59, 2, 183382398, time.UTC),
		time.Date(2021, time.July, 15, 23, 59, 3, 110516404, time.UTC),
		time.Date(2021, time.July, 16, 23, 59, 4, 0, time.UTC),
		time.Date(2021, time.July, 17, 23, 59, 5, 364826174, time.UTC),
		time.Date(2021, time.July, 18, 23, 59, 7, 0, time.UTC),
		time.Date(2021, time.July, 19, 23, 59, 4, 613313208, time.UTC),
		time.Date(2021, time.July, 20, 23, 59, 5, 157727215, time.UTC),
		time.Date(2021, time.July, 21, 23, 59, 6, 256763794, time.UTC),
		time.Date(2021, time.July, 22, 23, 59, 8, 234239808, time.UTC),
		time.Date(2021, time.July, 23, 23, 59, 7, 321565106, time.UTC),
		time.Date(2021, time.July, 24, 23, 59, 8, 298329788, time.UTC),
		time.Date(2021, time.July, 25, 23, 59, 8, 633682537, time.UTC),
		time.Date(2021, time.July, 26, 23, 59, 10, 320324132, time.UTC),
		time.Date(2021, time.July, 27, 23, 59, 9, 558598892, time.UTC),
		time.Date(2021, time.July, 28, 23, 59, 10, 588790837, time.UTC),
		time.Date(2021, time.July, 29, 23, 59, 11, 805879082, time.UTC),
		time.Date(2021, time.July, 30, 23, 59, 12, 805879082, time.UTC),
		time.Date(2021, time.July, 31, 23, 59, 13, 805879082, time.UTC),
		time.Date(2021, time.August, 1, 23, 59, 14, 805879082, time.UTC),
		time.Date(2021, time.August, 2, 23, 59, 15, 805879082, time.UTC),
		time.Date(2021, time.August, 3, 23, 59, 16, 805879082, time.UTC),
		time.Date(2021, time.August, 4, 23, 59, 17, 805879082, time.UTC),
		time.Date(2021, time.August, 5, 23, 59, 18, 805879082, time.UTC),
		time.Date(2021, time.August, 6, 23, 59, 19, 805879082, time.UTC),
		time.Date(2021, time.August, 7, 23, 59, 20, 805879082, time.UTC),
		time.Date(2021, time.August, 8, 23, 59, 21, 805879082, time.UTC),
		time.Date(2021, time.August, 9, 23, 59, 22, 805879082, time.UTC),
		time.Date(2021, time.August, 10, 23, 59, 23, 805879082, time.UTC),
		time.Date(2021, time.August, 11, 23, 59, 24, 805879082, time.UTC),
		time.Date(2021, time.August, 12, 23, 59, 25, 805879082, time.UTC),
		time.Date(2021, time.August, 13, 23, 59, 26, 805879082, time.UTC),
		time.Date(2021, time.August, 14, 23, 59, 27, 805879082, time.UTC),
		time.Date(2021, time.August, 15, 23, 59, 28, 805879082, time.UTC),
		time.Date(2021, time.August, 16, 23, 59, 29, 805879082, time.UTC),
		time.Date(2021, time.August, 17, 23, 59, 30, 805879082, time.UTC),
		time.Date(2021, time.August, 18, 23, 59, 31, 805879082, time.UTC),
		time.Date(2021, time.August, 19, 23, 59, 32, 805879082, time.UTC),
		time.Date(2021, time.August, 20, 23, 59, 33, 805879082, time.UTC),
		time.Date(2021, time.August, 21, 23, 59, 34, 805879082, time.UTC),
		time.Date(2021, time.August, 22, 23, 59, 35, 805879082, time.UTC),
		time.Date(2021, time.August, 23, 23, 59, 36, 805879082, time.UTC),
		time.Date(2021, time.August, 24, 23, 59, 37, 805879082, time.UTC),
	}

	dailyGraph := Chart{
		Width:      1248,
		Canvas:     Style{ClassName: "bg inner"},
		Background: Style{ClassName: "bg outer"},
		XAxis:      XAxis{Name: "Date"},
		YAxis:      YAxis{Name: "Daily count"},
		Series: []Series{
			TimeSeries{
				Name:    "Daily installs",
				XValues: dates,
				YValues: dailyInstalls,
			},
			TimeSeries{
				Name:    "Daily updates",
				XValues: dates,
				YValues: dailyUpdates,
			},
			TimeSeries{
				Name:    "Daily views",
				XValues: dates,
				YValues: dailyViews,
			},
		},
	}

	dailyGraph.Elements = []Renderable{Legend(&dailyGraph)}

	b.ResetTimer()
	var buffer bytes.Buffer
	for i := 0; i < b.N; i++ {
		_ = dailyGraph.Render(SVG, &buffer)
	}
}
