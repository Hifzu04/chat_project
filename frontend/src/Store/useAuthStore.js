import { create } from "zustand";
import { axiosInstance } from "../lib/axios";
import toast from "react-hot-toast";
import { useChatStore } from "./useChatStore";

export const useAuthStore = create((set, get) => ({
  authUser: null,
  isSigningup: false,
  isLoggingin: false,
  isUpdatingProfile: false,
  onlineUsers: [],

  socket: null,
  isCheckingAuth: true,

  checkAuth: async () => {
    try {
      const res = await axiosInstance.get("/auth/check", { withCredentials: true });
      set({ authUser: res.data });
      get().connectSocket();
    } catch (error) {
      console.log("error in checkAuth:", error);
      set({ authUser: null });
    } finally {
      set({ isCheckingAuth: false });
    }
  },

  signup: async (data) => {
    set({ isSigningup: true });
    try {
      const res = await axiosInstance.post("/signup", data, { withCredentials: true });
      set({ authUser: res.data });
      toast.success("Account Created Successfully");
      get().connectSocket();
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    } finally {
      set({ isSigningup: false });
    }
  },

  login: async (data) => {
    set({ isLoggingin: true });
    try {
      // 🔥 FIX: actually perform the HTTP request
      const res = await axiosInstance.post("/login", data, { withCredentials: true });
      const { user } = res.data;       // now res is defined
      set({ authUser: user });
      toast.success("Logged in successfully");
      get().connectSocket();
      return true;
    } catch (error) {
      console.log("error during login:", error);
      toast.error(error.response?.data?.message || error.message);
      return false;
    } finally {
      set({ isLoggingin: false });
    }
  },

  logout: async () => {
    try {
      await axiosInstance.post("/logout", {}, { withCredentials: true });
      set({ authUser: null });
      toast.success("Logged out successfully");
      get().disconnectSocket();
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    }
  },

  updateProfile: async (data) => {
    set({ isUpdatingProfile: true });
    try {
      const res = await axiosInstance.put(
        "/user/update",
        data,
        { headers: { "Content-Type": "multipart/form-data" }, withCredentials: true }
      );
      set({ authUser: res.data });
      toast.success("Profile updated successfully");
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    } finally {
      set({ isUpdatingProfile: false });
    }
  },

  connectSocket: () => {
    const { authUser, socket: existing } = get();
    if (!authUser || existing) return;

    const base = import.meta.env.VITE_API_URL.replace(/\/$/, "");
    const wsUrl = base.replace(/^http/, "ws");
    const socket = new WebSocket(`${wsUrl}/ws?userId=${authUser._id}`);

    socket.onopen = () => console.log("WebSocket connected");

    socket.addEventListener("message", ({ data }) => {
      const { event: type, data: payload } = JSON.parse(data);
      if (type === "getOnlineUsers") {
        set({ onlineUsers: payload });
      }
      if (type === "newMessage") {
        useChatStore.getState().handleIncomingMessage(payload);
      }
    });

    socket.onclose = () => {
      console.log("WebSocket disconnected");
      set({ socket: null });
    };
    socket.onerror = (err) => console.error("WebSocket error:", err);

    set({ socket });
  },

  disconnectSocket: () => {
    const socket = get().socket;
    if (socket) {
      socket.close();
      set({ socket: null, onlineUsers: [] });
    }
  },
}));
