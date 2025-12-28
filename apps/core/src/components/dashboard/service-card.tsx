import { motion } from "motion/react";
import { LucideIcon } from "lucide-react";
import { Badge } from "@/components/ui/badge";

interface ServiceCardProps {
  icon: LucideIcon;
  title: string;
  description: string;
  status: "online" | "offline" | "maintenance";
  metrics?: {
    label: string;
    value: string;
  }[];
  delay?: number;
}

export function ServiceCard({ icon: Icon, title, description, status, metrics, delay = 0 }: ServiceCardProps) {
  const statusColors = {
    online: "text-accent border-accent/50 bg-accent/10",
    offline: "text-destructive border-destructive/50 bg-destructive/10",
    maintenance: "text-secondary border-secondary/50 bg-secondary/10",
  };

  const statusLabels = {
    online: "ONLINE",
    offline: "OFFLINE",
    maintenance: "MAINTENANCE",
  };

  return (
    <motion.div
      className="group relative overflow-hidden rounded-lg border border-border bg-card backdrop-blur-md p-6 transition-all duration-300 hover:border-primary/50 hover:shadow-lg hover:shadow-primary/20"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, delay }}
      whileHover={{ y: -5 }}
    >
      {/* Shimmer Effect on Hover */}
      <div className="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-500">
        <div
          className="absolute inset-0 w-1/2 h-full bg-gradient-to-r from-transparent via-primary/10 to-transparent"
          style={{
            animation: 'shimmer 2s ease-in-out infinite',
          }}
        />
      </div>

      <div className="relative">
        {/* Header */}
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-primary/10 border border-primary/30">
              <Icon className="w-5 h-5 text-primary" />
            </div>
            <div>
              <h3 className="text-card-foreground">{title}</h3>
              <p className="text-xs text-muted-foreground mt-1">{description}</p>
            </div>
          </div>
          <Badge className={`${statusColors[status]} border text-xs`}>
            <div className="w-1.5 h-1.5 rounded-full bg-current mr-1.5 animate-pulse" />
            {statusLabels[status]}
          </Badge>
        </div>

        {/* Metrics */}
        {metrics && metrics.length > 0 && (
          <div className="grid grid-cols-2 gap-3 pt-4 border-t border-border/50">
            {metrics.map((metric, idx) => (
              <div key={idx} className="flex flex-col">
                <span className="text-xs text-muted-foreground">{metric.label}</span>
                <span className="text-sm text-card-foreground mt-1">{metric.value}</span>
              </div>
            ))}
          </div>
        )}

        {/* Hover Accent Line */}
        <motion.div
          className="absolute bottom-0 left-0 h-0.5 bg-gradient-to-r from-primary via-accent to-secondary"
          initial={{ width: 0 }}
          whileHover={{ width: "100%" }}
          transition={{ duration: 0.3 }}
        />
      </div>
    </motion.div>
  );
}
