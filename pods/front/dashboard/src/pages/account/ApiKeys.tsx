
import React from 'react';
import {
  // Breadcrumbs as MuiBreadcrumbs,
  Card as MuiCard,
  CardContent as MuiCardContent,
  // Divider as MuiDivider,
  Paper as MuiPaper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  Button,
  Grid,
  CircularProgress
} from "@material-ui/core";
import styled from "styled-components/macro";

import { spacing } from "@material-ui/system";
import { useQuery, useMutation } from '@apollo/client';
import { LIST_API_KEYS, CREATE_NEW_KEY, DELETE_KEY, APIKey } from "./../../repositories/api-keys";
import { ConfirmationDialog } from "./../../components/Confirmation";

import {
  DeleteForever
} from "@material-ui/icons";

import moment from "moment";


const Card = styled(MuiCard)(spacing);

const CardContent = styled(MuiCardContent)(spacing);

// const Divider = styled(MuiDivider)(spacing);

// const Breadcrumbs = styled(MuiBreadcrumbs)(spacing);

const Paper = styled(MuiPaper)(spacing);


const SmallButton = styled(Button)`
  padding: 4px;
  min-width: 0;

  svg {
    width: 0.9em;
    height: 0.9em;
  }
`;



export function ApiKeys() {

  const { loading, error, data, refetch } = useQuery(LIST_API_KEYS);
  const [createNewApiKey] = useMutation(CREATE_NEW_KEY);
  const [deleteApiKey] = useMutation(DELETE_KEY);

  console.log("data", data);

  async function onClickCreateButton(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault()
    try {
      const result = await createNewApiKey();
      console.log("result", result);
      refetch();
    } catch (e) {
      console.log(e);
    }
  }

  async function onClickDeleteButton(e: React.MouseEvent<HTMLButtonElement, MouseEvent>, id: number) {

    const apiKey: APIKey = data?.listApiKeys?.[id];

    try {
      const result = await deleteApiKey({
        variables: {
          accessKey: apiKey?.accessKey
        }
      });
      console.log("result", result);
      refetch();
    } catch (e) {
      console.log(e);
    }
  }

  return (<>
    {
      loading &&  <Grid container justify="center"><CircularProgress /></Grid>
    }
    {
      error && "error"
    }
    {
      !loading && !error && data &&
      <>
        <Grid container justify="flex-end">
          <Button
            variant="contained"
            color="secondary"
            onClick={onClickCreateButton}
          >
            Create new API Key
        </Button>
        </Grid>
        <br />
        <Card mb={6}>
          <CardContent pb={1}>
            <Typography variant="h6" gutterBottom>
              API Keys
          </Typography>
            <Typography variant="body2" gutterBottom>
              Use those keys to connect to our Websocket API or REST API.
          </Typography>
          </CardContent>
          <Paper>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>access key</TableCell>
                  <TableCell>secret key</TableCell>
                  <TableCell align="right">created</TableCell>
                  <TableCell align="right">action</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {
                  data.listApiKeys?.map((row: APIKey, id: number) => {
                    return (<TableRow
                      key={id}
                    >
                      <TableCell component="th" scope="row">
                        {row.accessKey}
                      </TableCell>
                      <TableCell> {row.secretKey}</TableCell>
                      <TableCell align="right">{moment(row.created * 1000).format('MMMM Do YYYY, h:mm:ss a')}</TableCell>
                      <TableCell align="right">

                        <ConfirmationDialog
                          title="Deleting"
                          message="Are you sure you want to delete this API Key?"
                        >{
                            (confirm: any) => {
                              return (<SmallButton size="small" onClick={confirm((e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => onClickDeleteButton(e, id))}>
                                <DeleteForever />
                              </SmallButton>)
                            }
                          }</ConfirmationDialog>

                      </TableCell>
                    </TableRow>)
                  })
                }
              </TableBody>
            </Table>
          </Paper>
        </Card>
      </>
    }

  </>
  );
}