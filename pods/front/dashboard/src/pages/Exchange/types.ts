export type Subscription = {
    subscribed: boolean;
    level: number;
    serverLevel?: number;
    pair: string;
}
export type Order = {
    clientId: string;
    event: string;
    side: string;
    quantity: number;
    instrument: string;
}

export type OrderResult = {
    clientId: string;
    event: string;
    side: string;
    quantity: number;
    instrument: string;
    pricePerUnit: string;
    price: string;
}

export type Token = {
    value: string;
}

export type Subscribe = {
    event: string,
    instrument: string,
    levels: number[]
}

export type EventPriceUpdate = {
    event: string;
    ask: number;
    bid: number;
    level: number;
    market: string;
}

export type OrderState = {
    order: Order,
    result?: OrderResult
}

export type OrdersState = {
    [key: string]: OrderState
}
