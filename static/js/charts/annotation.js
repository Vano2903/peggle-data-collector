"use strict"
google.charts.load('current', { 'packages': ['annotatedtimeline'] });

function drawChart(who, rows) {
    let chart;
    let data = new google.visualization.DataTable();
    data.addColumn({ type: 'date', id: 'Date' });
    if (who == "s") {
        data.addColumn({ type: 'number', id: 'punti synergo' });
        chart = new google.visualization.AnnotatedTimeLine(document.getElementById('sPointsChart'));
    } else if (who == "r") {
        data.addColumn({ type: 'number', id: 'punti redez' });
        chart = new google.visualization.AnnotatedTimeLine(document.getElementById('rPointsChart'));
    } else {
        data.addColumn({ type: 'number', id: 'punti synergo' });
        data.addColumn({ type: 'number', id: 'punti redez' });
        chart = new google.visualization.AnnotatedTimeLine(document.getElementById('srPointsChart'));
    }
    data.addColumn({ type: 'string', id: 'titolo partita' });
    data.addRows(rows);

    var options = {
        displayAnnotations: true,
        allowHtml: true,
        table: {
            sortAscending: true
        }
    };

    chart.draw(data, options);
}