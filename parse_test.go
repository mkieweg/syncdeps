package syncdeps

import (
	"bytes"
	"io"
	"io/fs"
	"reflect"
	"testing"

	semver "github.com/blang/semver/v4"
)

func Test_dependency_Read(t *testing.T) {
	type fields struct {
		name    string
		version semver.Version
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		{
			name: "cloud.google.com/go v0.26.0",
			fields: fields{
				name: "cloud.google.com/go",
				version: semver.Version{
					Major: 0,
					Minor: 26,
				},
			},
			args:  args{p: []byte("cloud.google.com/go v0.26.0")},
			wantN: 27,
		},
		{
			name: "github.com/google/uuid v1.2.0",
			fields: fields{
				name: "github.com/google/uuid",
				version: semver.Version{
					Major: 1,
					Minor: 2,
				},
			},
			args:  args{p: []byte("github.com/google/uuid v1.2.0")},
			wantN: 29,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gm := &Dependency{
				Name:    tt.fields.name,
				Version: tt.fields.version,
			}
			p := make([]byte, 64)
			gotN, err := gm.Read(p)
			if (err != nil) != tt.wantErr {
				if err != io.EOF {
					t.Errorf("dependency.Read() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if gotN != tt.wantN {
				t.Errorf("dependency.Read() = %v, want %v", gotN, tt.wantN)
			}
			if !bytes.Equal(tt.args.p, p[:gotN]) {
				t.Errorf("dependency.Read() = %v, want %v", string(p[:gotN]), string(tt.args.p))
			}
		})
	}
}

func Test_dependency_Write(t *testing.T) {
	type fields struct {
		name    string
		version semver.Version
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		{
			name: "cloud.google.com/go v0.26.0",
			fields: fields{
				name: "cloud.google.com/go",
				version: semver.Version{
					Major: 0,
					Minor: 26,
				},
			},
			args:  args{p: []byte("cloud.google.com/go v0.26.0")},
			wantN: 27,
		},
		{
			name: "github.com/google/uuid v1.2.0",
			fields: fields{
				name: "github.com/google/uuid",
				version: semver.Version{
					Major: 1,
					Minor: 2,
				},
			},
			args:  args{p: []byte("github.com/google/uuid v1.2.0")},
			wantN: 29,
		},
		{
			name:    "invalid",
			fields:  fields{},
			args:    args{p: []byte("invalid")},
			wantN:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := &Dependency{
				Name:    tt.fields.name,
				Version: tt.fields.version,
			}
			got := &Dependency{}
			gotN, err := got.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("dependency.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("dependency.Write() = %v, want %v", gotN, tt.wantN)
			}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("dependency.Write() = %v, want %v", got, want)
			}
		})
	}
}

func TestScanFile(t *testing.T) {
	type args struct {
		file fs.File
	}
	tests := []struct {
		name    string
		args    args
		want    []Dependency
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScanFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trim(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=",
			args: args{in: []byte("cloud.google.com/go v0.26.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=")},
			want: []byte("cloud.google.com/go v0.26.0"),
		},
		{
			name: "github.com/google/uuid v1.2.0",
			args: args{in: []byte("github.com/google/uuid v1.2.0")},
			want: []byte("github.com/google/uuid v1.2.0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trim(tt.args.in)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("trim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSemVer(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    semver.Version
		wantErr bool
	}{
		{
			name: "1.0.0",
			args: args{version: "1.0.0"},
			want: semver.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
		},
		{
			name: "v1.0.0",
			args: args{version: "v1.0.0"},
			want: semver.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
		},
		{
			name: "v1.0.0+build",
			args: args{version: "v1.0.0+build"},
			want: semver.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
				Build: []string{"build"},
			},
		},
		{
			name:    "erroneous",
			args:    args{version: "erroneous"},
			want:    semver.Version{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSemVer(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSemVer() = %v, want %v", got, tt.want)
			}
		})
	}
}
