import { useChatStore } from "../Store/useChatStore";
import { useEffect, useRef } from "react";
import { useAuthStore } from "../Store/useAuthStore";

import ChatHeader from "./ChatHeader";
import MessageInput from "./MessageInput";
import MessageSkeleton from "./skeletons/MessageSkeleton";

import { formatMessageTime } from "../lib/utils";

const ChatContainer = () => {
  const {
    messages = [],
    getMessages,
    isMessagesLoading,
    selectedUser,
    // subscribeToMessages,
    // unsubscribeFromMessages,
  } = useChatStore();
  const { authUser } = useAuthStore();
  const messageEndRef = useRef(null);

  useEffect(() => {
    getMessages(selectedUser._id);

    // subscribeToMessages();

    // return () => unsubscribeFromMessages();

  }, [selectedUser._id, getMessages]);
  // console.log(`selected user is ,${messages.}`);

  useEffect(() => {
    if (messageEndRef.current && messages) {
      messageEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  if (isMessagesLoading) {
    return (
      <div className="flex-1 flex flex-col overflow-auto">
        <ChatHeader />
        <MessageSkeleton />
        <MessageInput />
      </div>
    );
  }


  return (
    <div className="flex-1 flex flex-col overflow-auto">
      <ChatHeader />

      <div className="flex-1 overflow-y-auto p-4 space-y-4">

        {messages === null ? (<div className="text-center text-zinc-500 mt-10">
          No messages yet. Start the conversation!
        </div>) : (messages.map((message) => (
          <div
            key={message._id}
            className={`chat ${message.sender_id === authUser.id ? "chat-end" : "chat-start"}`}
            ref={messageEndRef}

          >
            <div className=" chat-image avatar">
              <div className="size-10 rounded-full border">
                <img
                  src={
                    message.sender_id === authUser.id
                      ? authUser.profile_pic || "/avatar.png"
                      : selectedUser.profile_pic || "/avatar.png"
                  }
                  alt="profile pic"
                />
              </div>
            </div>
            <div className="chat-header mb-1">
              <time className="text-xs opacity-50 ml-1">
                {formatMessageTime(message.created_at)}
              </time>
            </div>
            <div className="chat-bubble flex flex-col">
              {message.images?.map((url, i) => (
                <img
                  key={i}
                  src={url}
                  alt="Attachment"
                  className="sm:max-w-[200px] rounded-md mb-2"
                />
              ))}
              {message.text && <p>{message.text}</p>}
            </div>
          </div>
        )))}

      </div>

      <MessageInput />
    </div>
  );
};
export default ChatContainer;