
import React, { useEffect, useCallback } from 'react';
import axios from "axios";

export function Home() {

    const getProfile = useCallback(async () => {
        try {

            const response = await axios.get("https://localhost:8081/protected/me")

            console.log("response", response);
        } catch (e) {
            console.log(e);
        }
    }, [])

    useEffect(() => {
        if (getProfile) {
            getProfile();
        }
    }, [getProfile]);

    return (
        <>Welcome</>
    );
}