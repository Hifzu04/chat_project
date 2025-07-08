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
      set({ messages: [...messages, messagePayload] });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
      throw error;
    }
  },
  subscribeToMessages: () => {
    const socket = useAuthStore.getState().socket;
    if (!socket) return
    socket.onmessage = ({ data }) => {
      const { event, data: payload } = JSON.parse(data)
      if (event === "newMessage" && payload.sender_id === get().selectedUser._id) {
        set({ messages: [...get().messages, payload] })
      }
    }
  },

  unsubscribeFromMessages: () => {
    const socket = useAuthStore.getState().socket;
    if (socket) socket.onmessage = null;
  },






  setSelectedUser: (selectedUser) => set({ selectedUser })




}))

