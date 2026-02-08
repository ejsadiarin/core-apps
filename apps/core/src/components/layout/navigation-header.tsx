"use client";

import { motion } from "motion/react";
import { Terminal, Settings, Bell, User, LogOut, Shield } from "lucide-react";
import { usePathname, useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
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
  const queryClient = useQueryClient();
  const { user, logout, isGuest, isAdmin } = useAuth();

  const handleLogout = async () => {
    await logout();
    queryClient.clear(); // Clear React Query cache to prevent data leakage between users
    router.push("/login");
  };

  const getRoleBadgeColor = () => {
    switch (user?.role) {
      case "admin":
        return "bg-destructive/10 text-destructive border-destructive/30";
      case "user":
        return "bg-primary/10 text-primary border-primary/30";
      case "guest":
        return "bg-muted text-muted-foreground border-muted";
      default:
        return "bg-primary/10 text-primary border-primary/30";
    }
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
              { name: "Budget", path: "/dashboard/budget" },
              { name: "Services", path: "/dashboard/services" },
              { name: "Analytics", path: "/analytics" },
              { name: "Logs", path: "/logs" },
              ...(isAdmin ? [{ name: "Users", path: "/dashboard/admin/users", isAdmin: true }] : []),
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
                      : "isAdmin" in item && item.isAdmin
                      ? "text-destructive hover:text-destructive hover:bg-destructive/10"
                      : "text-muted-foreground hover:text-primary hover:bg-primary/10"
                  }`}
                >
                  {"isAdmin" in item && item.isAdmin && <Shield className="w-3 h-3 mr-1" />}
                  {item.name}
                </Button>
              </motion.div>
            ))}
          </nav>

          {/* Action Buttons */}
          <div className="flex items-center gap-2">
            {isGuest && (
              <Badge variant="outline" className="text-xs bg-accent/10 text-accent border-accent/30">
                Guest Mode
              </Badge>
            )}
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
                    <span className="text-sm">{user?.email}</span>
                    <Badge className={`mt-2 w-fit text-xs ${getRoleBadgeColor()}`}>
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
                {isAdmin && (
                  <DropdownMenuItem onClick={() => router.push("/dashboard/admin/users")}>
                    <Shield className="w-4 h-4 mr-2" />
                    User Management
                  </DropdownMenuItem>
                )}
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
