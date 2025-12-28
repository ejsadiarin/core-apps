import type { Metadata } from "next";
import { Geist_Mono as GeistMono } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/contexts/auth-context";

const geistMono = GeistMono({
  variable: "--font-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "CORE System",
  description: "Homelab Management Platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body
        className={`${geistMono.variable} font-mono antialiased`}
      >
        <AuthProvider>
            {children}
        </AuthProvider>
      </body>
    </html>
  );
}