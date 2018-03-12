library(shiny)
library(dplyr)
library(ggplot2)
library(ggmap) # For maps
library(gmodels) # For CI calc
library(markdown) # For markdownToHTML function

# Forbered data

load("dk-map.RData")
data <- read.csv("data.csv")


# Pool byer
byer <- list(
    # Sjælland
    KBH=c("Amager", "Bispebjerg", "Frederiksberg", "Glostrup", "Gentofte", "Herlev", "Hvidovre", "Rigshospitalet"),
    Frederikssund = c("Frederikssund"),
    Helsingør = c("Helsingør"),
    Hillerød = c("Hillerød", "Nordsjællands Hospital"),
    Holbæk = c("Holbæk"),
    Køge = c("Køge", "Sjællands Universitetshospital"),
    Roskilde = c("Roskilde"),
    Slagelse = c("Slagelse"),
    Næstved = c("Næstved"),
    
    # Nordjylland
    Hjørring = c("Hjørring", "Vendsyssel"),
    Aalborg = c("Aalborg"),
    Thisted = c("Thisted"),
    
    # Midtjylland
    Aarhus = c("Aarhus", "Århus"),
    Holstebro = c("Holstebro"),
    Silkeborg = c("Silkeborg"),
    Viborg = c("Viborg"),
    Horsens = c("Horsens"),
    Herning = c("Herning"),
    Randers = c("Randers"),
    
    # Sønderjylland
    Esbjerg = c("Esbjerg"),
    Fredericia = c("Fredericia"),
    Kolding = c("Kolding"),
    Vejle = c("Vejle"),
    Haderslev = c("Haderslev"),
    Sønderborg = c("Sønderborg"),
    Aabenraa = c("Aabenraa"),
    
    # Øerne
    "Nykøbing Falster" = c("Nykøbing F."),
    Bornholm = c("Bornholm"),
    
    "Færøerne" = c("Færøerne"),
    
    # Fyn
    Odense = c("Odense"),
    Svendborg = c("Svendborg"))


for(i in 1:length(byer)) {
    matches <- grepl(paste(byer[[i]], collapse="|"), data$Uddannelsessted)
    data[matches, "Område"] <- names(byer[i])
}


data <- data %>% 
    group_by(runde) %>% 
    mutate(PlaceringIAar = 1-Nummer/max(Nummer)) %>% 
    group_by(Område) %>% 
    mutate(Popularitet=mean(PlaceringIAar))

#data %>% group_by(Uddannelsessted) %>% 
#    summarise(Gennemsnit = mean(Nummer), sd=sd(Nummer), n=n()) %>%
#    mutate(se=sd/n,
#           lower=Gennemsnit - qt(1 - (0.05 / 2), n - 1) * se,
#           upper = Gennemsnit + qt(1 - (0.05 / 2), n - 1) * se) 

server <- function(input, output) {
    
    # Map
    output$map <- renderPlot({
        ggmap(map) + 
            geom_point(data=subset(data, Område %in% input$byer), aes(x=lon, y=lat, color=Popularitet), size=3) +
            theme_void()
    })
    
    # Density
    densityPlot <-  reactive({
        ggplot(data=subset(data, Område %in% input$byer), aes(x=PlaceringIAar, color=Område, fill=Område)) + 
            geom_density(alpha=0.2) +
            ggtitle("Popularitet for valgte byer") + 
            labs(y="Popularitet (densitet)", x="Hvornår stedet er blevet valgt. 0=Sidst, 1=Først") +
            xlim(0,1)
        })
    
    # Nummeret:
    intercept <- reactive({1-as.numeric(input$nr)/as.numeric(input$nrTot)})
    output$nummer <- reactive({paste0("Dit nummer svarer til en værdi på ", round(intercept(),3))})
    
    # Regn areal under kurven (benyttes ikke p.t.)
    auc <- reactive({
   
        xy <- ggplot_build(densityPlot())
        
        xy <- xy$data[[1]]
        
        interval <- xy$x[2] - xy$x[1]
        
        xy$auc <- interval*xy$density
        
        xy <- xy %>% group_by(group) %>% mutate(By=input$byer[group])
        
        xy

        })
    
    # Density plot    
    output$density <- renderPlot({
        if(length(input$byer) > 0) {
            densityPlot() + geom_vline(xintercept=intercept())
        }
    })
    
    # Boxplot
    output$boxplot <- renderPlot({
        if(length(input$byer) > 0) {
            ggplot(data=subset(data, Område %in% input$byer), aes(x=Område, y=PlaceringIAar, fill=Område)) + 
                geom_boxplot(alpha=0.2) + 
                coord_flip() +
                ylim(0,1) +
                labs(y="Hvornår stedet er blevet valgt. 0=Sidst, 1=Først", x="") +
                geom_hline(yintercept=intercept())
        }
    })

    # Tabel over denstet - benyttes ikke
    output$dens <- renderTable({
        if(length(input$byer) > 0) {
            auc() %>% filter(x <= intercept()) %>% group_by(By) %>% summarise(Sum=sum(auc))
        }
    })
    
    # Boxplot for sidste plads i byen
    output$lastPlot <- renderPlot({
        if(length(input$byer) > 0) {
            mins <- data %>% filter(Område %in% input$byer) %>% 
                group_by(Område, runde) %>% 
                summarise(Min=min(PlaceringIAar))
            
            ggplot(data=mins, aes(x=Område, y=Min, fill=Område)) + 
                geom_boxplot(alpha=0.2) + 
                coord_flip() +
                labs(y="Hvornår sidste plads i byen blev valgt. 0=Sidst, 1=Først", x="") +
                ylim(0,1) +
                geom_hline(yintercept=intercept())
        }
    })
    
    # Tabel for sidste plads i byen
    output$lastTable <- renderTable({
        if(length(input$byer) > 0) {
             data %>% filter(Område %in% input$byer) %>% 
                group_by(Område, runde) %>% 
                summarise(Min=min(PlaceringIAar)) %>% 
                group_by(Område) %>% 
                summarise(
                    "Gennemsnitlig sidste position"=ci(Min)[1], 
                    "Nedre CI"=ci(Min)[2], 
                    "Øvre CI"=ci(Min)[3], 
                    SE=ci(Min)[4])
        }
    })
}

ui <- fluidPage(
    titlePanel("KBU-statistik 2010-2017"),
    
    HTML(markdownToHTML(fragment.only=TRUE, text=c(
"Denne lille interaktive hjemmeside indeholder statistik over KBU-fordelinger siden 2010. 
Data er hentet ud fra de tabeller, der ses på [basislaege.dk](http://basislaege.dk) under `Historik`.

Da der er et forskelligt antal numre hver lodtrækningsrunde, har jeg givet alle lodtrækningsnumre
en ny værdi ud fra formlen: `Placering = 1-nummer/(antal numre denne runde)`. Herved vil et 
tidligt valgt forløb have en værdi tæt på 1, mens et sent valgt forløb vil være tættere på 0.

Alle hospitaler er poolet i byer. Så alle hospitaler i København vil være at findes under KBH,
alle hospitaler i Aarhus vil ligeledes være under ét. 

P.t. er der ikke taget højde for specialer på nogen måde men udelukkende geografi.

I menuen til venstre bedes du indtaste dit (eller et fiktivt) nummer samt hvor mange numre,
der totalt er i denne runde. Dette vil tilføje nogle vertikale linjer på graferne, der viser
hvor det angivne nummer falder.
"
    ))),
    
    sidebarLayout(
        sidebarPanel(
            h4("Vælg dine input"),
            textInput("nr", "Dit nummer:", "1"),
            textInput("nrTot", "Hvor mange har trukket:", "643"),
            textOutput("nummer"),
            HTML("<br><br>"),
            checkboxGroupInput("byer", "Hvilke byer vil du se?", choices=sort(names(byer)))
        ),
        mainPanel(
                  h4("Placering på kort"),
                  plotOutput("map"), 
                  p("OBS: Lokaliteten er hentet via en Google ud fra hospitalets navn,
                    så det er ikke sikkert, denne er præcis."),
                  
                  HTML("<br><br>"),
                  h4("Hvornår vælges de valgte byer normalt?"),
                  plotOutput("density"), 
                  plotOutput("boxplot"),
                  div(
                      p("I graferne ovenfor vises populariteten af de valgte byer. 
                        Det øverste plot er et Kernel Density Plot. Der hvor Y-værdien er højest,
                        er stedet mest populært. Arealet under kurven (integralet) mellem to punkter
                        svarer til sandsynligheden for at stedet er valgt med et nummer i det interval"),
                      p("Samme data vises i boksplottet, der er mere intuitivt at få fornuft ud af:
                        Der vises 25%, 50% og 75% percentiler. Enderne af plottet er ved den mest ekstreme
                        værdi, der ikke er mere end 1.5 gange IQR.")
                  ),
                  
                  HTML("<br><br>"),
                  h4("Hvornår går den SIDSTE plads i byerne?"),
                  div(
                      p("Nedenstående data viser, hvornår den SIDSTE plads er gået for hver
                        by. I tabellen vises gennemsnittet for alle runder samt konfidensintervallet og SE hertil."),
                        
                      p("I boksplottet her vises fordelingen af det sidste nummer i hver runde, der har valgt en given by.")
                  ),
                  tableOutput("lastTable"), 
                  plotOutput("lastPlot")
          )
    )
)

shinyApp(ui = ui, server = server)