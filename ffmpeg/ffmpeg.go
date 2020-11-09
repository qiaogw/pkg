//流媒体处理，以便获取长度、截图等。
package ffmpeg

import (
	"github.com/astaxie/beego"
	//	"bufio"

	"os"
	"os/exec"
	"path/filepath"

	// "io"
	// "io/ioutil"
	"strconv"
	"strings"

	"github.com/qiaogw/pkg/pathtool"
)

var duration = 0
var allRes = ""
var lastPer = -1

//durToSec 时间字符串转换为秒数
func durToSec(dur string) (sec int) {
	durAry := strings.Split(dur, ":")
	if len(durAry) != 3 {
		return
	}
	hr, _ := strconv.Atoi(durAry[0])
	sec = hr * (60 * 60)
	min, _ := strconv.Atoi(durAry[1])
	sec += min * (60)
	second, _ := strconv.Atoi(durAry[2])
	sec += second
	return
}

//getRatio 获取视频文件的时长，返回时长秒数和时长字符串00:00:00
func getRatio(res string) (duration int, dur string) {
	i := strings.Index(res, "Duration")
	//	dur := ""
	if i >= 0 {

		dur = res[i+10:]
		if len(dur) > 8 {
			dur = dur[0:8]

			duration = durToSec(dur)
			beego.Debug("duration:", duration)
			allRes = ""
		}
	}
	beego.Debug("durdurdur:", dur)
	if duration == 0 {
		return
	}
	i = strings.Index(res, "time=")
	if i >= 0 {

		time := res[i+5:]
		if len(time) > 8 {
			time = time[0:8]
			sec := durToSec(time)
			per := (sec * 100) / duration
			if lastPer != per {
				lastPer = per
				beego.Debug("Percentage:", per)
			}
			allRes = ""
		}
	}
	return
}

//获取时长字符串
func getRat(res string) (dur string) {
	i := strings.Index(res, "Duration")
	//	dur := ""
	if i >= 0 {

		dur = res[i+10:]
		if len(dur) > 8 {
			dur = dur[0:8]

			duration = durToSec(dur)
			beego.Debug("duration:", duration)
			allRes = ""
		}
	}
	beego.Debug("durdurdur:", dur)

	return
}

//FfmpegActImag 执行 ffmpeg -ss 00:00:15 -i *.flv -f image2 -y *.jpg  2>&1 | grep Duration
//截取视频图片，获取视频长度
//输入流媒体文件名称
//输出错误、视频长度秒数、视频时长字符串、视频第3秒截图
func FfmpegActImag(file string) (err error, duration int, dur, img string) {
	_, err = os.Stat(file)

	if err != nil {
		beego.Error("file not found!", err)
		return
	}
	filename := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	filedir := filepath.Dir(file)
	imagdir := filedir + "/"
	_, er := os.Stat(imagdir)
	if er != nil {
		beego.Notice("dir not found")
		os.Mkdir(imagdir, os.ModePerm)
	}
	img = imagdir + filename + ".jpg"
	cmdName := "ffmpeg -ss 00:00:3 -i " + file + " -f image2 -y -s 320*240 " + img + "   2>&1 | grep Duration"
	//	log.Debug(cmdName)
	cmd := exec.Command("sh", "-c", cmdName)
	stdout, _ := cmd.StdoutPipe()
	//	log.Notice("cmd  is :", cmd)
	cmd.Start()
	oneByte := make([]byte, 8)
	for {
		_, err = stdout.Read(oneByte)
		if err != nil {
			beego.Debug("文件错误：", err.Error())
			break
		}
		allRes += string(oneByte)
		duration, dur = getRatio(allRes)
		if duration > 0 {
			break
		}
	}
	cmd.Wait()
	return
}

//FfmpegJoin 执行 	ffmpeg -f concat -i filelist.txt -c copy output
//合并视频文件
//输出错误、合并后视频文件信息
func FfmpegJoin(filename string) (outfile string, err error) {
	// var f *os.File

	filedir, _ := filepath.Split(filename)
	beego.Debug(filename, filedir)
	outfile = filedir + "outfile.flv"

	if pathtool.CheckFileIsExist(outfile) { //如果文件存在
		os.Remove(outfile) //创建文件
	}
	if err != nil {
		beego.Error(err)
		return
	}
	cmd := exec.Command("ffmpeg", "-f", "concat", "-i", filename, "-c", "copy", outfile)
	stdout, err := cmd.StdoutPipe()
	// 保证关闭输出流
	defer stdout.Close()
	if err != nil {
		beego.Error(err)
		return
	}
	// 运行命令
	if err = cmd.Run(); err != nil {
		beego.Error(err)
		return
	}
	return
}
