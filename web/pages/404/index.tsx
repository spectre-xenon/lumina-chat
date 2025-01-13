import { Button } from "~/components/ui/button";
import { FourZeroFour } from "./404";
import { Sparkles } from "lucide-react";
import { useLocation } from "wouter";

export function NotFound() {
  const [, navigate] = useLocation();

  return (
    <div className="flex h-screen w-screen flex-col items-center justify-center gap-16">
      <FourZeroFour width={450} height={250} />
      <div className="flex items-center gap-2">
        <p>Void ahead. Reroute to something </p>
        <Button variant="secondary" onClick={() => navigate("/")}>
          <Sparkles />
        </Button>
        <p>stellar!</p>
      </div>
    </div>
  );
}
