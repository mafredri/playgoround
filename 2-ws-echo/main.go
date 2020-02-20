// From https://github.com/gorilla/websocket/tree/master/examples/echo
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Headers: %v\n", r.Header)
	u, p, ok := r.BasicAuth()
	fmt.Printf("Basic: %q %q %v\n", u, p, ok)

	fmt.Println("ReqURI:", r.RequestURI)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		fmt.Printf("Headers: %v\n", r.Header)
		u, p, ok := r.BasicAuth()
		fmt.Printf("Basic: %q %q %v\n", u, p, ok)

		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	fmt.Printf("Basic: %q %q %v\n", u, p, ok)
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func check(user, pass string) bool {
	return false
}
func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		log.Println(user, pass, ok)
		if !check(user, pass) {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		fn(w, r)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", auth(echo))
	mux.HandleFunc("/", home)
	fmt.Printf("Listening on: https://%s\n", *addr)

	log.Fatal(ListenAndServeTLS(*addr, "cert.pem", "key.pem", mux))
}

func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	server := &Server{Server: &http.Server{Addr: addr, Handler: handler}}
	server.TLSConfig = &tls.Config{}
	return server.ListenAndServeTLS(certFile, keyFile)
}

type Server struct {
	*http.Server
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return srv.ServeTLS(tcpKeepAliveListener{ln.(*net.TCPListener)}, certFile, keyFile)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
    window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
