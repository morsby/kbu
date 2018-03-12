# Analyse af valg af KBU-forløb fra 2010 til 2017.

Dette repo indeholder koden, der ligger bag Shiny-app'en, der er tilgængelig på 
https://morsby.shinyapps.io/kbu-stats/.

Koden er tredelt.

1. Der er en NodeJS-app (`index.js`), der henter data fra http://basislaege.dk (under `Historik`) og eksporterer det til JSON.
2. Et R-script (`R/analyse.R`), der læser JSON-filen dannet ovenfor og henter koordinater fra Google Maps samt behandler data (se filen for detaljer).
3. En Shiny App (`R/app.R`), der tillader brugeren at vælge et antal byer og herefter kan se, hvordan denne by er blevet valgt i datasættet.

Al feedback er velkommen.
