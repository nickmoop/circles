function create_battle_functions() {
    const battle_functions_map = new Map();
    battle_functions_map.set('create_squad', tmp_create_squad_response);
    battle_functions_map.set('clock', tmp_clock_response);
    battle_functions_map.set('kill_units', tmp_kill_units_response);
    battle_functions_map.set('move_units', tmp_move_unit_response);
    battle_functions_map.set('move_squad', tmp_move_squad_response);
    battle_functions_map.set('change_direction', tmp_change_squad_direction_response);
    battle_functions_map.set('correct_formation', tmp_correct_formation_response);
    battle_functions_map.set('test_kill_units', tmp_test_kill_units_response);

    return battle_functions_map;
}


function tmp_create_squad_response(response_json) {
    if (response_json['NewSquad']['Units'] == null) {
        response_message = response_json['ResponseMessage'];
        alert(response_message);
        return
    };

    circles_ids = tmp_make_units_battle_response(response_json['NewSquad']['Units']);
    shade_circles_ids = tmp_make_units_shades_battle_response(response_json['NewSquad']['Units']);
    let squad = new Squad(response_json['NewSquad'], circles_ids);
    let shade_squad = new ShadeSquad(response_json['NewSquad'], shade_circles_ids);
    SQUADS.set(squad.name, squad);
    SHADE_SQUADS.set(shade_squad.name, shade_squad);
    squad.TMP_draw();
    shade_squad.TMP_draw();
    shade_squad.hide();
}


function tmp_clock_response(response_json) {
    if (response_json['ClockState'] == 'start') {
        MY_INTERFACE.buttons.get('Stop clock').unhide();
        MY_INTERFACE.buttons.get('Start clock').hide();
    } else if (response_json['ClockState'] == 'stop') {
        MY_INTERFACE.buttons.get('Stop clock').hide();
        MY_INTERFACE.buttons.get('Start clock').unhide();
    } else {
        alert('Unknown clock state in server response: ' + response_json['ClockState']);
    };
}


function tmp_kill_units_response(response_json) {
    for (var unit_id of Object.keys(response_json['DeadUnits'])) {
        killCircle(CIRCLES.get(unit_id));
    };
}


function tmp_move_unit_response(response_json) {
    for (let unit of response_json['UnitsList']) {
        var circle = CIRCLES.get(unit['Id']);
        circle.x = parseFloat(unit['X'])
        circle.y = parseFloat(unit['Y'])
    };
}


function tmp_move_squad_response(response_json) {
    console.log(response_json['ResponseMessage']);
}


function tmp_change_squad_direction_response(response_json) {
    squad = SQUADS.get(response_json['Name']);
    shade_squad = SHADE_SQUADS.get(response_json['Name']);
    squad._direction = response_json['Direction'];
    squad.x = response_json['X'];
    squad.y = response_json['Y'];
    squad._columns = response_json['Columns'];
    squad._lines = response_json['Lines'];
    squad.TMP_draw_shape();

    if (USER_NAME != squad.TMPOwnerName) {
        return;
    };

    if (shade_squad.is_hidden) {
        return;
    };

    if (squad.x == shade_squad.x && squad.y == shade_squad.y) {
        shade_squad.hide();
        move_squad_to(shade_squad.name, shade_squad.direction, shade_squad.columns, shade_squad.x, shade_squad.y);
    };
}


function tmp_test_kill_units_response(response_json) {
    console.log('Kill units starts');
}


function tmp_correct_formation_response(response_json) {
    squad = SQUADS.get(response_json['Name']);
    shade_squad = SHADE_SQUADS.get(response_json['Name']);
    squad._columns = response_json['Columns'];
    squad._lines = response_json['Lines'];
    squad._direction = response_json['Direction'];
    squad.x = response_json['X'];
    squad.y = response_json['Y'];
    squad._TMP_init_formation(response_json);
    squad.TMP_draw_shape();

    if (USER_NAME != squad.TMPOwnerName) {
        return;
    };

    if (shade_squad.is_hidden) {
        return;
    };

    if (squad.x == shade_squad.x && squad.y == shade_squad.y) {
        shade_squad.hide();
        move_squad_to(shade_squad.name, shade_squad.direction, shade_squad.columns, shade_squad.x, shade_squad.y);
    };
}
