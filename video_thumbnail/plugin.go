package video_thumbnail

import (
	"bufio"
	"fmt"
	"github.com/ccwawamiya/screengen"
	"github.com/chai2010/webp"
	flutter "github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const channelName = "video_thumbnail"

// VideoThumbnailPlugin implements flutter.Plugin and handles method.
type VideoThumbnailPlugin struct{}

var _ flutter.Plugin = &VideoThumbnailPlugin{} // compile-time type check

const (
	imageFormatJPEG = 0
	imageFormatPNG  = 1
	imageFormatWEBP = 2
)

// InitPlugin initializes the plugin.
func (p *VideoThumbnailPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	channel.HandleFunc("file", p.handleThumbnailFile)
	return nil
}

func (p *VideoThumbnailPlugin) handleThumbnailFile(arguments interface{}) (reply interface{}, err error) {
	argsMap := arguments.(map[interface{}]interface{})
	videoPath := argsMap["video"].(string)
	thumbnailPath := argsMap["path"].(string)
	imageFormat := int(argsMap["format"].(int32))
	maxHeight := int(argsMap["maxh"].(int32))
	maxWidth := int(argsMap["maxw"].(int32))
	timeMs := int(argsMap["timeMs"].(int32))
	quality := int(argsMap["quality"].(int32))
	var fileName string
	fileName, err = thumbnailFile(videoPath, thumbnailPath, imageFormat, maxHeight, maxWidth, timeMs, quality)
	if err != nil {
		fmt.Println(err.Error())
	}
	return fileName, nil
}

func thumbnailFile(videoPath, thumbnailPath string, imageFormat, maxHeight, maxWidth, timeMs, quality int) (fileName string, err error) {
	g := &screengen.Generator{}
	g, err = screengen.NewGenerator(videoPath)
	if err != nil {
		return
	}
	if quality > 100 {
		quality = 100
	} else if quality <= 0 {
		quality = 10
	}
	videoHeight := g.Height()
	videoWidth := g.Width()
	imgHeight := videoHeight
	imgWidth := videoWidth

	if videoHeight > maxHeight && maxHeight > 0 {
		imgHeight = maxHeight
	}
	if videoWidth > maxWidth && maxWidth > 0 {
		imgWidth = maxWidth
	}
	if imgHeight/videoHeight <= imgWidth/videoWidth {
		imgWidth = videoWidth * imgHeight / videoHeight
	} else {
		imgHeight = videoHeight * imgWidth / videoWidth
	}
	fmt.Println(imgWidth, imgHeight)
	var img image.Image
	img, err = g.ImageWxH(int64(timeMs), imgWidth, imgHeight)
	if err != nil {
		return
	}
	switch imageFormat {
	case imageFormatPNG:
		fileName, err = saveImg("png", thumbnailPath, quality, img)
		break
	case imageFormatJPEG:
		fileName, err = saveImg("jpeg", thumbnailPath, quality, img)
		break
	case imageFormatWEBP:
		fileName, err = saveImg("webp", thumbnailPath, quality, img)
		break
	default:
		fileName, err = saveImg("png", thumbnailPath, quality, img)
		break
	}

	return
}

func saveImg(format, savePath string, quality int, img image.Image) (fileName string, err error) {
	fileName = savePath + "/" + randString(6) + "." + format
	fmt.Println(fileName)
	var outFile *os.File
	outFile, err = os.Create(fileName)
	if err != nil {
		os.Exit(-1)
		return
	}
	defer outFile.Close()
	buff := bufio.NewWriter(outFile)

	if format == "jpeg" {
		err = jpeg.Encode(buff, img, &jpeg.Options{Quality: quality})
	} else if format == "webp" {
		err = webp.Encode(buff, img, &webp.Options{Quality: float32(quality)})
	} else {
		err = png.Encode(buff, img)
	}

	if err != nil {
		os.Exit(-1)
		return
	}
	err = buff.Flush()
	if err != nil {
		os.Exit(-1)
		return
	}
	return
}

func randString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}
