class BaseSquad {

    constructor(obj, circles_ids) {
        this.circles = this.TMP_get_circles(circles_ids);
        this.name = obj.Name;
        this.formation = [];
        this.columns = obj.Columns;
        this.lines = obj.Lines;
        this.circle_size = obj.UnitSize;
        this.interval = obj.Interval;
        this.x = obj.X;
        this.y = obj.Y;
        this.width = obj.Width;
        this.height = obj.Height;
        this.direction = obj.Direction;
        this.is_select = false;

        this._link_squad_to_circles();
        this._init_direction();
        this._init_corners();
        this._TMP_calculate_shape();
        this._TMP_init_formation(obj);
    }

    move_to(x, y) {
       this.x = x;
       this.y = y;

       this._init_formation();
       this.TMP_draw();
    }

    rotate_to(degree) {
        this.direction = degree;
        this.TMP_draw();
    }

    TMP_draw_shape() {
        this._TMP_calculate_shape();
        this._TMP_draw_corners();
        this._TMP_draw_direction();
    }

    TMP_draw() {
        for (let circle of this.circles) {
            circle.TMP_draw();
        }
        this.TMP_draw_shape()
    }

    _link_squad_to_circles() {
        for (let circle of this.circles) {
            circle.squad = this;
        }
    }

    _init_formation() {
        let line = 0;
        let column = 0;
        let x = 0.0;
        let y = 0.0;
        let result = [];
        let sin_ = Math.sin(to_radians(this.direction));
        let cos_ = Math.cos(to_radians(this.direction));
        this.formation = new Array(this.lines).fill(null).map(() => new Array(this.columns).fill(null));

        for (let circle of this.circles) {
            // initialize formation
            this.formation[line][column] = circle;
            x = this.first_x + (column * circle.size) + (column * this.interval) + (circle.size / 2);
            y = this.first_y + (line * circle.size) + (line * this.interval) + (circle.size / 2);

            result = change_coordinates_system_and_rotate(x, y, this.x, this.y, cos_, sin_);

            circle.x = result[0];
            circle.y = result[1];

            column += 1;
            if (column >= this.columns) {
                column = 0;
                line += 1;
            }
        }
    }

    _init_corners() {
        this.central_circle = new TMPCircle("center");
        this.top_left_circle = new TMPCircle("top_left_corner");
        this.bottom_right_circle = new TMPCircle("bottom_right_corner");
        this.top_right_circle = new TMPCircle("top_right_corner");
        this.bottom_left = new TMPCircle("bottom_left");

        let circle_size = 4;
        let corner_color = "#000000";

        this.top_left_circle.size = circle_size;
        this.top_left_circle.color = corner_color;
        this.central_circle.size = circle_size;
        this.central_circle.color = corner_color;
        this.bottom_right_circle.size = circle_size;
        this.bottom_right_circle.color = corner_color;
        this.top_right_circle.size = circle_size;
        this.top_right_circle.color = corner_color;
        this.bottom_left.size = circle_size;
        this.bottom_left.color = corner_color;

        this.central_circle.TMP_draw();
        this.top_left_circle.TMP_draw();
        this.bottom_right_circle.TMP_draw();
        this.top_right_circle.TMP_draw();
        this.bottom_left.TMP_draw();
    }

    _init_direction() {
        this._direction_line = new TMPLine(this.name + '_direction_line');
        this._direction_line.TMP_draw();
    }

    _TMP_calculate_shape() {
        this.width = this.columns * this.circle_size + (this.columns - 1) * this.interval;
        this.height = this.lines * this.circle_size + (this.lines - 1) * this.interval;
    }

    _TMP_draw_direction() {
        this._direction_line.x = this.x;
        this._direction_line.y = this.y;
        this._direction_line.length = this.height / 2 + this.circle_size;
        this._direction_line.direction = this._direction - 90;
    }

    _TMP_draw_corners() {
        let cos_ = Math.cos(to_radians(this.direction));
        let sin_ = Math.sin(to_radians(this.direction));

        this.central_circle.x = this.x;
        this.central_circle.y = this.y;

        result = change_coordinates_system_and_rotate(this.first_x, this.first_y, this.x, this.y, cos_, sin_);
        this.top_left_circle.x = result[0];
        this.top_left_circle.y = result[1];

        result = change_coordinates_system_and_rotate(this.last_x, this.last_y, this.x, this.y, cos_, sin_);
        this.bottom_right_circle.x = result[0];
        this.bottom_right_circle.y = result[1];

        result = change_coordinates_system_and_rotate(this.last_x, this.first_y, this.x, this.y, cos_, sin_);
        this.top_right_circle.x = result[0];
        this.top_right_circle.y = result[1];

        result = change_coordinates_system_and_rotate(this.first_x, this.last_y, this.x, this.y, cos_, sin_);
        this.bottom_left.x = result[0];
        this.bottom_left.y = result[1];
    }

    erase() {
        for (let circle of this.circles) {
            circle.erase();
        };

        this.central_circle.erase();
        this.top_left_circle.erase();
        this.bottom_right_circle.erase();
        this.top_right_circle.erase();
        this.bottom_left.erase();
        this._direction_line.erase();
    }

    hide() {
        for (let circle of this.circles) {
            circle.hide();
        };

        this.central_circle.hide();
        this.top_left_circle.hide();
        this.bottom_right_circle.hide();
        this.top_right_circle.hide();
        this.bottom_left.hide();
        this._direction_line.hide();
    }

    unhide() {
        for (let circle of this.circles) {
            circle.unhide();
        };

        this.central_circle.unhide();
        this.top_left_circle.unhide();
        this.bottom_right_circle.unhide();
        this.top_right_circle.unhide();
        this.bottom_left.unhide();
        this._direction_line.unhide();
    }

    get columns() {
        return this._columns;
    }

    set columns(value) {
        if (value < 2 || value > 30) {
            return;
        }

        this._columns = value;
        this._lines = Math.ceil(this.circles.length / value);
        this._TMP_calculate_shape();
        if (this.formation !== undefined) {
            this._init_formation();
        }
    }

    get lines() {
        return this._lines;
    }

    set lines(value) {
        if (value < 2 || value > 30) {
            return;
        }

        this._lines = value;
        this._columns = Math.ceil(this.circles.length / value);
        this._TMP_calculate_shape();
        if (this.formation !== undefined) {
            this._init_formation();
        }
    }

    get first_x() {
        return this.x - (this.width / 2);
    }

    get first_y() {
        return this.y - (this.height / 2);
    }

    get last_x() {
        return this.x + (this.width / 2);
    }

    get last_y() {
        return this.y + (this.height / 2);
    }

    get direction() {
        return this._direction;
    }

    set direction(value) {
        let cos_ = Math.cos(to_radians(value - this._direction));
        let sin_ = Math.sin(to_radians(value - this._direction));

        this._direction = value;

        for (let circle of this.circles) {
            result = change_coordinates_system_and_rotate(circle.x, circle.y, this.x, this.y, cos_, sin_);
            circle.x = result[0];
            circle.y = result[1];
        }
    }
}


class Squad extends BaseSquad {
    constructor(obj, circles_ids) {
        super(obj, circles_ids);
    }

    TMP_select() {
        this.is_select = true;
        for (let circle of this.circles) {
            circle.TMP_select();
        }
        SELECTED_SQUAD = this;
    }

    TMP_unselect() {
        this.is_select = false;
        for (let circle of this.circles) {
            circle.TMP_unselect();
        }
        SELECTED_SQUAD = null;
    }

    TMP_get_circles(circles_ids) {
        let circles = Array(circles_ids.length);

        for (var id = 0; id < circles_ids.length; id++) {
            circles[id] = CIRCLES.get(circles_ids[id]);
        }

        return circles;
    }

    kill_circle(circle) {
        var index_to_remove = this.circles.indexOf(circle);

        if (index_to_remove > -1) {
            this.circles.splice(index_to_remove, 1);
        };

        for (var line = 0; line < this.lines; line++) {
            index_to_remove = this.formation[line].indexOf(circle);
            if (index_to_remove > -1) {
                this.formation[line][index_to_remove] = null;
                break;
            };
        };

        if (this.circles == 0) {
            this.erase();
            SQUADS.delete(this.name);
        };
    }

    _TMP_init_formation(obj) {
        this.formation = new Array(this.lines).fill(null).map(() => new Array(this.columns).fill(null));

        for (var line = 0; line < this.lines; line++) {
            for (var column = 0; column < this.columns; column++) {
                let circle = obj.Formation[line][column]
                if (circle == null) {
                    continue;
                };
                this.formation[line][column] = CIRCLES.get(circle.Id);
            };
        };
    }
}


class ShadeSquad extends BaseSquad {
    constructor(obj, circles_ids) {
        super(obj, circles_ids);
    }

    TMP_select() {
        this.is_select = true;
        for (let circle of this.circles) {
            circle.TMP_select();
        }
        SELECTED_SHADE_SQUAD = this;
    }

    TMP_unselect() {
        this.is_select = false;
        for (let circle of this.circles) {
            circle.TMP_unselect();
        }
        SELECTED_SHADE_SQUAD = null;
    }

    TMP_get_circles(circles_ids) {
        let circles = Array(circles_ids.length);

        for (var id = 0; id < circles_ids.length; id++) {
            circles[id] = SHADE_CIRCLES.get(circles_ids[id]);
        }

        return circles;
    }

    kill_circle(circle) {
        var index_to_remove = this.circles.indexOf(circle);

        if (index_to_remove > -1) {
            this.circles.splice(index_to_remove, 1);
        };

        for (var line = 0; line < this.lines; line++) {
            index_to_remove = this.formation[line].indexOf(circle);
            if (index_to_remove > -1) {
                this.formation[line][index_to_remove] = null;
                break;
            };
        };

        if (this.circles == 0) {
            this.erase();
            SHADE_SQUADS.delete(this.name);
        };
    }

    _TMP_init_formation(obj) {
        this.formation = new Array(this.lines).fill(null).map(() => new Array(this.columns).fill(null));

        for (var line = 0; line < this.lines; line++) {
            for (var column = 0; column < this.columns; column++) {
                let circle = obj.Formation[line][column]
                if (circle == null) {
                    continue;
                };
                this.formation[line][column] = SHADE_CIRCLES.get(circle.Id);
            };
        };
    }
}
