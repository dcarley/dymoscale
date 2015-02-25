package main

import (
	"flag"
	"log"
	"time"

	"github.com/dcarley/dymoscale"
)

var (
	interval = flag.Duration("interval", time.Second, "Interval between reads")
)

func main() {
	flag.Parse()

	scale, err := dymoscale.NewScale()
	if err != nil {
		log.Fatalln(err)
	}
	defer scale.Close()

	for {
		payload, err := scale.ReadRaw()
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("%+v\n", payload)
		}

		time.Sleep(*interval)
	}
}
