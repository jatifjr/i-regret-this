import { Separator } from "@workspace/ui/components/separator"
import { SidebarTrigger } from "@workspace/ui/components/sidebar"

export function SiteHeader({ children }: { children: React.ReactNode }) {
    return (
        <header className="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear">
            <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
                <SidebarTrigger className="-ml-1" />
            </div>
        </header>
    )
}
