package schema_test

import (
	"strings"
	"testing"

	"github.com/atulkc/fabric-service-broker/schema"

	. "gopkg.in/go-playground/assert.v1"
)

func TestNewManifest(t *testing.T) {
	manifest, err := schema.NewManifest("../resources/templates/fabric.yml")

	stemcell := schema.Stemcell{
		Alias:   "default",
		Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
		Version: "latest",
	}

	Equal(t, err, nil)
	NotEqual(t, manifest, nil)

	Equal(t, manifest.Name, "hyperledger-fabric")
	Equal(t, manifest.DirectorUuid, "28539132-6d43-4e1b-bf40-f2ce032ee9f8")
	Equal(t, manifest.Stemcells[0], stemcell)
	Equal(t, manifest.Properties.Peer.Network, map[string]string{"id": "dev"})
	Equal(t, manifest.Properties.Peer.Consensus, map[string]string{"plugin": "pbft"})
}

func TestManifestToString(t *testing.T) {
	manifest, err := schema.NewManifest("../resources/templates/fabric.yml")

	Equal(t, err, nil)
	NotEqual(t, manifest, nil)

	manifest.Name = "test-deployment-name"

	Equal(t, strings.Contains(manifest.String(), "name: test-deployment-name"), true)
	Equal(t, strings.Contains(manifest.String(), "name: hyperledger-fabric"), false)
	Equal(t, strings.Contains(manifest.String(), "plugin: pbft"), true)
}
