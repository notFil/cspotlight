import { useParams, useSearchParams } from "react-router-dom"
import { columns } from "./columns"
import { DataTable } from "./data-table"
import { useReports } from "@/hooks/use-reports"
import { LoadingPage } from "@/components/common/loading-page"
import { ErrorPage } from "@/components/common/error-page"

export default function Reports() {
  const { projectId } = useParams()
  const [searchParams, setSearchParams] = useSearchParams()

  const page = Number(searchParams.get("page")) || 1
  const pageSize = Number(searchParams.get("pageSize")) || 10

  const { data: response, isLoading, error } = useReports(projectId || "", page, pageSize)

  const reports = response?.data || []
  const pagination = response?.pagination

  const handlePageChange = (newPage: number) => {
    setSearchParams(prev => {
      prev.set("page", newPage.toString())
      return prev
    })
  }

  if (isLoading) {
    return <LoadingPage />
  }

  if (error) {
    return <ErrorPage error={error} />
  }

  return (
    <div className="p-6 min-h-screen" style={{ backgroundColor: 'var(--color-background)' }}>
      <div className="max-w-7xl mx-auto">
        <div className="flex flex-col gap-4 mb-8">
          <h1 className="text-3xl font-bold tracking-tight">CSP Reports</h1>
          <p className="text-muted-foreground">
            Viewing reports for Report ID: {projectId}
          </p>
        </div>
        <DataTable
          columns={columns}
          data={reports}
          pagination={pagination}
          onPageChange={handlePageChange}
        />
      </div>
    </div>
  )
}
