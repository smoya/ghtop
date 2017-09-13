package testp

import (
	"context"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/smoya/ghtop/pkg/contributor"
	"github.com/smoya/ghtop/pkg/logx"
	"github.com/smoya/ghtop/pkg/server"
)

type topFeature struct {
	serverFeature
}

func topFeatureContext(s *godog.Suite) {
	f := &topFeature{}
	f.serverFeature.prepareScenario(s)

	s.Step(`^having these contributors "([^"]*)"$`, f.havingTheseContributors)
}

func (f *serverFeature) havingTheseContributors(stringContributors string) (err error) {
	return persistContributors(stringContributors, f.contributorRepo)
}

func persistContributors(stringContributors string, contributorRepo *contributorInMemoryRepository) error {
	contributorStr := strings.Split(stringContributors, ",")

	for _, s := range contributorStr {
		contributorRepo.data = append(contributorRepo.data, contributor.Contributor{
			Username:   s,
			VCS:        contributor.VCSGithub,
			ProfileURL: "https://github.com/" + s,
			AvatarURL:  "http://foo.bar",
		})
	}

	return nil
}

type contributorInMemoryRepository struct {
	data []contributor.Contributor
}

func (r *contributorInMemoryRepository) FindByLocation(ctx context.Context, location string, limit int, sortBy string) ([]contributor.Contributor, error) {
	if limit > len(r.data) {
		return r.data, nil
	}

	return r.data[:limit], nil
}

func (f *serverFeature) aServer() (err error) {
	f.server = server.NewServer(f.config, logx.NewDiscardAll(), f.contributorRepo)

	return
}
