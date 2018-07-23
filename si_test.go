package humanize

import (
	"math"
	"testing"
)

func TestSI(t *testing.T) {
	tests := []struct {
		name      string
		num       float64
		unit      string
		formatted string
	}{
		{"e-24", 1e-24, "F", "1 yF"},
		{"e-21", 1e-21, "F", "1 zF"},
		{"e-18", 1e-18, "F", "1 aF"},
		{"e-15", 1e-15, "F", "1 fF"},
		{"e-12", 1e-12, "F", "1 pF"},
		{"e-12", 2.2345e-12, "F", "2.2345 pF"},
		{"e-12", 2.23e-12, "F", "2.23 pF"},
		{"e-11", 2.23e-11, "F", "22.3 pF"},
		{"e-10", 2.2e-10, "F", "220 pF"},
		{"e-9", 2.2e-9, "F", "2.2 nF"},
		{"e-8", 2.2e-8, "F", "22 nF"},
		{"e-7", 2.2e-7, "F", "220 nF"},
		{"e-6", 2.2e-6, "F", "2.2 µF"},
		{"e-6", 1e-6, "F", "1 µF"},
		{"e-5", 2.2e-5, "F", "22 µF"},
		{"e-4", 2.2e-4, "F", "220 µF"},
		{"e-3", 2.2e-3, "F", "2.2 mF"},
		{"e-2", 2.2e-2, "F", "22 mF"},
		{"e-1", 2.2e-1, "F", "220 mF"},
		{"e+0", 2.2e-0, "F", "2.2 F"},
		{"e+0", 2.2, "F", "2.2 F"},
		{"e+1", 2.2e+1, "F", "22 F"},
		{"0", 0, "F", "0 F"},
		{"e+1", 22, "F", "22 F"},
		{"e+2", 2.2e+2, "F", "220 F"},
		{"e+2", 220, "F", "220 F"},
		{"e+3", 2.2e+3, "F", "2.2 kF"},
		{"e+3", 2200, "F", "2.2 kF"},
		{"e+4", 2.2e+4, "F", "22 kF"},
		{"e+4", 22000, "F", "22 kF"},
		{"e+5", 2.2e+5, "F", "220 kF"},
		{"e+6", 2.2e+6, "F", "2.2 MF"},
		{"e+6", 1e+6, "F", "1 MF"},
		{"e+7", 2.2e+7, "F", "22 MF"},
		{"e+8", 2.2e+8, "F", "220 MF"},
		{"e+9", 2.2e+9, "F", "2.2 GF"},
		{"e+10", 2.2e+10, "F", "22 GF"},
		{"e+11", 2.2e+11, "F", "220 GF"},
		{"e+12", 2.2e+12, "F", "2.2 TF"},
		{"e+15", 2.2e+15, "F", "2.2 PF"},
		{"e+18", 2.2e+18, "F", "2.2 EF"},
		{"e+21", 2.2e+21, "F", "2.2 ZF"},
		{"e+24", 2.2e+24, "F", "2.2 YF"},

		// special case
		{"1F", 1000 * 1000, "F", "1 MF"},
		{"1F", 1e6, "F", "1 MF"},
		{"5%", .05, "%", "0.05 %"},

		// negative number
		{"-100 F", -100, "F", "-100 F"},
	}

	for _, test := range tests {
		got := SI(test.num, test.unit)
		if got != test.formatted {
			t.Errorf("On %v (%v), got %v, wanted %v",
				test.name, test.num, got, test.formatted)
		}

		gotf, gotu, err := ParseSI(test.formatted)
		if err != nil {
			t.Errorf("Error parsing %v (%v): %v", test.name, test.formatted, err)
			continue
		}

		if math.Abs(1-(gotf/test.num)) > 0.01 {
			t.Errorf("On %v (%v), got %v, wanted %v (±%v)",
				test.name, test.formatted, gotf, test.num,
				math.Abs(1-(gotf/test.num)))
		}
		if gotu != test.unit {
			t.Errorf("On %v (%v), expected unit %v, got %v",
				test.name, test.formatted, test.unit, gotu)
		}
	}

}

func TestParseSI(t *testing.T) {
	tests := []struct {
		input    string
		num      float64
		unit     string
		hasError bool
	}{
		{"1.21kW", 1210.0, "W", false},
		{"1.21 kW", 1210.0, "W", false},
		{"1.21 KW", 1210.0, "W", false},
		{"1.21µW", 1.21e-6, "W", false},
		{"1.21uW", 1.21e-6, "W", false},
		{"1.21 uW", 1.21e-6, "W", false},
		{"6.8e-3 W", 6.8e-3, "W", false},
		{"6.8E-3 W", 6.8e-3, "W", false},
		{"6.8e-3W", 6.8e-3, "W", false},
		{"-6.8e-3W", -6.8e-3, "W", false},
		{"-6.8e-3 W", -6.8e-3, "W", false},
		{"-6.8e-3 kW", -6.8, "W", false},
		{"1000", 1000, "", false},
		{"1000W", 1000, "W", false},
		{"0W", 0, "W", false},
		{"0.6 pF", 6e-13, "F", false},
		{"100pF", 1e-10, "F", false},
		{"x1.21JW", 0, "", true},
	}

	for _, test := range tests {
		gotf, gotu, err := ParseSI(test.input)
		if test.hasError && err == nil {
			t.Errorf("Expected error on %s, got %v %v", test.input, gotf, gotu)
		}
		if !test.hasError && err != nil {
			t.Errorf("Expected no error on %s, got error %v", test.input, err)
		}
		if math.Abs(1-(gotf/test.num)) > 0.01 {
			t.Errorf("On %v got %v, wanted %v (±%v)",
				test.input, gotf, test.num,
				math.Abs(1-(gotf/test.num)))
		}
		if gotu != test.unit {
			t.Errorf("On %v expected %v got %v", test.input, test.unit, gotu)
		}
	}
}

func TestSIWithDigits(t *testing.T) {
	tests := []struct {
		name      string
		num       float64
		digits    int
		formatted string
	}{
		{"e-12", 2.234e-12, 0, "2 pF"},
		{"e-12", 2.234e-12, 1, "2.2 pF"},
		{"e-12", 2.234e-12, 2, "2.23 pF"},
		{"e-12", 2.234e-12, 3, "2.234 pF"},
		{"e-12", 2.234e-12, 4, "2.234 pF"},
	}

	for _, test := range tests {
		got := SIWithDigits(test.num, test.digits, "F")
		if got != test.formatted {
			t.Errorf("On %v (%v), got %v, wanted %v",
				test.name, test.num, got, test.formatted)
		}
	}
}

func BenchmarkParseSI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseSI("2.2346ZB")
	}
}
