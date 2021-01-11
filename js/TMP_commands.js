function create_squad(name, size, color, unit_size, units_interval) {
    if (!SOCKET_BATTLE) {
        console.log('have no socket battle');
        return;
    }

    let data = JSON.stringify({
        'Command': 'create_squad',
        'Parameters': {
            'Name': name,
            'Size': size,
            'Color': color,
            'UnitSize': unit_size,
            'UnitsInterval': units_interval,
        }
    });
    SOCKET_BATTLE.send(data);
}


function move_squad_to(name, direction, columns, x, y) {
    if (!SOCKET_BATTLE) {
        console.log('have no socket battle');
        return;
    }

    let data = JSON.stringify({
        'Command': 'move_squad',
        'Parameters': {
            'Name': name,
            'Direction': direction,
            'Columns': columns,
            'ToX': x,
            'ToY': y,
        }
    });
    SOCKET_BATTLE.send(data);
}
