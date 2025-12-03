import { AlertCircle } from "lucide-react"

interface ErrorPageProps {
  error?: Error | null
  message?: string
}

export function ErrorPage({ error, message }: ErrorPageProps) {
  const errorMessage = message || error?.message || "An unknown error occurred"

  return (
    <div className="flex h-[50vh] w-full flex-col items-center justify-center gap-2 text-destructive">
      <AlertCircle className="size-10" />
      <p className="text-lg font-medium">Error loading data</p>
      <p className="text-sm text-muted-foreground">{errorMessage}</p>
    </div>
  )
}
