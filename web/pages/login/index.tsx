import {
  Dispatch,
  FormEvent,
  SetStateAction,
  useEffect,
  useState,
} from "react";
import { toast } from "sonner";
import { Link, useLocation, useSearch } from "wouter";
import { Button } from "~/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { simpleFetch } from "~/hooks/useFetch";
import { apiCodesMap, genericErrorsMap } from "~/lib/statuscodes";
import { ApiResponse } from "~/types/api";
import { User } from "~/types/user";

export function LoginPage({
  setAuthed,
}: {
  setAuthed: Dispatch<SetStateAction<boolean>>;
}) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [, navigate] = useLocation();
  const params = useSearch();

  useEffect(() => {
    const errcode = new URLSearchParams(params).get("errcode");

    if (errcode) {
      requestAnimationFrame(() =>
        toast.error(apiCodesMap[Number(errcode)], {
          position: "top-center",
        }),
      );
    }
  }, [params]);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();

    // Generic checks
    if (email === "" || password === "")
      return toast.error(genericErrorsMap["emptyFormField"]);

    if (password.length < 8)
      return toast.error(genericErrorsMap["shortPassword"]);

    const body = new URLSearchParams({
      email,
      password,
    });

    const data = await simpleFetch<ApiResponse<User>>(
      "/v1/auth/login",
      navigate,
      {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: body,
      },
    );

    // if Error do nothing and notify
    if (data.err_code) {
      return toast.error(apiCodesMap[data.err_code], {
        position: "top-center",
      });
    }

    // Success go home
    setAuthed(true);
    navigate("/");
  }

  return (
    <div className="flex h-screen w-screen flex-col items-center justify-center gap-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Login</CardTitle>
          <CardDescription>
            Enter your email below to login to your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  placeholder="m@example.com"
                  onChange={(e) => setEmail(e.target.value)}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>
              <Button type="submit" className="w-full">
                Login
              </Button>
              <Button
                variant="outline"
                className="w-full"
                type="button"
                onClick={() => {
                  window.open("/v1/auth/login/google", "_self");
                }}
              >
                Login with Google
              </Button>
            </div>
            <div className="mt-4 text-center text-sm">
              Don&apos;t have an account?{" "}
              <Link to="/signup" className="underline underline-offset-4">
                Sign up
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
