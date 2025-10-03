# gocket
A local Pocket clone

This project will
 - Run in Go
 - Run a local web site
 - Will accept a web page url and save a snapshot of it for reading later
 - It will use a SQLite DB to hold the saved articles
 - The web page will allow browsing of the saved articles
 - The ultimate design is that a single docker container or kubernetes pod can be used to host an instance, with the SQLite DB file as the cached articles. This would allow it to scale to many users if hosted

 