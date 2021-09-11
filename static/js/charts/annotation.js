"use strict"
google.charts.load('current', { 'packages': ['annotatedtimeline'] });

function drawAnnotationChart(who, rows, which) {
    let chart = new google.visualization.AnnotatedTimeLine(document.getElementById(who + which + 'Chart'));
    let data = new google.visualization.DataTable();
    var options = new Object();
    data.addColumn('date', 'Date');

    if (who == "s") {
        options.colors = ["blue"]
        data.addColumn('number', 'punti synergo');
    } else if (who == "r") {
        options.colors = ["red"]
        data.addColumn('number', 'punti redez');
    } else {
        data.addColumn('number', 'punti synergo');
        data.addColumn('number', 'punti redez');
    }

    data.addColumn('string', 'titolo partita');
    data.addRows(rows);

    options.displayAnnotations = true;
    options.allowHtml = true;
    options.table = { sortAscending: true }

    chart.draw(data, options);
}