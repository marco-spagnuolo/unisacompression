
// +build cgo,!no_cgo_flate

package main

import "github.com/marco-spagnuolo/unisacompression/internal/cgo/flate"

func init() {
	RegisterEncoder(FormatFlate, "cgo", flate.NewWriter)
	RegisterDecoder(FormatFlate, "cgo", flate.NewReader)
}
