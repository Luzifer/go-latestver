package version

import "testing"

func TestNumericDotSeparatedCompareFunc(t *testing.T) {
	comp := numericDotSeparatedComparer{}

	for _, tc := range []struct {
		v1, v2 string
		res    compareResult
	}{
		{"2.1.0", "2.1.1", compareResultUpgrade},
		{"2.1.1", "2.1.0", compareResultDowngrade},
		{"2.1.1", "2.1.1", compareResultEqual},
		{"103.0.5060.134", "104.0.5112.79", compareResultUpgrade},
		{"2022.7.7", "2022.8.0", compareResultUpgrade},
	} {
		res, err := comp.Compare(tc.v1, tc.v2)
		if err != nil {
			t.Errorf("Comparing %q to %q: %s", tc.v1, tc.v2, err)
		}

		if res != tc.res {
			t.Errorf("Comparing %q to %q: expected %v, got %v", tc.v1, tc.v2, tc.res, res)
		}
	}
}
