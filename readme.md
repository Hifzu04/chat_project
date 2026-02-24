# 📬 ChatNest 

<h3 align="center">High-Performance Real-Time Messaging with AI Intelligence</h3>

<p align="center">
  <a href="https://chatnest-3y2m.onrender.com/"><strong>Live Demo »</strong></a>
  ·
  <a href="https://github.com/Hifzu04/chat_project/issues">Report Bug</a>
  ·
  <a href="https://github.com/Hifzu04/chat_project/pulls">Request Feature</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Maintained%3F-yes-green.svg?style=for-the-badge" />
  <img src="https://img.shields.io/badge/AI--Enhanced-NestMind-blueviolet?style=for-the-badge&logo=openai" />
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go" />
</p>


[![Tech Stack](https://img.shields.io/badge/Stack-Go%20|%20React%20|%20WebSockets%20|%20MongoDB-blueviolet?style=for-the-badge)](#-tech-stack)

## 📖 Overview

**ChatNest** is a full-stack messaging ecosystem designed for the modern era. By combining **Go's** legendary concurrency with a reactive **React 19** frontend, it delivers a sub-100ms messaging experience. It doesn't just send messages—it understands them through its integrated **AI** layer.

----

## 🚀 Key Features      

### 🧠 The AI Edge (NestMind)
* **Smart Replies:** Context-aware suggestions that learn from your conversation style.
* **Instant Summaries:** One-click "catch-up" feature that summarizes long threads you've missed.
* **Sentiment Analysis:** Visual indicators of the conversation's mood (Experimental).
* **AI Media Assist:** Automatic image captioning and descriptive alt-text for accessibility.

### 💬 Real-Time Messaging
* **Full-Duplex WebSockets:** Instant text and image delivery via Socket.io.
* **Presence Tracking:** Real-time online status and typing indicators.
* **Notification System:** Instant alerts for new messages via React-hot-toast.

### 🔐 Security & Infrastructure
* **JWT Security:** Authentication handled via secure, `httpOnly` cookies to mitigate XSS risks.
* **Media Hosting:** Robust image storage and optimization via Cloudinary CDN.
* **Responsive UI:** 15+ themes (Light, Dark, Retro, etc.) optimized for every device.
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
- 🤖 **Google gemini API** - AI powerd bot assistant 
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
          │   React Frontend   │ ◄────┐
          │ (Vite + Zustand)  │      │
          └─────────┬─────────┘      │
                    │ REST / WSS     │ AI Logic
                    ▼                │
        ┌────────────────────────┐   │
        │    Go Backend API       │───┘
        │ (Gorilla Mux + JWT)    │
        └─────────┬──────────────┘
                  │
        ┌─────────┼──────────────┐
        ▼         ▼              ▼
  ┌──────────┐ ┌──────────┐ ┌──────────────┐
  │ MongoDB  │ │Cloudinary│ │  OpenAI/LLM  │
  │ (Data)   │ │ (Images) │ │  (Gemini-API)  │
  └──────────┘ └──────────┘ └──────────────┘
----


