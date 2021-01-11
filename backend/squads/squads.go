package squads



import (
    "math"
    "log"

    "backend/units"
    "backend/utils"
)


type Squad struct {
    Name string
    Color string
    Size int
    Units map[string]*units.Unit
    Formation [][]*units.Unit
    TMPFormation [][][]float64
    UnitsToMove map[string]bool
    Columns int
    Lines int
    Interval int
    Width int
    Height int
    Direction float64
    DirectionUnit float64
    UnitSize int
    X float64
    Y float64
    ToX float64
    ToY float64
    RotateTo float64
    ReshapeTo int
    TMPOwnerName string
}


func NewSquad(name string, color string, size int, unit_size int, interval int, user_name string) *Squad {
    new_units := units.MakeUnits(size, color, unit_size, name)

    squad := Squad{
        Name: name,
        Color: color,
        Size: size,
        Units: new_units,
        Interval: interval,
        UnitSize: unit_size,
        X: 200,
        Y: 200,
        UnitsToMove: make(map[string]bool),
        Direction: 0,
        TMPOwnerName: user_name,
    }

    squad.ToX = squad.X
    squad.ToY = squad.Y
    squad.RotateTo = squad.Direction
    squad.SetColumns(int(math.Ceil(math.Sqrt(float64(len(new_units))))))
    squad.ReshapeTo = squad.Columns
    squad.CalculateShape()
    squad.InitTMPFormation()
    squad.RotateTMPFormation(squad.Direction)
    squad.InitFormation()

    return &squad
}


func (squad *Squad) TMPProcessAttack(session_squads map[string]*Squad, killed_units map[string]*units.DeadUnit) {
    units_made_attack := make(map[string]bool)

    for _, other_squad := range session_squads {
        if (squad.TMPOwnerName == other_squad.TMPOwnerName) {
            continue
        }

        distance_between_centers := utils.FindDirectionVectorLength(squad.X, squad.Y, other_squad.X, other_squad.Y)
        attack_distance := utils.FindDirectionVectorLength(0, 0, float64(squad.Width + other_squad.Width) / 2, float64(squad.Height + other_squad.Height) / 2) + float64(utils.ATTACK_RANGE)

        if (distance_between_centers > attack_distance) {
            continue
        }

        squad.AttackOtherSquad(other_squad, units_made_attack, killed_units)
    }
}


func (squad *Squad) AttackOtherSquad(other_squad *Squad, units_made_attack map[string]bool, killed_units map[string]*units.DeadUnit) {
    possible_dead_units := squad.FindOtherUnitsToAttack(other_squad)
    possible_attackers_units := other_squad.FindOtherUnitsToAttack(squad)

    for attacker_id, _ := range possible_attackers_units {
        if _, ok := units_made_attack[attacker_id]; ok {
            continue
        }

        for defender_id, defender_unit := range possible_dead_units {
            if _, ok := killed_units[defender_id]; ok {
                continue
            }

            distance_between_units := utils.FindDirectionVectorLength(squad.Units[attacker_id].X, squad.Units[attacker_id].Y, other_squad.Units[defender_id].X, other_squad.Units[defender_id].Y)
            if (distance_between_units <= float64(utils.ATTACK_RANGE)) {
                units_made_attack[attacker_id] = true
                killed_units[defender_id] = defender_unit
                break
            }
        }
    }
}


func (squad *Squad) FindOtherUnitsToAttack(other_squad *Squad) map[string]*units.DeadUnit {
    start_column := 0
    start_line := 0
    end_column := other_squad.Columns
    end_line := other_squad.Lines
    increment_column := 1
    increment_line := 1

    top_left := utils.FindDirectionVectorLength(squad.X, squad.Y, other_squad.TMPFormation[0][0][0], other_squad.TMPFormation[0][0][1])
    top_right := utils.FindDirectionVectorLength(squad.X, squad.Y, other_squad.TMPFormation[0][other_squad.Columns - 1][0], other_squad.TMPFormation[0][other_squad.Columns - 1][1])
    bottom_left := utils.FindDirectionVectorLength(squad.X, squad.Y, other_squad.TMPFormation[other_squad.Lines - 1][0][0], other_squad.TMPFormation[other_squad.Lines - 1][0][1])
    bottom_right := utils.FindDirectionVectorLength(squad.X, squad.Y, other_squad.TMPFormation[other_squad.Lines - 1][other_squad.Columns - 1][0], other_squad.TMPFormation[other_squad.Lines - 1][other_squad.Columns - 1][1])
    min_distance := top_left

    if (top_right < min_distance) {
        min_distance = top_right
        start_column = other_squad.Columns - 1
        start_line = 0
        end_column = 0
        end_line = other_squad.Lines
        increment_column = -1
        increment_line = 1
    }
    if (bottom_left < min_distance) {
        min_distance = bottom_left
        start_column = 0
        start_line = other_squad.Lines - 1
        end_column = other_squad.Columns
        end_line = 0
        increment_column = 1
        increment_line = -1
    }
    if (bottom_right < min_distance) {
        start_column = other_squad.Columns - 1
        start_line = other_squad.Lines - 1
        end_column = 0
        end_line = 0
        increment_column = -1
        increment_line = -1
    }

    attack_distance := utils.FindDirectionVectorLength(0, 0, float64(squad.Width + other_squad.Width) / 2, float64(squad.Height + other_squad.Height) / 2) + float64(utils.ATTACK_RANGE)
    possible_dead_units := make(map[string]*units.DeadUnit)

    for line := start_line; line != end_line; line += increment_line {
        find_in_line := false

        for column := start_column; column != end_column; column += increment_column {
            distance := utils.FindDirectionVectorLength(squad.X, squad.Y, other_squad.TMPFormation[line][column][0], other_squad.TMPFormation[line][column][1])
            if (distance <= attack_distance) {
                find_in_line = true

                unit := other_squad.Formation[line][column]
                if (unit == nil) {
                    continue
                }
                possible_dead_units[unit.Id] = units.NewDeadUnit(unit.Id, other_squad.Name, line, column)
            }
        }

        if (!find_in_line) {
            break
        }
    }

    for column := start_column; column != end_column; column += increment_column {
        find_in_column := false

        for line := start_line; line != end_line; line += increment_line {
            distance := utils.FindDirectionVectorLength(squad.X, squad.Y, other_squad.TMPFormation[line][column][0], other_squad.TMPFormation[line][column][1])
            if (distance <= attack_distance) {
                find_in_column = true

                unit := other_squad.Formation[line][column]
                if (unit == nil) {
                    continue
                }
                possible_dead_units[unit.Id] = units.NewDeadUnit(unit.Id, other_squad.Name, line, column)
            }
        }

        if (!find_in_column) {
            break
        }
    }

    return possible_dead_units
}


func (squad *Squad) TMPReshapeFormation() bool {
    if (squad.ReshapeTo == squad.Columns) {
        return false
    }

    old_columns := squad.Columns
    old_lines := squad.Lines

    if (squad.ReshapeTo - old_columns > 2) {
        squad.SetColumns(old_columns + 2)
    } else if (squad.ReshapeTo - old_columns < -2) {
        squad.SetColumns(old_columns - 2)
    } else {
        squad.SetColumns(squad.ReshapeTo)
    }

    if (old_columns == squad.Columns) {
        squad.ReshapeTo = squad.Columns
        return false
    }

    squad.Lines = utils.MaxInt(squad.Lines, old_lines)
    squad.CalculateShape()
    squad.InitTMPFormation()
    squad.RotateTMPFormation(squad.Direction)

    new_formation := make([][]*units.Unit, squad.Lines)

    if (squad.Columns > old_columns) {
        columns_delta := squad.Columns - old_columns
        left_add := make([]*units.Unit, columns_delta / 2)
        right_add := make([]*units.Unit, columns_delta / 2 + columns_delta % 2)

        for line := 0; line < old_lines; line++ {
            new_formation[line] = append(append(left_add, squad.Formation[line]...), right_add...)
        }
    } else if (squad.Columns < old_columns) {
        new_column_index := 0
        new_line_index := squad.Lines - 1

        for line := 0; line < squad.Lines; line++ {
            new_formation[line] = make([]*units.Unit, squad.Columns)
        }

        for line := old_lines - 1; line >= 0; line-- {
            for column := 0; column < old_columns; column++ {
                log.Println(line, column)
                unit := squad.Formation[line][column]
                if (unit == nil) {
                    continue
                }

                new_formation[new_line_index][new_column_index] = unit
                new_column_index += 1
                if (new_column_index == squad.Columns) {
                    new_column_index = 0
                    new_line_index -= 1
                }
            }
        }
    }

    squad.Formation = new_formation

    return true
}


func (squad *Squad) MoveSquad() bool {
    tmp_distance_between := utils.FindDirectionVectorLength(squad.X, squad.Y, squad.ToX, squad.ToY)
    if (tmp_distance_between == 0) {
        return false
    }

    move_degree := utils.FindVectorDegree(squad.X, squad.Y, squad.ToX, squad.ToY)
    if (math.Abs(squad.Direction - move_degree) > 1) {
        squad.RotateTo = move_degree
        return false
    }

    if (tmp_distance_between <= utils.UNIT_VECTOR_MODIFIER) {
        squad.X = squad.ToX
        squad.Y = squad.ToY

        squad.InitTMPFormation()
        squad.RotateTMPFormation(squad.Direction)
    } else {
        x_unit, y_unit := utils.FindDirectionUnitVector(squad.X, squad.Y, squad.ToX, squad.ToY)
        x_add := x_unit * utils.UNIT_VECTOR_MODIFIER
        y_add := y_unit * utils.UNIT_VECTOR_MODIFIER
        squad.X += x_add
        squad.Y += y_add

        for line := 0; line < squad.Lines; line++ {
            for column := 0; column < squad.Columns; column++ {
                squad.TMPFormation[line][column][0] += x_add
                squad.TMPFormation[line][column][1] += y_add
            }
        }
    }

    squad.TMPUnitsMoveToTMPFormation()
    return true
}


func (squad *Squad) RotateSquad() bool {
    direction_destination := utils.NormalizeDirectionAngle(squad.RotateTo)
	current_direction := utils.NormalizeDirectionAngle(squad.Direction)

    delta := direction_destination - current_direction
    step := 0.0

    if (delta == 0) {
        return false
    } else if (math.Abs(delta) <= squad.DirectionUnit) {
        squad.Direction = squad.RotateTo
        squad.InitTMPFormation()
        squad.RotateTMPFormation(squad.Direction)
        squad.TMPUnitsMoveToTMPFormation()
        return true
    } else if (delta < -180 || delta > 0 && delta <= 180) {
        step = squad.DirectionUnit
    } else if (delta > 180 || delta < 0 && delta >= -180) {
        step = -squad.DirectionUnit
    }

    squad.RotateTMPFormation(step)
    squad.Direction += step

    squad.TMPUnitsMoveToTMPFormation()

    return true
}


func (squad *Squad) KillUnit(dead_unit *units.DeadUnit) {
    delete(squad.Units, dead_unit.Id)
    squad.Formation[dead_unit.Line][dead_unit.Column] = nil
}


func (squad *Squad) TMPCorrectFormation() bool {
    max_lines := squad.Lines - 1
    moved_forward := false
    moved_from_center := false
    is_removed_lines := false

    for line := 0; line < max_lines; line++ {
        if (squad.MoveColumnsForward(line)) {
            moved_forward = true
        }

        removed_lines := squad.RemoveLastLines()
        if (removed_lines != 0) {
            max_lines -= removed_lines
            is_removed_lines = true
        }

        if (squad.MoveFromCenter(line)) {
            moved_from_center = true
        }
    }

    moved_to_center := squad.MoveToCenter(squad.Lines - 1)

    if (is_removed_lines || moved_to_center || moved_from_center || moved_forward) {
        return true
    }

    return false // !squad.IsUnitsInTMPFormation()
}


func (squad *Squad) IsUnitsInTMPFormation() bool {
    for line := 0; line < squad.Lines; line++ {
        for column := 0; column < squad.Columns; column++ {
            unit := squad.Formation[line][column]
            if (unit == nil) {
                continue
            }
            if (unit.X != squad.TMPFormation[line][column][0] || unit.Y != squad.TMPFormation[line][column][1]) {
                return false
            }
        }
    }

    return true
}


func (squad *Squad) MoveColumnsForward(line int) bool {
    moved_forward := false

    for column := 0; column < squad.Columns; column++ {
        unit := squad.Formation[line][column]
        if (unit != nil) {
            continue
        }

        for tmp_line := line + 1; tmp_line < squad.Lines; tmp_line++ {
            not_nil_unit := squad.Formation[tmp_line][column]
            if (not_nil_unit == nil) {
                continue
            }

            if (squad.TMPChangeUnitPositionInFormation(not_nil_unit, tmp_line, column, line, column)) {
                moved_forward = true
                break
            }
        }
    }

    return moved_forward
}


func (squad *Squad) MoveFromCenter(line int) bool {
	max_iteration_index := squad.Columns / 2 + squad.Columns % 2
	left_zero_indexes := []int{}
	right_zero_indexes := []int{}
	moved_from_center := false

	for index := 0; index < max_iteration_index; index++ {
		squad.TMPLineOperation(line, index, &left_zero_indexes)

		right_index := squad.Columns-index-1
		if (index == right_index) {
            continue
		}

        if (squad.TMPLineOperation(line, right_index, &right_zero_indexes)) {
            moved_from_center = true
        }
	}

	return moved_from_center
}


func (squad *Squad) MoveToCenter(line int) bool {
    max_iteration_index := squad.Columns / 2 - 1 + squad.Columns % 2
	left_zero_indexes := []int{}
	right_zero_indexes := []int{}
	moved_to_center := false

	for index := max_iteration_index; index >= 0; index-- {
		squad.TMPLineOperation(line, index, &left_zero_indexes)

		right_index := squad.Columns-index-1
		if (index == right_index) {
            continue
		}

		if (squad.TMPLineOperation(line, right_index, &right_zero_indexes)) {
		    moved_to_center = true
		}
	}

    if (squad.TMPMoveLastUnitsDiagonaly(line, max_iteration_index)) {
        moved_to_center = true
    }

    if (squad.TMPLineCentrize(line, len(left_zero_indexes) - len(right_zero_indexes), max_iteration_index)) {
        moved_to_center = true
    }

	return moved_to_center
}


func (squad *Squad) TMPMoveLastUnitsDiagonaly(line int, max_iteration_index int) bool {
    if (line == 0 || max_iteration_index == 0) {
        return false
    }

    for column := max_iteration_index; column <= max_iteration_index + 1; column++{
	    unit := squad.Formation[line][column]
	    if (unit == nil) {
	        continue
	    }

	    left_forward_unit := squad.Formation[line-1][column-1]
	    if (left_forward_unit == nil) {
	        if (squad.TMPChangeUnitPositionInFormation(unit, line, column, line-1, column-1)) {
	            return true
	        }
	    }

	    right_forward_unit := squad.Formation[line-1][column+1]
	    if (right_forward_unit == nil) {
	        if (squad.TMPChangeUnitPositionInFormation(unit, line, column, line-1, column+1)) {
	            return true
	        }
	    }
	}

	return false
}


func (squad *Squad) TMPLineCentrize(line int, delta_len int, max_iteration_index int) bool {
    is_unit_moved := false

	if (delta_len >= 2) {
	    for column := 2; column <= max_iteration_index + 1; column++ {
	        unit := squad.Formation[line][column]
	        if (unit == nil) {
	            continue
	        }

	        if (squad.TMPChangeUnitPositionInFormation(unit, line, column, line, column-1)) {
	            is_unit_moved = true
	        }
	    }
	} else if (delta_len <= -2) {
	    for column := squad.Columns-2; column >= max_iteration_index - 1; column-- {
	        unit := squad.Formation[line][column]
	        if (unit == nil) {
	            continue
	        }

	        if (squad.TMPChangeUnitPositionInFormation(unit, line, column, line, column+1)) {
	            is_unit_moved = true
	        }
	    }
	}

	return is_unit_moved
}


func (squad *Squad) TMPLineOperation(line int, column int, zero_indexes *[]int) bool {
    unit := squad.Formation[line][column]
    new_column := 0
    is_unit_moved := false

    if (unit == nil) {
        (*zero_indexes) = append(*zero_indexes, column)
    } else if (len(*zero_indexes) != 0) {
        (*zero_indexes) = append(*zero_indexes, column)
        new_column, (*zero_indexes) = (*zero_indexes)[0], (*zero_indexes)[1:]

        if (squad.TMPChangeUnitPositionInFormation(unit, line, column, line, new_column)) {
            is_unit_moved = true
        }
    }

    return is_unit_moved
}


func (squad *Squad) RemoveLastLines() int {
    removed_lines_counter := 0

    for column := 0; column < squad.Columns; column++ {
        unit := squad.Formation[squad.Lines-1][column]
        if (unit != nil) {
            return removed_lines_counter
        }
    }

    squad.Lines -= 1
    new_x := (squad.TMPFormation[0][0][0] + squad.TMPFormation[squad.Lines-1][squad.Columns-1][0]) / 2
    new_y := (squad.TMPFormation[0][0][1] + squad.TMPFormation[squad.Lines-1][squad.Columns-1][1]) / 2

    if (squad.X == squad.ToX && squad.Y == squad.ToY) {
        squad.ToX = new_x
        squad.ToY = new_y
    }

    squad.X = new_x
    squad.Y = new_y
    squad.CalculateShape()

    squad.Formation = squad.Formation[:squad.Lines]
    squad.TMPFormation = squad.TMPFormation[:squad.Lines]
    removed_lines_counter += 1
    removed_lines_counter += squad.RemoveLastLines()

    return removed_lines_counter
}


func (squad *Squad) MoveUnits() []*units.Unit {
    moved_units := make([]*units.Unit, 0)
    for unit_id, _ := range squad.UnitsToMove {
        moved_units = append(moved_units, squad.Units[unit_id])

        if (squad.Units[unit_id].Move()) {
            delete(squad.UnitsToMove, unit_id)
        }
    }

    return moved_units
}


func (squad *Squad) InitFormation() {
    line := 0
    column := 0

    squad.Formation = make([][]*units.Unit, squad.Lines)
    for i := 0; i < squad.Lines; i++ {
        squad.Formation[i] = make([]*units.Unit, squad.Columns)
    }

    for _, unit := range squad.Units {
        // initialize formation
        unit_position := squad.TMPFormation[line][column]
        unit.X = unit_position[0]
        unit.Y = unit_position[1]
        unit.TMPSetToCoordinates(unit_position[0], unit_position[1])

        squad.Formation[line][column] = unit

        column += 1
        if (column >= squad.Columns) {
            column = 0
            line += 1
        }
    }
}


func (squad *Squad) InitTMPFormation() {
    squad.TMPFormation = make([][][]float64, squad.Lines)

    for line := 0; line < squad.Lines; line++ {
        squad.TMPFormation[line] = make([][]float64, squad.Columns)

        for column :=0; column < squad.Columns; column ++ {
            // initialize formation
            x := squad.FirstX() + float64(column * squad.UnitSize) + float64(column * squad.Interval) + float64(squad.UnitSize / 2)
            y := squad.FirstY() + float64(line * squad.UnitSize) + float64(line * squad.Interval) + float64(squad.UnitSize / 2)

            squad.TMPFormation[line][column] = []float64{x, y}
        }
    }
}


func (squad *Squad) RotateTMPFormation(direction float64) {
    sin_ := math.Sin(utils.ToRadians(direction))
    cos_ := math.Cos(utils.ToRadians(direction))

    for line := 0; line < squad.Lines; line++ {
        for column := 0; column < squad.Columns; column++ {
            x_to, y_to := utils.ChangeCoordinatesSystemAndRotate(squad.TMPFormation[line][column][0], squad.TMPFormation[line][column][1], squad.X, squad.Y, cos_, sin_)

            squad.TMPFormation[line][column][0] = x_to
            squad.TMPFormation[line][column][1] = y_to
        }
    }
}


func (squad *Squad) CalculateDirectionUnit() {
    radius := utils.FindDirectionVectorLength(0, 0, float64(squad.Width / 2), float64(squad.Height / 2))
    squad.DirectionUnit = 180 * utils.UNIT_VECTOR_MODIFIER / math.Pi / radius
}


func (squad *Squad) TMPUnitsMoveToTMPFormation() {
    for line := 0; line < squad.Lines; line++ {
        for column := 0; column < squad.Columns; column++ {
            squad.TMPUnitFormationToTMPFormation(line, column)
        }
    }
}


func (squad *Squad) TMPChangeUnitPositionInFormation(unit *units.Unit, old_line int, old_column int, new_line int, new_column int) bool {
    squad.Formation[new_line][new_column] = unit
    squad.Formation[old_line][old_column] = nil

    if (squad.TMPUnitFormationToTMPFormation(new_line, new_column)) {
        return true
    }

    return false
}


func (squad *Squad) TMPUnitFormationToTMPFormation(line int, column int) bool {
    unit := squad.Formation[line][column]
    if (unit == nil) {
        return false
    }

    if (unit.TMPSetToCoordinates(squad.TMPFormation[line][column][0], squad.TMPFormation[line][column][1])) {
        squad.UnitsToMove[unit.Id] = true
        return true
    }

    return false
}


func (squad *Squad) CalculateShape() {
    squad.Width = squad.Columns * squad.UnitSize + (squad.Columns - 1) * squad.Interval
    squad.Height = squad.Lines * squad.UnitSize + (squad.Lines - 1) * squad.Interval
    squad.CalculateDirectionUnit()
}


func (squad *Squad) SetColumns(columns int) {
    if (columns < 2) {
        columns = 2
    } else if (columns > 30) {
        columns = 30
    }

    lines := int(math.Ceil(float64(len(squad.Units) / columns) + 0.5))

    if (columns * (lines - 1) >= len(squad.Units)) {
        lines -= 1
    }

    if (lines < 2) {
        squad.SetColumns(columns - 1)
    } else if (lines > 30) {
        squad.SetColumns(columns + 1)
    } else {
        squad.Columns = columns
        squad.Lines = lines
    }
}


func (squad *Squad) FirstX() float64 {
    return squad.X - float64(squad.Width / 2)
}


func (squad *Squad) FirstY() float64 {
    return squad.Y - float64(squad.Height / 2)
}


func (squad *Squad) LastX() float64 {
    return squad.X + float64(squad.Width / 2)
}


func (squad *Squad) LastY() float64 {
    return squad.Y + float64(squad.Height / 2)
}


func (squad *Squad) GetColumns() int {
    return squad.Columns
}


func (squad *Squad) GetLines() int {
    return squad.Lines
}
