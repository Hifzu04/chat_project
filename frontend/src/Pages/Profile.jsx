import React, { useState } from 'react'
import { useAuthStore } from '../Store/useAuthStore'
//import Profile_pic from '../assets/Profile_pic.png'
import { Camera, Mail, User } from 'lucide-react';

function Profile() {
  const { updateProfile, isUpdatingProfile, authUser } = useAuthStore();
  const [selectedImage, setSelectedImage] = useState();

  const handleImageUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    // 1) keep a preview URL
    const previewURL = URL.createObjectURL(file);
    setSelectedImage(previewURL);

    // 2) build FormData
    const form = new FormData();
    form.append("profilePic", file);

    // 3) call store with FormData
    await updateProfile(form);

  }
  // console.log(authUser.profile_pic);


  return (
    <div className='h-screen pt-20'>
      <div className='max-w-2xl mx-auto p-4 py-8'>
        <div className="bg-base-300 rounded-xl p-6 space-y-8">
          <div className='text-center'>
            <h1 className='text-2xl font-semibold'>Profile</h1>
            <p className='mt-2'>Your profile information</p>
          </div>


          {/* profile pic section */}
          <div className='items-center gap-4 flex flex-col '>
            <div className='relative' >
              <img
                //checkfordebugging .profilepic
                src={selectedImage || authUser.profile_pic || "/avatar.png"}
                alt='profilepic'
                className='rounded-full border-4  object-cover size-32  '
              />
              <label
                htmlFor="avatar-upload"
                className={`
                  absolute bottom-0 right-0 
                  bg-base-content hover:scale-110
                  p-2 rounded-full cursor-pointer 
                  transition-all duration-200
                  ${isUpdatingProfile ? "animate-pulse pointer-events-none" : ""}
                `}
              >
                <Camera className="w-5 h-5 text-base-200" />
                <input
                  type="file"
                  id="avatar-upload"
                  className="hidden"
                  accept="image/*"
                  onChange={handleImageUpload}
                  disabled={isUpdatingProfile}
                />
              </label>
            </div>
            <p className='text-sm text-zinc-500'>{isUpdatingProfile ? "uploading..." : "Click on camera to upload your profile."}</p>
          </div>

          {/* info section */}
          <div className='space-y-8'>
            <div className='space-y-1.5'>
              <div className='flex items-center gap-2 '>
                <User className='size-4' />
                Full Name
              </div>
              <p className='border rounded-lg py-2 px-5 bg-base-200 text-zinc-500 '>{authUser?.fullname || "Don't you know? :D"}</p>
            </div>

            <div className='space-y-1.5'>
              <div className='flex items-center gap-2 '>
                <Mail className='size-4' />
                Email
              </div>
              <p className='border rounded-lg py-2 px-5 bg-base-200 text-zinc-500 '>{authUser?.email || "huhhh..Your email?? think about it man"}</p>
            </div>
          </div>

          <div className='mt-6 bg-base-200 rounded-xl p-6'>
            <h2 className='text-lg font-medium mb-4 '>Account Information</h2>
            <div className='space-y-3 text-sm '>
              <div className='flex items-center justify-between py-2 border-b border-zinc-700' >
                <span>Member Since</span>
                <span>
                  {new Date(authUser.createdAt).toLocaleDateString("en-GB", {
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
      </div>
    </div>
  )
}

export default Profile