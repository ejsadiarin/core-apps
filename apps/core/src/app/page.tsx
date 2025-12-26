"use client"

import { useState } from "react"
import { motion, AnimatePresence } from "motion/react"
import { Particles } from "@/components/ui/particles"
import { RetroGrid } from "@/components/ui/retro-grid"
import { HyperText } from "@/components/ui/hyper-text"
import { ShimmerButton } from "@/components/ui/shimmer-button"
import { Meteors } from "@/components/ui/meteors"
import { Dashboard } from "@/components/dashboard/dashboard-view"

export default function Home() {
  const [isEntered, setIsEntered] = useState(false)

  return (
    <main className="relative min-h-screen w-full overflow-hidden bg-background">
      <AnimatePresence mode="wait">
        {!isEntered ? (
          <motion.div
            key="entrance"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0, scale: 1.1, filter: "blur(10px)" }}
            transition={{ duration: 0.8 }}
            className="relative flex min-h-screen flex-col items-center justify-center p-4 bg-[radial-gradient(ellipse_at_center,rgba(6,182,212,0.05),transparent_50%),radial-gradient(ellipse_at_bottom,rgba(34,197,94,0.05),transparent_50%)]"
          >
            {/* Background Effects */}
            <RetroGrid opacity={0.2} darkLineColor="#06b6d4" />
            <Particles
              className="absolute inset-0 z-0"
              quantity={200}
              color="#22d3ee"
              staticity={30}
            />
            <Meteors number={12} />

            {/* Content */}
            <div className="z-10 flex flex-col items-center gap-8 text-center">
              <div className="space-y-2">
                <HyperText
                  className="text-7xl font-bold tracking-tighter text-primary-glow sm:text-8xl md:text-9xl"
                  duration={1200}
                >
                  CORE
                </HyperText>
                <motion.p
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.5, duration: 0.8 }}
                  className="font-mono text-sm uppercase tracking-[0.3em] text-cyan-500/60"
                >
                  Homelab Management System v1.0.0
                </motion.p>
              </div>

              <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: 1, duration: 0.5 }}
              >
                <ShimmerButton
                  shimmerColor="#06b6d4"
                  shimmerDuration="2s"
                  background="rgba(6, 182, 212, 0.1)"
                  className="px-12 py-4 font-mono text-lg font-bold tracking-widest border-cyan-500/30 hover:scale-105 transition-transform"
                  onClick={() => setIsEntered(true)}
                >
                  INITIALIZE SYSTEM
                </ShimmerButton>
              </motion.div>
            </div>

            {/* Decorative Elements */}
            <div className="absolute bottom-12 left-12 font-mono text-[10px] text-cyan-500/20 uppercase tracking-widest vertical-rl">
              System Core // Protected Area
            </div>
            <div className="absolute top-12 right-12 font-mono text-[10px] text-cyan-500/20 uppercase tracking-widest">
              Coordinate: 0.0.0.0 // Localhost
            </div>
          </motion.div>
        ) : (
          <motion.div
            key="dashboard"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8, ease: "easeOut" }}
            className="w-full"
          >
            <Dashboard />
          </motion.div>
        )}
      </AnimatePresence>
    </main>
  )
}