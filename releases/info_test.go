package releases

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

func TestLoadReleaseInfo(t *testing.T) {
	rs, err := LoadReleaseInfo()
	if err != nil {
		t.Fatalf("error loading releases: %v", err)
	}
	t.Logf("%+v\n", rs)
}

func TestReleases_Latest(t *testing.T) {
	rs, err := LoadReleaseInfo()
	if err != nil {
		t.Fatalf("error loading releases: %v", err)
	}
	f := rs.Latest(runtime.GOOS, runtime.GOARCH, "archive")
	t.Logf("%+v\n", f)
}

func TestFile_Download(t *testing.T) {
	rs, err := LoadReleaseInfo()
	if err != nil {
		t.Fatalf("error loading releases: %v", err)
	}
	f := rs.Latest(runtime.GOOS, runtime.GOARCH, "archive")
	dwnfile, _ := ioutil.TempFile(os.TempDir(), f.Filename)
	defer os.Remove(dwnfile.Name())
	err = f.Download(dwnfile)
	if err != nil {
		t.Fatalf("error downloading golang: %v", err)
	}
}
