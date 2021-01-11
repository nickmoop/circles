package battle

import (
    "encoding/json"
    "log"

    "backend/squads"
)


type CreateSquadResponse struct {
    Response string
    ResponseMessage string
    NewSquad *squads.Squad
}

type CreateClockStateResponse struct {
    Response string
    ResponseMessage string
    ClockState string
}

type CreateTestKillUnitsResponse struct {
    Response string
    ResponseMessage string
}


func JsonMoveSquadResponse(status string, response_message string) []byte {
    response := CreateTestKillUnitsResponse{
        status,
        response_message,
    }
    response_json, err := json.Marshal(response)
    if err != nil {
		log.Println("create move squad response encode error:", err, "status:", status)
		return nil
	}
    return response_json
}


func JsonCreateSquadResponse(status string, response_message string, new_squad *squads.Squad) []byte {
    response := CreateSquadResponse{
        status,
        response_message,
        new_squad,
    }
    response_json, err := json.Marshal(response)
    if err != nil {
		log.Println("create squad response encode error:", err, "status:", status)
		return nil
	}
    return response_json
}


func JsonCreateClockStateResponse(status string, response_message string, clock_state string) []byte {
    response := CreateClockStateResponse{
        status,
        response_message,
        clock_state,
    }
    response_json, err := json.Marshal(response)
    if err != nil {
		log.Println("create clock state response encode error:", err, "status:", status)
		return nil
	}
    return response_json
}


func JsonTestKillUnitsResponse(status string, response_message string) []byte {
    response := CreateTestKillUnitsResponse{
        status,
        response_message,
    }
    response_json, err := json.Marshal(response)
    if err != nil {
		log.Println("create test kill units response encode error:", err, "status:", status)
		return nil
	}
    return response_json
}
