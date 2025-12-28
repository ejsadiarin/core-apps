"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "motion/react";
import { Terminal, Lock, User, AlertCircle, Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useAuth } from "@/contexts/auth-context";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      const success = await login(username, password);
      if (success) {
        router.push("/dashboard");
      } else {
        setError("Invalid credentials. Please try again.");
      }
    } catch (err) {
      setError("An error occurred. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center relative overflow-hidden px-4">
      {/* Animated Background */}
      <div className="absolute inset-0 opacity-30">
        <div
          className="absolute inset-0"
          style={{
            backgroundImage: `
              linear-gradient(rgba(0, 229, 204, 0.1) 1px, transparent 1px),
              linear-gradient(90deg, rgba(0, 229, 204, 0.1) 1px, transparent 1px)
            `,
            backgroundSize: "50px 50px",
            animation: "grid-glow 4s ease-in-out infinite",
          }}
        />
      </div>

      {/* Scanning Line */}
      <div className="absolute inset-0 pointer-events-none overflow-hidden">
        <div className="scan-line absolute w-full h-px bg-gradient-to-r from-transparent via-primary to-transparent opacity-50" />
      </div>

      {/* Login Card */}
      <motion.div
        className="relative z-10 w-full max-w-md"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <div className="rounded-lg border border-border bg-card/80 backdrop-blur-xl p-8 glow-box">
          {/* Header */}
          <div className="text-center mb-8">
            <motion.div
              className="inline-flex items-center justify-center w-16 h-16 rounded-lg bg-primary/10 border border-primary/30 mb-4"
              animate={{
                boxShadow: [
                  "0 0 20px rgba(0, 229, 204, 0.3)",
                  "0 0 40px rgba(0, 229, 204, 0.5)",
                  "0 0 20px rgba(0, 229, 204, 0.3)",
                ],
              }}
              transition={{ duration: 2, repeat: Infinity }}
            >
              <Terminal className="w-8 h-8 text-primary" />
            </motion.div>
            <h1 className="text-primary mb-2">SECURE ACCESS</h1>
            <p className="text-sm text-muted-foreground">
              Enter your credentials to access the Core System
            </p>
          </div>

          {/* Demo Credentials Info */}
          <div className="mb-6 p-3 rounded-lg bg-accent/10 border border-accent/30">
            <p className="text-xs text-accent mb-2">Demo Credentials:</p>
            <div className="text-xs text-muted-foreground space-y-1">
              <p>Admin: admin / admin123</p>
              <p>User: user / user123</p>
              <p>Viewer: viewer / viewer123</p>
            </div>
          </div>

          {/* Error Alert */}
          {error && (
            <motion.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              className="mb-6"
            >
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            </motion.div>
          )}

          {/* Login Form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="username" className="text-card-foreground">
                Username
              </Label>
              <div className="relative">
                <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="pl-10 bg-input-background border-border focus:border-primary"
                  placeholder="Enter username"
                  required
                  autoComplete="username"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="password" className="text-card-foreground">
                Password
              </Label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="pl-10 pr-10 bg-input-background border-border focus:border-primary"
                  placeholder="Enter password"
                  required
                  autoComplete="current-password"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-primary transition-colors"
                >
                  {showPassword ? (
                    <EyeOff className="w-4 h-4" />
                  ) : (
                    <Eye className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>

            <Button
              type="submit"
              className="w-full bg-primary text-primary-foreground hover:bg-primary/90 relative overflow-hidden"
              disabled={isLoading}
            >
              {isLoading ? (
                <motion.div
                  className="flex items-center gap-2"
                  animate={{ opacity: [1, 0.5, 1] }}
                  transition={{ duration: 1.5, repeat: Infinity }}
                >
                  <div className="w-4 h-4 border-2 border-primary-foreground border-t-transparent rounded-full animate-spin" />
                  <span>Authenticating...</span>
                </motion.div>
              ) : (
                "Access System"
              )}
            </Button>
          </form>

          {/* Footer */}
          <div className="mt-6 pt-6 border-t border-border/50 text-center">
            <p className="text-xs text-muted-foreground">
              Authorized personnel only. All access is logged and monitored.
            </p>
          </div>
        </div>

        {/* Status Indicator */}
        <motion.div
          className="mt-4 flex items-center justify-center gap-2 text-xs text-muted-foreground"
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{ duration: 2, repeat: Infinity }}
        >
          <div className="w-2 h-2 rounded-full bg-accent" />
          <span>SECURE CONNECTION ESTABLISHED</span>
        </motion.div>
      </motion.div>

      {/* Corner Decorations */}
      <div className="absolute top-8 left-8 w-16 h-16 border-l-2 border-t-2 border-primary/50" />
      <div className="absolute top-8 right-8 w-16 h-16 border-r-2 border-t-2 border-primary/50" />
      <div className="absolute bottom-8 left-8 w-16 h-16 border-l-2 border-b-2 border-primary/50" />
      <div className="absolute bottom-8 right-8 w-16 h-16 border-r-2 border-b-2 border-primary/50" />
    </div>
  );
}
