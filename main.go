package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	path := "/Users/johnnash/Downloads/"
	err := filepath.Walk(path, LRCtoSRTConverter)
	if err != nil {
		fmt.Println(err)
	}
}

func LRCtoSRTConverter(path string, info os.FileInfo, err error) error {
	if path[len(path)-3:] == "lrc" {
		fmt.Println(info.Name())
		if err != nil {
			fmt.Println(err)
		}
		//打开源文件
		fmt.Println(path)
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = file.Close()
		}()
		//打开待写入文件
		fDst, err := os.OpenFile(path[:len(path)-3]+"srt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		defer func() {
			_ = fDst.Close()
		}()
		_, _ = fDst.Write([]byte{0xEF, 0xBB, 0xBF})
		//添加utf8 bom
		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K
		//scanner := bufio.NewScanner(file)
		//
		//const maxCapacity int = longLineLen  // your required line length
		//buf := make([]byte, maxCapacity)
		//scanner.Buffer(buf, maxCapacity)
		//增加一行的读取上限，在一行中有较多字符时使用
		scanner.Scan()
		count := 0
		for true {
			lineCurrent := scanner.Text()
			startTime := "00:" + lineCurrent[1:6] + "," + lineCurrent[7:9] + "0"
			if scanner.Scan() == false {
				timeDS := strings.Replace(lineCurrent[1:9], ":", "m", 1)
				timeDS += "s"
				timeD, _ := time.ParseDuration(timeDS) //将时间换为time.duration格式
				t3 := timeD + time.Second*5            //增加5s, 以下分别计算增加5s之后对应的时、分、秒和毫秒
				hours := int(t3.Hours())
				t3 -= time.Hour * time.Duration(hours)
				minutes := int(t3.Minutes())
				t3 -= time.Minute * time.Duration(minutes)
				seconds := int(t3.Seconds())
				t3 -= time.Second * time.Duration(seconds)
				milliseconds := int(t3.Milliseconds())
				endTime := fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
				_, _ = fDst.WriteString(fmt.Sprint(count) + "\n")
				_, _ = fDst.WriteString(startTime + " --> " + endTime + "\n")
				_, _ = fDst.WriteString(lineCurrent[10:] + "\n")
				break
			}
			lineNext := scanner.Text()
			endTime := "00:" + lineNext[1:6] + "," + lineNext[7:9] + "0"
			_, _ = fDst.WriteString(fmt.Sprint(count) + "\n")
			_, _ = fDst.WriteString(startTime + " --> " + endTime + "\n")
			_, _ = fDst.WriteString(lineCurrent[10:] + "\n\n")
			count++
		}
	}
	return nil
}
