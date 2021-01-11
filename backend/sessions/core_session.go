package sessions


import (
    "log"
	"encoding/json"
	"net/http"
	"time"

	"backend/connections"
	"backend/messages"
	"backend/squads"
	"backend/units"
	"backend/utils"
)



type CoreSession struct {
    SessionId string
    Users map[string]UserSession
    Squads map[string]*squads.Squad
    Units map[string]*units.Unit
    KilledUnits map[string]*units.DeadUnit
    ClockState string
    TickCount int
    Ticker *time.Ticker
    TickChannel chan bool
}

type UserSession struct {
    UserName string
    SessionId string
    ChatConnection *connections.WSConnection
    BattleConnection *connections.WSConnection
}


var CORE_SESSIONS_POOL = make(map[string]*CoreSession)


func TMPSendBattleMessageToAllInSession(session_id string, sender_name string, message []byte) {
    for _, battle_connection := range getSessionBattleConnections(session_id) {
        battle_connection.WriteMessage(message)
    }
}


func TMPSendChatMessageToAllInSession(session_id string, sender_name string, message string) {
    json_message := messages.JsonChatMessage(sender_name, message)
    for _, chat_connection := range getSessionChatConnections(session_id) {
        chat_connection.WriteMessage(json_message)
    }
}


func GetCoreSessionById(session_id string) *CoreSession {
    return CORE_SESSIONS_POOL[session_id]
}


func getSessionChatConnections(session_id string) []*connections.WSConnection {
    session := GetCoreSessionById(session_id)

    chat_connections := []*connections.WSConnection{}
    for _, core_session := range session.Users {
        chat_connections = append(chat_connections, core_session.ChatConnection)
    }

    return chat_connections
}

func getSessionBattleConnections(session_id string) []*connections.WSConnection {
    session := GetCoreSessionById(session_id)

    chat_connections := []*connections.WSConnection{}
    for _, core_session := range session.Users {
        chat_connections = append(chat_connections, core_session.BattleConnection)
    }

    return chat_connections
}


func TMPAddToSessionsPool(session_json CoreSessionJson) *string {
    if exists_session, ok := CORE_SESSIONS_POOL[session_json.SessionId]; ok {

        if _, ok := exists_session.Users[session_json.UserName]; ok {
            var error_message = "session id and user name already exists."
            return &error_message

        } else {
            exists_session.Users[session_json.UserName] = NewUserSession(session_json.SessionId, session_json.UserName)
        }

    } else {
        core_session := NewCoreSession(session_json.SessionId, session_json.UserName)
        CORE_SESSIONS_POOL[session_json.SessionId] = core_session
    }

    return nil
}


func JoinChatConnectionToSession(connection *connections.WSConnection) {
    for {
		message, err := connection.ReadMessage()

		if err != nil {
		    log.Println("error on join chat connection to session:", err)
		}

        session_json := NewCoreSessionJsonFromBytes(message)

		if message != nil {
		    connection.SessionId = session_json.SessionId
		    connection.UserName = session_json.UserName
		    user_session := CORE_SESSIONS_POOL[session_json.SessionId].Users[session_json.UserName]
		    user_session.ChatConnection = connection
		    CORE_SESSIONS_POOL[session_json.SessionId].Users[session_json.UserName] = user_session
		    log.Println("chat connection added to session:", session_json.SessionId, "user name:", session_json.UserName)
            break
		}
	}
}


func JoinBattleConnectionToSession(connection *connections.WSConnection) {
    for {
		message, err := connection.ReadMessage()

		if err != nil {
		    log.Println("error on join battle connection to session:", err)
		}

        session_json := NewCoreSessionJsonFromBytes(message)

		if message != nil {
		    connection.SessionId = session_json.SessionId
		    connection.UserName = session_json.UserName
		    user_session := CORE_SESSIONS_POOL[session_json.SessionId].Users[session_json.UserName]
		    user_session.BattleConnection = connection
		    CORE_SESSIONS_POOL[session_json.SessionId].Users[session_json.UserName] = user_session
		    log.Println("battle connection added to session:", session_json.SessionId, "user name:", session_json.UserName)
            break
		}
	}
}


func RemoveChatConnectionFromSession(connection *connections.WSConnection) {
    delete(CORE_SESSIONS_POOL[connection.SessionId].Users, connection.UserName)
    log.Println("connection for session", connection.SessionId, "and user:", connection.UserName, "removed from session")
    connection = nil
}


func NewCoreSession(session_id string, user_name string) *CoreSession {
    session := new(CoreSession)
    session.SessionId = session_id
    users := map[string]UserSession{user_name: NewUserSession(session_id, user_name)}
    session.Users = users
    session.Squads = make(map[string]*squads.Squad)
    session.Units = make(map[string]*units.Unit)
    session.KilledUnits = make(map[string]*units.DeadUnit)
    session.ClockState = "stop"
    session.TickCount = 0
    session.Ticker = time.NewTicker(utils.MS_PER_TICK)
    session.TickChannel = make(chan bool)

    return session
}


func NewUserSession(session_id string, user_name string) UserSession {
    session := new(UserSession)
    session.UserName = user_name
    session.SessionId = session_id

    return UserSession{
        UserName: user_name,
        SessionId: session_id,
    }
}


func (session *CoreSessionJson) Decode(request *http.Request) {
    decoder := json.NewDecoder(request.Body)
    err := decoder.Decode(&session)
    if err != nil {
        panic(err)
    }
}


type CoreSessionJson struct {
    UserName string
    SessionId string
}


func NewCoreSessionJson(request *http.Request) CoreSessionJson {
    session_json := new(CoreSessionJson)
    session_json.Decode(request)

    return *session_json
}


func NewCoreSessionJsonFromBytes(message []byte) CoreSessionJson {
    session_json := new(CoreSessionJson)
    err := json.Unmarshal(message, session_json)
    if err != nil {
        panic(err)
    }

    return *session_json
}


func (session *CoreSession) KillUnits() {
    if (len(session.KilledUnits) == 0) {
        return
    }

    kill_message := JsonKillUnitResponse("kill_units", session.KilledUnits)

    for unit_id, dead_unit := range session.KilledUnits {
        session.Squads[dead_unit.SquadName].KillUnit(dead_unit)
        delete(session.Units, unit_id)
        delete(session.KilledUnits, unit_id)
    }

    TMPSendBattleMessageToAllInSession(session.SessionId, "Server", kill_message)
}


func (session *CoreSession) TMPStartTicking() {
    log.Println("start ticking. session id:", session.SessionId)

    go func() {
        for {
            select {
            case <-session.TickChannel:
                return
            case <-session.Ticker.C:
                session.TMPTick()
            }
        }
    }()
}


func (session *CoreSession) TMPStopTicking() {
    log.Println("stop ticking. session id:", session.SessionId)
    session.TickChannel <- true
}


func (session *CoreSession) TMPTick() {
    // move units if exists some units to move
    for _, squad := range session.Squads {
        session.ProcessSquadActivity(squad)
    }

    session.KillUnits()

    session.TickCount++
    log.Println("session id:", session.SessionId, "tick count:", session.TickCount)
}


func (session *CoreSession) ProcessSquadActivity(squad *squads.Squad) {
    if (squad.TMPCorrectFormation()) {
        correct_formation_message := JsonCorrectFormationResponse("correct_formation", squad)
        TMPSendBattleMessageToAllInSession(session.SessionId, "Server", correct_formation_message)
    } else if (squad.TMPReshapeFormation()) {
        correct_formation_message := JsonCorrectFormationResponse("correct_formation", squad)
        TMPSendBattleMessageToAllInSession(session.SessionId, "Server", correct_formation_message)
    } else if (squad.RotateSquad()) {
        change_direction_message := JsonChangeDirectionResponse("change_direction", squad)
        TMPSendBattleMessageToAllInSession(session.SessionId, "Server", change_direction_message)
    } else if (squad.MoveSquad()) {
        change_direction_message := JsonChangeDirectionResponse("change_direction", squad)
        TMPSendBattleMessageToAllInSession(session.SessionId, "Server", change_direction_message)
    }

    squad.TMPProcessAttack(session.Squads, session.KilledUnits)

    if (len(squad.UnitsToMove) == 0) {
        return
    }

    moved_units := squad.MoveUnits()
    move_message := JsonMoveUnitsResponse("move_units", moved_units)
    TMPSendBattleMessageToAllInSession(session.SessionId, "Server", move_message)
}


type CreateCorrectFormationResponse struct {
    Response string
    Name string
    Formation [][]*units.Unit
    X float64
    Y float64
    Direction float64
    Columns int
    Lines int
}


func JsonCorrectFormationResponse(status string, squad *squads.Squad) []byte {
    response := CreateCorrectFormationResponse{
        status,
        squad.Name,
        squad.Formation,
        squad.X,
        squad.Y,
        squad.Direction,
        squad.Columns,
        squad.Lines,
    }
    response_json, err := json.Marshal(response)
    if err != nil {
		log.Println("create correct formation response encode error:", err, "status:", status)
		return nil
	}
    return response_json
}


type CreateChangeDirectionResponse struct {
    Response string
    Name string
    Direction float64
    X float64
    Y float64
    Columns int
    Lines int
}


func JsonChangeDirectionResponse(status string, squad *squads.Squad) []byte {
    response := CreateChangeDirectionResponse{
        status,
        squad.Name,
        squad.Direction,
        squad.X,
        squad.Y,
        squad.Columns,
        squad.Lines,
    }
    response_json, err := json.Marshal(response)
    if err != nil {
		log.Println("create change squad direction response encode error:", err, "status:", status)
		return nil
	}
    return response_json
}


type CreateKillUnitsResponse struct {
    Response string
    DeadUnits map[string]*units.DeadUnit
}


func JsonKillUnitResponse(status string, dead_units map[string]*units.DeadUnit) []byte {
    response := CreateKillUnitsResponse{
        status,
        dead_units,
    }
    response_json, err := json.Marshal(response)
    if err != nil {
		log.Println("create kill units response encode error:", err, "status:", status)
		return nil
	}
    return response_json
}


type UnitJson struct {
    Id string
    X float64
    Y float64
}

type CreateMoveUnitResponse struct {
    Response string
    UnitsList []*UnitJson
}


func JsonMoveUnitsResponse(status string, units_to_move []*units.Unit) []byte {
    units_json := make([]*UnitJson, len(units_to_move))
    for index, unit := range units_to_move {
        units_json[index] = &UnitJson{unit.Id, unit.X, unit.Y}
    }

    response := CreateMoveUnitResponse{
        status,
        units_json,
    }

    response_json, err := json.Marshal(response)
    if err != nil {
		log.Println("create move units response encode error:", err, "status:", status)
		return nil
	}
    return response_json
}
