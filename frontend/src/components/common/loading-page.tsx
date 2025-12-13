import { Spinner } from "@/components/ui/spinner"
import { cn } from "@/lib/utils"

interface LoadingPageProps {
  className?: string
}

export function LoadingPage({ className }: LoadingPageProps) {
  return (
    <div className={cn("flex h-[50vh] w-full items-center justify-center", className)}>
      <Spinner className="size-10 text-primary" />
    </div>
  )
}
