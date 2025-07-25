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

  // 1) Check existing session on page load
  checkAuth: async () => {
    try {
      const res = await axiosInstance.get("/auth/check", { withCredentials: true });
      const u = res.data;
      u._id = u.id;
      set({ authUser: u });
      get().connectSocket();
    } catch (error) {
      console.log("error in checkAuth:", error);
      set({ authUser: null });
    } finally {
      set({ isCheckingAuth: false });
    }
  },

  // 2) Signup new user
  signup: async (data) => {
    set({ isSigningup: true });
    try {
      const res = await axiosInstance.post("/signup", data, { withCredentials: true });
      set({ authUser: res.data });
      toast.success("Account Created Successfully");
      get().connectSocket();
      return true;
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
      return false;
    } finally {
      set({ isSigningup: false });
    }
  },

  // 3) Login existing user
  login: async (data) => {
    set({ isLoggingin: true });
    try {
      // • send credentials
      const res = await axiosInstance.post("/login", data, { withCredentials: true });
      // • grab token + user
      const { token, user } = res.data;
      // • set fallback header for subsequent requests
      axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      // • update state + connect websocket

      user._id = user.id;
      set({ authUser: user }); toast.success("Logged in successfully");
      get().connectSocket();
      return true;
    } catch (error) {
      console.error("login failed:", error);
      toast.error(error.response?.data?.message || error.message);
      return false;
    } finally {
      set({ isLoggingin: false });
    }
  },

  // 4) Logout
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

  // 5) Update profile
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

  // 6) WebSocket setup
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
