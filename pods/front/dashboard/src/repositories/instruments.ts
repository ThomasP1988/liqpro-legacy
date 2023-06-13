import { gql } from '@apollo/client';

export const LIST_INSTRUMENTS = gql`
  query listIntruments {
    listIntruments {
      pair,
      levels
    }
  }
`;

export type Instrument = {
  pair: string,
  levels: number[]
}