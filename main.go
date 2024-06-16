package main

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"strconv"
	"strings"
)

var dirs = []string{
	"route-1.kml",
	"route-2.kml",
	"route-4.kml",
	"route-5.kml",
}

type KML struct {
	XMLName  xml.Name `xml:"kml"`
	Document Document `xml:"Document"`
}

type Document struct {
	Placemark []Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name        string     `xml:"name"`
	Description string     `xml:"description"`
	Point       Point      `xml:"Point"`
	LineString  LineString `xml:"LineString"`
}

type Point struct {
	Coordinates string `xml:"coordinates"`
}

type LineString struct {
	Coordinates string `xml:"coordinates"`
}

func processKML(filePath string, latLongChan chan<- []float64) error {
	// Open the KML file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	// Read the file content
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	// Unmarshal the XML data
	var kml KML
	err = xml.Unmarshal(byteValue, &kml)
	if err != nil {
		return fmt.Errorf("error unmarshaling XML from file %s: %v", filePath, err)
	}

	// Extract coordinates
	var coordinatesList []string
	for _, placemark := range kml.Document.Placemark {
		if placemark.Point.Coordinates != "" {
			coordinatesList = append(coordinatesList, placemark.Point.Coordinates)
		}
		if placemark.LineString.Coordinates != "" {
			coordinatesList = append(coordinatesList, placemark.LineString.Coordinates)
		}
	}

	// Process coordinates
	for _, coordString := range coordinatesList {
		coordPairs := strings.Split(strings.TrimSpace(coordString), " ")
		for _, pair := range coordPairs {
			if pair != "" {
				coords := strings.Split(pair, ",")
				if len(coords) >= 2 {
					lat, lon := coords[1], coords[0]
					floatLat, err := strconv.ParseFloat(lat, 64)
					if err != nil {
						return fmt.Errorf("error parsing latitude %s: %v", lat, err)
					}
					floatLong, err := strconv.ParseFloat(lon, 64)
					if err != nil {
						return fmt.Errorf("error parsing longitude %s: %v", lon, err)
					}
					latLongChan <- []float64{floatLat, floatLong}
				}
			}
		}
	}
	return nil
}

func main() {
	var g errgroup.Group
	latLongChan := make(chan []float64)

	for _, dir := range dirs {
		dir := "./files/" + dir // create new instance for the closure
		g.Go(func() error {
			return processKML(dir, latLongChan)
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			fmt.Println("Error processing KML files:", err)
		}
		close(latLongChan)
	}()

	var latLongPairs [][]float64
	for pair := range latLongChan {
		latLongPairs = append(latLongPairs, pair)
		fmt.Printf("Lat: %f, Lon: %f\n", pair[0], pair[1])
	}

	//fmt.Println(latLongPairs)
}
