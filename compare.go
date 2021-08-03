package syncdeps

import (
	"sort"

	"github.com/blang/semver/v4"
)

func Compare(baseline, target map[string][]semver.Version) []Dependency {
	return extractHighestVersion(baseline, target)
}

func extractHighestVersion(b, t map[string][]semver.Version) []Dependency {
	out := make([]Dependency, 0)
	for k, v := range t {
		versions, ok := b[k]
		if !ok {
			continue
		}
		if v[0].LT(versions[0]) {
			out = append(out, Dependency{k, versions[0]})
		}
	}
	if len(out) != 0 {
		return out
	}
	return nil
}

func checkAndAdd(m map[string][]semver.Version, d Dependency) {
	versions, ok := m[d.Name]
	if !ok {
		m[d.Name] = []semver.Version{d.Version}
	} else {
		versions = append(versions, d.Version)
		sort.Slice(versions, func(i, j int) bool {
			return versions[i].GT(versions[j])
		})
		m[d.Name] = versions
	}
}
