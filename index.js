//const HtmlTableToJson = require('html-table-to-json');
//const DOMParser = require('xmldom').DOMParser;

const request = require("async-request");
const fs = require("fs");
const _ = require("lodash");
//const curl = require('curl');

const analyseHistoricRound = require("./analyseHistoricRound")
    .analyseHistoricRound;

const analyseDraft = require("./analyseLatestRound").analyseDraft;
const analyseLatestRound = require("./analyseLatestRound").analyseLatestRound;

const kbus = [
    "https://kbu.logbog.net/Ajax_get2010v2.asp",
    "https://kbu.logbog.net/Ajax_get2011v1.asp",
    "https://kbu.logbog.net/Ajax_get2011v2.asp",
    "https://kbu.logbog.net/Ajax_get2012v1.asp",
    "https://kbu.logbog.net/Ajax_get2012v2.asp",
    "https://kbu.logbog.net/Ajax_get2013v1.asp",
    "https://kbu.logbog.net/Ajax_get2013v2.asp",
    "https://kbu.logbog.net/Ajax_get2014v1.asp",
    "https://kbu.logbog.net/Ajax_get2014v2.asp",
    "https://kbu.logbog.net/Ajax_get2015v1.asp",
    "https://kbu.logbog.net/Ajax_get2015v2.asp",
    "https://kbu.logbog.net/Ajax_get2016v1.asp",
    "https://kbu.logbog.net/Ajax_get2016v2.asp",
    "https://kbu.logbog.net/Ajax_get2017v1.asp",
    "https://kbu.logbog.net/Ajax_get2017v2.asp",
    "https://kbu.logbog.net/Ajax_get2018v1.asp",
    "https://kbu.logbog.net/Ajax_get2018v2.asp",
    "https://kbu.logbog.net/Ajax_get2019v1.asp",
    "https://kbu.logbog.net/Ajax_get2019v2.asp",
    "https://kbu.logbog.net/Ajax_get2020v1.asp",
];

const getPicks = kbus.map(async url => {
    let response = await request(url);
    let picks = analyseHistoricRound(url, response.body);
    return picks;
});

const getLatestRound = async url => {
    let draftHtml = await request("https://kbu.logbog.net/AJAX_Draft.asp"),
        roundHtml = await request(
            "https://kbu.logbog.net/AJAX_Timelines.asp"
        ),
        draft = analyseDraft(
            "https://kbu.logbog.net/AJAX_Draft.asp",
            draftHtml.body
        ),
        round = analyseLatestRound(
            "https://kbu.logbog.net/AJAX_Timelines.asp",
            roundHtml.body
        );

    round.map(choice => {
        let chooser = _.find(draft, { Valgt_id: choice.id }) || {},
            uni;

        switch (chooser["Universitet"]) {
            case "Aarhus Universitet":
                uni = "AU";
                break;
            case "Aalborg Universitet":
                uni = "AAU";
                break;
            case "Københavns Universitet":
                uni = "KU";
                break;
            case "Syddansk Universitet":
                uni = "SDU";
                break;
            default:
                uni = "";
        }
        choice["Lodtr."] = choice["Lodtr."].replace("?", uni);

        choice["Valgt"] = chooser["Tid"]
            ? chooser["Tid"].replace("kl. ", "")
            : "";
        return choice;
    });
    console.log("Draft: ", draft[0], "\nRound: ", round[0]);
    return round;
};

allPicks = [];
Promise.all(getPicks).then(async completed => {
    allPicks = [];
    for (i = 0; i < completed.length; i++) {
        for (n = 0; n < completed[i].length; n++) {
            allPicks.push(completed[i][n]);
        }
    }
    let latestRound = await getLatestRound();
    latestRound.map(el => newPicks.push(el));

    fs.writeFile("data.json", JSON.stringify(allPicks), () => {
        console.log("Wrote file");
    });


    //console.log(completed.flatten());
});
