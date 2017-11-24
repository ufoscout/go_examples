// GOAUTOPATH POC
// This code is in public domain.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Example usage:
// GOBIN=`pwd` GOAUTOPATH=1 \
// goa install  ../../../../gogs/util/src/tesa/dummy.go

// chkAutoPath appends GOPATH with the path to the build dir
// estabilished automatically from the path to file or dir
// given as the last argument to the main
func chkAutoPath() (err error) {

	/*
		if os.Getenv("GOAUTOPATH") == "" {
			return
		}
	*/
	a := len(os.Args) - 1
	v := os.Args[a] // get last
	os.Args = os.Args[:a]
	s := string(os.PathSeparator)

	v, err = filepath.Abs(filepath.Clean(v)) // normalize to cwd
	if err != nil {
		err = fmt.Errorf("Bad src path given: %s [%s]", v, err)
		return
	}
	if strings.HasSuffix(v, ".go") {
		v = filepath.Dir(v)
	}

	// find build top
	a = strings.LastIndex(v, s+"src"+s)
	if a < 0 {
		err = fmt.Errorf("No src updir found in: %s", v)
		return
	}
	bd := v[:a]   // build top
	mp := v[a+1:] // module path
	dt := ""      // temp

	// Autopath-ed go cmd may not pollute the tree.
	for _, cd := range [...]string{"", "pkg", "bin", mp} {
		dt = bd + s + cd
		if err = os.Chdir(dt); err != nil {
			err = fmt.Errorf("Can not enter into directory %s.", dt)
			return
		}
	}

	// append GOPATH with bd
	dt = filepath.Clean(os.Getenv("GOPATH"))

	if !filepath.IsAbs(dt) || len(dt) < 3 {
		err = fmt.Errorf("GOPATH is bad: %s", dt)
		return
	}
	dt = dt + string(os.PathListSeparator) + bd

	if err = os.Setenv("GOPATH", dt); err != nil {
		err = fmt.Errorf("Can not set GOPATH to %s [%s]", dt, err)
		return
	}
	return
}

func main() {
	if err := chkAutoPath(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n-----\n", err)
	}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "GO") {
			fmt.Fprintf(os.Stderr, "ENV: %s\n", v)
		}
	}
	fmt.Fprintf(os.Stderr, "CMD: ")
	for _, v := range os.Args {
		fmt.Fprintf(os.Stderr, "%s ", v)
	}
	fmt.Fprintf(os.Stderr, "\n")

	os.Setenv("GOPATH", "custom")

	cmd := exec.Command("go", "env")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
