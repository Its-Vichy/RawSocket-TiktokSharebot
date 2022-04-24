package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/zenthangplus/goccm"
)

var (
	domains = []string{"api19.tiktokv.com", "api.toutiao50.com", "api19.toutiao50.com", "api19-core-c-alisg.tiktokv.com"}

	sent   int = 0
	errors int = 0
	rpm    int = 0
	rps    int = 0
)

func get_device_id() string {
	var id string

	for i := 1; i <= 19; i++ {
		id += strconv.Itoa(rand.Intn(9-1) + 9)
	}

	return id
}

func connect_and_send(api_addr string, payload []byte) {
	conn, err := net.Dial("tcp", api_addr+":80")

	if err != nil {
		return
	}

	_, err = conn.Write(payload)

	if err != nil {
		errors++
	} else {
		sent++
	}
}

func rpmCounter() {
	for {
		before := sent
		time.Sleep(6000 * time.Millisecond)
		after := sent

		rpm = (after - before) * 10
	}
}

func rpsCounter() {
	for {
		before := sent
		time.Sleep(1 * time.Second)
		after := sent

		rps = (after - before) * 10
	}
}

func update_counter() {
	for {
		time.Sleep(time.Millisecond * 500)
		fmt.Printf("--> Sent: %d - err: %d | rpm: %d - rps: %d     \r", sent, errors, rpm, rps)
	}
}

func main() {
	var video_id string
	var threads int

	fmt.Print("[>] Threads: ")
	fmt.Scanln(&threads)

	fmt.Print("[>] ID: ")
	fmt.Scanln(&video_id)

	fmt.Print("\n\n")

	go rpmCounter()
	go rpsCounter()
	go update_counter()

	c := goccm.New(threads)

	for {
		c.Wait()

		go func() {
			addr := domains[rand.Intn(len(domains))]
			connect_and_send(addr, []byte(fmt.Sprintf("POST /aweme/v1/aweme/stats/?channel=tiktok_web&device_type=iPhone1,1&deviceID=161612101510139131516111091112141411&os_version=20&version_code=220400&appName=tiktok_web&device_platform=iphone&aid=1988 HTTP/1.1\r\nUser-Agent: Mozilla/5.0 (Linux; Android 7.1.1; Moto E (4) Plus) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.111 Mobile Safari/537.36\r\nHost: %s\r\nContent-Type: application/x-www-form-urlencoded; charset=UTF-8\r\nContent-Length: 41\r\n\r\nitem_id=%s&share_delta=1", addr, video_id)))
			c.Done()
		}()
	}

	c.WaitAllDone()
}
