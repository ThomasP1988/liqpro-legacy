import React from 'react';
// import './App.css';
import useMediaQuery from '@material-ui/core/useMediaQuery';
import { ThemeProvider, jssPreset, StylesProvider } from '@material-ui/core/styles';
import { ThemeProvider as StyledComponent } from "styled-components/macro";
import { ApolloProvider } from '@apollo/client';

import { create } from "jss";

import CssBaseline from '@material-ui/core/CssBaseline';
import SuperTokens from "supertokens-auth-react";
import { doesSessionExist  } from 'supertokens-auth-react/recipe/session';
import { createTheme } from "./theme/theme";
import {
  BrowserRouter as Router
} from "react-router-dom";
import { Routes } from "./Routes";
import SuperTokensRequest from 'supertokens-website';
import {
  ApolloClient,
  InMemoryCache,
  HttpLink,
  NormalizedCacheObject,
} from '@apollo/client';
import Session, {addAxiosInterceptors} from "supertokens-auth-react/recipe/session";
import axios from "axios";
addAxiosInterceptors(axios);


// interception is on by default
SuperTokensRequest.init({
  apiDomain: "127.0.0.1:8081"
});

export const Client = new ApolloClient<NormalizedCacheObject>({
  link: new HttpLink({
      uri: "https://localhost:8081/protected/graphql",  // change this depending on your path
  }),
  cache: new InMemoryCache(),  // change this depending on your preference
  // ... other params
});


SuperTokens.init({
  appInfo: {
    appName: "liqpro",
    apiDomain: "https://localhost:8081",
    websiteDomain: "https://localhost:3000"
  },
  recipeList: [
    Session.init(
      {
        autoAddCredentials: true
      }
    )
  ]

});

const jss = create({
  ...jssPreset(),
  insertionPoint: document.getElementById("jss-insertion-point")!,
});


function App() {

  const prefersDarkMode: boolean = useMediaQuery('(prefers-color-scheme: dark)');
  console.log("doesSessionExist", doesSessionExist());
  const theme = React.useMemo(
    () =>
      createTheme(prefersDarkMode),
    [prefersDarkMode],
  );

  return (
    <div className="App">
      <StylesProvider jss={jss}>
        <ThemeProvider theme={theme}>
          <StyledComponent theme={theme}>
          <ApolloProvider client={Client}>
            <Router>
              <CssBaseline />
              <Routes />
            </Router>
            </ApolloProvider>
          </StyledComponent>
        </ThemeProvider>
      </StylesProvider>
    </div>
  );
}

export default App;
