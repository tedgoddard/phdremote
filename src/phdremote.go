// Copyright 2014 Ted Goddard. All rights reserved.
// Use of this source code is governed the Apache 2.0
// license that can be found in the LICENSE file.

package main

import "fmt"
import "net"
import "os"
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
                "<script>" +
                "var ws = new WebSocket('ws://' + location.host + '/echo/');" +
                "ws.onmessage = function(msg) {console.log(msg.data);" +
                "    var msgJSON = JSON.parse(msg.data);" +
                "    console.log(msgJSON.Event);" +
                "    var marker = document.getElementById('marker');"+
                "    if ('LoopingExposures' == msgJSON.Event)  {" +
                "        var camImg = document.getElementById('cam');" +
                "        camImg.src = 'cam.jpg?' + new Date().getTime();" +
                "    };" +
                "    if ('StartCalibration' == msgJSON.Event)  {" +
                "       marker.firstElementChild.setAttribute('stroke', 'yellow');" +
                "    };" +
                "    if ('GuideStep' == msgJSON.Event)  {" +
                "       marker.firstElementChild.setAttribute('stroke', 'green');" +
                "       marker.firstElementChild.style['stroke-dasharray'] = null;" +
                "    };" +
                "};" +
    
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
                "    marker.firstElementChild.style['stroke-dasharray'] = '2 2';" +
                "};" +
                "function guide() {" +
                "    console.log('guide');" +
                "    ws.send(JSON.stringify({method:'guide'," +
                "        params:[{pixels:1.5, time:8, timeout:40}, false], id:1}));" +
                "};" +
                "</script>" +
            "</head>" +
            "<body>" +
            "<div style='position: relative; left: 0; top: 0;'>" +
                "<img id='cam' src='cam.jpg' onclick='imageClick(event)' style='transform: scaleY(-1);-webkit-filter:brightness(140%)contrast(300%);position: relative; top: 0; left: 0;'>" +
                "<svg id='marker' width='20' height='20' style='position: absolute; top: 0; left: 0;'>" +
                "    <rect x='0' y='0' width='20' height='20' stroke='green' stroke-width='4' fill='none' />" +
                "</svg>" +
            "</div>" +
            "<button style='position:fixed;bottom:0;left:0' onclick='guide()'>GUIDE</button>" +

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


