package main

import (
	_ "fmt"
	"net"
	_ "strconv"
	_ "strings"
	_ "time"
)

Pooja.  We can do $65/hr C2C remote for a minimum of 2 months, then $75/hr C2C onsite Irving

type Tvpnpacket []byte

//This code is to process our VPN encrypted traffic
//When a program opens the tun0 ip address, this traffic is sent to this program
//so it can send it onward

func (gpack *Tvpnpacket) IPver() int {
	if 4 == ((*gpack)[2] >> 4) {
		return 4
	}
	if 6 == ((*gpack)[2] >> 4) {
		return 6
	}
	return 0
}

func (gpack *Tvpnpacket) Dst() [4]byte {
	return [4]byte{(*gpack)[18], (*gpack)[19], (*gpack)[20], (*gpack)[21]}
}

func (gpack *Tvpnpacket) DstV4() net.IP {
	return net.IPv4((*gpack)[18], (*gpack)[19], (*gpack)[20], (*gpack)[21])
}

func (gpack *Tvpnpacket) Src() [4]byte {
	return [4]byte{(*gpack)[14], (*gpack)[15], (*gpack)[16], (*gpack)[17]}
}

func (gpack *Tvpnpacket) IsMulticast() bool {
	return ((*gpack)[18] > 223) && ((*gpack)[18] < 240)
}

//func (gvpn *Tvpn) Trafficrx(proto string, port int, iface *water.Interface) {
func (gvpn *Tvpn) Trafficrx() {

	/*
		conn, err := reuseport.NewReusableUDPPortConn(proto, fmt.Sprintf(":%v", port))
		if nil != err {
			log.Fatalln("Unable to get UDP socket:", err)
		}

		buf := make([]byte, BUFFERSIZE)
		for {
			n, _, err := conn.ReadFrom(buf)

			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}

			// ReadFromUDP can return 0 bytes on timeout
			if 0 == n {
				continue
			}

			if n%aes.BlockSize != 0 {
				fmt.Println("packet size ", n, " is not a multiple of the block size")
				continue
			}

			iv := buf[:aes.BlockSize]
			ciphertext := buf[aes.BlockSize:n]

			conf := config.Load().(VPNState)

			mode := cipher.NewCBCDecrypter(conf.Main.block, iv)

			var size int

			if conf.Main.hasalt {

				// if we have alternative key we need store orig packet for second try

				pcopy := make([]byte, n)
				copy(pcopy, buf[:n])

				mode.CryptBlocks(ciphertext, ciphertext)

				size = int(ciphertext[0]) + (256 * int(ciphertext[1]))
				if (n-aes.BlockSize-2)-size > 16 || (n-aes.BlockSize-2)-size < 0 || 4 != ((ciphertext)[2]>>4) {
					// don't looks like anything is ok, trying second key

					copy(buf[:n], pcopy)
					cipher.NewCBCDecrypter(conf.Main.altblock, iv).CryptBlocks(ciphertext, ciphertext)

					size = int(ciphertext[0]) + (256 * int(ciphertext[1]))
					if (n-aes.BlockSize-2)-size > 16 || (n-aes.BlockSize-2)-size < 0 || 4 != ((ciphertext)[2]>>4) {
						fmt.Println("Invalid size field or IPv4 id in decrypted message", size, (n - aes.BlockSize - 2))
						continue
					}
				}

			} else {

				mode.CryptBlocks(ciphertext, ciphertext)

				size = int(ciphertext[0]) + (256 * int(ciphertext[1]))
				if (n-aes.BlockSize-2)-size > 16 || (n-aes.BlockSize-2)-size < 0 {
					fmt.Println("Invalid size field in decrypted message", size, (n - aes.BlockSize - 2))
					continue
				}

				if 4 != ((ciphertext)[2] >> 4) {
					fmt.Println("Non IPv4 packet after decryption, possible corupted packet")
					continue
				}
			}

			iface.Write(ciphertext[2 : 2+size])

		}
	*/
}

//func (gvpn *Tvpn) Traffictx(conn *net.UDPConn, iface *water.Interface) {
func (gvpn *Tvpn) Traffictx() {

	/*
		// first time fill with random numbers
		ivbuf := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, ivbuf); err != nil {
			log.Fatalln("Unable to get rand data:", err)
		}

		var packet Tvpnpacket = make([]byte, BUFFERSIZE)
		for {
			plen, err := iface.Read(packet[2 : MTU+2])
			if err != nil {
				break
			}

			if 4 != packet.IPver() {
				header, _ := ipv4.ParseHeader(packet[2:])
				log.Printf("Non IPv4 packet [%+v]\n", header)
				continue
			}

			// each time get pointer to (probably) new config
			c := config.Load().(VPNState)

			dst := packet.Dst()

			wanted := false

			addr, ok := c.remotes[dst]

			if ok {
				wanted = true
			}

			if dst == c.Main.bcastIP || packet.IsMulticast() {
				wanted = true
			}

			// very ugly and useful only for a limited numbers of routes!
			if !wanted {
				ip := packet.DstV4()
				for n, s := range c.routes {
					if n.Contains(ip) {
						addr = s
						ok = true
						wanted = true
						break
					}
				}
			}

			if wanted {
				// store orig packet len
				packet[0] = byte(plen % 256)
				packet[1] = byte(plen / 256)

				// encrypt
				clen := plen + 2

				if clen%aes.BlockSize != 0 {
					clen += aes.BlockSize - (clen % aes.BlockSize)
				}

				if clen > len(packet) {
					log.Println("clen > len(package)", clen, len(packet))
					continue
				}

				ciphertext := make([]byte, aes.BlockSize+clen)
				iv := ciphertext[:aes.BlockSize]

				copy(iv, ivbuf)

				mode := cipher.NewCBCEncrypter(c.Main.block, iv)
				mode.CryptBlocks(ciphertext[aes.BlockSize:], packet[:clen])

				// save new iv
				copy(ivbuf, ciphertext[clen-aes.BlockSize:])

				if ok {
					n, err := conn.WriteToUDP(ciphertext, addr)
					if nil != err {
						log.Println("Error sending package:", err)
					}
					if n != len(ciphertext) {
						log.Println("Only ", n, " bytes of ", len(ciphertext), " sent")
					}
				} else {
					// multicast or broadcast
					for _, addr := range c.remotes {
						n, err := conn.WriteToUDP(ciphertext, addr)
						if nil != err {
							log.Println("Error sending package:", err)
						}
						if n != len(ciphertext) {
							log.Println("Only ", n, " bytes of ", len(ciphertext), " sent")
						}
					}
				}
			} else {
				fmt.Println("Unknown dst", dst)
			}
		}
	*/

}
