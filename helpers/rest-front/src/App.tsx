import React from 'react';

import './App.css';

import { Button } from '@material-ui/core';
import CryptoJS from "crypto-js";
import { v4 as uuidv4 } from 'uuid';

const APIKey: string = "lVKxR8F0dSFrL3gt3QIstQ==";
const SecretKey: string = "6QJY5tT_39BtkfXJdZjs2ZbL2xlQHEQIkT6fpPJQbYo=";


function App() {

  async function buy() {
    const nonce: string = uuidv4();
    const signature: string = CryptoJS.HmacSHA256(nonce, SecretKey).toString(CryptoJS.enc.Base64);
    try {
      await fetch("http://127.0.0.1:8083/order",
        {
          method: "POST",
          mode: 'cors',
          headers: {
            "Nonce": nonce,
            "Signature": signature,
            "API-Key": APIKey
          },
          body: JSON.stringify({
            "instrument": "btceur",
            "quantity": 0.3,
            "side": "0"
          })
        });
    } catch (e) {
      console.log(e);
    }

  }

  async function sell() {
    const nonce: string = uuidv4();
    const signature: string = CryptoJS.HmacSHA256(nonce, SecretKey).toString(CryptoJS.enc.Base64);
    try {
      await fetch("http://127.0.0.1:8083/order",
        {
          headers: {
            "Nonce": nonce,
            "Signature": signature,
            "API-Key": APIKey
          },
          mode: 'cors',
          method: "POST",
          body: JSON.stringify({
            "instrument": "btceur",
            "quantity": 0.3,
            "side": "1"
          })
        });
    } catch (e) {
      console.log(e);
    }

  }


  return (
    <div className="App">
      <Button
        variant="contained"
        color="primary"
        onClick={buy}
      >
        Buy
            </Button>  <Button
        variant="contained"
        color="primary"
        onClick={sell}
      >
        Sell
            </Button>
    </div>
  );
}

export default App;
