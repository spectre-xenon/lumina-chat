import { Switch, Route, useLocation } from "wouter";
import { useAuth } from "~/hooks/useAuth";
import { NotFound } from "./pages/404";
import { LoginPage } from "./pages/login";
import { Toaster } from "./components/ui/sonner";
import { SignupPage } from "./pages/signup";
import { Home } from "./pages/home";

const allowedLocations: { [key: string]: boolean } = {
  "/login": true,
  "/signup": true,
};

export function Router() {
  const [location, navigate] = useLocation();
  const [authed, setAuthed, loading] = useAuth();

  if (loading) {
    return "loading";
  }

  if (!authed && !allowedLocations[location]) {
    navigate("/login", { replace: true });
  }

  return (
    <>
      <Switch>
        {/* Auth */}
        <Route path="/login">
          <LoginPage setAuthed={setAuthed} />
        </Route>
        <Route path="/signup">
          <SignupPage setAuthed={setAuthed} />
        </Route>

        {/* index */}
        <Route path="/:chatId?">
          <Home />
        </Route>

        {/* Default route in a switch */}
        <Route>
          <NotFound />
        </Route>
      </Switch>

      <Toaster richColors />
    </>
  );
}
