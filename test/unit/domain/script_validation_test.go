package domain_test

import (
	"testing"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestCreateScriptRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     domain.CreateScriptRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     "root",
				Password: "password",
				Template: "template1",
				Port:     2222,
			},
			wantErr: false,
		},
		{
			name: "valid hostname",
			req: domain.CreateScriptRequest{
				Host:     "example.com",
				User:     "root",
				Password: "password",
				Template: "template1",
				Port:     22,
			},
			wantErr: false,
		},
		{
			name: "empty host",
			req: domain.CreateScriptRequest{
				Host:     "",
				User:     "root",
				Password: "password",
				Template: "template1",
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "invalid host",
			req: domain.CreateScriptRequest{
				Host:     "invalid_host!@#",
				User:     "root",
				Password: "password",
				Template: "template1",
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "empty user",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     "",
				Password: "password",
				Template: "template1",
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "password too short (min 8)",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     "root",
				Password: "123",
				Template: "template1",
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "empty template",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     "root",
				Password: "password",
				Template: "",
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "host too long (>255)",
			req: domain.CreateScriptRequest{
				Host:     string(make([]byte, 300)),
				User:     "root",
				Password: "password",
				Template: "template1",
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "user too long (>100)",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     string(make([]byte, 150)),
				Password: "password",
				Template: "template1",
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "password too long (>128)",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     "root",
				Password: string(make([]byte, 200)),
				Template: "template1",
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "template too long (>50)",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     "root",
				Password: "password",
				Template: string(make([]byte, 100)),
				Port:     22,
			},
			wantErr: true,
		},
		{
			name: "invalid port (0)",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     "root",
				Password: "password",
				Template: "template1",
				Port:     0,
			},
			wantErr: true,
		},
		{
			name: "invalid port (>65535)",
			req: domain.CreateScriptRequest{
				Host:     "127.0.0.1",
				User:     "root",
				Password: "password",
				Template: "template1",
				Port:     99999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCallbackRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     domain.CallbackRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: domain.CallbackRequest{
				User:       "root",
				ScriptPath: "/opt/script-monitor/scripts/test.sh",
				Action:     "execute",
				Time:       "2026-08-13T08:09:36Z",
			},
			wantErr: false,
		},
		{
			name: "valid action modify",
			req: domain.CallbackRequest{
				User:       "root",
				ScriptPath: "/opt/script-monitor/scripts/test.sh",
				Action:     "modify",
				Time:       "2026-08-13T08:09:36Z",
			},
			wantErr: false,
		},
		{
			name: "empty user",
			req: domain.CallbackRequest{
				User:       "",
				ScriptPath: "/opt/script-monitor/scripts/test.sh",
				Action:     "execute",
				Time:       "2026-08-13T08:09:36Z",
			},
			wantErr: true,
		},
		{
			name: "user too long (>100)",
			req: domain.CallbackRequest{
				User:       string(make([]byte, 150)),
				ScriptPath: "/opt/script-monitor/scripts/test.sh",
				Action:     "execute",
				Time:       "2026-08-13T08:09:36Z",
			},
			wantErr: true,
		},
		{
			name: "empty script path",
			req: domain.CallbackRequest{
				User:       "root",
				ScriptPath: "",
				Action:     "execute",
				Time:       "2026-08-13T08:09:36Z",
			},
			wantErr: true,
		},
		{
			name: "script path too long (>4096)",
			req: domain.CallbackRequest{
				User:       "root",
				ScriptPath: string(make([]byte, 5000)),
				Action:     "execute",
				Time:       "2026-08-13T08:09:36Z",
			},
			wantErr: true,
		},
		{
			name: "invalid action",
			req: domain.CallbackRequest{
				User:       "root",
				ScriptPath: "/opt/script-monitor/scripts/test.sh",
				Action:     "invalid",
				Time:       "2026-08-13T08:09:36Z",
			},
			wantErr: true,
		},
		{
			name: "invalid time format",
			req: domain.CallbackRequest{
				User:       "root",
				ScriptPath: "/opt/script-monitor/scripts/test.sh",
				Action:     "execute",
				Time:       "invalid-time",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
