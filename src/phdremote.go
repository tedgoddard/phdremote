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

    wsClientHTML :=
        "<html>" +
            "<head>" +
                "<style>" +
                "  .controls {" +
                "    position:fixed;"+
                "    bottom:0;left:0;"+
                "    left:0;"+
                "  }" +
                "  .controls button {" +
                "    height:40px;"+
                "  }" +
                "</style>" +
                "<script>" +
                "var ws = new WebSocket('ws://' + location.host + '/echo/');" +
                "ws.onmessage = function(msg) {console.log(msg.data);" +
                "    var msgJSON = JSON.parse(msg.data);" +
                "    console.log(msgJSON.Event);" +
                "    var marker = document.getElementById('marker');"+
                "    if ('LoopingExposures' == msgJSON.Event)  {" +
                "        updateCam();" +
                "    };" +
                "    if ('StartCalibration' == msgJSON.Event)  {" +
                "       showMarker('calib');" +
                "    };" +
                "    if ('GuideStep' == msgJSON.Event)  {" +
                "       updateCam();" +
                "       showMarker('guide');" +
                "    };" +
                "    if ('StarLost' == msgJSON.Event)  {" +
                "       showMarker('lost');" +
                "    };" +
                "};" +
    
                "function updateCam() {" +
                "    var camImg = document.getElementById('cam');" +
                "    camImg.src = 'cam.jpg?' + new Date().getTime();" +
                "}" +
                "function showMarker(name) {" +
                "    clearMarkers();" +
                "    document.getElementById('m-' + name).style['opacity'] = 1.0;" +
                "}" +
                "function clearMarkers() {" +
                "    var marker = document.getElementById('marker');"+
                "    for (i = 0; i < marker.childNodes.length; i++)  {" +
                "       if (!marker.childNodes[i].style) { continue; };" +
                "       marker.childNodes[i].style['opacity'] = 0;" +
                "    }" +
                "}" +
                "function getClickPosition(e) {" +
                "    var parentPosition = getPosition(e.currentTarget);" +
                "    return {" +
                "        x: e.clientX - parentPosition.x," +
                "        y: e.clientY - parentPosition.y" +
                "    }" +
                "}" +
                "function getPosition(element) {" +
                "    var x = 0;" +
                "    var y = 0;" +
                "    while (element) {" +
                "        x += (element.offsetLeft - element.scrollLeft +" +
                "            element.clientLeft);" +
                "        y += (element.offsetTop - element.scrollTop +" +
                "            element.clientTop);" +
                "        element = element.offsetParent;" +
                "    }" +
                "    return { x: x, y: y };" +
                "}" +

                "function imageClick(event) {" +
                "    var imgClick = getClickPosition(event);" +
                "    ws.send(JSON.stringify({method: 'set_lock_position'," +
                "        params: [imgClick.x, imgClick.y], id: 42}));" +
                "    var marker = document.getElementById('marker');"+
                "    marker.style.top = imgClick.y - 10;" +
                "    marker.style.left = imgClick.x - 10;" +
                "    showMarker('select');" +
                "};" +
                "function guide() {" +
                "    console.log('guide');" +
                "    ws.send(JSON.stringify({method:'guide'," +
                "        params:[{pixels:1.5, time:8, timeout:40}, false], id:1}));" +
                "};" +
                "function stop() {" +
                "    console.log('stop');" +
                "    ws.send(JSON.stringify({'method':'set_paused','params':[true,'full'],'id':2}));" +
                "};" +
                "function loop() {" +
                "    console.log('loop');" +
                "    ws.send(JSON.stringify({method:'loop', id:3}));" +
                "};" +
                "</script>" +
            "</head>" +
            "<body>" +
            "<div style='position: relative; left: 0; top: 0;'>" +
                "<img id='cam' src='cam.jpg' onclick='imageClick(event)' style='transform: scaleY(-1);-webkit-filter:brightness(140%)contrast(300%);position: relative; top: 0; left: 0;'>" +
                "<svg id='marker' width='20' height='20' style='position: absolute; top: 0; left: 0;'>" +
                "    <g id='m-select' style='opacity:0'>" +
                "        <rect x='-4' y='-4' width='10' height='10' stroke='white' stroke-width='2' fill='none' />" +
                "        <rect x='14' y='-4' width='10' height='10' stroke='white' stroke-width='2' fill='none' />" +
                "        <rect x='-4' y='14' width='10' height='10' stroke='white' stroke-width='2' fill='none' />" +
                "        <rect x='14' y='14' width='10' height='10' stroke='white' stroke-width='2' fill='none' />" +
                "    </g>" +
                "    <g id='m-calib' style='opacity:0'>" +
                "        <rect x='0' y='0' width='20' height='20' stroke='yellow' stroke-width='4' stroke-dasharray='2 2' fill='none' />" +
                "    </g>" +
                "    <g id='m-guide' style='opacity:0'>" +
                "        <line x1='10' y1='0' x2='10' y2='20' stroke='green' stroke-width='2' />" +
                "        <line x1='0' y1='10' x2='20' y2='10' stroke='green' stroke-width='2' />" +
                "        <rect x='4' y='4' width='12' height='12' stroke='green' stroke-width='2' fill='none' />" +
                "    </g>" +
                "    <g id='m-lost'  style='opacity:0'>" +
                "        <line x1='0' y1='0' x2='20' y2='20' stroke='red' stroke-width='2' />" +
                "        <line x1='20' y1='0' x2='0' y2='20' stroke='red' stroke-width='2' />" +
                "        <rect x='0' y='0' width='20' height='20' stroke='red' stroke-width='4' fill='none' />" +
                "    </g>" +
                "</svg>" +
            "</div>" +
            "<div class='controls' style='position:fixed;bottom:0;left:0'>" +
            "    <button onclick='guide()'>GUIDE</button>" +
            "    <button onclick='stop()'>STOP</button>" +
            "    <button onclick='loop()'>LOOP</button>" +
            "</div>" +

            "</body>" +
        "</html>"

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


