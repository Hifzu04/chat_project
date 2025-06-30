import { create } from "zustand"
import { axiosInstance } from "../lib/axios"
import toast from "react-hot-toast";


export const useChatStore = create((set, get) => ({
  users: [],
  messages: [],

  isUsersLoading: false,
  isMessagesLoading: false,
  selectedUser: null,




  getUsers: async () => {
    try {
      set({isUsersLoading:true})
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
      const res = await axiosInstance.post("/messages/send", formData, {
        headers: { "Content-Type": "multipart/form-data" }
      });
      // assuming your Go handler returns the saved message:
      set({ messages: [...messages, res.data] });
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
      throw error;
    }
  },




  setSelectedUser: (selectedUser) => set({ selectedUser })




}))

