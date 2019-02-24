"use strict";
const cheerio = require("cheerio");

const analyseDraft = (url, html) => {
    const $ = cheerio.load(html);
    let draft = [];
    let headers = [];

    /**
     * Generate headers
     * [ 'Nr',
     * 'Lodtr.',
     * 'Universitet',
     * 'Tid',
     * 'Type',
     * 'Status(forløb)' ]
     */
    $("tbody > tr")
        .first()
        .children()
        .each((i, elem) => {
            let header = $(elem)
                .text()
                .trim();
            // Omdøber sidste header fra Status --> Valgt_id
            if (i === 5) {
                header = "Valgt_id";
            }
            headers.push(header);
        });
    $("body > table > tbody > tr + tr").each((i, elem) => {
        let draftee = {};

        $(elem)
            .children("td")
            .each((i, elem) => {
                let content = $(elem)
                    .text()
                    .trim();
                // Trækker id'et ud af status
                //draftee[headers[i]]
                if (i === 5) {
                    content = content.match(/([0-9]+)/);
                }

                draftee[headers[i]] = Array.isArray(content)
                    ? content[0]
                    : content;
            });

        draft.push(draftee);
    });
    // Slet sidste
    draft.splice(-1, 1);
    return draft;
};

const analyseLatestRound = (url, html) => {
    const $ = cheerio.load(html);
    let picks = [];
    let headers = [];
    /**
     * Generate headers
     * ['id',
     * 'Valgt',
     * 'Lodtr.',
     * 'Region',
     * 'Startdato',
     * 'Uddannelsessted',
     * 'Afdeling',
     * 'Speciale',
     * 'Uddannelsessted2',
     * 'Afdeling2',
     * 'Speciale2',
     * 'Runde' ]
     */
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

    $("table#tlTable > tbody > tr + tr").each((i, elem) => {
        let pick = {};
        pick.url = url;

        $(elem)
            .children("td")
            .each((i, elem) => {
                // Valgt dato (i=0) og som nummer (i=1)
                if (i < 2) {
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

                            // Region
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
    analyseDraft: analyseDraft,
    analyseLatestRound: analyseLatestRound
};
