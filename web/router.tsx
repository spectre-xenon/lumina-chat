import { Switch, Route, useLocation } from "wouter";
import { useAuth } from "~/hooks/useAuth";
import { NotFound } from "./pages/404";

const allowedLocations: { [key: string]: boolean } = {
  "/login": true,
  "/signup": true,
};

const notAllowedLocations: { [key: string]: boolean } = {
  "/": true,
};

export function Router() {
  const [location, navigate] = useLocation();
  const [authed, loading] = useAuth();

  if (loading) {
    return <h1>loading</h1>;
  }

  if (!authed && !allowedLocations[location] && notAllowedLocations[location]) {
    navigate("/login");
  }

  return (
    <Switch>
      <Route path="/">home</Route>
      <Route path="/login">login</Route>
      <Route path="/signup" />

      {/* Default route in a switch */}
      <Route component={NotFound}></Route>
    </Switch>
  );
}
