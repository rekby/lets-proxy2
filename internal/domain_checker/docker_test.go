package domain_checker

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/domain"
	"github.com/rekby/lets-proxy2/internal/th"

	"github.com/rekby/lets-proxy2/internal/docker"
)

//go:generate minimock -i internalDocker -o internal_docker_mock_test.go -g

type internalDocker interface {
	docker.Interface
}

func TestDocker(t *testing.T) {
	const domainName = "domain"

	td := testdeep.NewT(t)

	ctx, flush := th.TestContext(td)
	defer flush()

	mc := minimock.NewController(td)

	table := []struct {
		Name               string
		Domain             string
		DockerTargetResult *docker.DomainInfo
		DockerTargetError  error
		Result             bool
		ErrorString        string
	}{
		{"ok", "domain", &docker.DomainInfo{TargetAddress: "asd"}, nil, true, ""},
		{"not found", "domain", nil, docker.ErrNotFound, false, ""},
		{"bad domain", "domain:adsfasdf:dfasdf", nil, nil, false, "normalize domain"},
		{"docker-error", "domain", nil, errors.New("test-err"), false, "test-err"},
	}

	for _, test := range table {
		dockerMock := NewInternalDockerMock(mc)
		if test.DockerTargetResult != nil || test.DockerTargetError != nil {
			testDomain, err := domain.NormalizeDomain(test.Domain)
			td.CmpNoError(err)
			dockerMock.GetTargetMock.Expect(ctx, testDomain).Return(test.DockerTargetResult, test.DockerTargetError)
		}

		checker := NewDockerChecker(dockerMock)
		res, err := checker.IsDomainAllowed(ctx, test.Domain)
		td.Cmp(res, test.Result)
		if diff := th.ErrorSubstringCmp(err, test.ErrorString); diff != "" {
			td.Errorf("Test name: '%v', diff error: %v", test.Name, diff)
		}
	}
}
