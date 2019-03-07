package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/httpway/httpway"
)

func TestConfigValid(t *testing.T) {
	tests := []struct {
		name      string
		config    config
		assertion require.ErrorAssertionFunc
	}{
		{
			"valid #1",
			config{},
			require.NoError,
		},
		{
			"valid #2",
			config{
				Server: []configServer{
					{},
				},
			},
			require.NoError,
		},
		{
			"multiple servers",
			config{
				Server: []configServer{
					{},
					{},
				},
			},
			require.Error,
		},
		{
			"multiple actions inside a server",
			config{
				Server: []configServer{
					{
						Action: []configServerAction{
							{},
							{},
						},
					},
				},
			},
			require.Error,
		},
		{
			"multiple http inside a server",
			config{
				Server: []configServer{
					{
						HTTP: []configServerHTTP{
							{},
							{},
						},
					},
				},
			},
			require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.config.valid())
		})
	}
}

func TestConfigToGenerateConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   config
		expected []httpway.GenerateConfigHandler
	}{
		{
			"success #1",
			config{},
			[]httpway.GenerateConfigHandler{},
		},
		{
			"success #2",
			config{
				Handler: []configHandler{
					{
						Path:    "path1",
						Version: "version1",
						Alias:   "alias1",
					},
					{
						Path:    "path2",
						Version: "version2",
						Alias:   "alias2",
					},
				},
			},
			[]httpway.GenerateConfigHandler{
				{
					Path:    "path1",
					Version: "version1",
					Alias:   "alias1",
				},
				{
					Path:    "path2",
					Version: "version2",
					Alias:   "alias2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.config.toGenerateConfig().Handler
			require.ElementsMatch(t, tt.expected, actual)
		})
	}
}

func TestConfigToClientConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   config
		expected httpway.ClientConfig
	}{
		{
			"success #1",
			config{
				Host:   []configHost{},
				Server: []configServer{},
			},
			httpway.ClientConfig{
				Host:   []httpway.ClientConfigHost{},
				Server: httpway.ClientConfigServer{},
			},
		},
		{
			"success #2",
			config{
				Host: []configHost{
					{
						Endpoint: "endpoint1",
						Handler:  "handler1",
						Origin:   "origin1",
					},
					{
						Endpoint: "endpoint2",
						Handler:  "handler2",
						Origin:   "origin2",
					},
				},
				Server: []configServer{
					{
						HTTP: []configServerHTTP{
							{
								Port: 80,
							},
						},
						Action: []configServerAction{
							{
								NotFound: "notFound",
								Panic:    "panic",
							},
						},
					},
				},
			},
			httpway.ClientConfig{
				Host: []httpway.ClientConfigHost{
					{
						Endpoint: "endpoint1",
						Handler:  "handler1",
						Origin:   "origin1",
					},
					{
						Endpoint: "endpoint2",
						Handler:  "handler2",
						Origin:   "origin2",
					},
				},
				Server: httpway.ClientConfigServer{
					HTTP: httpway.ClientConfigServerHTTP{
						Port: 80,
					},
					Action: httpway.ClientConfigServerAction{
						NotFound: "notFound",
						Panic:    "panic",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.config.toClientConfig()
			actual.AsyncErrHandler = tt.expected.AsyncErrHandler
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestConfigCtxShutdown(t *testing.T) {
	tests := []struct {
		name     string
		config   config
		deadline time.Duration
		exist    require.BoolAssertionFunc
	}{
		{
			"with deadline",
			config{
				Server: []configServer{
					{
						GracefulShutdown: "1s",
					},
				},
			},
			time.Second,
			require.True,
		},
		{
			"without deadline",
			config{},
			0,
			require.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, ctxCancel := tt.config.ctxShutdown()
			defer ctxCancel()

			deadline, ok := ctx.Deadline()
			tt.exist(t, ok)
			if !ok {
				return
			}

			tdiff := time.Now().Add(tt.deadline).Sub(deadline)
			require.True(t, (tdiff >= 0))
		})
	}
}