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

  getUsers: async () => {
    set({ isUsersLoading: true });
    try {
      const res = await axiosInstance.get("/users", { withCredentials: true });
      const users = res.data.map((u) => ({ ...u, _id: u.id }));
      set({ users });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    } finally {
      set({ isUsersLoading: false });
    }
  },

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
    const { messages } = get();
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

      set({ messages: [...(Array.isArray(messages) ? messages : []), messagePayload] });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
      throw error;
    }
  },

  handleIncomingMessage: (msg) => {
    const { sender_id } = msg;
    const { selectedUser, users, unseenCounts } = get();

    // 1) Reorder users
    const rest = users.filter((u) => u._id !== sender_id);
    const sender = users.find((u) => u._id === sender_id);
    const newUsers = sender ? [sender, ...rest] : users;

    // 2) Increment unseen
    const newCounts = { ...unseenCounts };
    if (!selectedUser || selectedUser._id !== sender_id) {
      newCounts[sender_id] = (newCounts[sender_id] || 0) + 1;
    }

    // 3) Append to open chat
    if (selectedUser && selectedUser._id === sender_id) {
      set((s) => ({ messages: [...s.messages, msg] }));
    }

    // 4) Commit updates
    set({ users: newUsers, unseenCounts: newCounts });
  },

  setSelectedUser: (selectedUser) => set({ selectedUser, unseenCounts: { ...get().unseenCounts, [selectedUser._id]: 0 } }),
}));
