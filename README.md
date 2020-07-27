# Sales Discovery (SD) Api

###Project goal: 
- REST Api for Go backend of Sales Discovery (SD) tool

###SD Architecture:
- Front-end (in progress): Typescript + React.js -- https://github.com/Lemmart/sd
- Back-end: Golang with REST Api to scrape user-provided list of company websites
- Containerization (future): Docker-based in Kubernetes cluster managed by Helm

###Endpoint(s): 
- `GET` `/salesData`: returns json-formatted list of requested companies and their available offers and offer codes 