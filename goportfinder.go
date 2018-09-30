// goportfinder : Port discovery tool in golang
// written by : @x1m_martijn
//
// Codebase example used: https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type PortFinder struct {
	ip   string
	lock *semaphore.Weighted
}

func Ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}
	
	s := strings.TrimSpace(string(out))
	
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	
	return i
}

func ScanPort(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		} else {
			//fmt.Println(port, "closed") // this gives a lot of spam
		}
		return
	}

	conn.Close()
	fmt.Println("[+] Found open port:", port, "on target:", ip)
	//fmt.Println("Total ports found:")
}

func (ps *PortFinder) Start(f, l int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.ip, port, timeout)
		}(port)
	}
}

func main() {
	var ip_string string
	flag.StringVar(&ip_string, "ip", "127.0.0.1", "insert ip")
	flag.Parse()
	ipaddress := net.ParseIP(ip_string)
	fmt.Println(ipaddress)

	ps := &PortFinder {
		ip:   ip_string, // input target IP, needs to be changed to user input (single or file)
		lock: semaphore.NewWeighted(Ulimit()),
	}
	ps.Start(1, 3000, 3000*time.Millisecond) // fiddling around with this a bit, 1, 3000 is the portrange and can be adjusted.

	// work in progress based on the script from Kent Gruber: https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d
}