package docker

import (
	"context"

	"github.com/rekby/lets-proxy2/internal/domain"
)

type DomainInfo struct {
	TargetAddress string
}

type Interface interface {
	GetTarget(ctx context.Context, domain domain.DomainName) (*DomainInfo, error)
}
