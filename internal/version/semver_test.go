package version

import "testing"

func TestSemVerCompareFunc(t *testing.T) {
	comp := semVerComparer{}

	for _, tc := range []struct {
		v1, v2 string
		res    compareResult
	}{
		{"1.0.0", "2.0.0", compareResultUpgrade},
		{"2.0.0", "2.1.0", compareResultUpgrade},
		{"2.1.0", "2.1.1", compareResultUpgrade},
		{"2.1.1", "2.1.0", compareResultDowngrade},
		{"2.1.1", "2.1.1", compareResultEqual},
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
