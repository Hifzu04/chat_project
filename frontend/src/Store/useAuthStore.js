
//Zustand is a small, fast, and scalable global  state management library for React.
//It provides a simple and minimal API to manage global or shared state in React applications —
//often as an alternative to Redux, Context API, or Recoil.
import { create } from "zustand"
import { axiosInstance } from "../lib/axios"

export const useAuthStore = create((set) => ({
    authUser: null,
    isSigningup: false,
    isLoggingin: false,
    isUpdatingProfile: false,

    isCheckingAuth: true,


    //when we refresh page
  checkAuth: async () => {
        try {
            const res = await axiosInstance.get("/auth/check");  //auth/check are from backend
            set({ authUser: res.data });
        } catch (error) {
            console.log("error in check auth:", error);
            set({ authUser: null });
        } finally {
            set({ isCheckingAuth: false });
        }
    },

}))