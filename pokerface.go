package main

import (
	"errors"
	"flag"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/luizbranco/identico"
)

var (
	dir  = flag.String("images-dir", "./images", "Images directory")
	port = flag.String("port", ":8080", "Server port")
)

func init() {
	flag.Parse()
}

func main() {
	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)

	images, err := loadImages(*dir)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		mask := randomMask(images)
		avatar := identico.Classic(mask, colorful.WarmColor(), colorful.HappyColor())
		err = png.Encode(res, avatar)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(*port, nil))
}

func loadImages(path string) ([]image.Image, error) {
	names := []string{}
	files, _ := ioutil.ReadDir(*dir)
	for _, f := range files {
		names = append(names, f.Name())
	}

	images := []image.Image{}
	for _, name := range names {
		file, err := os.Open(path + string(filepath.Separator) + name)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	if len(images) == 0 {
		return nil, errors.New("Images list is blank")
	}
	return images, nil
}

func randomMask(images []image.Image) image.Image {
	i := rand.Intn(len(images) - 1)
	return images[i]
}
