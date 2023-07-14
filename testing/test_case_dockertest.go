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

func (tc TestCase) network(
	t *testing.T,
	networkName string,
) *docker.Network {
	t.Helper()

	var (
		pool    = tc.Pool(t)
		network *docker.Network
	)

	networks, err := pool.NetworksByName(networkName)
	require.NoError(t, err)

	if len(networks) > 0 {
		network = networks[0].Network
	} else {
		network, err = pool.Client.CreateNetwork(docker.CreateNetworkOptions{
			Name: networkName,
		})
	}

	require.NoError(t, err)
	require.NotNil(t, network)

	if !tc.Debug {
		t.Cleanup(func() {
			if err := pool.Client.RemoveNetwork(network.ID); err != nil {
				panic(err)
			}
		})
	}

	return network
}

func (tc TestCase) container(
	t *testing.T,
	opts *dockertest.RunOptions,
) (
	*dockertest.Resource,
	bool,
) {
	t.Helper()

	var (
		pool           = tc.Pool(t)
		container, err = pool.RunWithOptions(opts)
		created        = true
	)

	if errors.Is(err, docker.ErrContainerAlreadyExists) {
		c, ok := pool.ContainerByName(opts.Name)
		require.True(t, ok)

		container = c
		err = nil
		created = false
	}

	require.NoError(t, err)
	require.NotNil(t, container)
	require.NotNil(t, container.Container)
	require.True(t, container.Container.State.Running)

	if !tc.Debug {
		t.Cleanup(func() {
			if err := pool.Purge(container); err != nil {
				panic(err)
			}
		})
	}

	return container, created
}
