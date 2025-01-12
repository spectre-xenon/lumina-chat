import { Switch, Route, useLocation } from "wouter";
import { useAuth } from "~/hooks/useAuth";

const allowedLocations: { [key: string]: boolean } = {
  "/login": true,
  "/signup": true,
};

export function Router() {
  const [location, navigate] = useLocation();
  const [authed, loading] = useAuth();

  if (loading) {
    return <h1>loading</h1>;
  }

  if (!authed && !allowedLocations[location]) {
    navigate("/login");
  }

  return (
    <>
      <Switch>
        <Route path="/">home</Route>
        <Route path="/login">login</Route>
        <Route path="/signup" />

        {/* Default route in a switch */}
        <Route>404: No such page!</Route>
      </Switch>
    </>
  );
}
