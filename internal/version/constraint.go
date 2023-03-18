package version

import "github.com/pkg/errors"

type (
	// Constraint document how a version update should be handled
	Constraint struct {
		AllowDowngrade  bool `yaml:"allow_downgrade"`
		AllowPrerelease bool `yaml:"allow_prerelease"`

		Type string `yaml:"type"`
	}

	comparer interface {
		Compare(oldVersion, newVersion string) (compareResult, error)
		IsPrerelease(newVersion string) (bool, error)
	}

	compareResult uint
)

const (
	compareResultInvalid compareResult = iota
	compareResultEqual
	compareResultDowngrade
	compareResultUpgrade
)

// ShouldApply checks whether a new version should overwrite the old
// one given the parameters inside the Constraint
func (c Constraint) ShouldApply(oldVersion, newVersion string) (bool, error) {
	if oldVersion == "" && newVersion != "" {
		// The old version does not exist, the new one does, update it!
		return true, nil
	}

	comp := c.getComparer()
	if comp == nil {
		return false, errors.New("invalid version type specified")
	}

	// Compare versions and check for UpgradeOnly flag
	compResult, err := comp.Compare(oldVersion, newVersion)
	if err != nil {
		return false, errors.Wrap(err, "comparing versions")
	}

	if !c.AllowDowngrade && compResult != compareResultUpgrade {
		return false, nil
	}

	// check for forbidden pre-releases
	isPreR, err := comp.IsPrerelease(newVersion)
	if err != nil {
		return false, errors.Wrap(err, "checking pre-release")
	}

	if !c.AllowPrerelease && isPreR {
		return false, nil
	}

	return true, nil
}

func (c Constraint) getComparer() comparer {
	switch c.Type {
	case "numeric_dot":
		return numericDotSeparatedComparer{}

	case "semver":
		return semVerComparer{}

	default:
		return nil
	}
}
