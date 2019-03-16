package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}

type Team struct {
	Team string `json:"team"`
	Funs string `json:"funs"`
}

type Message struct {
	Type int `json:"type"`         				// Тип цвета, по нему распознаём, что делать дальше.
	Payload json.RawMessage `json:"payload"` 	// Внутренний JSON, который будет парситься
}

type Answer struct {
	Message string `json:"message"`
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", home)
	http.Handle("/",router)
	http.HandleFunc("/echo", echo)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)

}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mess := Message{}
		err := c.ReadJSON(&mess)


		if err != nil {
			log.Println("read:", err)
			break
		}

		fmt.Println("Type: ", mess.Type)

		if 1 == mess.Type {
			data := Team{}
			err = json.Unmarshal(mess.Payload, &data)
			fmt.Println(data)
		}

		if 2 == mess.Type {
			data := User{}
			err = json.Unmarshal(mess.Payload, &data)
			fmt.Println(data)
		}

		answer := Answer{"I've got it!"}

		if err = c.WriteJSON(answer); err != nil {
			fmt.Println(err)
		}

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
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
    document.getElementById("send_team").onclick = function(evt) {
        if (!ws) {
            return false;
        }
		print("SEND Team...");
		team = JSON.stringify({"type": 1, "payload": {"team": "Boston Red Sox", "funs": "many"}});
		console.log(team)
        ws.send(team);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
	document.getElementById("send_id").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND ID...");

		user = JSON.stringify({"type": 2, "payload": {"id": 555, "name": "Alex", "age": 44}});
		console.log(user)
        ws.send(user);
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
<button id="send_team">Send Team</button>
<button id="send_id">Send Id</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))


