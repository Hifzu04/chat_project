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
  // useEffect(() => {
  //   console.log("🔥 imageFile state changed:", imageFile);
  // }, [imageFile]);

  const fileInputRef = useRef(null);
  const { sendMessage, selectedUser } = useChatStore();

  const handleImageChange = (e) => {
    // Debug #1: see if this ever runs
    console.log(" handleImageChange fired, files:", e.target.files);
    // Debug #2: Check the first file
    const file = e.target.files?.[0];
    //console.log("⚡️ file[0]:", file);
    //// 2. Validation: Stop if it's not an image (e.g., if they pick a PDF)
    if (!file?.type.startsWith("image/")) {
      toast.error("Please select an image file");
      console.log("No file is selected or an unexpected error during image upload")
      return;
    }

    setImageFile(file);
    //// URL.createObjectURL() creates a temporary fake URL pointing to the file 
    // in the browser's memory. This string is not null, so the UI updates.
    setImagePreview(URL.createObjectURL(file));
  };

  const removeImage = () => {
    setImagePreview(null);
    setImageFile(null);
    //fileInputRef.current.value = "" resets the hidden input.\
    //  If you don't do this, and the user tries to upload the exact same photo again immediately after removing it,
    //  the browser won't trigger onChange, and nothing will happen.
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const handleSendMessage = async (e) => {
    e.preventDefault();   //e.preventDefault(); // Stop page reload

    if (!text.trim() && !imageFile) {
      toast.error("Cannot send an empty message");
      console.log("Attempted to send an empty message. Text:", text, "ImageFile:", imageFile);
      return;
    }  
     

    // Build a real FormData :FormData is a special browser API that formats data as multipart/form-data
    const formData = new FormData();
    formData.append("receiver_id", selectedUser._id);
    if (text.trim()) formData.append("text", text.trim());
    if (imageFile) formData.append("images", imageFile);

    try {
      // Delegate to your store’s sendMessage (which posts + broadcasts)
      await sendMessage(formData);
    } catch (error) {
      console.error(" Error sending message:", error);
      toast.error("Failed to send message");
    } finally {
      setText("");
      removeImage();
    }
  };

  return (
    <div className="p-4 w-full">

      {imagePreview && (
        <div className="mb-3 flex items-center gap-2 animate-pulse ">
          <div className="relative">
            <img
              src={imagePreview}
              alt="Preview"
              className="w-20 h-20 object-cover rounded-lg border border-zinc-700"
            />
            <button onClick={removeImage} type="button"
              //relative (Parent) & absolute (Button): This allows the "X" button to float on top of the image.
              //-top-1.5 -right-1.5: This creates the "badge" effect, pushing the button slightly outside the top-right corner of the image.
              className="absolute -top-1.5 -right-1.5 w-5 h-5 rounded-full bg-base-300 flex items-center justify-center"
            >
              <X className="size-3 " />
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
            className="w-full bg-base-300 input input-bordered rounded-lg input-sm sm:input-md"
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
            //onClick: When you click this nice button, React secretly clicks the hidden file input (ref.current.click()), opening the file explorer.
            onClick={() => fileInputRef.current?.click()}
            className={` sm:flex btn btn-circle ${imagePreview ? "text-emerald-500" : "text-zinc-400"
              }`}
          >
            <Image size={20} />
          </button>
        </div>


        <button
          type="submit"
          className="btn btn-sm btn-circle  text-white hover:bg-emerald-600"
          //Disabled Logic: The button is grayed out and unclickable unless there is either text typed OR an image selected
          // This prevents sending empty "ghost" messages.
          disabled={!text.trim() && !imagePreview}
        >
          <Send size={20} />
        </button>
      </form>
    </div>
  );
};

export default MessageInput;









// // useState is for React's internal data,
// //  while useRef allows you to reach outside of React and forcefully manipulate the raw browser DOM.