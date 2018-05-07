# godownload

Simple command line tool that will download the latest version of Go.

It may seem a little chicken and egg for a tool written in Go to be used to download Go.
However due to the statically compile nature of Go there is no dependency on the Go language package to run this tool.

The purpose of this tool is to be used in automated build pipelines to download the latest version of Go before compiling your Go project with it.

The checksum of the file is verified as part of the download process.

## Build
On a host that already has Go installed:

```go get -u github.com/jcmturner/godownload```

## Run
To simply download the latest Go version to the current working directory for the OS and architecture ```godownload``` is run on:

```
./godownload
```

The OS, architecture and path to download to can be specified with the following arguments:

```
Usage of godownload:
  -arch string
    	Architecture type of release to download
  -info
    	Show latest Go version only, do not download
  -kind string
    	Kind of release to download (default "archive")
  -os string
    	OS type of release to download
  -path string
    	Path into which to write download (default "./")
```
If the GOOS and GOARCH environment variables are set these will be used if the ```-os``` and ```-arch``` are not provided.

To get the latest version information without downloading use the ```-info``` switch.