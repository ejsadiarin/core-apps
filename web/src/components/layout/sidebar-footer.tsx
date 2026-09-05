import { motion, AnimatePresence } from "motion/react";
import { LogOut, User, ChevronUp } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { cn } from "@/lib/utils";
import { useAuth } from "@/contexts/auth-context";
import { useSidebar } from "./sidebar-context";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

function getRoleBadgeClasses(role?: string) {
  switch (role) {
    case "admin":
      return "bg-destructive/10 text-destructive border-destructive/30";
    case "user":
      return "bg-primary/10 text-primary border-primary/30";
    case "guest":
      return "bg-accent/10 text-accent border-accent/30";
    default:
      return "bg-primary/10 text-primary border-primary/30";
  }
}

export function SidebarFooter() {
  const { user, logout, isGuest } = useAuth();
  const { isCollapsed } = useSidebar();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const handleLogout = async () => {
    await logout();
    queryClient.clear();
    navigate("/login");
  };

  // collapsed: show icon with dropdown
  if (isCollapsed) {
    return (
      <div className="border-t border-sidebar-border p-2">
        <DropdownMenu>
          <Tooltip delayDuration={0}>
            <TooltipTrigger asChild>
              <DropdownMenuTrigger asChild>
                <button className="flex w-full items-center justify-center rounded-md p-2 text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground">
                  <div className="relative">
                    <User className="h-[18px] w-[18px]" />
                    {isGuest && (
                      <div className="absolute -top-0.5 -right-0.5 h-2 w-2 rounded-full bg-accent" />
                    )}
                  </div>
                </button>
              </DropdownMenuTrigger>
            </TooltipTrigger>
            <TooltipContent side="right" className="font-mono text-xs">
              {user?.email ?? "User"}
            </TooltipContent>
          </Tooltip>
          <DropdownMenuContent side="right" align="end" className="w-56 mb-2">
            <DropdownMenuLabel>
              <div className="flex flex-col gap-1">
                <span className="text-xs text-muted-foreground truncate">{user?.email}</span>
                <div className="flex items-center gap-2">
                  <Badge className={cn("text-[10px]", getRoleBadgeClasses(user?.role))}>
                    {user?.role?.toUpperCase()}
                  </Badge>
                  {isGuest && (
                    <Badge variant="outline" className="text-[10px] bg-accent/10 text-accent border-accent/30">
                      GUEST
                    </Badge>
                  )}
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleLogout} className="text-destructive">
              <LogOut className="mr-2 h-3.5 w-3.5" />
              Logout
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    );
  }

  // expanded: show full user info
  return (
    <div className="border-t border-sidebar-border p-3">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button className="flex w-full items-center gap-3 rounded-md p-2 text-left transition-colors hover:bg-sidebar-accent group">
            <div className="relative flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-muted/50 border border-sidebar-border">
              <User className="h-4 w-4 text-muted-foreground" />
              {isGuest && (
                <div className="absolute -top-0.5 -right-0.5 h-2 w-2 rounded-full bg-accent" />
              )}
            </div>
            <div className="flex-1 min-w-0">
              <p className="truncate text-xs text-sidebar-foreground">
                {user?.email ?? "Unknown"}
              </p>
              <div className="flex items-center gap-1.5 mt-0.5">
                <Badge className={cn("text-[10px] px-1.5 py-0", getRoleBadgeClasses(user?.role))}>
                  {user?.role?.toUpperCase()}
                </Badge>
                <AnimatePresence>
                  {isGuest && (
                    <motion.div
                      initial={{ opacity: 0, scale: 0.8 }}
                      animate={{ opacity: 1, scale: 1 }}
                      exit={{ opacity: 0, scale: 0.8 }}
                    >
                      <Badge variant="outline" className="text-[10px] px-1.5 py-0 bg-accent/10 text-accent border-accent/30">
                        GUEST
                      </Badge>
                    </motion.div>
                  )}
                </AnimatePresence>
              </div>
            </div>
            <ChevronUp className="h-4 w-4 text-muted-foreground/50 group-hover:text-muted-foreground transition-colors" />
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent side="top" align="start" className="w-[calc(var(--radix-dropdown-menu-trigger-width))] mb-1">
          <DropdownMenuItem onClick={handleLogout} className="text-destructive">
            <LogOut className="mr-2 h-3.5 w-3.5" />
            Logout
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
