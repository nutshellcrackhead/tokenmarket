package values

import "math"

func leftLegLogic(level int64, position int64) int64 {
	index := (position-1)*int64(math.Pow(float64(2), float64(level))) + 1
	return int64(index)
}

func rightLegLogic(level int64, position int64) int64 {
	index := position * int64(math.Pow(float64(2), float64(level)))
	return int64(index)
}

var Legs map[string]func(level int64, position int64) int64 = map[string]func(level int64, position int64) int64{
	"left":  leftLegLogic,
	"right": rightLegLogic,
}

func GetIndexForLeftLeg(level int64, position int64) int64 {
	return Legs["left"](level, position)
}

func GetIndexForRightLeg(level int64, position int64) int64 {
	return Legs["right"](level, position)
}

func StartEndPositionInStructure(level int64, position int64) (int64, int64) {
	return GetIndexForLeftLeg(level, position), GetIndexForRightLeg(level, position)
}
