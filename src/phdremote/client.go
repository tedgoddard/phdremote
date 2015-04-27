package phdremote

    const ClientHTML = `
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
          .trcontrols {
            position:fixed;
            top:10px;
            right:10px;
          }
          .brcontrols a, .trcontrols a {
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
              .tcontrols {
                position:fixed;
                top:100px;
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
              .trcontrols {
                position:fixed;
                top:10px;
                right:10px;
              }
          }
          @media (min-width: 641px) {
              .tcontrols {
                position:fixed;
                top:20px;
                left:50%%;
              }
              .bcontrols {
                position:fixed;
                bottom:20px;
                left:50%%;
              }
              .bcinner {
                margin-left:-50%%;
              }
              .tcinner {
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
        var currentInfo = { };
        var currentTarget = "";
        var targetRA = 0.0;
        var targetDEC = 0.0;
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
        var startX = 0;
        var startY = 0;
        var newX = 0;
        var newY = 0;
        var camContrast = 3.0;
        var camBrightness = 1.4;
        var startContrast = 3.0;
        var startBrightness = 1.4;
        function adjustStart(event) {
            startX = event.pageX;
            startY = event.pageY;
            startContrast = camContrast;
            startBrightness = camBrightness;
        }
        function adjustImage(event) {
            var deltaX = event.pageX - startX;
            var deltaY = event.pageY - startY;
            camContrast = startContrast + deltaX / 100.0;
            camBrightness = startBrightness + deltaY / 100.0;
            var camElement = document.getElementById("cam");
            camElement.style.webkitFilter =
                "brightness(" + camBrightness + ") contrast(" + camContrast + ")";
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
        function toggleSolved() {
            var solvedElement = document.getElementById("solvedfield");
            var solvedSpinner = document.getElementById("solvedspinner");
            var newOpacity = 0.5 - solvedElement.style["opacity"];
            if (newOpacity > 0) {
                solvedSpinner.beginElement();
                solvedElement.src = "solved.jpg?" + new Date().getTime();
                solvedElement.onload = function() {
                    solvedElement.style["opacity"] = newOpacity;
                    solvedSpinner.endElement();
                    getAndDisplayInfo();
               }
            } else {
                solvedElement.style["opacity"] = newOpacity;
            }
        }
        function processInfo(data) {
            var newInfo = data.split("\n");
            for (index in newInfo) {
                var entry = newInfo[index].split(" ");
                currentInfo[entry[0]] = entry[1];
            }
            if (currentTarget) {
                findField(currentTarget);
            }
        }
        function getAndDisplayInfo() {
            httpGet("solved.wcsinfo?" + new Date().getTime(), processInfo);
        }
        function testMarkers() {
            displayMarker("Medusa", 112.48, 13.20694);
            displayMarker("HD 61199", 114.5745875, 4.94234722);
            displayMarker("HD 61112", 114.45529167, 4.53981389);
            displayMarker("HD 60803", 114.1445625, 5.86169722);
            displayMarker("HD 61696", 115.195175, 5.00018611);
            displayMarker("Procyon", 115.028, 5.187222);
            displayMarker("HD 61664", 115.17254583, 6.62442778);
        }
        function processLookup(data) {
            console.log("lookup response " + data);
            var coords = JSON.parse(data);
            displayMarker("*", coords.ra, coords.dec);

        }
        function findField(targetText) {
            targetText = targetText.toLowerCase();
            currentTarget = targetText;
            var targetCoords = targetText.split(",");
            if (targetCoords.length > 1) {
                targetRA = parseFloat(targetCoords[0]);
                targetDEC = parseFloat(targetCoords[1]);
                displayMarker("+", targetRA, targetDEC);
            } else {
                httpGet("lookup?o=" + targetText, processLookup);
            }

        }
        function displayMarker(label, targetRA, targetDEC) {
            var arrowElement = document.getElementById("arrow");
            var camElement = document.getElementById("cam");
            arrowElement.style.opacity = 1.0;
            var deltaRA = targetRA - currentInfo.ra_center;
            var deltaDEC = targetDEC - currentInfo.dec_center;
            var rads =  currentInfo.orientation / 180.0 * Math.PI;
            var pointRA = deltaRA * Math.cos(rads) - deltaDEC * Math.sin(rads);
            var pointDEC = deltaRA * Math.sin(rads) + deltaDEC * Math.cos(rads);
            var norm = Math.sqrt(pointRA * pointRA + pointDEC * pointDEC);
            var scaledRA = pointRA / norm * camElement.height * 0.4;
            var scaledDEC = pointDEC / norm * camElement.height * 0.4;
            console.log("scaledRA " + scaledRA + " scaledDEC " + scaledDEC);
            var rotation = (90 + 180 * Math.atan(scaledDEC / scaledRA) / Math.PI);
            if (scaledRA < 0) {
                rotation = rotation + 180;
            }
            arrowElement.style.top = ((camElement.height / 2 + scaledDEC)) + "px";
            arrowElement.style.left = (camElement.width / 2 + scaledRA) + "px";
            arrowElement.firstElementChild.setAttribute("transform",
                    "rotate(" + rotation + " 20 20)");
        }
        function depositMarker(label, scaledRA, scaledDEC) {
            var camElement = document.getElementById("cam");
            var marker = document.createElement("div");
            marker.innerHTML = label;
            marker.style.top = ((camElement.height / 2 + scaledDEC)) + "px";
            marker.style.left = (camElement.width / 2 + scaledRA) + "px";
            marker.style.position = "absolute";
            document.body.appendChild(marker);
        }
        function imageDispatch(event) {
            var solvedElement = document.getElementById("solvedfield");
            if (solvedElement.style["opacity"] > 0) {
                solvedClick(event);
            } else {
                imageClick(event);
            }
        }
        var solvedClickGal = true;
        function solvedClick(event) {
            var camElement = document.getElementById("cam");
            var galIcon = document.getElementById("galIcon");
            var starIcon = document.getElementById("starIcon");
            var destIcon = document.getElementById("destIcon");
            var theIcon = galIcon;

            destIcon.style["opacity"] = 0.0;
            if (!solvedClickGal) {
                theIcon = starIcon;
                destIcon.style["opacity"] = 1.0;
            }

            var pos = getClickPosition(event);
            var iconSize = svgSize(theIcon);
            var destIconSize = svgSize(destIcon);

            theIcon.style.top = pos.y - (iconSize / 2);
            theIcon.style.left = pos.x - (iconSize / 2);
            theIcon.style["opacity"] = 1.0;
            solvedClickGal = !solvedClickGal;

            var totalExtra = svgSize(galIcon) / 2;
            destIcon.style.top = svgTop(starIcon) +
                    ((camElement.height / 2) - svgTop(galIcon)) - totalExtra;
            destIcon.style.left = svgLeft(starIcon) -
                    (svgLeft(galIcon) - (camElement.width / 2)) - totalExtra;

        };
        function httpGet(url, callback) {
            var xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function() {
                if ((xhr.readyState) == 4 && (xhr.status == 200)) {
                   callback(xhr.responseText);
                }
            }
            xhr.open("GET", url, true);
            xhr.send();
        }
        function svgTop(elm) {
            return parseFloat(elm.style.top);
        }
        function svgLeft(elm) {
            return parseFloat(elm.style.left);
        }
        function svgSize(elm) {
            return parseFloat(elm.getAttribute("height"));
        }
        function adjustSizes() {
            var bullseyeElement = document.getElementById("bull");
            var camElement = document.getElementById("cam");
            bullseyeElement.style.width = camElement.width;
            bullseyeElement.style.height = camElement.height;
            var solvedElement = document.getElementById("solvedfield");
            solvedElement.style.width = camElement.width;
            solvedElement.style.height = camElement.height;
        }
        window.onresize = function(event)  {
            adjustSizes();
        }
        </script>
    </head>
    <body>
    <div class="imgBox" onclick="imageDispatch(event)">
        <img id="cam" src="cam.jpg" onclick="imageClick(event)" onload="adjustSizes()"
            style="-webkit-filter:brightness(140%%)contrast(300%%);position: relative; top: 0; left: 0;">
        <img id="solvedfield" onload="adjustSizes()" onclick="solvedClick(event)"
            onerror="this.style.display='none';"
            style="position: absolute; top: 0; left: 0;">
        <svg id="bull" width="100%%" height="100%%" style="opacity:0; position: absolute; top: 0; left: 0;">
            <g >
                <line x1="0px" y1="50%%" x2="100%%" y2="50%%" stroke="red" stroke-width="1" />
                <line x1="50%%" y1="0px" x2="50%%" y2="100%%" stroke="red" stroke-width="1" />
                <circle cx="50%%" cy="50%%" r="10%%" stroke="red" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="4%%" stroke="red" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="2%%" stroke="red" stroke-width="1" fill="none" />
            </g>
        </svg>
        <svg id="arrow" width="40" height="40" style="opacity: 0; position: absolute; top: 0; left: 0;" >
            <polygon points="20,0 25,20 22,20 23,40 17,40 18,20 15,20" stroke="red" stroke-width="1.0" fill="firebrick"/>
        </svg>
        <svg id="galIcon" width="40" height="40" style="opacity: 0; position: absolute; top: 0; left: 0;" >
            <circle cx="50%%" cy="50%%" r="48%%" stroke="red" stroke-width="1.0" fill="none"/>
        </svg>
        <svg id="starIcon" width="10" height="10" style="opacity: 0; position: absolute; top: 0; left: 0;" >
            <circle cx="50%%" cy="50%%" r="48%%" stroke="red" stroke-width="1.0" stroke-dasharray="2 2" fill="none"/>
        </svg>
        <svg id="destIcon" width="10" height="10" style="opacity: 0; position: absolute; top: 0; left: 0;">
            <circle cx="50%%" cy="50%%" r="48%%" stroke="red" stroke-width="1.0" fill="none"/>
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
    <div class="tcontrols" >
      <div class="tcinner" >
        <input onchange="findField(this.value);" >
      </div>
    </div>
    <div class="trcontrols" >
      <div class="trinner" >
        <a draggable="true"
            ontouchstart="adjustStart(event)" ondragstart="adjustStart(event)"
            ondrag="adjustImage(event)" ontouchmove="adjustImage(event)">
          <svg width="60px" height="60px">
            <g >
                <path d="M30,10 L30,50 A20,20 0 0,1 30,10 z" fill="black" />
                <path d="M30,50 L30,10 A20,20 0 0,1 30,50 z" fill="firebrick" />
            </g>
          </svg>
        </a>
      </div>
    </div>
    <div class="brcontrols" >
      <div class="brinner" >
        <a onclick="toggleSolved()">
            <svg width="40px" height="40px">
            <g >
                <animateTransform id="solvedspinner"
                    attributeName="transform"
                    attributeType="XML"
                    type="rotate"
                    from="0 20 20"
                    to="360 20 20"
                    dur="10s"
                    begin="indefinite"
                    repeatCount="indefinite"/>
                <line x1="60%%" y1="30%%" x2="20%%" y2="60%%" stroke="black" stroke-width="1" />
                <line x1="20%%" y1="60%%" x2="80%%" y2="80%%" stroke="black" stroke-width="1" />
                <line x1="80%%" y1="80%%" x2="60%%" y2="30%%" stroke="black" stroke-width="1" />
                <circle cx="60%%" cy="30%%" r="8%%" stroke="black" stroke-width="1" fill="firebrick" />
                <circle cx="20%%" cy="60%%" r="8%%" stroke="black" stroke-width="1" fill="firebrick" />
                <circle cx="80%%" cy="80%%" r="8%%" stroke="black" stroke-width="1" fill="firebrick" />
            </g>

            </svg>
        </a>
        <a onclick="toggleBullseye()">
            <svg width="40px" height="40px">
            <g >
                <line x1="0px" y1="50%%" x2="100%%" y2="50%%" stroke="black" stroke-width="1" />
                <line x1="50%%" y1="0px" x2="50%%" y2="100%%" stroke="black" stroke-width="1" />
                <circle cx="50%%" cy="50%%" r="20%%" stroke="black" stroke-width="1" fill="none" />
                <circle cx="50%%" cy="50%%" r="10%%" stroke="black" stroke-width="1" fill="none" />
            </g>
            </svg>
        </a>
      </div>
    </div>
    </body>
</html>
`
