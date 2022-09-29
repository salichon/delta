package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	componentMake = iota
	componentModel
	componentType
	componentNumber
	componentSubsource
	componentDip
	componentAzimuth
	componentTypes
	componentResponse
	componentLast
)

type Component struct {
	Make      string
	Model     string
	Type      string
	Number    int
	Subsource string
	Dip       float64
	Azimuth   float64
	Types     string
	Response  string

	number  string
	dip     string
	azimuth string
}

// Less compares Component structs suitable for sorting.
func (c Component) Less(comp Component) bool {

	switch {
	case strings.ToLower(c.Make) < strings.ToLower(comp.Make):
		return true
	case strings.ToLower(c.Make) > strings.ToLower(comp.Make):
		return false
	case strings.ToLower(c.Model) < strings.ToLower(comp.Model):
		return true
	case strings.ToLower(c.Model) > strings.ToLower(comp.Model):
		return false
	case c.Number < comp.Number:
		return true
	case c.Number > comp.Number:
		return false
	default:
		return true
	}
}

type ComponentList []Component

func (s ComponentList) Len() int           { return len(s) }
func (s ComponentList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ComponentList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s ComponentList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Type",
		"Number",
		"Subsource",
		"Dip",
		"Azimuth",
		"Types",
		"Response",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.number),
			strings.TrimSpace(v.Subsource),
			strings.TrimSpace(v.dip),
			strings.TrimSpace(v.azimuth),
			strings.TrimSpace(v.Types),
			strings.TrimSpace(v.Response),
		})
	}
	return data
}
func (s *ComponentList) decode(data [][]string) error {
	var components []Component

	if !(len(data) > 1) {
		return nil
	}

	for _, d := range data[1:] {
		if len(d) != componentLast {
			return fmt.Errorf("incorrect pin of installed component fields")
		}

		var number int
		if s := strings.TrimSpace(d[componentNumber]); s != "" {
			v, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			number = v
		}

		var dip float64
		if s := strings.TrimSpace(d[componentDip]); s != "" {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			dip = v
		}

		var azimuth float64
		if s := strings.TrimSpace(d[componentAzimuth]); s != "" {
			v, err := strconv.ParseFloat(d[componentAzimuth], 64)
			if err != nil {
				return err
			}
			azimuth = v
		}

		components = append(components, Component{
			Make:      strings.TrimSpace(d[componentMake]),
			Model:     strings.TrimSpace(d[componentModel]),
			Type:      strings.TrimSpace(d[componentType]),
			Number:    number,
			Subsource: strings.TrimSpace(d[componentSubsource]),
			Dip:       dip,
			Azimuth:   azimuth,
			Types:     strings.TrimSpace(d[componentTypes]),
			Response:  strings.TrimSpace(d[componentResponse]),

			number:  strings.TrimSpace(d[componentNumber]),
			dip:     strings.TrimSpace(d[componentDip]),
			azimuth: strings.TrimSpace(d[componentAzimuth]),
		})
	}

	*s = ComponentList(components)

	return nil
}

func LoadComponents(path string) ([]Component, error) {
	var s []Component

	if err := LoadList(path, (*ComponentList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(ComponentList(s))

	return s, nil
}