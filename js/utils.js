function change_coordinates_system_and_rotate(x, y, x_center, y_center, cos_, sin_) {
    // convert coordinate system. from global to local
    // rotation center is a formation center
    // rotate formation relative to formation center
    result = rotate_in_current_system(x-x_center, y-y_center, cos_, sin_);

    // convert coordinate system. from local to global
    return [result[0]+x_center, result[1]+y_center];
}


function rotate_in_current_system(x, y, cos_, sin_) {
    // rotate vectors relative to "current" coordinates system center
    // x=0 and y=0 is center of current coordinates system
    return [x * cos_ - y * sin_, x * sin_ + y * cos_];
}


function to_radians(degree) {
    return degree * Math.PI / 180;
}


function to_degree(radians) {
    return radians * 180 / Math.PI;
}


function timer(ms) {
    return new Promise(res => setTimeout(res, ms));
}


async function load () { // We need to wrap the loop into an async function for this to work
    for (var i = 0; i < 100; i++) {
        console.log(i);
        circle_1.x += 1;
        circle_1.size += 1;
        console.log(circle_1.x);
        console.log(circle_1.size);
        await timer(100); // then the created Promise can be awaited
    }
}


function tmp_make_units_battle_response(units) {
    let TMP_circles_ids = Array(units.length);
    let new_circle = null;
    let i = 0;

    for (var [unit_id, unit] of Object.entries(units)) {
        new_circle = createCircle(unit);
        CIRCLES.set(unit_id, new_circle);
        TMP_circles_ids[i] = unit_id;
        i += 1;
    }

    return TMP_circles_ids;
}


function tmp_make_units_shades_battle_response(units) {
    let TMP_circles_ids = Array(units.length);
    let new_circle = null;
    let i = 0;

    for (var [unit_id, unit] of Object.entries(units)) {
        new_circle = createShadeCircle(unit);
        SHADE_CIRCLES.set(unit_id, new_circle);
        TMP_circles_ids[i] = unit_id;
        i += 1;
    }

    return TMP_circles_ids;
}
