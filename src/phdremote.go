// Copyright 2014 Ted Goddard. All rights reserved.
// Use of this source code is governed the Apache 2.0
// license that can be found in the LICENSE file.

package main

import "fmt"
import "net"
import "bufio"
import "net/http"
import "websocket"
import "log"
import "fits"
import "encoding/json"

func main() {
	fmt.Println("Starting server")

    currentImagePath := "RCA.fit"

    wsClientHTML :=
        "<html>" +
            "<head>" +
                "<script>" +
                "var ws = new WebSocket('ws://' + location.host + '/echo/');" +
                "ws.onmessage = function(msg) {console.log(msg.data);" +
                "    var msgJSON = JSON.parse(msg.data);" +
                "    console.log(msgJSON.Event);" +
                "    if ('LoopingExposures' == msgJSON.Event)  {" +
                "        var camImg = document.getElementById('cam');" +
                "        camImg.src = 'cam.png?' + new Date().getTime();" +
                "    };" +
                "};" +
                "function imageClick(event) {" +
                "   console.log('click' + event.x + ' ' + event.y);" +
                "   ws.send(JSON.stringify({method: 'set_lock_position', params: [event.x, event.y], id: 42}));" +
                "};" +
                "</script>" +
            "</head>" +
            "<body>" +
            "<img id='cam' src='cam.png' onclick='imageClick(event)'>" +
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
                            currentImagePath = result["filename"].(string)
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
        fits.Convert(currentImagePath, w)
    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8080", nil))

}


