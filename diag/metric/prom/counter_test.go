package pl_prom_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	pl_prom "github.com/agurinov/gopl.git/diag/metric/prom"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

func TestCounter(t *testing.T) {
	pl_testing.Init(t)

	_, err := pl_prom.NewCounter()
	require.NoError(t, err)

	t.Fail()
}
