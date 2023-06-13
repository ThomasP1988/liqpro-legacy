
import React, { useState, useCallback } from 'react';
import useWebSocket from 'react-use-websocket';
import {
  // Breadcrumbs as MuiBreadcrumbs,
  // Card as MuiCard,
  // CardContent as MuiCardContent,
  // Divider as MuiDivider,
  Paper as MuiPaper,
  Button,
  IconButton,
  Chip,
  Typography,
  Grid,
  Dialog,
  DialogContent,
  DialogTitle,
  DialogActions,
  DialogContentText,
  TextField,
  LinearProgress
} from "@material-ui/core";
import styled from "styled-components/macro";

import AddIcon from '@material-ui/icons/Add';
import RemoveIcon from '@material-ui/icons/Remove';

import { spacing } from "@material-ui/system";
import { useQuery } from '@apollo/client';

import { PriceButton } from "./../../components/PriceButton";

import './../../price.css';
import { Instrument, LIST_INSTRUMENTS } from '../../repositories/instruments';
import { GET_TOKEN } from '../../repositories/tokens';
import { Subscribe, Subscription, Token, Order, OrderResult, EventPriceUpdate, OrdersState } from "./types";
import { CryptoIcon } from '../../components/CryptoIcon';
import PairsInfo from "./../../pairInfos.json";
import { TYPE_PRICE } from '../../components/PriceButton/types';
import { useStyles } from "./styles";
import { useFormFields } from './../../hooks/form';
// import { useInterval } from '../../hooks/helpers';

// const Card = styled(MuiCard)(spacing);

// const CardContent = styled(MuiCardContent)(spacing);

// const Divider = styled(MuiDivider)(spacing);

// const Breadcrumbs = styled(MuiBreadcrumbs)(spacing);

const Paper = styled(MuiPaper)(spacing);


// const SmallButton = styled(Button)`
//   padding: 4px;
//   min-width: 0;

//   svg {
//     width: 0.9em;
//     height: 0.9em;
//   }
// `;



export function Exchange() {

  const { data: dataInstr } = useQuery(LIST_INSTRUMENTS);
  const { refetch: refetchToken } = useQuery(GET_TOKEN);
  // const [socket, setSocket] = useState<WebSocket>();
  // const socket = useRef<WebSocket | null>(null)
  const [prices, setPrices] = useState<{ [id: string]: number }>({});
  const [subscriptions, setSubscriptions] = useState<{ [id: string]: Subscription }>({});
  const [openChangeLevel, setOpenChangeLevel] = useState<boolean>(false);
  const [openOrderDialog, setOpenOrderDialog] = useState<boolean>(false);
  const [orderInProgress, setOrderInProgress] = useState<Order | null>();

  const [ordersState, setOrdersState] = useState<OrdersState>({});

  const [levelChangePair, setLevelChangePair] = useState<string>("");
  const [fields, handleFieldChange, setFields] = useFormFields({
    level: "",
  });
  const classes = useStyles();

  const getSocketUrl = useCallback(async (): Promise<string> => {
    let token: Token;
    try {
      const response = await refetchToken();
      token = response.data.getToken
    } catch (e) {
      console.log(e);
      return Promise.reject(e);
    }

    const params: URLSearchParams = new URLSearchParams();
    params.set("Token", token.value);
    const url: string = "ws://127.0.0.1:8080/ws?" + params.toString();

    return Promise.resolve(url);
  }, [refetchToken]);

  const {
    sendJsonMessage
  } = useWebSocket(getSocketUrl, {
    onOpen: () => console.log('opened'),
    //Will attempt to reconnect on all close events, such as server shutting down
    // shouldReconnect: (closeEvent) => true,
    onMessage,
  });


  function onMessage(message: WebSocketEventMap['message']) {
    console.log("message", message);
    console.log("subscriptions", subscriptions);

    console.log("typeof message", typeof message);
    const data = JSON.parse(message.data);

    // "ask":57135.481649167144,"bid":57131.51926357042,"market":"btcusd","level":1}
    if (data.event === "price") {
      UpdatePrice(data as EventPriceUpdate)
    } else if (data.event === "orderResult") {


      const newOrdersState = Object.assign({}, ordersState);
      newOrdersState[(data as OrderResult).clientId].result = data as OrderResult;
      setOrdersState(newOrdersState);
    }

  };

  const unsubscribe = (pair: string, level: number) => {
    scribe("unsubscribe", pair, level);
    const subscription: Subscription = {
      pair,
      level,
      subscribed: false
    };

    const newSubscriptions = Object.assign({}, subscriptions);
    newSubscriptions[pair] = subscription;
    setSubscriptions(newSubscriptions);
  }

  const subscribe = (pair: string, level: number) => {
    let serverLevel: number = 0;
    if (!dataInstr?.listIntruments?.[pair]?.[level]) {
      for (let lvl in dataInstr?.listIntruments?.[pair]) {
        if (level > (+lvl)) {
          serverLevel = +lvl;
        }
      }
    }

    scribe("subscribe", pair, serverLevel ? serverLevel : level);

    const subscription: Subscription = {
      pair,
      level,
      subscribed: true
    };

    if (serverLevel) {
      subscription.serverLevel = serverLevel;
    }

    const newSubscriptions = Object.assign({}, subscriptions);
    newSubscriptions[pair] = subscription;
    setSubscriptions(newSubscriptions);
  }

  const order = (instrument: string, quantity: number, side: string) => {

    const clientId = String(Math.floor(Math.random() * 100000));

    const orderPayload: Order = {
      event: "order", // should always be first argument
      side,
      instrument,
      quantity,
      clientId,
    }
    sendJsonMessage(orderPayload);
    setOrderInProgress(orderPayload);

    const newOrdersState = Object.assign({}, ordersState);
    newOrdersState[clientId] = {
      order: orderPayload
    };
    setOrdersState(newOrdersState);

    setOpenOrderDialog(true);

  }

  const scribe = (event: string, pair: string, level: number) => {
    const subscribePayload: Subscribe = {
      event: event,
      instrument: pair,
      levels: [+level]
    }
    sendJsonMessage(subscribePayload);
  }

  const toggleSubscription = (e: React.MouseEventHandler<HTMLDivElement>, pair: string): void => {

    let subscription: Subscription = subscriptions[pair];
    if (!subscription || !subscription.subscribed) {
      // subscribe
      subscribe(pair, 1)
    } else {
      // unsubscribe
      unsubscribe(pair, subscription.level)
    }

  }

  const changeLevel = (): void => {
    let subscription: Subscription = subscriptions[levelChangePair];
    unsubscribe(levelChangePair, subscription.level);
    subscribe(levelChangePair, fields.level);
    handleClose();
  }

  const handleClose = () => {
    setOpenChangeLevel(false);
    setLevelChangePair("");
  }

  const handleCloseOrder = () => {
    setOpenOrderDialog(false);
    setOrderInProgress(null);
  }

  const triggerChangeLevelDialog = (pair: string) => {
    setOpenChangeLevel(true);
    setLevelChangePair(pair);
  }

  const reduceFloatDecimal = (value: number, decimal: number): number => {
    return +value.toFixed(decimal);
  }

  const UpdatePrice = (event: EventPriceUpdate) => {
    console.log("subscriptions", subscriptions);
    const subLevel: number = subscriptions[event.market].level;
    const localLevel: number = event.level;

    const ratio: number = subLevel !== localLevel ? localLevel / subLevel : 1;
    const updatedPrices = Object.assign({}, prices);

    updatedPrices[`${event.market}0`] = reduceFloatDecimal(event.bid, (PairsInfo as any)[event.market].decimalQuote * ratio);
    updatedPrices[`${event.market}1`] = reduceFloatDecimal(event.ask, (PairsInfo as any)[event.market].decimalQuote * ratio);
    setPrices(updatedPrices);
  }

  console.log("changeLevel", changeLevel)

  // useInterval(() => {
  //   setPrices({
  //     "btcusd0": Math.floor(Math.random() * 100000),
  //     "btcusd1": Math.floor(Math.random() * 100000),
  //   })
  // }, 1000);

  return (<>
    {
      dataInstr?.listIntruments?.map((item: Instrument) => {
        return <><Chip
          icon={subscriptions[item.pair]?.subscribed ? <RemoveIcon /> : <AddIcon />}
          label={(PairsInfo as any)[item.pair]?.display}
          onClick={(e: any) => toggleSubscription(e, item.pair)}
          color={subscriptions[item.pair]?.subscribed ? "secondary" : "default"}
          variant="outlined"
          className={classes.pairTopToggle}
        />
        </>
      })
    }
    {/* <Button onClick={() => order("", 1, "0")} >test</Button> */}
    {
      subscriptions && Object.values(subscriptions).filter((i: Subscription) => i.subscribed).map((item: Subscription, key: number) => {
        return <Paper key={key} className={classes.pricesContainer}>
          <Grid container>
            <Grid item xs={9}>
              <CryptoIcon className={classes.currencyIcon} symbol={(PairsInfo as any)[item.pair]?.base} />
              <Typography display="inline" className={classes.pairContainer}>
                {(PairsInfo as any)[item.pair]?.display}
              </Typography>
              <PriceButton className={classes.price} type={TYPE_PRICE.BUY} newPrice={prices[`${item.pair}0`]} onClick={() => order(item.pair, Number(item.level), "0")} />
              <PriceButton className={classes.price} type={TYPE_PRICE.SELL} newPrice={prices[`${item.pair}1`]} onClick={() => order(item.pair,  Number(item.level), "1")} />
              <Typography display="inline">
                price for {item.level} {(PairsInfo as any)[item.pair]?.baseTitle}
              </Typography>
            </Grid>
            <Grid item xs={3} container alignItems="center" alignContent="center" justify="flex-end">
              <Button variant="outlined" onClick={() => triggerChangeLevelDialog(item.pair)}>Change level</Button>
              <IconButton aria-label="delete" color="primary" size="small" className={classes.removeIcon}>
                <RemoveIcon onClick={() => unsubscribe(item.pair, item.level)} />
              </IconButton>
            </Grid>
          </Grid>
        </Paper>
      })
    }

    <Dialog open={openChangeLevel} onClose={handleCloseOrder} aria-labelledby="form-dialog-title" fullWidth>
      <DialogTitle id="form-dialog-title">Change level</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Please indicate the new level for this pair.
          </DialogContentText>
        <TextField
          autoFocus
          margin="dense"
          id="level"
          label="Level"
          type="number"
          InputProps={{ inputProps: { min: 0 } }}
          onChange={handleFieldChange}
          value={fields.level}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          Cancel
          </Button>
        <Button onClick={changeLevel} color="primary">
          Change
          </Button>
      </DialogActions>
    </Dialog>
    <Dialog open={openOrderDialog} onClose={handleClose} aria-labelledby="form-dialog-title" fullWidth>
      <DialogTitle id="form-dialog-title">Placing order {
        orderInProgress ? "on " + (PairsInfo as any)[orderInProgress?.instrument]?.display : ""
      }</DialogTitle>
      <DialogContent>

        {
          orderInProgress && ordersState[orderInProgress.clientId]?.result ? <DialogContentText>
            <CryptoIcon className={classes.currencyIcon} symbol={(PairsInfo as any)[orderInProgress?.instrument]?.base} />
            Successfully {
              orderInProgress?.side === "1" ? "Sold" : "Bought"
            } {
              ordersState[orderInProgress.clientId]?.result?.quantity
            } {
              (PairsInfo as any)[(ordersState[orderInProgress.clientId]?.result as OrderResult).instrument]?.baseTitle
            } for {
              ordersState[orderInProgress.clientId]?.result?.price
            } {
              (PairsInfo as any)[(ordersState[orderInProgress.clientId]?.result as OrderResult).instrument]?.quoteTile
            } ({
              ordersState[orderInProgress.clientId]?.result?.pricePerUnit
            } per unit)
          </DialogContentText> : <DialogContentText>
            {
              orderInProgress && <>
                <CryptoIcon className={classes.currencyIcon} symbol={(PairsInfo as any)[orderInProgress?.instrument]?.base} />
                {
                  orderInProgress?.side === "1" ? "Selling" : "Buying"
                } {
                  orderInProgress?.quantity
                } {
                  (PairsInfo as any)[orderInProgress?.instrument]?.baseTitle
                }
              </>
            }
            <LinearProgress />
          </DialogContentText>
        }
      </DialogContent>
      {
        orderInProgress && ordersState[orderInProgress.clientId]?.result && <DialogActions>
          <Button onClick={handleCloseOrder} color="primary">
            Close
            </Button>

        </DialogActions>
      }
    </Dialog>
  </>
  );
}