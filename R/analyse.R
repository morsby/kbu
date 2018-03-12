library(jsonlite)
library(dplyr)
library(tidyr)
library(ggplot2)
library(ggmap)
library(readtext)

setwd("~/Applications/Git/KBU-Stats/R")

# Read data
data <- fromJSON("../kbu.json")
# Split columns
data <- data %>% separate(Lodtr., c("Nummer", "Uni"), "[\\s]")
data$Nummer <- as.numeric(data$Nummer)

# Fjern uvalgte forløb
data <- data[-(which((data$Nummer == 0 || is.na(data$Nummer)) & is.na(data$Uni))), ]
# Sæt forhåndsvalgte forløb til valgt som nr. 0
data[is.na(data$Nummer), "Nummer"] <- 0

# Alt med google kan kræve et par forsøg (dvs. geocode og get_map funktioner)

# Skaf DK-kort
dk <- geocode("Denmark", output="more")
map <- get_map(c(
    dk[1,"west"],
    dk[1, "south"],
    dk[1, "east"],
    dk[1, "north"]))

register_google(key=readtext("googlekey.txt")$text)

# Steder
# Find unikke steder fra data
steder <- unique(data$Uddannelsessted)

# For hvert sted; find geokoden
geocodes <- geocode(steder)

# Tilføj stednavne til geocodes
geocodes$Uddannelsessted <- steder

# Kombiner data med geocodes ud fra stednavn
combined <- merge(data, geocodes)

saveData <- function(map, data) {
    save(map, "dk-map.RData")
    write.csv(data, "data.csv", row.names = F)
}