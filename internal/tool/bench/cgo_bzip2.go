
// +build cgo,!no_cgo_bzip2

package main

import "github.com/marco-spagnuolo/unisacompression/internal/cgo/bzip2"

func init() {
	RegisterEncoder(FormatBZ2, "cgo", bzip2.NewWriter)
	RegisterDecoder(FormatBZ2, "cgo", bzip2.NewReader)
}
