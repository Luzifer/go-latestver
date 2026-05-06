package version

import (
	"errors"
	"fmt"

	"github.com/blang/semver/v4"
)

type (
	semVerComparer struct{}
)

var _ comparer = semVerComparer{}

func (semVerComparer) Compare(oldVersion, newVersion string) (compareResult, error) {
	oldS, err := semver.Make(oldVersion)
	if err != nil {
		return compareResultInvalid, fmt.Errorf("parsing old version: %w", err)
	}

	newS, err := semver.Make(newVersion)
	if err != nil {
		return compareResultInvalid, fmt.Errorf("parsing new version: %w", err)
	}

	switch oldS.Compare(newS) {
	case -1:
		// oldS < newS
		return compareResultUpgrade, nil

	case 0:
		// oldS == newS
		return compareResultEqual, nil

	case 1:
		// oldS > newS
		return compareResultDowngrade, nil

	default:
		// WTF, that does not exist according to lib docs
		return compareResultInvalid, errors.New("invalid compare result")
	}
}

func (semVerComparer) IsPrerelease(newVersion string) (bool, error) {
	newS, err := semver.Make(newVersion)
	if err != nil {
		return false, fmt.Errorf("parsing version: %w", err)
	}

	return newS.Pre != nil, nil
}
