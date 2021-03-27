//nolint:golint
package domain_checker

import (
	"context"

	"github.com/rekby/lets-proxy2/internal/domain"
	"golang.org/x/xerrors"

	"github.com/rekby/lets-proxy2/internal/docker"
)

type DockerChecker struct {
	client docker.Interface
}

func (d DockerChecker) IsDomainAllowed(ctx context.Context, dn string) (bool, error) {
	domainName, err := domain.NormalizeDomain(dn)
	if err != nil {
		return false, xerrors.Errorf("normalize domain in docker domain checker: :w", err)
	}

	_, err = d.client.GetTarget(ctx, domainName)
	if err == nil {
		return true, nil
	}
	if xerrors.Is(err, docker.ErrNotFound) {
		return false, nil
	}
	return false, xerrors.Errorf("get target from docker: %w", err)
}

func NewDockerChecker(client docker.Interface) DockerChecker {
	return DockerChecker{
		client: client,
	}
}
