import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Signup from "./Pages/Signup";
import Home from "./Pages/Home";
import Login from "./Pages/Login";
import Profile from "./Pages/Profile";
import Setting from "./Pages/Setting";
import Navbar from "./Components/Navbar";


export default function App() {
  return (
    <div>
    
      <Navbar />

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/login" element={<Login />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/setting" element={<Setting />} />
      </Routes>

    </div>

  );
}

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);