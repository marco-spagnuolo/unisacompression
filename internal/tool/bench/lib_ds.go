
// +build !no_lib_ds

package main

import (
	"io"

	"github.com/marco-spagnuolo/unisacompression/brotli"
	"github.com/marco-spagnuolo/unisacompression/bzip2"
	"github.com/marco-spagnuolo/unisacompression/flate"
)

func init() {
	RegisterDecoder(FormatBrotli, "ds",
		func(r io.Reader) io.ReadCloser {
			zr, err := brotli.NewReader(r, nil)
			if err != nil {
				panic(err)
			}
			return zr
		})
	RegisterDecoder(FormatFlate, "ds",
		func(r io.Reader) io.ReadCloser {
			zr, err := flate.NewReader(r, nil)
			if err != nil {
				panic(err)
			}
			return zr
		})
	RegisterEncoder(FormatBZ2, "ds",
		func(w io.Writer, lvl int) io.WriteCloser {
			zw, err := bzip2.NewWriter(w, &bzip2.WriterConfig{Level: lvl})
			if err != nil {
				panic(err)
			}
			return zw
		})
	RegisterDecoder(FormatBZ2, "ds",
		func(r io.Reader) io.ReadCloser {
			zr, err := bzip2.NewReader(r, nil)
			if err != nil {
				panic(err)
			}
			return zr
		})
}
