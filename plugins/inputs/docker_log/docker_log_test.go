package docker_log

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/stretchr/testify/require"

	"github.com/extremenetworks/telegraf"
	"github.com/extremenetworks/telegraf/config"
	"github.com/extremenetworks/telegraf/testutil"
)

type MockClient struct {
	ContainerListF    func() ([]types.Container, error)
	ContainerInspectF func() (types.ContainerJSON, error)
	ContainerLogsF    func() (io.ReadCloser, error)
}

func (c *MockClient) ContainerList(context.Context, container.ListOptions) ([]types.Container, error) {
	return c.ContainerListF()
}

func (c *MockClient) ContainerInspect(context.Context, string) (types.ContainerJSON, error) {
	return c.ContainerInspectF()
}

func (c *MockClient) ContainerLogs(context.Context, string, container.LogsOptions) (io.ReadCloser, error) {
	return c.ContainerLogsF()
}

type Response struct {
	io.Reader
}

func (r *Response) Close() error {
	return nil
}

func MustParse(layout, value string) time.Time {
	tm, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return tm
}

func Test(t *testing.T) {
	tests := []struct {
		name     string
		client   *MockClient
		expected []telegraf.Metric
	}{
		{
			name: "no containers",
			client: &MockClient{
				ContainerListF: func() ([]types.Container, error) {
					return nil, nil
				},
			},
		},
		{
			name: "one container tty",
			client: &MockClient{
				ContainerListF: func() ([]types.Container, error) {
					return []types.Container{
						{
							ID:    "deadbeef",
							Names: []string{"/telegraf"},
							Image: "extremenetworks/telegraf:1.11.0",
						},
					}, nil
				},
				ContainerInspectF: func() (types.ContainerJSON, error) {
					return types.ContainerJSON{
						Config: &container.Config{
							Tty: true,
						},
					}, nil
				},
				ContainerLogsF: func() (io.ReadCloser, error) {
					return &Response{Reader: bytes.NewBuffer([]byte("2020-04-28T18:43:16.432691200Z hello\n"))}, nil
				},
			},
			expected: []telegraf.Metric{
				testutil.MustMetric(
					"docker_log",
					map[string]string{
						"container_name":    "telegraf",
						"container_image":   "extremenetworks/telegraf",
						"container_version": "1.11.0",
						"stream":            "tty",
						"source":            "deadbeef",
					},
					map[string]interface{}{
						"container_id": "deadbeef",
						"message":      "hello",
					},
					MustParse(time.RFC3339Nano, "2020-04-28T18:43:16.432691200Z"),
				),
			},
		},
		{
			name: "one container multiplex",
			client: &MockClient{
				ContainerListF: func() ([]types.Container, error) {
					return []types.Container{
						{
							ID:    "deadbeef",
							Names: []string{"/telegraf"},
							Image: "extremenetworks/telegraf:1.11.0",
						},
					}, nil
				},
				ContainerInspectF: func() (types.ContainerJSON, error) {
					return types.ContainerJSON{
						Config: &container.Config{
							Tty: false,
						},
					}, nil
				},
				ContainerLogsF: func() (io.ReadCloser, error) {
					var buf bytes.Buffer
					w := stdcopy.NewStdWriter(&buf, stdcopy.Stdout)
					_, err := w.Write([]byte("2020-04-28T18:42:16.432691200Z hello from stdout"))
					return &Response{Reader: &buf}, err
				},
			},
			expected: []telegraf.Metric{
				testutil.MustMetric(
					"docker_log",
					map[string]string{
						"container_name":    "telegraf",
						"container_image":   "extremenetworks/telegraf",
						"container_version": "1.11.0",
						"stream":            "stdout",
						"source":            "deadbeef",
					},
					map[string]interface{}{
						"container_id": "deadbeef",
						"message":      "hello from stdout",
					},
					MustParse(time.RFC3339Nano, "2020-04-28T18:42:16.432691200Z"),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var acc testutil.Accumulator
			plugin := &DockerLogs{
				Timeout:          config.Duration(time.Second * 5),
				newClient:        func(string, *tls.Config) (Client, error) { return tt.client, nil },
				containerList:    make(map[string]context.CancelFunc),
				IncludeSourceTag: true,
			}

			err := plugin.Init()
			require.NoError(t, err)

			err = plugin.Gather(&acc)
			require.NoError(t, err)

			acc.Wait(len(tt.expected))
			plugin.Stop()

			require.Nil(t, acc.Errors) // no errors during gathering

			testutil.RequireMetricsEqual(t, tt.expected, acc.GetTelegrafMetrics())
		})
	}
}
