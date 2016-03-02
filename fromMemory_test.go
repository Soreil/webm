package webm

import (
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

//const dataDirectory = "../testData/"

func TestWebm(t *testing.T) {
	//const filename = dataDirectory + "wafel.webm"
	const filename = "test.webm"
	if _, err := os.Stat(filename); err != nil {
		t.Fatal(err)
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	img := decode(file)
	out, err := os.Create("lmao.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(out, img); err != nil {
		t.Fatal(err)
	}
}
