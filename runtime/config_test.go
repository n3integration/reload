package runtime

import (
	"testing"

	"github.com/n3integration/reload/test"
)

func Test_LoadConfig(t *testing.T) {
	config, err := LoadConfig("testdata/config.json")

	test.Expect(t, err, nil)
	test.Expect(t, config.Port, 5678)
	test.Expect(t, config.ProxyTo, "http://localhost:3000")
}

func Test_LoadConfig_WithNonExistentFile(t *testing.T) {
	_, err := LoadConfig("im/not/here.json")

	test.Refute(t, err, nil)
	test.Expect(t, err.Error(), "unable to read configuration file im/not/here.json")
}

func Test_LoadConfig_WithMalformedFile(t *testing.T) {
	_, err := LoadConfig("testdata/bad_config.json")

	test.Refute(t, err, nil)
	test.Expect(t, err.Error(), "unable to parse configuration file testdata/bad_config.json")
}
