package testutils

import (
	"log"

	"github.com/nyaxt/otaru/facade"
	"github.com/nyaxt/otaru/gcloud/auth"
)

var testConfigCached *facade.Config

func TestConfig() *facade.Config {
	if testConfigCached != nil {
		return testConfigCached
	}

	cfg, err := facade.NewConfig(facade.DefaultConfigDir())
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	testConfigCached = cfg
	return testConfigCached
}

func TestClientSource() auth.ClientSource {
	cfg := TestConfig()

	clisrc, err := auth.GetGCloudClientSource(cfg.CredentialsFilePath, cfg.TokenCacheFilePath, false)
	if err != nil {
		log.Fatalf("Failed to create testClientSource: %v", err)
	}
	return clisrc
}

func TestBucketName() string {
	name := TestConfig().TestBucketName
	if len(name) == 0 {
		log.Fatalf("Please specify \"test_bucket_name\" in config.toml to run gcloud tests.")
	}
	return name
}