const gameData = []

const stats = {
    "generic": {
        "totTimeWatched": 4984,
        "totEpisodesStored": 9,
        "collaborators": [
            "vano",
            "MoraGames"
        ]
    },
    "synergo": {
        "totPoints": 1219290,
        "totn25": 33,
        "totWins": 4,
        "FEstats": {
            "n5000": 0,
            "n25000": 1,
            "n50000": 0,
            "totPointsMade": 25000
        },
        "charStats": {
            "cas": 0,
            "uni": 6,
            "zuc": 0,
            "gat": 4,
            "ali": 2,
            "gra": 2,
            "gir": 6,
            "dra": 0,
            "con": 3,
            "guf": 4,
            "sep": 0
        }
    },
    "redez": {
        "totPoints": 1205555,
        "totn25": 17,
        "totWins": 5,
        "FEstats": {
            "n5000": 2,
            "n25000": 0,
            "n50000": 3,
            "totPointsMade": 160000
        },
        "charStats": {
            "cas": 2,
            "uni": 2,
            "zuc": 4,
            "gat": 2,
            "ali": 3,
            "gra": 1,
            "gir": 0,
            "dra": 4,
            "con": 5,
            "guf": 4,
            "sep": 0
        }
    }
}

//generic
var usersPfps = [];

//together
var srPointsData = [];
var sr25Data = [];

//synergo
var sPointsData = [];
var s25Data = [];
var sCharData = [];
var sFEData = [];

//redez
var rPointsData = [];
var r25Data = [];
var rCharData = [];
var rFEData = [];

"use strict"

async function getStatsData() {
    var res = await fetch("/games/search?limit=-1")
    var resp = await res.json();
    return resp
}

async function getGameData() {
    var res = await fetch("/stats/all")
    var resp = await res.json();
    return resp
}

async function getUsersPfp() {
    url = "/users/pfp/";
    for (let i = 0; i < stats.generic.collaborators.length; i++) {
        if (i == 0) {
            url += stats.generic.collaborators;
        } else {
            url += ";" + stats.generic.collaborators;
        }
    }
    var res = await fetch(url);
    usersPfps = await res.json();
}

function secondToHHMMSS(sec_num) {
    var days = Math.floor(sec_num / 86400);
    var hours = Math.floor((sec_num - (days * 86400)) / 3600);
    var minutes = Math.floor((sec_num - (days * 86400) - (hours * 3600)) / 60);
    var seconds = sec_num - (days * 86400) - (hours * 3600) - (minutes * 60);
    return { "days": days, "hours": hours, "minutes": minutes, "seconds": seconds }
}

function genChartsData() {
    gameData.forEach((game) => {
        let date = new Date(game.videoData.uploadDate);
        let sPoint = game.stats.synergo.overall.tPoints;
        let s25 = game.stats.synergo.overall.t25;
        let rPoint = game.stats.redez.overall.tPoints;
        let r25 = game.stats.redez.overall.t25;
        let annotation = "<a href='/" + game.videoData.id + "'>" + game.videoData.title + "</a>"
        srPointsData.push([new Date(date), sPoint, rPoint, annotation])
        sPointsData.push([new Date(date), sPoint, annotation])
        rPointsData.push([new Date(date), rPoint, annotation])

        sr25Data.push([new Date(date), s25, r25, annotation]);
        s25Data.push([new Date(date), s25, annotation]);
        r25Data.push([new Date(date), r25, annotation]);
    })

    var chars = ["castoro", "unicorno", "zucca", "gatto", "alieno", "granchio", "girasole", "drago", "coniglio", "gufo", "seppia"]
    var fe = ["numero di 5000", "numero di 25000", "numero di 50000"]

    sCharData = Object.entries(stats.synergo.charStats);
    sCharData = replaceNames(sCharData, chars)
    sCharData.unshift(["Tipo di Festa Estrema", "valore"]);


    rCharData = Object.entries(stats.redez.charStats);
    rCharData = replaceNames(rCharData, chars)
    rCharData.unshift(["Tipo di Festa Estrema", "valore"]);


    sFEData = Object.entries(stats.synergo.FEstats);
    sFEData.pop();
    sFEData = replaceNames(sFEData, fe)
    sFEData.unshift(["Tipo di Festa Estrema", "valore"]);


    rFEData = Object.entries(stats.redez.FEstats);
    rFEData.pop();
    rFEData = replaceNames(rFEData, fe)
    rFEData.unshift(["Tipo di Festa Estrema", "valore"]);
}

function replaceNames(toReplace, fromReplace) {
    for (let i = 0; i < toReplace.length; i++) {
        toReplace[i][0] = fromReplace[i]
    }
    return toReplace
}

function initDataInHtml() {
    gameData = await getGameData();
    //data synergo main section
    document.getElementById("spoint").innerHTML = stats.synergo.totPoints;
    document.getElementById("s25").innerHTML = stats.synergo.totn25;
    document.getElementById("sFE").innerHTML = stats.synergo.FEstats.totPointsMade;
    document.getElementById("sWins").innerHTML = stats.synergo.totWins;

    //data redez main section
    document.getElementById("rpoint").innerHTML = stats.redez.totPoints;
    document.getElementById("r25").innerHTML = stats.redez.totn25;
    document.getElementById("rFE").innerHTML = stats.redez.FEstats.totPointsMade;
    document.getElementById("rWins").innerHTML = stats.redez.totWins;

    //give crown based on wins 
    if (stats.redez.totWins < stats.synergo.totWins) {
        document.getElementById("scrown").style.display = "block";
    } else if (stats.redez.totWins > stats.synergo.totWins) {
        document.getElementById("rcrown").style.display = "block";
    } else {
        document.getElementById("scrown").style.display = "block";
        document.getElementById("rcrown").style.display = "block";
    }
    runAnimations(".countup")
    genChartsData()
    drawAllCharts()
}

$(window).resize(drawAllCharts);

function drawAllCharts() {
    drawAnnotationChart("s", sPointsData, "Points");
    drawAnnotationChart("r", rPointsData, "Points");
    drawAnnotationChart("sr", srPointsData, "Points");
    drawAnnotationChart("s", s25Data, "25");
    drawAnnotationChart("r", r25Data, "25");
    drawAnnotationChart("sr", sr25Data, "25");
    drawPieChart("s", sCharData, "Char");
    drawPieChart("r", rCharData, "Char");
    drawPieChart("s", sFEData, "FE");
    drawPieChart("r", rFEData, "FE");
}

function toggleChart(isAnnotation, chart) {
    let sep = $("#sep" + chart + "Charts");
    let sin = $("#sin" + chart + "Charts");
    let button = $("#" + chart + "Button");
    if (sep.css('display') == 'flex') {
        sep.hide();
        sin.show();
        button.text("Dividi i grafici")
        if (isAnnotation) {
            drawAnnotationChart("sr", window["sr" + chart + "Data"], chart)
        } else {
            drawPieChart("sr", "", chart)
        }
    } else {
        sep.show();
        sin.hide();
        button.text("Unisci")
        if (isAnnotation) {
            drawAnnotationChart("s", window["s" + chart + "Data"], chart)
            drawAnnotationChart("r", window["r" + chart + "Data"], chart)
        } else {
            drawPieChart("s", window["s" + chart + "Data"], chart)
            drawPieChart("r", window["r" + chart + "Data"], chart)
        }

    }
}