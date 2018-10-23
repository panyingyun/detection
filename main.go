package main

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/corona10/goimagehash"
)

func main() {
	fmt.Println("hello")

	file1, _ := os.Open("a.jpg")
	file2, _ := os.Open("b.jpg")
	defer file1.Close()
	defer file2.Close()

	img1, _ := jpeg.Decode(file1)
	img2, _ := jpeg.Decode(file2)
	hash1, _ := goimagehash.AverageHash(img1)
	hash2, _ := goimagehash.AverageHash(img2)
	distance, _ := hash1.Distance(hash2)
	fmt.Printf("Distance between images: %v\n", distance)
}
