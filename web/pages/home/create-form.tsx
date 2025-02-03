import { Label } from "@radix-ui/react-label";
import { BadgePlus } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { useLocation } from "wouter";
import {
  AlertDialog,
  AlertDialogTrigger,
  AlertDialogContent,
  AlertDialogTitle,
  AlertDialogCancel,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogDescription,
} from "~/components/ui/alert-dialog";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { simpleFetch } from "~/hooks/useFetch";
import { apiCodesMap, genericErrorsMap } from "~/lib/statuscodes";
import { ApiResponse } from "~/types/api";
import { Chat } from "~/types/chat";

type FormData = {
  name: string;
  picture: File | null;
};

const allowedMimes = {
  "image/jpeg": true,
  "image/png": true,
  "image/webp": true,
};

export function CreateForm({ addChat }: { addChat: (chat: Chat) => void }) {
  const [open, setOpen] = useState(false);
  const [formData, setFormData] = useState<FormData>({
    name: "",
    picture: null,
  });

  const [, navigate] = useLocation();

  function closeDialog() {
    setFormData({ name: "", picture: null });
    setOpen(false);
  }

  function handlePictureChange(e: React.ChangeEvent<HTMLInputElement>) {
    const files = e.target.files;

    if (!files) {
      return;
    }

    if (!(files[0].type in allowedMimes)) {
      e.target.value = "";
      return toast.error("This image type is now allowed!", {});
    }

    // 1MB
    if (files[0].size > 1e6) {
      e.target.value = "";
      return toast.error("The image size can't be bigger than 1 megabyte!", {});
    }

    setFormData({ ...formData, picture: files[0] });
  }

  async function handleSubmit() {
    console.log(formData.name);
    if (formData.name === "")
      return toast.error(genericErrorsMap["emptyFormField"]);
    if (formData.name.length < 3)
      return toast.error(genericErrorsMap["shortChatName"]);

    const body = new FormData();

    if (formData["picture"]) body.append("picture", formData["picture"]);
    body.append("name", formData["name"]);

    const data = await simpleFetch<ApiResponse<Chat[]>>("/v1/chat", navigate, {
      method: "POST",
      body: body,
    });

    if (data["err_code"]) {
      return toast.error(apiCodesMap[data["err_code"]], {
        position: "top-center",
      });
    }

    if (data["data"]) addChat(data["data"][0]);

    closeDialog();
  }

  return (
    <form>
      <AlertDialog open={open}>
        <AlertDialogTrigger asChild>
          <Button type="button" variant="ghost" onClick={() => setOpen(true)}>
            <BadgePlus />
          </Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Create a new chat.</AlertDialogTitle>
            <AlertDialogDescription>
              Create a new chat with a name and a picture (which is optional).
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="flex flex-col gap-6">
            <div className="grid gap-2">
              <Label htmlFor="picture">Picture</Label>
              <Input
                id="picture"
                name="picture"
                type="file"
                accept="image/*"
                className="cursor-pointer"
                onChange={handlePictureChange}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="name">Chat name*</Label>
              <Input
                id="name"
                name="name"
                type="text"
                placeholder="Starry night"
                onChange={(e) => {
                  setFormData({ ...formData, name: e.target.value });
                }}
              />
            </div>
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => closeDialog()}>
              Cancel
            </AlertDialogCancel>
            <Button type="submit" onClick={handleSubmit}>
              Submit
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </form>
  );
}
