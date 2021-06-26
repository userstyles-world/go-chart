package chart

import (
	"testing"

	"github.com/userstyles-world/go-chart/v2/testutil"
)

func TestXAxisGetTicks(t *testing.T) {
	// replaced new assertions helper

	r, err := PNG(1024, 1024)
	testutil.AssertNil(t, err)

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	xa := XAxis{}
	xr := &ContinuousRange{Min: 10, Max: 100, Domain: 1024}
	styleDefaults := Style{
		Font:     f,
		FontSize: 10.0,
	}
	vf := FloatValueFormatter
	ticks := xa.GetTicks(r, xr, styleDefaults, vf)
	testutil.AssertLen(t, ticks, 16)
}

func TestXAxisGetTicksWithUserDefaults(t *testing.T) {
	// replaced new assertions helper

	r, err := PNG(1024, 1024)
	testutil.AssertNil(t, err)

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	xa := XAxis{
		Ticks: []Tick{{Value: 1.0, Label: "1.0"}},
	}
	xr := &ContinuousRange{Min: 10, Max: 100, Domain: 1024}
	styleDefaults := Style{
		Font:     f,
		FontSize: 10.0,
	}
	vf := FloatValueFormatter
	ticks := xa.GetTicks(r, xr, styleDefaults, vf)
	testutil.AssertLen(t, ticks, 1)
}

func TestXAxisMeasure(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)
	style := Style{
		Font:     f,
		FontSize: 10.0,
	}
	r, err := PNG(100, 100)
	testutil.AssertNil(t, err)
	ticks := []Tick{{Value: 1.0, Label: "1.0"}, {Value: 2.0, Label: "2.0"}, {Value: 3.0, Label: "3.0"}}
	xa := XAxis{}
	xab := xa.Measure(r, NewBox(0, 0, 100, 100), &ContinuousRange{Min: 1.0, Max: 3.0, Domain: 100}, style, ticks)
	testutil.AssertEqual(t, 122, xab.Width())
	testutil.AssertEqual(t, 21, xab.Height())
}
