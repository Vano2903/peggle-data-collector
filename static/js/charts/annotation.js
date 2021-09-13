"use strict"
google.charts.load('current', { 'packages': ['annotatedtimeline'] });

//given who ("s" if synergo, "r" if redez, whatever if both)
//rows is the table formatter

/**
 * this function will draw the annotation chart (in stats page)
 * @param {[string]} who if "s" the func will draw the annotation chart of synergo, if "r" the redez one, if whatever the combination of both
 * @param {[Array[]]} rows the data fromatted in google chart notation (will be a array of arrays that has a date, number and string)
 * @param {[string]} which define which section in the stats page the function will draw the graph
 */
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