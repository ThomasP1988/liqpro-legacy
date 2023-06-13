import React from 'react';

import {
    Switch,
    Route,
    Redirect
} from "react-router-dom";
import { SignIn } from "./auth/SignIn";
import { SignUp } from "./auth/SignUp";
import { Home } from "./pages/Home";
import { ApiKeys } from "./pages/account/ApiKeys";
import { Portfolio } from "./pages/Portfolio";
import { Orders } from "./pages/Orders";
import { Exchange } from "./pages/Exchange";
import Dashboard from "./layout/Dashboard";
import { doesSessionExist } from 'supertokens-auth-react/recipe/session';

function PrivateRoute({ children, ...rest }: any) {
    return (
        <Route {...rest} render={() => {
            return doesSessionExist() === true
                ? children
                : <Redirect to='/login' />
        }} />
    )
}


export function Routes() {

    return (

        <Switch>
            <Route exact path="/">
                <Redirect to='/dashboard' />
            </Route>
            <Route path="/login">
                <SignIn />
            </Route>
            <Route path="/register">
                <SignUp />
            </Route>
            <PrivateRoute path="/dashboard">
                <Dashboard>
                    <Home /> 
                </Dashboard>
            </PrivateRoute>
            <PrivateRoute path="/account/apikeys">
                <Dashboard>
                    <ApiKeys /> 
                </Dashboard>
            </PrivateRoute>
            <PrivateRoute path="/portfolio">
                <Dashboard>
                    <Portfolio /> 
                </Dashboard>
            </PrivateRoute>
            <PrivateRoute path="/orders">
                <Dashboard>
                    <Orders /> 
                </Dashboard>
            </PrivateRoute>
            <PrivateRoute path="/exchange">
                <Dashboard>
                    <Exchange /> 
                </Dashboard>
            </PrivateRoute>
        </Switch>
    );
}





