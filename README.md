phdremote
=========

Remote control for OpenPHD Guiding 


Building
========

  * Install go http://golang.org
  * Install libcfitsio http://heasarc.gsfc.nasa.gov/fitsio/
  * Install OpenPHD http://openphdguiding.org (snapshot 1165 or later)
  
  * go build src/phdremote.go

Running
=======

  * Start PHD and connect to your autoguider
  * ./phdremote
  * Connect to http://localhost:8080/phdremote/ or use the IP address for remote viewing
