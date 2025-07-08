import { X } from "lucide-react";


import { useChatStore } from "../Store/useChatStore";
import { useAuthStore } from "../Store/useAuthStore";


const ChatHeader = () => {
  const { selectedUser, setSelectedUser } = useChatStore();
   const { onlineUsers } = useAuthStore();

  return (
    <div className="">
      <div className="flex items-center justify-between bg-base-300 p-2">
        <div className="flex items-center gap-3">
          {/* Avatar */}
          <div className="avatar">
            <div className="size-10 rounded-full relative">
              <img src={selectedUser.profile_pic || "/avatar.png"} alt={selectedUser.fullName} />
            </div>
          </div>

          {/* User info */}
          <div>
            <h3 className="font-medium  ">{selectedUser.fullname}</h3>
            <p className="text-sm text-base-content/70">
              {onlineUsers.includes(selectedUser._id) ? "Online" : "Offline"}
            </p> 
          </div>
        </div>

        {/* Close button */}
        <button onClick={() => setSelectedUser(null)}>
          <X />
        </button>
      </div>
    </div>
  );
};
export default ChatHeader;