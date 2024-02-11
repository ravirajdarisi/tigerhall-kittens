package utils

import (
	"image"
	"log"
	"math"

	"github.com/nfnt/resize"
)

// haversine formula calculates the distance between two points on the earth
// given their latitude and longitude in degrees.
func bckCalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371 // km

	// Convert latitude and longitude from degrees to radians.
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusKm * c

	return distance
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const earthRadiusKm = 6371 // km

    // Convert latitude and longitude from degrees to radians.
    lat1Rad := lat1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    lon1Rad := lon1 * math.Pi / 180
    lon2Rad := lon2 * math.Pi / 180

    // Haversine formula
    dLat := lat2Rad - lat1Rad
    dLon := lon2Rad - lon1Rad
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1Rad)*math.Cos(lat2Rad)*
            math.Sin(dLon/2)*math.Sin(dLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    distance := earthRadiusKm * c

    // Log the components of the distance calculation
    log.Printf("dLat: %v, dLon: %v, a: %v, c: %v", dLat, dLon, a, c)

    return distance
}




// ResizeImage takes an image and new width and height values,
// and returns a new image resized to the specified width and height.
func ResizeImage(img image.Image, newWidth, newHeight uint) image.Image {
	// Resize the image to the specified width and height using Lanczos resampling
	// and preserve aspect ratio.
	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return resizedImage
}
