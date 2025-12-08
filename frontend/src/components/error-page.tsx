import { AlertCircle } from "lucide-react"
import { cn } from "@/lib/utils"

interface ErrorPageProps {
  error?: Error | null
  message?: string
  className?: string
}

export function ErrorPage({ error, message, className }: ErrorPageProps) {
  const errorMessage = message || error?.message || "An unknown error occurred"

  return (
    <div className={cn("flex h-[50vh] w-full flex-col items-center justify-center gap-2 text-destructive", className)}>
      <AlertCircle className="size-10" />
      <p className="text-lg font-medium">Error loading data</p>
      <p className="text-sm text-muted-foreground">{errorMessage}</p>
    </div>
  )
}
