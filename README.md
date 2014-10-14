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

Environment
===========

Depending on your operating system and your go workspace, you may need to configure the environment for the go and c compilers, for instance:

    export C_INCLUDE_PATH=/usr/local/include
    export LIBRARY_PATH=/usr/local/lib/     
    export PATH=$PATH:/usr/local/go/bin/    
    export GOPATH=`pwd`                     


