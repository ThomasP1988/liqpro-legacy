import { gql } from '@apollo/client';

export const PORTFOLIO = gql`
  query portfolio {
    portfolio {
      id
      currency
      userId
      total
      lastModified
    }
  }
`;

export type Hold = {
    id: string,
    currency: string,
    userId: string,
    total: number,
    lastModified: number
  }