package ffmpeg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"

	"clc.hmu/app/public/log"
	"github.com/gwaylib/errors"
)

var debug = true

type FFMPEG struct {
	debug   bool
	cmdPath string
}

func NewCmd(cmdPath string) *FFMPEG {
	return &FFMPEG{
		cmdPath: cmdPath,
	}
}

// 设置调试开关
func (ff *FFMPEG) Debug(on bool) *FFMPEG {
	ff.debug = on
	return ff
}

// 详细可用参数参考以下链接
// https://www.jianshu.com/p/049d03705a81
func (ff *FFMPEG) Exec(args ...string) error {
	ffCmd := exec.Command(ff.cmdPath, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	ffCmd.Stdout = &stdout
	ffCmd.Stderr = &stderr
	defer func() {
		if ff.debug {
			fmt.Printf("stdout:\n%s\n", stdout.String())
			fmt.Printf("stderr:\n%s\n", stderr.String())
		}
	}()

	if err := ffCmd.Run(); err != nil {
		log.Warning(errors.As(err))
	}
	return nil

}

//
// 从指定的uri中读取视频图片
//
// 参数
// uri -- 格式如："rtsp://admin:hyiot123@192.168.1.64:554/Streaming/Channels/101"
// toFile -- 格式如, 需确保目录是存在的："/tmp/1.jpg"
// toSize -- 视频尺寸，如：320*240
// toFormat -- 指定输出文件的格式，如: image2, mjpeg, gif
// posTime -- 视频的抓取时间点, 格式如：00:00:02
// vFrames -- 视频抓取时间点之后的帧数
//
// 返回
// error -- 若执行错误，返回此错误值
func (ff *FFMPEG) FetchImage(uri, toFile, toSize, toFormat, posTime string, vFrames int) error {
	return ff.Exec(
		"-i", uri,
		"-s", toSize,
		"-f", toFormat,
		"-ss", posTime,
		"-vframes", strconv.Itoa(vFrames),
		"-y", toFile,
	)
}

//
// 从指定的uri截取视频数据
//
// 参数
// uri -- 格式如："rtsp://admin:hyiot123@192.168.1.64:554/Streaming/Channels/101"
// toFile -- 格式如, 需确保目录是存在的："/tmp/1.jpg"
// toSize -- 视频尺寸，如：320*240
// toFormat -- 需要转出的格式, mp4, gif
// posTime -- 起始指定的时间 [-]hh:mm:ss[.xxx]的格式也支持
// duration -- 截取的时长, 单位为秒
// fps -- 需要转出的帧频
//
// 返回
// error -- 若执行错误，返回此错误值
func (ff *FFMPEG) FetchVideo(uri, toFile, toSize, toFormat, posTime string, duration int, fps int) error {
	return ff.Exec(
		"-i", uri,
		"-s", toSize,
		"-f", toFormat,
		"-ss", posTime,
		"-t", fmt.Sprintf("%d", duration),
		"-r", strconv.Itoa(fps),
		"-y", toFile,
	)
}
