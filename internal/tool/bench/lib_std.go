
// +build !no_lib_std

package main

import (
	"io"
	"io/ioutil"

	"compress/bzip2"
	"compress/flate"
)

func init() {
	RegisterEncoder(FormatFlate, "std",
		func(w io.Writer, lvl int) io.WriteCloser {
			zw, err := flate.NewWriter(w, lvl)
			if err != nil {
				panic(err)
			}
			return zw
		})
	RegisterDecoder(FormatFlate, "std",
		func(r io.Reader) io.ReadCloser {
			return flate.NewReader(r)
		})
	RegisterDecoder(FormatBZ2, "std",
		func(r io.Reader) io.ReadCloser {
			return ioutil.NopCloser(bzip2.NewReader(r))
		})
}
