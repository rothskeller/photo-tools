package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: analyze-groups directory\n")
		os.Exit(2)
	}
	if err := os.Chdir(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	fh, err := os.Open("groups")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	defer fh.Close()
	startExiftool()
	scan := bufio.NewScanner(fh)
	var paths []string
	for scan.Scan() {
		path := scan.Text()
		if path == "" {
			analyzeGroup(paths)
			paths = paths[:0]
		} else {
			paths = append(paths, path)
		}
	}
	if len(paths) != 0 {
		analyzeGroup(paths)
	}
	stopExiftool()
}

var unreadableImage = image.NewUniform(color.Black)

func analyzeGroup(paths []string) {
	if len(paths) == 0 {
		return
	}
	var (
		images    []image.Image
		dates     []string
		authors   []string
		titles    []string
		captions  []string
		keywords  []string
		faces     []string
		coords    []string
		locations []string
	)
	fmt.Println("IDATCKFGL Path")
	for _, path := range paths {
		image, date, author, title, caption, keyword, face, coord, location := getMetadata(path)
		if image == nil {
			fmt.Print(" ")
		} else if image == unreadableImage {
			images = append(images, image)
			fmt.Printf("%d", len(images))
		} else {
			var found bool
			for i, img := range images {
				if img != unreadableImage && sameImage(img, image) {
					fmt.Printf("%d", i+1)
					found = true
					break
				}
			}
			if !found {
				images = append(images, image)
				fmt.Printf("%d", len(images))
			}
		}
		dates = compareStrings(dates, date)
		authors = compareStrings(authors, author)
		titles = compareStrings(titles, title)
		captions = compareStrings(captions, caption)
		keywords = compareStrings(keywords, keyword)
		faces = compareStrings(faces, face)
		coords = compareStrings(coords, coord)
		locations = compareStrings(locations, location)
		fmt.Printf(" %s\n", path)
	}
	fmt.Println("---------")
	for i, date := range dates {
		fmt.Printf("D%d        %s\n", i+1, date)
	}
	for i, author := range authors {
		fmt.Printf(" A%d       %s\n", i+1, author)
	}
	for i, title := range titles {
		fmt.Printf("  T%d      %s\n", i+1, title)
	}
	for i, caption := range captions {
		fmt.Printf("   C%d     %s\n", i+1, caption)
	}
	for i, keyword := range keywords {
		fmt.Printf("    K%d    %s\n", i+1, keyword)
	}
	for i, face := range faces {
		fmt.Printf("     F%d   %s\n", i+1, face)
	}
	for i, coord := range coords {
		fmt.Printf("      G%d  %s\n", i+1, coord)
	}
	for i, location := range locations {
		fmt.Printf("       L%d %s\n", i+1, location)
	}
	fmt.Println()
}

func compareStrings(list []string, val string) []string {
	if val == "" {
		fmt.Print(" ")
		return list
	}
	for i, s := range list {
		if s == val {
			fmt.Printf("%d", i+1)
			return list
		}
	}
	list = append(list, val)
	fmt.Printf("%d", len(list))
	return list
}

func sameImage(a, b image.Image) bool {
	ab := a.Bounds()
	bb := b.Bounds()
	if !ab.Eq(bb) {
		return false
	}
	for x := ab.Min.X; x < ab.Max.X; x++ {
		for y := ab.Min.Y; y < ab.Max.Y; y++ {
			ar, ag, ab, aa := a.At(x, y).RGBA()
			br, bg, bb, ba := b.At(x, y).RGBA()
			if ar != br || ag != bg || ab != bb || aa != ba {
				return false
			}
		}
	}
	return true
}
