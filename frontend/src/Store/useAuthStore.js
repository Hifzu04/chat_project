import { create } from "zustand";
import { axiosInstance } from "../lib/axios";
import toast from "react-hot-toast";
import { useChatStore } from "./useChatStore";

export const useAuthStore = create((set, get) => ({
  authUser: null,
  onlineUsers: [],
  socket: null,
  isCheckingAuth: true,

  checkAuth: async () => {
    try {
      const res = await axiosInstance.get("/auth/check", { withCredentials: true });
      const u = res.data.user || res.data;
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

  signup: async (data) => {
    try {
      const res = await axiosInstance.post("/signup", data, { withCredentials: true });
      const u = res.data;
      u._id = u.id;
      set({ authUser: u });
      toast.success("Account Created Successfully");
      get().connectSocket();
      return true;
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
      return false;
    }
  },

  login: async (data) => {
    try {
      const res = await axiosInstance.post("/login", data, { withCredentials: true });
      const { token, user } = res.data;
      user._id = user.id;
      set({ authUser: user });
      axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      toast.success("Logged in successfully");
      get().connectSocket();
      return true;
    } catch (error) {
      console.error("login failed:", error);
      toast.error(error.response?.data?.message || error.message);
      return false;
    }
  },

  logout: async () => {
    try {
      await axiosInstance.post("/logout", {}, { withCredentials: true });
      get().disconnectSocket();
      set({ authUser: null });
      toast.success("Logged out successfully");
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    }
  },

  connectSocket: () => {
    const { authUser, socket: existing } = get();
    //console.log("🧪 connectSocket() called. authUser:", authUser, "existing socket:", !!existing);
    if (!authUser || existing) return;


    //for loacl host remove comment and commet the line 80-83
    // const wsUrl = `ws://localhost:8000/ws?userId=${authUser._id}`;
    // console.log("👉 Opening WS:", wsUrl);

    const apiUrl = import.meta.env.VITE_API_URL || window.location.origin;
    const { host, protocol } = new URL(apiUrl);
    const wsProtocol = protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${wsProtocol}//${host}/ws?userId=${authUser._id}`;


    const socket = new WebSocket(wsUrl);

    socket.onopen = () => console.log("✅ WS open");
    socket.onerror = e => console.error("❌ WS error", e);
    socket.onclose = () => console.log("⚠️ WS closed");

    socket.onmessage = ({ data }) => {
     // console.log("📨 raw message:", data);
      let parsed;
      try {
        parsed = JSON.parse(data);
      } catch (err) {
        return console.error("❌ WS parse error:", err, data);
      }
      const { event, data: payload } = parsed;
     // console.log("🟦 WS event:", event, "payload:", payload);

      if (event === "getOnlineUsers") {
        console.log("🟩 setting onlineUsers to:", payload);
        set({ onlineUsers: payload });
      }
      if (event === "newMessage") {
        console.log("🟪 dispatching newMessage:", payload);
        useChatStore.getState().handleIncomingMessage(payload);
      }
    };

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
