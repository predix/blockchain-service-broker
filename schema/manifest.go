package schema

import (
	"gopkg.in/yaml.v2"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("manifest-schema")

const defaultManifest = `
---
name: GIVE-ME-A-NAME
director_uuid: REPLACE-ME
stemcells:
- alias: default
  name: USE-IAAS-SPECIFIC-STEMCELL
  version: latest
releases:
- name: fabric-release
  version: latest
update:
  canaries: 1
  canary_watch_time: 5000-120000
  max_in_flight: 3
  serial: false
  update_watch_time: 5000-120000
jobs:
- instances: 3
  azs: [z1, z2]
  name: peer
  networks:
  - name: peer
  persistent_disk: 1024
  vm_type: small
  stemcell: default
  templates:
  - name: peer
    release: fabric-release
  - name: docker
    release: fabric-release
properties:
  peer:
    network:
      id: GENERATED
    consensus:
      plugin: pbft
`

type Manifest struct {
	Name         string     `yaml:"name"`
	DirectorUuid string     `yaml:"director_uuid"`
	Stemcells    Stemcells  `yaml:"stemcells"`
	Releases     Releases   `yaml:"releases"`
	Update       Update     `yaml:"update"`
	Jobs         Jobs       `yaml:"jobs"`
	Properties   Properties `yaml:"properties"`
}

type Stemcells []Stemcell

type Stemcell struct {
	Alias   string `yaml:"alias"`
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Releases []Release

type Release struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Update struct {
	Canaries        uint   `yaml:"canaries"`
	CanaryWatchTime string `yaml:"canary_watch_time"`
	MaxInFlight     uint   `yaml:"max_in_flight"`
	Serial          bool   `yaml:"serial"`
	UpdateWatchTime string `yaml:"update_watch_time"`
}

type Jobs []Job

type Job struct {
	Instances      uint                `yaml:"instances"`
	AZs            []string            `yaml:"azs"`
	Name           string              `yaml:"name"`
	Networks       []map[string]string `yaml:"networks"`
	PersistentDisk uint                `yaml:"persistent_disk"`
	VmType         string              `yaml:"vm_type"`
	Stemcell       string              `yaml:"stemcell"`
	Templates      []map[string]string `yaml:"templates"`
}

type PeerProperties struct {
	Network   map[string]string `yaml:"network"`
	Consensus map[string]string `yaml:"consensus"`
}

type Properties struct {
	Peer PeerProperties `yaml:"peer"`
}

func NewManifest() (*Manifest, error) {
	manifest := Manifest{}

	err := yaml.Unmarshal([]byte(defaultManifest), &manifest)
	if err != nil {
		log.Error("Error unmarshalling manifest file", err)
		return nil, err
	}

	return &manifest, nil
}

func (m *Manifest) String() string {
	d, err := yaml.Marshal(m)
	if err != nil {
		log.Error("Error marshalling manifest", err)
		return ""
	}

	return string(d)
}
