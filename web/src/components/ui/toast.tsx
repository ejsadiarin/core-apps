
import { createContext, useContext, useState, useCallback, ReactNode } from "react";
import { motion, AnimatePresence } from "motion/react";
import { X, Info } from "lucide-react";

interface Toast {
  id: string;
  message: string;
  type: "info" | "warning" | "error" | "success";
}

interface ToastContextType {
  showToast: (message: string, type?: Toast["type"]) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error("useToast must be used within ToastProvider");
  }
  return context;
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const showToast = useCallback((message: string, type: Toast["type"] = "info") => {
    const id = Math.random().toString(36).substring(7);
    setToasts((prev) => [...prev, { id, message, type }]);
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 4000);
  }, []);

  const removeToast = (id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  };

  const getToastStyles = (type: Toast["type"]) => {
    switch (type) {
      case "error":
        return "bg-destructive/90 text-destructive-foreground border-destructive";
      case "warning":
        return "bg-yellow-500/90 text-yellow-foreground border-yellow-500";
      case "success":
        return "bg-green-500/90 text-green-foreground border-green-500";
      default:
        return "bg-primary/90 text-primary-foreground border-primary";
    }
  };

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
        <AnimatePresence>
          {toasts.map((toast) => (
            <motion.div
              key={toast.id}
              initial={{ opacity: 0, y: 20, scale: 0.9 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -20, scale: 0.9 }}
              className={`flex items-center gap-3 px-4 py-3 rounded-lg shadow-lg border backdrop-blur-md ${getToastStyles(toast.type)}`}
            >
              <Info className="w-5 h-5 shrink-0" />
              <p className="text-sm font-medium">{toast.message}</p>
              <button
                onClick={() => removeToast(toast.id)}
                className="ml-2 hover:opacity-70 transition-opacity"
              >
                <X className="w-4 h-4" />
              </button>
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </ToastContext.Provider>
  );
}
