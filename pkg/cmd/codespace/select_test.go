package codespace

import (
	"context"
	"errors"
	"os"
	"fmt"
	"testing"

	"github.com/cli/cli/v2/internal/codespaces/api"
	"github.com/cli/cli/v2/pkg/iostreams"
)

const CODESPACE_NAME = "monalisa-cli-cli-abcdef"
const OUTPUT_FILE_PATH = "../../../bin/codespace-selection-test.log"

func TestApp_Select(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		opts 	selectOptions
		wantErr bool
		wantStdout string
		wantStderr string
		wantFileContents string
	}{
		{
			name: "Select a codespace",
			arg: CODESPACE_NAME,
			wantErr: false,
			wantStdout: fmt.Sprintf("%s\n", CODESPACE_NAME),
		},
		{
			name: "Select a codespace error",
			arg: "non-existent-codespace-name",
			wantErr: true,
		},
		{
			name: "Select a codespace",
			arg: CODESPACE_NAME,
			wantErr: false,
			wantFileContents: CODESPACE_NAME,
			opts: selectOptions { filePath: OUTPUT_FILE_PATH },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, stdout, stderr := iostreams.Test()
			io.SetStdinTTY(true)
			io.SetStdoutTTY(true)
			a := NewApp(io, nil, testSelectApiMock(), nil)

			if err := a.Select(context.Background(), tt.arg, tt.opts); (err != nil) != tt.wantErr {
				t.Errorf("App.Select() error = %v, wantErr %v", err, tt.wantErr)
			}

			if out := stdout.String(); out != tt.wantStdout {
				t.Errorf("stdout = %q, want %q", out, tt.wantStdout)
			}
			if out := sortLines(stderr.String()); out != tt.wantStderr {
				t.Errorf("stderr = %q, want %q", out, tt.wantStderr)
			}

			if tt.wantFileContents != "" {
				if tt.opts.filePath == "" {
					t.Errorf("wantFileContents is set but opts.filePath is not")
				}

				dat, err := os.ReadFile(tt.opts.filePath)
				if err != nil {
					panic(err)
				}

				if string(dat) != tt.wantFileContents {
					t.Errorf("file contents = %q, want %q", string(dat), CODESPACE_NAME)
				}
			}
		})
	}
}

func testSelectApiMock() *apiClientMock {
	testingCodespace := &api.Codespace{
		Name: CODESPACE_NAME,
	}
	return &apiClientMock{
		GetCodespaceFunc: func(_ context.Context, name string, includeConnection bool) (*api.Codespace, error) {
			if name == CODESPACE_NAME {
				return testingCodespace, nil
			}

			return nil, errors.New("Cannot find codespace.")
		},
	}
}
