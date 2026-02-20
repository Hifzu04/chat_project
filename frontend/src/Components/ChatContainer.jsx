import { useChatStore } from "../Store/useChatStore";
import { useEffect, useRef } from "react";
import { useAuthStore } from "../Store/useAuthStore";

import ChatHeader from "./ChatHeader";
import MessageInput from "./MessageInput";
import MessageSkeleton from "./skeletons/MessageSkeleton";

import { formatMessageTime } from "../lib/utils";

const ChatContainer = () => {
  const {
    //messages: rawMessages: This renames the messages array from your store to rawMessages locally.
    messages: rawMessages,
    getMessages,
    isMessagesLoading,
    selectedUser,
    isAiTyping,

  } = useChatStore();
  const { authUser } = useAuthStore();

  // Ensure messages is always an array
  //Why: If your backend accidentally returns null or a string instead of an array, 
  //running .map() on it later will crash your entire React app.
  //This line guarantees that messages is always an array, even if the server messes up.
  const messages = Array.isArray(rawMessages) ? rawMessages : [];





  //What it does: Every time you click a different friend in the sidebar (changing selectedUser._id), 
  //this effect runs and fetches the chat history for that specific friend.
  useEffect(() => {
    if (selectedUser?._id) {
      getMessages(selectedUser._id);

    }
  }, [
    selectedUser?._id,
    getMessages,

  ]);

  const messageEndRef = useRef(null);


  //The Logic: Later in the HTML, you put ref={messageEndRef} on the message <div>. Because it is inside the .map() loop,
  //the laser gets moved to every single message as React builds the list. By the time the loop finishes, the laser is pointing at the very last message.
  //the Effect: Whenever the messages array changes (like when you send or receive a new message),
  //this tells the browser to smoothly scroll down to wherever that laser is pointing.
  useEffect(() => {
    if (messageEndRef.current && messages.length) {
      messageEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  if (isMessagesLoading) {
    return (
      //flex-1: Takes up all remaining space to the right of your Sidebar.
      //overflow-auto: Prevents the entire page from scrolling. Only the message area inside will scroll.
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
                key={message.id}
                //chat: The base DaisyUI class that sets up the chat bubble grid.
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
                      ? "bg-primary/50 text-primary-content"
                      : "bg-base-300 text-base-content"}
                  `}
                >
                  {/* //message.images?.map: The ?. is a safety check. It says: "If the images array exists, loop through it.
                 If it doesn't exist (because it's just a text message), don't crash, just ignore this part." */}

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

        {/* 👇 NEW TYPING INDICATOR BUBBLE 👇 */}
        {isAiTyping && selectedUser?.email === "bot@chatnest.com" && (
          <div className="chat chat-start" ref={messageEndRef}>
            <div className="chat-image avatar">
              <div className="size-10 rounded-full border">
                <img
                  src={selectedUser.profile_pic || "/avatar.png"}
                  alt="bot profile pic"
                />
              </div>
            </div>

            <div className="chat-header mb-1">
              <span className="text-xs opacity-50 ml-1">NestBot is typing</span>
            </div>

            <div className="chat-bubble bg-base-300 text-base-content flex items-center justify-center p-3 h-10 w-16">
              {/* DaisyUI Loading Dots */}
              <span className="loading loading-dots loading-sm opacity-60"></span>
            </div>
          </div>
        )}




      </div>



      <MessageInput />
    </div>
  );
};

export default ChatContainer;





{/*
 
  Question 1: Why useRef instead of useState for scrolling?

const messageEndRef = useRef(null);
You used useRef here because you need to physically touch a DOM element (an HTML <div>) without causing React to redraw the screen.

Here is why useState would be the wrong tool for this specific job:

1. The Rule of Re-rendering (useState's job)

Every time you update a useState variable, React says: "The data changed! I must completely re-draw (re-render) this component on the screen to show the new data."

This is great for things like text inputs or opening menus.

2. The "Remote Control" (useRef's job)

useRef is like a secret laser pointer. It points directly at an HTML element on the screen so you can control it (like telling it to scrollIntoView()).

Crucially: Updating a useRef does NOT trigger a re-render. What if you used useState?
If you tried to store the HTML <div> inside useState, every time a new message appeared and the div moved, React would try to update the state. 
That state update would trigger a re-render. That re-render would recreate the div, which would trigger another state update... creating an infinite loop that crashes the browser!

Summary: * Use useState when you want the visual data on the screen to change.

Use useRef when you want a direct hook into an HTML element to do things like scrolling, focusing an input, or playing a video, without interrupting React's rendering cycle.



  
  */}