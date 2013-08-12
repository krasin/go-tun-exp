package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	IFNAMSIZ = 16

	IFF_NO_PI = 0x1000
	IFF_TUN   = 0x0001
)
const TUNSETIFF = 1074025674

type ifreq struct {
	name  [IFNAMSIZ]byte
	flags uint16
}

func ioctl(fd uintptr, req int, data unsafe.Pointer) (err syscall.Errno) {
	_, _, err = syscall.RawSyscall(syscall.SYS_IOCTL, fd, uintptr(req), uintptr(data))
	return
}

func main() {
	tun, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Opened tun device, f: %v\n", tun.Fd())

	req := ifreq{
		name:  [16]byte{'t', 'u', 'n', '2', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		flags: IFF_TUN | IFF_NO_PI,
	}
	errno := ioctl(tun.Fd(), TUNSETIFF, unsafe.Pointer(&req))
	if errno != 0 {
		log.Fatalf("Failed to ioctl tun device, errno: %d", errno)
	}
	fmt.Printf("ioctl succeeded!\n")

	if err = exec.Command("/sbin/ip", "link", "set", "tun2", "up").Run(); err != nil {
		log.Fatal("Unable to ip link up: ", err)
	}

	if err = exec.Command("/sbin/ip", "addr", "add", "10.0.0.1/24", "dev", "tun2").Run(); err != nil {
		log.Fatal("Unable to set ipv4 addr: ", err)
	}

	buf := make([]byte, 4096)
	for {
		n, err := tun.Read(buf)
		if err != nil {
			log.Fatal("Error reading from tun fd: %v", err)
		}
		toshow := n
		suffix := ""
		if toshow > 200 {
			toshow = 200
			suffix = "... <truncated> ...\n"
		}
		log.Printf("Read %d bytes from tun:\n%s%s", n, hex.Dump(buf[:toshow]), suffix)
	}
}
