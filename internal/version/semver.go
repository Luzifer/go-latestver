package version

import (
	"github.com/blang/semver/v4"
	"github.com/pkg/errors"
)

type (
	semVerComparer struct{}
)

var _ comparer = semVerComparer{}

func (semVerComparer) Compare(oldVersion, newVersion string) (compareResult, error) {
	oldS, err := semver.Make(oldVersion)
	if err != nil {
		return compareResultInvalid, errors.Wrap(err, "parsing old version")
	}

	newS, err := semver.Make(newVersion)
	if err != nil {
		return compareResultInvalid, errors.Wrap(err, "parsing new version")
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
		return false, errors.Wrap(err, "parsing version")
	}

	return newS.Pre != nil, nil
}
