package tuikit

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FromIterm2 parses an iTerm2 colour scheme XML plist and returns a Theme.
// Recognised keys: "Foreground Color", "Background Color", "Ansi 1 Color"
// (red/negative), "Ansi 2 Color" (green/positive), "Ansi 4 Color"
// (blue/accent), "Cursor Color", "Selection Color".
func FromIterm2(xmlData []byte) (Theme, error) {
	colors, err := parsePlist(xmlData)
	if err != nil {
		return Theme{}, fmt.Errorf("FromIterm2: %w", err)
	}

	toHex := func(c plistColorComponents) lipgloss.Color {
		r := uint8(math.Round(c.r * 255))
		g := uint8(math.Round(c.g * 255))
		b := uint8(math.Round(c.b * 255))
		return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
	}

	t := DefaultTheme()
	if c, ok := colors["Foreground Color"]; ok {
		t.Text = toHex(c)
	}
	if c, ok := colors["Background Color"]; ok {
		t.TextInverse = toHex(c)
		t.Border = toHex(c)
	}
	if c, ok := colors["Ansi 1 Color"]; ok {
		t.Negative = toHex(c)
	}
	if c, ok := colors["Ansi 2 Color"]; ok {
		t.Positive = toHex(c)
	}
	if c, ok := colors["Ansi 4 Color"]; ok {
		t.Accent = toHex(c)
	}
	if c, ok := colors["Cursor Color"]; ok {
		t.Cursor = toHex(c)
	}
	if c, ok := colors["Selection Color"]; ok {
		t.Muted = toHex(c)
	}
	return t, nil
}

type plistColorComponents struct{ r, g, b float64 }

// parsePlist parses an iTerm2 plist XML into a map of colour name -> components.
func parsePlist(data []byte) (map[string]plistColorComponents, error) {
	result := map[string]plistColorComponents{}
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	var currentKey string
	depth := 0

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			if t.Name.Local == "key" && depth == 3 {
				var k string
				_ = dec.DecodeElement(&k, &t)
				depth--
				currentKey = k
			}
			if t.Name.Local == "dict" && depth == 3 && currentKey != "" {
				var comp plistColorComponents
				for {
					st, serr := dec.Token()
					if serr != nil {
						break
					}
					se, ok := st.(xml.StartElement)
					if !ok {
						if ee, ok2 := st.(xml.EndElement); ok2 && ee.Name.Local == "dict" {
							break
						}
						continue
					}
					if se.Name.Local == "key" {
						var compKey string
						_ = dec.DecodeElement(&compKey, &se)
						nst, nerr := dec.Token()
						if nerr != nil {
							break
						}
						if nse, ok3 := nst.(xml.StartElement); ok3 && nse.Name.Local == "real" {
							var v float64
							_ = dec.DecodeElement(&v, &nse)
							switch compKey {
							case "Red Component":
								comp.r = v
							case "Green Component":
								comp.g = v
							case "Blue Component":
								comp.b = v
							}
						}
					}
				}
				result[currentKey] = comp
				currentKey = ""
				depth--
			}
		case xml.EndElement:
			depth--
		}
	}
	return result, nil
}

// goghScheme is the JSON structure of a Gogh theme file.
type goghScheme struct {
	Name       string `json:"name"`
	Foreground string `json:"foreground"`
	Background string `json:"background"`
	Color02    string `json:"color_02"` // red / negative
	Color03    string `json:"color_03"` // green / positive
	Color05    string `json:"color_05"` // magenta / accent
	Color07    string `json:"color_07"` // cyan / cursor fallback
	Cursor     string `json:"cursor"`
}

// FromGogh parses a Gogh JSON colour scheme and returns a Theme.
func FromGogh(jsonData []byte) (Theme, error) {
	var s goghScheme
	if err := json.Unmarshal(jsonData, &s); err != nil {
		return Theme{}, fmt.Errorf("FromGogh: %w", err)
	}
	t := DefaultTheme()
	if s.Foreground != "" {
		t.Text = lipgloss.Color(s.Foreground)
	}
	if s.Background != "" {
		t.TextInverse = lipgloss.Color(s.Background)
	}
	if s.Color02 != "" {
		t.Negative = lipgloss.Color(s.Color02)
	}
	if s.Color03 != "" {
		t.Positive = lipgloss.Color(s.Color03)
	}
	if s.Color05 != "" {
		t.Accent = lipgloss.Color(s.Color05)
	}
	if s.Cursor != "" {
		t.Cursor = lipgloss.Color(s.Cursor)
	} else if s.Color07 != "" {
		t.Cursor = lipgloss.Color(s.Color07)
	}
	return t, nil
}

// FromAlacritty parses an Alacritty TOML colour scheme and returns a Theme.
// Only [colors.primary], [colors.normal], [colors.cursor], and [colors.bright]
// sections are used. No external TOML library is required.
func FromAlacritty(tomlData []byte) (Theme, error) {
	kvs, err := parseSimpleTOML(tomlData)
	if err != nil {
		return Theme{}, fmt.Errorf("FromAlacritty: %w", err)
	}
	t := DefaultTheme()
	if v, ok := kvs["colors.primary.foreground"]; ok {
		t.Text = lipgloss.Color(v)
	}
	if v, ok := kvs["colors.primary.background"]; ok {
		t.TextInverse = lipgloss.Color(v)
	}
	if v, ok := kvs["colors.normal.red"]; ok {
		t.Negative = lipgloss.Color(v)
	}
	if v, ok := kvs["colors.normal.green"]; ok {
		t.Positive = lipgloss.Color(v)
	}
	if v, ok := kvs["colors.normal.blue"]; ok {
		t.Accent = lipgloss.Color(v)
	}
	if v, ok := kvs["colors.cursor.cursor"]; ok {
		t.Cursor = lipgloss.Color(v)
	}
	if v, ok := kvs["colors.normal.black"]; ok {
		t.Border = lipgloss.Color(v)
	}
	if v, ok := kvs["colors.normal.yellow"]; ok {
		t.Flash = lipgloss.Color(v)
	}
	if v, ok := kvs["colors.bright.black"]; ok {
		t.Muted = lipgloss.Color(v)
	}
	return t, nil
}

// parseSimpleTOML handles flat [section] headers and key = "value" pairs.
// Returns a map of "section.key" -> unquoted value.
func parseSimpleTOML(data []byte) (map[string]string, error) {
	result := map[string]string{}
	var section string
	for _, rawLine := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.Trim(line, "[]")
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(strings.Trim(strings.TrimSpace(parts[1]), `"`))
		fullKey := key
		if section != "" {
			fullKey = section + "." + key
		}
		result[fullKey] = val
	}
	return result, nil
}
