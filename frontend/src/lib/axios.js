import axios from "axios";


//Automatically handles JSON parsing, errors, and headers.
const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8000";

export const axiosInstance = axios.create({
  baseURL: API_URL,
  withCredentials: true,
});


