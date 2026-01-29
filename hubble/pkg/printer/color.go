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
		{color: red,     escapeCode: "\x1b[31m"},
		{color: green,   escapeCode: "\x1b[32m"},
		{color: blue,    escapeCode: "\x1b[34m"},
		{color: cyan,    escapeCode: "\x1b[36m"},
		{color: magenta, escapeCode: "\x1b[35m"},
		{color: yellow,  escapeCode: "\x1b[33m"},
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
