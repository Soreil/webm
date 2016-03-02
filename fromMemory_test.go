package webm

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

const dataDirectory = "testData/"

func TestWebmToPNG(t *testing.T) {
	const filename = dataDirectory + "wafel.webm"
	if _, err := os.Stat(filename); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	out, err := os.Create(filename + ".png")
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(out, img); err != nil {
		t.Fatal(err)
	}
}

func TestWebmToJPG(t *testing.T) {
	const filename = dataDirectory + "wafel.webm"
	if _, err := os.Stat(filename); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	out, err := os.Create(filename + ".jpg")
	if err != nil {
		t.Fatal(err)
	}
	if err := jpeg.Encode(out, img, nil); err != nil {
		t.Fatal(err)
	}
}
