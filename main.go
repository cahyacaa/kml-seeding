package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

const (
	Route1 string = "route-1.kml"
	Route2 string = "route-2.kml"
	Route3 string = "route-3.kml"
	Route4 string = "route-4.kml"
	Route5 string = "route-5.kml"
)

var dirs = []string{
	"route-1.kml",
	"route-2.kml",
	"route-3.kml",
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

func processKML(filePath, key string, latLongChan chan<- map[string][]float64) error {
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

					var data []float64
					data = append(data, []float64{floatLat, floatLong}...)
					latLongChan <- map[string][]float64{
						key: data,
					}
				}
			}
		}
	}
	return nil
}

func main() {
	var g errgroup.Group
	latLongChan := make(chan map[string][]float64)
	var latLongPairs = make(map[string][][]float64)

	for _, fileName := range dirs {
		dir := "./files/" + fileName
		g.Go(func() error {
			return processKML(dir, fileName, latLongChan)
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			fmt.Println("Error processing KML files:", err)
		}
		close(latLongChan)
	}()

	for pair := range latLongChan {
		for k, v := range pair {
			switch k {
			case Route1, Route2, Route3, Route4, Route5:
				latLongPairs[k] = append(latLongPairs[k], v)
			default:
				fmt.Println("Warning processing KML files : there is lat long not belongs to any kml file")
			}
		}
	}

	fmt.Println(len(latLongPairs[Route1]), len(latLongPairs[Route2]), len(latLongPairs[Route3]), len(latLongPairs[Route4]), len(latLongPairs[Route5]))
}
