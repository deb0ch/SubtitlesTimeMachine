package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func addTimeToTime(time string, offsetMilli int) string {
	hh := time[0:2]
	mm := time[3:5]
	ss := time[6:8]
	ms := time[9:12]
	h, _ := strconv.Atoi(hh)
	m, _ := strconv.Atoi(mm)
	s, _ := strconv.Atoi(ss)
	l, _ := strconv.Atoi(ms)
	if offsetMilli > 0 {
		l += offsetMilli
		s += l / 1000
		l %= 1000
		m += s / 60
		s %= 60
		h += m / 60
		m %= 60
	} else {
		l += offsetMilli
		if l < 0 {
			s += l/1000 - 1
			l = 1000 + l%1000
		}
		if s < 0 {
			m += s/60 - 1
			s = 60 + s%60
		}
		if m < 0 {
			h += m/60 - 1
			m = 60 + m%60
		}
		if h < 0 {
			fmt.Println("Error: negative times. Please try a smaller offset.")
			os.Exit(-1)
		}
	}
	if h >= 100 {
		h = 99
	}
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, l)
}

func addTimeToLine(line string, offsetMilli int) string {
	time1 := line[0:12]
	time2 := line[17:29]

	res := addTimeToTime(time1, offsetMilli)
	res += " --> "
	res += addTimeToTime(time2, offsetMilli) + line[29:] + "\r\n"

	return res
}

func processLine(line string, outFile *bytes.Buffer, offsetMilli int) {
	r, _ := regexp.Compile("^([0-9]{2}):([0-9]{2}):([0-9]{2}),([0-9]{3}) --> ([0-9]{2}):([0-9]{2}):([0-9]{2}),([0-9]{3})")
	if r.MatchString(line) {
		outFile.WriteString(addTimeToLine(line, offsetMilli))
	} else {
		outFile.WriteString(line + "\r\n")
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <file.srt> <time offset>\n", os.Args[0]) // Todo: Potential segfault if execve() was called without arguments in argv.
		return
	}
	file, err := os.OpenFile(os.Args[1], os.O_RDWR, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // By default, splits arround '\n' characters
	var newFile bytes.Buffer

	offsetMilli, _ := strconv.Atoi(os.Args[2])
	for scanner.Scan() {
		processLine(scanner.Text(), &newFile, offsetMilli)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.Write(bytes.Trim(newFile.Bytes(), "\x00"))
	if err != nil {
		fmt.Println(err)
		return
	}
}
