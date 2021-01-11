package battle


import (
    "log"
    "math/rand"
    "time"

    "backend/sessions"
    "backend/squads"
    "backend/units"
)


var COMMANDS_MAP = map[string]func(map[string]interface{}, string, string) []byte {
    "create_squad": commandCreateSquad,
    "clock": commandClock,
    "kill_units": commandKillUnits,
    "move_squad": commandMoveSquad,
}


type BattleCommand struct {
    UserName string
    SessionId string
    Command string
    Parameters map[string]interface{}
}


func (battle_command *BattleCommand) executeCommand() []byte {
    command := COMMANDS_MAP[battle_command.Command]
    response := command(battle_command.Parameters, battle_command.UserName, battle_command.SessionId)

    return response
}


func commandMoveSquad(parameters map[string]interface{}, user_name string, session_id string) []byte {
    session := sessions.GetCoreSessionById(session_id)
    name := parameters["Name"].(string)
    direction := parameters["Direction"].(float64)
    columns := int(parameters["Columns"].(float64))
    to_x := parameters["ToX"].(float64)
    to_y := parameters["ToY"].(float64)

    squad := session.Squads[name]
    squad.ReshapeTo = columns
    squad.RotateTo = direction
    squad.ToX = to_x
    squad.ToY = to_y

    response := "move_squad"
    response_message := "Squad '" + name + "' start moving"

    return JsonMoveSquadResponse(response, response_message)
}


func commandCreateSquad(parameters map[string]interface{}, user_name string, session_id string) []byte {
    session := sessions.GetCoreSessionById(session_id)
    new_squad_name := parameters["Name"].(string)
    response := "create_squad"

    if _, exist := session.Squads[new_squad_name]; exist{
        response_message := "Squad: '" + new_squad_name + "' already exists"
        return JsonCreateSquadResponse(response, response_message, &squads.Squad{})
    }

    new_squad_color := parameters["Color"].(string)
    new_squad_size := int(parameters["Size"].(float64))
    new_squad_unit_size := int(parameters["UnitSize"].(float64))
    new_squad_units_interval := int(parameters["UnitsInterval"].(float64))

    new_squad := squads.NewSquad(new_squad_name, new_squad_color, new_squad_size, new_squad_unit_size, new_squad_units_interval, user_name)
    session.Squads[new_squad_name] = new_squad

    for unit_id, unit  := range new_squad.Units{
        session.Units[unit_id] = unit
    }

    response_message := "Squad: '" + new_squad_name + "' created"
    return JsonCreateSquadResponse(response, response_message, new_squad)
}


func commandClock(parameters map[string]interface{}, user_name string, session_id string) []byte {
    session := sessions.GetCoreSessionById(session_id)
    clock_state := parameters["State"].(string)
    response := "clock"

    log.Println("new clock state:", clock_state)
    session.ClockState = clock_state
    log.Println("session clock state:", session.ClockState)

    if session.ClockState == "start" {
        session.TMPStartTicking()
    } else if session.ClockState == "stop" {
        session.TMPStopTicking()
    } else {
        log.Println("Unknown clock state:", session.ClockState)
        session.TMPStopTicking()
    }

    response_message := "Session clock state now: '" + session.ClockState + "'"
    return JsonCreateClockStateResponse(response, response_message, clock_state)
}


func commandKillUnits(parameters map[string]interface{}, user_name string, session_id string) []byte {
    session := sessions.GetCoreSessionById(session_id)

    for _, squad := range session.Squads {
        go func() {
            previous_tick := -1
            for {
                time.Sleep(1 * time.Second)
                ticks_delay := session.TickCount - previous_tick
                if (ticks_delay >= 10) {
                    line_to_kill := rand.Intn(squad.Lines)
                    column_to_kill := rand.Intn(squad.Columns)

                    unit_to_kill := squad.Formation[line_to_kill][column_to_kill]
                    previous_tick = session.TickCount

                    if (unit_to_kill == nil) {
                        continue
                    }

                    session.KilledUnits[unit_to_kill.Id] = units.NewDeadUnit(unit_to_kill.Id, unit_to_kill.SquadName, line_to_kill, column_to_kill)

                    if (len(squad.Units) == 0) {
                        return
                    }
                }
            }
        }()
    }

    response := "test_kill_units"
    response_message := "Units kills starts"

    return JsonTestKillUnitsResponse(response, response_message)
}
