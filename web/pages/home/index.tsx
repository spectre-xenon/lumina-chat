import { StringRouteParams, useParams } from "wouter";
import { Logo } from "../../components/logo";
import { ScrollArea } from "~/components/ui/scroll-area";
import { Chat } from "~/types/chat";
import { useFetch } from "~/hooks/useFetch";
import { ApiResponse } from "~/types/api";
import { ChatLink, ChatLinkSkeleton } from "~/components/chat-link";
import { CreateForm } from "./create-form";
import { toast } from "sonner";
import { apiCodesMap } from "~/lib/statuscodes";
import { sortChatsByTime } from "~/lib/time";

export function Home() {
  const { chatId } = useParams<StringRouteParams<"/:chatId?">>();

  const [chats, chatsLoading, setChats] =
    useFetch<ApiResponse<Chat[]>>("/v1/chats");

  function addChat(chat: Chat) {
    const oldChats = chats.data ? chats.data : [];
    setChats({ ...chats, data: sortChatsByTime([...oldChats, chat]) });
  }

  if (chats) {
    requestAnimationFrame(() => {
      if (chats["err_code"])
        toast.error(apiCodesMap[chats["err_code"]], {
          position: "top-center",
        });
    });
  }

  return (
    <main className="flex h-screen w-screen gap-4 p-2">
      <div className="flex h-full w-[25vw] flex-col gap-4">
        <Logo />

        <div className="flex w-full flex-col gap-2 overflow-hidden">
          {/* create chat button */}
          <div className="flex items-center justify-between">
            <h1 className="text-lg">Chats</h1>
            <CreateForm addChat={addChat} />
          </div>

          {/* chats */}
          <ScrollArea>
            <div className="flex flex-col px-3">
              {chatsLoading ? (
                // Chat skeletons
                [...Array(7).keys()].map((key) => (
                  <ChatLinkSkeleton key={key} />
                ))
              ) : chats["data"] ? (
                // Chats if exists
                chats["data"].map((chat) => (
                  <ChatLink key={chat.id} chat={chat} />
                ))
              ) : (
                // Has no chats
                <div className="flex h-full w-full flex-col items-center justify-center gap-2 text-muted-foreground">
                  <h1 className="text-lg">Look like you have no chats!</h1>
                  <p className="text-center">
                    Create a chat from the button above or get an invite link
                    from a friend.
                  </p>
                </div>
              )}
            </div>
          </ScrollArea>
        </div>
      </div>

      {chatId && (
        <div className="flex flex-col gap-4">
          <div>{chatId}</div>
        </div>
      )}
    </main>
  );
}
