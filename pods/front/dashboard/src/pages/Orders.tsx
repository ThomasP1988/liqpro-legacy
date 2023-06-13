
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
    // Button,
    Grid,
    CircularProgress
} from "@material-ui/core";
import styled from "styled-components/macro";

import { spacing } from "@material-ui/system";
import { useQuery } from '@apollo/client';
import { ORDERS, Order } from "../repositories/orders";
import { CryptoIcon } from "./../components/CryptoIcon";

import PairsInfo from "./../pairInfos.json";
import moment from "moment";

const Card = styled(MuiCard)(spacing);

const CardContent = styled(MuiCardContent)(spacing);

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



export function Orders() {

    const { loading, error, data } = useQuery(ORDERS, {
        fetchPolicy: "no-cache"
    });
    console.log("data", data)
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
                <Card mb={6}>
                    <CardContent pb={1}>
                        <Typography variant="h6" gutterBottom>
                            Orders
                        </Typography>
                    </CardContent>
                    <Paper>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell>
                                    </TableCell>
                                    <TableCell>instrument</TableCell>
                                    <TableCell>action</TableCell>
                                    <TableCell>Quantity Asked</TableCell>
                                    <TableCell>Quantity Executed</TableCell>
                                    <TableCell>Price Per Unit</TableCell>
                                    <TableCell>Total Price</TableCell>
                                    <TableCell>status</TableCell>
                                    <TableCell align="right">Executed on</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {
                                    data.orders?.map((row: Order, id: number) => {
                                        return (<TableRow
                                            key={id}
                                        >
                                            <TableCell component="th" scope="row">
                                                <CryptoIcon symbol={(PairsInfo as any)[row.instrument]?.base} />
                                            </TableCell>
                                            <TableCell component="th" scope="row">
                                                
                                                {(PairsInfo as any)[row.instrument]?.display}
                                            </TableCell>
                                            <TableCell>{
                                                row.action ? <>SELL</> : <>BUY</>
                                            }</TableCell>
                                            <TableCell>{row.quantityAsked}</TableCell>
                                            <TableCell>{row.quantityExecuted}</TableCell>
                                            <TableCell>{row.pricePerUnit?.toFixed(2)}</TableCell>
                                            <TableCell>{row.totalPrice?.toFixed(2)}</TableCell>
                                            <TableCell>{row.status}</TableCell>
                                            <TableCell align="right">{moment(row.dateExecuted / 1000 / 1000).format('MMMM Do YYYY, h:mm:ss a')}</TableCell>
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