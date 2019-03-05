package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	t := time.Now()
	hostname, _ := os.Hostname()
	filename := hostname + ".txt"
	out, _ := exec.Command("df", "-alk").Output()
	outlines := strings.Split(string(out), "\n")
	l := len(outlines)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, line := range outlines[1 : l-1] {
		parsedLine := strings.Fields(line)
		row := strconv.FormatInt(t.Unix(), 10) + "," + hostname + "," + strings.Join(parsedLine, ",")
		fmt.Fprintln(file, row)
	}
}
