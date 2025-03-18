import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    outDir: "../static", // กำหนดให้ build ไปที่โฟลเดอร์ static (อยู่นอก frontend)
    emptyOutDir: true, // ลบไฟล์เก่าก่อน build ใหม่
  },
});
