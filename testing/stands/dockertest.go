package stands

// TODO(a.gurinov): Cleanup based on debug (pass from testcase)

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

var (
	pool         *dockertest.Pool
	errPoolInit  error
	poolInitOnce sync.Once
	poolMaxWait  = 10 * time.Second
)

func Pool(t *testing.T) *dockertest.Pool {
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

func network(t *testing.T) *docker.Network {
	t.Helper()

	var (
		p           = Pool(t)
		networkName = "gopl_" + hash(t)
		network     *docker.Network
	)

	networks, err := p.NetworksByName(networkName)
	require.NoError(t, err)

	if len(networks) > 0 {
		network = networks[0].Network
	} else {
		network, err = p.Client.CreateNetwork(docker.CreateNetworkOptions{
			Name: networkName,
		})
	}

	require.NoError(t, err)
	require.NotNil(t, network)

	t.Cleanup(func() {
		p.RemoveNetwork(&dockertest.Network{ //nolint:errcheck
			Network: network,
		})
	})

	return network
}

func container(
	t *testing.T,
	opts *dockertest.RunOptions,
) (
	*dockertest.Resource,
	bool,
) {
	t.Helper()

	var (
		p          = Pool(t)
		hostConfig = func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		}
		container, err = p.RunWithOptions(opts, hostConfig)
		created        = true
	)

	if errors.Is(err, docker.ErrContainerAlreadyExists) {
		c, ok := p.ContainerByName(opts.Name)
		require.True(t, ok)

		container = c
		err = nil
		created = false
	}

	require.NoError(t, err)
	require.NotNil(t, container)
	require.NotNil(t, container.Container)
	require.True(t, container.Container.State.Running)

	t.Cleanup(func() {
		if purgeErr := p.Purge(container); purgeErr != nil {
			panic(purgeErr)
		}
	})

	return container, created
}

func containerExec(
	t *testing.T,
	container *dockertest.Resource,
	stdin io.Reader,
	cmd ...string,
) {
	t.Helper()

	require.NotNil(t, container)
	require.NotEmpty(t, cmd)

	var (
		p           = Pool(t)
		execOptions = dockertest.ExecOptions{StdIn: stdin}
	)

	require.NoError(t,
		p.Retry(func() error {
			exitCode, err := container.Exec(cmd, execOptions)
			if err != nil {
				return err
			}

			if exitCode != 0 {
				return fmt.Errorf(
					"container %q exec: unexpected exit code: %d",
					container.Container.Name,
					exitCode,
				)
			}

			return nil
		}),
	)
}
