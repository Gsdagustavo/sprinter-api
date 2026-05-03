package rules

import (
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
)

func ValidateActivityPoints(points *[]entities.Point) error {
	for _, point := range *points {
		if point.Latitude == 0 || point.Longitude == 0 {
			return derr.InvalidGeoPoint
		}
	}

	return nil
}
