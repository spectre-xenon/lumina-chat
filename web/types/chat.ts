export type Chat = {
  id: string;
  name: string;
  invite_link: string;
  picture: string;
  last_message: {
    id: string;
    sender: string;
    content: string;
    sent_at: string;
  };
};
