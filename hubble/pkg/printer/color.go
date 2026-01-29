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

type colorer struct {
	colors  []*color.Color
	seqName string
	red     sprinter
	green   sprinter
	blue    sprinter
	cyan    sprinter
	magenta sprinter
	yellow  sprinter
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

	c.colors = []*color.Color{
		red, green, blue,
		cyan, magenta, yellow,
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

// globally cached map of possible sequences, introduced for performance reasons, to avoid recalculating it for each colorer instance
var SEQUENCES = map[string][]string{
	"auto":   newColorer("auto").calculateSequences(),
	"never":  newColorer("never").calculateSequences(),
	"always": newColorer("always").calculateSequences(),
}

func (c *colorer) auto() {
	for _, v := range c.colors {
		if color.NoColor { // NoColor is global and set dynamically
			v.DisableColor()
		} else {
			v.EnableColor()
		}
	}
	c.seqName = "auto"
}

func (c *colorer) enable() {
	for _, v := range c.colors {
		v.EnableColor()
	}
	c.seqName = "always"
}

func (c *colorer) disable() {
	for _, v := range c.colors {
		v.DisableColor()
	}
	c.seqName = "never"
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
func (c *colorer) calculateSequences() []string {
	unique := make(map[string]struct{})
	for _, v := range c.colors {
		seq := v.Sprint("|")
		split := strings.Split(seq, "|")
		if len(split) != 2 {
			// should never happen
			continue
		}
		unique[split[0]] = struct{}{}
		unique[split[1]] = struct{}{}
	}
	return slices.Collect(maps.Keys(unique))
}

func (c *colorer) sequences() []string {
	return SEQUENCES[c.seqName]
}
