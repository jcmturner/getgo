package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jcmturner/godownload/releases"
)

const (
	output = `Latest Go Version Information:
Version: %s
OS: %s
Arch: %s
Kind: %s
URL: %s
SHA256: %s
Size (bytes): %d
`
)

func main() {
	osptr := flag.String("os", "", "OS type of release to download")
	archptr := flag.String("arch", "", "Architecture type of release to download")
	kind := flag.String("kind", releases.KindArchive, "Kind of release to download")
	path := flag.String("path", "./", "Path into which to write download")
	info := flag.Bool("info", false, "Show latest Go version only, do not download")
	flag.Parse()

	goos := *osptr
	if goos == "" && *kind != "source" {
		if os.Getenv("GOOS") != "" {
			goos = os.Getenv("GOOS")
		} else {
			goos = runtime.GOOS
		}
	}

	arch := *archptr
	if arch == "" && *kind != "source" {
		if os.Getenv("GOARCH") != "" {
			arch = os.Getenv("GOARCH")
		} else {
			arch = runtime.GOARCH
		}
	}

	// Source has no specific OS or Arch
	if *kind == "source" {
		arch = ""
		goos = ""
	}

	validArgs := true
	var argsErr []string
	if !releases.ValidOS(goos) {
		argsErr = append(argsErr, "invalid OS type")
		validArgs = false
	}
	if !releases.ValidArch(arch) {
		argsErr = append(argsErr, "invalid arch type")
		validArgs = false
	}
	if !releases.ValidKind(*kind) {
		argsErr = append(argsErr, "invalid kind")
		validArgs = false
	}
	if !validArgs {
		fmt.Fprintf(os.Stderr, "Error invalid arguments: %s\n", strings.Join(argsErr, "; "))
		os.Exit(1)
	}

	rs, err := releases.LoadReleaseInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading release information: %v\n", err)
		os.Exit(1)
	}
	f := rs.Latest(goos, arch, *kind)

	if f.Filename == "" {
		fmt.Fprintf(os.Stderr, "Error could not find release that matches: OS:%s, Arch:%s, Kind:%s\n", goos, arch, *kind)
		os.Exit(1)
	}

	fmt.Printf(output, f.Version, f.OS, f.Arch, f.Kind, f.URL, f.ChecksumSHA256, f.Size)
	if !*info {
		err = download(*path, f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

func download(path string, f releases.File) error {
	fn := strings.TrimRight(path, "/") + "/" + f.Filename
	dwnfile, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("Error opening destination file: %v\n", err)
	}
	fmt.Printf("Downloading to: %s\n...\n", dwnfile.Name())
	err = f.Download(dwnfile)
	if err != nil {
		return fmt.Errorf("Error downloading Go release: %v\n", err)
	}
	fmt.Println("Download Complete")
	return nil
}
