phdremote
=========

Remote control for OpenPHD Guiding. A plate solver script can be used to solve and display an overlay for the current image.


Building
========

  * Install go http://golang.org
  * Install libcfitsio http://heasarc.gsfc.nasa.gov/fitsio/
  * Install OpenPHD http://openphdguiding.org (snapshot 1165 or later)
  
  * go build src/phdremote.go
  
  * Modify the zsh script plotfield according to your astrometry installation
  
  * To build on Windows: open MINGW32 window with gcc configured, export CPATH and LIBRARY_PATH to cfitsio distribution directory, export GOPATH to phdremote directory

Running
=======

  * Start PHD and connect to your autoguider
  * ./phdremote -solver plotfield
  * Connect to http://localhost:8080/phdremote/ or use the IP address for remote viewing

Environment
===========

Depending on your operating system and your go workspace, you may need to configure the environment for the go and c compilers, for instance:

    export C_INCLUDE_PATH=/usr/local/include
    export LIBRARY_PATH=/usr/local/lib/     
    export PATH=$PATH:/usr/local/go/bin/    
    export GOPATH=`pwd`                     


