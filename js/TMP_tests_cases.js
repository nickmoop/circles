function test_create_session() {
    let user_name = 'Champion_test';
    let session_id = 'test_id';
    let create_session_url = 'http://' + HOSTNAME + '/createSession';
    let create_session_parameters = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({'UserName': user_name, 'SessionId': session_id})
    };

    let func_creation_error = function(err) {
        console.log('Something wrong with session creation: ' + err);
    }
    let func_created = function(res) {
        console.log('Session with id: ' + res.SessionId + '. For user: ' + res.UserName + '. Was created.');
        join_ws_session(res.SessionId, res.UserName);
    }
    let func_create_session = function(res) {
        if (res.status != 200) {
            return res.text().then(
                function(res) {
                    console.log('Session creation error: ' + res);
                }, func_creation_error
            );
        }

        return res.json().then(func_created, func_creation_error);
    }

    fetch(create_session_url, create_session_parameters).then(
        func_create_session, func_creation_error
    );
}


function test_squad_creation() {
    create_squad('Awesome squad', 10, '#FF0AFA', 10, 5);
}


function test_units_kill() {
    if (!SOCKET_BATTLE) {
        console.log('have no socket battle');
        return;
    }

    let data = JSON.stringify({
        'Command': 'kill_units',
        'Parameters': {}
    });
    SOCKET_BATTLE.send(data);
}


function test_move_squad() {
    move_squad_to('Awesome squad', 90, 4, 200, 200);
}
