class Interface {
    constructor() {
        this.interface = document.getElementById("interface");

        this.chat_window = create_chat_window(this);
        this.buttons = create_buttons();
        this.TMP_draw_buttons();
    }

    TMP_draw_buttons() {
        for (let button of this.buttons.values()) {
            button.TMP_draw(this);
        };
    }
}


class ChatWindow {
    constructor() {
        this.div = document.createElement('div');
        this.div.className = 'chat_window';
        this.x = 0;
        this.y = 0;
    }

    TMP_draw(interface_div) {
        interface_div.interface.appendChild(this.div);
    }

    add_text(value) {
        this.div.innerHTML += value + '\n';
    }

    set x(value) {
        this.div.style.right = value + 'px';
    }

    set y(value) {
        this.div.style.bottom = value + 'px';
    }
}


function create_chat_window(interface_div) {
    let chat_window = new ChatWindow;
    chat_window.x = 5;
    chat_window.y = 15;
    chat_window.TMP_draw(interface_div);

    return chat_window;
}


function create_buttons() {
    let buttons = new Map()

    let create_squad_button = new TMPButton;
    create_squad_button.x = 1250;
    create_squad_button.text = 'Create Squad';
    create_squad_button.onclick_function = create_squad_button_click;
    buttons.set("Create Squad", create_squad_button);

    let start_clock_button = new TMPButton();
    start_clock_button.x = 1250;
    start_clock_button.y = 40;
    start_clock_button.text = 'Start clock';
    start_clock_button.onclick_function = clock_start;
    buttons.set("Start clock", start_clock_button);

    let stop_clock_button = new TMPButton();
    stop_clock_button.x = 1250;
    stop_clock_button.y = 40;
    stop_clock_button.text = 'Stop clock';
    stop_clock_button.onclick_function = clock_stop;
    stop_clock_button.hide();
    buttons.set("Stop clock", stop_clock_button);


    let close_connection_button = new TMPButton;
    close_connection_button.x = 1250;
    close_connection_button.y = 80;
    close_connection_button.text = 'Close connection';
    close_connection_button.onclick_function = close_connection;
    buttons.set("Close connection", close_connection_button);

    let create_session_button = new TMPButton;
    create_session_button.x = 1250;
    create_session_button.y = 120;
    create_session_button.text = 'Create session';
    create_session_button.onclick_function = create_new_session;
    buttons.set("Create session", create_session_button);

    let send_chat_message_button = new TMPButton;
    send_chat_message_button.x = 1250;
    send_chat_message_button.y = 160;
    send_chat_message_button.text = 'Send message';
    send_chat_message_button.onclick_function = send_chat_message;
    buttons.set("Send message", send_chat_message_button);

    let test_create_session_button = new TMPButton;
    test_create_session_button.x = 20;
    test_create_session_button.y = 20;
    test_create_session_button.text = 'Test create_session';
    test_create_session_button.onclick_function = test_create_session;
    buttons.set("Test create session", test_create_session_button);

    let test_squad_creation_button = new TMPButton;
    test_squad_creation_button.x = 140;
    test_squad_creation_button.y = 20;
    test_squad_creation_button.text = 'Test squad creation';
    test_squad_creation_button.onclick_function = test_squad_creation;
    buttons.set("Test squad creation", test_squad_creation_button);

    let test_units_kill_button = new TMPButton;
    test_units_kill_button.x = 260;
    test_units_kill_button.y = 20;
    test_units_kill_button.text = 'Test units kill';
    test_units_kill_button.onclick_function = test_units_kill;
    buttons.set("Test units kill", test_units_kill_button);

    let test_move_squad_button = new TMPButton;
    test_move_squad_button.x = 380;
    test_move_squad_button.y = 20;
    test_move_squad_button.text = 'Test move squad';
    test_move_squad_button.onclick_function = test_move_squad;
    buttons.set("Test move squad", test_move_squad_button);

    return buttons;
};
