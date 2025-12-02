import { useParams } from "react-router-dom"
import { columns } from "./columns"
import { DataTable } from "./data-table"
import { useReports } from "@/hooks/use-reports"

export default function ReportsPage() {
  const { projectId } = useParams()
  const { data: reports, isLoading, error } = useReports(projectId || "")

  if (isLoading) {
    return <div className="container mx-auto py-10">Loading reports...</div>
  }

  if (error) {
    return <div className="container mx-auto py-10 text-red-500">Error loading reports</div>
  }

  return (
    <div className="container mx-auto py-10">
      <div className="flex flex-col gap-4 mb-8">
        <h1 className="text-3xl font-bold tracking-tight">CSP Reports</h1>
        <p className="text-muted-foreground">
          Viewing reports for Report ID: {projectId}
        </p>
      </div>
      <DataTable columns={columns} data={reports || []} />
    </div>
  )
}
