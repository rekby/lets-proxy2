package docker

import (
	"errors"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep"

	"github.com/docker/docker/api/types/filters"
	"github.com/gojuno/minimock/v3"

	"github.com/docker/docker/api/types/network"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"

	"github.com/rekby/lets-proxy2/internal/domain"

	"github.com/docker/docker/api/types"
	"github.com/rekby/lets-proxy2/internal/th"
)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/docker.dockerClientInterface -o ./docker_client_interface_mock_test.go -g

func TestFindContainer(t *testing.T) {
	td := testdeep.NewT(t)

	const label = "label"
	ctx, flush := th.TestContext(td)
	defer flush()

	table := []struct {
		Name       string
		Containers []types.Container
		Domain     string
		ResultID   string
	}{
		{"one", []types.Container{{ID: "1", Labels: map[string]string{label: "domain"}}}, "domain", "1"},
		{"two", []types.Container{{ID: "1", Labels: map[string]string{label: "dom1,dom2"}}}, "dom2", "1"},
		{"second", []types.Container{
			{ID: "1", Labels: map[string]string{label: "domain"}},
			{ID: "2", Labels: map[string]string{label: "domain2"}},
		}, "domain2", "2"},
		{"first", []types.Container{
			{ID: "2", Labels: map[string]string{label: "domain2"}},
			{ID: "1", Labels: map[string]string{label: "domain"}},
		}, "domain2", "2"},
		{"nil", []types.Container{
			{ID: "2", Labels: map[string]string{label: "domain2"}},
			{ID: "1", Labels: map[string]string{label: "domain"}},
		}, "domain3", ""},
		{"domain with error", []types.Container{
			{ID: "1", Labels: map[string]string{label: "domain2:asd:fdsad"}},
			{ID: "2", Labels: map[string]string{label: "domain3"}},
		}, "domain3", "2"},
	}

	for _, item := range table {
		td.NotEmpty(item.Name)
		ctx := zc.WithLogger(ctx, zc.L(ctx).With(zap.String("test_name", item.Name)))

		// check test case
		ids := make(map[string]bool)
		for _, c := range item.Containers {
			td.NotEmpty(c.ID)
			if ids[c.ID] {
				td.Fatalf("Duplicate container id. testName: %v, ID: '%v'", item.Name, c.ID)
			}
			ids[c.ID] = true
		}

		needDomain, err := domain.NormalizeDomain(item.Domain)
		td.CmpNoError(err)

		resContainer := findContainer(ctx, item.Containers, label, needDomain)
		var resID string
		if resContainer != nil {
			resID = resContainer.ID
		}
		td.Cmp(resID, item.ResultID)
	}
}

func TestGetTarget(t *testing.T) {
	td := testdeep.NewT(t)

	const label = "label"
	ctx, flush := th.TestContext(td)
	defer flush()

	dName, err := domain.NormalizeDomain("domain")
	if err != nil {
		t.Fatal(err)
	}

	mc := minimock.NewController(td)

	cont := func(ips ...string) []types.Container {
		c := types.Container{}
		c.Labels = map[string]string{label: dName.String()}
		c.NetworkSettings = &types.SummaryNetworkSettings{}
		c.NetworkSettings.Networks = map[string]*network.EndpointSettings{}
		for _, ip := range ips {
			c.NetworkSettings.Networks["net-"+ip] = &network.EndpointSettings{IPAddress: ip}
		}
		return []types.Container{c}
	}

	table := []struct {
		Name              string
		ContainersForFind []types.Container
		DockerError       error
		Target            string
		TargetError       string
	}{
		{
			"ok",
			cont("1.2.3.4"),
			nil,
			"1.2.3.4:80",
			"",
		},
		{
			"docker-error",
			nil,
			errors.New("test-error"),
			"",
			"test-error",
		},
		{
			"container not found",
			nil,
			nil,
			"",
			ErrNotFound.Error(),
		},
		{
			"container without networks",
			cont(),
			nil,
			"",
			"no IP address",
		},
		{
			"container with many networks",
			cont("1.2.3.4", "2.3.4.5"),
			nil,
			"",
			"connected to many networks",
		},
	}

	for _, test := range table {
		td.NotEmpty(test.Name)
		dockerClientMock := NewDockerClientInterfaceMock(mc)
		dockerClientMock.ContainerListMock.Expect(ctx, types.ContainerListOptions{
			Limit:   containerListLimit,
			Filters: filters.NewArgs(filters.KeyValuePair{Key: "label", Value: label}),
		}).Return(test.ContainersForFind, test.DockerError)

		d := newDocker(Config{
			LabelDomain:     label,
			DefaultHttpPort: 80,
		}, dockerClientMock)

		target, err := d.GetTarget(ctx, dName)

		var errString string
		if err != nil {
			errString = err.Error()
		}
		if test.TargetError == "" && errString != "" || test.TargetError != "" && !strings.Contains(errString, test.TargetError) {
			td.Errorf("Got: '%v' Expect: '%v'", errString, test.DockerError)
		}

		var targetAddress string
		if target != nil {
			targetAddress = target.TargetAddress
		}
		td.Cmp(targetAddress, test.Target)
	}
}
