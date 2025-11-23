import { useSidebar } from "@/components/ui/sidebar"
import { cn } from "@/lib/utils"

export default function SidebarBackdrop() {
  const { open, setOpen, isMobile } = useSidebar()

  // Only show backdrop on desktop when sidebar is open
  // On mobile, the Sheet component handles the backdrop automatically
  if (isMobile || !open) return null

  return (
    <div
      className={cn(
        "fixed inset-0 z-20 bg-black/20 backdrop-blur-sm transition-all duration-200",
        "data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"
      )}
      onClick={() => setOpen(false)}
    />
  )
}
