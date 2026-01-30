// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Hubble

package printer

import (
	"maps"
	"slices"
	"strings"

	"github.com/fatih/color"
)

type sprinter interface {
	Sprint(a ...any) string
}

type extColor struct {
	color      *color.Color
	escapeCode string
}

type colorer struct {
	colors  []extColor
	red     sprinter
	green   sprinter
	blue    sprinter
	cyan    sprinter
	magenta sprinter
	yellow  sprinter
	enabled bool
}

func getCode(c *color.Color) string {
	c.EnableColor()
	seq := c.Sprint("|")
	split := strings.Split(seq, "|")
	return split[0]
}

func newColorer(when string) *colorer {
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	cyan := color.New(color.FgCyan)
	magenta := color.New(color.FgMagenta)
	yellow := color.New(color.FgYellow)

	c := &colorer{
		red:     red,
		green:   green,
		blue:    blue,
		cyan:    cyan,
		magenta: magenta,
		yellow:  yellow,
	}

	c.colors = []extColor{
		{color: red, escapeCode: getCode(red)},
		{color: green, escapeCode: getCode(green)},
		{color: blue, escapeCode: getCode(blue)},
		{color: cyan, escapeCode: getCode(cyan)},
		{color: magenta, escapeCode: getCode(magenta)},
		{color: yellow, escapeCode: getCode(yellow)},
	}
	switch strings.ToLower(when) {
	case "always":
		c.enable()
	case "never":
		c.disable()
	case "auto":
		c.auto()
	}
	return c
}

func (c *colorer) auto() {
	for _, v := range c.colors {
		if color.NoColor { // NoColor is global and set dynamically
			v.color.DisableColor()
			c.enabled = false
		} else {
			v.color.EnableColor()
			c.enabled = true
		}
	}
}

func (c *colorer) enable() {
	for _, v := range c.colors {
		v.color.EnableColor()
	}
	c.enabled = true
}

func (c *colorer) disable() {
	for _, v := range c.colors {
		v.color.DisableColor()
	}
	c.enabled = false
}

func (c colorer) port(a any) string {
	return c.yellow.Sprint(a)
}

func (c colorer) host(a any) string {
	return c.cyan.Sprint(a)
}

func (c colorer) identity(a any) string {
	return c.magenta.Sprint(a)
}

func (c colorer) verdictForwarded(a any) string {
	return c.green.Sprint(a)
}

func (c colorer) verdictDropped(a any) string {
	return c.red.Sprint(a)
}

func (c colorer) verdictAudit(a any) string {
	return c.yellow.Sprint(a)
}

func (c colorer) verdictTraced(a any) string {
	return c.yellow.Sprint(a)
}

func (c colorer) verdictTranslated(a any) string {
	return c.yellow.Sprint(a)
}

func (c colorer) authTestAlwaysFail(a any) string {
	return c.red.Sprint(a)
}

func (c colorer) authIsEnabled(a any) string {
	return c.green.Sprint(a)
}

// compute the list of unique ANSI escape sequences for this colorer.
func (c *colorer) sequences() []string {
	unique := make(map[string]struct{})
	for _, v := range c.colors {
		if c.enabled {
			unique[v.escapeCode] = struct{}{}
			unique["\x1b[0m"] = struct{}{} // reset code
		}
	}
	return slices.Collect(maps.Keys(unique))
}
