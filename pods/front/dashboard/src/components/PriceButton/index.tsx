import { useEffect, useState, useCallback } from "react";
// import { useInterval } from '../../hooks/helpers';

import {
    Button, ButtonGroup
} from "@material-ui/core";
import _ from "lodash";

import "./style.css";
import { TYPE_PRICE } from "./types";

type Args = {
    newPrice: number,
    onClick: () => void,
    type: TYPE_PRICE,
    className?: string
}

const positiveChange: string = "positive-change";
const negativeChange: string = "negative-change";

export function PriceButton({ newPrice, onClick, type, className }: Args) {
    const [price, setPrice] = useState<number>(0);
    const [extraClass, setExtraClass] = useState<string>("");

    const resetClass = useCallback(_.throttle(() => {
        console.log("resetClass")
        setExtraClass("");
    }, 3000), [setExtraClass])

    useEffect(() => {
        if (newPrice < price) {
            setExtraClass(positiveChange);
            setPrice(newPrice);
            resetClass();
        } else if (newPrice > price) {
            setExtraClass(negativeChange);
            setPrice(newPrice);
            resetClass();
        }
    }, [newPrice, price, setPrice])

    return <ButtonGroup variant="contained" color="primary" className={(className ? className : "") + " price " + extraClass}>
        {
            type === TYPE_PRICE.BUY && <Button onClick={onClick} 
            disabled={!Boolean(newPrice)}
            >Buy</Button>
        }
        <Button className="priceButton" disabled>{price}</Button>
        {
            type === TYPE_PRICE.SELL && <Button onClick={onClick} disabled={!Boolean(newPrice)}>Sell</Button>
        }
    </ButtonGroup>
}
