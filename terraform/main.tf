terraform {
  required_providers {
    grafana = {
      source  = "grafana/grafana"
      version = "~> 2.3.0"
    }
  }
}

provider "grafana" {
  url  = "http://localhost:3000"
  auth = "admin:admin"
}

resource "grafana_folder" "app_metrics" {
  title = "Application Metrics"
}

resource "grafana_dashboard" "requests_dashboard" {
  folder      = grafana_folder.app_metrics.id
  config_json = jsonencode({
    "title" = "API Metrics Dashboard"
    "editable" = true
    "refresh" = "10s"
     "time" = {
      "from" = "now-5m"
      "to" = "now"
    }
    "panels" = [
      {
        "title" = "Request Count by Path"
        "type" = "timeseries"
        "gridPos" = {
          "h" = 8
          "w" = 12
          "x" = 0
          "y" = 0
        }
        "datasource" = {
          "type" = "prometheus",
          "uid" = "prometheus"
        }
        "targets" = [
          {
            "expr" = "sum by(path) (rate(requests_total[5m]))"
            "legendFormat" = "{{path}}"
            "refId" = "A"
          }
        ]
      },
      {
        "title" = "Error Rate"
        "type" = "timeseries"
        "gridPos" = {
          "h" = 8
          "w" = 12
          "x" = 12
          "y" = 0
        }
        "datasource" = {
          "type" = "prometheus",
          "uid" = "prometheus"
        }
        "targets" = [
          {
            "expr" = "sum(rate(requests_errors_total[5m])) / sum(rate(requests_total[5m]))"
            "legendFormat" = "Error Rate"
            "refId" = "A"
          }
        ]
      },
      {
        "title" = "HTTP Status Distribution"
        "type" = "timeseries"
        "gridPos" = {
          "h" = 8
          "w" = 12
          "x" = 0
          "y" = 8
        }
        "datasource" = {
          "type" = "prometheus",
          "uid" = "prometheus"
        }
        "targets" = [
          {
            "expr" = "sum by(status) (rate(requests_total[5m]))"
            "legendFormat" = "{{status}}"
            "refId" = "A"
          }
        ]
      },
      {
        "title" = "Errors by Path"
        "type" = "timeseries"
        "gridPos" = {
          "h" = 8
          "w" = 12
          "x" = 12
          "y" = 8
        }
        "datasource" = {
          "type" = "prometheus",
          "uid" = "prometheus"
        }
        "targets" = [
          {
            "expr" = "sum by(path) (rate(requests_errors_total[5m]))"
            "legendFormat" = "{{path}}"
            "refId" = "A"
          }
        ]
      },
      {
        "title" = "Average Request Duration"
        "type" = "timeseries"
        "gridPos" = {
          "h" = 8
          "w" = 12
          "x" = 0
          "y" = 16
        }
        "datasource" = {
          "type" = "prometheus",
          "uid" = "prometheus"
        }
        "targets" = [
          {
            "expr" = "sum(rate(request_duration_seconds_sum[5m])) by (path) / sum(rate(request_duration_seconds_count[5m])) by (path)"
            "legendFormat" = "{{path}}"
            "refId" = "A"
          }
        ],
        "options" = {
          "tooltip" = {
            "mode" = "single",
            "sort" = "none"
          }
        },
        "fieldConfig" = {
          "defaults" = {
            "color" = {
              "mode" = "palette-classic"
            },
            "custom" = {
              "axisLabel" = "seconds",
              "lineInterpolation" = "linear",
              "fillOpacity" = 10
            },
            "unit" = "s"
          }
        }
      },
      {
        "title" = "95th Percentile Request Duration"
        "type" = "timeseries"
        "gridPos" = {
          "h" = 8
          "w" = 12
          "x" = 12
          "y" = 16
        }
        "datasource" = {
          "type" = "prometheus",
          "uid" = "prometheus"
        }
        "targets" = [
          {
            "expr" = "histogram_quantile(0.95, sum(rate(request_duration_seconds_bucket[5m])) by (le, path))"
            "legendFormat" = "{{path}}"
            "refId" = "A"
          }
        ],
        "options" = {
          "tooltip" = {
            "mode" = "single",
            "sort" = "none"
          }
        },
        "fieldConfig" = {
          "defaults" = {
            "color" = {
              "mode" = "palette-classic"
            },
            "custom" = {
              "axisLabel" = "seconds",
              "lineInterpolation" = "linear",
              "fillOpacity" = 10
            },
            "unit" = "s"
          }
        }
      },
      {
        "title" = "Request Duration Heatmap"
        "type" = "heatmap"
        "gridPos" = {
          "h" = 8
          "w" = 24
          "x" = 0
          "y" = 24
        }
        "datasource" = {
          "type" = "prometheus",
          "uid" = "prometheus"
        }
        "targets" = [
          {
            "expr" = "sum(rate(request_duration_seconds_bucket[1m])) by (le)",
            "format" = "heatmap",
            "legendFormat" = "{{le}}",
            "refId" = "A"
          }
        ],
        "options" = {
          "calculate" = false,
          "calculation" = "",
          "cellGap" = 1,
          "color" = {
            "exponent" = 0.5,
            "fill" = "dark-orange",
            "mode" = "scheme",
            "scale" = "exponential",
            "scheme" = "Oranges"
          },
          "yAxis" = {
            "decimals" = 0,
            "unit" = "s"
          }
        }
      }
    ]
  })
  depends_on = [grafana_data_source.prometheus]
}

resource "grafana_data_source" "prometheus" {
  type          = "prometheus"
  name          = "Prometheus"
  uid           = "prometheus"  // This must match the uid used in your dashboard
  url           = "http://prometheus:9090"
  is_default    = true
}