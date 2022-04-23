package releases

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	infoURL     = "https://go.dev/dl/?mode=json"
	downloadURL = "https://dl.google.com/go/%s"

	KindArchive   = "archive"
	KindInstaller = "installer"
	KindMSI       = "installer"
	KindMacOSPkg  = "installer"
	KindSource    = "source"
)

func ValidOS(s string) bool {
	// The known operating systems. Copied from github.com/golang/go/src/cmd/dist/build.go with "" added for "source"
	var okgoos = []string{
		"darwin",
		"dragonfly",
		"linux",
		"android",
		"solaris",
		"freebsd",
		"nacl",
		"netbsd",
		"openbsd",
		"plan9",
		"windows",
		"",
	}
	for _, os := range okgoos {
		if os == s {
			return true
		}
	}
	return false
}

func ValidArch(s string) bool {
	// The known operating systems. Copied from github.com/golang/go/src/cmd/dist/build.go with "" added for "source"
	var okgoarch = []string{
		"386",
		"amd64",
		"amd64p32",
		"arm",
		"arm64",
		"mips",
		"mipsle",
		"mips64",
		"mips64le",
		"ppc64",
		"ppc64le",
		"s390x",
		"",
	}
	for _, a := range okgoarch {
		if a == s {
			return true
		}
	}
	return false
}

func ValidKind(s string) bool {
	var okkind = []string{
		"archive",
		"installer",
		"source",
	}
	for _, k := range okkind {
		if k == s {
			return true
		}
	}
	return false
}

func LoadReleaseInfo() (Releases, error) {
	var r Releases
	resp, err := http.Get(infoURL)
	if err != nil {
		return r, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return r, err
	}
	sort.Sort(r)
	for _, rf := range r {
		sort.Sort(rf.Files)
		for i := range rf.Files {
			rf.Files[i].URL = fmt.Sprintf(downloadURL, rf.Files[i].Filename)
		}
	}
	return r, nil
}

// File represents a file on the go.dev downloads page.
type File struct {
	Filename       string `json:"filename"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	Version        string `json:"version"`
	ChecksumSHA256 string `json:"sha256"`
	Size           int64  `json:"size"`
	Kind           string `json:"kind"` // "archive", "installer", "source"
	URL            string `json:"-"`
}

type Files []File

type Release struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   Files  `json:"files"`
}

type Releases []Release

func (r *Release) VersionNumbers() (maj, min, patch int) {
	maj, min, patch = parseVersion(r.Version)
	return
}

func parseVersion(v string) (maj, min, patch int) {
	p := strings.Split(strings.TrimPrefix(v, "go"), ".")
	maj, _ = strconv.Atoi(p[0])
	if len(p) < 2 {
		return
	}
	min, _ = strconv.Atoi(p[1])
	if len(p) < 3 {
		return
	}
	patch, _ = strconv.Atoi(p[2])
	return
}

func (r Releases) Len() int      { return len(r) }
func (r Releases) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Releases) Less(i, j int) bool {
	a, b := r[i], r[j]
	// Put stable releases first.
	if a.Stable != b.Stable {
		return a.Stable
	}
	if av, bv := a.Version, b.Version; av != bv {
		return versionLess(av, bv)
	}
	return true
}

func (f Files) Len() int      { return len(f) }
func (f Files) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f Files) Less(i, j int) bool {
	a, b := f[i], f[j]
	if av, bv := a.Version, b.Version; av != bv {
		return versionLess(av, bv)
	}
	if a.OS != b.OS {
		return a.OS < b.OS
	}
	if a.Arch != b.Arch {
		return a.Arch < b.Arch
	}
	if a.Kind != b.Kind {
		return a.Kind < b.Kind
	}
	return a.Filename < b.Filename
}

func versionLess(a, b string) bool {
	maja, mina, pa := parseVersion(a)
	majb, minb, pb := parseVersion(b)
	if maja == majb {
		if mina == minb {
			return pa >= pb
		}
		return mina >= minb
	}
	return maja >= majb
}

func (rs Releases) Latest(os, arch, kind string) File {
	sort.Sort(rs)
	for _, rf := range rs {
		sort.Sort(rf.Files)
	}
	for _, f := range rs[0].Files {
		if f.OS == os && f.Arch == arch && f.Kind == kind {
			return f
		}
	}
	return File{}
}

func (f File) Download(w io.Writer) error {
	resp, err := http.Get(f.URL)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)
	n, err := io.Copy(w, tee)
	if err != nil {
		return err
	}
	if n != f.Size {
		return fmt.Errorf("downloaded size %d, expected %d", n, f.Size)
	}
	hash := sha256.New()
	hash.Write(buf.Bytes())
	if hex.EncodeToString(hash.Sum(nil)) != f.ChecksumSHA256 {
		return errors.New("checksum of download does not match")
	}
	return nil
}
