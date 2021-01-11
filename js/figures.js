class TMPCircle{

    constructor(id) {
        this.div = document.createElement('div');
        this.id = id;
        this.class_ = 'circle';
        this.color = '#FF0AFA';
        this.size = 30;
        this.x = 15;
        this.y = 15;
        this.squad = undefined;
    }

    hide() {
        this.div.style.display = "none"
    }

    unhide() {
        this.div.style.display = "block"
    }

    TMP_draw() {
        MY_BODY.field.append(this.div);
    }

    TMP_select() {
        this.div.style.borderStyle = 'dotted';
    }

    TMP_unselect() {
        this.div.style.borderStyle = 'none';
    }

    erase() {
        this.div.remove();
    }

    get class_() {
        return this.div.className;
    }

    set class_(value) {
        this.div.className = value;
    }

    get id() {
        return this.div.id;
    }

    set id(value) {
        this.div.id = value;
    }

    get x() {
        return this._x;
    }

    set x(value) {
        this._x = value;
        this.div.style.left = value - this.size / 2 + 'px';
    }

    get y() {
        return this._y;
    }

    set y(value) {
        this._y = value;
        this.div.style.top = value - this.size / 2 + 'px';
    }

    get color() {
        return this.div.style.background;
    }

    set color(value) {
        this.div.style.background = value;
    }

    get size() {
        return this._size;
    }

    set size(value) {
        this._size = value;
        this.x = this.x;
        this.y = this.y;
        this.div.style.width = value + 'px';
        this.div.style.height = value + 'px';
    }
}


function killCircle(circle) {
    circle.erase();
    circle.squad.kill_circle(circle);
    CIRCLES.delete(circle.id);
}


function createCircle(obj) {
    let circle = new TMPCircle(obj.Id);
    circle.color = obj.Color;
    circle.size = obj.Size;
    circle.x = obj.X;
    circle.y = obj.Y;

    return circle;
};


function createShadeCircle(obj) {
    let circle = new TMPCircle(obj.Id);
    circle.color = obj.Color+'88';
    circle.size = obj.Size;
    circle.x = obj.X;
    circle.y = obj.Y;

    return circle;
};


class TMPLine{
    constructor(id) {
        this.div = document.createElement('div');
        this.id = id;
        this.class_ = 'line';
        this.color = '#000000';
        this.length = 100;
        this.size = 2;
        this.x = 15;
        this.y = 15;
        this.direction = 0;
    }

    hide() {
        this.div.style.display = "none"
    }

    unhide() {
        this.div.style.display = "block"
    }

    TMP_draw() {
        MY_BODY.field.append(this.div);
    }

    erase() {
        this.div.remove();
    }

    get size() {
        return this._size;
    }

    set size(value) {
        this._size = value;
        this.x = this.x;
        this.y = this.y;
        this.div.style.borderStyle = 'solid';
        this.div.style.borderBottomWidth = value + 'px';
    }

    get direction() {
        return this._direction;
    }

    set direction(value) {
        this._direction = value;
        this.div.style.transform = 'rotate(' + value + 'deg)';
    }

    get id() {
        return this.div.id;
    }

    set id(value) {
        this.div.id = value;
    }

    get color() {
        return this.div.style.borderColor;
    }

    set color(value) {
        this.div.style.borderColor = value;
    }

    get class_() {
        return this.div.className;
    }

    set class_(value) {
        this.div.className = value;
    }

    get x() {
        return this._x;
    }

    set x(value) {
        this._x = value;
        this.div.style.left = value - this.size + 'px';
    }

    get y() {
        return this._y;
    }

    set y(value) {
        this._y = value;
        this.div.style.top = value - this.size + 'px';
    }

    get length() {
        return this._length;
    }

    set length(value) {
        this._length = value;
        this.div.style.width = value + 'px';
    }
}
