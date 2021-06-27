package chart

import (
	"bytes"
	"testing"

	"github.com/userstyles-world/go-chart/v2/testutil"
)

func TestPieChart(t *testing.T) {
	// replaced new assertions helper

	pie := PieChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 10, Label: "Blue"},
			{Value: 9, Label: "Green"},
			{Value: 8, Label: "Gray"},
			{Value: 7, Label: "Orange"},
			{Value: 6, Label: "HEANG"},
			{Value: 5, Label: "??"},
			{Value: 2, Label: "!!"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	err := pie.Render(PNG, b)
	testutil.AssertNil(t, err)
	testutil.AssertNotZero(t, b.Len())
}

func TestPieChartDropsZeroValues(t *testing.T) {
	// replaced new assertions helper

	pie := PieChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 5, Label: "Blue"},
			{Value: 5, Label: "Green"},
			{Value: 0, Label: "Gray"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	err := pie.Render(PNG, b)
	testutil.AssertNil(t, err)
}

func TestPieChartAllZeroValues(t *testing.T) {
	// replaced new assertions helper

	pie := PieChart{
		Canvas: Style{
			FillColor: ColorLightGray,
		},
		Values: []Value{
			{Value: 0, Label: "Blue"},
			{Value: 0, Label: "Green"},
			{Value: 0, Label: "Gray"},
		},
	}

	b := bytes.NewBuffer([]byte{})
	err := pie.Render(PNG, b)
	testutil.AssertNotNil(t, err)
}
