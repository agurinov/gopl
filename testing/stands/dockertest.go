package stands

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/backoff"
	"github.com/agurinov/gopl/backoff/strategies"
)

var (
	pool               *dockertest.Pool
	errPoolInit        error
	poolInitOnce       sync.Once
	defaultPoolMaxWait = 30 * time.Second
	poolMaxWait        time.Duration
	poolNoCleanup      bool
	execVerbose        bool
)

func init() {
	// TODO(a.gurinov): flag.Parse() panics here
	flag.DurationVar(&poolMaxWait, "dockertest-pool-max-wait", defaultPoolMaxWait, "")
	flag.BoolVar(&poolNoCleanup, "dockertest-pool-no-cleanup", false, "")
	flag.BoolVar(&execVerbose, "dockertest-exec-verbose", false, "")
}

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
	require.NotNil(t, pool.Client)
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
		if poolNoCleanup {
			return
		}

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
		if poolNoCleanup {
			return
		}

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
	require.NotEmpty(t, container.Container.Name)

	b, err := backoff.New(
		backoff.WithExponentialStrategy(
			strategies.WithMaxDelay(poolMaxWait),
		),
	)
	require.NoError(t, err)
	require.NotNil(t, b)

	execOp := func() error {
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)

		defer func() {
			if !execVerbose {
				return
			}

			t.Logf(
				"stdout from %q exec(%s): %s",
				container.Container.Name,
				cmd,
				stdout.String(),
			)

			t.Logf(
				"stderr from %q exec(%s): %s",
				container.Container.Name,
				cmd,
				stderr.String(),
			)
		}()

		exitCode, execErr := container.Exec(cmd, dockertest.ExecOptions{
			StdIn:  stdin,
			StdOut: &stdout,
			StdErr: &stderr,
		})
		if execErr != nil {
			return execErr
		}

		if exitCode != 0 {
			return fmt.Errorf(
				"container %q exec: unexpected exit code: %d",
				container.Container.Name,
				exitCode,
			)
		}

		return nil
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		poolMaxWait,
	)
	t.Cleanup(cancel)

	for {
		execErr := execOp()
		if execErr == nil {
			break
		}

		if _, backoffErr := b.Wait(ctx); backoffErr != nil {
			require.NoError(t, execErr)
		}
	}
}
