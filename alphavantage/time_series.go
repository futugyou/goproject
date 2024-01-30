package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type TimeSeries struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func ReadTimeSeries(r io.Reader) ([]*TimeSeries, error) {
	reader := csv.NewReader(r)
	reader.ReuseRecord = true
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	fmt.Println(1111)
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	result := make([]*TimeSeries, 0)

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
		value, err := readTimeSeriesItem(record)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil

}

func readTimeSeriesItem(s []string) (*TimeSeries, error) {
	const (
		timestamp = iota
		open
		high
		low
		close
		volume
	)

	value := &TimeSeries{}

	d, err := time.Parse("2006-01-02 15:04:05", s[timestamp])
	if err != nil {
		return nil, fmt.Errorf("error parsing timestamp %s", s[timestamp])
	}
	value.Time = d

	f, err := strconv.ParseFloat(s[open], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing open %s", s[open])
	}
	value.Open = f

	f, err = strconv.ParseFloat(s[high], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing high %s", s[high])
	}
	value.High = f

	f, err = strconv.ParseFloat(s[low], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing low %s", s[low])
	}
	value.Low = f

	f, err = strconv.ParseFloat(s[close], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing close %s", s[close])
	}
	value.Close = f

	f, err = strconv.ParseFloat(s[volume], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing volume %s", s[volume])
	}
	value.Volume = f

	return value, nil
}
