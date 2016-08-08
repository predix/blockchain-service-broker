package schema_test

import (
	"strings"
	"testing"

	"github.com/atulkc/fabric-service-broker/schema"

	. "gopkg.in/go-playground/assert.v1"
)

func TestNewManifest(t *testing.T) {
	manifest, err := schema.NewManifest()

	stemcell := schema.Stemcell{
		Alias:   "default",
		Name:    "USE-IAAS-SPECIFIC-STEMCELL",
		Version: "latest",
	}

	Equal(t, err, nil)
	NotEqual(t, manifest, nil)

	Equal(t, manifest.Name, "GIVE-ME-A-NAME")
	Equal(t, manifest.DirectorUuid, "REPLACE-ME")
	Equal(t, manifest.Stemcells[0], stemcell)
	Equal(t, manifest.Properties.Peer.Network, map[string]string{"id": "GENERATED"})
	Equal(t, manifest.Properties.Peer.Consensus, map[string]string{"plugin": "pbft"})
}

func TestManifestToString(t *testing.T) {
	manifest, err := schema.NewManifest()

	Equal(t, err, nil)
	NotEqual(t, manifest, nil)

	manifest.Name = "test-deployment-name"

	Equal(t, strings.Contains(manifest.String(), "name: test-deployment-name"), true)
	Equal(t, strings.Contains(manifest.String(), "name: hyperledger-fabric"), false)
	Equal(t, strings.Contains(manifest.String(), "plugin: pbft"), true)
}
