package utils

import (
    "math"
    "time"
)


var TICKS_PER_SECOND int = 10
var CELL_SIZE int = 50
var CELLS_PER_SECOND int = 1
var ATTACK_RANGE int = 15
var MS_PER_TICK time.Duration = time.Duration(1000 / TICKS_PER_SECOND) * time.Millisecond
var UNIT_VECTOR_MODIFIER float64 = float64(CELLS_PER_SECOND * CELL_SIZE / TICKS_PER_SECOND)


func FindVectorDegree(x_start float64, y_start float64, x_finish float64, y_finish float64) float64 {
    vector_length := FindDirectionVectorLength(x_start, y_start, x_finish, y_finish)
    norm_x := x_finish - x_start
    norm_y := y_finish - y_start
    direction_cos := (0 * norm_x - 1 * norm_y) / (1 * vector_length)

    if (math.Signbit(norm_x)) {
        return -1 * ToDegree(math.Acos(direction_cos))
    } else {
        return ToDegree(math.Acos(direction_cos))
    }
}


func ChangeCoordinatesSystemAndRotate(x float64, y float64, x_center float64, y_center float64, cos_ float64, sin_ float64) (float64, float64) {
    // convert coordinate system. from global to local
    // rotation center is a formation center
    // rotate formation relative to formation center
    x_, y_ := RotateInCurrentSystem(x-x_center, y-y_center, cos_, sin_)

    // convert coordinate system. from local to global
    return x_ + x_center, y_ + y_center
}


func RotateInCurrentSystem(x float64, y float64, cos_ float64, sin_ float64) (float64, float64) {
    // rotate vectors relative to "current" coordinates system center
    // x=0 and y=0 is center of current coordinates system
    return x * cos_ - y * sin_, x * sin_ + y * cos_
}


func FindDirectionUnitVector(x_start float64, y_start float64, x_finish float64, y_finish float64) (float64, float64) {
    x_, y_ := FindDirectionVector(x_start, y_start, x_finish, y_finish)
    direction_vector_length := FindDirectionVectorLength(x_start, y_start, x_finish, y_finish)

    return x_ / direction_vector_length, y_ / direction_vector_length
}


func FindDirectionVector(x_start float64, y_start float64, x_finish float64, y_finish float64) (float64, float64) {
    x_direction := x_finish - x_start
    y_direction := y_finish - y_start

    return x_direction, y_direction
}


func FindDirectionVectorLength(x_start float64, y_start float64, x_finish float64, y_finish float64) float64 {
    x_direction := x_finish - x_start
    y_direction := y_finish - y_start

    return math.Sqrt(math.Pow(x_direction, 2) + math.Pow(y_direction, 2))
}


func NormalizeDirectionAngle(angle float64) float64 {
    return angle - (math.Ceil(angle / 360) - 1) * 360
}


func ToRadians(degree float64) float64 {
    return degree * math.Pi / 180
}


func ToDegree(radians float64) float64 {
    return radians * 180 / math.Pi
}


func MaxInt(a int, b int) int {
    if (a > b) {
        return a
    }

    return b
}
