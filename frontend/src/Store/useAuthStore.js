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
      const u = res.data.user || res.data;
      u.id = u.id;     // already there
      u._id = u.id;    // add Mongo‑style alias

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
      user._id = user.id;
      set({ authUser: user });
      // • set fallback header for subsequent requests
      axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      // • update state + connect websocket


      toast.success("Logged in successfully");
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

 connectSocket: () => {
  const { authUser, socket: existing } = get()
  if (!authUser || existing) return

  // parse your API base so you never include a path or trailing slash
  const apiUrl = import.meta.env.VITE_API_URL      // e.g. "https://chatnest-49i2.onrender.com"
  const { host, protocol } = new URL(apiUrl)

  // pick ws:// vs wss://
  const wsProto = protocol === "https:" ? "wss:" : "ws:"
  const wsUrl   = `${wsProto}//${host}/ws?userId=${authUser._id}`

  console.log("👉 Opening WS:", wsUrl)
  const socket = new WebSocket(wsUrl)

  socket.onopen    = () => console.log("✅ WS open")
  socket.onerror   = e => console.error("❌ WS error", e)
  socket.onclose   = () => console.log("⚠️ WS closed")
  socket.onmessage = ({ data }) => {
    const { event, data: payload } = JSON.parse(data)
    if (event === "getOnlineUsers") set({ onlineUsers: payload })
    if (event === "newMessage") useChatStore.getState().handleIncomingMessage(payload)
  }

  set({ socket })
},


  disconnectSocket: () => {
    const socket = get().socket;
    if (socket) {
      socket.close();
      set({ socket: null, onlineUsers: [] });
    }
  },
}));
