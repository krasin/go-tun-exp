### Linux TUN device explorations.

For the explanation what happens here, please, read this awesome post: [Tun/Tap interface tutorial](http://backreference.org/2010/03/26/tuntap-interface-tutorial/).

To use this (assuming that you already have Go 1.1.1 or newer installed):

```
sudo apt-get install sendip openvpn
git clone https://github.com/krasin/go-tun-exp
cd go-tun-exp
./create-iface.sh # will ask for sudo permissions
go build -o go-tun-exp main.go
./go-tun-exp
```

Let's send a UDP packet from a different terminal:

```
sudo sendip -p ipv4 -is 192.168.1.81 -p udp -us 5070 -ud 12233 -d "Hello" -v 10.0.0.2
```

You should see the following output from go-tun-exp:

```
2013/08/11 13:43:31 Read 33 bytes from tun:
00000000  45 00 00 21 84 ae 00 00  ff 11 6b 22 c0 a8 01 51  |E..!......k"...Q|
00000010  0a 00 00 02 13 ce 2f c9  00 0d cc 6f 48 65 6c 6c  |....../....oHell|
00000020  6f                                                |o|
```

You can also try to ping 10.0.0.2 and see ICMP packets:

```
2013/08/11 13:44:30 Read 84 bytes from tun:
00000000  45 00 00 54 00 00 40 00  40 01 26 a7 0a 00 00 01  |E..T..@.@.&.....|
00000010  0a 00 00 02 08 00 b1 ff  69 7b 00 01 ae f7 07 52  |........i{.....R|
00000020  00 00 00 00 64 67 03 00  00 00 00 00 10 11 12 13  |....dg..........|
00000030  14 15 16 17 18 19 1a 1b  1c 1d 1e 1f 20 21 22 23  |............ !"#|
00000040  24 25 26 27 28 29 2a 2b  2c 2d 2e 2f 30 31 32 33  |$%&'()*+,-./0123|
00000050  34 35 36 37                                       |4567|
```

When you're done, stop go-tun-exp and delete tun2 interface:

```
./delete-iface.sh
```

That's all very hacky now, sorry about that.



