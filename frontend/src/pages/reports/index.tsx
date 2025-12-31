import { useParams, useSearchParams } from "react-router-dom"
import { columns } from "./columns"
import { DataTable } from "./data-table"
import { useReports } from "@/hooks/use-reports"
import { LoadingPage } from "@/components/common/loading-page"
import { ErrorPage } from "@/components/common/error-page"
import type { ReportFilters } from "@/types"

export default function Reports() {
  const { projectId } = useParams()
  const [searchParams, setSearchParams] = useSearchParams()

  const page = Number(searchParams.get("page")) || 1
  const pageSize = Number(searchParams.get("pageSize")) || 10

  const { data: response, isLoading, error } = useReports(projectId || "", page, pageSize, {
    directive: searchParams.get("directive") || undefined,
    disposition: searchParams.get("disposition") || undefined,
    blockedURL: searchParams.get("blocked_url") || undefined,
    userAgent: searchParams.get("user_agent") || undefined,
    documentURL: searchParams.get("document_url") || undefined,
  })

  const reports = response?.data || []
  const pagination = response?.pagination

  const handlePageChange = (newPage: number) => {
    setSearchParams(prev => {
      prev.set("page", newPage.toString())
      return prev
    })
  }

  const handlePageSizeChange = (newPageSize: number) => {
    setSearchParams(prev => {
      prev.set("pageSize", newPageSize.toString())
      prev.set("page", "1")
      return prev
    })
  }

  const handleFilter = (newFilters: ReportFilters) => {
    setSearchParams(prev => {
      prev.set("page", "1")
      if (newFilters.directive) prev.set("directive", newFilters.directive); else prev.delete("directive")
      if (newFilters.disposition) prev.set("disposition", newFilters.disposition); else prev.delete("disposition")
      if (newFilters.blockedURL) prev.set("blocked_url", newFilters.blockedURL); else prev.delete("blocked_url")
      if (newFilters.userAgent) prev.set("user_agent", newFilters.userAgent); else prev.delete("user_agent")
      if (newFilters.documentURL) prev.set("document_url", newFilters.documentURL); else prev.delete("document_url")
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
            Viewing reports for Project ID: {projectId}
          </p>
        </div>
        <DataTable
          columns={columns}
          data={reports}
          pagination={pagination}
          onPageChange={handlePageChange}
          onPageSizeChange={handlePageSizeChange}
          filters={{
            directive: searchParams.get("directive") || undefined,
            disposition: searchParams.get("disposition") || undefined,
            blockedURL: searchParams.get("blocked_url") || undefined,
            userAgent: searchParams.get("user_agent") || undefined,
            documentURL: searchParams.get("document_url") || undefined,
          }}
          onFilter={handleFilter}
        />
      </div>
    </div>
  )
}
