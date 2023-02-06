const net = require('net');
const crypto = require('crypto');
const {
    time
} = require('console');

// Simple HTTP server responds with a simple WebSocket client test
const httpServer = net.createServer((connection) => {
    connection.on('data', () => {
        let content = `<!DOCTYPE html>
    <html>
    
    <head>
        <meta charset="UTF-8" />
        <title>weed</title>
        <style>
            body {
                color: white;
                background-color: #181818;
                margin: 0;
                padding: 0;
                font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            }
    
            .app {
                margin: 100px auto;
                width: 800px;
            }
    
            .header {
                display: flex;
                flex-direction: row;
                align-items: center;
                background-color: #252525;
            }
    
            .inp-text {
                flex-grow: 1;
                padding: 1.5em;
            }
    
            .inp-text+.inp-text {
                border-left: 1px solid #616161;
            }
    
            .inp {
                width: 100%;
                font-size: 18px;
                color: white;
                background-color: transparent;
                border: none;
                border-bottom: 1px solid white;
            }
    
            .inp:focus-visible {
                outline: none;
            }
    
            .btn {
                margin-right: 2em;
                padding: .75em 2em;
                background-color: #7cf058;
                border: none;
                border-radius: 2px;
                cursor: pointer;
                transition: .2s;
            }
    
            .btn:hover {
                background-color: #6ddb4b;
            }
    
            .chat {
                max-height: 40em;
                overflow-y: auto;
                margin-top: 1em;
                padding: .5em 1em;
                border: 1px solid #3e3e3e;
                display: flex;
                flex-direction: column;
            }
    
            .message {
                padding: .5em 0;
            }
    
            .message+.message {
                border-top: 1px solid #e3e3e3;
            }
        </style>
    </head>
    
    <body>
        <div class="app">
            <h2 class="title">
                WebSocket chat page
            </h2>
            <div class="header">
                <div class="inp-text">
                    <input class="inp unm" type="text" placeholder="Who are you?">
                </div>
                <div class="inp-text">
                    <input class="inp message" type="text" placeholder="Message the boys">
                </div>
                <button class="btn">Send</button>
            </div>
            <div class="chat">
            </div>
            <script>
                let input = document.querySelector(".message")
                let username = document.querySelector(".unm")
                let send = document.querySelector("button")
                let chat = document.querySelector('.chat')
                let ws = new WebSocket('ws://localhost:3001');
    
                send.onclick = () => {
                    if (!input.value || !username.value) {
                        return
                    }
                    ws.send(JSON.stringify({
                        username: username.value,
                        message: input.value
                    }))
                    input.value = ""
                    input.focus()
                }
    
                ws.onmessage = (e) => {
                    let d = document.createElement('div')
                    let data = JSON.parse(e.data)
                    d.innerHTML = data.username + ': ' + data.message
                    d.className = "message"
                    //chat.appendChild(d)
                    chat.prepend(d)
                }
            </script>
        </div>
    </body>
    
    </html>`;
        connection.write('HTTP/1.1 200 OK\r\nContent-Length: ' + content.length + '\r\n\r\n' + content);
    });
});
httpServer.listen(3000, () => {
    console.log('HTTP server listening on port 3000');
});

let connections = []

const wsServer = net.createServer((connection) => {
    console.log('Client connected');

    let first = true;
    let uid = crypto.randomUUID()
    connections.push({
        uid,
        connection
    });

    connection.on('data', (data) => {
        if (first) {
            let key;
            data.toString().split('\n').forEach(l => {
                if (l.startsWith('Sec-WebSocket-Key')) {
                    key = l.split(':')[1].trim()
                }
            })

            let shasum = crypto.createHash('sha1')
            shasum.update(key)
            shasum.update("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
            let accept = shasum.digest('base64')

            connection.write('HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: ' + accept + '\r\n\r\n')
            first = false
            return
        }

        // signals close
        if ((data[0] & 0x8) === 0x8) {
            console.log("should close");
            //connection.write(Buffer.from([136, 2, 1000]))
            connection.end()
            connections = connections.filter(c => c.uid !== uid)
            return
        }

        let dataLength = data[1] & 0x7f
        let maskKey = [data[2], data[3], data[4], data[5]]
        let decoded = []
        let decodedString = ''
        for (let i = 0; i < dataLength; i++) {
            let e = data[6 + i] ^ maskKey[i % 4]
            decoded.push(e)
            decodedString += String.fromCharCode(e)
        }
        console.log("message: " + decodedString);

        let msg = Buffer.from([129, dataLength, ...decoded])

        connections.forEach(conn => conn.connection.write(msg))
    });

    connection.on('error', (error) => {
        console.log(error);
    })

    connection.on('end', () => {
        connections = connections.filter(c => c.uid !== uid)
        console.log('Client disconnected');
    });
});

wsServer.on('error', (error) => {
    console.error('Error: ', error);
});
wsServer.listen(3001, () => {
    console.log('WebSocket server listening on port 3001');
});