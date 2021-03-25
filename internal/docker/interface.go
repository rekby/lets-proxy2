package docker

import "context"

type DomainInfo struct {
	TargetAddress string
}

type Interface interface {
	GetTarget(ctx context.Context, domain string) (*DomainInfo, error)
}
