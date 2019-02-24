library(jsonlite)
library(dplyr)
library(tidyr)
library(ggplot2)
library(googleway)
library(ggmap)
library(readtext)
library(stringr)
setwd("~/Applications/Git/KBU-Stats/R")

# Read data
data <- fromJSON("../data.json")
# Split columns
data <- data %>% separate(Lodtr., c("Nummer", "Uni"), "[\\s]")
data$Nummer <- as.numeric(data$Nummer)

# Fjern uvalgte forløb
data <-
  data[-(which((data$Nummer == 0 ||
                  is.na(data$Nummer)) & is.na(data$Uni))),]
# Sæt forhåndsvalgte forløb til valgt som nr. 0
data[is.na(data$Nummer), "Nummer"] <- 0

# Beregn runden ud fra startdato
findRunde <- function(startdato) {
  efteraar <- c("jul", "aug", "sep", "okt", "nov")
  foraar <- c("feb", "mar", "apr", "maj")
  
  split <- strsplit(startdato, " ")[[1]]
  maaned <- split[2]
  aar <- split[3]
  if (maaned %in% efteraar) {
    aarstid <- "Efterår"
  } else {
    aarstid <- "Forår"
  }
  return(paste(aarstid, aar))
}

findRundestart <- function(runde) {
  aar <- str_extract(runde,"[0-9]+")
  if(grepl("Efterår", runde)) {
    return(paste0("20",aar,"-08-01"))
  } else {
    return(paste0("20",aar,"-02-01"))
  }
}

data$Runde <- mapply(findRunde, data$Startdato)
data$Rundestart <- mapply(findRundestart, data$Runde)
data$Rundestart <- as.Date(data$Rundestart)

data$Startdato <- as.Date(data$Startdato, "%d. %b %y")

# Alt med google kan kræve et par forsøg (dvs. geocode og get_map funktioner)
# Skaf DK-kort
#dk <- geocode("Denmark", output="more")
#map <- get_map(c(
#    dk[1,"west"],
#    dk[1, "south"],
#    dk[1, "east"],
#    dk[1, "north"]))

register_google(key = readtext("googlekey.txt")$text)

# Steder
# Find unikke steder fra data
steder <- unique(data$Uddannelsessted)

# For hvert sted; find geokoden
#geocodes <- geocode(steder, output="latlon", source="google")

# Tilføj stednavne til geocodes
#geocodes$Uddannelsessted <- steder

# Kombiner data med geocodes ud fra stednavn
combined <- merge(data, geocodes)

saveData <- function(map = NULL, data = NULL) {
  if (!is.null(map)) {
    save(map, file = "dk-map.RData")
  }
  if (!is.null(data)) {
    write.csv(data, "data.csv", row.names = F)
  }
}
