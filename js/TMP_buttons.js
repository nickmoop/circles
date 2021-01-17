class TMPButton {
    constructor() {
        this.div = document.createElement('BUTTON');
        this.div.className = 'button';
        this.text = 'Button Text';
        this.onclick_function = dummyButtonClick;
        this.x = 0;
        this.y = 0;
    }

    TMP_draw(interface_div) {
        interface_div.interface.appendChild(this.div);
    }

    hide() {
        this.div.style.display = "none"
    }

    unhide() {
        this.div.style.display = "block"
    }

    set x(value) {
        this.div.style.left = value + 'px';
    }

    set y(value) {
        this.div.style.top = value + 'px';
    }

    set text(value) {
        this.div.innerHTML = value;
        this.div.id = value;
    }

    set onclick_function(value) {
        this.div.onclick = value;
    }
}


function create_squad_button_click() {
    if (!SOCKET_BATTLE) {
        alert('Empty socket_battle');
        return;
    }

    let name = prompt('Squad name:', 'Awesome Squad');

    if (SQUADS.has(name)) {
        alert('Squad with name: ' + name + ' already exists.');
        return;
    }

    let size = parseFloat(prompt('Squad size:', '13'));
    let color = prompt('Squad color:', '#FF0AFA');
    let circles_ids = Array(parseInt(size));

    if (SQUADS.has(name)) {
        alert('Squad with name: ' + name + ' already exists.');
        return;
    }

    if (size > 900 || size < 4) {
        alert('Size should be in interval 4 - 900. Size is ' + size);
        return;
    }

    create_squad(name, size, color, 10, 5);
}


function join_ws_session(session_id, user_name) {
    if (!SOCKET_CHAT) {
        SOCKET_CHAT = new TMPSocket('chatJoin', session_id, user_name);
        SOCKET_CHAT.socket.onmessage = on_chat_message;
    }

    if (!SOCKET_BATTLE) {
        SOCKET_BATTLE = new TMPSocket('battleJoin', session_id, user_name);
        SOCKET_BATTLE.socket.onmessage = on_battle_message;
    }
}


function create_new_session() {
    let user_name = prompt('User name:', 'Champion');
    let session_id = prompt('Session id:', '1');
    let create_session_url = 'http://' + HOSTNAME + '/createSession';
    let create_session_parameters = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({'UserName': user_name, 'SessionId': session_id})
    };

    let func_creation_error = function(err) {
        alert('Something wrong with session creation: ' + err);
    }
    let func_created = function(res) {
        console.log('Session with id: ' + res.SessionId + '. For user: ' + res.UserName + '. Was created.');
        USER_NAME = res.UserName;
        join_ws_session(res.SessionId, res.UserName);
    }
    let func_create_session = function(res) {
        if (res.status != 200) {
            return res.text().then(
                function(res) {
                    alert('Session creation error: ' + res);
                }, func_creation_error
            );
        }

        return res.json().then(func_created, func_creation_error);
    }

    fetch(create_session_url, create_session_parameters).then(
        func_create_session, func_creation_error
    );
}


function close_connection() {
    if (SOCKET_BATTLE) {
        SOCKET_BATTLE.close();
        SOCKET_BATTLE = null;
    }

    if (SOCKET_CHAT) {
        SOCKET_CHAT.close();
        SOCKET_CHAT = null;
    }
}


function send_chat_message() {
    if (!SOCKET_CHAT) {
        alert('have no socket chat');
        return;
    }

    let data = prompt('Data to send:', 'message');

    SOCKET_CHAT.send(data);
}


function clock_start() {
    if (!SOCKET_BATTLE) {
        alert('have no socket battle');
        return;
    }

    let data = JSON.stringify({
        'Command': 'clock',
        'Parameters': {
            'State': 'start',
        }
    });
    SOCKET_BATTLE.send(data);
};


function clock_stop() {
    if (!SOCKET_BATTLE) {
        alert('have no socket battle');
        return;
    }

    let data = JSON.stringify({
        'Command': 'clock',
        'Parameters': {
            'State': 'stop',
        }
    });
    SOCKET_BATTLE.send(data);
};


function dummyButtonClick() {
    alert('button clicked');
};
