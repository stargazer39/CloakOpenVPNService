package main

import (
	"flag"
	"log"
	"net"
	"time"
)

func main() {
	cloakHost := flag.String("cloak-host", "127.0.0.1", "Cloak IP")
	cloakPort := flag.String("cloak-port", "443", "Cloak port")
	cloakConfig := flag.String("cloak-config", "ckclient-zoom.json", "Cloak config")
	openvpnConfig := flag.String("openvpn-config", "profile.ovpn", "Open VPN config")

	flag.Parse()

	defaultTimeout := time.Second * 5

	// Two services
	cloakProc := NewService("ck-client", "-c", *cloakConfig, "-s", *cloakHost)
	openvpnProc := NewService("openvpn", *openvpnConfig)

	log.Println("Starting cloak + openvpn")

	// Start checking in seperate thread
	aliveChan := make(chan bool)
	boost := true

	go checkAliveThread(aliveChan, &boost, *cloakHost, *cloakPort, defaultTimeout)
	for {
		alive := <-aliveChan

		if alive {
			// log.Printf("Host is alive, trying to start processes\n")
		startcloak:
			if !cloakProc.Running {
				log.Printf("Cloak isn't runnig, trying to start\n")
				sErr := cloakProc.Start()
				check(sErr)
			}

			if !openvpnProc.Running {
				log.Printf("OpenVPN isn't runnig, trying to start\n")
				if !cloakProc.Running {
					log.Println("Turns out, cloak is dead, go to startcloak")
					goto startcloak
				}
				oErr := openvpnProc.Start()
				check(oErr)
			}
		} else { //log.Println("Killing procesess")
			if openvpnProc.Running {
				// log.Println("Killing openvpn")
				openvpnProc.Kill()
			}

			if cloakProc.Running {
				cloakProc.Kill()
			}
		}
	}
}

func check(e error) {
	if e != nil {
		log.Panicln(e)
	}
}

func checkAliveThread(alive chan bool, boost *bool, host string, port string, timeout time.Duration) {
	for {
		if err := isHostAlive(host, port, timeout); err != nil {
			log.Printf("Host %s is not reachable\n", host)
			alive <- false
		} else {
			log.Printf("Host %s is alive\n", host)
			alive <- true
		}
		time.Sleep(time.Second * 5)
	}
}

func isHostAlive(host string, port string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)

	if err != nil {
		return err
	}

	if conn != nil {
		defer conn.Close()
	}
	return nil
}
