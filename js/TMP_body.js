class Body {
    constructor() {
        this.div = document.body;
        this.field = document.getElementById("field");
        this.prevent_default_right_click();
        this.onmousedown_ = mousedown_handler;
        this.onmouseup_ = mouseup_handler;
        this.onkeypress_ = onkeypress_handler;
        this.move_start_x = 0;
        this.move_start_y = 0;
        this.move_end_x = 0;
        this.move_end_y = 0;
        this.TMP_old_columns = 0;
        this.field_size = 2000;
        this.zoom_level = 1;

        this.field.style.width = this.field_size + 'px'
        this.field.style.height = this.field_size + 'px'
        this.field.style.borderStyle = 'solid'
        this.field.addEventListener('mousemove', TMP_mouse_move);
    }

    left_mouse_move(x_unscaled, y_unscaled) {
        if (SELECTED_SHADE_SQUAD) {
            let scaled_coordinates = this.scale_field_coordinates(x_unscaled, y_unscaled);
            let x = scaled_coordinates[0];
            let y = scaled_coordinates[1];
            let length = ((SELECTED_SHADE_SQUAD.x - x) ** 2 + (SELECTED_SHADE_SQUAD.y - y) ** 2) ** 0.5;
            let norm_x = x - SELECTED_SHADE_SQUAD.x;
            let norm_y = y - SELECTED_SHADE_SQUAD.y;
            let direction_cos = (0 * norm_x - 1 * norm_y) / (1 * length);
            let direction_degree = Math.sign(norm_x) * to_degree(Math.acos(direction_cos));
            let step = 10 / this.zoom_level;
            let length_ratio = 0;

            if (length > 80) {
                length_ratio = Math.ceil((140 - length)/step);
            }

            let new_columns = this.TMP_old_columns + length_ratio

            SELECTED_SHADE_SQUAD.columns = new_columns;
            SELECTED_SHADE_SQUAD.rotate_to(direction_degree);
        }
    }

    TMP_mouse_move(x_unscaled, y_unscaled) {
        if (SELECTED_SHADE_SQUAD) {
            let scaled_coordinates = this.scale_field_coordinates(x_unscaled, y_unscaled);

            SELECTED_SHADE_SQUAD.unhide();
            SELECTED_SHADE_SQUAD.move_to(scaled_coordinates[0], scaled_coordinates[1]);
            SELECTED_SHADE_SQUAD.TMP_draw_shape();
        }
    }

    right_click_field(x_unscaled, y_unscaled) {
        if (SELECTED_SHADE_SQUAD) {
            let scaled_coordinates = this.scale_field_coordinates(x_unscaled, y_unscaled);
            SELECTED_SHADE_SQUAD.move_to(scaled_coordinates[0], scaled_coordinates[1]);
        }
    }

    scale_field_coordinates(x_unscaled, y_unscaled) {
        return [x_unscaled / this.zoom_level, y_unscaled / this.zoom_level];
    }

    left_click_field(event) {
        if (SELECTED_SQUAD) {
            SELECTED_SQUAD.TMP_unselect();
        }
        if (SELECTED_SHADE_SQUAD) {
            SELECTED_SHADE_SQUAD.TMP_unselect();
        }
    }

    left_click_circle(event) {
        let squad = CIRCLES.get(event.target.id).squad;
        let shade_squad = SHADE_CIRCLES.get(event.target.id).squad;

        if (squad === undefined && shade_squad === undefined) {
            alert('squad undefined for circle: ' + event.target.id);
            return;
        };

        if (USER_NAME != squad.TMPOwnerName) {
            if (SELECTED_SQUAD) {
                SELECTED_SQUAD.TMP_unselect();
            };
            if (SELECTED_SHADE_SQUAD) {
                SELECTED_SHADE_SQUAD.TMP_unselect();
            };

            return;
        };

        if (squad.is_select === false) {
            squad.TMP_select();
            shade_squad.TMP_select();
        } else {
            squad.TMP_unselect();
            shade_squad.TMP_unselect();
        };
    }

    zoom_in() {
        this.zoom_level = this.zoom_level + 0.2;
    }

    zoom_out() {
        this.zoom_level = this.zoom_level - 0.2;
    }

    get zoom_level() {
        return this.zoom_level_;
    }

    set zoom_level(value) {
        this.zoom_level_ = value;
        this.field.style.zoom = this.zoom_level_;
    }

    set onmousedown_(value) {
        this.field.onmousedown = value;
    }

    set onmouseup_(value) {
        this.field.onmouseup = value;
    }

    set onkeypress_(value) {
        this.div.onkeypress = value;
    }

    prevent_default_right_click() {
        window.addEventListener('contextmenu', function (e) {
            e.preventDefault();
        }, false);
    }
}


function mousedown_handler() {
    if (event.target.className === 'button') {
        return;
    }

    MY_BODY.move_start_x = event.pageX;
    MY_BODY.move_start_y = event.pageY;
    if (SELECTED_SQUAD) {
        MY_BODY.TMP_old_columns = SELECTED_SQUAD.columns;
    }

    if (event.which === 1) {
        MY_BODY.field.addEventListener('mousemove', left_mouse_move);
        MY_BODY.field.removeEventListener('mousemove', TMP_mouse_move);
    }
};


function mouseup_handler() {
    if (event.target.className === 'button') {
        return;
    }

    MY_BODY.move_end_x = event.pageX;
    MY_BODY.move_end_y = event.pageY;

    if (event.which === 1) {
        MY_BODY.field.removeEventListener('mousemove', left_mouse_move);

        if (SELECTED_SHADE_SQUAD != null) {
            move_squad_to(SELECTED_SHADE_SQUAD.name, SELECTED_SQUAD.direction, SELECTED_SHADE_SQUAD.columns, SELECTED_SHADE_SQUAD.x, SELECTED_SHADE_SQUAD.y);
        }

        left_click_handler(event);
        MY_BODY.field.addEventListener('mousemove', TMP_mouse_move);
    } else if (event.which === 3) {
        right_click_handler(event);
        if (SELECTED_SHADE_SQUAD != null) {
            move_squad_to(SELECTED_SHADE_SQUAD.name, SELECTED_SQUAD.direction, SELECTED_SHADE_SQUAD.columns, SELECTED_SHADE_SQUAD.x, SELECTED_SHADE_SQUAD.y);
            SELECTED_SQUAD.TMP_unselect();
            SELECTED_SHADE_SQUAD.TMP_unselect();
        }
    }
};


function onkeypress_handler(event) {
    if (event.code === 'Equal') {
        MY_BODY.zoom_in();
    } else if (event.code === 'Minus') {
        MY_BODY.zoom_out();
    }
};


function left_mouse_move(event) {
    MY_BODY.left_mouse_move(event.pageX, event.pageY);
};


function TMP_mouse_move(event) {
    MY_BODY.TMP_mouse_move(event.pageX, event.pageY);
};


function left_click_handler(event) {
    if (event.target.className === 'circle') {
        MY_BODY.left_click_circle(event);
    } else {
        MY_BODY.left_click_field(event);
    }
};


function right_click_handler(event) {
    if (event.target.className === 'circle') {
        return;
    } else {
        MY_BODY.right_click_field(event.pageX, event.pageY);
    }
};
