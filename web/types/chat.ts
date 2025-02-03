import { Message } from "./message";

export type Chat = {
  id: string;
  name: string;
  invite_link: string;
  picture: string | null;
  message: Message;
};
