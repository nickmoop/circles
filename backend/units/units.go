package units

import (
    "strconv"

    "backend/utils"
)


type Unit struct {
    Id string
    Color string
    Size int
    X float64
    Y float64
    ToX float64
    ToY float64
    SquadName string
}


func (unit *Unit) Move() bool {
    distance_left := utils.FindDirectionVectorLength(unit.X, unit.Y, unit.ToX, unit.ToY)
    if (distance_left <= utils.UNIT_VECTOR_MODIFIER) {
        unit.X = unit.ToX
        unit.Y = unit.ToY

        return true
    }

    x_unit, y_unit := utils.FindDirectionUnitVector(unit.X, unit.Y, unit.ToX, unit.ToY)
    unit.X += x_unit * utils.UNIT_VECTOR_MODIFIER
    unit.Y += y_unit * utils.UNIT_VECTOR_MODIFIER

    return false
}


func (unit *Unit) TMPSetToCoordinates(to_x float64, to_y float64) bool {
    unit.ToX = to_x
    unit.ToY = to_y

    if (unit.X == unit.ToX && unit.Y == unit.ToY) {
        return false
    }

    return true
}


func NewUnit(id string, color string, size int, squad_name string) *Unit{
    new_unit := Unit{
        Id: id,
        Color: color,
        Size: size,
        X: 0.0,
        Y: 0.0,
        ToX: 0.0,
        ToY: 0.0,
        SquadName: squad_name,
    }

    return &new_unit
}


func MakeUnits(number_of_units int, color string, size int, squad_name string) map[string]*Unit {
    TMP_units_ids := map[string]*Unit{}

    for id := 0; id < number_of_units; id++ {
        new_unit_id := squad_name + "_" + strconv.Itoa(id)
        TMP_units_ids[new_unit_id] = NewUnit(new_unit_id, color, size, squad_name)
    }

    return TMP_units_ids
}


type DeadUnit struct {
    Id string
    SquadName string
    Line int
    Column int
}


func NewDeadUnit(id string, squad_name string, line int, column int) *DeadUnit{
    new_dead_unit := DeadUnit{
        Id: id,
        SquadName: squad_name,
        Line: line,
        Column: column,
    }

    return &new_dead_unit
}
