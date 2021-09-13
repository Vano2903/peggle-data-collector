"use strict"
google.charts.load('current', { 'packages': ['corechart'] });

/**
 * this function will draw the pie chart (in stats page)
 * @param {[string]} who if "s" the func will draw the annotation chart of synergo, if "r" the redez one, if whatever the combination of both
 * @param {[Array[]]} rows the data fromatted in google chart notation (will be a array of arrays that has a date, number and string)
 * @param {[string]} which define which section in the stats page the function will draw the graph
 */
function drawPieChart(who, rows, which) {
    let chart = new google.visualization.PieChart(document.getElementById(who + which + 'Chart'));
    if (who != "sr") {
        var data = new google.visualization.arrayToDataTable(rows);
    }
    var options = new Object();

    let name;
    if (who == "s") {
        name = "da synergo";
    } else if (who == "r") {
        name = "da redez";
    } else {
        name = "(u = synergo, current = redez)"
        if (which == "Char") {
            data = chart.computeDiff(new google.visualization.arrayToDataTable(sCharData), new google.visualization.arrayToDataTable(rCharData));
        } else {
            data = chart.computeDiff(new google.visualization.arrayToDataTable(sFEData), new google.visualization.arrayToDataTable(rFEData));
        }
    }

    if (which == "Char") {
        options = { title: "Percentuale di scelta dei personaggi " + name }
    } else {
        options = { title: "Percentuale di feste estreme fatte " + name }
    }

    chart.draw(data, options);
}