import {
    Sidebar,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
} from "@workspace/ui/components/sidebar"
import { FileText, History, LayoutDashboard, Pencil, Wallet } from "lucide-react"

// User Data
const user = {
    nim: "NIM.123456",
    name: "Jane Doe",
    email: "janedoe@example.com"
}

// Menu items.
const items = [
    {
        title: "Dashboard",
        url: "#",
        icon: LayoutDashboard,
    },
    {
        title: "Pendaftaran",
        url: "#",
        icon: Pencil,
    },
    {
        title: "Keuangan",
        url: "#",
        icon: Wallet,
    },
    {
        title: "Riwayat",
        url: "#",
        icon: History,
    },
    {
        title: "Dokumen",
        url: "#",
        icon: FileText,
    },
]

export function AppSidebar() {
    return (
        <Sidebar>
            <SidebarHeader>
                <h1>UNW TOEFL</h1>
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            {items.map((item) => (
                                <SidebarMenuItem key={item.title}>
                                    <SidebarMenuButton asChild>
                                        <a href={item.url}>
                                            <item.icon />
                                            <span>{item.title}</span>
                                        </a>
                                    </SidebarMenuButton>
                                </SidebarMenuItem>
                            ))}
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
            </SidebarContent>
        </Sidebar>
    )
}