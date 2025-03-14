package mongobin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/skyne/memongo/mongobin"
	"github.com/stretchr/testify/assert"
)

// Change this to true to issue a HEAD request in each test to make
// sure the file is there and accessible. We leave this off for reliability,
// but it can be turned on if you want to test that the generated URLs point
// to real files.
const testHTTPHead = true

func TestGetDownloadURL(t *testing.T) {
	mongoVersionsToTest := []string{
		"3.2.0", "3.2.22", "3.4.0", "3.4.19", "3.6.0", "3.6.10", "4.0.0", "4.0.13", "4.2.1", "4.4.7", "5.0.0",
	}

	tests := map[string]struct {
		spec          *mongobin.DownloadSpec
		mongoVersions []string

		expectedURL string
	}{
		"mac-ssl": {
			spec: &mongobin.DownloadSpec{
				Platform:       "osx",
				Arch:           "x86_64",
				OSName:         "",
				SSLBuildNeeded: true,
			},
			mongoVersions: []string{
				"3.2.0", "3.2.22", "3.4.0", "3.4.19", "3.6.0", "3.6.10", "4.0.0", "4.0.13",
			},
			expectedURL: "https://fastdl.mongodb.org/osx/mongodb-osx-ssl-x86_64-VERSION.tgz",
		},
		"mac": {
			spec: &mongobin.DownloadSpec{
				Platform: "osx",
				Arch:     "x86_64",
				OSName:   "",
			},
			mongoVersions: []string{
				"4.2.1",
			},
			expectedURL: "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-VERSION.tgz",
		},
		"arm64 mac": {
			spec: &mongobin.DownloadSpec{
				Platform: "osx",
				Arch:     "arm64",
				OSName:   "",
			},
			mongoVersions: []string{
				"6.0.0",
			},
			expectedURL: "https://fastdl.mongodb.org/osx/mongodb-macos-arm64-VERSION.tgz",
		},
		"ubuntu 18.04": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "ubuntu1804",
			},
			mongoVersions: []string{"4.0.1", "4.0.13", "4.2.1"},
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1804-VERSION.tgz",
		},
		"arm64 ubuntu 18.04": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "aarch64",
				OSName:   "ubuntu1804",
			},
			mongoVersions: []string{"4.2.1", "4.4.0", "6.0.0"},
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-aarch64-ubuntu1804-VERSION.tgz",
		},
		"arm64 ubuntu 20.04": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "aarch64",
				OSName:   "ubuntu2004",
			},
			mongoVersions: []string{"4.4.0", "6.0.0"},
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-aarch64-ubuntu2004-VERSION.tgz",
		},
		"arm64 ubuntu 22.04": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "aarch64",
				OSName:   "ubuntu2204",
			},
			mongoVersions: []string{"6.0.4"},
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-aarch64-ubuntu2204-VERSION.tgz",
		},
		"ubuntu 16.04": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "ubuntu1604",
			},
			mongoVersions: []string{
				"3.2.7", "3.4.0", "3.4.19", "3.6.0", "3.6.10", "4.0.0", "4.0.13", "4.2.1",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1604-VERSION.tgz",
		},
		"arm64 ubuntu 16.04": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "arm64",
				OSName:   "ubuntu1604",
			},
			mongoVersions: []string{
				"3.4.0", "3.4.19", "3.6.0", "3.6.10", "4.0.0", "4.0.13",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-arm64-ubuntu1604-VERSION.tgz",
		},
		"ubuntu 14.04": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "ubuntu1404",
			},
			mongoVersions: []string{
				"3.2.0", "3.2.22", "3.4.0", "3.4.19", "3.6.0", "3.6.10", "4.0.0", "4.0.13",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1404-VERSION.tgz",
		},
		"SUSE 12": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "suse12",
			},
			mongoVersions: mongoVersionsToTest,
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-suse12-VERSION.tgz",
		},
		"ARM64 RHEL 8.2": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "aarch64",
				OSName:   "rhel82",
			},
			mongoVersions: []string{"4.4.4", "5.0.1", "6.0.4"},
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-aarch64-rhel82-VERSION.tgz",
		},
		"RHEL 7": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "rhel70",
			},
			mongoVersions: mongoVersionsToTest,
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-rhel70-VERSION.tgz",
		},
		"RHEL 6": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "rhel62",
			},
			mongoVersions: []string{
				"3.2.0", "3.2.22", "3.4.0", "3.4.19", "3.6.0", "3.6.10", "4.0.0", "4.0.13", "4.2.1", "4.4.7",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-rhel62-VERSION.tgz",
		},
		"Debian buster": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "debian10",
			},
			mongoVersions: []string{
				"4.2.1", "4.4.7", "5.0.0",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-debian10-VERSION.tgz",
		},
		"Debian stretch": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "debian92",
			},
			mongoVersions: []string{
				"3.6.5", "3.6.10", "4.0.0", "4.0.13", "4.2.1",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-debian92-VERSION.tgz",
		},
		"Debian jessie": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "debian81",
			},
			mongoVersions: []string{
				"3.2.8", "3.2.22", "3.4.0", "3.4.19", "3.6.0", "3.6.10", "4.0.0", "4.0.13",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-debian81-VERSION.tgz",
		},
		"Amazon Linux": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "amazon",
			},
			mongoVersions: mongoVersionsToTest,
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-amazon-VERSION.tgz",
		},
		"Amazon Linux 2": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "amazon2",
			},
			mongoVersions: []string{
				"4.0.0", "4.0.13", "4.2.1",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-amazon2-VERSION.tgz",
		},
		"ARM64 Amazon Linux 2": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "aarch64",
				OSName:   "amazon2",
			},
			mongoVersions: []string{"4.4.4", "5.0.1", "6.0.4"},
			expectedURL:   "https://fastdl.mongodb.org/linux/mongodb-linux-aarch64-amazon2-VERSION.tgz",
		},
		"Other Linux": {
			spec: &mongobin.DownloadSpec{
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "",
			},
			mongoVersions: []string{
				"3.2.0", "3.2.22", "3.4.0", "3.4.19", "3.6.0", "3.6.10", "4.0.0", "4.0.13",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-VERSION.tgz",
		},
	}

	for testName, test := range tests {
		for _, mongoVersion := range test.mongoVersions {
			t.Run(testName+"_"+mongoVersion, func(t *testing.T) {
				spec := &mongobin.DownloadSpec{
					Version:        mongoVersion,
					Platform:       test.spec.Platform,
					Arch:           test.spec.Arch,
					OSName:         test.spec.OSName,
					SSLBuildNeeded: test.spec.SSLBuildNeeded,
				}

				expectedURL := strings.Replace(test.expectedURL, "VERSION", mongoVersion, -1)
				actualURL := spec.GetDownloadURL()

				if testHTTPHead {
					resp, err := http.Head(actualURL)
					assert.NoError(t, err)
					assert.Equal(t, 200, resp.StatusCode)
				}

				assert.Equal(t, expectedURL, actualURL)
			})
		}
	}
}
