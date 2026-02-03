import { useState } from "react";

export const useLode = () => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    return {
        loading,
        error,
        success,
        setLoading,
        setError,
        setSuccess,
    }
}