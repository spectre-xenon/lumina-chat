import { Link } from "wouter";
import { Chat } from "~/types/chat";
import { Skeleton } from "./ui/skeleton";

const MAX_CHARS = 25;

export function ChatLink({ chat }: { chat: Chat }) {
  const senderLength = chat["last_message"].sender.length;
  const lastMessage =
    chat["last_message"].sender +
    ": " +
    chat["last_message"].content.substring(0, MAX_CHARS - senderLength) +
    "...";

  return (
    <Link
      to={`/${chat["id"]}`}
      className="flex items-center justify-center gap-3 rounded-lg bg-background p-4 hover:bg-muted"
    >
      <img src={chat["picture"]} className="h-14 w-14 rounded-full" />

      <div className="flex flex-grow flex-col gap-2">
        <div className="flex justify-between">
          <p className="font-bold">{chat["name"]}</p>
          <p className="text-muted-foreground">09:18 pm</p>
        </div>
        <div>
          <p className="overflow-hidden whitespace-nowrap text-muted-foreground">
            {lastMessage}
          </p>
        </div>
      </div>
    </Link>
  );
}

export function ChatLinkSkeleton() {
  return (
    <div className="flex items-center justify-center gap-3 rounded-lg bg-background p-4 hover:bg-muted">
      <Skeleton className="h-14 w-14 rounded-full" />

      <div className="flex flex-grow flex-col gap-4">
        <Skeleton className="h-4 w-full" />
        <div>
          <Skeleton className="h-4 w-full" />
        </div>
      </div>
    </div>
  );
}
