package sock

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	testHomeDir, err := os.MkdirTemp("", "controller_test")
	if err != nil {
		panic(err)
	}
	os.Setenv("HOME", testHomeDir)
	code := m.Run()
	os.RemoveAll(testHomeDir)
	os.Exit(code)
}

func TestStartAndShutdownServer(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_server_start_shutdown")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	unixServer, err := NewServer(
		&Config{
			Addr: tmpFile.Name(),
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
		})
	require.NoError(t, err)

	client := Client{Addr: tmpFile.Name()}
	listen := make(chan error)
	go func() {
		for range listen {
		}
	}()

	go func() {
		err = unixServer.Serve(listen)
		assert.True(t, errors.Is(ErrServerRequestedShutdown, err))
	}()

	time.Sleep(time.Millisecond * 50)

	ret, err := client.Request(http.MethodPost, "/")
	assert.Equal(t, "OK", ret)

	unixServer.Shutdown()

	time.Sleep(time.Millisecond * 50)
	_, err = client.Request(http.MethodPost, "/")
	assert.Error(t, err)
}

func TestNoResponse(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_error_response")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	unixServer, err := NewServer(
		&Config{
			Addr: tmpFile.Name(),
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
		})
	require.NoError(t, err)

	client := Client{Addr: tmpFile.Name()}
	listen := make(chan error)
	go func() {
		for range listen {
		}
	}()

	go func() {
		err = unixServer.Serve(listen)
		defer unixServer.Shutdown()
	}()

	time.Sleep(time.Millisecond * 50)

	_, err = client.Request(http.MethodGet, "/")
	require.Error(t, err)
}

func TestErrorResponse(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_error_response")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	unixServer, err := NewServer(
		&Config{
			Addr:        tmpFile.Name(),
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {},
		})
	require.NoError(t, err)

	client := Client{Addr: tmpFile.Name()}
	listen := make(chan error)
	go func() {
		for range listen {
		}
	}()

	go func() {
		err = unixServer.Serve(listen)
		defer unixServer.Shutdown()
	}()

	time.Sleep(time.Millisecond * 50)

	_, err = client.Request(http.MethodGet, "/")
	require.Error(t, err)
}

func TestResponseWriter(t *testing.T) {
	w := NewHttpResponseWriter(nil)
	require.Equal(t, make(http.Header), w.Header())
}
