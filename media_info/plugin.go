package media_info

import (
	"bufio"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"gitlab.com/opennota/screengen"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"
)

const channelName = "asia.ivity.flutter/media_info"

// MediaInfoPlugin implements flutter.Plugin and handles method.
type MediaInfoPlugin struct{}

var _ flutter.Plugin = &MediaInfoPlugin{} // compile-time type check

// InitPlugin initializes the plugin.
func (p *MediaInfoPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	channel.HandleFunc("getMediaInfo", p.getMediaInfo)
	channel.HandleFunc("generateThumbnail", p.generateThumbnail)
	return nil
}

func (p *MediaInfoPlugin) getMediaInfo(arguments interface{}) (reply interface{}, err error) {
	reply = map[interface{}]interface{}{
		"mimeType":   "-",
		"width":      int32(0),
		"height":     int32(0),
		"frameRate":  float64(0),
		"durationMs": int32(0),
	}
	videoPath := arguments.(string)
	mime := &mimetype.MIME{}
	mime, err = mimetype.DetectFile(videoPath)
	if err != nil {
		return
	}
	if strings.Index(mime.String(), "audio/") == -1 || strings.Index(mime.String(), "video/") == -1 {
		g := &screengen.Generator{}
		g, err = screengen.NewGenerator(videoPath)
		if err != nil {
			return
		}
		reply = map[interface{}]interface{}{
			"mimeType":   mime.String(),
			"width":      int32(g.Width()),
			"height":     int32(g.Height()),
			"frameRate":  g.FPS,
			"durationMs": int32(g.Duration),
		}
	}
	return
}

func (p *MediaInfoPlugin) generateThumbnail(arguments interface{}) (reply interface{}, err error) {
	reply = ""
	argsMap := arguments.(map[interface{}]interface{})
	videoPath := argsMap["path"].(string)
	target := argsMap["target"].(string)
	width := int(argsMap["width"].(int32))
	height := int(argsMap["height"].(int32))
	mime := &mimetype.MIME{}
	mime, err = mimetype.DetectFile(videoPath)
	if err != nil {
		return
	}
	if strings.Index(mime.String(), "audio/") == -1 || strings.Index(mime.String(), "video/") == -1 {
		g := &screengen.Generator{}
		g, err = screengen.NewGenerator(videoPath)
		if err != nil {
			return
		}
		var img image.Image
		img, err = g.ImageWxH(0, width, height)
		if err != nil {
			return
		}
		err = os.MkdirAll(path.Dir(target), 0755)
		if err != nil {
			return
		}
		var outFile *os.File
		outFile, err = os.Create(target)
		if err != nil {
			os.Exit(-1)
			return
		}
		defer outFile.Close()
		buff := bufio.NewWriter(outFile)
		err = jpeg.Encode(buff, img, &jpeg.Options{Quality: 100})
		if err == nil {
			reply = target
		}
	}
	return
}
