### Linux TUN device explorations

For the explanation what happens here, please, read this awesome post: [Tun/Tap interface tutorial](http://backreference.org/2010/03/26/tuntap-interface-tutorial/).

To use this (assuming that you already have Go 1.1.1 or newer installed):

```
go get github.com/krasin/go-tun-exp
sudo `which go-tun-exp`
```

At this point we have tun-exp network interface:
```
$ ifconfig
...
tun-exp   Link encap:UNSPEC  HWaddr 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00  
          inet addr:10.0.0.1  P-t-P:10.0.0.1  Mask:255.255.255.0
          UP POINTOPOINT RUNNING NOARP MULTICAST  MTU:1500  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:500 
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
```

Btw, you can change the interface name by using -iface option to go-tun-exp.
The name must be shorter than 16 symbols (Linux kernel restriction, see IFNAMSIZ definition).

Let's try to ping 10.0.0.2 (from a different terminal):

```
2013/08/11 13:44:30 Read 84 bytes from tun:
00000000  45 00 00 54 00 00 40 00  40 01 26 a7 0a 00 00 01  |E..T..@.@.&.....|
00000010  0a 00 00 02 08 00 b1 ff  69 7b 00 01 ae f7 07 52  |........i{.....R|
00000020  00 00 00 00 64 67 03 00  00 00 00 00 10 11 12 13  |....dg..........|
00000030  14 15 16 17 18 19 1a 1b  1c 1d 1e 1f 20 21 22 23  |............ !"#|
00000040  24 25 26 27 28 29 2a 2b  2c 2d 2e 2f 30 31 32 33  |$%&'()*+,-./0123|
00000050  34 35 36 37                                       |4567|
```

Let's send a UDP packet from a different terminal (you will need to ```sudo apt-get install sendip```:

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




