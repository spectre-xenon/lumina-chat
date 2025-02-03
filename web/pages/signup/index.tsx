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

export function SignupPage({
  setAuthed,
}: {
  setAuthed: Dispatch<SetStateAction<boolean>>;
}) {
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
  });

  function updateFormData(key: keyof typeof formData, value: string) {
    setFormData({ ...formData, [key]: value });
  }

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
    if (
      formData["username"] === "" ||
      formData["email"] === "" ||
      formData["password"] === ""
    )
      return toast.error(genericErrorsMap["emptyFormField"]);

    if (formData["username"].length < 3)
      return toast.error(genericErrorsMap["shortUsername"]);

    if (formData["username"].length > 27)
      return toast.error(genericErrorsMap["shortUsername"]);

    if (formData["password"].length < 8)
      return toast.error(genericErrorsMap["shortPassword"]);

    const body = new URLSearchParams(formData);

    const data = await simpleFetch<ApiResponse<User>>(
      "/v1/auth/signup",
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
          <CardTitle className="text-2xl">Signup</CardTitle>
          <CardDescription>
            Enter your email below to create a new account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-2">
                <Label htmlFor="username">Username*</Label>
                <Input
                  id="username"
                  name="username"
                  type="username"
                  placeholder="jhin"
                  onChange={(e) => updateFormData("username", e.target.value)}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="email">Email*</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  placeholder="m@example.com"
                  onChange={(e) => updateFormData("email", e.target.value)}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="password">Password*</Label>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  onChange={(e) => updateFormData("password", e.target.value)}
                />
              </div>
              <Button type="submit" className="w-full">
                Signup
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
              Already have an account?{" "}
              <Link to="/login" className="underline underline-offset-4">
                login
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
