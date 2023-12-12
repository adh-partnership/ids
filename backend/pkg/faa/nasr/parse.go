package nasr

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/adh-partnership/ids/backend/pkg/network"
	"github.com/adh-partnership/ids/backend/pkg/utils"
	"github.com/gocarina/gocsv"
)

type Airport struct {
	State              string  `csv:"STATE_CODE"`
	FAAID              string  `csv:"ARPT_ID"`
	ICAO               string  `csv:"ICAO_ID"`
	City               string  `csv:"CITY"`
	Name               string  `csv:"ARPT_NAME"`
	Latitude           float64 `csv:"LAT_DECIMAL"`
	Longitude          float64 `csv:"LONG_DECIMAL"`
	Elevation          float32 `csv:"ELEVATION"`
	MagnaticVariation  int     `csv:"MAG_VARN"`
	MagneticHemisphere string  `csv:"MAG_HEMIS"`
}

func ProcessAirports() (map[string]*Airport, error) {
	airports := make(map[string]*Airport)

	data, err := downloadAndExtract()
	if err != nil {
		return nil, fmt.Errorf("unable to download and extract NASR data: %s", err)
	}

	// Parse the CSV file
	var apts []Airport
	if err := gocsv.UnmarshalBytes(data, &apts); err != nil {
		return nil, fmt.Errorf("unable to parse NASR data: %s", err)
	}

	for _, apt := range apts {
		airports[apt.FAAID] = &apt
	}

	return airports, nil
}

func downloadAndExtract() ([]byte, error) {
	// Download the file
	url := URL(GetCurrentCycle(utils.Now()), "APT")

	data, err := network.DownloadFile(url)
	if err != nil {
		return nil, fmt.Errorf("unable to download file %s: %s", url, err)
	}

	// We need to extract APT_BASE.csv from the zip file in the data byte slice
	// convert byte slice to reader
	reader := bytes.NewReader(data)
	z, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("unable to read zip file: %s", err)
	}

	// Find the APT_BASE.csv file
	for _, f := range z.File {
		if f.Name == "APT_BASE.csv" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("unable to open file %s: %s", f.Name, err)
			}
			defer rc.Close()
			// Read and return the file
			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("unable to find APT_BASE.csv in zip file")
}
