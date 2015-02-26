package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dcarley/dymoscale"
)

const (
	outputResult = "Result: %+v"
	outputError  = "Error: %s"
)

var (
	ValidModes = []string{"raw", "parsed", "grams"}
)

func main() {
	var (
		mode     = flag.String("mode", "grams", fmt.Sprintf("Mode to print: %v", ValidModes))
		interval = flag.Duration("interval", time.Second, "Interval between reads")
	)

	flag.Parse()
	if err := validateMode(*mode); err != nil {
		log.Fatalln(err)
	}

	scale, err := dymoscale.NewScale()
	if err != nil {
		log.Fatalln(err)
	}
	defer scale.Close()

	for {
		log.Println(getResultOrError(scale, *mode))
		time.Sleep(*interval)
	}
}

func validateMode(mode string) error {
	for _, validMode := range ValidModes {
		if mode == validMode {
			return nil
		}
	}

	return fmt.Errorf("invalid mode: %s", mode)
}

func getResultOrError(scale dymoscale.Scaler, mode string) string {
	var result string
	var err error

	switch mode {
	case "grams":
		var parsed dymoscale.Measurementer
		var payload int

		if parsed, err = scale.ReadMeasurement(); err == nil {
			payload, err = parsed.Grams()
			result = fmt.Sprintf(outputResult, payload)
		}
	case "parsed":
		var payload dymoscale.Measurementer

		payload, err = scale.ReadMeasurement()
		result = fmt.Sprintf(outputResult, payload)
	case "raw":
		var payload []byte

		payload, err = scale.ReadRaw()
		result = fmt.Sprintf(outputResult, payload)
	}

	if err != nil {
		return fmt.Sprintf(outputError, err)
	}

	return result
}
