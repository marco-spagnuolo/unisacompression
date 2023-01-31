package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/marco-spagnuolo/unisacompression/bzip2"
)

func main() {
	// Set up a slice to store the performance data for each block size
	var data []PerformanceData

	for num := 1; num < 10; num++ {
		// Generate a file with a unique name
		s := "test" + strconv.Itoa(num) + ".txt"
		generateFile(1000, s)

		// Compress the file and measure the time it takes
		startTime := time.Now()
		compressWithBlockSize(s, num)
		compressTime := time.Since(startTime)

		// Decompress the file and measure the time it takes
		startTime = time.Now()
		decompress(s + ".bz2")
		decompressTime := time.Since(startTime)
		//compare file
		if compareFile(s, s+".unbz2") {
			fmt.Println("File ", s, " and ", s+".unbz2", " are the same")
		} else {
			fmt.Println("File ", s, " and ", s+".unbz2", " are different")
		}
		fmt.Println()

		// Gather the size of the original file, the size of the compressed file, and the times for compression and decompression
		d := PerformanceData{
			BlockSize:      num,
			OriginalSize:   int(getSize(s)),
			CompressedSize: int(getSize(s + ".bz2")),
			CompressTime:   compressTime,
			DecompressTime: decompressTime,
		}
		data = append(data, d)

		// Clean up the test file and compressed versions
		clean()
	}

	// print result in a beauty way
	GenerateDiagram(data)

	//save data
	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "CompressionTime: Blue ,DecompressionTime: Green ,Compressed size: yellow",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "BlockSize",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Time in µs",
		}),
	)
	// Put data into instance
	bar.SetXAxis([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}).
		AddSeries("Category A", generateBarItemsCompressionTime(data)).
		AddSeries("Category B", generateBarItemsDecompressTime(data)).
		AddSeries("Category C", generateBarItemsCompressedSize(data)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)

	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)

}

// generate a random file named @s of size @n byte with only 4 characters
func generateFile(n int, s string) {
	f, err := os.Create(s)
	log.Println("Sto generando il file... ", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	var letters = []rune("abcd")
	//s := make([]rune, 1)
	for i := 0; i < n; i++ {
		f.WriteString(string(letters[rand.Intn(len(letters))]))
	}
}

//normal compression
func compress(s string) error {

	fin, err := os.Open(s)
	log.Println("Sto aprendo il file      ", s)
	if err != nil {
		return err
	}
	defer fin.Close()

	s = s + ".bz2"

	f2, err := os.Create(s)
	log.Println("Sto comprimedo il file     ", s)
	if err != nil {
		return err
	}
	defer f2.Close()

	w, err := bzip2.NewWriter(f2,
		&bzip2.WriterConfig{
			Level: bzip2.BestCompression,
		})
	if err != nil {
		return err

	}

	if _, err := io.Copy(w, fin); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

//normal decompression
func decompress(s string) error {

	fin, err := os.Open(s)
	log.Println("Sto aprendo il file      ", s)
	if err != nil {
		return err
	}
	defer fin.Close()

	replacedText := strings.Replace(s, "bz2", "unbz2", 3)

	f2, err := os.Create(replacedText)
	log.Println("Sto decomprimedo il file ", replacedText)
	if err != nil {
		return err
	}
	defer f2.Close()

	r, err := bzip2.NewReader(fin,
		&bzip2.ReaderConfig{})
	if err != nil {
		return err
	}
	if _, err := io.Copy(f2, r); err != nil {
		return err
	}
	return nil

}

//compression with different block size given as parameter
func compressWithBlockSize(s string, n int) error {

	fin, err := os.Open(s)
	log.Println("Sto aprendo il file      ", s)

	if err != nil {
		return err
	}
	defer fin.Close()

	s = s + ".bz2"

	f2, err := os.Create(s)
	log.Println("Sto comprimedo il file   ", s)

	if err != nil {
		return err
	}
	defer f2.Close()

	w, err := bzip2.NewWriter(f2,
		&bzip2.WriterConfig{
			Level: n,
		})
	if err != nil {
		return err

	}

	if _, err := io.Copy(w, fin); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil

}

//clean dir from gen file
func clean() {
	// Get a list of all files in the current directory
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Loop through the files
	for _, file := range files {
		// Check if the file has a .txt or .bzip2 extension
		if filepath.Ext(file.Name()) == ".txt" || filepath.Ext(file.Name()) == ".bz2" || filepath.Ext(file.Name()) == ".unbz2" {
			// Remove the file
			err := os.Remove(file.Name())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

//get file size
func getSize(s string) int64 {
	fi, err := os.Stat(s)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return fi.Size()
}

// compare file
func compareFile(s1, s2 string) bool {
	f1, err := os.Open(s1)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f1.Close()
	f2, err := os.Open(s2)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f2.Close()
	b1 := make([]byte, 1)
	b2 := make([]byte, 1)
	for {
		_, err1 := f1.Read(b1)
		_, err2 := f2.Read(b2)
		if err1 == io.EOF && err2 == io.EOF {
			return true
		}
		if err1 != nil || err2 != nil {
			return false
		}
		if b1[0] != b2[0] {
			return false
		}
	}
}

// PerformanceData represents the data gathered for each block size
type PerformanceData struct {
	BlockSize      int
	OriginalSize   int
	CompressedSize int
	CompressTime   time.Duration
	DecompressTime time.Duration
}

// GenerateTensorDiagram generates a tensor diagram of the performance data for each block size and save into a file
func GenerateDiagram(data []PerformanceData) {
	for _, d := range data {
		fmt.Printf("Block size %d: Original size = %d, Compressed size = %d, Compress time = %v, Decompress time = %v\n", d.BlockSize, d.OriginalSize, d.CompressedSize, d.CompressTime, d.DecompressTime)
	}
}

// generate random data for bar chart
func generateBarItemsCompressionTime(data []PerformanceData) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 9; i++ {
		items = append(items, opts.BarData{Value: int(data[i].CompressTime)})
	}
	return items
}

func generateBarItemsDecompressTime(data []PerformanceData) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 9; i++ {
		items = append(items, opts.BarData{Value: int(data[i].DecompressTime)})
	}
	return items
}

func generateBarItemsCompressedSize(data []PerformanceData) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 9; i++ {
		items = append(items, opts.BarData{Value: int(data[i].CompressedSize)})
	}
	return items
}

// func main() {
// 	for num := 1; num < 10; num++ {
// 		// Generate a file with a unique name
// 		s := "test" + strconv.Itoa(num) + ".txt"
// 		generateFile(10000, s)

// 		// Compress it with a different block size
// 		fmt.Println("Compressing with block size ", num)
// 		start := time.Now()

// 		compressWithBlockSize(s, num)

// 		elapsed := time.Since(start)
// 		fmt.Println("Compressing took ", elapsed)
// 		// Decompress it
// 		decompress(s + ".bz2")
// 		fmt.Println("The original file has size:", getSize(s))
// 		fmt.Println("The compressed file has size:", getSize(s+".bz2"))

// 		// Compare the original file to the decompressed file
// 		fmt.Println("The original file is equal to the decompressed file?")

// 		fmt.Println(compareFile(s, s+".unbz2"))
// 		fmt.Println("--------------------------------------------------------------------")

// 		// Clean up the test file and compressed versions
// 		clean()
// 	}
// }
