package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jcmturner/godownload/releases"
)

func main() {
	goos := flag.String("os", "", "OS type of release to download")
	arch := flag.String("arch", "", "Architecture type of release to download")
	kind := flag.String("kind", releases.KindArchive, "Kind of release to download")
	path := flag.String("path", "./", "Path into which to write download")
	info := flag.Bool("info", false, "Show latest Go version only, do not download")
	flag.Parse()

	if *goos == "" {
		if os.Getenv("GOOS") != "" {
			t := os.Getenv("GOOS")
			goos = &t
		} else {
			t := runtime.GOOS
			goos = &t
		}
	}

	if *arch == "" {
		if os.Getenv("GOARCH") != "" {
			t := os.Getenv("GOARCH")
			goos = &t
		} else {
			t := runtime.GOARCH
			goos = &t
		}
	}

	rs, err := releases.LoadReleaseInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading release information: %v\n", err)
		os.Exit(1)
	}
	f := rs.Latest(*goos, *arch, *kind)

	if *info {
		fmt.Printf(`Latest Go Version Information:
Version: %s
OS: %s
Arch: %s
Kind: %s
`, f.Version, f.OS, f.Arch, f.Kind)
	}
	err = download(*path, f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func download(path string, f releases.File) error {
	fn := strings.TrimRight(path, "/") + "/" + f.Filename
	dwnfile, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("Error opening destination file: %v\n", err)
	}
	err = f.Download(dwnfile)
	if err != nil {
		return fmt.Errorf("Error downloading Go release: %v\n", err)
	}
	fmt.Printf(`Go release download successful:
Location: %s
Version: %s
OS: %s
Arch: %s
Kind: %s
`, dwnfile.Name(), f.Version, f.OS, f.Arch, f.Kind)
	return nil
}
