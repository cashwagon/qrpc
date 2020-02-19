package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var protocGenGoTests = []struct {
	name   string
	args   []string
	result string
	golden string
}{
	{
		"Basic",
		[]string{
			"-I=./test/proto",
			"--go_out=paths=source_relative:./test/result",
			"./test/proto/test_api.proto",
		},
		"test_api.pb.go",
		"test_api.pb.go.basic.golden",
	},
	{
		"GRPC",
		[]string{
			"-I=./test/proto",
			"--go_out=plugins=grpc,paths=source_relative:./test/result",
			"./test/proto/test_api.proto",
		},
		"test_api.pb.go",
		"test_api.pb.go.grpc.golden",
	},
	{
		"QRPC",
		[]string{
			"-I=./test/proto",
			"--go_out=plugins=qrpc,paths=source_relative:./test/result",
			"./test/proto/test_api.proto",
		},
		"test_api.pb.go",
		"test_api.pb.go.qrpc.golden",
	},
}

func Test_ProtocGenGo(t *testing.T) {
	protocBin := getProtoc(t)
	projectPath, err := filepath.Abs("..")
	require.NoError(t, err, "Cannot get project path")

	binPath := filepath.Join(projectPath, "bin")

	for _, tt := range protocGenGoTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := os.Mkdir("./result", os.ModeDir|os.ModePerm|755)
			require.NoError(t, err, "cannot create test results directory")

			defer func() {
				require.NoError(t, os.RemoveAll("./result"))
			}()

			pathEnv := os.Getenv("PATH")

			cmd := exec.Command(protocBin, tt.args...) // nolint:gosec // Variable is safe
			cmd.Dir = projectPath
			cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s:%s", binPath, pathEnv))

			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Cannot execute protoc '%s':\n%s", protocBin, output)

			result := loadResult(t, tt.result)
			golden := loadGolden(t, tt.golden)

			assert.Equal(t, string(golden), string(result))
		})
	}
}

func loadGolden(t *testing.T, name string) []byte {
	t.Helper()

	return loadFile(t, filepath.Join("golden", filepath.Clean(name)))
}

func loadResult(t *testing.T, name string) []byte {
	t.Helper()

	return loadFile(t, filepath.Join("result", filepath.Clean(name)))
}

func loadFile(t *testing.T, name string) []byte {
	t.Helper()

	file, err := os.Open(filepath.Clean(name))
	require.NoErrorf(t, err, "Cannot open file %s", name)

	defer func() {
		require.NoError(t, file.Close())
	}()

	data, err := ioutil.ReadAll(file)
	require.NoError(t, err, "Cannot read file %s", name)

	return data
}

func getProtoc(t *testing.T) string {
	t.Helper()

	binaryPath := os.Getenv("PROTOC")
	if binaryPath != "" {
		return binaryPath
	}

	return "protoc"
}
