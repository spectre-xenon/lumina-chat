import { format } from "date-fns";
import { Chat } from "~/types/chat";

export function formatDate(date: string) {
  const dateObj = new Date(date);
  const dateDifferenceInTime = new Date().getTime() - dateObj.getTime();

  // conerting milli seconds to days
  // (1000 milliseconds * (60 seconds * 60 minutes) * 24 hours)
  const dateDifferenceInDays = dateDifferenceInTime / (1000 * 60 * 60 * 24);

  //After returning in particular formats as of our convinent
  if (dateDifferenceInDays < 1) {
    return format(dateObj, "p"); // 10:04 am
  } else if (dateDifferenceInDays < 2) {
    return "Yesterday"; // just YesterDay
  } else if (dateDifferenceInDays <= 7) {
    return format(dateObj, "eee"); //like monday , tuesday , wednesday ....
  } else {
    return format(dateObj, "P"); // if it was more than a week before it will returns as like 05/23/2022
  }
}

export function sortChatsByTime(chats: Chat[]) {
  return chats.sort((a, b) =>
    b.message.sent_at.localeCompare(a.message.sent_at),
  );
}
