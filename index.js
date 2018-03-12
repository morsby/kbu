//const HtmlTableToJson = require('html-table-to-json');
//const DOMParser = require('xmldom').DOMParser;
const cheerio = require('cheerio');
const request = require('async-request');
const fs = require('fs');
//const curl = require('curl');

const kbus = [
	'https://www.basislaege.dk/Ajax_get2010v2.asp',
	'https://www.basislaege.dk/Ajax_get2011v1.asp',
	'https://www.basislaege.dk/Ajax_get2011v2.asp',
	'https://www.basislaege.dk/Ajax_get2012v1.asp',
	'https://www.basislaege.dk/Ajax_get2012v2.asp',
	'https://www.basislaege.dk/Ajax_get2013v1.asp',
	'https://www.basislaege.dk/Ajax_get2013v2.asp',
	'https://www.basislaege.dk/Ajax_get2014v1.asp',
	'https://www.basislaege.dk/Ajax_get2014v2.asp',
	'https://www.basislaege.dk/Ajax_get2015v1.asp',
	'https://www.basislaege.dk/Ajax_get2015v2.asp',
	'https://www.basislaege.dk/Ajax_get2016v1.asp',
	'https://www.basislaege.dk/Ajax_get2016v2.asp',
	'https://www.basislaege.dk/Ajax_get2017v1.asp',
	'https://www.basislaege.dk/Ajax_get2017v2.asp'
];

const getPicks = kbus.map(async url => {
	let response = await request(url);
	let picks = analyseTable(url, response.body);
	return picks;
});

const analyseTable = (url, html) => {
	const $ = cheerio.load(html);
	let picks = [];
	let headers = [];
	$('tbody > tr')
		.first()
		.children()
		.each((i, elem) => {
			headers.push(
				$(elem)
					.text()
					.trim()
			);
		});

	headers.push('Uddannelsessted2', 'Afdeling2', 'Speciale2', 'Runde');

	$('table#tlTable > tbody > tr + tr').each((i, elem) => {
		let pick = {};
		pick.runde = url;

		$(elem)
			.children('td')
			.each((i, elem) => {
				// Valgt dato og som nummer
				if (i < 2) {
					pick[headers[i]] = $(elem)
						.text()
						.replace(/\s+$/, '') // fjern trailing whitespace
						.replace(/^\s/, ''); // fjern første ledende space
				}
				// Selve valget
				if (i == 2) {
					$(elem)
						.find('tr td')
						.each((n, elem) => {
							// Region
							if (n == 0) {
								if ($(elem.text))
									// hvis ikke tom
									pick[headers[n + 2]] = $(elem)
										.text()
										.trim();
							}

							// Udd.sted, afdeling, speciale
							if (n >= 2 && n < 5) {
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
allPicks = [];
Promise.all(getPicks).then(completed => {
	allPicks = [];
	for (i = 0; i < completed.length; i++) {
		for (n = 0; n < completed[i].length; n++) {
			allPicks.push(completed[i][n]);
		}
	}

	console.log(allPicks[0]);

	fs.writeFile('data.json', JSON.stringify(allPicks), () => {
		console.log('Wrote file');
	});

	//console.log(completed.flatten());
});
