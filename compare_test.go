package syncdeps

import (
	"reflect"
	"testing"

	"github.com/blang/semver/v4"
	log "github.com/golang/glog"
)

var bigger semver.Version
var smaller semver.Version

func init() {
	var err error
	bigger, err = semver.Parse("0.10.1")
	if err != nil {
		log.Fatal(err)
	}
	smaller, err = semver.Parse("0.10.0")
	if err != nil {
		log.Fatal(err)
	}
}

func Test_extractHighestVersion(t *testing.T) {
	type args struct {
		b map[string][]semver.Version
		t map[string][]semver.Version
	}
	tests := []struct {
		name string
		args args
		want []Dependency
	}{
		{
			name: "default",
			args: args{
				b: map[string][]semver.Version{
					"test0": {bigger, smaller},
					"test1": {bigger, smaller},
					"test2": {bigger},
				},
				t: map[string][]semver.Version{
					"test0": {bigger, smaller},
					"test1": {smaller},
					"test2": {smaller},
				},
			},
			want: []Dependency{
				{"test1", bigger},
				{"test2", bigger},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractHighestVersion(tt.args.b, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractHighestVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addToMap(t *testing.T) {
	type args struct {
		m map[string][]semver.Version
		d Dependency
	}
	tests := []struct {
		name string
		args args
		want map[string][]semver.Version
	}{
		{
			name: "add on empty",
			args: args{
				m: make(map[string][]semver.Version),
				d: Dependency{"test", smaller},
			},
			want: map[string][]semver.Version{"test": {smaller}},
		},
		{
			name: "add on exist",
			args: args{
				m: map[string][]semver.Version{"test": {smaller}},
				d: Dependency{"test", bigger},
			},
			want: map[string][]semver.Version{
				"test": {bigger, smaller},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addToMap(tt.args.m, tt.args.d)
			if !reflect.DeepEqual(tt.args.m, tt.want) {
				t.Errorf("extractHighestVersion() = %v, want %v", tt.args.m, tt.want)
			}
		})
	}
}
