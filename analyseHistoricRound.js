"use strict";
const cheerio = require("cheerio");
const analyseHistoricRound = (url, html) => {
    const $ = cheerio.load(html);
    let picks = [];
    let headers = [];
    $("tbody > tr")
        .first()
        .children()
        .each((i, elem) => {
            headers.push(
                $(elem)
                    .text()
                    .trim()
            );
        });

    headers.push("Uddannelsessted2", "Afdeling2", "Speciale2", "Runde");
    /**
     * Headers
     [ 'Valgt',
       'Lodtr.',
       'Region',
       'Startdato',
       'Uddannelsessted',
       'Afdeling',
       'Speciale',
       'Uddannelsessted2',
       'Afdeling2',
       'Speciale2',
       'Runde' ]
     */
    $("table#tlTable > tbody > tr + tr").each((i, elem) => {
        let pick = {};
        pick.url = url;

        $(elem)
            .children("td")
            .each((i, elem) => {
                // Valgt dato (i=0) og som nummer (i=1)
                if (i < 2) {
                    if (url.indexOf("Timelines") === -1) {
                        pick[headers[i]] = $(elem)
                            .text()
                            .replace(/\s+$/, "") // fjern trailing whitespace
                            .replace(/^\s/, ""); // fjern første ledende space
                    } else {
                        if (i === 0) {
                            pick.id = $(elem)
                                .text()
                                .match(/([0-9]+)/)[0];
                        }
                        if (i === 1) {
                            let text = $(elem).text();

                            let nr = text.match(/^([0-9]+)/g)[0];
                            let valgt = text.substring(nr.length);
                            pick["Valgt"] = valgt;
                            pick["Lodtr."] = nr + " ?";
                        }
                    }
                }
                // Selve valget
                if (i == 2) {
                    $(elem)
                        .find("tr td")
                        .each((n, elem) => {
                            /**
                             * 0: Region
                             * 1: Startdato
                             * 2: Uddannelsessted
                             * 3: Afdeling
                             * 4: Speciale
                             * 5: blank
                             * 6: blank
                             * 7: Afdeling 2
                             * 8: Speciale 2
                             */

                            if (n == 0) {
                                if ($(elem.text))
                                    // hvis ikke tom
                                    pick[headers[n + 2]] = $(elem)
                                        .text()
                                        .trim();
                            }

                            // Startdato, Udd.sted, afdeling, speciale
                            if (n >= 1 && n < 5) {
                                pick[headers[n + 2]] = $(elem)
                                    .text()
                                    .trim();
                            }

                            if (n >= 7) {
                                pick[headers[n]] = $(elem)
                                    .text()
                                    .trim();
                            }
                        });
                }
            });

        picks.push(pick);
    });
    // Slet sidste
    picks.splice(-1, 1);
    return picks;
};

module.exports = {
    analyseHistoricRound: analyseHistoricRound
};
