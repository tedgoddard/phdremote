// Copyright 2014 Ted Goddard. All rights reserved.
// Use of this source code is governed the Apache 2.0
// license that can be found in the LICENSE file.

package main

import "fmt"
import "flag"
import "net"
import "io/ioutil"
import "os"
import "os/signal"
import "os/exec"
import "path/filepath"
import "strings"
import "bufio"
import "net/http"
import "websocket"
import "log"
import "fits"
import "phdremote"
import "encoding/json"

func main() {
	fmt.Println("Starting server")
    var userScriptPath = flag.String("solver", "",
            "path to solver script arguments inFile outFile")
    var fakeImagePath = flag.String("fake", "",
            "path to fake image for testing")
    flag.Parse()

    tmpDir, err := ioutil.TempDir("", "phdremote")
    if (nil != err)  {
        log.Print("could not find temp directory ", err)
    }

    previousImagePath := ""
    currentImagePath := ""

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
    go func(){
        for _ = range sigChan {
            fmt.Println("killed, orphan ", currentImagePath)
//            os.Remove(currentImagePath)
            os.Exit(0)
        }
    }()


    http.HandleFunc("/phdremote/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, phdremote.ClientHTML)
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
                            if (*fakeImagePath != "") {
                                currentImagePath = *fakeImagePath
                            } else {
                                newImagePath := result["filename"].(string)
                                _, name := filepath.Split(newImagePath)
                                currentImagePath = filepath.Join(tmpDir, name)
                                os.Rename(newImagePath, currentImagePath)
                                if ("" != previousImagePath)  {
    //                                os.Remove(previousImagePath)
                                }
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

    UserScript := func(inPath string) string {
        if (*userScriptPath == "") {
            return ""
        }
        outPath := inPath + "sol.jpg"
        cmd := exec.Command(*userScriptPath, inPath, outPath)
        buf, err := cmd.CombinedOutput()
        if err != nil {
            log.Print("Unable to execute user script ", outPath, " ", err)
            log.Print(string(buf))
            return "error.jpg"
        }
        log.Print(string(buf))
        return outPath
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

    http.HandleFunc("/phdremote/solved.jpg", func(w http.ResponseWriter, r *http.Request) {
log.Print("returning solved jpg image")
        momentaryImagePath := currentImagePath
log.Print("run script ", *userScriptPath, " on ", momentaryImagePath)
        outPath := UserScript(momentaryImagePath)
        if (outPath == "") {
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "image/jpeg")
        http.ServeFile(w, r, outPath)
    })

    http.HandleFunc("/phdremote/solved.wcsinfo", func(w http.ResponseWriter, r *http.Request) {
        momentaryImagePath := currentImagePath
        outPath := filepath.Join(filepath.Dir(momentaryImagePath),
                "previous.wcsinfo")
        log.Print("returning solved wcs info for ", outPath)
        w.Header().Set("Content-Type", "text/plain")
        http.ServeFile(w, r, outPath)
    })

    log.Print("http.ListenAndServe")
    log.Fatal(http.ListenAndServe(":8080", nil))

}


