import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/common/app-sidebar"
import ThemeToggle from "@/components/common/theme-toggle"

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider defaultOpen={false}>
      <ThemeToggle />
      <AppSidebar />
      <main className="w-full bg-background">
        <SidebarTrigger className="md:hidden fixed top-4 left-4 z-50" />
        {children}
      </main>
    </SidebarProvider>
  )
}