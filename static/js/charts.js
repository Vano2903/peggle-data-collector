google.charts.load("current", { packages: ["calendar"] });
google.charts.setOnLoadCallback(resizeChart);

function genCalendarCommits() {
    fetch("/")
}

function drawCommitChart(opt, dataSet) {
    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn({ type: 'date', id: 'Date' });
    dataTable.addColumn({ type: 'number', id: 'Won/Loss' });
    // dataTable.addRows(genCalendarCommits());
    dataTable.addRows([
        [new Date(2012, 3, 13), 37032],
        [new Date(2012, 3, 14), 38024],
        [new Date(2012, 3, 15), 38024],
        [new Date(2012, 3, 16), 38108],
        [new Date(2012, 11, 17), 38229]
    ]);

    var chart = new google.visualization.Calendar(document.getElementById('commits-graph'));

    // var options = {
    //     calendar: { cellSize: 20},
    //     title: "Commits",
    //     // height: 800
    // };

    // var options = opt
    //calendar: { cellSize: 10 },
    // var options = {
    //     'title': 'Google results searching for "data" 2002 - 2012',
    //     'width': '100%',
    //     'height': 900,
    //     'curveType': "function",
    //     'backgroundColor': "#f1f2f2",
    //     'chartArea': { 'width': '100%' },
    //     'axisTitlesPosition': 'in',
    //     'vAxis': { 'textPosition': 'in' },
    //     'titleTextStyle': { 'color': '#809ECE' },
    //     'hAxis.textStyle': { 'color': '#809ECE' },
    //     'colors': ['#809ECE']
    // };


    chart.draw(dataTable, opt);
}

function resizeChart() {
    var width = $(window).width();
    var calendarOptions = genCalendarOptions(width, "vano")

    drawCommitChart(calendarOptions, width);
}

$(window).resize(function () {
    resizeChart()
});

function genCalendarOptions(width, name) {
    var opt = {}
    opt.title = "Commits " + name;
    opt.height = 400;

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