package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// diskusage
	out, _ := exec.Command("df", "-alk").Output()
	outlines := strings.Split(string(out), "\n")
	l := len(outlines)
	for _, line := range outlines[1 : l-1] {
		parsedLine := strings.Fields(line)
		row := strconv.FormatInt(t.Unix(), 10) + "," + hostname + "," + strings.Join(parsedLine, ",") + ",diskusage"
		fmt.Fprintln(file, row)
	}

	// top_result(load_average)
	out2, _ := exec.Command("sh", "-c", "top -l 1 | head -n 10").Output()
	outlines2 := strings.Split(string(out2), "\n")
	for _, line := range outlines2[1:len(outlines2)] {
		parsedLine := strings.Split(line, ":")
		if 1 != len(parsedLine) {
			parsedLine2 := strings.Split(parsedLine[1], ",")
			row := strconv.FormatInt(t.Unix(), 10) + "," + hostname + "," + strings.Join(parsedLine2, ",") + "," + parsedLine[0]
			fmt.Fprintln(file, row)
		}
	}

	// Local Ip Address
	out3, _ := exec.Command("sh", "-c", "ifconfig -a | grep inet | grep -v inet6 | grep -v 127.0.0.1").Output()
	outlines3 := strings.Split(string(out3), "\n")
	localip := strings.Fields(outlines3[0])
	src_lip := localip[1]

	// form values
	values := url.Values{}
	values.Add("hostname", hostname)
	values.Add("src_lip",src_lip)
	values.Encode()

	pcw_host := os.Getenv("PCW_HOST")
	res, err := http.PostForm(pcw_host + "/api/heartbeat", values)
	if err != nil {
		log.Fatal(err)
	}

	// header
	fmt.Printf("[status] %d\n", res.StatusCode)
	for k, v := range res.Header {
		fmt.Print("[header] " + k)
		fmt.Println(": " + strings.Join(v, ","))
	}

	// body
	defer res.Body.Close()
	body, error := ioutil.ReadAll(res.Body)
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println("[body] " + string(body))

}
