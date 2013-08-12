package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

var iface = flag.String("iface", "tun-exp", "Name of tun network interface")

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
	flag.Parse()
	if len([]byte(*iface)) > IFNAMSIZ-1 {
		log.Fatalf("Too long -iface name. It must be shorter than %d chars", IFNAMSIZ)
	}

	tun, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Opened tun device, f: %v\n", tun.Fd())

	req := ifreq{
		flags: IFF_TUN | IFF_NO_PI,
	}
	copy(req.name[:], []byte(*iface))

	errno := ioctl(tun.Fd(), TUNSETIFF, unsafe.Pointer(&req))
	if errno != 0 {
		log.Fatalf("Failed to ioctl %s device, errno: %d", *iface, errno)
	}
	fmt.Printf("ioctl succeeded!\n")

	if err = exec.Command("/sbin/ip", "link", "set", *iface, "up").Run(); err != nil {
		log.Fatal("Unable to ip link up: ", err)
	}

	if err = exec.Command("/sbin/ip", "addr", "add", "10.0.0.1/24", "dev", *iface).Run(); err != nil {
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
