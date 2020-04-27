package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/builtin/lambda"
)

var (
	fPath    = flag.String("input", "/input", "path for the input application")
	fPre     = flag.String("pre", "", "run prehooks")
	fExtract = flag.String("extract", "", "extract written data using this path as manifest")
	fLayer   = flag.String("layer", "", "path to store layer output if observed")
)

var cleanupPaths = []string{
	"/var/cache/yum",
}

const preHook = ".devflow/pre.sh"

func main() {
	flag.Parse()

	L := hclog.L()

	if *fExtract != "" {
		var files []lambda.MappedFile

		f, err := os.Open(*fExtract)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		err = json.NewDecoder(f).Decode(&files)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("extracting %d files...\n", len(files))

		w, err := os.Create(*fExtract + ".zip")
		if err != nil {
			log.Fatal(err)
		}

		defer w.Close()

		tw := zip.NewWriter(w)

		var (
			ltw *zip.Writer
			lw  *os.File
		)

		if *fLayer != "" {
			lw, err = os.Create(*fLayer + ".zip")
			if err != nil {
				log.Fatal(err)
			}

			defer lw.Close()

			ltw = zip.NewWriter(lw)
		}

		var layerFiles int

		for _, path := range files {
			f, err := os.Open(path.Source)
			if err != nil {
				log.Fatal(err)
			}

			stat, err := f.Stat()
			if err != nil {
				log.Fatal(err)
			}

			if stat.IsDir() {
				continue
			}

			hdr, err := zip.FileInfoHeader(stat)
			if err != nil {
				log.Fatal(err)
			}

			hdr.Modified = time.Time{}
			hdr.ModifiedTime = 0
			hdr.ModifiedDate = 0

			h := sha1.New()

			if ltw != nil && strings.HasPrefix(path.Dest, "_layer/") {
				hdr.Name = path.Dest[len("_layer/"):]

				body, err := ltw.CreateHeader(hdr)
				if err != nil {
					log.Fatal(err)
				}

				io.Copy(io.MultiWriter(body, h), f)

				layerFiles++

			} else {
				hdr.Name = path.Dest

				body, err := tw.CreateHeader(hdr)
				if err != nil {
					log.Fatal(err)
				}

				io.Copy(io.MultiWriter(body, h), f)
			}

			f.Close()

			// fmt.Printf("added %s (%s)\n", hdr.Name, hex.EncodeToString(h.Sum(nil))[:10])
		}

		tw.Close()
		if ltw != nil {
			ltw.Close()
		}
		w.Close()

		if ltw != nil && layerFiles == 0 {
			os.Remove(lw.Name())
		}

		L.Info("finished gathering files")

		os.Exit(0)
	}

	L.Info("Starting Devflow Lambda Builder...")

	if *fPre != "" {
		L.Info("Executing prehooks...")

		f, err := os.Open(*fPre)
		if err != nil {
			L.Error("error opening pre", "error", err)
			os.Exit(1)
		}

		to, err := os.Create("/tmp/pre.sh")
		if err != nil {
			L.Error("error creating pre temp", "error", err)
			os.Exit(1)
		}

		io.Copy(to, f)

		to.Close()
		f.Close()

		pre := exec.Command("bash", "/tmp/pre.sh")
		pre.Stdout = os.Stdout
		pre.Stderr = os.Stderr
		pre.Dir = "/tmp"

		err = pre.Run()
		if err != nil {
			L.Error("Error executing pre.sh", "error", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	if _, err := os.Stat(*fPath); err != nil {
		L.Error("error checking on input path", "error", err, "path", *fPath)
		os.Exit(1)
	}

	cmd := exec.Command("/buildpack/bin/detect")
	cmd.Dir = *fPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		L.Error("error in detection process", "error", err)
		os.Exit(1)
	}

	L.Info("Detected properly, building...")

	os.MkdirAll("/var/task", 0755)

	cp := exec.Command("cp", "-r", *fPath+"/.", "/var/task")
	cp.Stdout = os.Stdout
	cp.Stderr = os.Stderr

	err = cp.Run()
	if err != nil {
		L.Error("error copying application into image", "error", err)
		os.Exit(1)
	}

	cmd = exec.Command("/buildpack/bin/build")
	cmd.Dir = "/var/task"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		L.Error("error in build process", "error", err)
		os.Exit(1)
	}

	L.Info("Cleaning up image")

	cmd = exec.Command("yum", "clean", "all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		L.Error("error in build process", "error", err)
		os.Exit(1)
	}
}
