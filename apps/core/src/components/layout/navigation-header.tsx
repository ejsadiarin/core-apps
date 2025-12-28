"use client";

import { motion } from "motion/react";
import { Terminal, Settings, Bell, User, LogOut } from "lucide-react";
import { usePathname, useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useAuth } from "@/contexts/auth-context";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function NavigationHeader() {
  const router = useRouter();
  const pathname = usePathname();
  const { user, logout } = useAuth();

  const handleLogout = () => {
    logout();
    router.push("/");
  };
  return (
    <motion.header
      className="sticky top-0 z-50 border-b border-border bg-background/80 backdrop-blur-lg"
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      transition={{ duration: 0.5 }}
    >
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          {/* Logo Section */}
          <div className="flex items-center gap-3">
            <div className="relative">
              <Terminal className="w-6 h-6 text-primary glow-text" />
              <motion.div
                className="absolute inset-0 border-2 border-primary rounded"
                animate={{
                  scale: [1, 1.2, 1],
                  opacity: [0.5, 0, 0.5],
                }}
                transition={{
                  duration: 2,
                  repeat: Infinity,
                }}
              />
            </div>
            <div>
              <h1 className="text-sm text-primary">CORE SYSTEM</h1>
              <p className="text-xs text-muted-foreground">Homelab Control Center</p>
            </div>
          </div>

          {/* Navigation Links */}
          <nav className="hidden md:flex items-center gap-1">
            {[
              { name: "Dashboard", path: "/dashboard" },
              { name: "Services", path: "/services" },
              { name: "Analytics", path: "/analytics" },
              { name: "Logs", path: "/logs" },
              ...(user?.role === "admin" ? [{ name: "Admin", path: "/admin" }] : []),
            ].map((item, idx) => (
              <motion.div
                key={item.name}
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: idx * 0.1 }}
              >
                <Button
                  variant="ghost"
                  onClick={() => router.push(item.path)}
                  className={`text-sm transition-colors ${
                    pathname === item.path
                      ? "text-primary bg-primary/10"
                      : item.name === "Admin"
                      ? "text-destructive hover:text-destructive hover:bg-destructive/10"
                      : "text-muted-foreground hover:text-primary hover:bg-primary/10"
                  }`}
                >
                  {item.name}
                </Button>
              </motion.div>
            ))}
          </nav>

          {/* Action Buttons */}
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="relative text-muted-foreground hover:text-primary hover:bg-primary/10"
            >
              <Bell className="w-4 h-4" />
              <Badge className="absolute -top-1 -right-1 w-4 h-4 p-0 flex items-center justify-center bg-destructive text-xs border-0">
                3
              </Badge>
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="text-muted-foreground hover:text-primary hover:bg-primary/10"
            >
              <Settings className="w-4 h-4" />
            </Button>
            <div className="w-px h-6 bg-border mx-2" />
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-muted-foreground hover:text-primary hover:bg-primary/10"
                >
                  <User className="w-4 h-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel>
                  <div className="flex flex-col">
                    <span className="text-sm">{user?.username}</span>
                    <span className="text-xs text-muted-foreground">{user?.email}</span>
                    <Badge className="mt-2 w-fit text-xs bg-primary/10 text-primary border-primary/30">
                      {user?.role.toUpperCase()}
                    </Badge>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => router.push("/dashboard")}>
                  <Terminal className="w-4 h-4 mr-2" />
                  Dashboard
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Settings className="w-4 h-4 mr-2" />
                  Settings
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout} className="text-destructive">
                  <LogOut className="w-4 h-4 mr-2" />
                  Logout
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      {/* Accent Line */}
      <motion.div
        className="h-px bg-gradient-to-r from-transparent via-primary to-transparent"
        initial={{ opacity: 0 }}
        animate={{ opacity: [0.3, 0.6, 0.3] }}
        transition={{ duration: 3, repeat: Infinity }}
      />
    </motion.header>
  );
}
