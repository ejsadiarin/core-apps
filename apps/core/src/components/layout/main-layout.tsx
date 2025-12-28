"use client";

import { motion } from "motion/react";
import { NavigationHeader } from "./navigation-header";

export function MainLayout({
    children,
  }: Readonly<{
    children: React.ReactNode;
  }>) {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <NavigationHeader />
      <motion.main
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5 }}
      >
        {children}
      </motion.main>

      {/* Footer */}
      <motion.footer
        className="mt-12 py-8 border-t border-border/50 text-center text-xs text-muted-foreground"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.8 }}
      >
        <div className="flex items-center justify-center gap-2 mb-2">
          <div className="w-2 h-2 rounded-full bg-accent animate-pulse" />
          <span>CORE SYSTEM v1.0.0 | MATRIX PROTOCOL ACTIVE</span>
        </div>
        <p>Homelab Infrastructure Management System © 2025</p>
      </motion.footer>
    </div>
  );
}
