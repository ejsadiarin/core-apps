"use client";

import { motion } from "motion/react";
import { Server, Database, Cloud, Lock, FileText, Video, Download, Globe, Mail, Code, HardDrive, Wifi, LucideIcon } from "lucide-react";
import { SystemStatus } from "@/components/dashboard/system-status";
import { QuickActions } from "@/components/dashboard/quick-actions";
import { ServiceCard } from "@/components/dashboard/service-card";

const services = [
    {
        icon: Server,
        title: "Docker Manager",
        description: "Container orchestration",
        status: "online" as const,
        metrics: [
            { label: "Containers", value: "24 running" },
            { label: "Images", value: "87 cached" },
        ],
    },
    {
        icon: Database,
        title: "PostgreSQL",
        description: "Primary database cluster",
        status: "online" as const,
        metrics: [
            { label: "Connections", value: "12/100" },
            { label: "Queries/s", value: "234" },
        ],
    },
    {
        icon: Cloud,
        title: "Nextcloud",
        description: "Cloud storage platform",
        status: "online" as const,
        metrics: [
            { label: "Storage", value: "1.2TB used" },
            { label: "Users", value: "5 active" },
        ],
    },
    {
        icon: Lock,
        title: "Vault",
        description: "Secrets management",
        status: "online" as const,
        metrics: [
            { label: "Secrets", value: "156 stored" },
            { label: "Last Seal", value: "2h ago" },
        ],
    },
    {
        icon: FileText,
        title: "Paperless-NGX",
        description: "Document management",
        status: "online" as const,
        metrics: [
            { label: "Documents", value: "1,247" },
            { label: "Tags", value: "34" },
        ],
    },
    {
        icon: Video,
        title: "Jellyfin",
        description: "Media streaming server",
        status: "online" as const,
        metrics: [
            { label: "Movies", value: "342" },
            { label: "TV Shows", value: "87" },
        ],
    },
    {
        icon: Download,
        title: "qBittorrent",
        description: "Download manager",
        status: "online" as const,
        metrics: [
            { label: "Active", value: "3 torrents" },
            { label: "Speed", value: "12.4 MB/s" },
        ],
    },
    {
        icon: Globe,
        title: "Nginx Proxy",
        description: "Reverse proxy manager",
        status: "online" as const,
        metrics: [
            { label: "Hosts", value: "18 proxied" },
            { label: "SSL Certs", value: "12 valid" },
        ],
    },
    {
        icon: Mail,
        title: "Mail Server",
        description: "Email infrastructure",
        status: "maintenance" as const,
        metrics: [
            { label: "Mailboxes", value: "8" },
            { label: "Queue", value: "0 pending" },
        ],
    },
    {
        icon: Code,
        title: "Gitea",
        description: "Git repository hosting",
        status: "online" as const,
        metrics: [
            { label: "Repositories", value: "45" },
            { label: "Commits", value: "2.3K" },
        ],
    },
    {
        icon: HardDrive,
        title: "TrueNAS",
        description: "Storage management",
        status: "online" as const,
        metrics: [
            { label: "Total Space", value: "8TB" },
            { label: "Free Space", value: "2.4TB" },
        ],
    },
    {
        icon: Wifi,
        title: "Pi-hole",
        description: "Network-wide ad blocking",
        status: "online" as const,
        metrics: [
            { label: "Blocked", value: "23.4% ads" },
            { label: "Queries", value: "45.2K today" },
        ],
    },
];

export default function Dashboard() {
    return (
        <div className="container mx-auto px-4 py-8">
            {/* Welcome Banner */}
            <motion.div
                className="mb-8 p-6 rounded-lg border border-primary/30 bg-gradient-to-r from-primary/10 via-accent/10 to-secondary/10 backdrop-blur-sm"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
            >
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-primary mb-2">DASHBOARD ONLINE</h1>
                        <p className="text-sm text-muted-foreground">
                            All critical systems are operational. Last sync: <span className="text-accent">2 minutes ago</span>
                        </p>
                    </div>
                    <div className="hidden md:flex items-center gap-4 text-xs">
                        <div className="flex flex-col items-end">
                            <span className="text-muted-foreground">Uptime</span>
                            <span className="text-card-foreground">42d 13h 27m</span>
                        </div>
                        <div className="w-px h-8 bg-border" />
                        <div className="flex flex-col items-end">
                            <span className="text-muted-foreground">Load Avg</span>
                            <span className="text-card-foreground">1.24, 1.18, 1.05</span>
                        </div>
                    </div>
                </div>
            </motion.div>

            {/* Quick Actions */}
            <div className="mb-8">
                <QuickActions />
            </div>

            {/* System Status */}
            <div className="mb-8">
                <SystemStatus />
            </div>

            {/* Services Grid */}
            <div className="mb-8">
                <motion.div
                    className="flex items-center justify-between mb-6"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.5, delay: 0.2 }}
                >
                    <h2 className="text-primary">ACTIVE SERVICES</h2>
                    <span className="text-xs text-muted-foreground">
                        {services.filter((s) => s.status === "online").length} of {services.length} operational
                    </span>
                </motion.div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                    {services.map((service, idx) => (
                        <ServiceCard key={service.title} {...service} delay={0.1 + idx * 0.05} />
                    ))}
                </div>
            </div>
        </div>
    );
}
