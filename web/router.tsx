import { Switch, Route, useLocation } from "wouter";
import { useAuth } from "~/hooks/useAuth";
import { NotFound } from "./pages/404";
import { LoginPage } from "./pages/login";
import { Toaster } from "./components/ui/sonner";
import { SignupPage } from "./pages/signup";

const allowedLocations: { [key: string]: boolean } = {
  "/login": true,
  "/signup": true,
};

const notAllowedLocations: { [key: string]: boolean } = {
  "/": true,
};

export function Router() {
  const [location, navigate] = useLocation();
  const [authed, setAuthed, loading] = useAuth();

  if (loading) {
    return "loading";
  }

  if (!authed && !allowedLocations[location] && notAllowedLocations[location]) {
    navigate("/login", { replace: true });
  }

  return (
    <>
      <Toaster richColors />

      <Switch>
        <Route path="/">home</Route>
        <Route path="/login">{() => <LoginPage setAuthed={setAuthed} />}</Route>
        <Route path="/signup">
          {() => <SignupPage setAuthed={setAuthed} />}
        </Route>

        {/* Default route in a switch */}
        <Route component={NotFound}></Route>
      </Switch>
    </>
  );
}
