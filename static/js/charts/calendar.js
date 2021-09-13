google.charts.load("current", { packages: ["calendar"] });
google.charts.setOnLoadCallback(calendar);

//array of all the years with at least a commit it
var years;
//array of commits
var commits;
//delay for to run on resize (check onResize handler)
var delay = 250;
//used to run the function if the delay has passed
var throttled = false;

/**
 * given users information fetch what years the user committed in
 * @param {object} user user information, saved on local storage and generated on user log in
 * @returns array of strings containing all the years the user has at least a commit in
 */
async function getUsersCommitsYears(user) {
    var res = await fetch('/commit/years', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    var resp = await res.text();
    var years = resp.split(";");
    console.log(years);
    return years
}

/**
 * given a user informations and a specific year (the user must have at least a commit in the given year)
 * the function will generate the dataset that will be used by the google charts api to draw the calendar chart
 * @param {object} user user information, saved on local storage and generated on user log in
 * @param {string} year the year you want to fetch from the user 
 * @returns dataset used by google chart to draw the calendar chart
 */
async function genCalendarCommits(user, year) {
    user.year = parseInt(year);
    var res = await fetch('/commit/year', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    var respJson = await res.json();
    console.log(respJson);
    var dataSet = [];

    respJson.forEach(function (item) {
        dataSet.push([new Date(item.date), item.totals])
    });

    return dataSet;
}

/**
 * this functino will draw the calendar chart on the stats page in the users page
 * @param {object} opt the option for the google chart
 * @param {array[]} dataSet the dataset for the google chart
 */
function drawCommitChart(opt, dataSet) {
    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn({ type: 'date', id: 'Date' });
    dataTable.addColumn({ type: 'number', id: 'Won/Loss' });
    dataTable.addRows(dataSet);

    var chart = new google.visualization.Calendar(document.getElementById('commits-graph'));

    chart.draw(dataTable, opt);
}

/**
 * given a user information, the function will return the number of commits made
 * @param {object} user user information, saved on local storage and generated on user log in
 * @returns a number which is the ammount of commits
 */
async function getTotalCommits(user) {
    var res = await fetch('/commit/totCommits', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    var resp = await res.text();
    var respJson = JSON.parse(resp)
    console.log(respJson.totalCommits)
    return respJson.totalCommits
}

/**
 * this function will define what the google chart will render, 
 * if a specific year or the last year a user made a commit
 */
async function calendar() {
    var width = $(window).width();
    var calendarOptions = genCalendarOptions(width, user.username)

    if (commits == undefined) {
        commits = await genCalendarCommits(user, years[years.length - 1])
    }

    drawCommitChart(calendarOptions, commits);
}

/**
 * this function will draw the buttons that let the user choose which year check for the commits made 
 * (the button will render the calendar with a specific year)
 * @param {array} years array of strings with all the possible years a user can choose
 */
function drawYearsButtons(years) {
    var buttonContainer = document.getElementById("calendar-buttons");
    buttonContainer.innerHTML = "";
    years.forEach(function (item) {
        var but = document.createElement("button");
        but.innerHTML = item;
        but.value = item;
        but.className = "btn btn-primary";
        but.addEventListener('click', async function (e) {
            commits = await genCalendarCommits(user, item)
            calendar()
        }, false);
        buttonContainer.appendChild(but);
        buttonContainer.appendChild(document.createElement("br"));
        buttonContainer.appendChild(document.createElement("br"));
    });
}

/**
 * this function will run on resize of the windows, the reason if to make the calendar responsive 
 * without needing to reaload the page the function to draw a calendar seems to require a lot of cpu 
 * from the user and this function optimise the amount of time the function is called on a resize event
 */
window.addEventListener('resize', function () {
    if (!throttled) {
        if (locationSaver() == "stats") {
            calendar()
        }
        throttled = true;
        setTimeout(function () {
            throttled = false;
        }, delay);
    }
});

/**
 * this function initialize the stats page with the title and call
 * drawYearsButtons given the array of years to draw the buttons,
 * the calendar function to draw the chart and tells the user the ammount of commits made
 */
async function initCalendar() {
    years = await getUsersCommitsYears(user);
    drawYearsButtons(years)
    calendar()
    document.getElementById("total-commits").innerHTML = "numero totale di contributi: " + await getTotalCommits(user)
}

/**
 * on load check if the last page saw was the stats page,
 * if so run initCalendar
 */
$(document).ready(async function () {
    if (locationSaver() == "stats") {
        await initCalendar()
    }
})

/**
 * this function is used to make the calendar responsive
 * @param {number} width integer which is the width of the page the calendar graph is inside
 * @param {string} name is the name of the user
 * @returns a object that contains the charts options (look google charts options)
 */
function genCalendarOptions(width, name) {
    var opt = {}
    opt.title = "Commits di " + name;
    opt.height = 250;

    var calendar = {};
    if (width < 500) {
        calendar.cellSize = 5;
    } else if (width < 600) {
        calendar.cellSize = 8;
    } else if (width < 800) {
        calendar.cellSize = 11;
    } else if (width < 1000) {
        calendar.cellSize = 12;
    } else if (width < 1200) {
        calendar.cellSize = 16;
    } else {
        calendar.cellSize = 19;
    }

    opt.calendar = Object.assign(calendar)
    return opt
}