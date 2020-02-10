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

type protocGenGoTestExepectation struct {
	resultFile string
	goldenFile string
}

var protocGenGoTests = []struct {
	name         string
	args         []string
	expectations []protocGenGoTestExepectation
}{
	{
		"QRPC",
		[]string{
			"-I=./test/proto",
			"--qrpcgo_out=./test/result",
			"./test/proto/test_api.proto",
		},
		[]protocGenGoTestExepectation{
			{
				"caller/test_api.pb.go",
				"caller/test_api.pb.go.golden",
			},
			{
				"handler/test_api.pb.go",
				"handler/test_api.pb.go.golden",
			},
		},
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

			for _, exp := range tt.expectations {
				result := loadResult(t, exp.resultFile)
				golden := loadGolden(t, exp.goldenFile)

				assert.Equal(t, string(golden), string(result))
			}
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
