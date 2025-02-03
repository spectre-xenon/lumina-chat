import { Link } from "wouter";
import { Chat } from "~/types/chat";
import { Skeleton } from "./ui/skeleton";
import { formatDate } from "~/lib/time";
import { Sparkles } from "lucide-react";

const MAX_CHARS = 25;

export function ChatLink({ chat }: { chat: Chat }) {
  const lastMessage = chat["message"].content.substring(0, MAX_CHARS) + "...";

  return (
    <Link
      to={`/${chat["id"]}`}
      className="flex items-center justify-center gap-3 rounded-lg bg-background p-4 hover:bg-muted"
    >
      {chat["picture"] ? (
        <img
          src={chat["picture"]}
          className="h-14 w-14 rounded-full border object-cover object-center"
        />
      ) : (
        <div className="flex h-14 w-14 items-center justify-center rounded-full border">
          <Sparkles />
        </div>
      )}

      <div className="flex flex-grow flex-col gap-2">
        <div className="flex justify-between">
          <p className="font-bold">{chat["name"]}</p>
          <p className="text-muted-foreground">
            {formatDate(chat["message"].sent_at)}
          </p>
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
