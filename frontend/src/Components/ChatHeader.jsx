import { X } from "lucide-react";
import { Link } from "react-router-dom";


import { useChatStore } from "../Store/useChatStore";
import { useAuthStore } from "../Store/useAuthStore";


const ChatHeader = () => {
  const { selectedUser, setSelectedUser } = useChatStore();
  const { onlineUsers } = useAuthStore();



  if (!selectedUser) {
    // nothing to show
    return null;
  }
  return (

    <div className="flex items-center justify-between bg-base-300 p-2">

      <div className="flex items-center gap-3">
        {/* Avatar */}
        <div className="avatar">
          <div className="size-10 rounded-full ">
            <img src={selectedUser.profile_pic || "/avatar.png"} alt={selectedUser.fullName} />
          </div>
        </div>

        {/* User info */}
        <Link to={"/SelectedUserProfile"} >
          <div>
            <h3 className="font-medium  ">{selectedUser.fullname}</h3>
            <p className="text-sm text-base-content/70">
              
              {onlineUsers.includes(selectedUser._id) || selectedUser?.email === "bot@chatnest.com"  ? "Online" : "Offline"}
            </p>
          </div>
        </Link>
      </div>


      {/* Close button */}

      <button onClick={() => setSelectedUser(null)} className="pr-2">
        <X />
      </button>

    </div>

  );
};
export default ChatHeader;