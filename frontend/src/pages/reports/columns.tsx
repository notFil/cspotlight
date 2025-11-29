import type { ColumnDef } from "@tanstack/react-table"
import type { CSPReport } from "@/types"
import { ArrowUpDown } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Chrome, Globe, Smartphone } from "lucide-react"

const getBrowserIcon = (userAgent: string) => {
  const ua = userAgent.toLowerCase()
  if (ua.includes("chrome")) return <Chrome className="h-4 w-4" />
  if (ua.includes("safari") && !ua.includes("chrome")) return <Globe className="h-4 w-4" /> // Lucide doesn't have Safari, using Globe
  if (ua.includes("firefox")) return <Globe className="h-4 w-4" /> // Lucide doesn't have Firefox
  if (ua.includes("mobile")) return <Smartphone className="h-4 w-4" />
  return <Globe className="h-4 w-4" />
}

export const columns: ColumnDef<CSPReport>[] = [
  {
    accessorKey: "lastSeen",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Last Seen
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
    cell: ({ row }) => {
      return new Date(row.getValue("lastSeen")).toLocaleString()
    },
  },
  {
    accessorKey: "url",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          URL
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
    filterFn: "includesString",
  },
  {
    accessorKey: "directive",
    header: "Directive",
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: "ipAddress",
    header: "IP Address",
  },
  {
    accessorKey: "userAgent",
    header: "Browser",
    cell: ({ row }) => {
      const userAgent = row.getValue("userAgent") as string
      return (
        <div className="flex items-center gap-2" title={userAgent}>
          {getBrowserIcon(userAgent)}
          <span className="text-muted-foreground text-xs truncate max-w-[150px]">
            {userAgent}
          </span>
        </div>
      )
    },
    filterFn: "includesString",
  },
  {
    accessorKey: "count",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Count
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
  },
  {
    accessorKey: "raw",
    header: "Raw Report",
    cell: ({ row }) => {
      const raw = row.getValue("raw") as string
      let formattedRaw = raw
      try {
        formattedRaw = JSON.stringify(JSON.parse(raw), null, 2)
      } catch (e) {
        // keep original raw if parse fails
      }

      return (
        <Dialog>
          <DialogTrigger asChild>
            <Button variant="outline" size="sm">See raw</Button>
          </DialogTrigger>
          <DialogContent className="max-w-[800px] max-h-[80vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>Raw CSP Report</DialogTitle>
            </DialogHeader>
            <pre className="bg-muted p-4 rounded-md whitespace-pre-wrap break-all text-xs">
              {formattedRaw}
            </pre>
          </DialogContent>
        </Dialog>
      )
    },
  },
]
