import { gql } from '@apollo/client';

export const ORDERS = gql`
  query orders {
    orders {
      id
      action
      instrument
      quantityAsked
      quantityExecuted
      pricePerUnit
      totalPrice
      dateExecuted
      status
    }
  }
`;

export enum STATUS {
    Success = "success",
    Failed = "failed",
}

export type Order = {
    id: string;
    action: number;
    instrument: string;
    quantityAsked: number;
    quantityExecuted: number;
    pricePerUnit: number;
    totalPrice: number;
    dateExecuted: number;
    status: STATUS;
}