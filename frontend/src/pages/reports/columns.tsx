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
import { Badge } from "@/components/ui/badge"

const getBrowserIcon = (userAgent: string) => {
  const ua = userAgent.toLowerCase()
  if (ua.includes("chrome")) return <Chrome className="h-4 w-4" />
  if (ua.includes("safari") && !ua.includes("chrome")) return <Globe className="h-4 w-4" />
  if (ua.includes("firefox")) return <Globe className="h-4 w-4" />
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
    accessorKey: "disposition",
    header: "Disposition",
    cell: ({ row }) => {
      const disposition = row.getValue("disposition") as string
      return (
        <Badge
          variant={disposition === "enforce" ? "default" : "secondary"}
          className={disposition === "enforce" ? "bg-green-500 hover:bg-green-600" : ""}
        >
          {disposition}
        </Badge>
      )
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: "documentURL",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Document URL
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
    accessorKey: "sourceIP",
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
    accessorKey: "body",
    header: "Raw Report",
    cell: ({ row }) => {
      const body = row.getValue("body")
      let formattedRaw = ""

      if (typeof body === "object" && body !== null) {
        formattedRaw = JSON.stringify(body, null, 2)
      } else {
        try {
          formattedRaw = JSON.stringify(JSON.parse(body as string), null, 2)
        } catch (e) {
          formattedRaw = body as string
        }
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
