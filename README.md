# LED Website Monitoring 
![header](https://alonsoarteaga.me/content/images/2019/11/example.gif "led monitoring")

Monitor the status of websites with a raspberry pi and with some LEDs, 
change the color depending on the website status.

Read the Build process @ https://alonsoartega.me

## Initial Setup
change PinsIndex in cmd/flowingLED/main.go to match your GPIO setup.
run `go build`
`$ ./flowingLED &`

## todos
*[ ] create systemctl service
*[ ] swaggerAPI documentation
*[ ] fix while loop for healthcheck
*[ ] find way to make website monitoring dynamic
*[ ] web gui management site
*[ ] fix healthchecks
*[ ] feature: fix hours of operation (i.e lights off between midnight but turn red in case of downtime)