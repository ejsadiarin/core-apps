"use client"

import {
  LayoutDashboard,
  Wallet,
  Shirt,
  ShieldCheck,
  Cpu,
  Activity,
  Server,
  Network,
  Globe as GlobeIcon,
  Search,
  Bell,
  Settings,
  Terminal as TerminalIcon,
} from "lucide-react"
import { BentoGrid, BentoCard } from "@/components/ui/bento-grid"
import { Globe } from "@/components/ui/globe"
import { Dock, DockIcon } from "@/components/ui/dock"
import { NeonGradientCard } from "@/components/ui/neon-gradient-card"
import { AnimatedCircularProgressBar } from "@/components/ui/animated-circular-progress-bar"
import { Meteors } from "@/components/ui/meteors"
import { motion } from "motion/react"

const services = [
  {
    name: "Budget Tracker",
    description: "Manage your personal finances and expenses.",
    href: "https://budget.core.local",
    cta: "Open System",
    Icon: Wallet,
    className: "col-span-3 lg:col-span-1",
    background: <div className="absolute inset-0 bg-gradient-to-br from-cyan-500/10 to-transparent" />,
  },
  {
    name: "Wardrobe Manager",
    description: "Catalog and manage your digital wardrobe.",
    href: "https://wardrobe.core.local",
    cta: "Open System",
    Icon: Shirt,
    className: "col-span-3 lg:col-span-1",
    background: <div className="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-transparent" />,
  },
  {
    name: "API Gateway",
    description: "Core traffic management and authentication.",
    href: "#",
    cta: "View Status",
    Icon: ShieldCheck,
    className: "col-span-3 lg:col-span-1",
    background: <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/10 to-transparent" />,
  },
]

export function Dashboard() {
  return (
    <div className="relative min-h-screen bg-[#020617] text-slate-50 font-mono selection:bg-cyan-500/30">
      {/* Meteor effect in dashboard too for fantasy feel */}
      <Meteors number={10} className="opacity-20" />

      {/* Header */}
      <header className="sticky top-0 z-50 flex items-center justify-between border-b border-white/5 bg-[#020617]/80 px-8 py-4 backdrop-blur-md">
        <div className="flex items-center gap-4">
          <motion.div 
            animate={{ rotate: 360 }}
            transition={{ duration: 20, repeat: Infinity, ease: "linear" }}
            className="size-8 rounded-full bg-cyan-500/20 flex items-center justify-center border border-cyan-500/50"
          >
            <Cpu className="size-4 text-cyan-400" />
          </motion.div>
          <div>
            <h1 className="text-sm font-bold tracking-widest text-cyan-400">CORE_SYSTEM <span className="text-purple-500 font-normal opacity-50">// LVL_99</span></h1>
            <p className="text-[10px] text-cyan-500/40">RANK: S-CLASS // STATUS: AWAKENED</p>
          </div>
        </div>

        <div className="flex items-center gap-6">
          <div className="hidden md:flex items-center gap-4 text-[10px] text-cyan-500/40">
            <div className="flex items-center gap-1">
              <Activity className="size-3" /> 14ms
            </div>
            <div className="flex items-center gap-1 text-emerald-400/60">
              <Server className="size-3" /> 98% UP
            </div>
            <div className="flex items-center gap-1">
              <Network className="size-3" /> 1.2 Gbps
            </div>
          </div>
          <div className="flex items-center gap-3">
            <Search className="size-4 text-cyan-500/60 cursor-pointer hover:text-cyan-400 transition-colors" />
            <Bell className="size-4 text-cyan-500/60 cursor-pointer hover:text-cyan-400 transition-colors" />
            <Settings className="size-4 text-cyan-500/60 cursor-pointer hover:text-cyan-400 transition-colors" />
          </div>
        </div>
      </header>

      <div className="mx-auto max-w-7xl px-8 py-8 pb-32">
        <div className="grid grid-cols-12 gap-8">
          {/* Main Stats Column */}
          <div className="col-span-12 lg:col-span-8 space-y-8">
            {/* Hero Globe Section */}
            <div className="relative h-[400px] w-full overflow-hidden rounded-2xl border border-white/5 bg-cyan-950/5 group">
              <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(6,182,212,0.15),transparent_70%)] opacity-0 group-hover:opacity-100 transition-opacity duration-1000" />
              <div className="absolute inset-0 flex flex-col items-center justify-center p-8 text-center z-10 pointer-events-none">
                <motion.h2 
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="text-4xl font-bold tracking-tighter text-cyan-400/80 drop-shadow-[0_0_15px_rgba(6,182,212,0.5)]"
                >
                  REALM_OVERVIEW
                </motion.h2>
                <p className="max-w-xs text-xs text-cyan-500/40 mt-2">Monitoring shadow-network nodes across the physical domain.</p>
              </div>
              <Globe className="top-20 scale-110" />
            </div>

            {/* Applications Grid */}
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="h-px w-8 bg-cyan-500/30" />
                  <h3 className="text-xs font-bold tracking-[0.2em] text-cyan-500/60">CONNECTED_DUNGEONS</h3>
                </div>
                <span className="text-[10px] text-cyan-500/30 underline cursor-pointer hover:text-cyan-400 transition-colors italic">SUMMON NEW NODE</span>
              </div>
              <BentoGrid>
                {services.map((service) => (
                  <BentoCard key={service.name} {...service} />
                ))}
              </BentoGrid>
            </div>
          </div>

          {/* Sidebar Metrics Column */}
          <div className="col-span-12 lg:col-span-4 space-y-8">
            <NeonGradientCard 
              className="min-h-[200px]"
              neonColors={{ firstColor: "#06b6d4", secondColor: "#8b5cf6" }}
            >
              <div className="space-y-6">
                <div className="flex items-center justify-between">
                  <h3 className="text-xs font-bold tracking-[0.2em] text-cyan-400">MANA_LEVELS</h3>
                  <div className="text-[8px] px-1.5 py-0.5 rounded border border-cyan-500/30 text-cyan-500/60">AUTO_REGEN</div>
                </div>
                <div className="flex justify-around">
                  <div className="flex flex-col items-center gap-2">
                    <AnimatedCircularProgressBar
                      value={42}
                      gaugePrimaryColor="#06b6d4"
                      gaugeSecondaryColor="rgba(6, 182, 212, 0.1)"
                      className="size-24 text-sm font-bold"
                    />
                    <span className="text-[10px] text-cyan-500/60 tracking-tighter">PROCESS_POWER</span>
                  </div>
                  <div className="flex flex-col items-center gap-2">
                    <AnimatedCircularProgressBar
                      value={68}
                      gaugePrimaryColor="#8b5cf6"
                      gaugeSecondaryColor="rgba(139, 92, 246, 0.1)"
                      className="size-24 text-sm font-bold"
                    />
                    <span className="text-[10px] text-purple-500/60 tracking-tighter">ESSENCE_STORAGE</span>
                  </div>
                </div>
              </div>
            </NeonGradientCard>

            <div className="rounded-2xl border border-white/5 bg-white/5 p-6 space-y-4 relative overflow-hidden">
              <div className="absolute top-0 right-0 p-2 opacity-10">
                <TerminalIcon className="size-12" />
              </div>
              <h3 className="text-xs font-bold tracking-[0.2em] text-cyan-500/60 uppercase">System_Quest_Logs</h3>
              <div className="space-y-3 font-mono text-[10px]">
                {[
                  { time: "14:20:01", msg: "SHADOW_LOGIN: user eisen verified", type: "success" },
                  { time: "14:19:45", msg: "REALM_WARDROBE: armor_set_04 synchronized", type: "info" },
                  { time: "14:18:12", msg: "GATEWAY: blocking unauthorized malice", type: "warning" },
                  { time: "14:15:30", msg: "HEARTBEAT: all nodes pulse normal", type: "success" },
                ].map((log, i) => (
                  <div key={i} className="flex gap-3 items-start group">
                    <span className="text-cyan-500/20 group-hover:text-cyan-500/40 transition-colors whitespace-nowrap">[{log.time}]</span>
                    <span className={
                      log.type === "success" ? "text-emerald-400/80" : 
                      log.type === "warning" ? "text-amber-400/80" : 
                      "text-cyan-400/70"
                    }>
                      {log.msg}
                    </span>
                  </div>
                ))}
              </div>
            </div>

            <motion.div 
              whileHover={{ scale: 1.02 }}
              className="rounded-2xl border border-cyan-500/20 bg-cyan-500/5 p-6 border-dashed group cursor-pointer"
            >
              <div className="flex items-center gap-3">
                <div className="size-2 rounded-full bg-cyan-400 animate-pulse shadow-[0_0_10px_#22d3ee]" />
                <span className="text-[10px] text-cyan-400 uppercase tracking-widest font-bold group-hover:text-cyan-300 transition-colors">Awaiting Command_</span>
              </div>
            </motion.div>
          </div>
        </div>
      </div>

      {/* Navigation Dock */}
      <div className="fixed bottom-8 left-1/2 -translate-x-1/2 z-50">
        <Dock className="bg-[#020617]/40 border-white/5 backdrop-blur-xl">
          <DockIcon>
            <LayoutDashboard className="size-5 text-cyan-400" />
          </DockIcon>
          <DockIcon>
            <GlobeIcon className="size-5 text-cyan-400/60 hover:text-cyan-400 transition-colors" />
          </DockIcon>
          <DockIcon>
            <ShieldCheck className="size-5 text-cyan-400/60 hover:text-cyan-400 transition-colors" />
          </DockIcon>
          <DockIcon>
            <Settings className="size-5 text-cyan-400/60 hover:text-cyan-400 transition-colors" />
          </DockIcon>
        </Dock>
      </div>
    </div>
  )
}