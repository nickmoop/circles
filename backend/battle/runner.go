package battle


import (
	"encoding/json"
    "log"
)


func ExecuteBattleCommand(user_name string, session_id string, message string) []byte {
    battle_command := ParseBattleMessage(user_name, session_id, message)
    log.Println(battle_command)

    response := battle_command.executeCommand()

    return response
}


func ParseBattleMessage(user_name string, session_id string, message string) BattleCommand {
    battle_command := BattleCommand{UserName: user_name, SessionId: session_id}
    json.Unmarshal([]byte(message), &battle_command)

    return battle_command
}
