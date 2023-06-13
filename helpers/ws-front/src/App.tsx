import React, { useEffect, useState } from 'react';

import './App.css';

import { Button } from '@material-ui/core';
import CryptoJS from "crypto-js";
import { v4 as uuidv4 } from 'uuid';

const APIKey: string = "lVKxR8F0dSFrL3gt3QIstQ==";
const SecretKey: string = "6QJY5tT_39BtkfXJdZjs2ZbL2xlQHEQIkT6fpPJQbYo=";


function App() {
  const [socket, setSocket] = useState<WebSocket>();
  useEffect(() => {
    if (socket) {
      socket.onopen = function (e) {
        console.log("[open] Connection established");
        console.log("Sending to server");

      };

      socket.onmessage = function (event) {
        console.log(`[message] Data received from server: ${event.data}`);
      };

      socket.onclose = function (event) {
        if (event.wasClean) {
          console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
        } else {
          // e.g. server process killed or network down
          // event.code is usually 1006 in this case
          console.log('[close] Connection died');
        }
      };

      socket.onerror = function (error: any) {
        console.log(`[error] ${error.message}`);
      };
    } else {

      const nonce: string = uuidv4();
      const signature: string = CryptoJS.HmacSHA256(nonce, SecretKey).toString(CryptoJS.enc.Base64);
      const params: URLSearchParams = new URLSearchParams();
      params.set("Signature", signature);
      params.set("Nonce", nonce);
      params.set("API-Key", APIKey);
      params.toString()
      setSocket(new WebSocket("ws://127.0.0.1:8080/ws?" + params.toString()));
    }


  }, [socket, setSocket]);

  const subscribe = () => {
    socket?.send(JSON.stringify({
      "event": "subscribe",
      "instrument": "btceur",
      "levels": [0.5, 1]
    }));
  }


  return (
    <div className="App">
      <button onClick={() => subscribe()}>subscribe</button>
    </div>
  );
}

export default App;
