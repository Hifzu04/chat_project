import { create } from "zustand";
import { axiosInstance } from "../lib/axios";
import toast from "react-hot-toast";
import { useAuthStore } from "./useAuthStore";

export const useChatStore = create((set, get) => ({
  users: [],
  messages: [],
  isUsersLoading: false,
  isMessagesLoading: false,
  selectedUser: null,
  unseenCounts: {},
  isAiTyping: false,

  //Get the list of all the users and put them put them in an array
  getUsers: async () => {
    set({ isUsersLoading: true });
    try {
      const res = await axiosInstance.get("/users", { withCredentials: true });
      //// Normalization: Ensure every user has an `_id` property matching `id`

      let users = res.data.map(u => ({ ...u, _id: u.id }));


      //NEW SORTING LOGIC to pin the bot at top of the  list 
      // 2. Find the exact index where NestBot is sitting
      const botIndex = users.findIndex((u) => u.email === "bot@chatnest.com");

      // 3. If the bot exists in the list...
      if (botIndex !== -1) {
        // Remove the bot from its current spot in the array
        const botUser = users.splice(botIndex, 1)[0];
        // .unshift() pushes the bot into the very first slot (index 0)
        users.unshift(botUser);
      }
      set({ users });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    } finally {
      set({ isUsersLoading: false });
    }
  },
  //When you click on "Bob," this function runs. It fetches the conversation history between You and Bob.
  getMessages: async (userId) => {
    set({ isMessagesLoading: true });
    try {
      const res = await axiosInstance.get(`/messages/${userId}`, { withCredentials: true });
      set({ messages: res.data });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    } finally {
      set({ isMessagesLoading: false });
    }
  },

  sendMessage: async (formData) => {
    const { messages, selectedUser } = get();
    try {
      const res = await axiosInstance.post("/messages/send", formData, {
        headers: { "Content-Type": "multipart/form-data" },
        withCredentials: true,
      });
      const messagePayload = res.data;

      const socket = useAuthStore.getState().socket;
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ event: "sendMessage", data: messagePayload }));
      }

      if (selectedUser?.email === "bot@chatnest.com") {
        set({ isAiTyping: true });
      }
      //// 3. UI UPDATE: Add the new message to the bottom of the chat list
      set({ messages: [...(Array.isArray(messages) ? messages : []), messagePayload] });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
      throw error;
    }
  },

  //YET TO IMPLEMENT: 
  //This function is called whenever a new message arrives via the WebSocket connection. It updates the chat UI and unseen message counts accordingly.

  handleIncomingMessage: (msg) => {
    const { sender_id } = msg;
    const { selectedUser, users, unseenCounts } = get();
    // LOGIC 1: Re-order the Sidebar (Move sender to top)
    const rest = users.filter(u => u._id !== sender_id);
    const sender = users.find(u => u._id === sender_id);
    //if an old friend messages you, their name jumps to the top of the list. That is what [sender, ...rest] does.
    const newUsers = sender ? [sender, ...rest] : users;
    // LOGIC 2: Update Unseen Counts (Increment if sender isn't currently selected)
    const newCounts = { ...unseenCounts };
    // If I am NOT currently talking to this person...
    if (!selectedUser || selectedUser._id !== sender_id) {
      newCounts[sender_id] = (newCounts[sender_id] || 0) + 1;
    }
    // LOGIC 3: Update Messages in Chat Window (Only if sender is currently selected)
    if (selectedUser && selectedUser._id === sender_id) {
      set(state => ({ messages: [...state.messages, msg] }));
    }
    //// Apply sidebar changes
    set({ users: newUsers, unseenCounts: newCounts, isAiTyping: false });
  },

  //When we click on a user in the sidebar, we want to do two things: 1) Mark that user as "selected" so we know who we're chatting with.
  // 2) Reset the unseen message count for that user to zero, since we're now viewing their messages.
  setSelectedUser: (user) => {
    if (!user) {
      // clear selection (and reset unseen counts if you like)
      return set({ selectedUser: null });
    }
    // otherwise select normally
    set({
      selectedUser: user,

      unseenCounts: {
        ...get().unseenCounts,
        [user._id]: 0,
      },
    });
  },
}));
