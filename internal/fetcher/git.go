package fetcher

import (
	"context"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go_helpers/fieldcollection"
)

/*
 * @module git_tag
 * @module_desc Reads git tags (annotated and leightweight) from a remote repository and returns the newest one
 */

type (
	// GitTagFetcher implements the fetcher interface to monitor tags in a git repository
	GitTagFetcher struct{}
)

func init() { registerFetcher("git_tag", func() Fetcher { return &GitTagFetcher{} }) }

// FetchVersion retrieves the latest version for the catalog entry
func (g GitTagFetcher) FetchVersion(_ context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "opening in-mem repo")
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{attrs.MustString("remote", nil)},
	})
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "adding remote")
	}

	if err = repo.Fetch(&git.FetchOptions{
		Depth:      1,
		RefSpecs:   []config.RefSpec{"+refs/tags/*:refs/remotes/origin/tags/*"},
		RemoteName: "origin",
	}); err != nil {
		return "", time.Time{}, errors.Wrap(err, "fetching remote")
	}

	tags, err := repo.Tags()
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "listing tags")
	}

	var (
		latestTag     *plumbing.Reference
		latestTagTime time.Time
	)
	if err = tags.ForEach(func(ref *plumbing.Reference) error {
		tt, err := g.tagRefToTime(repo, ref)
		if err != nil {
			return errors.Wrap(err, "fetching time for tag")
		}

		if latestTag == nil || tt.After(latestTagTime) {
			latestTag = ref
			latestTagTime = tt
		}

		return nil
	}); err != nil {
		return "", time.Time{}, errors.Wrap(err, "iterating tags")
	}

	if latestTag == nil {
		return "", time.Time{}, ErrNoVersionFound
	}

	return latestTag.Name().Short(), latestTagTime, nil
}

// Links retrieves a collection of links for the fetcher
func (GitTagFetcher) Links(_ *fieldcollection.FieldCollection) []database.CatalogLink {
	return nil
}

// Validate validates the configuration given to the fetcher
func (GitTagFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
	// @attr remote required string "" Repository remote to fetch the tags from (should accept everything you can use in `git remote set-url` command)
	if v, err := attrs.String("remote"); err != nil || v == "" {
		return errors.New("remote is expected to be non-empty string")
	}

	return nil
}

func (GitTagFetcher) tagRefToTime(repo *git.Repository, tag *plumbing.Reference) (time.Time, error) {
	tagObj, err := repo.TagObject(tag.Hash())
	if err == nil {
		// Annotated tag: Take the time of the tag
		return tagObj.Tagger.When, nil
	}

	commitObj, err := repo.CommitObject(tag.Hash())
	if err == nil {
		// Lightweight tag: Take the time of the commit
		return commitObj.Committer.When, nil
	}

	return time.Time{}, errors.New("reference points neither to tag nor to commit")
}
