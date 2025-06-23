import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Signup from "./Pages/Signup";
import Home from "./Pages/Home";
import { useThemeStore } from "./Store/useThemeStore";
import Login from "./Pages/Login";
import Profile from "./Pages/Profile";
import Setting from "./Pages/Setting";
import Navbar from "./Components/Navbar";
import { Loader } from 'lucide-react';
import { useEffect } from "react";
import { useAuthStore } from "./Store/useAuthStore"
import { Toaster } from "react-hot-toast";
//import { useThemeStore } from "./Store/useThemeStore";



export default function App() {
  const { checkAuth, authUser, isCheckingAuth } = useAuthStore();
  const {theme} = useThemeStore();
  //  const { theme } = useThemeStore();


  //when we will refresh the page , checkAuth will check whether user is loggedIn. 
     useEffect (() => {
       checkAuth();
     },[checkAuth]) 

     
    console.log({authUser});




 //loader when the user is not logged in (//loader rounding from lucide react)
 if(!authUser && isCheckingAuth)
  return(
    <div className="flex items-center justify-center h-screen ">
      <Loader className="size-10 animate-spin"/>
    </div>
  )



  return (

    <div data-theme={theme}>

      <Navbar />

      <Routes>
        <Route path="/" element={authUser ? <Home /> : <Navigate to="/login" />} />
        <Route path="/signup" element={!authUser ? <Signup /> : <Navigate to="/" />} />
        <Route path="/login" element={!authUser ? <Login /> : <Navigate to="/" />} />
        <Route path="/profile" element={authUser ? <Profile /> : <Navigate to="/login" />} />
        <Route path="/setting" element={<Setting />} />
      </Routes>



      <Toaster />
    </div>
  );
};