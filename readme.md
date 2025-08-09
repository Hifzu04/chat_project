# 📬 ChatNest

[![Live Demo](https://img.shields.io/badge/Live%20Demo-Click%20Here-brightgreen?style=for-the-badge&logo=vercel)](https://chatnest-3y2m.onrender.com/)

[![Tech Stack](https://img.shields.io/badge/Stack-Go%20|%20React%20|%20WebSockets%20|%20MongoDB-blueviolet?style=for-the-badge)](#-tech-stack)

**ChatNest** is a full-featured, production-ready **real-time chat application** built with **Go (Gorilla Mux)**, **React**, **WebSockets**, and **MongoDB**.  
It supports **secure authentication**, **real-time messaging**, **image sharing**, **user search**, **theme switching**, and more — all deployed on **Render**.

---



---

## 📸 Features

✅ **Real-time messaging** — Instant text & image communication via WebSockets.  
✅ **Authentication & Authorization** — Secure JWT tokens in httpOnly cookies.  
✅ **Online presence tracking** — See who's online + "Show Online Only" toggle.  
✅ **User search** — Quickly find other users.  
✅ **Profile management** — Update username, avatar, and theme preferences.  
✅ **Theme changer** — Light/dark + custom themes.  
✅ **Responsive UI** — Mobile & desktop friendly.  
✅ **Production-ready deployment** — Fully hosted on Render.  

---

## 🛠 Tech Stack

### 💻 Frontend
- ⚛️ **React 19**
- ⚡ **Vite**
- 🌊 **TailwindCSS** + 🎨 **DaisyUI**
- 🐻 **Zustand** (state management)
- 🛣 **React Router DOM**
- 🔌 **Socket.IO Client**
- 📡 **Axios**
- 🎯 **Lucide-react** (icons)
- 🔔 **React-hot-toast** (notifications)

### 🛠 Backend
- 🐹 **Go (Gorilla Mux)** — REST API + WebSockets
- 🔑 **JWT authentication** (httpOnly cookies) & authorization
- 🍃 **MongoDB**
- ☁️ **Cloudinary**
- ⚙️ **godotenv**
- 🔒 **x/crypto**
- 🌐 **CORS handling**
- 🆔 **UUID generation**

---

## 🏗 Architecture

```plaintext
          ┌───────────────────┐
          │   React Frontend   │
          │ (Vite, Zustand,    │
          │  Socket.IO Client) │
          └─────────┬─────────┘
                    │ HTTP (REST) + WebSockets
                    ▼
        ┌────────────────────────┐
        │    Go Backend API       │
        │ (Gorilla Mux, JWT Auth, │
        │   Socket.IO Server)     │
        └─────────┬──────────────┘
                  │
        ┌─────────┼─────────┐
        ▼                   ▼
  ┌──────────────┐   ┌───────────────┐
  │   MongoDB     │   │  Cloudinary   │
  │ (Users, Msgs) │   │ (Image Store) │
  └──────────────┘   └───────────────┘
