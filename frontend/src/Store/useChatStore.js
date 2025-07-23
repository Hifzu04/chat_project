import { create } from "zustand"
import { axiosInstance } from "../lib/axios"
import toast from "react-hot-toast";
import { useAuthStore } from "./useAuthStore";


export const useChatStore = create((set, get) => ({
  users: [],
  messages: [],

  isUsersLoading: false,
  isMessagesLoading: false,
  selectedUser: null,
  unseenCounts: {},




  getUsers: async () => {
    try {
      set({ isUsersLoading: true })
      const res = await axiosInstance.get("/users");
      const users = res.data.map((u) => ({
        ...u,
        _id: u.id,       // normalize
      }));
      set({ users });
    } catch (error) {
      toast.error(error.response.data.message);
    } finally {
      set({ isUsersLoading: false });
    }
  },


  // 2) Load message history for a given user
  getMessages: async (userId) => {
    set({ isMessagesLoading: true });
    try {
      const res = await axiosInstance.get(`/messages/${userId}`);
      set({ messages: res.data });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    } finally {
      set({ isMessagesLoading: false });
    }
  },

  // 3) Send a new message
  sendMessage: async (formData) => {
    const { messages } = get();
    try {
      // POST to your Go backend
      const res = await axiosInstance.post("/messages/send", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      // The server’s saved message payload (with .text, .images, .sender_id, .receiver_id, ._id, etc.)
      const messagePayload = res.data;

      // Real‑time: notify the receiver
      const socket = useAuthStore.getState().socket;
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(
          JSON.stringify({
            event: "sendMessage",
            data: messagePayload,
          })
        );
      }

      // Locally display it
      const updatedMessages = Array.isArray(messages) ? messages : [];
      set({ messages: [...updatedMessages, messagePayload] });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
      throw error;
    }
  },
 handleIncomingMessage: (msg) => {
    const { sender_id } = msg;
    const { selectedUser, users, unseenCounts } = get();

    // 1) Reorder users so sender jumps to front
    const rest = users.filter(u => u._id !== sender_id);
    const sender = users.find(u => u._id === sender_id);
    const newUsers = sender ? [sender, ...rest] : users;

    // 2) Increment unseen if not currently chatting
    const newCounts = { ...unseenCounts };
    if (!selectedUser || selectedUser._id !== sender_id) {
      newCounts[sender_id] = (newCounts[sender_id] || 0) + 1;
    }

    // 3) If that chat is open, also append the message
    if (selectedUser && selectedUser._id === sender_id) {
      set(state => ({ messages: [...state.messages, msg] }));
    }

    // 4) Commit both changes
    set({ users: newUsers, unseenCounts: newCounts });
  },










  setSelectedUser: (selectedUser) => set({ selectedUser })




}))

