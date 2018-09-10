package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	startReloadServer()
	go startLogViewer()
	// done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
			sendReload()
		}
	}()
	for {
		err = watcher.Add("/tmp/duo.json")
		if err != nil {
			time.Sleep(time.Second)
		}
	}
	// <-done
}

func startLogViewer() {
	port := flag.String("p", "8100", "port to serve on")
	file := flag.String("f", "/tmp/duo.json", "the file of static file to host")
	flag.Parse()
	loghandler := func(w http.ResponseWriter, req *http.Request) {
		j, err := ioutil.ReadFile(*file)
		if j == nil {
			j = []byte("{}")
		}
		J := string(j)
		if J == "" {
			j = []byte("{}")
		}
		j = []byte(strings.Replace(string(j), `\t`, ``, -1))
		j = []byte(strings.Replace(string(j), `\"`, "`", -1))
		c, err := ioutil.ReadFile("./jsoneditor.min.css")
		js, err := ioutil.ReadFile("./jsoneditor.min.js")
		if err == nil {
		}
		ic, err := ioutil.ReadFile("./jsoneditor-icons.svg")
		if err == nil {
			icons := "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(ic)
			if err == nil {
				css := strings.Replace(string(c), "img/jsoneditor-icons.svg", icons, -1)
				io.WriteString(w, `
			<head>
				<title>JSON log viewer</title>
				
				<style>`+
					css+`
						</style>
				</head>
					<body>
						<div id="jsoneditor" style="width: 100%;"></div>
				<script>
				
				function tryConnectToReload(address) {
					var conn = new WebSocket(address);
			
					conn.onclose = function () {
					  setTimeout(function () {
						 tryConnectToReload(address);
					  }, 2000);
					};
			
					conn.onmessage = function (evt) {
					
					  // If we uncomment this line, then the page will refresh every time a message is received.
					  location.reload()
					};
				 }
			
				 try {
					if (window["WebSocket"]) {
					  // The reload endpoint is hosted on a statically defined port.
					  try {
						 tryConnectToReload("ws://localhost:12450/reload");
					  } catch (ex) {
						 // If an exception is thrown, that means that we couldn't connect to to WebSockets because of mixed content
						 // security restrictions, so we try to connect using wss.
						 tryConnectToReload("wss://localhost:12451/reload");
					  }
					} else {
					  console.log("Your browser does not support WebSockets, cannot connect to the Reload service.");
					}
				 } catch (ex) {
					console.error('Exception during connecting to Reload:', ex);
				 }

				`+string(js)+`
				var container = document.getElementById("jsoneditor");
				var options = { modes: ['tree','view','form','code','text'], mode: 'view', name: 'logs' };
				var editor = new JSONEditor(container, options);
				var json = `+string(j)+`;editor.set(json);
				document.getElementById('getJSON').onclick = function () {
					var json = editor.get();
					alert(JSON.stringify(json, null, 2));
				};
				document.querySelectorAll('.jsoneditor-string').forEach(function(div) {
						
				});
					</script>
						</body>
				`)
			}
		}
	}
	http.HandleFunc("/", loghandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is an middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("An error happened when reading from the Websocket client: %v", err)
			}
			break
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Client) write(mt int, payload []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}

var (
	hub *Hub
	// The port on which we are hosting the reload server has to be hardcoded on the client-side too.
	reloadAddress    = ":12450"
	reloadAddressTLS = ":12451"
)

const (
	certContent = `-----BEGIN CERTIFICATE-----
MIICkzCCAfwCCQCbmnQ2PFatzzANBgkqhkiG9w0BAQsFADCBjTELMAkGA1UEBhMC
TkwxEzARBgNVBAgMClNvbWUtU3RhdGUxEjAQBgNVBAcMCUFtc3RlcmRhbTEPMA0G
A1UECgwGVHJhdml4MQ0wCwYDVQQLDARDb3JlMRIwEAYDVQQDDAlsb2NhbGhvc3Qx
ITAfBgkqhkiG9w0BCQEWEm12aW5jemVAdHJhdml4LmNvbTAeFw0xNjEwMTUxNzEy
NTVaFw0xOTA4MDUxNzEyNTVaMIGNMQswCQYDVQQGEwJOTDETMBEGA1UECAwKU29t
ZS1TdGF0ZTESMBAGA1UEBwwJQW1zdGVyZGFtMQ8wDQYDVQQKDAZUcmF2aXgxDTAL
BgNVBAsMBENvcmUxEjAQBgNVBAMMCWxvY2FsaG9zdDEhMB8GCSqGSIb3DQEJARYS
bXZpbmN6ZUB0cmF2aXguY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCf
OSB7LkkaRd6WXTplUHfD2k+EHoVi9flKcmbUlye9zHFzWVtCUQjhFjiZL1rNRQGn
9VMUqzpc55RyzTEy2KpyZ+7INR1ZAuqqXMxpNDzXeq+UQuAFnJrHnbwtiSYPiJ45
5EvysllYb5j6ihXEVZt+6QdMINFB+Gz0Xfrhug0+0QIDAQABMA0GCSqGSIb3DQEB
CwUAA4GBADrH8ibFye3iXHR6RkwVNBgeKyvL0kxs4C8785uYqjRJWVjAg2xJQyyZ
R3IHuvKqkmjs5i5d5CT9QT4t8Mlorg1XSnRz/HLf5zrRJlVzqrpd9N2+859TmTVD
9A91NtEwCNgBSGDGSCndjQ/dkPhbJFs28/ICujLySxbYswOGHGbK
-----END CERTIFICATE-----
`
	keyContent = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCfOSB7LkkaRd6WXTplUHfD2k+EHoVi9flKcmbUlye9zHFzWVtC
UQjhFjiZL1rNRQGn9VMUqzpc55RyzTEy2KpyZ+7INR1ZAuqqXMxpNDzXeq+UQuAF
nJrHnbwtiSYPiJ455EvysllYb5j6ihXEVZt+6QdMINFB+Gz0Xfrhug0+0QIDAQAB
AoGAIo9SxonwYhyCSN7peu4xYLh1A/df+m/rcUZNnZ1FigPjKCdgEI/oPnsFQ/Ks
Ydu1lVBBfT4BSAMYDKcPI7s1m5Hf++2TAWXuE/GiMmfmQq8QHVwdRERIzGo7BSIW
alA5tC4+dIe5gUKjR38MpG9VCEa3FBkNxlRQ2U1tIAoM9/ECQQDLWvbShPYpfKCM
8WlAGeWwgHJrjdmatMLsJepxFjGShxK1uhLy6mIMaVVCV0dFPk2Y81ACAirmev99
bqMd3sbtAkEAyHFgTZzQUrezQQhnfFcEDOaUrCwRBVERHFou6wHEwTLObJeedAuo
emRRpQkOp+wJq8y9eOI2pv0jpSI8pTKW9QJAdOuzOG1sX4Qhh4gSHOIG90mTABYK
BHJkFITkW+sHy5jQAB6hYHu0rjAt7jviZYSh9wwGd3Epm2Ui2sqvDLCXLQJBAKAk
NNTNXIM50TU8CbIFs267Kj0EV/Tvd8Q3KRUJLLFObi3EVQxR5CEk1TYNrm/q3S8t
PJO/5/oydLASUnGJoaECQGyPpJ6lVJb10yJKjcGtouwa+HFRJh9BxIQUHZRTbmHX
k7iRrF0Vcllo8k/Mos5PVPP0WIyS1l0lh4GZ+w8gA80=
-----END RSA PRIVATE KEY-----
`
)

func createCertFiles() (cert string, key string) {
	tempFolder, _ := ioutil.TempDir("", "reload")

	cert = tempFolder + "/reload-cert.pem"
	key = tempFolder + "/reload-key.pem"

	ioutil.WriteFile(cert, []byte(certContent), 0644)
	ioutil.WriteFile(key, []byte(keyContent), 0644)

	return cert, key
}

func startReloadServer() {
	hub = newHub()
	go hub.run()
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	go startServer()
	go startServerTLS()
	log.Println("Reload server listening at", reloadAddress)
}

func startServer() {
	err := http.ListenAndServe(reloadAddress, nil)

	if err != nil {
		log.Println("Failed to start up the Reload server: ", err)
		return
	}
}

func startServerTLS() {
	cert, key := createCertFiles()
	err := http.ListenAndServeTLS(reloadAddressTLS, cert, key, nil)

	if err != nil {
		log.Println("Failed to start up the Reload server with TLS: ", err)
		return
	}
}

func sendReload() {
	message := bytes.TrimSpace([]byte("reload"))
	hub.broadcast <- message
}
