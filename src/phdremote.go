// Copyright 2014 Ted Goddard. All rights reserved.
// Use of this source code is governed the Apache 2.0
// license that can be found in the LICENSE file.

package main

import "fmt"
import "net"
import "os"
import "strings"
import "bufio"
import "net/http"
import "websocket"
import "log"
import "fits"
import "encoding/json"

func main() {
	fmt.Println("Starting server")

    previousImagePath := ""
    currentImagePath := ""

    wsClientHTML := `
<html>
    <head>
        <meta name="viewport" content="initial-scale=0.5, width=640, user-scalable=no">
        <meta name="apple-mobile-web-app-status-bar-style" content="black">
        <meta name="apple-mobile-web-app-capable" content="yes">
        <style>
          body {
            background-color: #202020;
          }
          .imgBox {
            position: relative;+
            left: 0;
            top: 0;
          }
          .brcontrols {
            position:fixed;
            bottom:10px;
            right:10px;
          }
          .brcontrols a {
            display:block;
            padding:10px;
            margin:10px;
            font-size:20px;
            border-radius:8px;
            background:red;
          }
          @media (max-width: 640px) {
              .bcontrols {
                position:fixed;
                bottom:100px;
                left:60px;
              }
              .bcinner {
              }
              .rcontrols {
                position:fixed;
                top:100px;
                right:10px;
              }
              .rcinner {
              }
              .bcontrols a {
                height:40px;
                padding:10px;
                margin:10px;
                font-size:40px;
                border-radius:8px;
                background:red;
              }
              .rcontrols a {
                display:block;
                padding:10px;
                margin:10px;
                font-size:40px;
                border-radius:8px;
                background:red;
              }
              .brcontrols {
                position:fixed;
                bottom:100px;
                right:10px;
              }
          }
          @media (min-width: 641px) {
              .bcontrols {
                position:fixed;
                bottom:20px;
                left:50%%;
              }
              .bcinner {
                margin-left:-50%%;
              }
              .rcontrols {
                position:fixed;
                top:50%%;
                right:0px;
              }
              .rcinner {
                margin-top: -50%%;
              }
              .bcontrols a {
                height:40px;
                padding:10px;
                margin:20px;
                font-size:20px;
                border-radius:8px;
                background:red;
              }
              .rcontrols a {
                display:block;
                padding:10px;
                margin:10px;
                font-size:20px;
                border-radius:8px;
                background:red;
              }
          }
        </style>
        <script>
        var ws = new WebSocket("ws://" + location.host + "/echo/");
        ws.onmessage = function(msg) {console.log(msg.data);
            var msgJSON = JSON.parse(msg.data);
            console.log(msgJSON.Event);
            var marker = document.getElementById("marker");
            if ("LoopingExposures" == msgJSON.Event)  {
                updateCam();
            };
            if ("StartCalibration" == msgJSON.Event)  {
               showMarker("calib");
            };
            if ("GuideStep" == msgJSON.Event)  {
               updateCam();
               showMarker("guide");
            };
            if ("StarLost" == msgJSON.Event)  {
               showMarker("lost");
            };
        };

        function updateCam() {
            var camImg = document.getElementById("cam");
            camImg.src = "cam.jpg?" + new Date().getTime();
        }
        function showMarker(name) {
            clearMarkers();
            document.getElementById("m-" + name).style["opacity"] = 1.0;
        }
        function clearMarkers() {
            var marker = document.getElementById("marker");
            for (i = 0; i < marker.childNodes.length; i++)  {
               if (!marker.childNodes[i].style) { continue; };
               marker.childNodes[i].style["opacity"] = 0;
            }
        }
        function getClickPosition(e) {
            var parentPosition = getPosition(e.currentTarget);
            return {
                x: e.clientX - parentPosition.x,
                y: e.clientY - parentPosition.y
            }
        }
        function getPosition(element) {
            var x = 0;
            var y = 0;
            while (element) {
                x += (element.offsetLeft - element.scrollLeft +
                    element.clientLeft);
                y += (element.offsetTop - element.scrollTop +
                    element.clientTop);
                element = element.offsetParent;
            }
            return { x: x, y: y };
        }

        function imageClick(event) {
            var imgClick = getClickPosition(event);
            ws.send(JSON.stringify({method: "set_lock_position",
                params: [imgClick.x, imgClick.y], id: 42}));
            var marker = document.getElementById("marker");
            marker.style.top = imgClick.y - 10;
            marker.style.left = imgClick.x - 10;
            showMarker("select");
        };
        function guide() {
            console.log("guide");
            ws.send(JSON.stringify({method:"guide",
                params:[{pixels:1.5, time:8, timeout:40}, false], id:1}));
        };
        function stop() {
            console.log("stop");
            ws.send(JSON.stringify({"method":"set_paused","params":[true,"full"],"id":2}));
        };
        function loop() {
            console.log("loop");
            ws.send(JSON.stringify({method:"loop", id:3}));
        };
        function expose(t) {
            console.log("expose" + t);
            ws.send(JSON.stringify({method:"set_exposure", params:[t], id:4}));
        };
        function toggleBullseye() {
            var bullseyeElement = document.getElementById("bull");
            bullseyeElement.style["opacity"] = 1.0 - bullseyeElement.style["opacity"];
        }
        function adjustSizes() {
            var bullseyeElement = document.getElementById("bull");
            var camElement = document.getElementById("cam");
            bullseyeElement.style.width = camElement.width;
            bullseyeElement.style.height = camElement.height;
        }
        window.onresize = function(event)  {
            adjustSizes();
        }
        </script>
    </head>
    <body>
    <div class="imgBox">
        <img id="cam" src="cam.jpg" onclick="imageClick(event)" onload="adjustSizes()"
            style="-webkit-filter:brightness(140%%)contrast(300%%);position: relative; top: 0; left: 0;">
        <svg id="bull" width="100%%" height="100%%" style="opacity:0; position: absolute; top: 0; left: 0;">
            <g >
                <line x1="0px" y1="50%%" x2="100%%" y2="50%%" stroke="red" stroke-width="1" />
                <line x1="50%%" y1="0px" x2="50%%" y2="100%%" stroke="red" stroke-width="1" />
                <circle cx="50%%" cy="50%%" r="10%%" stroke="red" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="4%%" stroke="red" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="2%%" stroke="red" stroke-width="1" fill="none" />
            </g>
        </svg>
        <svg id="marker" width="20" height="20" style="position: absolute; top: 0; left: 0;">
            <g id="m-select" style="opacity:0">
                <rect x="-4" y="-4" width="10" height="10" stroke="white" stroke-width="2" fill="none" />
                <rect x="14" y="-4" width="10" height="10" stroke="white" stroke-width="2" fill="none" />
                <rect x="-4" y="14" width="10" height="10" stroke="white" stroke-width="2" fill="none" />
                <rect x="14" y="14" width="10" height="10" stroke="white" stroke-width="2" fill="none" />
            </g>
            <g id="m-calib" style="opacity:0">
                <rect x="0" y="0" width="20" height="20" stroke="yellow" stroke-width="4" stroke-dasharray="2 2" fill="none" />
            </g>
            <g id="m-guide" style="opacity:0">
                <line x1="10" y1="0" x2="10" y2="20" stroke="green" stroke-width="2" />
                <line x1="0" y1="10" x2="20" y2="10" stroke="green" stroke-width="2" />
                <rect x="4" y="4" width="12" height="12" stroke="green" stroke-width="2" fill="none" />
            </g>
            <g id="m-lost"  style="opacity:0">
                <line x1="0" y1="0" x2="20" y2="20" stroke="red" stroke-width="2" />
                <line x1="20" y1="0" x2="0" y2="20" stroke="red" stroke-width="2" />
                <rect x="0" y="0" width="20" height="20" stroke="red" stroke-width="4" fill="none" />
            </g>
        </svg>
    </div>
    <div class="rcontrols" >
      <div class="rcinner" >
        <a onclick="expose(500)">0.5s</a>
        <a onclick="expose(1000)">1.0s</a>
        <a onclick="expose(2000)">2.0s</a>
      </div>
    </div>
    <div class="bcontrols" >
      <div class="bcinner" >
        <a onclick="guide()">GUIDE</a>
        <a onclick="stop()">STOP</a>
        <a onclick="loop()">LOOP</a>
      </div>
    </div>
    <div class="brcontrols" >
      <div class="brinner" >
        <a onclick="toggleBullseye()">
            <svg width="40px" height="40px">
            <g >
                <line x1="0px" y1="50%%" x2="100%%" y2="50%%" stroke="black" stroke-width="1" />
                <line x1="50%%" y1="0px" x2="50%%" y2="100%%" stroke="black" stroke-width="1" />
                <circle cx="50%%" cy="50%%" r="20%%" stroke="black" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="10%%" stroke="black" stroke-width="1" fill="none" />
            </g>

            </svg></a>
      </div>
    </div>
    </body>
</html>
`
    http.HandleFunc("/phdremote/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, wsClientHTML)
    })

    //only one websocket allowed for now
    var wsConn *websocket.Conn
    var phdWrite *bufio.Writer
    phdDone := make(chan bool)

    GuideWatch := func (conn *net.Conn)  {
        connRead := bufio.NewReader(*conn)
        status := ""
        var err error
        for (err == nil)  {
            status, err = connRead.ReadString('\n')
            log.Print(status)
            var phdMessage map[string]interface{}
            err = json.Unmarshal([]byte(status), &phdMessage)
            if (nil != err)  {
                log.Print("jsonrpc ERROR, applying remove backslash hack ", err)
                status = strings.Replace(status, "\\", "\\\\", -1)
                err = json.Unmarshal([]byte(status), &phdMessage)
            }
            if (nil == err)  {
                if (nil != phdMessage["jsonrpc"]) {
                    log.Print("jsonrpc contents", status)
                    switch result := phdMessage["result"].(type)  {
                        case map[string]interface{}:
                            previousImagePath = currentImagePath
                            currentImagePath = result["filename"].(string)
                            if ("" != previousImagePath)  {
                                os.Remove(previousImagePath)
                            }
                        case float64:
                            log.Print("float64 jsonrpc result")
                    }
                }
                fmt.Println(phdMessage["jsonrpc"])
            }

            if (nil != wsConn)  {
                log.Print("writing to WebSocket")
                (*wsConn).Write([]byte(status))
            }
        }
        phdDone <- true
    }

    SocketWatch := func()  {
        var err error
        for (err == nil)  {
            var msg = make([]byte, 512)
            var n int
            n, err = wsConn.Read(msg)
            fmt.Printf("WEBSOCKET Received: %s.\n", msg[:n])
            if (nil != phdWrite)  {
                fmt.Fprintf(phdWrite, string(msg))
                phdWrite.Flush()
            }
        }
    }

    EchoServer := func(newConn *websocket.Conn) {
        log.Print("EchoServer started")
        wsConn = newConn
        go SocketWatch ()
        echoDone := <-phdDone
        if (echoDone)  {
            log.Print("EchoServer done")
       }
    }

    log.Print("websocket.Handler")
    wsHandler := websocket.Handler(EchoServer)
	http.Handle("/echo/", wsHandler)

    conn, err := net.Dial("tcp", "localhost:4400")
    if (err == nil) {
        phdWrite = bufio.NewWriter(conn)
        go GuideWatch (&conn)
    } else {
        log.Print("Unable to connect to PHD")
    }

    http.HandleFunc("/phdremote/cam.png", func(w http.ResponseWriter, r *http.Request) {
log.Print("returning png image")
        if (nil != phdWrite)  {
            fmt.Fprintf(phdWrite, "{\"method\":\"save_image\",\"id\":123}\n")
            phdWrite.Flush()
        }
        w.Header().Set("Content-Type", "image/png")
        momentaryImagePath := currentImagePath
        if ("" == momentaryImagePath)  {
            momentaryImagePath = "RCA.fit"
        }
        fits.ConvertPNG(momentaryImagePath, w)
    })

    http.HandleFunc("/phdremote/cam.jpg", func(w http.ResponseWriter, r *http.Request) {
log.Print("returning jpg image")
        if (nil != phdWrite)  {
            fmt.Fprintf(phdWrite, "{\"method\":\"save_image\",\"id\":123}\n")
            phdWrite.Flush()
        }
        w.Header().Set("Content-Type", "image/jpeg")
        momentaryImagePath := currentImagePath
        if ("" == momentaryImagePath)  {
            momentaryImagePath = "RCA.fit"
        }
        fits.ConvertJPG(momentaryImagePath, w)
    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8080", nil))

}


