import React, { useState } from 'react'

import { Eye, EyeOff, Loader2 } from 'lucide-react';
//import Meetme from '../assets/Meetme.png'
import { Link } from 'react-router-dom';

import { useAuthStore } from '../Store/useAuthStore';
import { useNavigate } from "react-router-dom";

import toast, { Toaster } from 'react-hot-toast';




function Login() {
      const navigate = useNavigate();
    const [showPassword , setShowPassword] = useState(false);
    const [formData, setFormData] = useState({
        email: "",
        password: "",

    });

    const { login, isLoggingin } = useAuthStore();


    const handleSubmit = async  (e) => {
       e.preventDefault();

        const sucess = await login(formData);

        // if (sucess){
        //     navigate("/");   
               
        // }

    }



    return (

        <div className="min-h-screen flex flex-col md:flex-row">
            {/* {left panel} */}
            <div className="w-full  md:w-1/2 flex items-center justify-center p-8 bg-base-300">

                <div className="max-w-md w-full space-y-6">

                    {/* logo */}
                    <div className="flex justify-center">
                        <img src={"/Meetme.png"} alt="App Logo" className="h-12 w-auto" />
                    </div>
                    <h2 className="mt-6 text-center text-3xl font-extrabold ">
                        Welcome Back
                    </h2>
                    <h1 className='text-center text-bold text-base-content/70'>Sign in to your account</h1>

                    <form className="mt-8 space-y-4" onSubmit={handleSubmit}>
                        <div>
                            <label className="block text-sm font-medium">
                                Email address
                            </label>
                            <input
                                name="email"
                                type="email"
                                placeholder='you@gmail.com'
                                required
                                value={formData.email}
                                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                className="mt-1 block w-full px-3 py-2 border rounded-lg shadow-sm  focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            />
                        </div>
                        <div className="relative">
                            <label className="block text-sm font-medium">
                                Password
                            </label>
                            <input
                                name="password"
                                placeholder='Enter your password'
                                type={showPassword ? 'text' : 'password'}
                                required
                                value={formData.password}
                                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                                className="mt-1 block w-full px-3 py-3 border rounded-lg shadow-sm  focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute inset-y-0 right-0 pr-4 pt-4 flex items-center text-sm "
                            >
                                {showPassword ? (
                                    <EyeOff className="size-5 text-base-content/40" />
                                ) : (
                                    <Eye className='size-5 text-base-content/40' />
                                )}
                            </button>
                        </div>
                        <button
                            type="submit"
                            disabled={isLoggingin}
                            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm  bg-indigo-600 text-base-300 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                        >
                            {isLoggingin ? (
                                <>
                                    <Loader2 className='size-5 animate-spin' />
                                    Loading...
                                </>
                            ) : ("log In")}
                        </button>
                    </form>


                    <p className="mt-4 text-center text-sm ">
                        New to  ChatterNest?{' '}
                        <Link to="/signup" className="link link-primary ">
                            Go to Signup
                        </Link>
                    </p>


                </div>
            </div>

            {/* Right panel */}
            <div className="hidden md:flex w-1/2 items-center justify-center bg-gradient-to-br from-indigo-500 to-purple-600 text-white p-8">
                <div className="max-w-lg text-center space-y-6 ">
                    <blockquote className='text-9xl font-bold text '>
                        ChatNest
                    </blockquote>
                    <blockquote className="text-lg  italic  ">
                        “Pop in, say hi, and keep the convo rolling!" ☕💬
                        {/* <footer className="mt-4 text-sm">— Steve Jobs</footer> */}
                    </blockquote>
                </div>
            </div>

        </div>


    )
}

export default Login


