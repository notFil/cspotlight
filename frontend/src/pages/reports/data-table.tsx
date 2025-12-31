import * as React from "react"
import {
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table"
import type {
  ColumnDef,
  SortingState,
} from "@tanstack/react-table"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import type { Pagination as PaginationType, ReportFilters } from "@/types"
import { DIRECTIVES } from "@/constants"
import { Button } from "@/components/ui/button"

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  pagination?: PaginationType
  onPageChange?: (page: number) => void
  onPageSizeChange?: (pageSize: number) => void
  filters?: ReportFilters
  onFilter?: (filters: ReportFilters) => void
}

export function DataTable<TData, TValue>({
  columns,
  data,
  pagination,
  onPageChange,
  onPageSizeChange,
  filters: initialFilters,
  onFilter,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = React.useState<SortingState>([])
  const [localFilters, setLocalFilters] = React.useState<ReportFilters>(initialFilters || {})

  React.useEffect(() => {
    if (initialFilters) {
      setLocalFilters(initialFilters)
    }
  }, [initialFilters])

  const handleFilter = () => {
    onFilter?.(localFilters)
  }

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    onSortingChange: setSorting,
    getSortedRowModel: getSortedRowModel(),
    state: {
      sorting,
    },
    manualPagination: true,
    pageCount: pagination?.totalPages ?? -1,
  })

  return (
    <div>
      <div className="flex items-center py-4 gap-4">
        <Input
          placeholder="Filter Document URL..."
          value={localFilters.documentURL || ""}
          onChange={(event) =>
            setLocalFilters(prev => ({ ...prev, documentURL: event.target.value }))
          }
          className="max-w-sm"
        />
        <Input
          placeholder="Filter User Agent..."
          value={localFilters.userAgent || ""}
          onChange={(event) =>
            setLocalFilters(prev => ({ ...prev, userAgent: event.target.value }))
          }
          className="max-w-sm"
        />
        <Select
          value={localFilters.directive || "all"}
          onValueChange={(value) =>
            setLocalFilters(prev => ({ ...prev, directive: value === "all" ? "" : value }))
          }
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select Directive" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Directives</SelectItem>
            {DIRECTIVES.map((directive) => (
              <SelectItem key={directive} value={directive}>
                {directive}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select
          value={localFilters.disposition || "all"}
          onValueChange={(value) =>
            setLocalFilters(prev => ({ ...prev, disposition: value === "all" ? "" : value }))
          }
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select Disposition" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Dispositions</SelectItem>
            <SelectItem value="enforce">Enforce</SelectItem>
            <SelectItem value="report">Report Only</SelectItem>
          </SelectContent>
        </Select>
        <Button onClick={handleFilter}>Filter</Button>
      </div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                    </TableHead>
                  )
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <div className="py-4 flex items-center justify-between">
        {pagination && onPageSizeChange && (
          <div className="flex items-center space-x-2">
            <p className="text-sm font-medium">Rows</p>
            <Select
              value={`${pagination.pageSize}`}
              onValueChange={(value) => {
                onPageSizeChange(Number(value))
              }}
            >
              <SelectTrigger className="h-8 w-[70px]">
                <SelectValue placeholder={pagination.pageSize} />
              </SelectTrigger>
              <SelectContent side="top">
                {[10, 20, 30, 50].map((pageSize) => (
                  <SelectItem key={pageSize} value={`${pageSize}`}>
                    {pageSize}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}
        {pagination && onPageChange && (
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() => onPageChange(Math.max(1, pagination.page - 1))}
                  className={pagination.page <= 1 ? "pointer-events-none opacity-50" : "cursor-pointer"}
                />
              </PaginationItem>
              {Array.from({ length: pagination.totalPages }, (_, i) => i + 1).map((page) => {
                if (
                  page === 1 ||
                  page === pagination.totalPages ||
                  (page >= pagination.page - 1 && page <= pagination.page + 1)
                ) {
                  return (
                    <PaginationItem key={page}>
                      <PaginationLink
                        isActive={page === pagination.page}
                        onClick={() => onPageChange(page)}
                        className="cursor-pointer"
                      >
                        {page}
                      </PaginationLink>
                    </PaginationItem>
                  )
                }
                if (
                  page === pagination.page - 2 ||
                  page === pagination.page + 2
                ) {
                  return (
                    <PaginationItem key={page}>
                      <PaginationEllipsis />
                    </PaginationItem>
                  )
                }
                return null
              })}

              <PaginationItem>
                <PaginationNext
                  onClick={() => onPageChange(Math.min(pagination.totalPages, pagination.page + 1))}
                  className={pagination.page >= pagination.totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        )}
      </div>
    </div>
  )
}
