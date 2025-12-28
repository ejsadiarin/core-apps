"use client";

import { createContext, useContext, useState, ReactNode, useEffect } from "react";

export type UserRole = "admin" | "user" | "viewer";

export interface User {
    id: string;
    username: string;
    email: string;
    role: UserRole;
    avatar?: string;
}

interface AuthContextType {
    user: User | null;
    login: (username: string, password: string) => Promise<boolean>;
    logout: () => void;
    isAuthenticated: boolean;
    hasPermission: (requiredRole: UserRole) => boolean;
    isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Mock users for demo
const MOCK_USERS: Record<string, { password: string; user: User }> = {
    admin: {
        password: "admin123",
        user: {
            id: "1",
            username: "admin",
            email: "admin@homelab.local",
            role: "admin",
        },
    },
    user: {
        password: "user123",
        user: {
            id: "2",
            username: "user",
            email: "user@homelab.local",
            role: "user",
        },
    },
    viewer: {
        password: "viewer123",
        user: {
            id: "3",
            username: "viewer",
            email: "viewer@homelab.local",
            role: "viewer",
        },
    },
};

const ROLE_HIERARCHY: Record<UserRole, number> = {
    admin: 3,
    user: 2,
    viewer: 1,
};

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    const login = async (username: string, password: string): Promise<boolean> => {
        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 800));

        const mockUser = MOCK_USERS[username.toLowerCase()];
        if (mockUser && mockUser.password === password) {
            setUser(mockUser.user);
            localStorage.setItem("user", JSON.stringify(mockUser.user));
            return true;
        }
        return false;
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem("user");
    };

    const hasPermission = (requiredRole: UserRole): boolean => {
        if (!user) return false;
        return ROLE_HIERARCHY[user.role] >= ROLE_HIERARCHY[requiredRole];
    };

    useEffect(() => {
        const storedUser = localStorage.getItem("user");
        if (storedUser) {
            try {
                setUser(JSON.parse(storedUser));
            } catch (e) {
                localStorage.removeItem("user");
            }
        }
        setIsLoading(false);
    }, []);

    return (
        <AuthContext.Provider
            value={{
                user,
                login,
                logout,
                isAuthenticated: !!user,
                hasPermission,
                isLoading,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
