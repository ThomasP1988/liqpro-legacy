import { gql } from '@apollo/client';


export const LIST_API_KEYS = gql`
  query listApiKeys {
    listApiKeys {
      id,
      accessKey,
      secretKey,
      created
    }
  }
`;

export const CREATE_NEW_KEY = gql`
  mutation createApiKey {
    createApiKey {
      id,
      accessKey,
      secretKey,
      created
    }
  }
`;

export const DELETE_KEY = gql`
  mutation deleteApiKey($accessKey: String!) {
    deleteApiKey(accessKey: $accessKey) {
      success
    }
  }
`;

export type APIKey = {
  id: string,
  accessKey: string,
  userId: string,
  secretKey: string,
  created: number
}