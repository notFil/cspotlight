import { ClipboardList, Folder, Home } from "lucide-react"
import { Link } from "react-router-dom"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { useAuth } from "@/hooks/use-auth"
import { NavUser } from "./nav-user"

// Menu items.
const items = [
  {
    title: "Dashboard",
    url: "/",
    icon: Home,
  },
  {
    title: "Projects",
    url: "/projects",
    icon: Folder,
  },
  {
    title: "Management",
    url: "/management",
    icon: ClipboardList,
  },
]


export function AppSidebar() {
  const { user } = useAuth()

  return (
    <Sidebar collapsible="icon" side="left" className="z-40">
      <SidebarHeader className="flex flex-row items-center gap-2">
        <SidebarTrigger className="hidden md:flex size-8" />
        <h2 className="text-lg font-semibold w-full group-data-[collapsible=icon]:hidden" style={{ color: "var(--color-text)" }}>
          cspotlight
        </h2>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild tooltip={item.title} className="h-12 text-base hover:bg-primary/10 hover:text-primary transition-colors">
                    <Link to={item.url}>
                      <item.icon className="h-5 w-5" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        {user && <NavUser user={user} />}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}