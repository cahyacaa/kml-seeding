package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	// Open the KML file
	file, err := os.Open("./files/Barito River Coordinate Point.kml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file content
	byteValue, _ := io.ReadAll(file)

	// Unmarshal the XML data
	var kml KML
	err = xml.Unmarshal(byteValue, &kml)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return
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
	var latLongPairs [][]float64
	for _, coordString := range coordinatesList {
		coordPairs := strings.Split(strings.TrimSpace(coordString), " ")
		for _, pair := range coordPairs {
			if pair != "" {
				coords := strings.Split(pair, ",")
				if len(coords) >= 2 {
					lat, lon := coords[1], coords[0]
					floatLat, _ := strconv.ParseFloat(lat, 64)
					floatLong, _ := strconv.ParseFloat(lat, 64)
					fmt.Printf("Lat: %s, Lon: %s\n", lat, lon)
					latLongPairs = append(latLongPairs, []float64{floatLat, floatLong})
				}
			}
		}
	}

	fmt.Println(latLongPairs)
}
