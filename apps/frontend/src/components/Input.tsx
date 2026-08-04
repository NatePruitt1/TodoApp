import type React from "react";
import { forwardRef, useImperativeHandle, useState } from "react";

import "./Input.css"

export type InputHandle = {
    setError: (message: string) => void;
};

export const Input = forwardRef<InputHandle, React.ComponentProps<"input">>((props, ref) => {
    const [error, setError] = useState("");

    useImperativeHandle(ref, () => ({
        setError,
    }));

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setError("");
        props.onChange?.(e);
    };

    return (
        <div className="form-input-container">
            <input 
                className="form-input"
                id={props.id}
                value={props.value}
                onChange={handleChange}
                type={props.type}
                name={props.name}
                placeholder={props.placeholder}
            />
        <p className="form-input-error">{error && <><b>(!)</b> {error}</>}</p>
        </div>
    )
});

Input.displayName = "Input";