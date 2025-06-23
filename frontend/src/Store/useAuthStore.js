
//Zustand is a small, fast, and scalable global  state management library for React.
//It provides a simple and minimal API to manage global or shared state in React applications —
//often as an alternative to Redux, Context API, or Recoil.
import { create } from "zustand"
import { axiosInstance } from "../lib/axios"
import toast from "react-hot-toast";
import axios from "axios";

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
            console.log("error in check auth check ur axios or BE:", error);
            set({ authUser: null });
        } finally {
            set({ isCheckingAuth: false });
        }
    },

    //signup 
    signup: async (data) => {
        set({ isSigningup: true })
        try {
            const res = await axiosInstance.post("/signup", data);
            set({ authUser: res.data });
            toast.success("Account Created Sucessfully")

        } catch (error) {
            const message =
                error.response?.data?.message ||
                error.message ||
                "Something went wrong";
            toast.error(message);
        } finally {
            set({ isSigningup: false });
        }
    },

    //login 
    login: async (data) => {
        set({ isLoggingin: true })
        try {
            const res = await axiosInstance.post("/login", data);
            const { token, user } = res.data;
            set({ authUser: user });
            axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
            toast.success("Loggedin sucessfully");
           //return true; // Indicate successful login
        } catch (error) {
            console.log("some thing went wrong while logging in ")
            const message = error.response?.data?.message || error.message || "something went wrong"
            toast.error(message)
            return false; // Indicate failed login
            
        } finally {
            set({ isLoggingin: false })
        }
    },

    logout: async () => {
        try {
            await axiosInstance.post("/logout");
            set({ authUser: null });
            toast.success("logged out sucessfully")
        } catch (error) {
            toast.error(error.response.data.message)
        }
    },



    updateProfile: async (data) => {
        set({ isUpdatingProfile: true });
        try {
            const res = await axiosInstance.put("/user/update", data, { headers: { "Content-Type": "multipart/form-data" } });
            set({ authUser: res.data });

            console.log(res.data);
            toast.success("profile updated sucessfully");
        } catch (error) {
            console.error("error in update profile", error);
            toast.error("Failed to update profile");
        } finally{
            set({isUpdatingProfile : false})
        }
    },

}))