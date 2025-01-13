import { useEffect, useState } from "react";
import { useLocation } from "wouter";

export function useFetch<T>(url: string, opts?: RequestInit) {
  const [data, setData] = useState<T | null>();
  const [loading, setLoading] = useState(true);

  const [, navigate] = useLocation();

  useEffect(() => {
    async function fetchData() {
      const res = await fetch(url, opts);

      if (res.status === 401) {
        return navigate("/login");
      }

      const data = (await res.json()) as T;

      setData(data);
      setLoading(false);
    }

    fetchData();
  }, [url, opts, navigate]);

  return { data, loading };
}

export async function simpleFetch<T>(
  url: string,
  navigate: (
    to: string | URL,
    options?: {
      replace?: boolean;
      state?: unknown;
    },
  ) => void,
  opts?: RequestInit,
): Promise<T> {
  const res = await fetch(url, opts);

  if (res.status === 401) {
    navigate("/login");
  }

  const data = (await res.json()) as T;

  return data;
}
