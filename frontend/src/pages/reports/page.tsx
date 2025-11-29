import { useParams } from "react-router-dom"
import type { CSPReport } from "@/types"
import { columns } from "./columns"
import { DataTable } from "./data-table"

const MOCK_DATA: CSPReport[] = [
  {
    url: "https://example.com/admin",
    directive: "script-src",
    ipAddress: "192.168.1.1",
    raw: '{"csp-report":{"document-uri":"https://example.com/admin","referrer":"","violated-directive":"script-src","effective-directive":"script-src","original-policy":"default-src \'none\'; script-src \'self\';","blocked-uri":"https://evil.com/malicious.js","status-code":200}}',
    userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    count: 5,
    lastSeen: "2023-11-29T10:00:00Z",
  },
  {
    url: "https://example.com/login",
    directive: "img-src",
    ipAddress: "10.0.0.5",
    raw: '{"csp-report":{"document-uri":"https://example.com/login","referrer":"","violated-directive":"img-src","effective-directive":"img-src","original-policy":"default-src \'none\'; img-src \'self\';","blocked-uri":"http://insecure.com/image.png","status-code":200}}',
    userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
    count: 12,
    lastSeen: "2023-11-29T09:45:00Z",
  },
  {
    url: "https://example.com/dashboard",
    directive: "style-src",
    ipAddress: "172.16.0.2",
    raw: '{"csp-report":{"document-uri":"https://example.com/dashboard","referrer":"","violated-directive":"style-src","effective-directive":"style-src","original-policy":"default-src \'none\'; style-src \'self\';","blocked-uri":"inline","status-code":200}}',
    userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
    count: 3,
    lastSeen: "2023-11-28T15:30:00Z",
  },
  {
    url: "https://example.com/api/data",
    directive: "connect-src",
    ipAddress: "192.168.1.100",
    raw: '{"csp-report":{"document-uri":"https://example.com/api/data","referrer":"","violated-directive":"connect-src","effective-directive":"connect-src","original-policy":"default-src \'none\'; connect-src \'self\';","blocked-uri":"https://analytics.tracker.com","status-code":200}}',
    userAgent: "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36",
    count: 8,
    lastSeen: "2023-11-28T12:00:00Z",
  },
  {
    url: "https://example.com/",
    directive: "frame-src",
    ipAddress: "10.10.10.10",
    raw: '{"csp-report":{"document-uri":"https://example.com/","referrer":"","violated-directive":"frame-src","effective-directive":"frame-src","original-policy":"default-src \'none\'; frame-src \'self\';","blocked-uri":"https://ads.adserver.com","status-code":200}}',
    userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",
    count: 1,
    lastSeen: "2023-11-27T08:15:00Z",
  },
]

export default function ReportsPage() {
  const { projectId } = useParams()

  return (
    <div className="container mx-auto py-10">
      <div className="flex flex-col gap-4 mb-8">
        <h1 className="text-3xl font-bold tracking-tight">CSP Reports</h1>
        <p className="text-muted-foreground">
          Viewing reports for Project ID: {projectId}
        </p>
      </div>
      <DataTable columns={columns} data={MOCK_DATA} />
    </div>
  )
}
