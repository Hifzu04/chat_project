import React from 'react'

import toast, { Toaster } from 'react-hot-toast';

const notify = () => toast('Here is your toast.');

function Login() {
    return (
        <div>
            <div>Login</div>
            <div>
                <button className='bg-green-600' onClick={notify}>Make me a toast</button>
                <Toaster />
            </div>
        </div>
    )
}

export default Login


