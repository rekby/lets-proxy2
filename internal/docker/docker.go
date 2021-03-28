package docker

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/rekby/lets-proxy2/internal/domain"
	"github.com/rekby/lets-proxy2/internal/log"
	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const containerListLimit = math.MaxInt32

var ErrNotFound = errors.New("target not found")

type Config struct {
	DefaultHttpPort int
	LabelDomain     string
}

type dockerClientInterface interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
}

type Docker struct {
	client      dockerClientInterface
	domainLabel string
	portSuffix  string
}

func New(cfg Config) (*Docker, error) {
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		return nil, xerrors.Errorf("create docker client: %w", err)
	}
	return newDocker(cfg, dockerClient), nil
}

func newDocker(cfg Config, dockerClient dockerClientInterface) *Docker {
	return &Docker{
		client:      dockerClient,
		domainLabel: cfg.LabelDomain,
		portSuffix:  ":" + strconv.Itoa(cfg.DefaultHttpPort),
	}
}

func (d *Docker) GetTarget(ctx context.Context, dn domain.DomainName) (*DomainInfo, error) {
	logger := zc.L(ctx)
	list, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		Limit:   containerListLimit,
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "label", Value: d.domainLabel}),
	})

	log.DebugError(logger, err, "Got docker images list", zap.String("label", d.domainLabel), zap.Int("containers_count", len(list)))
	if err != nil {
		return nil, xerrors.Errorf("get containers list: %w", err)
	}

	container := findContainer(ctx, list, d.domainLabel, dn)
	if container == nil {
		return nil, ErrNotFound
	}

	if container.NetworkSettings == nil || len(container.NetworkSettings.Networks) == 0 {
		logger.Warn("Container found, but it has no IP address", zap.String("id", container.ID), zap.Strings("name", container.Names))
		return nil, xerrors.Errorf("container found, but it has no IP address")
	}

	if len(container.NetworkSettings.Networks) > 1 {
		logger.Warn("Container found, but it connected to many networks, can't determine right network for connect")
		return nil, xerrors.Errorf("container found, but it connected to many networks, can't determine right network for connect")
	}

	// it has exactly one network now
	for _, net := range container.NetworkSettings.Networks {
		return &DomainInfo{TargetAddress: net.IPAddress + d.portSuffix}, nil
	}
	logger.DPanic("Impossible situation for detect docker container")
	return nil, xerrors.Errorf("impossible situation for detect docker container")
}

func findContainer(ctx context.Context, containers []types.Container, label string, need domain.DomainName) *types.Container {
	logger := zc.L(ctx)
	for i := range containers {
		container := &containers[i]
		labelValue := container.Labels[label]
		domains := strings.Split(labelValue, ",")
		for _, containerDomain := range domains {
			dn, err := domain.NormalizeDomain(containerDomain)
			logger.Debug("Normalize container domain", zap.String("source domain", containerDomain), domain.LogDomain(dn), zap.Error(err))
			if err != nil {
				continue
			}
			if dn == need {
				logger.Debug("Found container", zap.String("id", container.ID), zap.Strings("names", container.Names))
				return container
			}
		}
	}
	logger.Debug("Doesn't find container for domain", domain.LogDomain(need))
	return nil
}
