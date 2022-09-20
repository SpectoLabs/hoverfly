package gofakeit

import "math/rand"

// Color will generate a random color string
func Color() string { return color(globalFaker.Rand) }

// Color will generate a random color string
func (f *Faker) Color() string { return color(f.Rand) }

func color(r *rand.Rand) string { return getRandValue(r, []string{"color", "full"}) }

// SafeColor will generate a random safe color string
func SafeColor() string { return safeColor(globalFaker.Rand) }

// SafeColor will generate a random safe color string
func (f *Faker) SafeColor() string { return safeColor(f.Rand) }

func safeColor(r *rand.Rand) string { return getRandValue(r, []string{"color", "safe"}) }

// HexColor will generate a random hexadecimal color string
func HexColor() string { return hexColor(globalFaker.Rand) }

// HexColor will generate a random hexadecimal color string
func (f *Faker) HexColor() string { return hexColor(f.Rand) }

func hexColor(r *rand.Rand) string {
	color := make([]byte, 6)
	hashQuestion := []byte("?#")
	for i := 0; i < 6; i++ {
		color[i] = hashQuestion[r.Intn(2)]
	}

	return "#" + replaceWithHexLetters(r, replaceWithNumbers(r, string(color)))
}

// RGBColor will generate a random int slice color
func RGBColor() []int { return rgbColor(globalFaker.Rand) }

// RGBColor will generate a random int slice color
func (f *Faker) RGBColor() []int { return rgbColor(f.Rand) }

func rgbColor(r *rand.Rand) []int {
	return []int{randIntRange(r, 0, 255), randIntRange(r, 0, 255), randIntRange(r, 0, 255)}
}

func addColorLookup() {
	AddFuncLookup("color", Info{
		Display:     "Color",
		Category:    "color",
		Description: "Random color",
		Example:     "MediumOrchid",
		Output:      "string",
		Generate: func(r *rand.Rand, m *MapParams, info *Info) (interface{}, error) {
			return color(r), nil
		},
	})

	AddFuncLookup("safecolor", Info{
		Display:     "Safe Color",
		Category:    "color",
		Description: "Random safe color",
		Example:     "black",
		Output:      "string",
		Generate: func(r *rand.Rand, m *MapParams, info *Info) (interface{}, error) {
			return safeColor(r), nil
		},
	})

	AddFuncLookup("hexcolor", Info{
		Display:     "Hex Color",
		Category:    "color",
		Description: "Random hex color",
		Example:     "#a99fb4",
		Output:      "string",
		Generate: func(r *rand.Rand, m *MapParams, info *Info) (interface{}, error) {
			return hexColor(r), nil
		},
	})

	AddFuncLookup("rgbcolor", Info{
		Display:     "RGB Color",
		Category:    "color",
		Description: "Random rgb color",
		Example:     "[152 23 53]",
		Output:      "string",
		Generate: func(r *rand.Rand, m *MapParams, info *Info) (interface{}, error) {
			return rgbColor(r), nil
		},
	})
}
