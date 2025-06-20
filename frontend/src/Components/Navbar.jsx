import React from 'react'
import { useAuthStore } from '../Store/useAuthStore'
import { Link } from 'react-router-dom';
import { LogOut, LogOutIcon, MessageCircle, MessageSquare, Settings, User } from 'lucide-react';
import toast from 'react-hot-toast';

//logoandAppname(clickabletohomepage) ,settings 
//when authenticated : profile and logout 

function Navbar() {
  const { logout, authUser } = useAuthStore();
  return (
    <div className="navbar bg-base-100 shadow-sm  justify-between">
      <div className="navbar-start px-4">
        <div className="dropdown">
          <div tabIndex={0} role="button" className="btn btn-ghost btn-circle">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"> <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h7" /> </svg>
          </div>
          <ul
            tabIndex={0}
            className="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow">
            <li>  <Link to="/">
              Home
            </Link></li>
            <li><Link to={"/setting"}>Settings</Link></li>

          </ul>
        </div>
      </div>
      <div className="navbar-center  ">
        <Link to={"/"} className="btn btn-ghost text-xl">ChatterNest</Link>
      </div>

      {authUser && (
      <div className="navbar-end gap-8 px-6 ">
        
          <Link to={"/profile"} className='flex items-center gap-1' >
            <User className='size-5' />
            <span className='hidden md:inline'>Profile</span>

          </Link>
      
        <button className="flex items-center gap-1 " onClick={logout}>
          <LogOut className='size-5'/>
          <span className='hidden md:inline'>Logout</span>
        </button>
      </div>

      )}
    </div>
  )
}

export default Navbar