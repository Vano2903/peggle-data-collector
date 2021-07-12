let buttons = [document.getElementById("gameBut1"), document.getElementById("gameBut2"), document.getElementById("gameBut3")]
let corrects = []

function dangerOutlineAllButtons() {
    for (let i = 0; i < buttons.length; i++) {
        if (corrects.includes(i)) {
            buttons[i].classList = "btn btn-outline-success";
        } else {
            buttons[i].classList = "btn btn-outline-danger";
        }
    }
}

function buttonsCorrectArea() {
    for (let i = 0; i < buttons.length; i++) {
        if (checkIfAreaIsComplete(i) && !corrects.includes(i)) {
            corrects.push(i)
            buttons[i].classList = "btn btn-outline-success"
        }
    }
}

function onSection(index) {
    dangerOutlineAllButtons()
    buttons[index].classList = "btn btn-info";
}
