import { useState, useEffect } from "react";
import { useLocation } from "wouter";

export function useAuth() {
  const [location] = useLocation();
  const [authed, setAuthed] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function checkAuthed() {
      const res = await fetch("/v1/auth");

      // Unauthorized
      if (res.status == 401) {
        setAuthed(false);
        setLoading(false);
        return;
      }

      setAuthed(true);
      setLoading(false);
    }

    checkAuthed();
  }, [location]);

  return [authed, loading];
}
