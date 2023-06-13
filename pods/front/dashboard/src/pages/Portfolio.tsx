
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
import { PORTFOLIO, Hold } from "../repositories/portfolio";
import { CryptoIcon } from "./../components/CryptoIcon";
import moment from "moment";
// import {
//     DeleteForever
// } from "@material-ui/icons";

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



export function Portfolio() {

    const { loading, error, data } = useQuery(PORTFOLIO, {
        fetchPolicy: "no-cache"
    });


    return (<>
        {
            loading && <Grid container justify="center"><CircularProgress /></Grid>
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
                            Portfolio
                        </Typography>
                    </CardContent>
                    <Paper>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell>
                                    </TableCell>
                                    <TableCell>currency</TableCell>
                                    <TableCell>total</TableCell>
                                    <TableCell align="right">last modified</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {
                                    data.portfolio?.map((row: Hold, id: number) => {
                                        return (<TableRow
                                            key={id}
                                        >
                                            <TableCell component="th" scope="row">
                                                <CryptoIcon symbol={row.currency} />
                                            </TableCell>
                                            <TableCell>{row.currency}</TableCell>
                                            <TableCell>{
                                                row.currency === "usd" || row.currency === "eur" ?
                                                    row.total.toFixed(2) : row.total}</TableCell>
                                            <TableCell align="right"> {moment(row.lastModified * 1000).format('MMMM Do YYYY, h:mm:ss a')}</TableCell>
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