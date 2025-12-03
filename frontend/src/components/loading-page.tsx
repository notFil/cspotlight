import { Spinner } from "@/components/ui/spinner"

export function LoadingPage() {
  return (
    <div className="flex h-[50vh] w-full items-center justify-center">
      <Spinner className="size-10 text-primary" />
    </div>
  )
}
