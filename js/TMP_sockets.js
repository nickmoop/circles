class TMPSocket {
    constructor(socket_type, session_id, user_name) {
        this.socket = new WebSocket('ws://' + HOSTNAME + '/' + socket_type);
        this.socket.session_id = session_id;
        this.socket.user_name = user_name;
        this.socket.onopen = this.on_open;
        this.socket.onmessage = on_message;
        this.socket.onclose = on_close;
        this.socket.onerror = on_error;
    }

    close() {
        this.socket.close()
    }

    send(data) {
        console.log('[send]:' + data);
        this.socket.send(data)
    }

    on_open(event) {
        console.log('[open] Connection open');

        let on_open_message = JSON.stringify({
            'UserName': this.user_name,
            'SessionId': this.session_id
        });

        console.log('[open] Send on open message: ' + on_open_message);

        this.send(on_open_message);
    }
}


function on_open(event) {
    console.log('[open] Connection open');
}


function on_message(event) {
    console.log('[message] Data recv: ' + event.data);
}


function on_chat_message(event) {
    income_chat_data = JSON.parse(event.data);
    MY_INTERFACE.chat_window.add_text(income_chat_data.UserName + ": " + income_chat_data.Message);
}


function on_battle_message(event) {
    income_battle_data = JSON.parse(event.data);
    console.log('[recv battle]: ' + income_battle_data['Response']);
    response = income_battle_data['Response'];
    battle_function = BATTLE_FUNCTIONS.get(response)
    battle_function(income_battle_data);
}


function on_close(event) {
    if (event.wasClean) {
        console.log('[close] Connection closed. ' + event.reason);
    } else {
        console.log('[close] Connection closed with error');
    }
}


function on_error(error) {
    console.log('[error] ' + error.message);
}
