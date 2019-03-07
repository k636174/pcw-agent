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
	filename := hostname + ".tmp"
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
		row := strconv.FormatInt(t.Unix(), 10) + "," + hostname + "," + strings.Join(parsedLine, ",") + ",diskusage"
		fmt.Fprintln(file, row)
	}

	out2, _ := exec.Command("sh", "-c", "top -l 1 | head -n 10").Output()
	outlines2 := strings.Split(string(out2), "\n")
	l2 := len(outlines2)
	for _, line := range outlines2[1 : l2-0] {
		parsedLine := strings.Split(line, ":")
		if 1 != len(parsedLine) {
			parsedLine2 := strings.Split(parsedLine[1], ",")
			row := strconv.FormatInt(t.Unix(), 10) + "," + hostname + "," + strings.Join(parsedLine2, ",") + "," + parsedLine[0]
			fmt.Fprintln(file, row)
		}

	}

}
