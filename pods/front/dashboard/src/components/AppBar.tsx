import * as React from "react";
import styled, { withTheme } from "styled-components/macro";
// import { darken } from "polished";
import { Power } from "react-feather";
import axios from "axios";
import { useHistory } from "react-router-dom";

import {
  Grid,
  Hidden,
  // InputBase,
  AppBar as MuiAppBar,
  IconButton as MuiIconButton,
  Toolbar,
  Tooltip,
} from "@material-ui/core";

import { Menu as MenuIcon } from "@material-ui/icons";
import { attemptRefreshingSession  } from 'supertokens-auth-react/recipe/session';


const AppBar = styled(MuiAppBar)`
  background: ${(props: any) => props.theme.header.background};
  color: ${(props: any) => props.theme.header.color};
`;

const IconButton = styled(MuiIconButton)`
  svg {
    width: 22px;
    height: 22px;
  }
`;

// const Search = styled.div`
//   border-radius: 2px;
//   background-color: ${(props: any) => props.theme.header.background};
//   display: none;
//   position: relative;
//   width: 100%;

//   &:hover {
//     background-color: ${(props: any) => darken(0.05, props.theme.header.background)};
//   }

//   ${(props: any) => props.theme.breakpoints.up("md")} {
//     display: block;
//   }
// `;

// const SearchIconWrapper = styled.div`
//   width: 50px;
//   height: 100%;
//   position: absolute;
//   pointer-events: none;
//   display: flex;
//   align-items: center;
//   justify-content: center;

//   svg {
//     width: 22px;
//     height: 22px;
//   }
// `;

// const Input = styled(InputBase)`
//   color: inherit;
//   width: 100%;

//   > input {
//     color: ${(props: any) => props.theme.header.search.color};
//     padding-top: ${(props: any) => props.theme.spacing(2.5)}px;
//     padding-right: ${(props: any) => props.theme.spacing(2.5)}px;
//     padding-bottom: ${(props: any) => props.theme.spacing(2.5)}px;
//     padding-left: ${(props: any) => props.theme.spacing(12)}px;
//     width: 160px;
//   }
// `;

type AppBarProps = {
  theme: {};
  onDrawerToggle: React.MouseEventHandler<HTMLElement>;
};

const AppBarComponent: React.FC<AppBarProps> = ({ onDrawerToggle }) => {
  const history = useHistory();

  async function logOut(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();
    try {

      const response = await axios.get("https://localhost:8081/logout");

      await attemptRefreshingSession()

      history.push("/");
      console.log("response", response);
    } catch (e) {
      console.log(e);
    }
  }


  return (
    <React.Fragment>
      <AppBar position="sticky" elevation={0}>
        <Toolbar>
          <Grid container alignItems="center">
            <Hidden mdUp>
              <Grid item>
                <IconButton
                  color="inherit"
                  aria-label="Open drawer"
                  onClick={onDrawerToggle}
                >
                  <MenuIcon />
                </IconButton>
              </Grid>
            </Hidden>
            <Grid item>
            </Grid>
            <Grid item xs />
            <Grid item>
              <Tooltip title="Log out">
                <IconButton
                  onClick={logOut}
                  color="inherit"
                >
                  <Power />
                </IconButton>
              </Tooltip>
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
    </React.Fragment>
  )
};

export default withTheme(AppBarComponent);
