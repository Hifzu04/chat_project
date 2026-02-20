import React, { useState } from 'react'
import { useAuthStore } from '../Store/useAuthStore'
//import Profile_pic from '../assets/Profile_pic.png'
import {  Mail, User, X } from 'lucide-react';
import { useChatStore } from '../Store/useChatStore';

function SelectedUserProfile() {
    
    const { selectedUser } = useChatStore();





    return (
        <div className='h-screen pt-20 '>
            {/* mx-auto: on block elements (like divs), you can use margin auto to center them horizontally within their parent container. */}
            <div className='max-w-3xl mx-auto  p-4 py-8 relative'>
                <div className="bg-base-300 rounded-xl p-8 space-y-8">
                    <div className='text-center'>
                        <h1 className='text-2xl font-semibold'>Profile</h1>
                        <p className='mt-2'>{selectedUser?.fullname || "Selected User"}'s profile information</p>
                    </div>


                    {/* profile pic section */}
                    <div className='items-center gap-4 flex flex-col '>
                        <div className='relative' >
                            <img
                                //checkfordebugging .profilepic
                                src={selectedUser?.profile_pic || "/avatar.png"}
                                alt='profilepic'
                                className='rounded-full border-3 object-cover size-40  '
                            />


                        </div>

                    </div>

                    {/* info section */}
                    <div className='space-y-8'>
                        <div className='space-y-1.5'>
                            <div className='flex items-center gap-2 '>
                                <User className='size-4' />
                                Full Name
                            </div>
                            <p className='border rounded-lg py-2 px-5 bg-base-200 text-zinc-500 '>{selectedUser?.fullname || "Don't you know? :D"}</p>
                        </div>

                        <div className='space-y-1.5'>
                            <div className='flex items-center gap-2 '>
                                <Mail className='size-4' />
                                Email
                            </div>
                            <p className='border rounded-lg py-2 px-5 bg-base-200 text-zinc-500 '>{ "Hahaha... You are not allowed to see any one's email"}</p>
                        </div>
                    </div>


                    <div className='mt-6 bg-base-200 rounded-xl p-6'>
                        <h2 className='text-lg font-medium mb-4 '>Account Information</h2>
                        <div className='space-y-3 text-sm '>
                            <div className='flex items-center justify-between py-2 border-b border-zinc-700' >
                                <span>Member Since</span>
                                <span>
                                    {new Date(selectedUser?.createdAt).toLocaleDateString("en-GB", {
                                        day: "numeric", month: "short", year: "numeric"
                                    })}
                                </span>
                            </div>
                            <div className="flex items-center justify-between py-2">
                                <span>Account Status</span>
                                <span className="text-green-500">Active</span>
                            </div>
                        </div>
                    </div>
                </div>
                 <button className="absolute -top-1.5 -right-1.5 " onClick={() => window.history.back()} > <X /> </button>
            </div>
        
        </div>
    )
}

export default SelectedUserProfile