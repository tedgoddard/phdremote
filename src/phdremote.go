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
                "</script>" +
            "</head>" +
            "<body>" +
            "<img id='cam' src='cam.png'>" +
            "</body>" +
        "</html>"

    http.HandleFunc("/phdremote/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, wsClientHTML)
    })

    //only one websocket allowed for now
    var wsConn *websocket.Conn
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
                    rpcResult := phdMessage["result"].(map[string]interface{})
                    currentImagePath = rpcResult["filename"].(string)
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

    EchoServer := func(newConn *websocket.Conn) {
        log.Print("EchoServer started")
        wsConn = newConn
        echoDone := <-phdDone
        if (echoDone)  {
            log.Print("EchoServer done")
       }
    }

    log.Print("websocket.Handler")
    wsHandler := websocket.Handler(EchoServer)
	http.Handle("/echo/", wsHandler)

    conn, err := net.Dial("tcp", "localhost:4400")
    if err == nil {
        go GuideWatch (&conn)
    } else {
        log.Print("Unable to connect to PHD")
    }

    http.HandleFunc("/phdremote/cam.png", func(w http.ResponseWriter, r *http.Request) {
log.Print("returning image")
        connWrite := bufio.NewWriter(conn)
        fmt.Fprintf(connWrite, "{\"method\":\"save_image\",\"id\":123}\n")
        connWrite.Flush()
        w.Header().Set("Content-Type", "image/png")
        fits.Convert(currentImagePath, w)
    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8080", nil))

}


