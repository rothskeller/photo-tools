// photo-titles reads a file containing blocks of text, and renders a set of
// images with that text.  These are intended to be interspersed in a set of
// photos, providing running commentary in a photo album.
//
// The input text file (named on the command line, falling back to standard
// input if there are no command line arguments) contains blocks of text,
// separated by lines of three or more dashes.  Each block can have one of three
// forms:
//
// 1.  A single line of text.  This is treated as a section heading.  It is
//     rendered in a large font, centered in the image.
// 2.  Multiple lines of text, possibly split into paragraphs with blank lines.
//     These are rendered as paragraphs, left-justified, vertically centered.
// 3.  A single line of text, followed by a blank line, followed by multiple
//     lines of text, possibly split into paragraphs with blank lines.  These
//     are rendered as paragraphs, left-justified, vertically centered, except
//     that the first line is rendered as a title, horizontally centered above
//     the paragraphs.
//
// The output images are 1000 pixel square, black background JPEG images named
// title01.jpg, title02.jpg, etc., created in the current working directory, and
// overwriting any such files already present.  Section headings and titles are
// displayed in white; paragraphs are displayed in gray.  Section headings are
// in 120pt bold font, reduced as necessary to make them fit (using 72dpi).
// Paragraphs are in 36pt regular font, reduced as necessary.  Titles are bold
// and one-sixth larger font size than paragraphs, reduced as necessary.  Images
// have a 40 pixel margin.  Section headings and titles are also stored in the
// JPEG metadata as titles, and paragraphs are stored as image descriptions.
//
// photo-titles depends on having "exiftool" in the path in order to set image
// metadata.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/golang/freetype"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Maintenance note: the Go libraries for rendering text in an image are in a
// state of interrupted transition.  At present, the new library
// (golang.org/x/image/font/opentype) is the only one that can give font metrics
// and measure strings without rendering them, while the old library
// (github.com/golang/freetype) is the only one that can actually render them.
// Unfortunately that means this code has to use both.

const imageSize = 1000
const margin = 40

var separator = regexp.MustCompile(`^-+$`)
var regularFT, _ = freetype.ParseFont(goregular.TTF)
var boldFT, _ = freetype.ParseFont(gobold.TTF)
var regular, _ = sfnt.Parse(goregular.TTF)
var bold, _ = sfnt.Parse(gobold.TTF)
var black = image.NewUniform(color.Black)
var white = image.NewUniform(color.White)
var gray = image.NewUniform(color.Gray{Y: 0x99})
var transparent = image.NewUniform(color.Transparent)

func main() {
	var (
		input  *bufio.Scanner
		number int
		lines  []string
		more   = true
	)
	if len(os.Args) > 1 {
		fh, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		input = bufio.NewScanner(fh)
	} else {
		input = bufio.NewScanner(os.Stdin)
	}
	for more {
		number++
		lines = nil
		more = false
		for input.Scan() {
			var line = input.Text()
			if separator.MatchString(line) {
				more = true
				break
			}
			lines = append(lines, strings.TrimSpace(line))
		}
		if err := input.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		renderTitle(number, lines)
	}
}

func renderTitle(number int, lines []string) {
	var (
		img         draw.Image
		title       string
		description string
		buf         bytes.Buffer
		args        []string
		cmd         *exec.Cmd
		filename    string
		err         error
	)
	for len(lines) > 0 && lines[0] == "" { // remove empty lines at start
		lines = lines[1:]
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" { // remove empty lines at end
		lines = lines[:len(lines)-1]
	}
	// Create a grayscale image.
	img = image.NewGray(image.Rect(0, 0, imageSize, imageSize))
	draw.Draw(img, img.Bounds(), black, image.ZP, draw.Src)
	switch {
	case len(lines) == 1:
		renderHeading(img, lines[0])
		title = lines[0]
	case len(lines) > 1:
		title, description = renderParagraphs(img, lines)
	}
	// Encode the image into JPEG format.
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: jpeg encoding: %s\n", err)
		os.Exit(1)
	}
	// Write it to a file.
	filename = fmt.Sprintf("title%02d.jpg", number)
	if err = ioutil.WriteFile(filename, buf.Bytes(), 0666); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	// Set the metadata and modification time on the file.
	if title != "" {
		args = append(args, "-title="+title)
	}
	if description != "" {
		args = append(args, "-imagedescription="+description)
	}
	args = append(args, "-overwrite_original", "-q", filename)
	cmd = exec.Command("exiftool", args...)
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: exiftool: %s\n", err)
		os.Exit(1)
	}
}

// renderHeading renders a single line of text as a section heading title,
// horizontally and vertically centered in a large, bold font in white.
func renderHeading(img draw.Image, heading string) {
	var (
		width fixed.Int26_6
		draw  font.Drawer
		ctx   *freetype.Context
		err   error
		size  = 120.0
	)
LOOP:
	if draw.Face, err = opentype.NewFace(bold, &opentype.FaceOptions{DPI: 72.0, Size: size, Hinting: font.HintingFull}); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: NewFace: %s\n", err)
		os.Exit(1)
	}
	width = font.MeasureString(draw.Face, heading)
	if width.Ceil() > imageSize-2*margin {
		size -= 6.0
		goto LOOP
	}
	ctx = freetype.NewContext()
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetFont(boldFT)
	ctx.SetFontSize(size)
	ctx.SetHinting(font.HintingFull)
	ctx.SetSrc(white)
	if _, err = ctx.DrawString(heading, fixed.Point26_6{
		X: (fixed.I(imageSize) - width) / 2,
		Y: (fixed.I(imageSize) + draw.Face.Metrics().CapHeight) / 2,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: DrawString: %s\n", err)
		os.Exit(1)
	}
}

func renderParagraphs(img draw.Image, lines []string) (title, description string) {
	var (
		paragraphs []string
		rows       []string
		width      fixed.Int26_6
		height     fixed.Int26_6
		yshift     fixed.Int26_6
		theight    fixed.Int26_6
		twidth     fixed.Int26_6
		draw       font.Drawer
		ctx        *freetype.Context
		pt         fixed.Point26_6
		err        error
		size       = 36.0
	)
	if lines[1] == "" {
		title = lines[0]
		lines = lines[2:]
	}
	j := 0
	paragraphs = []string{""}
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			j++
			paragraphs = append(paragraphs, "")
		} else if paragraphs[j] == "" {
			paragraphs[j] = lines[i]
		} else {
			paragraphs[j] = paragraphs[j] + " " + lines[i]
		}
	}
	description = strings.Join(paragraphs, "\n\n")
SIZE:
	// If we have a title to display, figure out its line height, assuming a
	// size slightly greater than the base text size, and a bold font.
	height = 0
	rows = []string{}
	if title != "" {
		if draw.Face, err = opentype.NewFace(bold, &opentype.FaceOptions{DPI: 72.0, Size: size + 6.0, Hinting: font.HintingFull}); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: NewFace: %s\n", err)
			os.Exit(1)
		}
		height = draw.Face.Metrics().Height
		theight = draw.Face.Metrics().Height
		twidth = font.MeasureString(draw.Face, title)
		yshift = draw.Face.Metrics().Height - draw.Face.Metrics().Ascent
	}
	// Prepare the regular font in the size we're trying out.
	if draw.Face, err = opentype.NewFace(regular, &opentype.FaceOptions{DPI: 72.0, Size: size, Hinting: font.HintingFull}); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: NewFace: %s\n", err)
		os.Exit(1)
	}
	if title == "" {
		yshift = draw.Face.Metrics().Height - draw.Face.Metrics().Ascent
	}
	// Now measure each paragraph.
	for i, p := range paragraphs {
		var (
			line      string
			remainder string
			idx       int
		)
		if i != 0 || title != "" {
			rows = append(rows, "")
		}
		// Start by putting the whole paragraph into line.
		line, remainder = p, ""
		for line != "" {
			// Check to see if the line fits.
			width = font.MeasureString(draw.Face, line)
			if width.Ceil() <= imageSize-2*margin {
				// It does fit.  Add this line to the list of
				// rows, and move on to the remainder (or the
				// next paragraph, if there is no remainder).
				rows = append(rows, line)
				line, remainder = remainder, ""
				continue
			}
			// It doesn't fit, so move the last word of line into
			// remainder and try again.
			if idx = strings.LastIndexByte(line, ' '); idx < 0 {
				// Uh oh.  There was only one word.  This font
				// size won't do.
				size -= 6.0
				goto SIZE
			}
			remainder = strings.TrimSpace(line[idx:] + " " + remainder)
			line = strings.TrimSpace(line[:idx])
		}
	}
	// Does the result fit in the frame?
	height += draw.Face.Metrics().Height * fixed.Int26_6(len(rows))
	if height.Ceil() > imageSize-2*margin {
		// No, so we'll need to try a smaller font.
		size -= 6.0
		goto SIZE
	}
	// Yes, it does.  Draw the text.
	ctx = freetype.NewContext()
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetHinting(font.HintingFull)
	pt.Y = (fixed.I(imageSize)-height)/2 - yshift
	if title != "" {
		// Draw the title.
		pt.Y += theight
		pt.X = (fixed.I(imageSize) - twidth) / 2
		ctx.SetFont(boldFT)
		ctx.SetFontSize(size + 6.0)
		ctx.SetSrc(white)
		if _, err = ctx.DrawString(title, pt); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: DrawString: %s\n", err)
			os.Exit(1)
		}
	}
	ctx.SetFont(regularFT)
	ctx.SetFontSize(size)
	ctx.SetSrc(gray)
	pt.X = fixed.I(margin)
	for _, r := range rows {
		// Draw each row.
		pt.Y += draw.Face.Metrics().Height
		if _, err = ctx.DrawString(r, pt); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: DrawString: %s\n", err)
			os.Exit(1)
		}
	}
	return title, description
}
