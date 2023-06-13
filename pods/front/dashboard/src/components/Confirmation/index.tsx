import React, { ReactElement, useState } from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';

interface IConfirmationProps {
    message: string | ReactElement;
    title?: string;
    buttonText?: string | null;
    children: any;
}

const defaultCallbackEvent: Function = (): void => { };
let callbackFc: Function = defaultCallbackEvent;

export function ConfirmationDialog({ message = "", title = "", buttonText, children }: IConfirmationProps) {
    const [isOpen, setIsOpen] = useState(false);

    const show = (callback: any) => (event: any) => {
        event.preventDefault()

        event = {
            ...event,
            target: { ...event.target, value: event.target.value }
        }

        callbackFc = (): void => callback(event);
        setIsOpen(true);
    }

    const onConfirm = () => {
        callbackFc()
        handleClose();
    }

    const handleClose = () => {
        setIsOpen(false);
        callbackFc = defaultCallbackEvent;
    };

    return (
        <>
            {children(show)}
            <Dialog
                open={isOpen}
                keepMounted
                fullWidth
                onClose={handleClose}
                aria-labelledby="alert-dialog-slide-title"
                aria-describedby="alert-dialog-slide-description"
            >
                <DialogTitle id="alert-dialog-title" style={{ justifyContent: "center" }}>{title || "Confirmation"}</DialogTitle>
                <DialogContent>
                    <Grid container direction="row" justify="center">
                        {message}
                    </Grid>
                    <Grid container direction="row" justify="center" spacing={4}>
                        <Grid item xs={5} container justify="flex-end">
                            <Button onClick={handleClose} autoFocus variant="contained">
                                cancel
                            </Button>
                        </Grid>
                        <Grid item xs={5} container justify="flex-start">
                            <Button onClick={onConfirm} color="secondary" variant="contained">
                                {
                                    buttonText || "delete"
                                }
                            </Button>
                        </Grid>
                    </Grid>
                </DialogContent>
            </Dialog>
        </>
    )
}