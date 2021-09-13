//array with the 3 buttons used to change the game section
let buttons = [document.getElementById("gameBut1"), document.getElementById("gameBut2"), document.getElementById("gameBut3")]
//array that stores the button with a completed section
let corrects = []

/**
 * this function will set to danger (red) all the button that are not stored in the corrects array
 */
function dangerOutlineAllButtons() {
    for (let i = 0; i < buttons.length; i++) {
        if (corrects.includes(i)) {
            buttons[i].classList = "btn btn-outline-success";
        } else {
            buttons[i].classList = "btn btn-outline-danger";
        }
    }
}

/**
 * set to green the buttons with a completed section that are not already in the corrects array
 *  (when a new button becomes correct will be added to the corrects array)
 */
function buttonsCorrectArea() {
    for (let i = 0; i < buttons.length; i++) {
        if (checkIfAreaIsComplete(i) && !corrects.includes(i)) {
            corrects.push(i)
            buttons[i].classList = "btn btn-outline-success"
        }
    }
}

/**
 * will set to blue (info) the button which the user is looking at
 * @param {number} index the index of the game section
 */
function onSection(index) {
    dangerOutlineAllButtons()
    buttons[index].classList = "btn btn-info";
}
