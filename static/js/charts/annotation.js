"use strict"
google.charts.load('current', { 'packages': ['annotatedtimeline'] });

function drawChart(who, rows) {
    let chart;
    let data = new google.visualization.DataTable();
    var options = new Object();
    data.addColumn('date', 'Date');
    if (who == "s") {
        // options = ["blue"]
        data.addColumn('number', 'punti synergo');
        chart = new google.visualization.AnnotatedTimeLine(document.getElementById('sPointsChart'));
    } else if (who == "r") {
        // options = ["red"]
        data.addColumn('number', 'punti redez');
        chart = new google.visualization.AnnotatedTimeLine(document.getElementById('rPointsChart'));
    } else {
        // options = ["red", "blue"]
        data.addColumn('number', 'punti synergo');
        data.addColumn('number', 'punti redez');
        chart = new google.visualization.AnnotatedTimeLine(document.getElementById('srPointsChart'));
    }
    data.addColumn('string', 'titolo partita');
    data.addRows(rows);

    options.displayAnnotations = true;
    options.allowHtml = true;
    options.table = { sortAscending: true }

    chart.draw(data, options);
}