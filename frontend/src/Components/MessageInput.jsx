import { useRef, useState, useEffect } from "react";
import { useChatStore } from "../Store/useChatStore";
import { Image, Send, X } from "lucide-react";
import toast from "react-hot-toast";
import { axiosInstance } from "../lib/axios";

const MessageInput = () => {
  const [text, setText] = useState("");
  const [imagePreview, setImagePreview] = useState(null);
  const [imageFile, setImageFile] = useState(null);

  // Debug: log whenever imageFile changes
  useEffect(() => {
    console.log("🔥 imageFile state changed:", imageFile);
  }, [imageFile]);

  const fileInputRef = useRef(null);
  const { sendMessage, selectedUser } = useChatStore();

  const handleImageChange = (e) => {
    // Debug #1: see if this ever runs
    console.log("⚡️ handleImageChange fired, files:", e.target.files);

    const file = e.target.files?.[0];
    console.log("⚡️ file[0]:", file);

    if (!file?.type.startsWith("image/")) {
      toast.error("Please select an image file");
      return;
    }

    setImageFile(file);
    setImagePreview(URL.createObjectURL(file));
  };

  const removeImage = () => {
    setImagePreview(null);
    setImageFile(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

 const handleSendMessage = async (e) => {
  e.preventDefault();
  if (!text.trim() && !imageFile) return;

  // Build a real FormData
  const formData = new FormData();
  formData.append("receiver_id", selectedUser._id);
  if (text.trim()) formData.append("text", text.trim());
  if (imageFile) formData.append("images", imageFile);

  try {
    // Delegate to your store’s sendMessage (which posts + broadcasts)
    await sendMessage(formData);
  } catch (error) {
    console.error("❌ Error sending message:", error);
    toast.error("Failed to send message");
  } finally {
    setText("");
    removeImage();
  }
};

  return (
    <div className="p-4 w-full">
      {imagePreview && (
        <div className="mb-3 flex items-center gap-2">
          <div className="relative">
            <img
              src={imagePreview}
              alt="Preview"
              className="w-20 h-20 object-cover rounded-lg border border-zinc-700"
            />
            <button onClick={removeImage} type="button"
              className="absolute -top-1.5 -right-1.5 w-5 h-5 rounded-full bg-base-300 flex items-center justify-center"
            >
              <X className="size-3" />
            </button>
          </div>
        </div>
      )}

      <form onSubmit={handleSendMessage} className="flex items-center gap-2">
        <div className="flex-1 flex gap-2">
          <input
            type="text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="Type a message..."
            className="w-full input input-bordered rounded-lg input-sm sm:input-md"
          />
          <input
            type="file"
            accept="image/*"
            className="hidden"
            ref={fileInputRef}
            onChange={handleImageChange}
          />

          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className={` sm:flex btn btn-circle ${
              imagePreview ? "text-emerald-500" : "text-zinc-400"
            }`}
          >
            <Image size={20} />
          </button>
        </div>
        <button
          type="submit"
          className="btn btn-sm btn-circle"
          disabled={!text.trim() && !imagePreview}
        >
          <Send size={22} />
        </button>
      </form>
    </div>
  );
};

export default MessageInput;
