package main

import (
	"flag"
	"fmt"
	"os"

	"clc.hmu/app/public/ffmpeg"
)

func main() {
	uri := flag.String("uri", "rtsp://admin:hyiot123@192.168.1.64:554/Streaming/Channels/101", "rtsp uri")
	cmd := flag.String("cmd", "ffmpeg", "cmd path")
	flag.Parse()
	fmt.Println(*uri, *cmd)

	cameraUri := *uri

	ff := ffmpeg.NewCmd(*cmd).Debug(true)
	if err := ff.FetchImage(
		cameraUri,
		os.TempDir()+fmt.Sprintf("/testing_%d.jpg", 0),
		"640x480", // size
		"mjpeg",   // format
		"0",       // position time
		1,         // vframes
	); err != nil {
		panic(err)
	}

	if err := ff.FetchVideo(
		cameraUri,
		os.TempDir()+fmt.Sprintf("/testing_%d.mp4", 0),
		"640x480", // size
		"mp4",     // format
		"0",       // position time
		5,         // 5 seconds
		20,        // 20 fps
	); err != nil {
		panic(err)
	}

	if err := ff.FetchVideo(
		cameraUri,
		os.TempDir()+fmt.Sprintf("/testing_%d.gif", 0),
		"640x480", // size
		"gif",     // format
		"0",       // position time
		5,         // 5s
		10,        // 10 fps
	); err != nil {
		panic(err)
	}

}
