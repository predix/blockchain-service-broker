package schema

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("manifest-schema")

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

func NewManifest(manifestFile string) (*Manifest, error) {
	data, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		log.Error("Error reading manifest file", err)
		return nil, err
	}

	manifest := Manifest{}

	err = yaml.Unmarshal(data, &manifest)
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
