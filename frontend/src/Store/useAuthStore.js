import { create } from "zustand";
import { axiosInstance } from "../lib/axios";
import toast from "react-hot-toast";
import { useChatStore } from "./useChatStore";



//In Zustand(alternative for redux), when you define a store, you are writing a function that returns your state object.
//Zustand passes two helper tools into this function: set and get.
//set (The Writer): //Why you need it: When a user logs in successfully, you need to change authUser: null to authUser: { name: "Bob" }.
//  You cannot just say this.authUser = user.
//get() gives us access to "this" store instance.
//You must use set({ authUser: user }) so React knows the data changed and re-renders the UI.
//Why: Zustand uses a "Create" function to build a custom React Hook.
//Why not Redux? Redux requires wrapping your app in a <Provider>, creating "reducers", and "actions"
export const useAuthStore = create((set, get) => ({
  authUser: null,
  onlineUsers: [],
  socket: null,
  isCheckingAuth: true,
  isUpdatingProfile: false,
  //async : jubb takk fetch ho raha hai other kaam dekho. Khali nhi baithe dunga 
  //axios.get returns a Promise If you remove await, res will be the Promise object itself, not the actual user data. 
  // await(.then in non async ) forces the code to pause right here until the server responds.
  checkAuth: async () => {
    try {
      // withCredentials: true  : send this request with cookies to BE

      const res = await axiosInstance.get("/auth/check", { withCredentials: true });
      // console.log("checkAuth response:", res.data._id);

      const u = res.data.user || res.data;
      u._id = u.id;
      set({ authUser: u });
      get().connectSocket();
    } catch (error) {
      console.log("error in checkAuth:", error);
      set({ authUser: null });
    }
    // always executed
    finally {
      set({ isCheckingAuth: false });
    }
  },

  signup: async (data) => {
    try {
      //async/await	To pause execution while waiting for slow server responses without freezing the browser.
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
      delete axiosInstance.defaults.headers.common["Authorization"];
      toast.success("Logged out successfully");
    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
    }
  },





  updateProfile: async (formData) => {
    try {
      set({ isUpdatingProfile: true });
      const res = await axiosInstance.put("/user/update", formData, {
        headers: { "Content-Type": "multipart/form-data" },
        withCredentials: true,
      });
      toast.success("Profile updated successfully");




    } catch (error) {
      toast.error(error.response?.data?.message || error.message);
      console.log("error in updateProfile:", error);
      return false;
    } finally {
      set({ isUpdatingProfile: false });
    }

  },



  connectSocket: () => {
    const { authUser, socket: existing } = get();         //existing: If you are already connected, don't open a second connection.
    //This prevents creating 100 duplicate connections if the user clicks around fast.
    //console.log("connectSocket() called. authUser:", authUser, "existing socket:", !!existing);
    if (!authUser || existing) return;


    // for loacl host remove comment and commet the line 80-83
    //const wsUrl = `ws://localhost:8000/ws?userId=${authUser._id}`;
    //console.log("👉 Opening WS:", wsUrl);

    const apiUrl = import.meta.env.VITE_API_URL || window.location.origin;
    const { host, protocol } = new URL(apiUrl);
    const wsProtocol = protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${wsProtocol}//${host}/ws?userId=${authUser._id}`;


    const socket = new WebSocket(wsUrl);

    socket.onopen = () => console.log("✅ WS open");
    socket.onerror = e => console.error(" WS error", e);
    socket.onclose = () => console.log("WS closed");

    socket.onmessage = ({ data }) => {
      // console.log("📨 raw message:", data);
      let parsed;
      try {
        parsed = JSON.parse(data);
      } catch (err) {
        return console.error(" WS parse error:", err, data);
      }
      const { event, data: payload } = parsed;
      // console.log("🟦 WS event:", event, "payload:", payload);

      if (event === "getOnlineUsers") {
        console.log("🟩 setting onlineUsers to:", payload);
        set({ onlineUsers: payload });
      }
      if (event === "newMessage") {
        // console.log("🟪 dispatching newMessage:", payload);
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
