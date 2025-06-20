//import meetme from '.../assets/meetme.png'
import  { useState } from 'react';
import { Eye, EyeOff, Loader2 } from 'lucide-react';
import Meetme from '../assets/Meetme.png'
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { useAuthStore } from '../Store/useAuthStore';

export default function Signup() {
  const [showPassword, setShowPassword] = useState(false);
  const [formData, setFormData] = useState({
    fullname: '',
    email: '',
    password: '',
  });
  const { signup, isSigningup } = useAuthStore()

  const validateForm = () => {
        if(!formData.fullname.trim()) return toast.error("full name is required");
        if(!formData.email.trim()) return toast.error("email is required");
        if (!/\S+@\S+\.\S+/.test(formData.email)) return toast.error("invalid email address");
        if(!formData.password.trim()) return toast.error("password is required");
        if(formData.password.length<6 ) return toast.error("Password must be atleast 6 character");

        return true;
  }

  
  // }
  const handleSubmit = e => {

    e.preventDefault();
    
    // TODO: add your validation here
    
    const sucesss = validateForm()

    if(sucesss===true){
      signup(formData)
    }
  };
  return (
    <div className="min-h-screen flex flex-col md:flex-row">
      {/* {left panel} */}
      <div className="w-full  md:w-1/2 flex items-center justify-center p-8 bg-base-200">

        <div className="max-w-md w-full space-y-6">

          {/* logo */}
          <div className="flex justify-center">
            <img src={Meetme} alt="App Logo" className="h-12 w-auto" />
          </div>
          <h2 className="mt-6 text-center text-3xl font-extrabold ">
            Create your account
          </h2>

          <form className="mt-8 space-y-4" onSubmit={handleSubmit}>
            <div>
              <label className="block text-sm font-medium ">
                Full Name
              </label>
              <input
                name="fullname"
                placeholder='Enter your name'
                type="text"
                required
                value={formData.fullname}
                onChange={(e) => setFormData({ ...formData, fullname: e.target.value })}
                className="mt-1 block w-full px-3 py-2 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
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
                className="mt-1 block w-full px-3 py-2 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div className="relative">
              <label className="block text-sm font-medium ">
                Password
              </label>
              <input
                name="password"
                placeholder='Enter your password'
                type={showPassword ? 'text' : 'password'}
                required
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                className="mt-1 block w-full px-3 py-3 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute inset-y-0 right-0 pr-4 pt-4 flex items-center text-sm text-gray-300"
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
              disabled={isSigningup}
              className="w-full flex justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              {isSigningup ? (
                <>
                  <Loader2 className='size-5 animate-spin' />
                  Loading...
                </>
              ) : ("Create Account")}
            </button>
          </form>


          <p className="mt-4 text-center text-sm ">
            Already registered?{' '}
            <Link to="/login" className="link link-primary ">
              Go to login
            </Link>
          </p>


        </div>
      </div>

      {/* Right panel */}
      <div className="hidden md:flex w-1/2 items-center justify-center bg-gradient-to-br from-indigo-500 to-purple-600 text-white p-8">
        <div className="max-w-lg text-center space-y-6 ">
         <blockquote className='text-8xl  font-bold '>
           ChatterNest
         </blockquote>
          <blockquote className="text-xl italic font-light ">
            “Pop in, say hi, and keep the convo rolling!" ☕💬    
            {/* <footer className="mt-4 text-sm">— Steve Jobs</footer> */}
          </blockquote>
        </div>
      </div>

    </div>
  )
};