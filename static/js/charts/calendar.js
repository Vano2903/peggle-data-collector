google.charts.load("current", { packages: ["calendar"] });
google.charts.setOnLoadCallback(calendar);

var years
var commits

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

async function genCalendarCommits(user, year) {
    user.year = parseInt(year);
    var res = await fetch('/commit/year', {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    });
    var resp = await res.text();
    console.log(resp);
    var respJson = JSON.parse(resp);
    console.log(respJson);
    var dataSet = [];

    respJson.forEach(function (item) {
        dataSet.push([new Date(item.CreatedAt), item.Totals])
    });

    return dataSet;
}

function drawCommitChart(opt, dataSet) {
    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn({ type: 'date', id: 'Date' });
    dataTable.addColumn({ type: 'number', id: 'Won/Loss' });
    dataTable.addRows(dataSet);

    var chart = new google.visualization.Calendar(document.getElementById('commits-graph'));

    chart.draw(dataTable, opt);
}

async function calendar() {
    var width = $(window).width();
    var calendarOptions = genCalendarOptions(width, user.username)

    if (commits == undefined) {
        commits = await genCalendarCommits(user, years[years.length - 1])
    }

    drawCommitChart(calendarOptions, commits);
}

$(window).resize(function () {
    calendar()
});

function drawYearsButtons(years) {
    var buttonContainer = document.getElementById("calendar-buttons");
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

$(document).ready(async function () {
    years = await getUsersCommitsYears(user);
    drawYearsButtons(years)
})

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