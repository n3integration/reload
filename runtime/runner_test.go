package runtime

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/n3integration/reload/test"
)

func Test_NewRunner(t *testing.T) {
	filename := getBinFile()
	runner := NewRunner(filename)

	fi, _ := runner.Info()
	test.Expect(t, fi.Name(), filepath.Base(filename))
}

func Test_Runner_Run(t *testing.T) {
	runner := NewRunner(getBinFile())

	cmd, err := runner.Run()
	test.Expect(t, err, nil)
	test.Expect(t, cmd.Process == nil, false)
}

func Test_Runner_Kill(t *testing.T) {
	bin := getBinFile()
	runner := NewRunner(bin)

	cmd1, err := runner.Run()
	test.Expect(t, err, nil)

	_, err = runner.Run()
	test.Expect(t, err, nil)

	time.Sleep(time.Second * 1)
	os.Chtimes(bin, time.Now(), time.Now())
	if err != nil {
		t.Fatal("Error with Chtimes")
	}

	cmd3, err := runner.Run()
	test.Expect(t, err, nil)

	if runtime.GOOS != "windows" {
		// does not seem to work as expected on windows
		test.Refute(t, cmd1, cmd3)
	}
}

func getBinFile() string {
	bin := filepath.Join("testdata", "writing_output")
	if runtime.GOOS == "windows" {
		bin += ".bat"
	}
	return bin
}
