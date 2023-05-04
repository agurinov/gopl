package pl_testing

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

var (
	pool         *dockertest.Pool
	poolMaxWait  = 10 * time.Second
	errPoolInit  error
	poolInitOnce sync.Once
)

func (tc TestCase) Pool(t *testing.T) *dockertest.Pool {
	t.Helper()

	poolInitOnce.Do(func() {
		pool, errPoolInit = dockertest.NewPool("")

		if errPoolInit != nil {
			return
		}

		pool.MaxWait = poolMaxWait
	})

	require.NoError(t, errPoolInit)
	require.NotNil(t, pool)
	require.NoError(t, pool.Client.Ping())

	return pool
}

func (tc TestCase) network(t *testing.T, networkName string) *docker.Network {
	t.Helper()

	var (
		pool    = tc.Pool(t)
		network *docker.Network
	)

	networks, err := pool.NetworksByName(networkName)
	require.NoError(t, err)

	if len(networks) != 0 {
		network = networks[0].Network
	} else {
		network, err = pool.Client.CreateNetwork(docker.CreateNetworkOptions{
			Name: networkName,
		})
	}

	switch tc.Debug {
	case true:
		if errors.Is(err, docker.ErrNetworkAlreadyExists) {
			break
		}
	case false:
		require.NoError(t, err)
		require.NotNil(t, network)
		t.Cleanup(func() { pool.Client.RemoveNetwork(network.ID) }) //nolint:errcheck
	}

	return network
}

func (tc TestCase) container(t *testing.T, opts *dockertest.RunOptions) *dockertest.Resource {
	t.Helper()

	var (
		pool           = tc.Pool(t)
		container, err = pool.RunWithOptions(opts)
	)

	switch tc.Debug {
	case true:
		if errors.Is(err, docker.ErrContainerAlreadyExists) {
			break
		}
	case false:
		require.NoError(t, err)
		require.NotNil(t, container)
		t.Cleanup(func() { pool.Purge(container) }) //nolint:errcheck
	}

	return container
}
