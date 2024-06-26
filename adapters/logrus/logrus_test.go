package logrus

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axiomhq/axiom-go/axiom"
	"github.com/axiomhq/axiom-go/internal/test/adapters"
	"github.com/axiomhq/axiom-go/internal/test/testhelper"
)

// TestNew makes sure New() picks up the "AXIOM_DATASET" environment variable.
func TestNew(t *testing.T) {
	testhelper.SafeClearEnv(t)

	t.Setenv("AXIOM_TOKEN", "xaat-test")
	t.Setenv("AXIOM_ORG_ID", "123")

	handler, err := New()
	require.ErrorIs(t, err, ErrMissingDatasetName)
	require.Nil(t, handler)

	t.Setenv("AXIOM_DATASET", "test")

	handler, err = New()
	require.NoError(t, err)
	require.NotNil(t, handler)
	handler.Close()

	assert.Equal(t, "test", handler.datasetName)
}

func TestHook(t *testing.T) {
	now := time.Now()

	exp := fmt.Sprintf(`{"_time":"%s","severity":"info","key":"value","message":"my message"}`,
		now.Format(time.RFC3339Nano))

	var hasRun uint64
	hf := func(w http.ResponseWriter, r *http.Request) {
		zsr, err := zstd.NewReader(r.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zsr)
		assert.NoError(t, err)

		assert.JSONEq(t, exp, string(b))

		atomic.AddUint64(&hasRun, 1)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}

	logger, closeHook := adapters.Setup(t, hf, setup(t))

	logger.
		WithTime(now).
		WithField("key", "value").
		Info("my message")

	closeHook()

	assert.EqualValues(t, 1, atomic.LoadUint64(&hasRun))
}

func TestHook_NoPanicAfterClose(t *testing.T) {
	now := time.Now()

	exp := fmt.Sprintf(`{"_time":"%s","severity":"info","key":"value","message":"my message"}`,
		now.Format(time.RFC3339Nano))

	var lines uint64
	hf := func(w http.ResponseWriter, r *http.Request) {
		zsr, err := zstd.NewReader(r.Body)
		require.NoError(t, err)

		s := bufio.NewScanner(zsr)
		for s.Scan() {
			assert.JSONEq(t, exp, s.Text())
			atomic.AddUint64(&lines, 1)
		}
		assert.NoError(t, s.Err())

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}

	logger, closeHook := adapters.Setup(t, hf, setup(t))

	logger.
		WithTime(now).
		WithField("key", "value").
		Info("my message")

	closeHook()

	// This should be a no-op.
	logger.
		WithTime(now).
		WithField("key", "value").
		Info("my message")

	assert.EqualValues(t, 1, atomic.LoadUint64(&lines))
}

func setup(t *testing.T) func(dataset string, client *axiom.Client) (*logrus.Logger, func()) {
	return func(dataset string, client *axiom.Client) (*logrus.Logger, func()) {
		t.Helper()

		hook, err := New(
			SetClient(client),
			SetDataset(dataset),
		)
		require.NoError(t, err)
		t.Cleanup(hook.Close)

		logger := logrus.New()
		logger.AddHook(hook)

		// We don't want output in tests.
		logger.Out = io.Discard

		return logger, hook.Close
	}
}
