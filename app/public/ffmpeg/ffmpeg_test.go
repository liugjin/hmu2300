package ffmpeg

import (
	"fmt"
	"os"
	"testing"
)

const (
	testCameraUri = "rtsp://admin:hyiot123@192.168.1.64:554/Streaming/Channels/101"
)

func TestFetchImage(t *testing.T) {
	exec := NewCmd("ffmpeg").Debug(true)
	for i := 2; i > 0; i-- {
		if err := exec.FetchImage(
			testCameraUri,
			os.TempDir()+fmt.Sprintf("/testing_%d.jpg", i),
			"640x480",  // size
			"mjpeg",    // format
			"00:00:00", // position time
			1,          // vframes
		); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFetchVideo(t *testing.T) {
	exec := NewCmd("ffmpeg").Debug(true)
	for i := 2; i > 0; i-- {
		if err := exec.FetchVideo(
			testCameraUri,
			os.TempDir()+fmt.Sprintf("/testing_%d.mp4", i),
			"640x480",  // size
			"mp4",      // format
			"00:00:00", // position time
			5,          // 5 seconds
			20,         // 20 fps
		); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFetchGif(t *testing.T) {
	exec := NewCmd("ffmpeg").Debug(true)
	for i := 2; i > 0; i-- {
		if err := exec.FetchVideo(
			testCameraUri,
			os.TempDir()+fmt.Sprintf("/testing_%d.gif", i),
			"640x480",  // size
			"gif",      // format
			"00:00:00", // position time
			5,          // 5s
			10,         // 20 fps
		); err != nil {
			t.Fatal(err)
		}
	}
}
