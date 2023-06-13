// import SuperTokensRequest from 'supertokens-website';
import {
    ApolloClient,
    InMemoryCache,
    HttpLink,
    NormalizedCacheObject,
    DefaultOptions,
  } from '@apollo/client';

  const defaultOptions: DefaultOptions = {
    watchQuery: {
      fetchPolicy: 'no-cache',
      errorPolicy: 'ignore',
    },
    query: {
      fetchPolicy: 'no-cache',
      errorPolicy: 'all',
    },
  }

export const Client = new ApolloClient<NormalizedCacheObject>({
    link: new HttpLink({
        uri: "https://localhost:8081/protected/graphql",  // change this depending on your path
    }),
    cache: new InMemoryCache(),  // change this depending on your preference
    defaultOptions
    // ... other params
});
