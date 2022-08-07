package version

import (
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type (
	numericDotSeparatedComparer struct{}
)

var _ comparer = numericDotSeparatedComparer{}

func (n numericDotSeparatedComparer) Compare(oldVersion, newVersion string) (compareResult, error) {
	oldV, err := n.parse(oldVersion)
	if err != nil {
		return compareResultInvalid, errors.Wrap(err, "parsing old version")
	}

	newV, err := n.parse(newVersion)
	if err != nil {
		return compareResultInvalid, errors.Wrap(err, "parsing old version")
	}

	getSeg := func(v []int, i int) int {
		if i >= len(v) {
			return 0
		}
		return v[i]
	}

	for i := 0; i < int(math.Max(float64(len(oldV)), float64(len(newV)))); i++ {
		switch {
		case getSeg(oldV, i) < getSeg(newV, i):
			return compareResultUpgrade, nil

		case getSeg(oldV, i) > getSeg(newV, i):
			return compareResultDowngrade, nil

		default:
			continue
		}
	}

	return compareResultEqual, nil
}

func (numericDotSeparatedComparer) IsPrerelease(string) (bool, error) {
	// Numeric Dot has no marker for Pre-Releases
	return false, nil
}

func (numericDotSeparatedComparer) parse(ver string) ([]int, error) {
	var out []int

	for _, seg := range strings.Split(ver, ".") {
		segI, err := strconv.Atoi(seg)
		if err != nil {
			return nil, errors.Wrap(err, "parsing segment")
		}
		out = append(out, segI)
	}

	return out, nil
}
