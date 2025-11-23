import { ClipboardList, Folder, Home } from "lucide-react"

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
} from "@/components/ui/sidebar"
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

const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  }
}

export function AppSidebar() {
  return (
    <Sidebar collapsible="offcanvas" side="left" className="z-40">
      <SidebarHeader className="sidebar-brand-header">
        <h2 className="text-lg font-semibold w-full" style={{ color: "var(--color-text)" }}>
          cspotlight
        </h2>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild className="h-12 text-base hover:bg-primary/10 hover:text-primary transition-colors">
                    <a href={item.url}>
                      <item.icon className="h-5 w-5" />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter className="p-4">
        {/* <div className="flex items-center gap-3 p-2 rounded-lg bg-sidebar-accent/10">
          <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-semibold">
            SE
          </div>
          <div className="flex flex-col">
            <span className="text-sm font-medium">Security Engineer</span>
            <span className="text-xs text-muted-foreground">Admin</span>
          </div>
        </div> */}
        <NavUser user={data.user} />
      </SidebarFooter>
    </Sidebar>
  )
}