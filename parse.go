package syncdeps

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	semver "github.com/blang/semver/v4"
	log "github.com/sirupsen/logrus"
)

type Dependency struct {
	Name    string         `json:"name"`
	Version semver.Version `json:"version"`
}

func (dep *Dependency) Read(p []byte) (n int, err error) {
	str := strings.Builder{}
	str.WriteString(dep.Name)
	str.WriteString(" v")
	str.WriteString(dep.Version.String())
	for i, b := range str.String() {
		p[i] = byte(b)
	}
	return str.Len(), io.EOF
}

func (dep *Dependency) Write(p []byte) (n int, err error) {
	elements := strings.Split(string(p), " ")
	if len(elements) == 1 {
		return 0, fmt.Errorf("could not split")
	}
	dep.Name = elements[0]
	v, err := parseSemVer(elements[1])
	if err != nil {
		return 0, err
	}
	dep.Version = v
	return len(p), nil
}

func WriteGoMod(w io.Writer, deps []Dependency) error {
	b, err := json.Marshal(deps)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func ScanFile(r io.Reader) (map[string][]semver.Version, error) {
	deps := make(map[string][]semver.Version)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		d := Dependency{}
		_, err := d.Write(trim(scanner.Bytes()))
		if err != nil {
			return nil, err
		}
		addToMap(deps, d)
	}
	log.Infof("Scan found %v dependencies", len(deps))
	return deps, nil
}

func trim(in []byte) []byte {
	index := strings.Index(string(in), "/go.mod")
	if index == -1 {
		return in
	}
	return in[:index]
}

func parseSemVer(version string) (semver.Version, error) {
	version = strings.TrimLeft(version, "v")
	v, err := semver.Parse(version)
	if err != nil {
		return semver.Version{}, err
	}
	return v, nil
}
