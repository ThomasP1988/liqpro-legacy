import { useState, ChangeEvent } from "react";
import { useEffect, useRef } from "react";

export function useFormFields(initialState: any) {
    const [fields, setValues] = useState(initialState);

    return [
        fields,
        (event: ChangeEvent<HTMLInputElement>) => {
            setValues({
                ...fields,
                [event.target.id]: event.target.value,
            });
        },
        setValues,
    ];
}


export function useCheckboxes(initialState: any) {
    const [fields, setValues] = useState(initialState);

    return [
        fields,
        (event: ChangeEvent<HTMLInputElement>) => {
            setValues({
                ...fields,
                [event.target.id]: event.target.checked,
            });
        },
        setValues,
    ];
}

export const usePrevious = <T>(value: T): T | undefined => {
    const ref = useRef<T>();
    useEffect(() => {
      ref.current = value;
    });
    return ref.current;
  };

