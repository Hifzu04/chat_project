//This code creates a reusable Axios instance that simplifies making HTTP requests to your backend API
//  throughout your React application



import axios from "axios";

export const axiosInstance = axios.create({
    baseURL: "http://localhost:8000",
    

    //if error : check using baseurl  baseURL: "http://localhost:5001",
    withCredentials: true
})


//baseURL: All requests made using this instance will be prefixed with http://localhost:5001/api.

//withCredentials: true: Allows cookies and authorization headers to be sent along with requests
// (useful for sessions or JWT stored in cookies).