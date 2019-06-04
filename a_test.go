package bakefiles

import "testing"

type svgTest struct {
	svg []string
}

func TestConfig(t *testing.T) {
	parseConfig("tests/testdata.toml", &svgTest{})
}
