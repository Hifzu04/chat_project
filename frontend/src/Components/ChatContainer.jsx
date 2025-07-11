import { useChatStore } from "../Store/useChatStore";
import { useEffect, useRef } from "react";
import { useAuthStore } from "../Store/useAuthStore";

import ChatHeader from "./ChatHeader";
import MessageInput from "./MessageInput";
import MessageSkeleton from "./skeletons/MessageSkeleton";

import { formatMessageTime } from "../lib/utils";

const ChatContainer = () => {
  const {
    messages: rawMessages,
    getMessages,
    isMessagesLoading,
    selectedUser,
    subscribeToMessages,
    unsubscribeFromMessages,
  } = useChatStore();
  const { authUser } = useAuthStore();

  // Ensure messages is always an array
  const messages = Array.isArray(rawMessages) ? rawMessages : [];

  const messageEndRef = useRef(null);

  useEffect(() => {
    if (selectedUser?._id) {
      getMessages(selectedUser._id);
      subscribeToMessages();
      return () => unsubscribeFromMessages();
    }
  }, [
    selectedUser?._id,
    getMessages,
    subscribeToMessages,
    unsubscribeFromMessages,
  ]);

  useEffect(() => {
    if (messageEndRef.current && messages.length) {
      messageEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  if (isMessagesLoading) {
    return (
      <div className="flex-1 flex flex-col overflow-auto bg-base-200 text-base-content">
        <ChatHeader />
        <MessageSkeleton />
        <MessageInput />
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col overflow-auto bg-base-200 text-base-content">
      <ChatHeader />

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.length === 0 ? (
          <div className="text-center text-zinc-500 mt-10">
            No messages yet. Start the conversation!
          </div>
        ) : (
          messages.map((message) => {
            const isMe = message.sender_id === authUser.id;
            return (
              <div
                key={message._id}
                className={`chat ${isMe ? "chat-end" : "chat-start"}`}
                ref={messageEndRef}
              >
                <div className="chat-image avatar">
                  <div className="size-10 rounded-full border">
                    <img
                      src={
                        isMe
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

                <div
                  className={`
                    chat-bubble flex flex-col p-3 rounded-lg
                    ${isMe
                      ? "bg-primary text-primary-content"
                      : "bg-base-300 text-base-content"}
                  `}
                >
                  {message.images?.map((url, i) => (
                    <img
                      key={i}
                      src={url}
                      alt="Attachment"
                      className="sm:max-w-[200px] rounded-md  mb-2  "
                    />
                  ))}
                  {message.text && <p>{message.text}</p>}
                </div>
              </div>
            );
          })
        )}
      </div>

      <MessageInput />
    </div>
  );
};

export default ChatContainer;