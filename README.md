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
  
  * Use Chrome on Windows 
 
Translation Clone
=============

The translation clone feature helps to position dim objects in the field center. When the plate solving view is displayed, tap on the desired object, then tap on a visible star. The desired object will be highlighted with a large red circle and the position of the bright star that brings the desired object to the centre field will be highlighted with a dashed circle. Adjust the telescope position until the star is centered in the small circle.

Environment
===========

Depending on your operating system and your go workspace, you may need to configure the environment for the go and c compilers, for instance:

    export C_INCLUDE_PATH=/usr/local/include
    export LIBRARY_PATH=/usr/local/lib/     
    export PATH=$PATH:/usr/local/go/bin/    
    export GOPATH=`pwd`                     


